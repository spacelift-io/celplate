package celplate

// Evaluator evaluates expressions nested inside supported blocks (${{ ... }}).
type Evaluator interface {
	// Evaluate evaluates the given expression and returns its result, and an
	// error, if any.
	Evaluate(expression string) (string, error)
}
