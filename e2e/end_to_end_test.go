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
			"id":          1,
		},
		"context": {
			"datetime": time.Date(2022, time.April, 10, 1, 1, 1, 1, time.UTC),
		},
	})
	Expect(err).ToNot(HaveOccurred())

	input, err := os.ReadFile("fixtures/input.yaml")
	Expect(err).ToNot(HaveOccurred())

	output, err := os.ReadFile("fixtures/output.yaml")
	Expect(err).ToNot(HaveOccurred())

	out, err := celplate.NewScanner(evaluator).Transform(input)

	Expect(err).ToNot(HaveOccurred())
	Expect(string(out)).To(Equal(string(output)))
}
