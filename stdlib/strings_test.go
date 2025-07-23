package stdlib

import (
	"testing"

	"github.com/client9/sexpr/core"
)

func TestTakeAndDropStringFunctions(t *testing.T) {
	// Helper function to create integer list spec
	createIntList := func(nums ...int64) core.List {
		exprs := make([]core.Expr, len(nums)+1)
		exprs[0] = core.NewSymbolAtom("List")
		for i, num := range nums {
			exprs[i+1] = core.NewIntAtom(int(num))
		}
		return core.List{Elements: exprs}
	}

	// Helper function to extract string from result
	extractString := func(expr core.Expr) string {
		if str, ok := core.ExtractString(expr); ok {
			return str
		}
		return ""
	}

	// Helper function to check if result is an error
	isError := func(expr core.Expr) bool {
		str := expr.String()
		return str == "Error" || (len(str) > 7 && str[:7] == "$Failed")
	}

	tests := []struct {
		name      string
		function  string
		str       string
		arg       interface{} // int64 for n, core.List for [n] or [n,m]
		expected  string
		shouldErr bool
	}{
		// StringTakeByte tests (StringTakeByte(str, n))
		{"Take first 3 bytes", "StringTakeByte", "hello", int64(3), "hel", false},
		{"Take first 5 bytes", "StringTakeByte", "hello", int64(5), "hello", false},
		{"Take last 2 bytes", "StringTakeByte", "hello", int64(-2), "lo", false},
		{"Take last 3 bytes", "StringTakeByte", "hello", int64(-3), "llo", false},
		{"Take 0 bytes", "StringTakeByte", "hello", int64(0), "", false},
		{"Take more than available", "StringTakeByte", "hello", int64(10), "hello", false},
		{"Take from empty string", "StringTakeByte", "", int64(2), "", false},

		// StringTakeByteAt tests (StringTakeByteAt(str, [n]))
		{"Take single byte 1", "StringTakeByteAt", "hello", createIntList(1), "h", false},
		{"Take single byte 3", "StringTakeByteAt", "hello", createIntList(3), "l", false},
		{"Take single byte -1", "StringTakeByteAt", "hello", createIntList(-1), "o", false},
		{"Take single byte -2", "StringTakeByteAt", "hello", createIntList(-2), "l", false},
		{"Take single out of bounds", "StringTakeByteAt", "hello", createIntList(10), "", true},
		{"Take single from empty", "StringTakeByteAt", "", createIntList(1), "", true},

		// StringTakeByteRange tests (StringTakeByteRange(str, [n, m]))
		{"Take range [1,3]", "StringTakeByteRange", "hello", createIntList(1, 3), "hel", false},
		{"Take range [2,4]", "StringTakeByteRange", "hello", createIntList(2, 4), "ell", false},
		{"Take range [-3,-1]", "StringTakeByteRange", "hello", createIntList(-3, -1), "llo", false},
		{"Take range [-2,-2]", "StringTakeByteRange", "hello", createIntList(-2, -2), "l", false},
		{"Take range out of bounds", "StringTakeByteRange", "hello", createIntList(1, 10), "", true},
		{"Take range invalid order", "StringTakeByteRange", "hello", createIntList(3, 1), "", true},

		// StringDropByte tests (StringDropByte(str, n))
		{"Drop first 2 bytes", "StringDropByte", "hello", int64(2), "llo", false},
		{"Drop first 3 bytes", "StringDropByte", "hello", int64(3), "lo", false},
		{"Drop last 2 bytes", "StringDropByte", "hello", int64(-2), "hel", false},
		{"Drop last 3 bytes", "StringDropByte", "hello", int64(-3), "he", false},
		{"Drop 0 bytes", "StringDropByte", "hello", int64(0), "hello", false},
		{"Drop all bytes", "StringDropByte", "hello", int64(5), "", false},
		{"Drop more than available", "StringDropByte", "hello", int64(10), "", false},
		{"Drop from empty string", "StringDropByte", "", int64(2), "", false},

		// StringDropByteAt tests (StringDropByteAt(str, [n]))
		{"Drop single byte 1", "StringDropByteAt", "hello", createIntList(1), "ello", false},
		{"Drop single byte 3", "StringDropByteAt", "hello", createIntList(3), "helo", false},
		{"Drop single byte -1", "StringDropByteAt", "hello", createIntList(-1), "hell", false},
		{"Drop single byte -2", "StringDropByteAt", "hello", createIntList(-2), "helo", false},
		{"Drop single out of bounds", "StringDropByteAt", "hello", createIntList(10), "", true},
		{"Drop single from empty", "StringDropByteAt", "", createIntList(1), "", false},

		// StringDropByteRange tests (StringDropByteRange(str, [n, m]))
		{"Drop range [1,3]", "StringDropByteRange", "hello", createIntList(1, 3), "lo", false},
		{"Drop range [2,4]", "StringDropByteRange", "hello", createIntList(2, 4), "ho", false},
		{"Drop range [-3,-1]", "StringDropByteRange", "hello", createIntList(-3, -1), "he", false},
		{"Drop range [-2,-2]", "StringDropByteRange", "hello", createIntList(-2, -2), "helo", false},
		{"Drop range out of bounds", "StringDropByteRange", "hello", createIntList(1, 10), "", true},
		{"Drop range invalid order", "StringDropByteRange", "hello", createIntList(3, 1), "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result core.Expr

			switch tt.function {
			case "StringTakeByte":
				result = StringTakeByte(tt.str, tt.arg.(int64))
			case "StringTakeByteAt":
				result = StringTakeByteAt(tt.str, tt.arg.(core.List))
			case "StringTakeByteRange":
				result = StringTakeByteRange(tt.str, tt.arg.(core.List))
			case "StringDropByte":
				result = StringDropByte(tt.str, tt.arg.(int64))
			case "StringDropByteAt":
				result = StringDropByteAt(tt.str, tt.arg.(core.List))
			case "StringDropByteRange":
				result = StringDropByteRange(tt.str, tt.arg.(core.List))
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
					actual := extractString(result)
					if actual != tt.expected {
						t.Errorf("Expected %q, got %q", tt.expected, actual)
					}
				}
			}
		})
	}
}

func TestStringTakeDropEdgeCases(t *testing.T) {
	// Helper function to create integer list spec
	createIntList := func(nums ...int64) core.List {
		exprs := make([]core.Expr, len(nums)+1)
		exprs[0] = core.NewSymbolAtom("List")
		for i, num := range nums {
			exprs[i+1] = core.NewIntAtom(int(num))
		}
		return core.List{Elements: exprs}
	}

	t.Run("Take from single byte string", func(t *testing.T) {
		result := StringTakeByte("x", 1)
		if str, ok := core.ExtractString(result); !ok || str != "x" {
			t.Errorf("Expected \"x\", got %v", result)
		}
	})

	t.Run("Drop from single byte string", func(t *testing.T) {
		result := StringDropByte("x", 1)
		if str, ok := core.ExtractString(result); !ok || str != "" {
			t.Errorf("Expected empty string, got %v", result)
		}
	})

	t.Run("Take single with invalid spec", func(t *testing.T) {
		// Test with too many arguments in list spec
		invalidSpec := createIntList(1, 2, 3)
		result := StringTakeByteAt("abc", invalidSpec)
		str := result.String()
		if !(str == "Error" || (len(str) > 7 && str[:7] == "$Failed")) {
			t.Errorf("Expected error for invalid spec, got %v", result)
		}
	})

	t.Run("Drop range with zero index", func(t *testing.T) {
		// Test with zero index (invalid in 1-based indexing)
		zeroSpec := createIntList(0, 2)
		result := StringDropByteRange("abc", zeroSpec)
		str := result.String()
		if !(str == "Error" || (len(str) > 7 && str[:7] == "$Failed")) {
			t.Errorf("Expected error for zero index, got %v", result)
		}
	})

	t.Run("Test with multi-byte UTF-8 characters", func(t *testing.T) {
		// Test with string containing multi-byte UTF-8 characters
		// "café" = 4 Unicode characters but 5 bytes (é is 2 bytes: 0xC3 0xA9)
		utf8Str := "café"

		// Should operate on bytes, not Unicode characters
		// Taking 4 bytes gives us "caf" + first byte of é (0xC3)
		result := StringTakeByte(utf8Str, 4)
		if str, ok := core.ExtractString(result); !ok || str != "caf\xc3" {
			t.Errorf("Expected \"caf\\xc3\" (4 bytes from café), got %q", str)
		}

		// Taking 5 bytes should get the whole string
		result = StringTakeByte(utf8Str, 5)
		if str, ok := core.ExtractString(result); !ok || str != "café" {
			t.Errorf("Expected \"café\" (all 5 bytes), got %q", str)
		}

		// Taking first 3 bytes should give us just "caf"
		result = StringTakeByte(utf8Str, 3)
		if str, ok := core.ExtractString(result); !ok || str != "caf" {
			t.Errorf("Expected \"caf\" (3 bytes), got %q", str)
		}
	})
}

func TestStringTakeAndStringDropFunctions(t *testing.T) {
	// Helper function to create integer list spec
	createIntList := func(nums ...int64) core.List {
		exprs := make([]core.Expr, len(nums)+1)
		exprs[0] = core.NewSymbolAtom("List")
		for i, num := range nums {
			exprs[i+1] = core.NewIntAtom(int(num))
		}
		return core.List{Elements: exprs}
	}

	// Helper function to extract string from result
	extractString := func(expr core.Expr) string {
		if str, ok := core.ExtractString(expr); ok {
			return str
		}
		return ""
	}

	// Helper function to check if result is an error
	isError := func(expr core.Expr) bool {
		str := expr.String()
		return str == "Error" || (len(str) > 7 && str[:7] == "$Failed")
	}

	tests := []struct {
		name      string
		function  string
		str       string
		arg       interface{} // int64 for n, core.List for [n] or [n,m]
		expected  string
		shouldErr bool
	}{
		// StringTake tests (StringTake(str, n)) - RUNE-based
		{"StringTake first 3 runes", "StringTake", "café", int64(3), "caf", false},
		{"StringTake first 4 runes", "StringTake", "café", int64(4), "café", false},
		{"StringTake last 2 runes", "StringTake", "héllo", int64(-2), "lo", false},
		{"StringTake last 3 runes", "StringTake", "héllo", int64(-3), "llo", false},
		{"StringTake 0 runes", "StringTake", "café", int64(0), "", false},
		{"StringTake more than available", "StringTake", "hi", int64(10), "hi", false},
		{"StringTake from empty string", "StringTake", "", int64(2), "", false},

		// StringTakeAt tests (StringTakeAt(str, [n])) - RUNE-based
		{"StringTake single rune 1", "StringTakeAt", "café", createIntList(1), "c", false},
		{"StringTake single rune 4", "StringTakeAt", "café", createIntList(4), "é", false},
		{"StringTake single rune -1", "StringTakeAt", "café", createIntList(-1), "é", false},
		{"StringTake single rune -2", "StringTakeAt", "café", createIntList(-2), "f", false},
		{"StringTake single out of bounds", "StringTakeAt", "hi", createIntList(10), "", true},
		{"StringTake single from empty", "StringTakeAt", "", createIntList(1), "", true},

		// StringTakeRange tests (StringTake(str, [n, m])) - RUNE-based
		{"StringTake range [1,3]", "StringTakeRange", "café", createIntList(1, 3), "caf", false},
		{"StringTake range [2,4]", "StringTakeRange", "café", createIntList(2, 4), "afé", false},
		{"StringTake range [-3,-1]", "StringTakeRange", "héllo", createIntList(-3, -1), "llo", false},
		{"StringTake range [-2,-2]", "StringTakeRange", "café", createIntList(-2, -2), "f", false},
		{"StringTake range out of bounds", "StringTakeRange", "hi", createIntList(1, 10), "", true},
		{"StringTake range invalid order", "StringTakeRange", "café", createIntList(3, 1), "", true},

		// StringDrop tests (StringDrop(str, n)) - RUNE-based
		{"StringDrop first 2 runes", "StringDrop", "café", int64(2), "fé", false},
		{"StringDrop first 3 runes", "StringDrop", "héllo", int64(3), "lo", false},
		{"StringDrop last 2 runes", "StringDrop", "café", int64(-2), "ca", false},
		{"StringDrop last 3 runes", "StringDrop", "héllo", int64(-3), "hé", false},
		{"StringDrop 0 runes", "StringDrop", "café", int64(0), "café", false},
		{"StringDrop all runes", "StringDrop", "hi", int64(2), "", false},
		{"StringDrop more than available", "StringDrop", "hi", int64(10), "", false},
		{"StringDrop from empty string", "StringDrop", "", int64(2), "", false},

		// StringDropAt tests (StringDropAt(str, [n])) - RUNE-based
		{"StringDrop single rune 1", "StringDropAt", "café", createIntList(1), "afé", false},
		{"StringDrop single rune 4", "StringDropAt", "café", createIntList(4), "caf", false},
		{"StringDrop single rune -1", "StringDropAt", "café", createIntList(-1), "caf", false},
		{"StringDrop single rune -2", "StringDropAt", "café", createIntList(-2), "caé", false},
		{"StringDrop single out of bounds", "StringDropAt", "hi", createIntList(10), "", true},
		{"StringDrop single from empty", "StringDropAt", "", createIntList(1), "", false},

		// StringDropRange tests (StringDrop(str, [n, m])) - RUNE-based
		{"StringDrop range [1,3]", "StringDropRange", "café", createIntList(1, 3), "é", false},
		{"StringDrop range [2,4]", "StringDropRange", "café", createIntList(2, 4), "c", false},
		{"StringDrop range [-3,-1]", "StringDropRange", "héllo", createIntList(-3, -1), "hé", false},
		{"StringDrop range [-2,-2]", "StringDropRange", "café", createIntList(-2, -2), "caé", false},
		{"StringDrop range out of bounds", "StringDropRange", "hi", createIntList(1, 10), "", true},
		{"StringDrop range invalid order", "StringDropRange", "café", createIntList(3, 1), "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result core.Expr

			switch tt.function {
			case "StringTake":
				result = StringTake(tt.str, tt.arg.(int64))
			case "StringTakeAt":
				result = StringTakeAt(tt.str, tt.arg.(core.List))
			case "StringTakeRange":
				result = StringTakeRange(tt.str, tt.arg.(core.List))
			case "StringDrop":
				result = StringDrop(tt.str, tt.arg.(int64))
			case "StringDropAt":
				result = StringDropAt(tt.str, tt.arg.(core.List))
			case "StringDropRange":
				result = StringDropRange(tt.str, tt.arg.(core.List))
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
					actual := extractString(result)
					if actual != tt.expected {
						t.Errorf("Expected %q, got %q", tt.expected, actual)
					}
				}
			}
		})
	}
}

func TestStringTakeDropVsByteBasedComparison(t *testing.T) {
	// Test the difference between rune-based and byte-based operations
	utf8Str := "café" // 4 characters, 5 bytes (é = 2 bytes)

	t.Run("Compare StringTake vs Take on UTF-8", func(t *testing.T) {
		// StringTake works on runes
		result1 := StringTake(utf8Str, 4)
		if str, ok := core.ExtractString(result1); !ok || str != "café" {
			t.Errorf("StringTake(\"café\", 4) expected \"café\", got %q", str)
		}

		// Take works on bytes
		result2 := StringTakeByte(utf8Str, 4)
		if str, ok := core.ExtractString(result2); !ok || str != "caf\xc3" {
			t.Errorf("Take(\"café\", 4) expected \"caf\\xc3\", got %q", str)
		}
	})

	t.Run("Compare StringTake single vs Take single on UTF-8", func(t *testing.T) {
		createIntList := func(nums ...int64) core.List {
			exprs := make([]core.Expr, len(nums)+1)
			exprs[0] = core.NewSymbolAtom("List")
			for i, num := range nums {
				exprs[i+1] = core.NewIntAtom(int(num))
			}
			return core.List{Elements: exprs}
		}

		// StringTake gets 4th character (é)
		result1 := StringTakeAt(utf8Str, createIntList(4))
		if str, ok := core.ExtractString(result1); !ok || str != "é" {
			t.Errorf("StringTake(\"café\", [4]) expected \"é\", got %q", str)
		}

		// Take gets 4th byte (first byte of é which is 0xC3)
		result2 := StringTakeByteAt(utf8Str, createIntList(4))
		if str, ok := core.ExtractString(result2); !ok || str != "Ã" {
			t.Errorf("Take(\"café\", [4]) expected \"Ã\" (0xC3), got %q", str)
		}
	})

	t.Run("Test with complex Unicode", func(t *testing.T) {
		// String with various Unicode characters
		complexStr := "🙂café👍" // emoji (4 bytes each) + regular chars

		// StringTake should work on Unicode characters
		result := StringTake(complexStr, 3)
		if str, ok := core.ExtractString(result); !ok || str != "🙂ca" {
			t.Errorf("StringTake complex Unicode expected \"🙂ca\", got %q", str)
		}

		// Verify it's treating é as one character
		result2 := StringTake("café", -1)
		if str, ok := core.ExtractString(result2); !ok || str != "é" {
			t.Errorf("StringTake last character expected \"é\", got %q", str)
		}
	})
}
