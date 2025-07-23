package main

import (
	"fmt"
	//	"log"
	"reflect"
	"runtime"
	"strings"
)

// FunctionInfo contains analyzed information about a function signature
type FunctionInfo struct {
	Name         string   // Full function name from runtime
	IsVariadic   bool     // true if function takes variadic args
	ParamTypes   []string // Parameter types as strings
	ReturnType   string   // Primary return type as string
	ReturnsError bool     // true if second return type is error
}

// analyzeFunctionSignature uses reflection to analyze a function's signature
func analyzeFunctionSignature(fn interface{}) (FunctionInfo, error) {
	if fn == nil {
		return FunctionInfo{}, fmt.Errorf("function is nil")
	}

	t := reflect.TypeOf(fn)
	if t.Kind() != reflect.Func {
		return FunctionInfo{}, fmt.Errorf("expected function, got %s", t.Kind())
	}

	// Get function name from runtime
	fullName := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()

	// Extract parameter types
	paramTypes := make([]string, 0, t.NumIn())
	for i := 0; i < t.NumIn(); i++ {
		paramType := t.In(i)
		paramTypes = append(paramTypes, typeToString(paramType))
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

	return FunctionInfo{
		Name:         fullName,
		IsVariadic:   t.IsVariadic(),
		ParamTypes:   paramTypes,
		ReturnType:   returnType,
		ReturnsError: returnsError,
	}, nil
}

// typeToString converts a reflect.Type to a string representation suitable for code generation
func typeToString(t reflect.Type) string {
	// this switch is a problem
	switch t.Kind() {
	case reflect.Slice:
		// Handle variadic slice types like []int64 -> int64
		if t.Elem().Kind() == reflect.Int64 {
			return "int64"
		} else if t.Elem().Kind() == reflect.Float64 {
			return "float64"
		}
		// For other slices, return the element type
		return typeToString(t.Elem())
	}

	// Don't use "t.Kind()" since it removes types.
	//	switch t.Kind() {

	if t.Name() != "" {
		return t.Name()
	}

	// unclear how it falls to here
	return t.String()
	/*
	   switch t.Name() {
	   case "Number":

	   	return "Number"

	   case "string":

	   	return "string"

	   case "int64":

	   	return "int64"

	   case "float64":

	   	return "float64"

	   case "bool":

	   	return "bool"

	   case "":

	   	return t.String()

	   default:

	   		return t.Name()
	   	}
	*/
}

// extractFunctionName extracts just the function name from a full runtime name
// e.g., "github.com/client9/sexpr.PlusIntegers" -> "PlusIntegers"
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

	if info.IsVariadic && len(info.ParamTypes) > 0 {
		// For variadic functions, ParamType is the variadic parameter type
		spec.ParamType = info.ParamTypes[0]
		spec.ParamTypes = nil // Clear fixed param types
	} else {
		// For fixed-arity functions, use ParamTypes
		spec.ParamTypes = info.ParamTypes
		spec.ParamType = "" // Clear variadic param type
	}

	return nil
}
