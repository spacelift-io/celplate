package evaluator_test

import (
	"testing"
	"time"

	"github.com/franela/goblin"
	. "github.com/onsi/gomega"
	"github.com/spacelift-io/celplate/evaluator"
)

func TestCEL(t *testing.T) {
	g := goblin.Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("CEL", func() {
		var environment map[string]map[string]any
		var err error
		var sut *evaluator.CEL

		g.JustBeforeEach(func() {
			sut, err = evaluator.NewCEL(environment)
		})

		g.Describe("NewCEL", func() {
			g.Describe("with a valid environment", func() {
				environment = map[string]map[string]any{
					"input":   {"foo": "bar"},
					"context": {"time": time.Unix(1666960429, 0), "pi": 3.14},
				}
			})

			g.It("should return a new instance of CEL evaluator", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(sut).ToNot(BeNil())
			})
		})

		g.Describe("Evaluate", func() {
			var expression, result string

			g.JustBeforeEach(func() { result, err = sut.Evaluate(expression) })

			g.Describe("with a valid expression", func() {
				g.BeforeEach(func() {
					expression = `input.foo + "|" + string(context.time.getSeconds()) + "|" + string(context.pi)`
				})

				g.It("should return the result of the expression", func() {
					Expect(err).ToNot(HaveOccurred())
					Expect(result).To(Equal("bar|49|3.14"))
				})
			})

			g.Describe("with an invalid expression type", func() {
				g.BeforeEach(func() { expression = `context.time` })

				g.It("should return an output type error", func() {
					Expect(err).To(MatchError("expected \"2022-10-28 15:33:49 +0300 EEST\" to be of type string but it's google.protobuf.Timestamp"))
					Expect(result).To(BeEmpty())
				})
			})

			g.Describe("with an invalid expression", func() {
				g.BeforeEach(func() { expression = `<<<LLLdsf--dsdf` })

				g.It("should return a compilation error", func() {
					Expect(err.Error()).To(ContainSubstring("line 1, column 1: Syntax error"))
				})
			})

			g.Describe("with an invalid key", func() {
				g.BeforeEach(func() { expression = `unknown.var + input.bar` })

				g.It("should return a compilation error", func() {
					Expect(err).To(MatchError("line 1, column 1: undeclared reference to 'unknown' (in container '')"))
				})
			})
		})
	})
}
