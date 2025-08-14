// Package core implements s-expression regular expression matching.
// This system provides regex-like pattern matching for s-expressions using
// either fast-path direct matching or Thompson NFA for complex patterns.
//
// The design avoids recursive backtracking to prevent exponential blowups,
// instead using Thompson NFA construction with state set simulation for
// linear time complexity.
//
// Features:
// - Complete Thompson NFA implementation for Or, Not, and complex quantifiers
// - Named capture propagation through complex NFA structures
// - Multi-strategy optimization (Direct vs NFA) based on pattern complexity
// - Custom predicates for semantic matching beyond structural patterns
package core

import (
	"fmt"
)

// Pattern represents a compiled s-expression pattern
type Pattern interface {
	String() string
}

// Basic pattern types
type LiteralPattern struct {
	Value Expr
}

func (p *LiteralPattern) String() string {
	return fmt.Sprintf("Literal(%s)", p.Value.String())
}

type HeadPattern struct {
	HeadType string
}

func (p *HeadPattern) String() string {
	return fmt.Sprintf("MatchHead(%s)", p.HeadType)
}

type AnyPattern struct{}

func (p *AnyPattern) String() string {
	return "MatchAny"
}

// Composite patterns
type SequencePattern struct {
	Patterns []Pattern
}

func (p *SequencePattern) String() string {
	return fmt.Sprintf("Sequence(%v)", p.Patterns)
}

type OrPattern struct {
	Alternatives []Pattern
}

func (p *OrPattern) String() string {
	return fmt.Sprintf("Or(%v)", p.Alternatives)
}

type NotPattern struct {
	Inner Pattern
}

func (p *NotPattern) String() string {
	return fmt.Sprintf("Not(%s)", p.Inner.String())
}

// Quantifier patterns
type ZeroOrMorePattern struct {
	Inner  Pattern
	Greedy bool
}

func (p *ZeroOrMorePattern) String() string {
	if p.Greedy {
		return fmt.Sprintf("ZeroOrMore(%s)", p.Inner.String())
	}
	return fmt.Sprintf("ZeroOrMoreLazy(%s)", p.Inner.String())
}

type OneOrMorePattern struct {
	Inner  Pattern
	Greedy bool
}

func (p *OneOrMorePattern) String() string {
	if p.Greedy {
		return fmt.Sprintf("OneOrMore(%s)", p.Inner.String())
	}
	return fmt.Sprintf("OneOrMoreLazy(%s)", p.Inner.String())
}

type ZeroOrOnePattern struct {
	Inner  Pattern
	Greedy bool
}

func (p *ZeroOrOnePattern) String() string {
	if p.Greedy {
		return fmt.Sprintf("Optional(%s)", p.Inner.String())
	}
	return fmt.Sprintf("OptionalLazy(%s)", p.Inner.String())
}

// Named capture pattern
type NamedPattern struct {
	Name  string
	Inner Pattern
}

func (p *NamedPattern) String() string {
	return fmt.Sprintf("Named(%s, %s)", p.Name, p.Inner.String())
}

// Custom predicate pattern
type PredicatePattern struct {
	Inner     Pattern
	Predicate func(Expr) bool
	Name      string // Optional name for debugging
}

func (p *PredicatePattern) String() string {
	if p.Name != "" {
		return fmt.Sprintf("Predicate(%s, %s)", p.Name, p.Inner.String())
	}
	return fmt.Sprintf("Predicate(%s)", p.Inner.String())
}

// MatchResult represents the result of a pattern match
type MatchResult struct {
	Matched  bool
	Bindings map[string]Expr // Named captures
	Consumed int             // Number of expressions consumed
}

// ExecutionStrategy determines how pattern matching is performed
type ExecutionStrategy int

const (
	StrategyDirect ExecutionStrategy = iota // Fast path for simple patterns
	StrategyNFA                             // Full NFA for complex patterns
)

func (s ExecutionStrategy) String() string {
	switch s {
	case StrategyDirect:
		return "Direct"
	case StrategyNFA:
		return "NFA"
	default:
		return "Unknown"
	}
}

// Matcher represents a compiled pattern ready for execution
type Matcher interface {
	Match(expr Expr) MatchResult
	Strategy() ExecutionStrategy
}

// DirectMatcher implements fast-path matching for simple patterns
type DirectMatcher struct {
	pattern Pattern
}

// Match performs direct pattern matching without state maintenance
func (dm *DirectMatcher) Match(expr Expr) MatchResult {
	bindings := make(map[string]Expr)
	matched := dm.matchDirectWithBindings(expr, dm.pattern, bindings)
	return MatchResult{
		Matched:  matched,
		Bindings: bindings,
		Consumed: 1,
	}
}

// Strategy returns the execution strategy for DirectMatcher
func (dm *DirectMatcher) Strategy() ExecutionStrategy {
	return StrategyDirect
}

func (dm *DirectMatcher) matchDirectWithBindings(expr Expr, pattern Pattern, bindings map[string]Expr) bool {
	switch p := pattern.(type) {
	case *LiteralPattern:
		return expr.Equal(p.Value)

	case *HeadPattern:
		return expr.Head() == p.HeadType

	case *AnyPattern:
		return true

	case *NamedPattern:
		// First check if inner pattern matches
		if dm.matchDirectWithBindings(expr, p.Inner, bindings) {
			// If it matches, record the binding
			bindings[p.Name] = expr
			return true
		}
		return false

	case *SequencePattern:
		return dm.matchSequenceWithBindings(expr, p.Patterns, bindings)

	case *PredicatePattern:
		// First check if inner pattern matches
		if dm.matchDirectWithBindings(expr, p.Inner, bindings) {
			// If inner matches, test the predicate
			return p.Predicate(expr)
		}
		return false

	case *OrPattern:
		// Try each alternative in order until one matches
		for _, alternative := range p.Alternatives {
			if dm.matchDirectWithBindings(expr, alternative, bindings) {
				return true
			}
		}
		return false

	default:
		panic(fmt.Sprintf("DirectMatcher called on non-direct pattern: %T", pattern))
	}
}

func (dm *DirectMatcher) matchSequenceWithBindings(expr Expr, patterns []Pattern, bindings map[string]Expr) bool {
	list, ok := expr.(List)
	if !ok {
		return false
	}

	elements := list.Tail()
	return dm.matchElementsWithBindings(elements, patterns, bindings)
}

func (dm *DirectMatcher) matchElementsWithBindings(elements []Expr, patterns []Pattern, bindings map[string]Expr) bool {
	// Handle trailing quantifiers specially (optimization for common case)
	if len(patterns) > 0 {
		lastPattern := patterns[len(patterns)-1]

		// Unwrap Named patterns to get the actual quantifier type
		actualLastPattern := lastPattern
		if named, ok := lastPattern.(*NamedPattern); ok {
			actualLastPattern = named.Inner
		}

		switch actualLastPattern.(type) {
		case *ZeroOrMorePattern:
			return dm.matchWithTrailingQuantifierBindings(elements, patterns, "ZeroOrMore", bindings)
		case *OneOrMorePattern:
			return dm.matchWithTrailingQuantifierBindings(elements, patterns, "OneOrMore", bindings)
		case *ZeroOrOnePattern:
			return dm.matchWithTrailingQuantifierBindings(elements, patterns, "ZeroOrOne", bindings)
		}
	}

	// Check if any patterns are quantifiers - if so, use general quantifier matching
	hasQuantifiers := false
	for _, p := range patterns {
		// Unwrap Named patterns to check actual pattern type
		actualPattern := p
		if named, ok := p.(*NamedPattern); ok {
			actualPattern = named.Inner
		}

		switch actualPattern.(type) {
		case *ZeroOrMorePattern, *OneOrMorePattern, *ZeroOrOnePattern:
			hasQuantifiers = true
			break
		}
	}

	if hasQuantifiers {
		return dm.matchWithQuantifiers(elements, patterns, bindings)
	}

	// Simple 1:1 matching for pure sequences
	if len(elements) != len(patterns) {
		return false
	}

	for i, element := range elements {
		if !dm.matchDirectWithBindings(element, patterns[i], bindings) {
			return false
		}
	}

	return true
}

// General quantifier matching using greedy semantics
// This handles quantifiers in any position, not just trailing
func (dm *DirectMatcher) matchWithQuantifiers(elements []Expr, patterns []Pattern, bindings map[string]Expr) bool {
	// Use greedy matching: try to match as much as possible for each quantifier
	return dm.matchQuantifiersGreedy(elements, 0, patterns, 0, bindings)
}

// Recursive greedy matching algorithm
func (dm *DirectMatcher) matchQuantifiersGreedy(elements []Expr, elemIdx int, patterns []Pattern, patIdx int, bindings map[string]Expr) bool {
	// Base case: processed all patterns
	if patIdx >= len(patterns) {
		return elemIdx >= len(elements) // Success if consumed all elements too
	}

	// Base case: no more elements but still have patterns to match
	if elemIdx >= len(elements) {
		// Check if remaining patterns are all optional (ZeroOrMore, ZeroOrOne)
		for i := patIdx; i < len(patterns); i++ {
			switch patterns[i].(type) {
			case *ZeroOrMorePattern, *ZeroOrOnePattern:
				continue // These can match zero elements
			default:
				return false // Required pattern with no elements left
			}
		}
		return true
	}

	pattern := patterns[patIdx]

	// Unwrap Named patterns to get to the actual pattern type
	actualPattern := pattern
	var namedBinding string
	if named, ok := pattern.(*NamedPattern); ok {
		actualPattern = named.Inner
		namedBinding = named.Name
	}

	switch q := actualPattern.(type) {
	case *ZeroOrMorePattern:
		// Try matching 0, 1, 2, ... elements greedily
		// Start with maximum possible matches and work backwards
		maxMatches := len(elements) - elemIdx

		// Calculate minimum elements needed for remaining patterns
		minNeeded := 0
		for i := patIdx + 1; i < len(patterns); i++ {
			switch patterns[i].(type) {
			case *OneOrMorePattern:
				minNeeded++ // At least 1 element needed
			case *ZeroOrMorePattern, *ZeroOrOnePattern:
				// Optional, don't count
			default:
				minNeeded++ // Regular pattern needs 1 element
			}
		}

		// Try from maximum down to 0, ensuring we leave enough for remaining patterns
		for matchCount := maxMatches - minNeeded; matchCount >= 0; matchCount-- {
			// Check if all elements in this range match the quantifier pattern
			allMatch := true
			for i := 0; i < matchCount; i++ {
				if !dm.matchDirectWithBindings(elements[elemIdx+i], q.Inner, bindings) {
					allMatch = false
					break
				}
			}

			if allMatch {
				// Try continuing with remaining patterns
				if dm.matchQuantifiersGreedy(elements, elemIdx+matchCount, patterns, patIdx+1, bindings) {
					// Success! Record the Named binding if needed
					if namedBinding != "" {
						if matchCount == 0 {
							bindings[namedBinding] = nil // ZeroOrMore can match nothing
						} else if matchCount == 1 {
							bindings[namedBinding] = elements[elemIdx]
						} else {
							// Create a list of matched elements
							matchedElements := make([]Expr, matchCount)
							for i := 0; i < matchCount; i++ {
								matchedElements[i] = elements[elemIdx+i]
							}
							bindings[namedBinding] = NewList("List", matchedElements...)
						}
					}
					return true
				}
			}
		}
		return false

	case *OneOrMorePattern:
		// Similar to ZeroOrMore but requires at least 1 match
		maxMatches := len(elements) - elemIdx

		// Calculate minimum elements needed for remaining patterns
		minNeeded := 0
		for i := patIdx + 1; i < len(patterns); i++ {
			switch patterns[i].(type) {
			case *OneOrMorePattern:
				minNeeded++ // At least 1 element needed
			case *ZeroOrMorePattern, *ZeroOrOnePattern:
				// Optional, don't count
			default:
				minNeeded++ // Regular pattern needs 1 element
			}
		}

		// Try from maximum down to 1, ensuring we leave enough for remaining patterns
		for matchCount := maxMatches - minNeeded; matchCount >= 1; matchCount-- {
			// Check if all elements in this range match the quantifier pattern
			allMatch := true
			for i := 0; i < matchCount; i++ {
				if !dm.matchDirectWithBindings(elements[elemIdx+i], q.Inner, bindings) {
					allMatch = false
					break
				}
			}

			if allMatch {
				// Try continuing with remaining patterns
				if dm.matchQuantifiersGreedy(elements, elemIdx+matchCount, patterns, patIdx+1, bindings) {
					// Success! Record the Named binding if needed
					if namedBinding != "" {
						if matchCount == 1 {
							bindings[namedBinding] = elements[elemIdx]
						} else {
							// Create a list of matched elements
							matchedElements := make([]Expr, matchCount)
							for i := 0; i < matchCount; i++ {
								matchedElements[i] = elements[elemIdx+i]
							}
							bindings[namedBinding] = NewList("List", matchedElements...)
						}
					}
					return true
				}
			}
		}
		return false

	case *ZeroOrOnePattern:
		// Try matching 1 element first (greedy), then 0
		// Try 1 element
		if dm.matchDirectWithBindings(elements[elemIdx], q.Inner, bindings) {
			if dm.matchQuantifiersGreedy(elements, elemIdx+1, patterns, patIdx+1, bindings) {
				// Success! Record the Named binding if needed
				if namedBinding != "" {
					bindings[namedBinding] = elements[elemIdx]
				}
				return true
			}
		}
		// Try 0 elements
		if dm.matchQuantifiersGreedy(elements, elemIdx, patterns, patIdx+1, bindings) {
			// Success! Record the Named binding if needed (matched nothing)
			if namedBinding != "" {
				bindings[namedBinding] = nil // ZeroOrOne can match nothing
			}
			return true
		}
		return false

	default:
		// Regular pattern - must match exactly one element
		if dm.matchDirectWithBindings(elements[elemIdx], pattern, bindings) {
			return dm.matchQuantifiersGreedy(elements, elemIdx+1, patterns, patIdx+1, bindings)
		}
		return false
	}
}

// Unified trailing quantifier matching that handles Named patterns correctly
func (dm *DirectMatcher) matchWithTrailingQuantifierBindings(elements []Expr, patterns []Pattern, quantType string, bindings map[string]Expr) bool {
	if len(patterns) == 0 {
		return false
	}

	prefixPatterns := patterns[:len(patterns)-1]
	lastPattern := patterns[len(patterns)-1]

	// Extract the inner quantifier pattern
	var innerPattern Pattern
	var namedBinding string

	if named, ok := lastPattern.(*NamedPattern); ok {
		namedBinding = named.Name
		// Get the inner quantifier and then its inner pattern
		switch q := named.Inner.(type) {
		case *ZeroOrMorePattern:
			innerPattern = q.Inner
		case *OneOrMorePattern:
			innerPattern = q.Inner
		case *ZeroOrOnePattern:
			innerPattern = q.Inner
		default:
			return false // Shouldn't happen
		}
	} else {
		// Direct quantifier
		switch q := lastPattern.(type) {
		case *ZeroOrMorePattern:
			innerPattern = q.Inner
		case *OneOrMorePattern:
			innerPattern = q.Inner
		case *ZeroOrOnePattern:
			innerPattern = q.Inner
		default:
			return false // Shouldn't happen
		}
	}

	// Use the appropriate matching strategy
	var matched bool
	var matchedElements []Expr

	switch quantType {
	case "ZeroOrMore":
		matched = dm.matchTrailingZeroOrMore(elements, prefixPatterns, innerPattern, &matchedElements, bindings)
	case "OneOrMore":
		matched = dm.matchTrailingOneOrMore(elements, prefixPatterns, innerPattern, &matchedElements, bindings)
	case "ZeroOrOne":
		matched = dm.matchTrailingZeroOrOne(elements, prefixPatterns, innerPattern, &matchedElements, bindings)
	default:
		return false
	}

	// If successful and we have a Named binding, record it
	if matched && namedBinding != "" {
		if len(matchedElements) == 0 {
			if quantType == "ZeroOrMore" || quantType == "ZeroOrOne" {
				bindings[namedBinding] = nil // Can match nothing
			}
		} else if len(matchedElements) == 1 {
			bindings[namedBinding] = matchedElements[0]
		} else {
			bindings[namedBinding] = NewList("List", matchedElements...)
		}
	}

	return matched
}

// Helper methods that return the matched elements for Named binding
func (dm *DirectMatcher) matchTrailingZeroOrMore(elements []Expr, prefixPatterns []Pattern, quantPattern Pattern, matchedElements *[]Expr, bindings map[string]Expr) bool {
	// Must match prefix first
	if len(elements) < len(prefixPatterns) {
		return false
	}

	// Match prefix patterns
	for i, pattern := range prefixPatterns {
		if !dm.matchDirectWithBindings(elements[i], pattern, bindings) {
			return false
		}
	}

	// Remaining elements must all match the quantified pattern
	remaining := elements[len(prefixPatterns):]
	*matchedElements = make([]Expr, len(remaining))
	for i, element := range remaining {
		if !dm.matchDirectWithBindings(element, quantPattern, bindings) {
			return false
		}
		(*matchedElements)[i] = element
	}

	return true
}

func (dm *DirectMatcher) matchTrailingOneOrMore(elements []Expr, prefixPatterns []Pattern, quantPattern Pattern, matchedElements *[]Expr, bindings map[string]Expr) bool {
	// Need at least one element for the quantifier
	if len(elements) <= len(prefixPatterns) {
		return false
	}

	return dm.matchTrailingZeroOrMore(elements, prefixPatterns, quantPattern, matchedElements, bindings)
}

func (dm *DirectMatcher) matchTrailingZeroOrOne(elements []Expr, prefixPatterns []Pattern, quantPattern Pattern, matchedElements *[]Expr, bindings map[string]Expr) bool {
	// Can have 0 or 1 additional element
	if len(elements) < len(prefixPatterns) || len(elements) > len(prefixPatterns)+1 {
		return false
	}

	// Match prefix
	for i, pattern := range prefixPatterns {
		if !dm.matchDirectWithBindings(elements[i], pattern, bindings) {
			return false
		}
	}

	// If there's one more element, it must match the optional pattern
	if len(elements) == len(prefixPatterns)+1 {
		lastElement := elements[len(elements)-1]
		if !dm.matchDirectWithBindings(lastElement, quantPattern, bindings) {
			return false
		}
		*matchedElements = []Expr{lastElement}
	} else {
		*matchedElements = []Expr{} // Matched nothing
	}

	return true
}

// Optimized matching with bindings for patterns ending with quantifiers
func (dm *DirectMatcher) matchWithTrailingZeroOrMoreBindings(elements []Expr, prefixPatterns []Pattern, quantPattern Pattern, bindings map[string]Expr) bool {
	// Must match prefix first
	if len(elements) < len(prefixPatterns) {
		return false
	}

	// Match prefix patterns
	for i, pattern := range prefixPatterns {
		if !dm.matchDirectWithBindings(elements[i], pattern, bindings) {
			return false
		}
	}

	// Remaining elements must all match the quantified pattern
	remaining := elements[len(prefixPatterns):]
	for _, element := range remaining {
		if !dm.matchDirectWithBindings(element, quantPattern, bindings) {
			return false
		}
	}

	return true
}

func (dm *DirectMatcher) matchWithTrailingOneOrMoreBindings(elements []Expr, prefixPatterns []Pattern, quantPattern Pattern, bindings map[string]Expr) bool {
	// Need at least one element for the quantifier
	if len(elements) <= len(prefixPatterns) {
		return false
	}

	return dm.matchWithTrailingZeroOrMoreBindings(elements, prefixPatterns, quantPattern, bindings)
}

func (dm *DirectMatcher) matchWithTrailingZeroOrOneBindings(elements []Expr, prefixPatterns []Pattern, quantPattern Pattern, bindings map[string]Expr) bool {
	// Can have 0 or 1 additional element
	if len(elements) < len(prefixPatterns) || len(elements) > len(prefixPatterns)+1 {
		return false
	}

	// Match prefix
	for i, pattern := range prefixPatterns {
		if !dm.matchDirectWithBindings(elements[i], pattern, bindings) {
			return false
		}
	}

	// If there's one more element, it must match the optional pattern
	if len(elements) == len(prefixPatterns)+1 {
		lastElement := elements[len(elements)-1]
		if !dm.matchDirectWithBindings(lastElement, quantPattern, bindings) {
			return false
		}
	}

	return true
}

// Thompson NFA implementation for complex patterns

// NFAState represents a state in the Thompson NFA
type NFAState struct {
	ID          int
	Transitions []NFATransition
	IsAccept    bool
	// Group boundary tags for Named patterns
	GroupStart string // If non-empty, entering this state starts capturing for this group name
	GroupEnd   string // If non-empty, exiting this state ends capturing for this group name
}

// NFATransition represents a transition between states
type NFATransition struct {
	Type      NFATransitionType
	Target    int
	Condition NFACondition
}

// NFATransitionType defines the type of NFA transition
type NFATransitionType int

const (
	EpsilonTransition NFATransitionType = iota // ε-transition (no input consumed)
	MatchTransition                            // Match against input expression
	SplitTransition                            // Split to multiple states (for quantifiers)
)

// NFACondition defines what an NFA transition matches
type NFACondition interface {
	Matches(expr Expr, bindings map[string]Expr) bool
	String() string
}

// Specific condition types
type LiteralCondition struct {
	Value Expr
}

func (c *LiteralCondition) Matches(expr Expr, bindings map[string]Expr) bool {
	return expr.Equal(c.Value)
}

func (c *LiteralCondition) String() string {
	return fmt.Sprintf("Literal(%s)", c.Value.String())
}

type HeadCondition struct {
	HeadType string
}

func (c *HeadCondition) Matches(expr Expr, bindings map[string]Expr) bool {
	return expr.Head() == c.HeadType
}

func (c *HeadCondition) String() string {
	return fmt.Sprintf("Head(%s)", c.HeadType)
}

type AnyCondition struct{}

func (c *AnyCondition) Matches(expr Expr, bindings map[string]Expr) bool {
	return true
}

func (c *AnyCondition) String() string {
	return "Any"
}

type NotCondition struct {
	Inner NFACondition
}

func (c *NotCondition) Matches(expr Expr, bindings map[string]Expr) bool {
	// Create temporary bindings to avoid polluting the real ones
	tempBindings := make(map[string]Expr)
	return !c.Inner.Matches(expr, tempBindings)
}

func (c *NotCondition) String() string {
	return fmt.Sprintf("Not(%s)", c.Inner.String())
}

type PredicateCondition struct {
	Inner     NFACondition
	Predicate func(Expr) bool
	Name      string
}

func (c *PredicateCondition) Matches(expr Expr, bindings map[string]Expr) bool {
	// First check if inner condition matches
	if c.Inner.Matches(expr, bindings) {
		// If inner matches, test the predicate
		return c.Predicate(expr)
	}
	return false
}

func (c *PredicateCondition) String() string {
	if c.Name != "" {
		return fmt.Sprintf("PredicateCondition(%s, %s)", c.Name, c.Inner.String())
	}
	return fmt.Sprintf("PredicateCondition(%s)", c.Inner.String())
}

type ComplexNotCondition struct {
	InnerFragment NFAFragment
	Builder       *NFABuilder
}

func (c *ComplexNotCondition) Matches(expr Expr, bindings map[string]Expr) bool {
	// Create a copy of the builder's states and mark the accept state
	states := make([]NFAState, len(c.Builder.states))
	copy(states, c.Builder.states)

	// Ensure the accept state is marked
	states[c.InnerFragment.Accept].IsAccept = true

	// Create a temporary NFA executor for the inner pattern
	nfa := &NFA{
		States:      states,
		StartState:  c.InnerFragment.Start,
		AcceptState: c.InnerFragment.Accept,
	}
	executor := NewNFAExecutor(nfa)

	// Test if the inner pattern matches - if it does, we return false (NOT)
	result := executor.Match(expr)
	return !result.Matched
}

func (c *ComplexNotCondition) String() string {
	return fmt.Sprintf("ComplexNot([fragment])")
}

// NFA represents a Thompson NFA
type NFA struct {
	States      []NFAState
	StartState  int
	AcceptState int
}

// NFAFragment represents a partial NFA during construction
type NFAFragment struct {
	Start  int
	Accept int
}

// NFABuilder constructs Thompson NFAs from patterns
type NFABuilder struct {
	nextStateID int
	states      []NFAState
}

func NewNFABuilder() *NFABuilder {
	return &NFABuilder{
		nextStateID: 0,
		states:      make([]NFAState, 0),
	}
}

func (b *NFABuilder) newState() int {
	id := b.nextStateID
	b.nextStateID++
	b.states = append(b.states, NFAState{
		ID:          id,
		Transitions: make([]NFATransition, 0),
		IsAccept:    false,
	})
	return id
}

func (b *NFABuilder) addTransition(from int, transType NFATransitionType, to int, condition NFACondition) {
	b.states[from].Transitions = append(b.states[from].Transitions, NFATransition{
		Type:      transType,
		Target:    to,
		Condition: condition,
	})
}

func (b *NFABuilder) addEpsilonTransition(from, to int) {
	b.addTransition(from, EpsilonTransition, to, nil)
}

// propagatePredicateCondition wraps all match transitions in an NFA fragment with PredicateCondition
func (b *NFABuilder) propagatePredicateCondition(fragment NFAFragment, predicate func(Expr) bool, name string) {
	// Visit all states reachable from the fragment start
	visited := make(map[int]bool)
	b.visitStatesAndWrapPredicateConditions(fragment.Start, predicate, name, visited)
}

func (b *NFABuilder) visitStatesAndWrapPredicateConditions(stateID int, predicate func(Expr) bool, name string, visited map[int]bool) {
	if visited[stateID] {
		return
	}
	visited[stateID] = true

	state := &b.states[stateID]

	// Wrap all match transitions with PredicateCondition
	for i := range state.Transitions {
		transition := &state.Transitions[i]
		if transition.Type == MatchTransition && transition.Condition != nil {
			// Only wrap if not already a PredicateCondition to avoid double-wrapping
			if _, isPredicate := transition.Condition.(*PredicateCondition); !isPredicate {
				transition.Condition = &PredicateCondition{Inner: transition.Condition, Predicate: predicate, Name: name}
			}
		}

		// Recursively visit target states
		b.visitStatesAndWrapPredicateConditions(transition.Target, predicate, name, visited)
	}
}

// Thompson construction methods
func (b *NFABuilder) BuildLiteral(value Expr) NFAFragment {
	start := b.newState()
	accept := b.newState()
	b.addTransition(start, MatchTransition, accept, &LiteralCondition{Value: value})
	return NFAFragment{Start: start, Accept: accept}
}

func (b *NFABuilder) BuildHead(headType string) NFAFragment {
	start := b.newState()
	accept := b.newState()
	b.addTransition(start, MatchTransition, accept, &HeadCondition{HeadType: headType})
	return NFAFragment{Start: start, Accept: accept}
}

func (b *NFABuilder) BuildAny() NFAFragment {
	start := b.newState()
	accept := b.newState()
	b.addTransition(start, MatchTransition, accept, &AnyCondition{})
	return NFAFragment{Start: start, Accept: accept}
}

func (b *NFABuilder) BuildPredicate(inner NFAFragment, predicate func(Expr) bool, name string) NFAFragment {
	start := b.newState()
	accept := b.newState()

	// If the inner fragment is a simple match, wrap its condition with PredicateCondition
	if len(b.states[inner.Start].Transitions) == 1 &&
		b.states[inner.Start].Transitions[0].Type == MatchTransition {
		innerCondition := b.states[inner.Start].Transitions[0].Condition
		predicateCondition := &PredicateCondition{Inner: innerCondition, Predicate: predicate, Name: name}
		b.addTransition(start, MatchTransition, accept, predicateCondition)
		return NFAFragment{Start: start, Accept: accept}
	}

	// For complex inner patterns, we need to propagate the predicate through all match transitions
	b.propagatePredicateCondition(inner, predicate, name)

	// Connect with epsilon transitions
	b.addEpsilonTransition(start, inner.Start)
	b.addEpsilonTransition(inner.Accept, accept)
	return NFAFragment{Start: start, Accept: accept}
}

func (b *NFABuilder) BuildNamed(name string, inner NFAFragment) NFAFragment {
	// For named patterns, we tag the boundary states instead of wrapping transitions
	start := b.newState()
	accept := b.newState()

	// Tag the start state to begin group capture
	b.states[start].GroupStart = name

	// Tag the accept state to end group capture
	b.states[accept].GroupEnd = name

	// Connect with epsilon transitions - the inner NFA is preserved intact
	b.addEpsilonTransition(start, inner.Start)
	b.addEpsilonTransition(inner.Accept, accept)

	return NFAFragment{Start: start, Accept: accept}
}

func (b *NFABuilder) BuildNot(inner NFAFragment) NFAFragment {
	start := b.newState()
	accept := b.newState()

	// For NOT patterns, we need to invert the condition
	if len(b.states[inner.Start].Transitions) == 1 &&
		b.states[inner.Start].Transitions[0].Type == MatchTransition {
		innerCondition := b.states[inner.Start].Transitions[0].Condition
		notCondition := &NotCondition{Inner: innerCondition}
		b.addTransition(start, MatchTransition, accept, notCondition)
		return NFAFragment{Start: start, Accept: accept}
	}

	// For complex patterns, create a ComplexNotCondition that can evaluate the inner NFA
	complexNotCondition := &ComplexNotCondition{InnerFragment: inner, Builder: b}
	b.addTransition(start, MatchTransition, accept, complexNotCondition)
	return NFAFragment{Start: start, Accept: accept}
}

// Concatenation: f1 followed by f2
func (b *NFABuilder) BuildConcat(f1, f2 NFAFragment) NFAFragment {
	b.addEpsilonTransition(f1.Accept, f2.Start)
	return NFAFragment{Start: f1.Start, Accept: f2.Accept}
}

// Union: f1 OR f2 (Thompson's union construction)
func (b *NFABuilder) BuildUnion(f1, f2 NFAFragment) NFAFragment {
	start := b.newState()
	accept := b.newState()

	// Epsilon transitions from start to both alternatives
	b.addEpsilonTransition(start, f1.Start)
	b.addEpsilonTransition(start, f2.Start)

	// Epsilon transitions from both accepts to final accept
	b.addEpsilonTransition(f1.Accept, accept)
	b.addEpsilonTransition(f2.Accept, accept)

	return NFAFragment{Start: start, Accept: accept}
}

// Star: f* (zero or more)
func (b *NFABuilder) BuildStar(f NFAFragment) NFAFragment {
	start := b.newState()
	accept := b.newState()

	// Epsilon from start to accept (zero matches)
	b.addEpsilonTransition(start, accept)

	// Epsilon from start to f.start (enter loop)
	b.addEpsilonTransition(start, f.Start)

	// Epsilon from f.accept back to f.start (loop)
	b.addEpsilonTransition(f.Accept, f.Start)

	// Epsilon from f.accept to final accept (exit loop)
	b.addEpsilonTransition(f.Accept, accept)

	return NFAFragment{Start: start, Accept: accept}
}

// Plus: f+ (one or more)
func (b *NFABuilder) BuildPlus(f NFAFragment) NFAFragment {
	accept := b.newState()

	// Epsilon from f.accept back to f.start (loop)
	b.addEpsilonTransition(f.Accept, f.Start)

	// Epsilon from f.accept to final accept (exit)
	b.addEpsilonTransition(f.Accept, accept)

	return NFAFragment{Start: f.Start, Accept: accept}
}

// Question: f? (zero or one)
func (b *NFABuilder) BuildQuestion(f NFAFragment) NFAFragment {
	start := b.newState()
	accept := b.newState()

	// Epsilon from start to accept (zero matches)
	b.addEpsilonTransition(start, accept)

	// Epsilon from start to f.start (one match)
	b.addEpsilonTransition(start, f.Start)

	// Epsilon from f.accept to final accept
	b.addEpsilonTransition(f.Accept, accept)

	return NFAFragment{Start: start, Accept: accept}
}

func (b *NFABuilder) BuildPattern(pattern Pattern) (NFAFragment, error) {
	switch p := pattern.(type) {
	case *LiteralPattern:
		return b.BuildLiteral(p.Value), nil

	case *HeadPattern:
		return b.BuildHead(p.HeadType), nil

	case *AnyPattern:
		return b.BuildAny(), nil

	case *NamedPattern:
		inner, err := b.BuildPattern(p.Inner)
		if err != nil {
			return NFAFragment{}, err
		}
		return b.BuildNamed(p.Name, inner), nil

	case *NotPattern:
		inner, err := b.BuildPattern(p.Inner)
		if err != nil {
			return NFAFragment{}, err
		}
		return b.BuildNot(inner), nil

	case *PredicatePattern:
		inner, err := b.BuildPattern(p.Inner)
		if err != nil {
			return NFAFragment{}, err
		}
		return b.BuildPredicate(inner, p.Predicate, p.Name), nil

	case *SequencePattern:
		if len(p.Patterns) == 0 {
			return NFAFragment{}, fmt.Errorf("empty sequence pattern")
		}

		result, err := b.BuildPattern(p.Patterns[0])
		if err != nil {
			return NFAFragment{}, err
		}

		for i := 1; i < len(p.Patterns); i++ {
			next, err := b.BuildPattern(p.Patterns[i])
			if err != nil {
				return NFAFragment{}, err
			}
			result = b.BuildConcat(result, next)
		}
		return result, nil

	case *OrPattern:
		if len(p.Alternatives) == 0 {
			return NFAFragment{}, fmt.Errorf("empty or pattern")
		}

		result, err := b.BuildPattern(p.Alternatives[0])
		if err != nil {
			return NFAFragment{}, err
		}

		for i := 1; i < len(p.Alternatives); i++ {
			next, err := b.BuildPattern(p.Alternatives[i])
			if err != nil {
				return NFAFragment{}, err
			}
			result = b.BuildUnion(result, next)
		}
		return result, nil

	case *ZeroOrMorePattern:
		inner, err := b.BuildPattern(p.Inner)
		if err != nil {
			return NFAFragment{}, err
		}
		return b.BuildStar(inner), nil

	case *OneOrMorePattern:
		inner, err := b.BuildPattern(p.Inner)
		if err != nil {
			return NFAFragment{}, err
		}
		return b.BuildPlus(inner), nil

	case *ZeroOrOnePattern:
		inner, err := b.BuildPattern(p.Inner)
		if err != nil {
			return NFAFragment{}, err
		}
		return b.BuildQuestion(inner), nil

	default:
		return NFAFragment{}, fmt.Errorf("unsupported pattern type: %T", pattern)
	}
}

func (b *NFABuilder) Compile(pattern Pattern) (*NFA, error) {
	fragment, err := b.BuildPattern(pattern)
	if err != nil {
		return nil, err
	}

	// Mark the accept state
	b.states[fragment.Accept].IsAccept = true

	return &NFA{
		States:      b.states,
		StartState:  fragment.Start,
		AcceptState: fragment.Accept,
	}, nil
}

// NFAExecutor runs Thompson NFA simulation
type NFAExecutor struct {
	nfa *NFA
}

func NewNFAExecutor(nfa *NFA) *NFAExecutor {
	return &NFAExecutor{nfa: nfa}
}

func (ne *NFAExecutor) Match(expr Expr) MatchResult {
	// For single expression matching, treat as sequence of one
	var exprs []Expr
	if list, ok := expr.(List); ok {
		exprs = list.Tail() // Get list elements
	} else {
		exprs = []Expr{expr} // Single expression
	}

	return ne.matchSequence(exprs)
}

// Strategy returns the execution strategy for NFAExecutor
func (ne *NFAExecutor) Strategy() ExecutionStrategy {
	return StrategyNFA
}

// CaptureContext represents capture state for a specific execution path
type CaptureContext struct {
	GroupCaptures map[string][]Expr // Elements captured within each group
	ActiveGroups  []string          // Currently active groups
}

func (cc *CaptureContext) Clone() *CaptureContext {
	newContext := &CaptureContext{
		GroupCaptures: make(map[string][]Expr),
		ActiveGroups:  make([]string, len(cc.ActiveGroups)),
	}
	copy(newContext.ActiveGroups, cc.ActiveGroups)

	// Deep copy group captures
	for group, captures := range cc.GroupCaptures {
		newCaptures := make([]Expr, len(captures))
		copy(newCaptures, captures)
		newContext.GroupCaptures[group] = newCaptures
	}

	return newContext
}

// ExecutionPath represents a single execution path through the NFA
type ExecutionPath struct {
	StateID int
	Context *CaptureContext
}

func (ne *NFAExecutor) matchSequence(exprs []Expr) MatchResult {
	// Initialize with single path at start state
	initialContext := &CaptureContext{
		GroupCaptures: make(map[string][]Expr),
		ActiveGroups:  make([]string, 0),
	}

	currentPaths := []*ExecutionPath{{
		StateID: ne.nfa.StartState,
		Context: initialContext,
	}}

	// Apply epsilon closure to initial paths
	currentPaths = ne.expandPathsWithEpsilon(currentPaths)

	// Process each expression
	for _, expr := range exprs {
		currentPaths = ne.stepPaths(currentPaths, expr)
		if len(currentPaths) == 0 {
			// No valid paths remaining
			return MatchResult{Matched: false, Bindings: make(map[string]Expr)}
		}
	}

	// Collect all successful paths and select the best one (leftmost greedy)
	var successfulPaths []*ExecutionPath
	for _, path := range currentPaths {
		if ne.nfa.States[path.StateID].IsAccept {
			successfulPaths = append(successfulPaths, path)
		}
	}

	if len(successfulPaths) > 0 {
		// For greedy matching, select the path that captures elements in leftmost groups first
		bestPath := ne.selectGreedyPath(successfulPaths)
		return MatchResult{
			Matched:  true,
			Bindings: ne.buildBindingsFromContext(bestPath.Context),
			Consumed: len(exprs),
		}
	}

	return MatchResult{Matched: false, Bindings: make(map[string]Expr)}
}

func (ne *NFAExecutor) stepPaths(currentPaths []*ExecutionPath, expr Expr) []*ExecutionPath {
	var nextPaths []*ExecutionPath

	// For each active path, follow matching transitions
	for _, path := range currentPaths {
		state := ne.nfa.States[path.StateID]

		for _, transition := range state.Transitions {
			if transition.Type == MatchTransition {
				if transition.Condition.Matches(expr, make(map[string]Expr)) {
					// Create new path context for this transition
					newContext := path.Context.Clone()

					// Capture this expression in all active groups
					for _, groupName := range newContext.ActiveGroups {
						newContext.GroupCaptures[groupName] = append(newContext.GroupCaptures[groupName], expr)
					}

					newPath := &ExecutionPath{
						StateID: transition.Target,
						Context: newContext,
					}

					nextPaths = append(nextPaths, newPath)
				}
			}
		}
	}

	// Expand all paths with epsilon closure
	return ne.expandPathsWithEpsilon(nextPaths)
}

func (ne *NFAExecutor) expandPathsWithEpsilon(paths []*ExecutionPath) []*ExecutionPath {
	result := make([]*ExecutionPath, 0)
	visited := make(map[string]bool) // Use string key: "stateID:contextHash" to avoid duplicate paths

	for _, path := range paths {
		ne.expandSinglePathWithEpsilon(path, &result, visited)
	}

	return result
}

func (ne *NFAExecutor) expandSinglePathWithEpsilon(path *ExecutionPath, result *[]*ExecutionPath, visited map[string]bool) {
	// Create a unique key for this path to avoid infinite loops
	key := fmt.Sprintf("%d:%p", path.StateID, path.Context)
	if visited[key] {
		return
	}
	visited[key] = true

	state := ne.nfa.States[path.StateID]
	context := path.Context

	// Handle group start/end for this state
	if state.GroupStart != "" {
		// Start a new group - clone context to avoid affecting other paths
		context = context.Clone()
		context.ActiveGroups = append(context.ActiveGroups, state.GroupStart)
		if _, exists := context.GroupCaptures[state.GroupStart]; !exists {
			context.GroupCaptures[state.GroupStart] = make([]Expr, 0)
		}
	}

	if state.GroupEnd != "" {
		// End a group - clone context to avoid affecting other paths
		context = context.Clone()
		for i, group := range context.ActiveGroups {
			if group == state.GroupEnd {
				context.ActiveGroups = append(context.ActiveGroups[:i], context.ActiveGroups[i+1:]...)
				break
			}
		}
	}

	// Update the path with potentially modified context
	updatedPath := &ExecutionPath{
		StateID: path.StateID,
		Context: context,
	}

	// Add this path to results
	*result = append(*result, updatedPath)

	// Follow epsilon transitions
	for _, transition := range state.Transitions {
		if transition.Type == EpsilonTransition {
			newPath := &ExecutionPath{
				StateID: transition.Target,
				Context: context,
			}
			ne.expandSinglePathWithEpsilon(newPath, result, visited)
		}
	}
}

func (ne *NFAExecutor) buildBindingsFromContext(context *CaptureContext) map[string]Expr {
	bindings := make(map[string]Expr)

	for groupName, captures := range context.GroupCaptures {
		if len(captures) == 0 {
			continue // No captures for this group
		} else if len(captures) == 1 {
			// Single element - return as-is
			bindings[groupName] = captures[0]
		} else {
			// Multiple elements - create List
			bindings[groupName] = NewList("List", captures...)
		}
	}

	return bindings
}

// selectGreedyPath implements leftmost greedy matching semantics
// For patterns like ZeroOrMore(x), ZeroOrMore(y), prefer paths where leftmost groups capture more
func (ne *NFAExecutor) selectGreedyPath(paths []*ExecutionPath) *ExecutionPath {
	if len(paths) == 1 {
		return paths[0]
	}

	bestPath := paths[0]
	for _, candidate := range paths[1:] {
		if ne.isMoreGreedy(candidate, bestPath) {
			bestPath = candidate
		}
	}
	return bestPath
}

// isMoreGreedy determines if candidate path is more greedy than current best
// Implements leftmost-first greedy semantics: earlier groups should capture more elements
func (ne *NFAExecutor) isMoreGreedy(candidate, current *ExecutionPath) bool {
	// Get all group names in alphabetical order for consistent comparison
	// In practice, for patterns like Named("x", ...), Named("y", ...), this gives x before y
	allGroups := make(map[string]bool)
	for group := range candidate.Context.GroupCaptures {
		allGroups[group] = true
	}
	for group := range current.Context.GroupCaptures {
		allGroups[group] = true
	}

	var sortedGroups []string
	for group := range allGroups {
		sortedGroups = append(sortedGroups, group)
	}

	// Sort to ensure consistent comparison (x comes before y)
	for i := 0; i < len(sortedGroups)-1; i++ {
		for j := i + 1; j < len(sortedGroups); j++ {
			if sortedGroups[i] > sortedGroups[j] {
				sortedGroups[i], sortedGroups[j] = sortedGroups[j], sortedGroups[i]
			}
		}
	}

	// Compare group by group - leftmost groups should be more greedy
	for _, group := range sortedGroups {
		candidateLen := len(candidate.Context.GroupCaptures[group])
		currentLen := len(current.Context.GroupCaptures[group])

		if candidateLen != currentLen {
			// First differing group - prefer the path where this group captured more
			return candidateLen > currentLen
		}
	}

	// All groups captured the same amount - no preference
	return false
}

// Pattern builder functions
func MatchLiteral(expr Expr) Pattern {
	return &LiteralPattern{Value: expr}
}

func MatchHead(headType string) Pattern {
	return &HeadPattern{HeadType: headType}
}

func MatchAny() Pattern {
	return &AnyPattern{}
}

func MatchOr(patterns ...Pattern) Pattern {
	return &OrPattern{Alternatives: patterns}
}

func MatchNot(pattern Pattern) Pattern {
	return &NotPattern{Inner: pattern}
}

func MatchSequence(patterns ...Pattern) Pattern {
	return &SequencePattern{Patterns: patterns}
}

func ZeroOrMore(pattern Pattern) Pattern {
	return &ZeroOrMorePattern{Inner: pattern, Greedy: true}
}

func ZeroOrMoreLazy(pattern Pattern) Pattern {
	return &ZeroOrMorePattern{Inner: pattern, Greedy: false}
}

func OneOrMore(pattern Pattern) Pattern {
	return &OneOrMorePattern{Inner: pattern, Greedy: true}
}

func OneOrMoreLazy(pattern Pattern) Pattern {
	return &OneOrMorePattern{Inner: pattern, Greedy: false}
}

func Optional(pattern Pattern) Pattern {
	return &ZeroOrOnePattern{Inner: pattern, Greedy: true}
}

func OptionalLazy(pattern Pattern) Pattern {
	return &ZeroOrOnePattern{Inner: pattern, Greedy: false}
}

func Named(name string, pattern Pattern) Pattern {
	return &NamedPattern{Name: name, Inner: pattern}
}

func MatchPredicate(inner Pattern, predicate func(Expr) bool) Pattern {
	return &PredicatePattern{Inner: inner, Predicate: predicate}
}

func MatchPredicateNamed(inner Pattern, predicate func(Expr) bool, name string) Pattern {
	return &PredicatePattern{Inner: inner, Predicate: predicate, Name: name}
}

// PatternAnalyzer determines the optimal execution strategy for patterns
type PatternAnalyzer struct{}

// Analyze determines whether a pattern can use direct matching or needs NFA
func (pa *PatternAnalyzer) Analyze(pattern Pattern) ExecutionStrategy {
	return pa.analyzePattern(pattern, false)
}

func (pa *PatternAnalyzer) analyzePattern(pattern Pattern, insideQuantifier bool) ExecutionStrategy {
	switch p := pattern.(type) {
	case *LiteralPattern, *HeadPattern, *AnyPattern:
		return StrategyDirect

	case *SequencePattern:
		return pa.analyzeSequence(p.Patterns)

	case *OrPattern:
		// Simple OR patterns (only literals, heads, any) can use Direct matching
		// Complex OR patterns need NFA for branching
		if pa.isSimpleOrPattern(p) {
			return StrategyDirect
		}
		return StrategyNFA

	case *ZeroOrMorePattern, *OneOrMorePattern, *ZeroOrOnePattern:
		// All quantifiers need NFA for proper looping and state management
		return StrategyNFA

	case *NamedPattern:
		// Simple Named patterns can use Direct strategy with enhanced Direct matcher
		// Complex Named patterns (nested, quantifiers) need NFA
		innerStrategy := pa.analyzePattern(p.Inner, insideQuantifier)

		// If inner pattern needs NFA, Named wrapper also needs NFA
		if innerStrategy == StrategyNFA {
			return StrategyNFA
		}

		// Check if this is a complex nested case that needs NFA
		if pa.isComplexNamed(p) {
			return StrategyNFA
		}

		// Simple Named patterns can use enhanced Direct strategy
		return StrategyDirect

	case *PredicatePattern:
		// Predicate patterns can use direct matching if inner pattern can
		return pa.analyzePattern(p.Inner, insideQuantifier)

	case *NotPattern:
		// NOT patterns require NFA
		return StrategyNFA

	default:
		return StrategyNFA
	}
}

func (pa *PatternAnalyzer) analyzeSequence(patterns []Pattern) ExecutionStrategy {
	// New aggressive analysis: Direct matcher can handle most quantifier patterns
	// Only consecutive quantifiers that compete for the same elements need NFA

	var quantifierPositions []int

	for i, p := range patterns {
		// Extract actual pattern (unwrap Named)
		actualPattern := p
		if named, ok := p.(*NamedPattern); ok {
			actualPattern = named.Inner
		}

		switch actualPattern.(type) {
		case *OrPattern:
			// Simple Or patterns can use Direct matching, complex ones need NFA
			if orPattern, ok := actualPattern.(*OrPattern); ok {
				if !pa.isSimpleOrPattern(orPattern) {
					return StrategyNFA
				}
			}
		case *NotPattern:
			// Not patterns always need NFA
			return StrategyNFA
		case *ZeroOrMorePattern, *OneOrMorePattern, *ZeroOrOnePattern:
			quantifierPositions = append(quantifierPositions, i)

			// Check if quantifier inner pattern is complex
			switch q := actualPattern.(type) {
			case *ZeroOrMorePattern:
				if !pa.isSimpleInnerPattern(q.Inner) {
					return StrategyNFA
				}
			case *OneOrMorePattern:
				if !pa.isSimpleInnerPattern(q.Inner) {
					return StrategyNFA
				}
			case *ZeroOrOnePattern:
				if !pa.isSimpleInnerPattern(q.Inner) {
					return StrategyNFA
				}
			}
		default:
			// Check if non-quantifier patterns are complex
			if pa.analyzePattern(p, false) == StrategyNFA {
				return StrategyNFA
			}
		}
	}

	// Check for consecutive quantifiers that could compete
	if pa.hasConsecutiveQuantifiers(quantifierPositions) {
		return StrategyNFA
	}

	// Since we've already checked for consecutive quantifiers above,
	// any remaining quantifier patterns are non-consecutive and unambiguous.
	// The Direct matcher can handle these with proper greedy semantics.

	// Check for complex Named patterns that need NFA
	for _, p := range patterns {
		if named, ok := p.(*NamedPattern); ok {
			if pa.isComplexNamed(named) {
				return StrategyNFA
			}
		}
	}

	// Only trailing quantifiers or no quantifiers - Direct matcher can handle these
	return StrategyDirect
}

// isSimpleInnerPattern checks if a pattern is simple enough for Direct matching
func (pa *PatternAnalyzer) isSimpleInnerPattern(pattern Pattern) bool {
	switch pattern.(type) {
	case *LiteralPattern, *HeadPattern, *AnyPattern:
		return true
	default:
		return false
	}
}

// isSimpleOrPattern checks if an Or pattern contains only simple alternatives
func (pa *PatternAnalyzer) isSimpleOrPattern(orPattern *OrPattern) bool {
	for _, alternative := range orPattern.Alternatives {
		// Unwrap Named patterns first
		actualPattern := alternative
		if named, ok := alternative.(*NamedPattern); ok {
			actualPattern = named.Inner
		}

		// Check if the alternative is a simple pattern
		switch actualPattern.(type) {
		case *LiteralPattern, *HeadPattern, *AnyPattern:
			continue // Simple pattern, OK
		case *PredicatePattern:
			// Predicates with simple inner patterns are OK
			if predPattern, ok := actualPattern.(*PredicatePattern); ok {
				if pa.isSimpleInnerPattern(predPattern.Inner) {
					continue
				}
			}
			return false // Complex predicate
		default:
			return false // Complex pattern (sequences, quantifiers, nested Or, etc.)
		}
	}
	return true
}

// hasConsecutiveQuantifiers checks if there are consecutive quantifiers that could compete
func (pa *PatternAnalyzer) hasConsecutiveQuantifiers(positions []int) bool {
	for i := 0; i < len(positions)-1; i++ {
		if positions[i+1] == positions[i]+1 {
			return true // Consecutive quantifiers found
		}
	}
	return false
}

// isComplexNamed determines if a Named pattern requires NFA due to complexity
func (pa *PatternAnalyzer) isComplexNamed(named *NamedPattern) bool {
	// Check for nested Named patterns (most common complex case)
	if pa.hasNestedNamed(named.Inner) {
		return true
	}

	// Simple Named quantifiers can now use Direct matching
	// Only complex quantifiers need NFA
	switch inner := named.Inner.(type) {
	case *ZeroOrMorePattern:
		return !pa.isSimpleInnerPattern(inner.Inner)
	case *OneOrMorePattern:
		return !pa.isSimpleInnerPattern(inner.Inner)
	case *ZeroOrOnePattern:
		return !pa.isSimpleInnerPattern(inner.Inner)
	}

	// Simple Named patterns can use Direct
	return false
}

// hasNestedNamed checks if a pattern contains nested Named patterns
func (pa *PatternAnalyzer) hasNestedNamed(pattern Pattern) bool {
	switch p := pattern.(type) {
	case *NamedPattern:
		return true
	case *SequencePattern:
		for _, subPattern := range p.Patterns {
			if pa.hasNestedNamed(subPattern) {
				return true
			}
		}
	case *OrPattern:
		for _, subPattern := range p.Alternatives {
			if pa.hasNestedNamed(subPattern) {
				return true
			}
		}
	case *NotPattern:
		return pa.hasNestedNamed(p.Inner)
	case *ZeroOrMorePattern:
		return pa.hasNestedNamed(p.Inner)
	case *OneOrMorePattern:
		return pa.hasNestedNamed(p.Inner)
	case *ZeroOrOnePattern:
		return pa.hasNestedNamed(p.Inner)
	case *PredicatePattern:
		return pa.hasNestedNamed(p.Inner)
	}
	return false
}

// CompilePattern analyzes a pattern and selects the optimal execution strategy
func CompilePattern(pattern Pattern) (Matcher, error) {
	analyzer := &PatternAnalyzer{}
	strategy := analyzer.Analyze(pattern)

	switch strategy {
	case StrategyDirect:
		return &DirectMatcher{pattern: pattern}, nil

	case StrategyNFA:
		// Build Thompson NFA for the pattern
		builder := NewNFABuilder()
		nfa, err := builder.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("failed to compile NFA: %v", err)
		}
		return NewNFAExecutor(nfa), nil

	default:
		return nil, fmt.Errorf("unknown execution strategy: %v", strategy)
	}
}
