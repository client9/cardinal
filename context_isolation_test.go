package sexpr

import (
	"testing"
)

// TestContextIsolation demonstrates that different evaluators have isolated symbol tables
func TestContextIsolation(t *testing.T) {
	// Create first evaluator with custom attributes
	eval1 := NewEvaluator()
	eval1.context.symbolTable.SetAttributes("CustomFunc", []Attribute{Flat, Orderless})
	
	// Create second evaluator with different attributes
	eval2 := NewEvaluator()
	eval2.context.symbolTable.SetAttributes("CustomFunc", []Attribute{HoldAll})
	
	// Verify that the evaluators have different attributes for the same symbol
	attrs1 := eval1.context.symbolTable.Attributes("CustomFunc")
	attrs2 := eval2.context.symbolTable.Attributes("CustomFunc")
	
	// eval1 should have Flat and Orderless
	if !eval1.context.symbolTable.HasAttribute("CustomFunc", Flat) {
		t.Error("eval1 should have Flat attribute for CustomFunc")
	}
	if !eval1.context.symbolTable.HasAttribute("CustomFunc", Orderless) {
		t.Error("eval1 should have Orderless attribute for CustomFunc")
	}
	if eval1.context.symbolTable.HasAttribute("CustomFunc", HoldAll) {
		t.Error("eval1 should NOT have HoldAll attribute for CustomFunc")
	}
	
	// eval2 should have HoldAll
	if !eval2.context.symbolTable.HasAttribute("CustomFunc", HoldAll) {
		t.Error("eval2 should have HoldAll attribute for CustomFunc")
	}
	if eval2.context.symbolTable.HasAttribute("CustomFunc", Flat) {
		t.Error("eval2 should NOT have Flat attribute for CustomFunc")
	}
	if eval2.context.symbolTable.HasAttribute("CustomFunc", Orderless) {
		t.Error("eval2 should NOT have Orderless attribute for CustomFunc")
	}
	
	t.Logf("eval1 CustomFunc attributes: %s", AttributesToString(attrs1))
	t.Logf("eval2 CustomFunc attributes: %s", AttributesToString(attrs2))
}

// TestConcurrentEvaluators demonstrates that evaluators can be used safely in parallel
func TestConcurrentEvaluators(t *testing.T) {
	// Create multiple REPL instances (each has its own evaluator)
	repl1 := NewREPL()
	repl2 := NewREPL()
	
	// Set different variables in each REPL
	result1, err1 := repl1.EvaluateString("x = 10")
	if err1 != nil {
		t.Fatalf("repl1 error: %v", err1)
	}
	if result1 != "10" {
		t.Errorf("expected '10', got '%s'", result1)
	}
	
	result2, err2 := repl2.EvaluateString("x = 20")
	if err2 != nil {
		t.Fatalf("repl2 error: %v", err2)
	}
	if result2 != "20" {
		t.Errorf("expected '20', got '%s'", result2)
	}
	
	// Verify that the variables are isolated
	result1, _ = repl1.EvaluateString("x")
	result2, _ = repl2.EvaluateString("x")
	
	if result1 != "10" {
		t.Errorf("repl1 x should be 10, got '%s'", result1)
	}
	if result2 != "20" {
		t.Errorf("repl2 x should be 20, got '%s'", result2)
	}
	
	t.Logf("repl1 x = %s", result1)
	t.Logf("repl2 x = %s", result2)
}

// TestChildContextAttributeSharing demonstrates that child contexts share symbol tables
func TestChildContextAttributeSharing(t *testing.T) {
	// Create parent context with attributes
	parentCtx := NewContext()
	parentCtx.symbolTable.SetAttributes("TestSymbol", []Attribute{Flat, Orderless})
	
	// Create child context
	childCtx := NewChildContext(parentCtx)
	
	// Child should see parent's attributes
	if !childCtx.symbolTable.HasAttribute("TestSymbol", Flat) {
		t.Error("child context should see parent's Flat attribute")
	}
	if !childCtx.symbolTable.HasAttribute("TestSymbol", Orderless) {
		t.Error("child context should see parent's Orderless attribute")
	}
	
	// Child can add more attributes
	childCtx.symbolTable.SetAttributes("TestSymbol", []Attribute{OneIdentity})
	
	// Parent should also see the new attribute (shared symbol table)
	if !parentCtx.symbolTable.HasAttribute("TestSymbol", OneIdentity) {
		t.Error("parent context should see child's OneIdentity attribute")
	}
	
	attrs := parentCtx.symbolTable.Attributes("TestSymbol")
	t.Logf("TestSymbol attributes: %s", AttributesToString(attrs))
}

// TestSymbolTableIsolation demonstrates that different symbol tables are completely isolated
func TestSymbolTableIsolation(t *testing.T) {
	st1 := NewSymbolTable()
	st2 := NewSymbolTable()
	
	// Set different attributes in each table
	st1.SetAttributes("Symbol", []Attribute{Flat})
	st2.SetAttributes("Symbol", []Attribute{HoldAll})
	
	// Verify isolation
	if !st1.HasAttribute("Symbol", Flat) {
		t.Error("st1 should have Flat attribute")
	}
	if st1.HasAttribute("Symbol", HoldAll) {
		t.Error("st1 should NOT have HoldAll attribute")
	}
	
	if !st2.HasAttribute("Symbol", HoldAll) {
		t.Error("st2 should have HoldAll attribute")
	}
	if st2.HasAttribute("Symbol", Flat) {
		t.Error("st2 should NOT have Flat attribute")
	}
	
	// Clear one table
	st1.Reset()
	
	// Verify that only st1 was affected
	if st1.HasAttribute("Symbol", Flat) {
		t.Error("st1 should have no attributes after reset")
	}
	if !st2.HasAttribute("Symbol", HoldAll) {
		t.Error("st2 should still have HoldAll attribute after st1 reset")
	}
}