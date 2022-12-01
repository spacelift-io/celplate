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
					"input": {
						"foo": "bar",
					},
					"context": {
						"time":     time.Unix(1666960429, 0),
						"pi":       3.14,
						"unsigned": uint(1),
						"signed":   2,
						"boolean":  true,
					},
					"invalid": {
						"map":   map[int]int{1: 2},
						"slice": []int{1, 2},
						"func":  func() {},
					},
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

			g.Describe("with all expression values converted to string", func() {
				g.BeforeEach(func() {
					expression = `input.foo + "|" + string(context.time.getSeconds()) + "|" + string(context.pi)`
				})

				g.It("should return the result of the expression", func() {
					Expect(err).ToNot(HaveOccurred())
					Expect(result).To(Equal("bar|49|3.14"))
				})
			})

			g.Describe("with non string input values being converted to appropriate strings", func() {
				g.Describe("double", func() {
					g.BeforeEach(func() {
						expression = `context.pi`
					})

					g.It("should return the result of the expression", func() {
						Expect(err).ToNot(HaveOccurred())
						Expect(result).To(Equal("3.14"))
					})
				})
				g.Describe("int", func() {
					g.BeforeEach(func() {
						expression = `context.signed`
					})

					g.It("should return the result of the expression", func() {
						Expect(err).ToNot(HaveOccurred())
						Expect(result).To(Equal("2"))
					})
				})
				g.Describe("uint", func() {
					g.BeforeEach(func() {
						expression = `context.unsigned`
					})

					g.It("should return the result of the expression", func() {
						Expect(err).ToNot(HaveOccurred())
						Expect(result).To(Equal("1"))
					})
				})
				g.Describe("boolean", func() {
					g.BeforeEach(func() {
						expression = `context.boolean`
					})

					g.It("should return the result of the expression", func() {
						Expect(err).ToNot(HaveOccurred())
						Expect(result).To(Equal("true"))
					})
				})
				g.Describe("boolean", func() {
					g.BeforeEach(func() {
						expression = `context.boolean`
					})

					g.It("should return the result of the expression", func() {
						Expect(err).ToNot(HaveOccurred())
						Expect(result).To(Equal("true"))
					})
				})
				g.Describe("slice", func() {
					g.BeforeEach(func() {
						expression = `invalid.slice`
					})

					g.It("should fail as it cannot be a string", func() {
						Expect(err).To(HaveOccurred())
						Expect(result).To(Equal(""))
					})
				})
				g.Describe("map", func() {
					g.BeforeEach(func() {
						expression = `invalid.map`
					})

					g.It("should fail as it cannot be a string", func() {
						Expect(err).To(HaveOccurred())
						Expect(result).To(Equal(""))
					})
				})
				g.Describe("function", func() {
					g.BeforeEach(func() {
						expression = `invalid.func`
					})

					g.It("should fail as it cannot be a string", func() {
						Expect(err).To(HaveOccurred())
						Expect(result).To(Equal(""))
					})
				})
			})

			g.Describe("with an invalid expression type", func() {
				g.BeforeEach(func() { expression = `context.time` })

				g.It("should return an output type error", func() {
					Expect(err).To(MatchError("failed to cast value \"2022-10-28 15:33:49 +0300 EEST\" of type google.protobuf.Timestamp to a string"))
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
