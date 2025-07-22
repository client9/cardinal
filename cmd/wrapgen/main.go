package main

import (
	"flag"
	"fmt"
	"go/format"
	"log"
	"os"
	"strings"
	"text/template"
)

// FunctionSpec defines a function and its pattern for wrapper generation
type FunctionSpec struct {
	Pattern      string   // "Plus(x__Integer)"
	FunctionName string   // "PlusIntegers" 
	WrapperName  string   // "WrapPlusIntegers"
	OutputFile   string   // "arithmetic_wrappers.go" - can split functions across files
	IsVariadic   bool
	ParamType    string   // For variadic: "int64", "float64", etc.
	ParamTypes   []string // For fixed arity: ["Expr", "Expr"] or ["float64", "int64"]
	ReturnType   string
	ReturnsError bool     // If true, function returns (T, error) and wrapper handles error cases
}

// FunctionGroup groups functions by output file
type FunctionGroup struct {
	OutputFile string
	Functions  []FunctionSpec
}

// Organized function specifications with output file destinations
var functionSpecs = []FunctionSpec{
	// Arithmetic functions -> arithmetic_wrappers.go
	{
		Pattern:      "Plus(x__Integer)",
		FunctionName: "PlusIntegers", 
		WrapperName:  "WrapPlusIntegers",
		OutputFile:   "arithmetic_wrappers.go",
		IsVariadic:   true,
		ParamType:    "int64",
		ReturnType:   "int64",
	},
	{
		Pattern:      "Plus(x__Real)",
		FunctionName: "PlusReals", 
		WrapperName:  "WrapPlusReals",
		OutputFile:   "arithmetic_wrappers.go",
		IsVariadic:   true,
		ParamType:    "float64",
		ReturnType:   "float64",
	},
	{
		Pattern:      "Times(x__Integer)",
		FunctionName: "TimesIntegers",
		WrapperName:  "WrapTimesIntegers", 
		OutputFile:   "arithmetic_wrappers.go",
		IsVariadic:   true,
		ParamType:    "int64",
		ReturnType:   "int64",
	},
	{
		Pattern:      "Times(x__Real)",
		FunctionName: "TimesReals",
		WrapperName:  "WrapTimesReals", 
		OutputFile:   "arithmetic_wrappers.go",
		IsVariadic:   true,
		ParamType:    "float64",
		ReturnType:   "float64",
	},
	{
		Pattern:      "Plus(x__Number)",
		FunctionName: "PlusNumbers",
		WrapperName:  "WrapPlusNumbers",
		OutputFile:   "arithmetic_wrappers.go",
		IsVariadic:   true,
		ParamType:    "Expr",
		ReturnType:   "float64",
	},
	{
		Pattern:      "Times(x__Number)",
		FunctionName: "TimesNumbers", 
		WrapperName:  "WrapTimesNumbers",
		OutputFile:   "arithmetic_wrappers.go",
		IsVariadic:   true,
		ParamType:    "Expr",
		ReturnType:   "float64",
	},
	{
		Pattern:      "Power(base_Real, exp_Integer)", 
		FunctionName: "PowerReal",
		WrapperName:  "WrapPowerReal",
		OutputFile:   "arithmetic_wrappers.go",
		IsVariadic:   false,
		ParamTypes:   []string{"float64", "int64"},
		ReturnType:   "float64",
	},
	{
		Pattern:      "Subtract(x_Integer, y_Integer)",
		FunctionName: "SubtractIntegers",
		WrapperName:  "WrapSubtractIntegers",
		OutputFile:   "arithmetic_wrappers.go",
		IsVariadic:   false,
		ParamTypes:   []string{"int64", "int64"},
		ReturnType:   "int64",
	},
	{
		Pattern:      "Subtract(x_Number, y_Number)",
		FunctionName: "SubtractNumbers",
		WrapperName:  "WrapSubtractNumbers",
		OutputFile:   "arithmetic_wrappers.go",
		IsVariadic:   false,
		ParamTypes:   []string{"Expr", "Expr"},
		ReturnType:   "float64",
	},
	{
		Pattern:      "Divide(x_Integer, y_Integer)",
		FunctionName: "DivideIntegers",
		WrapperName:  "WrapDivideIntegers",
		OutputFile:   "arithmetic_wrappers.go",
		IsVariadic:   false,
		ParamTypes:   []string{"int64", "int64"},
		ReturnType:   "int64",
		ReturnsError: true,
	},
	{
		Pattern:      "Divide(x_Number, y_Number)",
		FunctionName: "DivideNumbers",
		WrapperName:  "WrapDivideNumbers",
		OutputFile:   "arithmetic_wrappers.go",
		IsVariadic:   false,
		ParamTypes:   []string{"Expr", "Expr"},
		ReturnType:   "float64",
		ReturnsError: true,
	},
	{
		Pattern:      "Power(x_Number, y_Number)",
		FunctionName: "PowerNumbers",
		WrapperName:  "WrapPowerNumbers",
		OutputFile:   "arithmetic_wrappers.go",
		IsVariadic:   false,
		ParamTypes:   []string{"Expr", "Expr"},
		ReturnType:   "float64",
		ReturnsError: true,
	},
	
	// Comparison functions -> comparison_wrappers.go
	{
		Pattern:      "Equal(x_, y_)",
		FunctionName: "EqualExprs",
		WrapperName:  "WrapEqualExprs",
		OutputFile:   "comparison_wrappers.go",
		IsVariadic:   false,
		ParamTypes:   []string{"Expr", "Expr"},
		ReturnType:   "bool",
	},
	{
		Pattern:      "Unequal(x_, y_)",
		FunctionName: "UnequalExprs",
		WrapperName:  "WrapUnequalExprs",
		OutputFile:   "comparison_wrappers.go",
		IsVariadic:   false,
		ParamTypes:   []string{"Expr", "Expr"},
		ReturnType:   "bool",
	},
	{
		Pattern:      "Less(x_, y_)",
		FunctionName: "LessExprs",
		WrapperName:  "WrapLessExprs",
		OutputFile:   "comparison_wrappers.go",
		IsVariadic:   false,
		ParamTypes:   []string{"Expr", "Expr"},
		ReturnType:   "bool",
	},
	{
		Pattern:      "Greater(x_, y_)",
		FunctionName: "GreaterExprs",
		WrapperName:  "WrapGreaterExprs",
		OutputFile:   "comparison_wrappers.go",
		IsVariadic:   false,
		ParamTypes:   []string{"Expr", "Expr"},
		ReturnType:   "bool",
	},
	{
		Pattern:      "LessEqual(x_, y_)",
		FunctionName: "LessEqualExprs",
		WrapperName:  "WrapLessEqualExprs",
		OutputFile:   "comparison_wrappers.go",
		IsVariadic:   false,
		ParamTypes:   []string{"Expr", "Expr"},
		ReturnType:   "bool",
	},
	{
		Pattern:      "GreaterEqual(x_, y_)",
		FunctionName: "GreaterEqualExprs",
		WrapperName:  "WrapGreaterEqualExprs",
		OutputFile:   "comparison_wrappers.go",
		IsVariadic:   false,
		ParamTypes:   []string{"Expr", "Expr"},
		ReturnType:   "bool",
	},
	{
		Pattern:      "SameQ(x_, y_)",
		FunctionName: "SameQExprs",
		WrapperName:  "WrapSameQExprs",
		OutputFile:   "comparison_wrappers.go",
		IsVariadic:   false,
		ParamTypes:   []string{"Expr", "Expr"},
		ReturnType:   "bool",
	},
	{
		Pattern:      "UnsameQ(x_, y_)",
		FunctionName: "UnsameQExprs",
		WrapperName:  "WrapUnsameQExprs",
		OutputFile:   "comparison_wrappers.go",
		IsVariadic:   false,
		ParamTypes:   []string{"Expr", "Expr"},
		ReturnType:   "bool",
	},
	
	// Type predicate functions -> type_predicate_wrappers.go
	{
		Pattern:      "IntegerQ(x_)",
		FunctionName: "IntegerQExpr",
		WrapperName:  "WrapIntegerQExpr",
		OutputFile:   "type_predicate_wrappers.go",
		IsVariadic:   false,
		ParamTypes:   []string{"Expr"},
		ReturnType:   "bool",
	},
	{
		Pattern:      "FloatQ(x_)",
		FunctionName: "FloatQExpr",
		WrapperName:  "WrapFloatQExpr",
		OutputFile:   "type_predicate_wrappers.go",
		IsVariadic:   false,
		ParamTypes:   []string{"Expr"},
		ReturnType:   "bool",
	},
	{
		Pattern:      "NumberQ(x_)",
		FunctionName: "NumberQExpr",
		WrapperName:  "WrapNumberQExpr",
		OutputFile:   "type_predicate_wrappers.go",
		IsVariadic:   false,
		ParamTypes:   []string{"Expr"},
		ReturnType:   "bool",
	},
	{
		Pattern:      "StringQ(x_)",
		FunctionName: "StringQExpr",
		WrapperName:  "WrapStringQExpr",
		OutputFile:   "type_predicate_wrappers.go",
		IsVariadic:   false,
		ParamTypes:   []string{"Expr"},
		ReturnType:   "bool",
	},
	{
		Pattern:      "BooleanQ(x_)",
		FunctionName: "BooleanQExpr",
		WrapperName:  "WrapBooleanQExpr",
		OutputFile:   "type_predicate_wrappers.go",
		IsVariadic:   false,
		ParamTypes:   []string{"Expr"},
		ReturnType:   "bool",
	},
	{
		Pattern:      "SymbolQ(x_)",
		FunctionName: "SymbolQExpr",
		WrapperName:  "WrapSymbolQExpr",
		OutputFile:   "type_predicate_wrappers.go",
		IsVariadic:   false,
		ParamTypes:   []string{"Expr"},
		ReturnType:   "bool",
	},
	{
		Pattern:      "ListQ(x_)",
		FunctionName: "ListQExpr",
		WrapperName:  "WrapListQExpr",
		OutputFile:   "type_predicate_wrappers.go",
		IsVariadic:   false,
		ParamTypes:   []string{"Expr"},
		ReturnType:   "bool",
	},
	{
		Pattern:      "AtomQ(x_)",
		FunctionName: "AtomQExpr",
		WrapperName:  "WrapAtomQExpr",
		OutputFile:   "type_predicate_wrappers.go",
		IsVariadic:   false,
		ParamTypes:   []string{"Expr"},
		ReturnType:   "bool",
	},
	{
		Pattern:      "Head(x_)",
		FunctionName: "HeadExpr",
		WrapperName:  "WrapHeadExpr",
		OutputFile:   "type_predicate_wrappers.go",
		IsVariadic:   false,
		ParamTypes:   []string{"Expr"},
		ReturnType:   "Expr",
	},
	{
		Pattern:      "FullForm(x_)",
		FunctionName: "FullFormExpr",
		WrapperName:  "WrapFullFormExpr",
		OutputFile:   "output_format_wrappers.go",
		IsVariadic:   false,
		ParamTypes:   []string{"Expr"},
		ReturnType:   "string",
	},
	{
		Pattern:      "InputForm(x_)",
		FunctionName: "InputFormExpr",
		WrapperName:  "WrapInputFormExpr",
		OutputFile:   "output_format_wrappers.go",
		IsVariadic:   false,
		ParamTypes:   []string{"Expr"},
		ReturnType:   "string",
	},
	{
		Pattern:      "Length(x_)",
		FunctionName: "LengthExpr",
		WrapperName:  "WrapLengthExpr",
		OutputFile:   "length_wrappers.go",
		IsVariadic:   false,
		ParamTypes:   []string{"Expr"},
		ReturnType:   "int64",
	},
	
	// String functions -> string_wrappers.go
	{
		Pattern:      "StringLength(x_String)",
		FunctionName: "StringLengthStr",
		WrapperName:  "WrapStringLengthStr",
		OutputFile:   "string_wrappers.go",
		IsVariadic:   false,
		ParamTypes:   []string{"string"},
		ReturnType:   "int64",
	},
	
	// List access functions -> list_access_wrappers.go
	{
		Pattern:      "First(x_List)",
		FunctionName: "FirstExpr",
		WrapperName:  "WrapFirstExpr",
		OutputFile:   "list_access_wrappers.go",
		IsVariadic:   false,
		ParamTypes:   []string{"List"},
		ReturnType:   "Expr",
	},
	{
		Pattern:      "Last(x_List)",
		FunctionName: "LastExpr",
		WrapperName:  "WrapLastExpr",
		OutputFile:   "list_access_wrappers.go",
		IsVariadic:   false,
		ParamTypes:   []string{"List"},
		ReturnType:   "Expr",
	},
	{
		Pattern:      "Rest(x_List)",
		FunctionName: "RestExpr",
		WrapperName:  "WrapRestExpr",
		OutputFile:   "list_access_wrappers.go",
		IsVariadic:   false,
		ParamTypes:   []string{"List"},
		ReturnType:   "Expr",
	},
	{
		Pattern:      "Most(x_List)",
		FunctionName: "MostExpr",
		WrapperName:  "WrapMostExpr",
		OutputFile:   "list_access_wrappers.go",
		IsVariadic:   false,
		ParamTypes:   []string{"List"},
		ReturnType:   "Expr",
	},
	
	// Association functions -> association_wrappers.go
	{
		Pattern:      "Association(x___Rule)",
		FunctionName: "AssociationRules",
		WrapperName:  "WrapAssociationRules",
		OutputFile:   "association_wrappers.go",
		IsVariadic:   true,
		ParamType:    "Expr", // Rules are passed as Expr (List expressions)
		ReturnType:   "Expr",
	},
	{
		Pattern:      "AssociationQ(x_)",
		FunctionName: "AssociationQExpr",
		WrapperName:  "WrapAssociationQExpr",
		OutputFile:   "association_wrappers.go",
		IsVariadic:   false,
		ParamTypes:   []string{"Expr"},
		ReturnType:   "bool",
	},
	{
		Pattern:      "Keys(x_Association)",
		FunctionName: "KeysExpr",
		WrapperName:  "WrapKeysExpr",
		OutputFile:   "association_wrappers.go",
		IsVariadic:   false,
		ParamTypes:   []string{"ObjectExpr"},
		ReturnType:   "Expr",
	},
	{
		Pattern:      "Values(x_Association)",
		FunctionName: "ValuesExpr",
		WrapperName:  "WrapValuesExpr",
		OutputFile:   "association_wrappers.go",
		IsVariadic:   false,
		ParamTypes:   []string{"ObjectExpr"},
		ReturnType:   "Expr",
	},
	
	// Logical functions -> logical_wrappers.go
	{
		Pattern:      "Not(x_)",
		FunctionName: "NotExpr",
		WrapperName:  "WrapNotExpr",
		OutputFile:   "logical_wrappers.go",
		IsVariadic:   false,
		ParamTypes:   []string{"Expr"},
		ReturnType:   "Expr",
	},
	
	// Part functions -> part_wrappers.go
	{
		Pattern:      "Part(x_List, i_Integer)",
		FunctionName: "PartList",
		WrapperName:  "WrapPartList",
		OutputFile:   "part_wrappers.go",
		IsVariadic:   false,
		ParamTypes:   []string{"List", "int64"},
		ReturnType:   "Expr",
	},
	{
		Pattern:      "Part(x_Association, y_)",
		FunctionName: "PartAssociation",
		WrapperName:  "WrapPartAssociation",
		OutputFile:   "part_wrappers.go",
		IsVariadic:   false,
		ParamTypes:   []string{"ObjectExpr", "Expr"},
		ReturnType:   "Expr",
	},
}

func main() {
	var (
		outputDir = flag.String("dir", ".", "Output directory for generated files")
		single    = flag.String("single", "", "Generate single file with all wrappers")
	)
	flag.Parse()

	if *single != "" {
		// Generate all functions in a single file
		err := generateSingleFile(*single, functionSpecs)
		if err != nil {
			log.Fatalf("Error generating single file: %v", err)
		}
		fmt.Printf("Generated %d wrappers in %s\n", len(functionSpecs), *single)
		return
	}

	// Group functions by output file
	groups := groupFunctionsByFile(functionSpecs)

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

{{range .}}
// {{.WrapperName}} wraps {{.FunctionName}} for the pattern system
// Generated from pattern: {{.Pattern}}
func {{.WrapperName}}(args []Expr, ctx *Context) Expr {
{{- if .IsVariadic}}
	{{- if ne .ParamType "Expr"}}
	funcName := "{{.Pattern | extractFuncName}}"
	{{- end}}
	
	// Handle empty case
	if len(args) == 0 {
		{{.FunctionName | getEmptyCase}}
	}
	
	// Handle single arg case  
	if len(args) == 1 {
		{{if eq .FunctionName "AssociationRules"}}result := AssociationRules(args[0]); return result{{else}}{{.ParamType | getSingleCase}}{{end}}
	}
	
	// Convert all args to {{.ParamType}}
	convertedArgs := make([]{{.ParamType}}, len(args))
	for i, arg := range args {
		{{.ParamType | getConversion}}
	}
	
	// Call business logic function
	{{- if .ReturnsError}}
	result, err := {{.FunctionName}}(convertedArgs...)
	if err != nil {
		return NewErrorExpr(err.Error(), err.Error(), args)
	}
	{{- else}}
	result := {{.FunctionName}}(convertedArgs...)
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
	result, err := {{.FunctionName}}({{.ParamTypes | getCallArgs}})
	if err != nil {
		return NewErrorExpr(err.Error(), err.Error(), args)
	}
	{{- else}}
	result := {{.FunctionName}}({{.ParamTypes | getCallArgs}})
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
			case "PlusIntegers":
				return "return NewIntAtom(0)"
			case "PlusReals":
				return "return NewFloatAtom(0.0)"
			case "PlusNumbers":
				return "return NewFloatAtom(0.0)"
			case "TimesIntegers":
				return "return NewIntAtom(1)"
			case "TimesReals":
				return "return NewFloatAtom(1.0)"
			case "TimesNumbers":
				return "return NewFloatAtom(1.0)"
			case "AssociationRules":
				return "result := AssociationRules(); return result"
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
				// Special case for Association functions, fallback for others
				return "return args[0]"
			default:
				return "return args[0]"
			}
		},
		"getConversion": func(paramType string) string {
			switch paramType {
			case "int64":
				return `if val, ok := ExtractInt64(arg); ok {
			convertedArgs[i] = val
		} else {
			// Type mismatch - return unchanged
			return CopyExprList(funcName, args)
		}`
			case "float64":
				return `if val, ok := ExtractFloat64(arg); ok {
			convertedArgs[i] = val
		} else {
			// Type mismatch - return unchanged
			return CopyExprList(funcName, args)
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
				return "return NewIntAtom(int(result))"
			case "float64":
				return "return NewFloatAtom(result)"
			case "string":
				return "return NewStringAtom(result)"
			case "bool":
				return "return NewBoolAtom(result)"
			case "Expr":
				return "return result"
			default:
				return "return result"
			}
		},
		"needsFallbackHandling": func(functionName string) bool {
			// Numeric comparison functions need fallback handling for non-numeric types
			fallbackFunctions := []string{"LessExprs", "GreaterExprs", "LessEqualExprs", "GreaterEqualExprs"}
			for _, name := range fallbackFunctions {
				if functionName == name {
					return true
				}
			}
			return false
		},
		"getFixedConversion": func(paramTypes []string) string {
			var conversions []string
			for i, paramType := range paramTypes {
				varName := fmt.Sprintf("arg%d", i)
				switch paramType {
				case "Expr":
					conversions = append(conversions, fmt.Sprintf("	%s := args[%d]", varName, i))
				case "int64":
					conversions = append(conversions, fmt.Sprintf("	%s, ok := ExtractInt64(args[%d])", varName, i))
					conversions = append(conversions, "	if !ok {")
					conversions = append(conversions, "		return CopyExprList(\"FUNC\", args)")
					conversions = append(conversions, "	}")
				case "float64":
					conversions = append(conversions, fmt.Sprintf("	%s, ok := ExtractFloat64(args[%d])", varName, i))
					conversions = append(conversions, "	if !ok {")
					conversions = append(conversions, "		return CopyExprList(\"FUNC\", args)")
					conversions = append(conversions, "	}")
				case "string":
					conversions = append(conversions, fmt.Sprintf("	%s, ok := ExtractString(args[%d])", varName, i))
					conversions = append(conversions, "	if !ok {")
					conversions = append(conversions, "		return CopyExprList(\"FUNC\", args)")
					conversions = append(conversions, "	}")
				case "bool":
					conversions = append(conversions, fmt.Sprintf("	%s, ok := ExtractBool(args[%d])", varName, i))
					conversions = append(conversions, "	if !ok {")
					conversions = append(conversions, "		return CopyExprList(\"FUNC\", args)")
					conversions = append(conversions, "	}")
				case "List":
					conversions = append(conversions, fmt.Sprintf("	%s, ok := args[%d].(List)", varName, i))
					conversions = append(conversions, "	if !ok {")
					conversions = append(conversions, "		return CopyExprList(\"FUNC\", args)")
					conversions = append(conversions, "	}")
				case "ObjectExpr":
					conversions = append(conversions, fmt.Sprintf("	%s, ok := args[%d].(ObjectExpr)", varName, i))
					conversions = append(conversions, "	if !ok {")
					conversions = append(conversions, "		return CopyExprList(\"FUNC\", args)")
					conversions = append(conversions, "	}")
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