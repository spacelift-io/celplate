package e2e_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/spacelift-io/celplate"
	"github.com/spacelift-io/celplate/evaluator"
	"github.com/spacelift-io/celplate/source"
)

// TestMultilineErrorReporting verifies that errors on different lines
// report the correct line numbers in error messages.
func TestMultilineErrorReporting(t *testing.T) {
	// Create an evaluator with some inputs, but missing the ones we'll reference
	eval, err := evaluator.NewCEL(map[string]map[string]any{
		"inputs": {
			"existing": "value",
		},
	})
	require.NoError(t, err)

	// Test case 1: Errors on different lines should report correct line numbers
	t.Run("errors on different lines", func(t *testing.T) {
		input := `first line # with a trailing comment
second line with ${{ inputs.missing1 }} expression
third line
fourth line with ${{ inputs.missing2 }} expression`

		scanner := celplate.NewScanner(eval)
		_, err := scanner.Transform([]byte(input))

		require.Error(t, err)

		// Check that we get source.Errors with multiple errors
		errs := source.GetErrors(err)
		assert.Len(t, errs, 2)

		// First error should be on line 2
		assert.Equal(t, 2, errs[0].Location.Line)
		assert.Contains(t, errs[0].Message, "no such key: missing1")

		// Second error should be on line 4
		assert.Equal(t, 4, errs[1].Location.Line)
		assert.Contains(t, errs[1].Message, "no such key: missing2")
	})

	// Test case 2: Multiple errors on the same line
	t.Run("multiple errors on same line", func(t *testing.T) {
		input := `first line
second line with ${{ inputs.missing1 }} and ${{ inputs.missing2 }}
third line`

		scanner := celplate.NewScanner(eval)
		_, err := scanner.Transform([]byte(input))

		require.Error(t, err)

		errs := source.GetErrors(err)
		assert.Len(t, errs, 2)

		// Both errors should be on line 2
		assert.Equal(t, 2, errs[0].Location.Line)
		assert.Equal(t, 2, errs[1].Location.Line)

		// Column numbers should be different
		assert.NotEqual(t, errs[0].Location.Column, errs[1].Location.Column)
	})

	// Test case 3: Error after comment lines
	t.Run("error after comment lines", func(t *testing.T) {
		input := `# This is a comment
# Another comment line
third line with ${{ inputs.missing }} expression`

		scanner := celplate.NewScanner(eval)
		_, err := scanner.Transform([]byte(input))

		require.Error(t, err)

		errs := source.GetErrors(err)
		assert.Len(t, errs, 1)

		// Error should be on line 3 (after two comment lines)
		assert.Equal(t, 3, errs[0].Location.Line)
		assert.Contains(t, errs[0].Message, "no such key: missing")
	})

	// Test case 4: Complex YAML-like structure
	t.Run("complex yaml structure", func(t *testing.T) {
		input := `stacks:
  - name: ${{ inputs.missing_name }}-stack
    description: ${{ inputs.missing_desc }}
    config:
      region: us-east-1
      environment: ${{ inputs.missing_env }}`

		scanner := celplate.NewScanner(eval)
		_, err := scanner.Transform([]byte(input))

		require.Error(t, err)

		errs := source.GetErrors(err)
		assert.Len(t, errs, 3)

		// Verify line numbers for each error
		assert.Equal(t, 2, errs[0].Location.Line, "first error should be on line 2")
		assert.Contains(t, errs[0].Message, "missing_name")

		assert.Equal(t, 3, errs[1].Location.Line, "second error should be on line 3")
		assert.Contains(t, errs[1].Message, "missing_desc")

		assert.Equal(t, 6, errs[2].Location.Line, "third error should be on line 6")
		assert.Contains(t, errs[2].Message, "missing_env")
	})

	// Test case 5: Verify error message format
	t.Run("error message format", func(t *testing.T) {
		input := `line one
line two with ${{ inputs.missing }} error
line three`

		scanner := celplate.NewScanner(eval)
		_, err := scanner.Transform([]byte(input))

		require.Error(t, err)

		// The error message should contain the line number
		assert.Contains(t, err.Error(), "line 2")

		errs := source.GetErrors(err)
		assert.Len(t, errs, 1)

		// The formatted error should show "line 2, column X"
		errorStr := errs[0].Error()
		assert.True(t, strings.HasPrefix(errorStr, "line 2, column"))
	})
}
