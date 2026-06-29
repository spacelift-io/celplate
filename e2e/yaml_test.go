package e2e_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/spacelift-io/celplate"
	"github.com/spacelift-io/celplate/evaluator"
)

func TestEndToEnd(t *testing.T) {
	cel, err := evaluator.NewCEL(map[string]map[string]any{
		"inputs": {
			"environment": "production",
			"region":      "us-east-1",
			"id":          1,
			"serial":      111111111111,
		},
		"context": {
			"datetime": time.Date(2022, time.April, 10, 1, 1, 1, 1, time.UTC),
		},
	})
	require.NoError(t, err)

	input, err := os.ReadFile("fixtures/input.yaml")
	require.NoError(t, err)

	expected, err := os.ReadFile("fixtures/output.yaml")
	require.NoError(t, err)

	out, err := celplate.NewScanner(cel).Transform(input)
	require.NoError(t, err)
	assert.Equal(t, string(expected), string(out))
}
