package stdlib

import (
	"testing"

	"github.com/client9/sexpr/core"
)

func TestTakeAndDropFunctions(t *testing.T) {
	// Helper function to create test lists
	createList := func(elements ...string) core.List {
		exprs := make([]core.Expr, len(elements)+1)
		exprs[0] = core.NewSymbol("List")
		for i, elem := range elements {
			exprs[i+1] = core.NewSymbol(elem)
		}
		return core.NewListFromExprs(exprs...)
	}

	// Helper function to create integer list spec
	createIntList := func(nums ...int64) core.List {
		exprs := make([]core.Expr, len(nums)+1)
		exprs[0] = core.NewSymbol("List")
		for i, num := range nums {
			exprs[i+1] = core.NewInteger(num)
		}
		return core.NewListFromExprs(exprs...)
	}

	// Helper function to extract list elements as strings (for comparison)
	extractElements := func(expr core.Expr) []string {
		if list, ok := expr.(core.List); ok {
			args := list.Tail()

			result := make([]string, len(args))
			for i, a := range args {
				result[i] = a.String()
			}
			return result
		}
		// For non-list expressions, return as single-element slice
		return []string{expr.String()}
	}

	// Helper function to check if result is an error
	isError := func(expr core.Expr) bool {
		str := expr.String()
		return str == "Error" || (len(str) > 7 && str[:7] == "$Failed")
	}

	testList := createList("a", "b", "c", "d", "e")
	emptyList := createList()

	tests := []struct {
		name      string
		function  string
		list      core.List
		arg       interface{} // int64 for n, core.List for [n] or [n,m]
		expected  []string
		shouldErr bool
	}{
		// TakeList tests (Take(expr, n))
		{"Take first 2", "TakeList", testList, int64(2), []string{"a", "b"}, false},
		{"Take first 3", "TakeList", testList, int64(3), []string{"a", "b", "c"}, false},
		{"Take last 2", "TakeList", testList, int64(-2), []string{"d", "e"}, false},
		{"Take last 3", "TakeList", testList, int64(-3), []string{"c", "d", "e"}, false},
		{"Take 0", "TakeList", testList, int64(0), []string{}, false},
		{"Take more than available", "TakeList", testList, int64(10), []string{"a", "b", "c", "d", "e"}, false},
		{"Take from empty list", "TakeList", emptyList, int64(2), []string{}, false},

		// TakeListSingle tests (Take(expr, [n]))
		{"Take single element 1", "TakeListSingle", testList, createIntList(1), []string{"a"}, false},
		{"Take single element 3", "TakeListSingle", testList, createIntList(3), []string{"c"}, false},
		{"Take single element -1", "TakeListSingle", testList, createIntList(-1), []string{"e"}, false},
		{"Take single element -2", "TakeListSingle", testList, createIntList(-2), []string{"d"}, false},
		{"Take single out of bounds", "TakeListSingle", testList, createIntList(10), nil, true},
		{"Take single from empty", "TakeListSingle", emptyList, createIntList(1), nil, true},

		// TakeListRange tests (Take(expr, [n, m]))
		{"Take range [1,3]", "TakeListRange", testList, createIntList(1, 3), []string{"a", "b", "c"}, false},
		{"Take range [2,4]", "TakeListRange", testList, createIntList(2, 4), []string{"b", "c", "d"}, false},
		{"Take range [-3,-1]", "TakeListRange", testList, createIntList(-3, -1), []string{"c", "d", "e"}, false},
		{"Take range [-2,-2]", "TakeListRange", testList, createIntList(-2, -2), []string{"d"}, false},
		{"Take range out of bounds", "TakeListRange", testList, createIntList(1, 10), nil, true},
		{"Take range invalid order", "TakeListRange", testList, createIntList(3, 1), nil, true},

		// DropList tests (Drop(expr, n))
		{"Drop first 2", "DropList", testList, int64(2), []string{"c", "d", "e"}, false},
		{"Drop first 3", "DropList", testList, int64(3), []string{"d", "e"}, false},
		{"Drop last 2", "DropList", testList, int64(-2), []string{"a", "b", "c"}, false},
		{"Drop last 3", "DropList", testList, int64(-3), []string{"a", "b"}, false},
		{"Drop 0", "DropList", testList, int64(0), []string{"a", "b", "c", "d", "e"}, false},
		{"Drop all", "DropList", testList, int64(5), []string{}, false},
		{"Drop more than available", "DropList", testList, int64(10), []string{}, false},
		{"Drop from empty list", "DropList", emptyList, int64(2), []string{}, false},

		// DropListSingle tests (Drop(expr, [n]))
		{"Drop single element 1", "DropListSingle", testList, createIntList(1), []string{"b", "c", "d", "e"}, false},
		{"Drop single element 3", "DropListSingle", testList, createIntList(3), []string{"a", "b", "d", "e"}, false},
		{"Drop single element -1", "DropListSingle", testList, createIntList(-1), []string{"a", "b", "c", "d"}, false},
		{"Drop single element -2", "DropListSingle", testList, createIntList(-2), []string{"a", "b", "c", "e"}, false},
		{"Drop single out of bounds", "DropListSingle", testList, createIntList(10), nil, true},
		{"Drop single from empty", "DropListSingle", emptyList, createIntList(1), []string{}, false},

		// DropListRange tests (Drop(expr, [n, m]))
		{"Drop range [1,3]", "DropListRange", testList, createIntList(1, 3), []string{"d", "e"}, false},
		{"Drop range [2,4]", "DropListRange", testList, createIntList(2, 4), []string{"a", "e"}, false},
		{"Drop range [-3,-1]", "DropListRange", testList, createIntList(-3, -1), []string{"a", "b"}, false},
		{"Drop range [-2,-2]", "DropListRange", testList, createIntList(-2, -2), []string{"a", "b", "c", "e"}, false},
		{"Drop range out of bounds", "DropListRange", testList, createIntList(1, 10), nil, true},
		{"Drop range invalid order", "DropListRange", testList, createIntList(3, 1), nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result core.Expr

			switch tt.function {
			case "TakeList":
				result = TakeList(tt.list, tt.arg.(int64))
			case "TakeListSingle":
				result = TakeListSingle(tt.list, tt.arg.(core.List))
			case "TakeListRange":
				result = TakeListRange(tt.list, tt.arg.(core.List))
			case "DropList":
				result = DropList(tt.list, tt.arg.(int64))
			case "DropListSingle":
				result = DropListSingle(tt.list, tt.arg.(core.List))
			case "DropListRange":
				result = DropListRange(tt.list, tt.arg.(core.List))
			default:
				t.Fatalf("Unknown function: %s", tt.function)
			}

			if tt.shouldErr {
				if !isError(result) {
					t.Errorf("Expected error, got: %v", result)
				}
			} else {
				if isError(result) {
					t.Errorf("Unexpected error: %v", result)
				} else {
					actual := extractElements(result)
					if len(actual) != len(tt.expected) {
						t.Errorf("Expected %d elements, got %d: %v vs %v",
							len(tt.expected), len(actual), tt.expected, actual)
					} else {
						for i, exp := range tt.expected {
							if i >= len(actual) || actual[i] != exp {
								t.Errorf("Expected %v, got %v", tt.expected, actual)
								break
							}
						}
					}
				}
			}
		})
	}
}

func TestTakeDropEdgeCases(t *testing.T) {
	// Helper function to create test lists
	createList := func(elements ...string) core.List {
		exprs := make([]core.Expr, len(elements)+1)
		exprs[0] = core.NewSymbol("List")
		for i, elem := range elements {
			exprs[i+1] = core.NewSymbol(elem)
		}
		return core.NewListFromExprs(exprs...)
	}

	// Helper function to create integer list spec
	createIntList := func(nums ...int64) core.List {
		exprs := make([]core.Expr, len(nums)+1)
		exprs[0] = core.NewSymbol("List")
		for i, num := range nums {
			exprs[i+1] = core.NewInteger(num)
		}
		return core.NewListFromExprs(exprs...)
	}

	singleElementList := createList("x")

	// Test edge cases
	t.Run("Take from single element list", func(t *testing.T) {
		result := TakeList(singleElementList, 1)
		if list, ok := result.(core.List); ok {
			args := list.AsSlice()
			if len(args) != 2 || args[1].String() != "x" {
				t.Errorf("Expected List(x), got %v", result)
			}
		} else {
			t.Errorf("Expected List, got %T", result)
		}
	})

	t.Run("Drop from single element list", func(t *testing.T) {
		result := DropList(singleElementList, 1)
		if list, ok := result.(core.List); ok {
			args := list.AsSlice()
			if len(args) != 1 { // Just the head
				t.Errorf("Expected empty list, got %v", result)
			}
		} else {
			t.Errorf("Expected List, got %T", result)
		}
	})

	t.Run("Take single with invalid spec", func(t *testing.T) {
		// Test with too many arguments in list spec
		invalidSpec := createIntList(1, 2, 3)
		result := TakeListSingle(createList("a", "b", "c"), invalidSpec)
		str := result.String()
		if !(str == "Error" || (len(str) > 7 && str[:7] == "$Failed")) {
			t.Errorf("Expected error for invalid spec, got %v", result)
		}
	})

	t.Run("Drop range with zero index", func(t *testing.T) {
		// Test with zero index (invalid in 1-based indexing)
		zeroSpec := createIntList(0, 2)
		result := DropListRange(createList("a", "b", "c"), zeroSpec)
		str := result.String()
		if !(str == "Error" || (len(str) > 7 && str[:7] == "$Failed")) {
			t.Errorf("Expected error for zero index, got %v", result)
		}
	})
}

func TestRotateFunctions(t *testing.T) {
	// Helper function to create test lists
	createList := func(elements ...int64) core.List {
		exprs := make([]core.Expr, len(elements)+1)
		exprs[0] = core.NewSymbol("List")
		for i, elem := range elements {
			exprs[i+1] = core.NewInteger(elem)
		}
		return core.NewListFromExprs(exprs...)
	}

	t.Run("RotateLeft basic tests", func(t *testing.T) {
		tests := []struct {
			name     string
			input    core.List
			n        int64
			expected string
		}{
			{
				name:     "rotate left by 1",
				input:    createList(1, 2, 3, 4, 5),
				n:        1,
				expected: "List(2, 3, 4, 5, 1)",
			},
			{
				name:     "rotate left by 2",
				input:    createList(1, 2, 3, 4, 5),
				n:        2,
				expected: "List(3, 4, 5, 1, 2)",
			},
			{
				name:     "rotate left by 0",
				input:    createList(1, 2, 3),
				n:        0,
				expected: "List(1, 2, 3)",
			},
			{
				name:     "rotate left by negative (equivalent to right)",
				input:    createList(1, 2, 3, 4, 5),
				n:        -1,
				expected: "List(5, 1, 2, 3, 4)",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := RotateLeft(tt.input, tt.n)
				if result.String() != tt.expected {
					t.Errorf("expected %s, got %s", tt.expected, result.String())
				}
			})
		}
	})

	t.Run("RotateRight basic tests", func(t *testing.T) {
		tests := []struct {
			name     string
			input    core.List
			n        int64
			expected string
		}{
			{
				name:     "rotate right by 1",
				input:    createList(1, 2, 3, 4, 5),
				n:        1,
				expected: "List(5, 1, 2, 3, 4)",
			},
			{
				name:     "rotate right by 2",
				input:    createList(1, 2, 3, 4, 5),
				n:        2,
				expected: "List(4, 5, 1, 2, 3)",
			},
			{
				name:     "rotate right by 0",
				input:    createList(1, 2, 3),
				n:        0,
				expected: "List(1, 2, 3)",
			},
			{
				name:     "rotate right by negative (equivalent to left)",
				input:    createList(1, 2, 3, 4, 5),
				n:        -1,
				expected: "List(2, 3, 4, 5, 1)",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := RotateRight(tt.input, tt.n)
				if result.String() != tt.expected {
					t.Errorf("expected %s, got %s", tt.expected, result.String())
				}
			})
		}
	})
}

func TestFlattenExpr(t *testing.T) {
	// Helper function to create test lists with mixed elements
	createMixedList := func(head string, elements ...core.Expr) core.List {
		exprs := make([]core.Expr, len(elements)+1)
		exprs[0] = core.NewSymbol(head)
		copy(exprs[1:], elements)
		return core.NewListFromExprs(exprs...)
	}

	// Helper function to create simple integer lists
	createIntList := func(head string, nums ...int64) core.List {
		exprs := make([]core.Expr, len(nums)+1)
		exprs[0] = core.NewSymbol(head)
		for i, num := range nums {
			exprs[i+1] = core.NewInteger(num)
		}
		return core.NewListFromExprs(exprs...)
	}

	tests := []struct {
		name     string
		input    core.Expr
		expected string
	}{
		// Basic flattening tests
		{
			name:     "Simple nested list",
			input:    createMixedList("List", core.NewInteger(1), core.NewInteger(2), createIntList("List", 3, 4)),
			expected: "List(1, 2, 3, 4)",
		},
		{
			name:     "Deeply nested list",
			input:    createMixedList("List", core.NewInteger(1), createMixedList("List", core.NewInteger(2), createIntList("List", 3, 4), core.NewInteger(5)), core.NewInteger(6)),
			expected: "List(1, 2, 3, 4, 5, 6)",
		},
		{
			name:     "Empty nested lists",
			input:    createMixedList("List", core.NewInteger(1), createIntList("List"), core.NewInteger(2)),
			expected: "List(1, 2)",
		},
		{
			name:     "Multiple nested lists",
			input:    createMixedList("List", createIntList("List", 1, 2), createIntList("List", 3, 4), createIntList("List", 5, 6)),
			expected: "List(1, 2, 3, 4, 5, 6)",
		},

		// Different head types
		{
			name:     "Zoo with nested Zoo",
			input:    createMixedList("Zoo", core.NewInteger(1), core.NewInteger(2), createIntList("Zoo", 3, 4)),
			expected: "Zoo(1, 2, 3, 4)",
		},
		{
			name:     "Mixed heads - Zoo with List inside",
			input:    createMixedList("Zoo", core.NewInteger(1), createIntList("List", 2, 3), core.NewInteger(4)),
			expected: "Zoo(1, List(2, 3), 4)",
		},
		{
			name:     "Mixed heads - List with Zoo inside",
			input:    createMixedList("List", core.NewInteger(1), createIntList("Zoo", 2, 3), core.NewInteger(4)),
			expected: "List(1, Zoo(2, 3), 4)",
		},

		// Edge cases
		{
			name:     "Empty list",
			input:    createIntList("List"),
			expected: "List()",
		},
		{
			name:     "Single element list",
			input:    createIntList("List", 42),
			expected: "List(42)",
		},
		{
			name:     "Already flat list",
			input:    createIntList("List", 1, 2, 3, 4),
			expected: "List(1, 2, 3, 4)",
		},
		{
			name:     "Nested with symbols and numbers",
			input:    createMixedList("List", core.NewSymbol("a"), createMixedList("List", core.NewInteger(1), core.NewSymbol("b")), core.NewInteger(2)),
			expected: "List(a, 1, b, 2)",
		},

		// Non-list inputs
		{
			name:     "Integer input",
			input:    core.NewInteger(42),
			expected: "42",
		},
		{
			name:     "Symbol input",
			input:    core.NewSymbol("x"),
			expected: "x",
		},
		{
			name:     "String input",
			input:    core.NewString("hello"),
			expected: "\"hello\"",
		},

		// Complex nesting
		{
			name: "Triple nesting",
			input: createMixedList("List",
				core.NewInteger(1),
				createMixedList("List",
					core.NewInteger(2),
					createMixedList("List",
						core.NewInteger(3),
						createIntList("List", 4, 5),
					),
					core.NewInteger(6),
				),
				core.NewInteger(7),
			),
			expected: "List(1, 2, 3, 4, 5, 6, 7)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FlattenExpr(tt.input)
			if result.String() != tt.expected {
				t.Errorf("FlattenExpr(%s) = %s, expected %s",
					tt.input.String(), result.String(), tt.expected)
			}
		})
	}
}
