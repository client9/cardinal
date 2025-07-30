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
	"flag"
	"fmt"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// FunctionInfo contains expanded information about a symbol's function
type FunctionInfo struct {
	SymbolName   string // "Plus"
	Pattern      string // "Plus(x__Integer)"
	FunctionName string // "PlusIntegers"
	WrapperName  string // "WrapPlusIntegers"
	IsVariadic   bool
	ParamType    string   // For variadic functions
	ParamTypes   []string // For fixed-arity functions
	ReturnType   string
	ReturnsError bool
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

	// Process all symbols and generate function info
	allFunctions, err := processSymbolSpecs(symbolSpecs)
	if err != nil {
		log.Fatalf("Error processing symbol specs: %v", err)
	}

	// Generate one file per symbol
	totalFunctions := 0
	for symbolName, symbol := range symbolSpecs {
		if len(symbol.Functions) == 0 {
			continue // Skip constants-only symbols
		}

		// Get functions for this symbol
		var symbolFunctions []FunctionInfo
		for _, fn := range allFunctions {
			if fn.SymbolName == symbolName {
				symbolFunctions = append(symbolFunctions, fn)
			}
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

	// Always generate builtin_setup.go file
	err = generateBuiltinSetupFile(*setupFile, symbolSpecs, allFunctions)
	if err != nil {
		log.Fatalf("Error generating setup file: %v", err)
	}

	// No longer need types.go file since we removed context parameter

	fmt.Printf("Generated %d wrappers across %d symbols in %s/\n", totalFunctions, len(symbolSpecs), *outputDir)
	fmt.Printf("Generated builtin setup file: %s\n", *setupFile)
}

// processSymbolSpecs converts symbol specs to function info using reflection
func processSymbolSpecs(specs map[string]SymbolSpec) ([]FunctionInfo, error) {
	var allFunctions []FunctionInfo

	for symbolName, symbol := range specs {
		for patternSuffix, function := range symbol.Functions {
			fullPattern := symbolName + patternSuffix

			// Create function spec for reflection analysis
			funcSpec := FunctionSpec{
				Pattern:  fullPattern,
				Function: function,
			}

			// Analyze with reflection
			err := funcSpec.fillFromReflection()
			if err != nil {
				return nil, fmt.Errorf("error analyzing %s: %v", fullPattern, err)
			}

			// Convert to FunctionInfo
			funcInfo := FunctionInfo{
				SymbolName:   symbolName,
				Pattern:      fullPattern,
				FunctionName: funcSpec.FunctionName,
				WrapperName:  funcSpec.WrapperName,
				IsVariadic:   funcSpec.IsVariadic,
				ParamType:    funcSpec.ParamType,
				ParamTypes:   funcSpec.ParamTypes,
				ReturnType:   funcSpec.ReturnType,
				ReturnsError: funcSpec.ReturnsError,
			}

			allFunctions = append(allFunctions, funcInfo)
		}
	}

	return allFunctions, nil
}

// generateSymbolFile generates a wrapper file for a single symbol
func generateSymbolFile(outputPath, symbolName string, functions []FunctionInfo, validationMode string) error {
	tmpl := `// Code generated by wrapgen; DO NOT EDIT.
// Symbol: {{.SymbolName}}
// Validation mode: {{.ValidationMode}}

package wrapped

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/stdlib"
{{- if eq .ValidationMode "debug"}}
	"fmt"
{{- end}}
)

{{range .Functions}}
// {{.WrapperName}} wraps {{.FunctionName}} for the pattern system
// Generated from pattern: {{.Pattern}}
func {{.WrapperName}}(args []core.Expr) core.Expr {
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
	{{- if eq .ParamType "Expr"}}
	result, err := stdlib.{{.FunctionName}}(args...)
	{{- else}}
	result, err := stdlib.{{.FunctionName}}(convertedArgs...)
	{{- end}}
	if err != nil {
		return core.NewErrorExpr(err.Error(), err.Error(), args)
	}
	{{- else}}
	{{- if eq .ParamType "Expr"}}
	result := stdlib.{{.FunctionName}}(args...)
	{{- else}}
	result := stdlib.{{.FunctionName}}(convertedArgs...)
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
		return core.NewErrorExpr("ArgumentError", 
			"{{.FunctionName}} expects {{len .ParamTypes}} arguments", args)
		{{- end}}
	}
	{{- end}}
	
{{getFixedConversionWithMode .ParamTypes $.ValidationMode $.SymbolName | raw}}
	
	// Call business logic function
	{{- if .ReturnsError}}
	result, err := stdlib.{{.FunctionName}}({{.ParamTypes | getCallArgs}})
	if err != nil {
		return core.NewErrorExpr(err.Error(), err.Error(), args)
	}
	{{- else}}
	result := stdlib.{{.FunctionName}}({{.ParamTypes | getCallArgs}})
	{{- end}}
	
	// Convert result back to Expr
	{{.ReturnType | getReturnConversion}}
{{- end}}
}

{{end}}`

	// Template data
	data := struct {
		SymbolName     string
		Functions      []FunctionInfo
		ValidationMode string
	}{
		SymbolName:     symbolName,
		Functions:      functions,
		ValidationMode: validationMode,
	}

	// Use same template functions as before
	funcMap := getTemplateFunctions()

	// Create and execute template
	t, err := template.New("symbol").Funcs(funcMap).Parse(tmpl)
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
func generateBuiltinSetupFile(outputFile string, symbols map[string]SymbolSpec, functions []FunctionInfo) error {
	tmpl := `// Code generated by wrapgen. DO NOT EDIT.

package sexpr

import (
	"fmt"
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/wrapped"
)

// setupBuiltinAttributes sets up standard attributes for built-in functions
func SetupBuiltinAttributes(symbolTable *SymbolTable) {
	// Reset attributes
	symbolTable.Reset()

{{range $name, $symbol := .Symbols}}{{if $symbol.Attributes}}	// {{$name}} attributes
	symbolTable.SetAttributes("{{$name}}", []Attribute{ {{range $i, $attr := $symbol.Attributes}}{{if $i}}, {{end}}{{$attr}}{{end}} })
{{end}}{{end}}
	// Pattern symbols  
	symbolTable.SetAttributes("Blank", []Attribute{Protected})
	symbolTable.SetAttributes("BlankSequence", []Attribute{Protected})
	symbolTable.SetAttributes("BlankNullSequence", []Attribute{Protected})
	symbolTable.SetAttributes("Pattern", []Attribute{Protected})
}

// registerDefaultBuiltins registers all built-in functions with their patterns
func registerDefaultBuiltins(registry *FunctionRegistry) {
	// Register built-in functions with pattern-based dispatch
	builtinPatterns := map[string]PatternFunc{
		// Generated pattern registrations
{{range .Functions}}		"{{.Pattern}}": func(args []core.Expr, ctx *Context) core.Expr {
			return wrapped.{{.WrapperName}}(args)
		}, // {{.FunctionName}}
{{end}}

		// Special attribute manipulation functions (require context)
		"Attributes(x_)":              WrapAttributesExpr,
		"SetAttributes(x_, y_List)":   WrapSetAttributesList,
		"SetAttributes(x_, y_)":       WrapSetAttributesSingle,
		"ClearAttributes(x_, y_List)": WrapClearAttributesList,
		"ClearAttributes(x_, y_)":     WrapClearAttributesSingle,
		
		// Functional programming functions (require context and evaluation)
		"Map(f_, list_)": WrapMapExpr,
		"Apply(f_, list_)": WrapApplyExpr,
		
		// Special debugging functions (require context and main package access)
		"PatternSpecificity(pattern_)":      WrapPatternSpecificity,
		"ShowPatterns(functionName_Symbol)": WrapShowPatterns,
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
		Symbols   map[string]SymbolSpec
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
