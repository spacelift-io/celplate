package evaluator_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/spacelift-io/celplate/evaluator"
)

func newTestCEL(t *testing.T) *evaluator.CEL {
	t.Helper()
	cel, err := evaluator.NewCEL(map[string]map[string]any{
		"input": {
			"foo": "bar",
		},
		"context": {
			"time":     time.Unix(1666960429, 0).UTC(),
			"pi":       3.14,
			"unsigned": uint(1),
			"signed":   2,
			"boolean":  true,
		},
		"complex": {
			"intmap":   map[any]any{1: 2},
			"mixedmap": map[any]any{1: "2"},
			"slice":    []int{1, 2},
		},
		"invalid": {
			"func": func() {},
		},
	})
	require.NoError(t, err)
	return cel
}

func TestCEL_NewCEL(t *testing.T) {
	cel, err := evaluator.NewCEL(map[string]map[string]any{
		"input": {"foo": "bar"},
	})
	require.NoError(t, err)
	assert.NotNil(t, cel)
}

func TestCEL_Evaluate(t *testing.T) {
	tests := []struct {
		name        string
		expression  string
		want        string
		wantErr     bool
		errContains string
	}{
		{
			name:       "all values converted to string",
			expression: `input.foo + "|" + string(context.time.getSeconds()) + "|" + string(context.pi)`,
			want:       "bar|49|3.14",
		},
		{
			name:       "double",
			expression: `context.pi`,
			want:       "3.14",
		},
		{
			name:       "int",
			expression: `context.signed`,
			want:       "2",
		},
		{
			name:       "uint",
			expression: `context.unsigned`,
			want:       "1",
		},
		{
			name:       "boolean",
			expression: `context.boolean`,
			want:       "true",
		},
		{
			name:       "slice joined properly",
			expression: `complex.slice`,
			want:       "[1 2]",
		},
		{
			name:       "intmap formatted properly",
			expression: `complex.intmap`,
			want:       "{1: 2}",
		},
		{
			name:       "mixedmap formatted properly",
			expression: `complex.mixedmap`,
			want:       "{1: 2}",
		},
		{
			name:    "function cannot be a string",
			expression: `invalid.func`,
			wantErr: true,
		},
		{
			name:        "invalid expression returns compilation error",
			expression:  `<<<LLLdsf--dsdf`,
			wantErr:     true,
			errContains: "line 1, column 1: Syntax error",
		},
		{
			name:        "invalid key returns compilation error",
			expression:  `unknown.var + input.bar`,
			wantErr:     true,
			errContains: "line 1, column 1: undeclared reference to 'unknown' (in container '')",
		},
		{
			name:       "join macro works with string lists",
			expression: `['1', '2'].join(', ')`,
			want:       "1, 2",
		},
		{
			name:       "timestamp becomes a string",
			expression: `timestamp('1972-01-01T10:00:20.021-05:00')`,
			want:       "1972-01-01T10:00:20.021-05:00",
		},
		{
			name:       "duration becomes a string",
			expression: `duration('1h5s')`,
			want:       "3605s",
		},
		{
			name:       "split becomes a list",
			expression: `"hello world".split(" ").join(", ")`,
			want:       "hello, world",
		},
		{
			name:       "replace replaces the string",
			expression: `'hello hello'.replace('he', 'we')`,
			want:       "wello wello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cel := newTestCEL(t)
			result, err := cel.Evaluate(tt.expression)
			if tt.wantErr {
				require.Error(t, err)
				assert.Empty(t, result)
				if tt.errContains != "" {
					assert.ErrorContains(t, err, tt.errContains)
				}
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, result)
		})
	}
}
