package evaluator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/ext"

	"github.com/spacelift-io/celplate/source"
)

var anyListType = reflect.TypeOf([]any{})
var anyMapType = reflect.TypeOf(map[any]any{})

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

	// There is a bunch of methods which isn't included in the default environment
	// like `charAt`, `join`, `split`, etc. Let's add them too.
	envOpts = append(envOpts, ext.Strings())

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

	return e.attemptConversionToString(out)
}

// attemptConversionToString tries to convert the outcome of the expression to a string.
func (e *CEL) attemptConversionToString(out ref.Val) (string, error) {
	switch out.Type() {
	// If it's a string, return it as is.
	case types.StringType:
		return out.Value().(string), nil

	// If it's a list, let's convert it to a string like this:
	// [1, 2, 3] -> "[1 2 3]"
	case types.ListType:
		asList, err := out.ConvertToNative(anyListType)
		if err != nil {
			return "", fmt.Errorf("failed to cast value %q of type %s to a list", out.Value(), out.Type().TypeName())
		}

		var items []string
		for _, item := range asList.([]any) {
			items = append(items, fmt.Sprint(item))
		}

		return fmt.Sprintf("[%s]", strings.Join(items, " ")), nil

	// If it's a map, let's convert it to a string like this:
	// {"a": 1, "b": 2} -> "{a: 1, b: 2}"
	case types.MapType:
		asMap, err := out.ConvertToNative(anyMapType)
		if err != nil {
			return "", fmt.Errorf("failed to cast value %q of type %s to a map", out.Value(), out.Type().TypeName())
		}

		var items []string
		for key, value := range asMap.(map[any]any) {
			items = append(items, fmt.Sprintf("%v: %v", key, value))
		}

		return fmt.Sprintf("{%s}", strings.Join(items, ", ")), nil

	// Otherwise, let's attempt a conversion.
	// For example: timestamp("2020-01-01T00:00:00Z") -> "2020-01-01T00:00:00Z" etc.
	default:
		converted := out.ConvertToType(types.StringType)
		if converted.Type() == types.ErrType {
			return "", fmt.Errorf("failed to cast value %q of type %s to a string", out.Value(), out.Type().TypeName())
		}
		return converted.Value().(string), nil
	}
}
