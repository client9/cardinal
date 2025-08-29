package core

func NewBool(value bool) Symbol {
	if value {
		return symbolTrue
	} else {
		return symbolFalse
	}
}

// IsBool checks if an expression is a boolean value (True/False symbol)
func IsBool(expr Expr) bool {
	// Check new Symbol type first
	if s, ok := expr.(Symbol); ok {
		return s == symbolTrue || s == symbolFalse
	}
	return false
}

// ExtractBool safely extracts a boolean value from an Expr
// Note: NewBool returns symbols "True"/"False", so we check for those
func ExtractBool(expr Expr) (bool, bool) {
	// Check new Symbol type first
	if s, ok := expr.(Symbol); ok {
		if s == symbolTrue {
			return true, true
		}
		if s == symbolFalse {
			return false, true
		}
	}
	return false, false
}
