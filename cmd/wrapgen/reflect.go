package main

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

// EngineFunc is the universal signature for all registered functions
// All functions in the system will eventually use this signature
type EngineFunc func(e interface{}, args []interface{}) interface{}

// ReflectionInfo contains analyzed information about a function signature
type ReflectionInfo struct {
	Name           string   // Full function name from runtime
	IsVariadic     bool     // true if function takes variadic args
	ParamTypes     []string // Parameter types as strings
	ReturnType     string   // Primary return type as string
	ReturnsError   bool     // true if second return type is error
	NeedsEvaluator bool     // true if first parameter is *Evaluator
	DirectCall     bool     // true if no conversion needed
}

// FunctionSpec is used for reflection analysis (legacy compatibility)
type FunctionSpec struct {
	Pattern        string   // "Plus(x__Integer)" - the full pattern
	Function       any      // Function reference for reflection
	FunctionName   string   // "PlusIntegers" - derived from Function name
	WrapperName    string   // "WrapPlusIntegers" - derived from FunctionName
	IsVariadic     bool     // derived from Function signature
	ParamType      string   // For variadic: derived from Function signature
	ParamTypes     []string // For fixed arity: derived from Function signature
	ReturnType     string   // derived from Function signature
	ReturnsError   bool     // derived from Function signature (has error return)
	NeedsEvaluator bool     // derived from Function signature (first param is *Evaluator)
	DirectCall     bool     // true if no wrapper needed
}

// analyzeFunctionSignature uses reflection to analyze a function's signature
func analyzeFunctionSignature(fn any) (ReflectionInfo, error) {
	if fn == nil {
		return ReflectionInfo{}, fmt.Errorf("function is nil")
	}

	t := reflect.TypeOf(fn)

	if t.Kind() != reflect.Func {
		return ReflectionInfo{}, fmt.Errorf("expected function, got %s", t.Kind())
	}

	// Get function name from runtime
	fullName := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()

	// Extract parameter types and check for evaluator dependency
	paramTypes := make([]string, 0, t.NumIn())
	needsEvaluator := false
	directCall := isEngineFunc(t)

	if t.NumIn() > 0 && isEvaluatorType(t.In(0)) {
		needsEvaluator = true
	}

	for i := 0; i < t.NumIn(); i++ {
		paramType := t.In(i)
		paramTypeStr := typeToString(paramType)
		paramTypes = append(paramTypes, paramTypeStr)
	}

	// Extract return type
	var returnType string
	var returnsError bool

	if t.NumOut() > 0 {
		returnType = typeToString(t.Out(0))

		// Check if second return type is error
		if t.NumOut() == 2 {
			errorType := reflect.TypeOf((*error)(nil)).Elem()
			if t.Out(1) == errorType {
				returnsError = true
			}
		}
	}

	return ReflectionInfo{
		Name:           fullName,
		IsVariadic:     t.IsVariadic(),
		ParamTypes:     paramTypes,
		ReturnType:     returnType,
		ReturnsError:   returnsError,
		NeedsEvaluator: needsEvaluator,
		DirectCall:     directCall,
	}, nil
}

// typeToString converts a reflect.Type to a string representation suitable for code generation
func typeToString(t reflect.Type) string {
	// Handle variadic slice types
	switch t.Kind() {
	case reflect.Slice:
		// Handle variadic slice types like []int64 -> int64
		if t.Elem().Kind() == reflect.Int64 {
			return "int64"
		} else if t.Elem().Kind() == reflect.Float64 {
			return "float64"
		}
		if t.Name() == "List" {
			return "List"
		}
		//	fmt.Printf("reflect: %v\n", t.Name())
		// For other slices, return the element type
		return typeToString(t.Elem())
	}
	//fmt.Printf("reflect2: %v\n", t)
	// Use type name if available
	if t.Name() != "" {
		return t.Name()
	}

	// Fall back to string representation
	return t.String()
}

// extractFunctionName extracts just the function name from a full runtime name
// e.g., "github.com/client9/sexpr/stdlib.PlusIntegers" -> "PlusIntegers"
func extractFunctionName(fullName string) string {
	parts := strings.Split(fullName, ".")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return fullName
}

// fillFromReflection populates auto-derived fields in FunctionSpec using reflection
func (spec *FunctionSpec) fillFromReflection() error {
	if spec.Function == nil {
		// No function reference provided, skip reflection
		return nil
	}

	info, err := analyzeFunctionSignature(spec.Function)
	if err != nil {
		return fmt.Errorf("failed to analyze function signature: %v", err)
	}

	// Extract and populate auto-derived fields
	functionName := extractFunctionName(info.Name)
	spec.FunctionName = functionName
	spec.WrapperName = "Wrap" + functionName
	spec.IsVariadic = info.IsVariadic
	spec.ReturnType = info.ReturnType
	spec.ReturnsError = info.ReturnsError
	spec.NeedsEvaluator = info.NeedsEvaluator
	spec.DirectCall = info.DirectCall

	// If function needs evaluator, skip the two parameter in conversion
	paramTypes := info.ParamTypes
	if info.NeedsEvaluator && len(paramTypes) > 0 {
		paramTypes = paramTypes[2:] // Skip (*engine.Evaluator, *engine.Context)
	}

	if info.IsVariadic && len(paramTypes) > 0 {
		// For variadic functions, ParamType is the variadic parameter type
		spec.ParamType = paramTypes[0]
		spec.ParamTypes = nil // Clear fixed param types
	} else {
		// For fixed-arity functions, use ParamTypes (excluding evaluator)
		spec.ParamTypes = paramTypes
		spec.ParamType = "" // Clear variadic param type
	}

	return nil
}

// is standard engine function, so no wrapper is needed
// TODO: instead of hard coding this we could look at the actually interface definition
func isEngineFunc(t reflect.Type) bool {
	if t.NumIn() != 3 {
		return false
	}
	if !isEvaluatorType(t.In(0)) {
		return false
	}
	if !isContextType(t.In(1)) {
		return false
	}
	if !isSliceExpr(t.In(2)) {
		return false
	}
	return true
}

// isEvaluatorType checks if a reflect.Type represents *Evaluator
func isEvaluatorType(t reflect.Type) bool {
	// Check for *Evaluator (pointer to struct)
	if t.Kind() == reflect.Ptr {
		elem := t.Elem()
		if elem.Kind() == reflect.Struct {
			typeName := elem.Name()
			return typeName == "Evaluator"
		}
	}
	return false
}

// isEvaluatorType checks if a reflect.Type represents *engine.Context
func isContextType(t reflect.Type) bool {
	// Check for *Evaluator (pointer to struct)
	if t.Kind() != reflect.Ptr {
		return false
	}
	elem := t.Elem()
	if elem.Kind() != reflect.Struct {
		return false
	}
	return elem.Name() == "Context"
}

func isSliceExpr(t reflect.Type) bool {
	if t.Kind() != reflect.Slice {
		return false
	}
	elem := t.Elem()
	if elem.Kind() != reflect.Interface {
		return false
	}
	return elem.Name() == "Expr"
}
