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
						"time":     time.Unix(1666960429, 0).UTC(),
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

			g.Describe("custom macros", func() {
				g.Describe("join", func() {
					g.Describe("works with string lists", func() {
						g.BeforeEach(func() {
							expression = `['1', '2'].join(', ')`
						})

						g.It("should return the result of the expression", func() {
							Expect(err).ToNot(HaveOccurred())
							Expect(result).To(Equal("1, 2"))
						})
					})

					g.Describe("works with int lists", func() {
						g.BeforeEach(func() {
							expression = `[1,2,3].join(', ')`
						})

						g.It("should return the result of the expression", func() {
							Expect(err).ToNot(HaveOccurred())
							Expect(result).To(Equal("1, 2, 3"))
						})
					})

					g.Describe("works with int timestamps and durations", func() {
						g.BeforeEach(func() {
							expression = `[timestamp('1972-01-01T10:00:20.021-05:00'), duration('5s')].join(', ')`
						})

						g.It("should return the result of the expression", func() {
							Expect(err).ToNot(HaveOccurred())
							Expect(result).To(Equal("1972-01-01 10:00:20.021 -0500 -0500, 5s"))
						})
					})
				})

				g.Describe("format", func() {
					g.Describe("case one", func() {
						g.Before(func() {
							expression = `timestamp('1972-01-01T10:00:20.021-05:00').format('2006-01-02')`
						})

						g.It("should return the result of the expression", func() {
							Expect(err).ToNot(HaveOccurred())
							Expect(result).To(Equal("1972-01-01"))
						})
					})

					g.Describe("case two", func() {
						g.Before(func() {
							expression = `context.time.format('2006-01-02')`
						})

						g.It("should return the result of the expression", func() {
							Expect(err).ToNot(HaveOccurred())
							Expect(result).To(Equal("2022-10-28"))
						})
					})

				})
			})

			g.Describe("convertible types", func() {
				g.Describe("timestamp", func() {
					g.BeforeEach(func() {
						expression = `timestamp('1972-01-01T10:00:20.021-05:00')`
					})

					g.It("becomes a string", func() {
						Expect(err).ToNot(HaveOccurred())
						Expect(result).To(Equal("1972-01-01T10:00:20.021-05:00"))
					})
				})

				g.Describe("duration", func() {
					g.BeforeEach(func() {
						expression = `duration('1h5s')`
					})

					g.It("becomes a string", func() {
						Expect(err).ToNot(HaveOccurred())
						Expect(result).To(Equal("3605s"))
					})
				})
			})
		})
	})
}
