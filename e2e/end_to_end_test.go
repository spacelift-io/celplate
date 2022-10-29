package e2e_test

import (
	"os"
	"testing"
	"time"

	. "github.com/onsi/gomega"

	"github.com/spacelift-io/celplate"
	"github.com/spacelift-io/celplate/evaluator"
)

func TestEndToEnd(t *testing.T) {
	RegisterTestingT(t)

	evaluator, err := evaluator.NewCEL(map[string]map[string]any{
		"inputs": {
			"environment": "production",
			"region":      "us-east-1",
		},
		"context": {
			"now": time.Now(),
		},
	})

	Expect(err).ToNot(HaveOccurred())

	input, err := os.ReadFile("fixtures/blueprint.yaml")
	Expect(err).ToNot(HaveOccurred())

	scanner := celplate.NewScanner(evaluator)

	out, err := scanner.Transform(input)
	Expect(err).ToNot(HaveOccurred())

	Expect(string(out)).To(Equal(`bacon`))
}
