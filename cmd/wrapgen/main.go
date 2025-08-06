// Package main implements the wrapgen code generator for s-expression wrapper functions.
//
// Wrapgen analyzes Go functions using reflection and generates type-safe wrapper
// functions that convert between Go types and s-expression types (core.Expr).
//
// # Validation Modes
//
// The system supports three validation modes controlled by the -validation flag:
//
//   - trust (production): Direct type assertions with no validation overhead.
//     Assumes input types are correct and performs direct casts like:
//     arg.(core.Integer) or arg.(core.Real)
//     This mode is fastest but will panic on type mismatches.
//
//   - debug (development): Full validation with panic on type mismatch.
//     Uses core.ExtractInt64(), core.ExtractFloat64(), etc. and panics
//     with detailed error messages on validation failures. Useful for
//     development and debugging.
//
//   - graceful (fallback): Full validation with graceful error handling.
//     Uses validation functions but returns unevaluated expressions
//     (core.CopyExprList) instead of panicking on type mismatches.
//     Allows the system to continue processing invalid inputs.
//
// # Code Generation
//
// For each symbol (e.g., Plus, Times), wrapgen generates:
//   - Individual wrapper file (e.g., wrapped/plus.go)
//   - Pattern-based dispatch registration in builtin_setup.go
//   - Type-safe parameter conversion and return value handling
//   - Validation code appropriate to the selected mode
//
// # Performance Optimizations
//
//   - Expr parameters: When Go functions accept core.Expr parameters,
//     no conversion is needed - args are passed directly
//   - Variadic optimization: Detects when conversions are unnecessary
//   - Trust mode: Eliminates all validation overhead for production builds
package main

import (
	"cmp"
	"flag"
	"fmt"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"text/template"
)

// FunctionInfo contains expanded information about a symbol's function
type FunctionInfo struct {
	SymbolName     string // "Plus"
	Pattern        string // "Plus(x__Integer)"
	FunctionName   string // "PlusIntegers"
	WrapperName    string // "WrapPlusIntegers"
	IsVariadic     bool
	ParamType      string   // For variadic functions
	ParamTypes     []string // For fixed-arity functions
	ReturnType     string
	ReturnsError   bool
	NeedsEvaluator bool   // true if function requires *Evaluator (EngineFunc)
	DirectCall     bool   // true if in same form as engine, no need for wrapper
	PackageName    string // "stdlib" or "builtins"
}

func main() {
	var (
		outputDir      = flag.String("dir", "wrapped", "Output directory for generated files")
		setupFile      = flag.String("setup", "builtin_setup.go", "Generate builtin setup file")
		validationMode = flag.String("validation", "trust", "Validation mode: trust (production), debug (panic), graceful (fallback)")
	)
	flag.Parse()

	// Create output directory
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		log.Fatalf("Error creating output directory %s: %v", *outputDir, err)
	}

	slices.SortFunc(symbolSpecs, func(a, b SymbolSpec) int {
		return cmp.Compare(a.Name, b.Name)
	})

	for i, s := range symbolSpecs {
		if !slices.Contains(s.Attributes, "Protected") {
			s.Attributes = append(s.Attributes, "Protected")
		}
		slices.SortFunc(s.Attributes, func(a, b string) int {
			return cmp.Compare(a, b)
		})
		symbolSpecs[i].Attributes = s.Attributes
	}

	// Process all symbols and generate function info
	allFunctions, err := processSymbolSpecs(symbolSpecs)
	if err != nil {
		log.Fatalf("Error processing symbol specs: %v", err)
	}

	// Generate one file per symbol
	totalFunctions := 0

	for _, symbol := range symbolSpecs {
		symbolName := symbol.Name
		if len(symbol.Functions) == 0 {
			continue // Skip constants-only symbols
		}

		// Get functions for this symbol
		var symbolFunctions []FunctionInfo
		for _, fn := range allFunctions {
			if fn.SymbolName == symbolName && !fn.DirectCall {
				symbolFunctions = append(symbolFunctions, fn)
			}
		}

		// Symbol has no functions that needed wrapping
		if len(symbolFunctions) == 0 {
			continue
		}

		// Generate file for this symbol
		filename := strings.ToLower(symbolName) + ".go"
		outputPath := filepath.Join(*outputDir, filename)

		err := generateSymbolFile(outputPath, symbolName, symbolFunctions, *validationMode)
		if err != nil {
			log.Fatalf("Error generating %s: %v", outputPath, err)
		}

		totalFunctions += len(symbolFunctions)
	}

	// Generate builtin_setup.go file in root directory for runtime registration
	setupPath := *setupFile
	err = generateBuiltinSetupFile(setupPath, symbolSpecs, allFunctions)
	if err != nil {
		log.Fatalf("Error generating setup file: %v", err)
	}

	fmt.Printf("Generated %d wrappers across %d symbols in %s/\n", totalFunctions, len(symbolSpecs), *outputDir)
	fmt.Printf("Generated builtin setup file: %s\n", *setupFile)
}

// processSymbolSpecs converts symbol specs to function info using reflection
func processSymbolSpecs(specs []SymbolSpec) ([]FunctionInfo, error) {
	var allFunctions []FunctionInfo

	for _, symbol := range specs {
		symbolName := symbol.Name
		for _, r := range symbol.Functions {
			patternSuffix := r.Pattern
			fullPattern := symbolName + patternSuffix

			// Create function spec for reflection analysis
			funcSpec := FunctionSpec{
				Pattern:  fullPattern,
				Function: r.Function,
			}

			// Analyze with reflection
			err := funcSpec.fillFromReflection()
			if err != nil {
				return nil, fmt.Errorf("error analyzing %s: %v", fullPattern, err)
			}

			// Determine package name based on function
			packageName := "stdlib" // default
			if funcSpec.NeedsEvaluator {
				packageName = "builtins"
			}

			// Convert to FunctionInfo
			funcInfo := FunctionInfo{
				SymbolName:     symbolName,
				Pattern:        fullPattern,
				FunctionName:   funcSpec.FunctionName,
				WrapperName:    funcSpec.WrapperName,
				IsVariadic:     funcSpec.IsVariadic,
				ParamType:      funcSpec.ParamType,
				ParamTypes:     funcSpec.ParamTypes,
				ReturnType:     funcSpec.ReturnType,
				ReturnsError:   funcSpec.ReturnsError,
				NeedsEvaluator: funcSpec.NeedsEvaluator,
				DirectCall:     funcSpec.DirectCall,
				PackageName:    packageName,
			}

			allFunctions = append(allFunctions, funcInfo)
		}
	}

	return allFunctions, nil
}

// generateSymbolFile generates a wrapper file for a single symbol
func generateSymbolFile(outputPath, symbolName string, functions []FunctionInfo, validationMode string) error {
	// Check if any function needs evaluator and what packages are needed
	hasEngineFunc := false
	hasStdlibFunc := false
	for _, fn := range functions {
		if fn.NeedsEvaluator {
			hasEngineFunc = true
		} else {
			hasStdlibFunc = true
		}
	}

	// Template data
	data := struct {
		SymbolName     string
		Functions      []FunctionInfo
		ValidationMode string
		HasEngineFunc  bool
		HasStdlibFunc  bool
	}{
		SymbolName:     symbolName,
		Functions:      functions,
		ValidationMode: validationMode,
		HasEngineFunc:  hasEngineFunc,
		HasStdlibFunc:  hasStdlibFunc,
	}

	// Use same template functions as before
	funcMap := getTemplateFunctions()

	// Create and execute template
	t, err := template.New("symbol").Funcs(funcMap).Parse(wrapperTemplate)
	if err != nil {
		return err
	}

	// Generate code
	var buf strings.Builder
	err = t.Execute(&buf, data)
	if err != nil {
		return err
	}

	// Format the generated code
	formatted, err := format.Source([]byte(buf.String()))
	if err != nil {
		// If formatting fails, write unformatted code
		formatted = []byte(buf.String())
	}

	// Write to file
	return os.WriteFile(outputPath, formatted, 0644)
}

// generateBuiltinSetupFile generates the builtin_setup.go file
func generateBuiltinSetupFile(outputFile string, symbols []SymbolSpec, functions []FunctionInfo) error {
	tmpl := `// Code generated by wrapgen. DO NOT EDIT.

package sexpr

import (
	"fmt"
	"github.com/client9/sexpr/engine"
	"github.com/client9/sexpr/builtins"
	"github.com/client9/sexpr/wrapped"
)

// SetupBuiltinAttributes sets up standard attributes for built-in functions
func SetupBuiltinAttributes(symbolTable *engine.SymbolTable) {
	// Reset attributes
	symbolTable.Reset()

{{range $symbol := .Symbols}}
	symbolTable.SetAttributes("{{$symbol.Name}}", []engine.Attribute{ 
{{- range $i, $attr := $symbol.Attributes -}}
{{- if $i}}, {{ end }}engine.{{$attr}}
{{- end -}}
})
{{ end }}
}

// RegisterDefaultBuiltins registers all built-in functions with their patterns
func RegisterDefaultBuiltins(registry *engine.FunctionRegistry) {
	// Register built-in functions with pattern-based dispatch
	builtinPatterns := map[string]engine.PatternFunc{
		// Generated pattern registrations
{{range .Functions}}
{{- if .DirectCall }}
          "{{.Pattern}}": builtins.{{.FunctionName}},
{{- else}}
	  "{{.Pattern}}": wrapped.{{.WrapperName}},
{{- end}}
{{end}}
	}

	// Register patterns with the function registry
	err := registry.RegisterPatternBuiltins(builtinPatterns)
	if err != nil {
		panic(fmt.Sprintf("Failed to register builtin patterns: %v", err))
	}
}
`

	// Template data
	data := struct {
		Symbols   []SymbolSpec
		Functions []FunctionInfo
	}{
		Symbols:   symbols,
		Functions: functions,
	}

	// Create and execute template
	t, err := template.New("setup").Parse(tmpl)
	if err != nil {
		return err
	}

	// Generate code
	var buf strings.Builder
	err = t.Execute(&buf, data)
	if err != nil {
		return err
	}

	// Format the generated code
	formatted, err := format.Source([]byte(buf.String()))
	if err != nil {
		// If formatting fails, write unformatted code
		formatted = []byte(buf.String())
	}

	// Write to file
	return os.WriteFile(outputFile, formatted, 0644)
}

// generateWrappedTypesFile function removed - no longer needed since ctx parameter eliminated

// getConversionWithMode generates type conversion code based on validation mode
func getConversionWithMode(paramType, validationMode, symbolName string) string {
	switch validationMode {
	case "trust":
		// Trust mode - direct type assertion without checks
		switch paramType {
		case "int64":
			return "convertedArgs[i] = int64(arg.(core.Integer))"
		case "float64":
			return "convertedArgs[i] = float64(arg.(core.Real))"
		case "Number":
			return "convertedArgs[i] = arg // Pass Number through as-is"
		case "Expr":
			return "convertedArgs[i] = arg"
		default:
			return "convertedArgs[i] = arg"
		}
	case "debug":
		// Debug mode - panic on type mismatch
		switch paramType {
		case "int64":
			return `if val, ok := core.ExtractInt64(arg); ok {
	convertedArgs[i] = val
} else {
	panic(fmt.Sprintf("Type mismatch in ` + symbolName + `: expected Integer, got %T", arg))
}`
		case "float64":
			return `if val, ok := core.ExtractFloat64(arg); ok {
	convertedArgs[i] = val
} else {
	panic(fmt.Sprintf("Type mismatch in ` + symbolName + `: expected Real, got %T", arg))
}`
		case "Number":
			return `if val, ok := stdlib.ExtractNumber(arg); ok {
	convertedArgs[i] = val
} else {
	panic(fmt.Sprintf("Type mismatch in ` + symbolName + `: expected Number, got %T", arg))
}`
		case "Expr":
			return "convertedArgs[i] = arg"
		default:
			return "convertedArgs[i] = arg"
		}
	default: // graceful mode
		switch paramType {
		case "int64":
			return `if val, ok := core.ExtractInt64(arg); ok {
	convertedArgs[i] = val
} else {
	// Type mismatch - return unchanged
	return core.CopyExprList(funcName, args)
}`
		case "float64":
			return `if val, ok := core.ExtractFloat64(arg); ok {
	convertedArgs[i] = val
} else {
	// Type mismatch - return unchanged
	return core.CopyExprList(funcName, args)
}`
		case "Number":
			return `if val, ok := stdlib.ExtractNumber(arg); ok {
	convertedArgs[i] = val
} else {
	// Type mismatch - return unchanged
	return core.CopyExprList(funcName, args)
}`
		case "Expr":
			return "convertedArgs[i] = arg"
		default:
			return "convertedArgs[i] = arg"
		}
	}
}

// getFixedConversionWithMode generates fixed-parameter conversion code based on validation mode
func getFixedConversionWithMode(paramTypes []string, validationMode, symbolName string) string {
	var conversions []string
	for i, paramType := range paramTypes {
		varName := fmt.Sprintf("arg%d", i)

		if validationMode == "trust" {
			// Trust mode - direct type assertions
			switch paramType {
			case "ByteArray":
				conversions = append(conversions, fmt.Sprintf("\t%s := args[%d].(core.ByteArray)", varName, i))
			case "Number":
				conversions = append(conversions, fmt.Sprintf("\t%s, _ := stdlib.ExtractNumber(args[%d]) // Trust Number extraction", varName, i))
			case "Association":
				conversions = append(conversions, fmt.Sprintf("\t%s := args[%d].(core.Association)", varName, i))
			case "Expr":
				conversions = append(conversions, fmt.Sprintf("\t%s := args[%d]", varName, i))
			case "int64":
				conversions = append(conversions, fmt.Sprintf("\t%s := int64(args[%d].(core.Integer))", varName, i))
			case "float64":
				conversions = append(conversions, fmt.Sprintf("\t%s := float64(args[%d].(core.Real))", varName, i))
			case "string":
				conversions = append(conversions, fmt.Sprintf("\t%s := string(args[%d].(core.String))", varName, i))
			case "bool":
				conversions = append(conversions, fmt.Sprintf("\t%s, _ := core.ExtractBool(args[%d]) // Trust bool extraction", varName, i))
			case "List":
				conversions = append(conversions, fmt.Sprintf("\t%s := args[%d].(core.List)", varName, i))
			case "ObjectExpr":
				conversions = append(conversions, fmt.Sprintf("\t%s := args[%d].(core.ObjectExpr)", varName, i))
			case "*engine.Evaluator":
				// This should not happen in wrappers - evaluator is passed separately
				log.Fatalf("*engine.Evaluator should not be a conversion parameter - use EngineFunc template instead")
			case "Evaluator":
				// This should not happen in wrappers - evaluator is passed separately
				log.Fatalf("Evaluator interface should not be a conversion parameter - use EngineFunc template instead")
			default:
				log.Fatalf("Unknown Parameter Type: %s", paramType)
			}
		} else {
			// Debug or graceful mode - with validation
			fallbackAction := fmt.Sprintf("\t\treturn core.CopyExprList(\"%s\", args)", symbolName)
			if validationMode == "debug" {
				fallbackAction = fmt.Sprintf("\t\tpanic(fmt.Sprintf(\"Type mismatch in %s: expected %s, got %%T\", args[%d]))", symbolName, paramType, i)
			}

			switch paramType {
			case "ByteArray":
				conversions = append(conversions, fmt.Sprintf("\t%s, ok := core.ByteArray(args[%d])", varName, i))
				conversions = append(conversions, "\tif !ok {")
				conversions = append(conversions, fallbackAction)
				conversions = append(conversions, "\t}")
			case "Number":
				conversions = append(conversions, fmt.Sprintf("\t%s, ok := stdlib.ExtractNumber(args[%d])", varName, i))
				conversions = append(conversions, "\tif !ok {")
				conversions = append(conversions, fallbackAction)
				conversions = append(conversions, "\t}")
			case "Association":
				conversions = append(conversions, fmt.Sprintf("\t%s, ok := core.ExtractAssociation(args[%d])", varName, i))
				conversions = append(conversions, "\tif !ok {")
				conversions = append(conversions, fallbackAction)
				conversions = append(conversions, "\t}")
			case "Expr":
				conversions = append(conversions, fmt.Sprintf("\t%s := args[%d]", varName, i))
			case "int64":
				conversions = append(conversions, fmt.Sprintf("\t%s, ok := core.ExtractInt64(args[%d])", varName, i))
				conversions = append(conversions, "\tif !ok {")
				conversions = append(conversions, fallbackAction)
				conversions = append(conversions, "\t}")
			case "float64":
				conversions = append(conversions, fmt.Sprintf("\t%s, ok := core.ExtractFloat64(args[%d])", varName, i))
				conversions = append(conversions, "\tif !ok {")
				conversions = append(conversions, fallbackAction)
				conversions = append(conversions, "\t}")
			case "string":
				conversions = append(conversions, fmt.Sprintf("\t%s, ok := core.ExtractString(args[%d])", varName, i))
				conversions = append(conversions, "\tif !ok {")
				conversions = append(conversions, fallbackAction)
				conversions = append(conversions, "\t}")
			case "bool":
				conversions = append(conversions, fmt.Sprintf("\t%s, ok := core.ExtractBool(args[%d])", varName, i))
				conversions = append(conversions, "\tif !ok {")
				conversions = append(conversions, fallbackAction)
				conversions = append(conversions, "\t}")
			case "List":
				conversions = append(conversions, fmt.Sprintf("\t%s, ok := args[%d].(core.List)", varName, i))
				conversions = append(conversions, "\tif !ok {")
				conversions = append(conversions, fallbackAction)
				conversions = append(conversions, "\t}")
			case "ObjectExpr":
				conversions = append(conversions, fmt.Sprintf("\t%s, ok := args[%d].(core.ObjectExpr)", varName, i))
				conversions = append(conversions, "\tif !ok {")
				conversions = append(conversions, fallbackAction)
				conversions = append(conversions, "\t}")
			case "*engine.Evaluator":
				// This should not happen in wrappers - evaluator is passed separately
				log.Fatalf("*engine.Evaluator should not be a conversion parameter - use EngineFunc template instead")
			case "Evaluator":
				// This should not happen in wrappers - evaluator is passed separately
				log.Fatalf("Evaluator interface should not be a conversion parameter - use EngineFunc template instead")
			default:
				log.Fatalf("Unknown Parameter Type: %s", paramType)
			}
		}
	}
	return strings.Join(conversions, "\n")
}

// getTemplateFunctions returns the template function map (same as before)
func getTemplateFunctions() template.FuncMap {
	return template.FuncMap{
		"extractFuncName": func(pattern string) string {
			for i, c := range pattern {
				if c == '(' {
					return pattern[:i]
				}
			}
			return pattern
		},
		"getConversion": func(paramType string) string {
			// Legacy function - replaced by getConversionWithMode
			return getConversionWithMode(paramType, "graceful", "Unknown")
		},
		"getConversionWithMode": func(paramType string, validationMode string, symbolName string) string {
			return getConversionWithMode(paramType, validationMode, symbolName)
		},
		"getFixedConversion": func(paramTypes []string) string {
			// Legacy function - replaced by getFixedConversionWithMode
			return getFixedConversionWithMode(paramTypes, "graceful", "Unknown")
		},
		"getFixedConversionWithMode": func(paramTypes []string, validationMode string, symbolName string) string {
			return getFixedConversionWithMode(paramTypes, validationMode, symbolName)
		},
		"getReturnConversion": func(returnType string) string {
			switch returnType {
			case "int64":
				return "return core.NewInteger(result)"
			case "float64":
				return "return core.NewReal(result)"
			case "string":
				return "return core.NewString(result)"
			case "bool":
				return "return core.NewBool(result)"
			case "Expr":
				return "return result"
			default:
				return "return result"
			}
		},
		"getCallArgs": func(paramTypes []string) string {
			var args []string
			for i := range paramTypes {
				args = append(args, fmt.Sprintf("arg%d", i))
			}
			return strings.Join(args, ", ")
		},
		"raw": func(s string) string {
			return s
		},
	}
}

var wrapperTemplate = `// Code generated by wrapgen; DO NOT EDIT.
// Symbol: {{.SymbolName}}
// Validation mode: {{.ValidationMode}}

package wrapped

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
{{- if .HasStdlibFunc}}
	"github.com/client9/sexpr/stdlib"
{{- else }}
	"github.com/client9/sexpr/builtins"
{{- end }}
{{- if eq .ValidationMode "debug"}}
	"fmt"
{{- end}}
)

{{range .Functions}}
// {{.WrapperName}} wraps {{.FunctionName}} for the pattern system
// Generated from pattern: {{.Pattern}}
func {{.WrapperName}}(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
{{- if .IsVariadic}}
	{{- if eq .ParamType "Expr"}}
	// No conversion needed - pass args directly
	{{- else}}
	{{- if ne $.ValidationMode "trust"}}
	funcName := "{{.Pattern | extractFuncName}}"
	{{- end}}
	
	// Convert all args to {{.ParamType}}
	convertedArgs := make([]{{.ParamType}}, len(args))
	for i, arg := range args {
		{{getConversionWithMode .ParamType $.ValidationMode $.SymbolName}}
	}
	{{- end}}
	
	// Call business logic function
	{{- if .ReturnsError}}
	{{- if .NeedsEvaluator}}
	{{- if eq .ParamType "Expr"}}
	result, err := {{.PackageName}}.{{.FunctionName}}(e, c, args...)
	{{- else}}
	result, err := {{.PackageName}}.{{.FunctionName}}(e, c, convertedArgs...)
	{{- end}}
	{{- else}}
	{{- if eq .ParamType "Expr"}}
	result, err := {{.PackageName}}.{{.FunctionName}}(args...)
	{{- else}}
	result, err := {{.PackageName}}.{{.FunctionName}}(convertedArgs...)
	{{- end}}
	{{- end}}
	if err != nil {
		return core.NewError(err.Error(), err.Error())
	}
	{{- else}}
	{{- if .NeedsEvaluator}}
	{{- if eq .ParamType "Expr"}}
	result := {{.PackageName}}.{{.FunctionName}}(e, c, args...)
	{{- else}}
	result := {{.PackageName}}.{{.FunctionName}}(e, c, convertedArgs...)
	{{- end}}
	{{- else}}
	{{- if eq .ParamType "Expr"}}
	result := {{.PackageName}}.{{.FunctionName}}(args...)
	{{- else}}
	result := {{.PackageName}}.{{.FunctionName}}(convertedArgs...)
	{{- end}}
	{{- end}}
	{{- end}}
	
	// Convert result back to Expr
	{{.ReturnType | getReturnConversion}}
{{- else}}
	{{- if eq $.ValidationMode "trust"}}
	// Trust mode - no validation, direct conversion
	{{- else}}
	// Validate argument count
	if len(args) != {{len .ParamTypes}} {
		{{- if eq $.ValidationMode "debug"}}
		panic(fmt.Sprintf("{{.FunctionName}} expects {{len .ParamTypes}} arguments, got %d", len(args)))
		{{- else}}
		return core.NewError("ArgumentError", 
			"{{.FunctionName}} expects {{len .ParamTypes}} arguments")
		{{- end}}
	}
	{{- end}}
	
{{getFixedConversionWithMode .ParamTypes $.ValidationMode $.SymbolName | raw}}
	
	// Call business logic function
	{{- if .ReturnsError}}
	{{- if .NeedsEvaluator}}
	result, err := {{.PackageName}}.{{.FunctionName}}(e, c, {{.ParamTypes | getCallArgs}})
	{{- else}}
	result, err := {{.PackageName}}.{{.FunctionName}}({{.ParamTypes | getCallArgs}})
	{{- end}}
	if err != nil {
		return core.NewError(err.Error(), err.Error())
	}
	{{- else}}
	{{- if .NeedsEvaluator}}
	result := {{.PackageName}}.{{.FunctionName}}(e, c, {{.ParamTypes | getCallArgs}})
	{{- else}}
	result := {{.PackageName}}.{{.FunctionName}}({{.ParamTypes | getCallArgs}})
	{{- end}}
	{{- end}}
	
	// Convert result back to Expr
	{{.ReturnType | getReturnConversion}}
{{- end}}
}

{{end}}`
