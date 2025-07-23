package main

import (
	"flag"
	"fmt"
	"go/format"
	"log"
	"os"
	"strings"
	"text/template"

	// Import stdlib functions for reflection
	"github.com/client9/sexpr/stdlib"
)

// FunctionSpec defines a function and its pattern for wrapper generation
type FunctionSpec struct {
	Pattern    string   // "Plus(x__Integer)" - MANUAL: domain-specific pattern
	OutputFile string   // "arithmetic_wrappers.go" - MANUAL: organization choice
	SymbolName string   // "Plus" - MANUAL: symbol name for attribute setup
	Attributes []string // ["Flat", "Orderless", "OneIdentity"] - MANUAL: domain knowledge

	// HYBRID: Either specify Function (for reflection) OR manual fields (legacy)
	Function interface{} // PlusIntegers - NEW: actual function reference for reflection

	// AUTO-DERIVED: These will be populated by reflection if Function is provided
	FunctionName string   // "PlusIntegers" - derived from Function name
	WrapperName  string   // "WrapPlusIntegers" - derived from FunctionName
	IsVariadic   bool     // derived from Function signature
	ParamType    string   // For variadic: derived from Function signature
	ParamTypes   []string // For fixed arity: derived from Function signature
	ReturnType   string   // derived from Function signature
	ReturnsError bool     // derived from Function signature (has error return)
}

// FunctionGroup groups functions by output file
type FunctionGroup struct {
	OutputFile string
	Functions  []FunctionSpec
}

// Organized function specifications with output file destinations
var functionSpecs = []FunctionSpec{
	{
		Pattern:    "Plus()",
		Function:   stdlib.PlusIdentity,
		OutputFile: "arithmetric_wrappers.go",
	},
	{
		Pattern:    "Times()",
		Function:   stdlib.TimesIdentity,
		OutputFile: "arithmetric_wrappers.go",
	},
	{
		Pattern:    "Plus(x__Integer)",
		Function:   stdlib.PlusIntegers,
		OutputFile: "arithmetic_wrappers.go",
		SymbolName: "Plus",
		Attributes: []string{"Flat", "Listable", "NumericFunction", "OneIdentity", "Orderless", "Protected"},
	},
	{
		Pattern:    "Plus(x__Real)",
		Function:   stdlib.PlusReals,
		OutputFile: "arithmetic_wrappers.go",
		SymbolName: "Plus",
		Attributes: []string{"Flat", "Listable", "NumericFunction", "OneIdentity", "Orderless", "Protected"},
	},
	{
		Pattern:    "Times(x__Integer)",
		Function:   stdlib.TimesIntegers,
		OutputFile: "arithmetic_wrappers.go",
		SymbolName: "Times",
		Attributes: []string{"Flat", "Orderless", "OneIdentity"},
	},
	{
		Pattern:    "Times(x__Real)",
		Function:   stdlib.TimesReals,
		OutputFile: "arithmetic_wrappers.go",
		SymbolName: "Times",
		Attributes: []string{"Flat", "Orderless", "OneIdentity"},
	},
	{
		Pattern:    "Plus(x__Number)",
		Function:   stdlib.PlusNumbers,
		OutputFile: "arithmetic_wrappers.go",
		SymbolName: "Plus",
	},
	{
		Pattern:    "Times(x__Number)",
		Function:   stdlib.TimesNumbers,
		OutputFile: "arithmetic_wrappers.go",
		SymbolName: "Times",
		Attributes: []string{"Flat", "Orderless", "OneIdentity"},
	},
	{
		Pattern:    "Power(base_Real, exp_Integer)",
		Function:   stdlib.PowerReal,
		OutputFile: "arithmetic_wrappers.go",
		SymbolName: "Power",
		Attributes: []string{"OneIdentity"},
	},
	{
		Pattern:    "Subtract(x_Integer, y_Integer)",
		Function:   stdlib.SubtractIntegers,
		OutputFile: "arithmetic_wrappers.go",
		SymbolName: "Subtract",
		Attributes: []string{},
	},
	{
		Pattern:    "Subtract(x_Number, y_Number)",
		Function:   stdlib.SubtractExprs,
		OutputFile: "arithmetic_wrappers.go",
		SymbolName: "Subtract",
		Attributes: []string{},
	},
	{
		Pattern:    "Divide(x_Integer, y_Integer)",
		Function:   stdlib.DivideIntegers,
		OutputFile: "arithmetic_wrappers.go",
		SymbolName: "Divide",
		Attributes: []string{},
	},
	{
		Pattern:    "Divide(x_Number, y_Number)",
		Function:   stdlib.DivideExprs,
		OutputFile: "arithmetic_wrappers.go",
		SymbolName: "Divide",
		Attributes: []string{},
	},
	{
		Pattern:    "Power(x_Number, y_Number)",
		Function:   stdlib.PowerExprs,
		OutputFile: "arithmetic_wrappers.go",
		SymbolName: "Power",
		Attributes: []string{"OneIdentity"},
	},
	{
		Pattern:    "Equal(x_Integer, y_Integer)",
		Function:   stdlib.EqualInts,
		OutputFile: "comparison_wrappers.go",
		SymbolName: "Equal",
		Attributes: []string{},
	},
	{
		Pattern:    "Equal(x_Real, y_Real)",
		Function:   stdlib.EqualFloats,
		OutputFile: "comparison_wrappers.go",
	},
	{
		Pattern:    "Equal(x_Number, y_Number)",
		Function:   stdlib.EqualNumbers,
		OutputFile: "comparison_wrappers.go",
	},
	{
		Pattern:    "Equal(x_String, y_String)",
		Function:   stdlib.EqualStrings,
		OutputFile: "comparison_wrappers.go",
	},
	{
		Pattern:    "Unequal(x_Integer, y_Integer)",
		Function:   stdlib.UnequalInts,
		OutputFile: "comparison_wrappers.go",
		SymbolName: "Unequal",
		Attributes: []string{},
	},
	{
		Pattern:    "Unequal(x_Real, y_Real)",
		Function:   stdlib.UnequalFloats,
		OutputFile: "comparison_wrappers.go",
	},
	{
		Pattern:    "Unequal(x_Number, y_Number)",
		Function:   stdlib.UnequalNumbers,
		OutputFile: "comparison_wrappers.go",
	},
	{
		Pattern:    "Unequal(x_String, y_String)",
		Function:   stdlib.UnequalStrings,
		OutputFile: "comparison_wrappers.go",
	},
	{
		Pattern:    "Less(x_Number, y_Number)",
		Function:   stdlib.LessNumber,
		OutputFile: "comparison_wrappers.go",
		SymbolName: "Less",
		Attributes: []string{},
	},
	{
		Pattern:    "Greater(x_Number, y_Number))",
		Function:   stdlib.GreaterNumber,
		OutputFile: "comparison_wrappers.go",
		SymbolName: "Greater",
		Attributes: []string{},
	},
	{
		Pattern:    "LessEqual(x_Number, y_Number)",
		Function:   stdlib.LessEqualNumber,
		OutputFile: "comparison_wrappers.go",
		SymbolName: "LessEqual",
		Attributes: []string{},
	},
	{
		Pattern:    "GreaterEqual(x_Number, y_Number)",
		Function:   stdlib.GreaterEqualNumber,
		OutputFile: "comparison_wrappers.go",
		SymbolName: "GreaterEqual",
		Attributes: []string{},
	},
	{
		Pattern:    "Equal(x_, y_)",
		Function:   stdlib.EqualExprs,
		OutputFile: "comparison_wrappers.go",
	},
	{
		Pattern:    "Unequal(x_, y_)",
		Function:   stdlib.UnequalExprs,
		OutputFile: "comparison_wrappers.go",
	},
	{
		Pattern:    "SameQ(x_, y_)",
		Function:   stdlib.SameQExprs,
		OutputFile: "comparison_wrappers.go",
		SymbolName: "SameQ",
		Attributes: []string{},
	},
	{
		Pattern:    "UnsameQ(x_, y_)",
		Function:   stdlib.UnsameQExprs,
		OutputFile: "comparison_wrappers.go",
		SymbolName: "UnsameQ",
		Attributes: []string{},
	},
	{
		Pattern:    "IntegerQ(x_)",
		Function:   stdlib.IntegerQExpr,
		OutputFile: "type_predicate_wrappers.go",
	},
	{
		Pattern:    "FloatQ(x_)",
		Function:   stdlib.FloatQExpr,
		OutputFile: "type_predicate_wrappers.go",
	},
	{
		Pattern:    "NumberQ(x_)",
		Function:   stdlib.NumberQExpr,
		OutputFile: "type_predicate_wrappers.go",
	},
	{
		Pattern:    "StringQ(x_)",
		Function:   stdlib.StringQExpr,
		OutputFile: "type_predicate_wrappers.go",
	},
	{
		Pattern:    "BooleanQ(x_)",
		Function:   stdlib.BooleanQExpr,
		OutputFile: "type_predicate_wrappers.go",
	},
	{
		Pattern:    "SymbolQ(x_)",
		Function:   stdlib.SymbolQExpr,
		OutputFile: "type_predicate_wrappers.go",
	},
	{
		Pattern:    "ListQ(x_)",
		Function:   stdlib.ListQExpr,
		OutputFile: "type_predicate_wrappers.go",
	},
	{
		Pattern:    "AtomQ(x_)",
		Function:   stdlib.AtomQExpr,
		OutputFile: "type_predicate_wrappers.go",
	},
	{
		Pattern:    "Head(x_)",
		Function:   stdlib.HeadExpr,
		OutputFile: "type_predicate_wrappers.go",
	},
	{
		Pattern:    "FullForm(x_)",
		Function:   stdlib.FullFormExpr,
		OutputFile: "output_format_wrappers.go",
	},
	{
		Pattern:    "InputForm(x_)",
		Function:   stdlib.InputFormExpr,
		OutputFile: "output_format_wrappers.go",
	},
	{
		Pattern:    "Length(x_)",
		Function:   stdlib.LengthExpr,
		OutputFile: "list_wrappers.go",
		SymbolName: "Length",
		Attributes: []string{},
	},
	{
		Pattern:    "First(x_List)",
		Function:   stdlib.FirstExpr,
		OutputFile: "list_wrappers.go",
		SymbolName: "First",
		Attributes: []string{},
	},
	{
		Pattern:    "Last(x_List)",
		Function:   stdlib.LastExpr,
		OutputFile: "list_wrappers.go",
		SymbolName: "Last",
		Attributes: []string{},
	},
	{
		Pattern:    "Rest(x_List)",
		Function:   stdlib.RestExpr,
		OutputFile: "list_wrappers.go",
		SymbolName: "Rest",
		Attributes: []string{},
	},
	{
		Pattern:    "Most(x_List)",
		Function:   stdlib.MostExpr,
		OutputFile: "list_wrappers.go",
		SymbolName: "Most",
		Attributes: []string{},
	},
	{
		Pattern:    "Take(x_List, n_Integer)",
		Function:   stdlib.TakeList,
		OutputFile: "list_wrappers.go",
		SymbolName: "Take",
		Attributes: []string{},
	},
	{
		Pattern:    "Take(x_List, List(n_Integer))",
		Function:   stdlib.TakeListSingle,
		OutputFile: "list_wrappers.go",
		SymbolName: "Take",
		Attributes: []string{},
	},
	{
		Pattern:    "Take(x_List, List(n_Integer, m_Integer))",
		Function:   stdlib.TakeListRange,
		OutputFile: "list_wrappers.go",
		SymbolName: "Take",
		Attributes: []string{},
	},
	{
		Pattern:    "Drop(x_List, n_Integer)",
		Function:   stdlib.DropList,
		OutputFile: "list_wrappers.go",
		SymbolName: "Drop",
		Attributes: []string{},
	},
	{
		Pattern:    "Drop(x_List, List(n_Integer))",
		Function:   stdlib.DropListSingle,
		OutputFile: "list_wrappers.go",
		SymbolName: "Drop",
		Attributes: []string{},
	},
	{
		Pattern:    "Drop(x_List, List(n_Integer, m_Integer))",
		Function:   stdlib.DropListRange,
		OutputFile: "list_wrappers.go",
		SymbolName: "Drop",
		Attributes: []string{},
	},
	{
		Pattern:    "Not(x_)",
		Function:   stdlib.NotExpr,
		OutputFile: "logical_wrappers.go",
		SymbolName: "Not",
		Attributes: []string{},
	},
	{
		Pattern:    "StringLength(x_String)",
		Function:   stdlib.StringLengthRunes,
		OutputFile: "string_wrappers.go",
		SymbolName: "StringLength",
		Attributes: []string{},
	},
	{
		Pattern:    "AssociationQ(x_)",
		Function:   stdlib.AssociationQExpr,
		OutputFile: "association_wrappers.go",
		SymbolName: "AssociationQ",
		Attributes: []string{},
	},
	{
		Pattern:    "Keys(x_Association)",
		Function:   stdlib.KeysExpr,
		OutputFile: "association_wrappers.go",
		SymbolName: "Keys",
		Attributes: []string{},
	},
	{
		Pattern:    "Values(x_Association)",
		Function:   stdlib.ValuesExpr,
		OutputFile: "association_wrappers.go",
		SymbolName: "Values",
		Attributes: []string{},
	},
	{
		Pattern:    "Association(x___Rule)",
		Function:   stdlib.AssociationRules,
		OutputFile: "association_wrappers.go",
		SymbolName: "Association",
		Attributes: []string{},
	},
	{
		Pattern:    "Part(x_Association, y_)",
		Function:   stdlib.PartAssociation,
		OutputFile: "association_wrappers.go",
		SymbolName: "Part",
		Attributes: []string{},
	},
}

// processFunctionSpecs processes all function specs, filling in auto-derived fields via reflection
func processFunctionSpecs(specs []FunctionSpec) ([]FunctionSpec, error) {
	processed := make([]FunctionSpec, len(specs))
	copy(processed, specs)

	for i := range processed {
		err := processed[i].fillFromReflection()
		if err != nil {
			return nil, fmt.Errorf("error processing spec %d (%s): %v", i, processed[i].Pattern, err)
		}
	}

	return processed, nil
}

func main() {
	var (
		outputDir = flag.String("dir", ".", "Output directory for generated files")
		single    = flag.String("single", "", "Generate single file with all wrappers")
		setupFile = flag.String("setup", "", "Generate builtin setup file (builtin_setup.go)")
	)
	flag.Parse()

	// Process function specs with reflection analysis
	processedSpecs, err := processFunctionSpecs(functionSpecs)
	if err != nil {
		log.Fatalf("Error processing function specs: %v", err)
	}

	if *setupFile != "" {
		// Generate builtin_setup.go file
		err := generateBuiltinSetupFile(*setupFile, processedSpecs)
		if err != nil {
			log.Fatalf("Error generating setup file: %v", err)
		}
		fmt.Printf("Generated builtin setup file: %s\n", *setupFile)
		return
	}

	if *single != "" {
		// Generate all functions in a single file
		err := generateSingleFile(*single, processedSpecs)
		if err != nil {
			log.Fatalf("Error generating single file: %v", err)
		}
		fmt.Printf("Generated %d wrappers in %s\n", len(processedSpecs), *single)
		return
	}

	// Group functions by output file
	groups := groupFunctionsByFile(processedSpecs)

	// Generate each file
	totalFunctions := 0
	for _, group := range groups {
		outputPath := fmt.Sprintf("%s/%s", strings.TrimSuffix(*outputDir, "/"), group.OutputFile)
		err := generateWrapperFile(outputPath, group.Functions)
		if err != nil {
			log.Fatalf("Error generating %s: %v", outputPath, err)
		}
		fmt.Printf("Generated %d wrappers in %s\n", len(group.Functions), outputPath)
		totalFunctions += len(group.Functions)
	}

	fmt.Printf("Total: %d wrappers across %d files\n", totalFunctions, len(groups))
}

// groupFunctionsByFile groups function specs by their output file
func groupFunctionsByFile(specs []FunctionSpec) []FunctionGroup {
	fileMap := make(map[string][]FunctionSpec)

	for _, spec := range specs {
		fileMap[spec.OutputFile] = append(fileMap[spec.OutputFile], spec)
	}

	var groups []FunctionGroup
	for outputFile, functions := range fileMap {
		groups = append(groups, FunctionGroup{
			OutputFile: outputFile,
			Functions:  functions,
		})
	}

	return groups
}

// generateSingleFile generates all wrappers in a single file
func generateSingleFile(outputFile string, specs []FunctionSpec) error {
	return generateWrapperFile(outputFile, specs)
}

// generateWrapperFile generates a wrapper file for the given functions
func generateWrapperFile(outputFile string, functions []FunctionSpec) error {
	tmpl := `// Code generated by wrapgen; DO NOT EDIT.

package sexpr

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/stdlib"
)

{{range .}}
// {{.WrapperName}} wraps {{.FunctionName}} for the pattern system
// Generated from pattern: {{.Pattern}}
func {{.WrapperName}}(args []core.Expr, ctx *Context) core.Expr {
{{- if .IsVariadic}}
	{{- if ne .ParamType "Expr"}}
	funcName := "{{.Pattern | extractFuncName}}"
	{{- end}}
	
	// Convert all args to {{.ParamType}}
	convertedArgs := make([]{{.ParamType}}, len(args))
	for i, arg := range args {
		{{.ParamType | getConversion}}
	}
	
	// Call business logic function
	{{- if .ReturnsError}}
	result, err := stdlib.{{.FunctionName}}(convertedArgs...)
	if err != nil {
		return NewErrorExpr(err.Error(), err.Error(), args)
	}
	{{- else}}
	result := stdlib.{{.FunctionName}}(convertedArgs...)
	{{- end}}
	
	// Convert result back to Expr
	{{.ReturnType | getReturnConversion}}
{{- else}}
	// Validate argument count
	if len(args) != {{len .ParamTypes}} {
		return NewErrorExpr("ArgumentError",
			"{{.FunctionName}} expects {{len .ParamTypes}} arguments", args)
	}
	
{{.ParamTypes | getFixedConversion | raw}}
	
	// Call business logic function
	{{- if .ReturnsError}}
	result, err := stdlib.{{.FunctionName}}({{.ParamTypes | getCallArgs}})
	if err != nil {
		return NewErrorExpr(err.Error(), err.Error(), args)
	}
	{{- else}}
	result := stdlib.{{.FunctionName}}({{.ParamTypes | getCallArgs}})
	{{- end}}
	
	// Convert result back to Expr
	{{.ReturnType | getReturnConversion}}
{{- end}}
}

{{end}}`

	// Custom template functions using the new public helper functions
	funcMap := template.FuncMap{
		"extractFuncName": func(pattern string) string {
			for i, c := range pattern {
				if c == '(' {
					return pattern[:i]
				}
			}
			return pattern
		},
		"getEmptyCase": func(funcName string) string {
			switch funcName {
			default:
				return "return CopyExprList(funcName, args)"
			}
		},
		"getSingleCase": func(paramType string) string {
			switch paramType {
			case "int64":
				return `if atom, ok := args[0].(Atom); ok && atom.AtomType == IntAtom {
			return args[0] // Return directly
		}
		// Fall back to original if not integer
		return CopyExprList(funcName, args)`
			case "float64":
				return `if atom, ok := args[0].(Atom); ok && atom.AtomType == FloatAtom {
			return args[0] // Return directly
		}
		// Fall back to original if not real
		return CopyExprList(funcName, args)`
			case "Expr":
				return "return args[0]"
			default:
				return "return args[0]"
			}
		},
		"getConversion": func(paramType string) string {
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
		},
		"getReturnConversion": func(returnType string) string {
			switch returnType {
			case "int64":
				return "return core.NewIntAtom(int(result))"
			case "float64":
				return "return core.NewFloatAtom(result)"
			case "string":
				return "return core.NewStringAtom(result)"
			case "bool":
				return "return core.NewBoolAtom(result)"
			case "Expr":
				return "return result"
			default:
				return "return result"
			}
		},
		"getFixedConversion": func(paramTypes []string) string {
			var conversions []string
			for i, paramType := range paramTypes {
				varName := fmt.Sprintf("arg%d", i)
				switch paramType {
				case "Number":

					conversions = append(conversions, fmt.Sprintf("	%s, ok := stdlib.ExtractNumber(args[%d])", varName, i))
					conversions = append(conversions, "	if !ok {")
					conversions = append(conversions, "		return core.CopyExprList(\"FUNC\", args)")
					conversions = append(conversions, "	}")
				case "AssociationValue":
					conversions = append(conversions, fmt.Sprintf("	%s, ok := stdlib.ExtractAssociation(args[%d])", varName, i))
					conversions = append(conversions, "	if !ok {")
					conversions = append(conversions, "		return core.CopyExprList(\"FUNC\", args)")
					conversions = append(conversions, "	}")
				case "Expr":
					conversions = append(conversions, fmt.Sprintf("	%s := args[%d]", varName, i))
				case "int64":
					conversions = append(conversions, fmt.Sprintf("	%s, ok := core.ExtractInt64(args[%d])", varName, i))
					conversions = append(conversions, "	if !ok {")
					conversions = append(conversions, "		return core.CopyExprList(\"FUNC\", args)")
					conversions = append(conversions, "	}")
				case "float64":
					conversions = append(conversions, fmt.Sprintf("	%s, ok := core.ExtractFloat64(args[%d])", varName, i))
					conversions = append(conversions, "	if !ok {")
					conversions = append(conversions, "		return core.CopyExprList(\"FUNC\", args)")
					conversions = append(conversions, "	}")
				case "string":
					conversions = append(conversions, fmt.Sprintf("	%s, ok := core.ExtractString(args[%d])", varName, i))
					conversions = append(conversions, "	if !ok {")
					conversions = append(conversions, "		return core.CopyExprList(\"FUNC\", args)")
					conversions = append(conversions, "	}")
				case "bool":
					conversions = append(conversions, fmt.Sprintf("	%s, ok := core.ExtractBool(args[%d])", varName, i))
					conversions = append(conversions, "	if !ok {")
					conversions = append(conversions, "		return core.CopyExprList(\"FUNC\", args)")
					conversions = append(conversions, "	}")
				case "List":
					conversions = append(conversions, fmt.Sprintf("	%s, ok := args[%d].(core.List)", varName, i))
					conversions = append(conversions, "	if !ok {")
					conversions = append(conversions, "		return core.CopyExprList(\"FUNC\", args)")
					conversions = append(conversions, "	}")
				case "ObjectExpr":
					conversions = append(conversions, fmt.Sprintf("	%s, ok := args[%d].(core.ObjectExpr)", varName, i))
					conversions = append(conversions, "	if !ok {")
					conversions = append(conversions, "		return core.CopyExprList(\"FUNC\", args)")
					conversions = append(conversions, "	}")
				default:
					log.Fatalf("Unknown Parameter Type: %s", paramType)
				}
			}
			return strings.Join(conversions, "\n")
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

	// Create and execute template
	t, err := template.New("wrappers").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		return err
	}

	// Generate code
	var buf strings.Builder
	err = t.Execute(&buf, functions)
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

// generateBuiltinSetupFile generates a builtin_setup.go file with both attribute setup and registration
func generateBuiltinSetupFile(outputFile string, functions []FunctionSpec) error {
	tmpl := `// Code generated by wrapgen. DO NOT EDIT.

package sexpr

import (
	"fmt"
)

// setupBuiltinAttributes sets up standard attributes for built-in functions
func setupBuiltinAttributes(symbolTable *SymbolTable) {
	// Reset attributes
	symbolTable.Reset()

{{range .UniqueSymbols}}{{if .Attributes}}	// {{.SymbolName}} attributes
	symbolTable.SetAttributes("{{.SymbolName}}", []Attribute{ {{range $i, $attr := .Attributes}}{{if $i}}, {{end}}{{$attr}}{{end}} })
{{end}}{{end}}
	// Constants
	symbolTable.SetAttributes("Pi", []Attribute{Constant, Protected})
	symbolTable.SetAttributes("E", []Attribute{Constant, Protected})
	symbolTable.SetAttributes("True", []Attribute{Constant, Protected})
	symbolTable.SetAttributes("False", []Attribute{Constant, Protected})

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
{{range .Functions}}		"{{.Pattern}}": {{.WrapperName}}, // {{.FunctionName}}
{{end}}

		// Special attribute manipulation functions (require context)
		"Attributes(x_)":           EvaluateAttributes,
		"SetAttributes(x_, y_)":    EvaluateSetAttributes,
		"ClearAttributes(x_, y_)":  EvaluateClearAttributes,
		
		// Pattern matching functions
		"MatchQ(x_, y_)":           EvaluateMatchQ,
	}

	// Register patterns with the function registry
	err := registry.RegisterPatternBuiltins(builtinPatterns)
	if err != nil {
		panic(fmt.Sprintf("Failed to register builtin patterns: %v", err))
	}
}
`

	// Create template data structure
	type TemplateData struct {
		Functions     []FunctionSpec
		UniqueSymbols []FunctionSpec
	}

	// Get unique symbols (deduplicated by SymbolName)
	symbolMap := make(map[string]FunctionSpec)
	for _, fn := range functions {
		if fn.SymbolName != "" {
			// Keep the one with the most attributes
			if existing, exists := symbolMap[fn.SymbolName]; !exists || len(fn.Attributes) > len(existing.Attributes) {
				symbolMap[fn.SymbolName] = fn
			}
		}
	}

	var uniqueSymbols []FunctionSpec
	for _, symbol := range symbolMap {
		uniqueSymbols = append(uniqueSymbols, symbol)
	}

	data := TemplateData{
		Functions:     functions,
		UniqueSymbols: uniqueSymbols,
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
