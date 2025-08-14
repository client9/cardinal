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

// CompiledPattern represents a compiled pattern ready for execution
type CompiledPattern struct {
	strategy      ExecutionStrategy
	directMatcher *DirectMatcher // Fast path for simple patterns
	nfaExecutor   *NFAExecutor   // Full NFA for complex patterns
}

// Match executes the pattern against an expression
func (cp *CompiledPattern) Match(expr Expr) MatchResult {
	switch cp.strategy {
	case StrategyDirect:
		return cp.directMatcher.Match(expr)
	case StrategyNFA:
		return cp.nfaExecutor.Match(expr)
	default:
		return MatchResult{Matched: false}
	}
}

// Strategy returns the execution strategy used by this compiled pattern
func (cp *CompiledPattern) Strategy() ExecutionStrategy {
	return cp.strategy
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
		switch q := lastPattern.(type) {
		case *ZeroOrMorePattern:
			return dm.matchWithTrailingZeroOrMoreBindings(elements, patterns[:len(patterns)-1], q.Inner, bindings)
		case *OneOrMorePattern:
			return dm.matchWithTrailingOneOrMoreBindings(elements, patterns[:len(patterns)-1], q.Inner, bindings)
		case *ZeroOrOnePattern:
			return dm.matchWithTrailingZeroOrOneBindings(elements, patterns[:len(patterns)-1], q.Inner, bindings)
		}
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

// StateSet represents a set of active NFA states
type StateSet map[int]bool

func (ne *NFAExecutor) matchSequence(exprs []Expr) MatchResult {
	// Track group boundaries and captured elements
	groupCaptures := make(map[string][]Expr) // Track elements captured within each group
	activeGroups := make(map[int][]string)   // Track which groups are active for each state

	// Initialize with epsilon closure of start state
	currentStates := make(StateSet)
	ne.addStateWithEpsilonAndGroups(currentStates, ne.nfa.StartState, activeGroups, groupCaptures)

	// Process each expression
	for _, expr := range exprs {
		currentStates = ne.stepWithGroups(currentStates, expr, activeGroups, groupCaptures)
		if len(currentStates) == 0 {
			// No valid states remaining
			return MatchResult{Matched: false, Bindings: ne.buildBindings(groupCaptures)}
		}
	}

	// Check if any current state is an accept state
	matched := false
	for stateID := range currentStates {
		if ne.nfa.States[stateID].IsAccept {
			matched = true
			break
		}
	}

	return MatchResult{
		Matched:  matched,
		Bindings: ne.buildBindings(groupCaptures),
		Consumed: len(exprs),
	}
}

func (ne *NFAExecutor) step(currentStates StateSet, expr Expr, bindings map[string]Expr) StateSet {
	nextStates := make(StateSet)

	// For each active state, follow matching transitions
	for stateID := range currentStates {
		state := ne.nfa.States[stateID]

		for _, transition := range state.Transitions {
			if transition.Type == MatchTransition {
				if transition.Condition.Matches(expr, bindings) {
					ne.addStateWithEpsilon(nextStates, transition.Target, bindings)
				}
			}
		}
	}

	return nextStates
}

func (ne *NFAExecutor) addStateWithEpsilon(states StateSet, stateID int, bindings map[string]Expr) {
	if states[stateID] {
		return // Already added
	}

	states[stateID] = true

	// Follow epsilon transitions
	state := ne.nfa.States[stateID]
	for _, transition := range state.Transitions {
		if transition.Type == EpsilonTransition {
			ne.addStateWithEpsilon(states, transition.Target, bindings)
		}
	}
}

func (ne *NFAExecutor) addStateWithEpsilonAndGroups(states StateSet, stateID int, activeGroups map[int][]string, groupCaptures map[string][]Expr) {
	if states[stateID] {
		return // Already added
	}

	states[stateID] = true
	state := ne.nfa.States[stateID]

	// Track groups for this state path
	var currentGroups []string
	if existing, ok := activeGroups[stateID]; ok {
		currentGroups = make([]string, len(existing))
		copy(currentGroups, existing)
	}

	// Handle group start
	if state.GroupStart != "" {
		currentGroups = append(currentGroups, state.GroupStart)
		// Initialize capture slice for this group if not exists
		if _, exists := groupCaptures[state.GroupStart]; !exists {
			groupCaptures[state.GroupStart] = make([]Expr, 0)
		}
	}

	// Handle group end
	if state.GroupEnd != "" {
		// Remove from active groups
		for i, group := range currentGroups {
			if group == state.GroupEnd {
				currentGroups = append(currentGroups[:i], currentGroups[i+1:]...)
				break
			}
		}
	}

	// Instead of overwriting, merge with any existing groups for this state
	if existing, exists := activeGroups[stateID]; exists {
		// Merge current groups with existing groups (avoid duplicates)
		merged := append([]string{}, existing...)
		for _, group := range currentGroups {
			found := false
			for _, existingGroup := range merged {
				if existingGroup == group {
					found = true
					break
				}
			}
			if !found {
				merged = append(merged, group)
			}
		}
		activeGroups[stateID] = merged
	} else {
		activeGroups[stateID] = currentGroups
	}

	// Follow epsilon transitions
	for _, transition := range state.Transitions {
		if transition.Type == EpsilonTransition {
			// Epsilon transition - manually propagate groups to target first

			// Ensure target state inherits our groups before recursing
			if len(currentGroups) > 0 {
				if existing, exists := activeGroups[transition.Target]; exists {
					// Merge groups
					merged := append([]string{}, existing...)
					for _, group := range currentGroups {
						found := false
						for _, existingGroup := range merged {
							if existingGroup == group {
								found = true
								break
							}
						}
						if !found {
							merged = append(merged, group)
						}
					}
					activeGroups[transition.Target] = merged
				} else {
					activeGroups[transition.Target] = append([]string{}, currentGroups...)
				}
			}

			ne.addStateWithEpsilonAndGroups(states, transition.Target, activeGroups, groupCaptures)
		}
	}
}

func (ne *NFAExecutor) stepWithGroups(currentStates StateSet, expr Expr, activeGroups map[int][]string, groupCaptures map[string][]Expr) StateSet {
	nextStates := make(StateSet)

	// For each active state, follow matching transitions
	for stateID := range currentStates {
		state := ne.nfa.States[stateID]

		for _, transition := range state.Transitions {
			if transition.Type == MatchTransition {
				if transition.Condition.Matches(expr, make(map[string]Expr)) { // Use dummy bindings since we handle groups separately
					// Capture this element in all active groups for this state
					if groups, ok := activeGroups[stateID]; ok {
						for _, groupName := range groups {
							groupCaptures[groupName] = append(groupCaptures[groupName], expr)
						}
					}

					// Propagate groups from source state to target state during match transition
					if sourceGroups, hasGroups := activeGroups[stateID]; hasGroups && len(sourceGroups) > 0 {
						if existing, exists := activeGroups[transition.Target]; exists {
							// Merge groups
							merged := append([]string{}, existing...)
							for _, group := range sourceGroups {
								found := false
								for _, existingGroup := range merged {
									if existingGroup == group {
										found = true
										break
									}
								}
								if !found {
									merged = append(merged, group)
								}
							}
							activeGroups[transition.Target] = merged
						} else {
							activeGroups[transition.Target] = append([]string{}, sourceGroups...)
						}
					}

					ne.addStateWithEpsilonAndGroups(nextStates, transition.Target, activeGroups, groupCaptures)
				}
			}
		}
	}

	return nextStates
}

func (ne *NFAExecutor) buildBindings(groupCaptures map[string][]Expr) map[string]Expr {
	bindings := make(map[string]Expr)

	for groupName, captures := range groupCaptures {
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
		// OR always needs NFA for branching
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
	hasQuantifier := false
	quantifierAtEnd := false

	for i, p := range patterns {
		// Helper function to check if a pattern contains a quantifier
		// (either directly or wrapped in Named)
		actualPattern := p
		if named, ok := p.(*NamedPattern); ok {
			actualPattern = named.Inner
		}

		switch actualPattern.(type) {
		case *ZeroOrMorePattern, *OneOrMorePattern, *ZeroOrOnePattern:
			hasQuantifier = true
			quantifierAtEnd = (i == len(patterns)-1)
		case *OrPattern, *NotPattern:
			// Complex patterns always need NFA
			return StrategyNFA
		}

		// For quantifier patterns, check if they can be handled directly
		switch q := actualPattern.(type) {
		case *ZeroOrMorePattern:
			if i == len(patterns)-1 {
				// For trailing quantifiers, check if inner pattern is simple
				switch q.Inner.(type) {
				case *LiteralPattern, *HeadPattern, *AnyPattern:
					// Simple inner patterns are fine for trailing quantifiers
				default:
					return StrategyNFA
				}
			} else {
				// Quantifier in middle always needs NFA
				return StrategyNFA
			}
		case *OneOrMorePattern:
			if i == len(patterns)-1 {
				// For trailing quantifiers, check if inner pattern is simple
				switch q.Inner.(type) {
				case *LiteralPattern, *HeadPattern, *AnyPattern:
					// Simple inner patterns are fine for trailing quantifiers
				default:
					return StrategyNFA
				}
			} else {
				// Quantifier in middle always needs NFA
				return StrategyNFA
			}
		case *ZeroOrOnePattern:
			if i == len(patterns)-1 {
				// For trailing quantifiers, check if inner pattern is simple
				switch q.Inner.(type) {
				case *LiteralPattern, *HeadPattern, *AnyPattern:
					// Simple inner patterns are fine for trailing quantifiers
				default:
					return StrategyNFA
				}
			} else {
				// Quantifier in middle always needs NFA
				return StrategyNFA
			}
		default:
			// Check if non-quantifier patterns are complex
			if pa.analyzePattern(p, false) == StrategyNFA {
				return StrategyNFA
			}
		}
	}

	// Only complex Named patterns require NFA for state-tagged group boundaries
	// Simple Named patterns in sequences can use Direct strategy
	for _, p := range patterns {
		if named, ok := p.(*NamedPattern); ok {
			if pa.isComplexNamed(named) {
				return StrategyNFA
			}
		}
	}

	// If quantifier is only at the end, we can use direct matching
	if hasQuantifier && quantifierAtEnd {
		return StrategyDirect
	}

	// If quantifier anywhere else, need NFA
	if hasQuantifier && !quantifierAtEnd {
		return StrategyNFA
	}

	// Pure sequence of simple patterns - direct match
	return StrategyDirect
}

// isComplexNamed determines if a Named pattern requires NFA due to complexity
func (pa *PatternAnalyzer) isComplexNamed(named *NamedPattern) bool {
	// Check for nested Named patterns (most common complex case)
	if pa.hasNestedNamed(named.Inner) {
		return true
	}

	// Named quantifiers need NFA for proper collection semantics
	switch named.Inner.(type) {
	case *ZeroOrMorePattern, *OneOrMorePattern, *ZeroOrOnePattern:
		return true
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
func CompilePattern(pattern Pattern) (*CompiledPattern, error) {
	analyzer := &PatternAnalyzer{}
	strategy := analyzer.Analyze(pattern)

	cp := &CompiledPattern{strategy: strategy}

	switch strategy {
	case StrategyDirect:
		cp.directMatcher = &DirectMatcher{pattern: pattern}

	case StrategyNFA:
		// Build Thompson NFA for the pattern
		builder := NewNFABuilder()
		nfa, err := builder.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("failed to compile NFA: %v", err)
		}
		cp.nfaExecutor = NewNFAExecutor(nfa)

	default:
		return nil, fmt.Errorf("unknown execution strategy: %v", strategy)
	}

	return cp, nil
}
