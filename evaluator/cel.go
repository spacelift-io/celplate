package evaluator

import (
	"fmt"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/parser"

	"github.com/spacelift-io/celplate/macros"
	"github.com/spacelift-io/celplate/source"
)

var customMacros = []parser.Macro{
	macros.GetJoinMacro(),
}

// CEL is an implementation of Evaluator that uses CEL expressions.
type CEL struct {
	env  *cel.Env
	vars map[string]any
}

// NewCEL returns a new instance of CEL evaluator.
func NewCEL(data map[string]map[string]any) (*CEL, error) {
	var envOpts []cel.EnvOption

	vars := make(map[string]any)

	for key, value := range data {
		envOpts = append(
			envOpts,
			cel.Variable(key, cel.MapType(cel.StringType, cel.AnyType)),
		)
		vars[key] = value
	}

	for _, macro := range customMacros {
		envOpts = append(envOpts, cel.Macros(macro))
	}

	env, err := cel.NewEnv(envOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create environment: %w", err)
	}

	return &CEL{env, vars}, nil
}

// Evaluate evaluates the given expression using Google CEL, and returns its
// result, plus an error, if any.
//
// It expects the final value to be one of types: "string", "int", "uint", "double", "bool".
func (e *CEL) Evaluate(expression string) (string, error) {
	ast, iss := e.env.Compile(expression)

	if errors := iss.Errors(); len(errors) > 0 {
		sourceErrors := &source.Errors{}

		for _, err := range errors {
			sourceErrors.Push(&source.Error{
				Location: source.Location{
					Line:   err.Location.Line(),
					Column: err.Location.Column() + 1, // CEL columns are 0-based
				},
				Message: err.Message,
			})
		}

		return "", sourceErrors.ErrorOrNil()
	}

	program, err := e.env.Program(ast)
	if err != nil {
		return "", fmt.Errorf("failed to create expression evaluator %w", err)
	}

	out, _, err := program.Eval(e.vars)
	if err != nil {
		return "", fmt.Errorf("failed to evaluate expression: %w", err)
	}

	// If it's a string, return it as is.
	if out.Type() == types.StringType {
		return out.Value().(string), nil
	}

	// Otherwise, let's attempt a conversion.
	converted := out.ConvertToType(types.StringType)
	if converted.Type() == types.ErrType {
		return "", fmt.Errorf("failed to cast value %q of type %s to a string", out.Value(), out.Type().TypeName())
	}

	return converted.Value().(string), nil
}
