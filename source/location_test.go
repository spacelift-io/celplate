package source_test

import (
	"testing"

	"github.com/franela/goblin"
	. "github.com/onsi/gomega"

	"github.com/spacelift-io/celplate/source"
)

func TestLocation(t *testing.T) {
	g := goblin.Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	var sut *source.Location

	g.Describe("Location", func() {
		g.BeforeEach(func() { sut = source.Start() })

		g.Describe("Start", func() {
			g.It("should return a starting location", func() {
				Expect(sut.Index).To(BeZero())
				Expect(sut.String()).To(Equal("line 1, column 1"))
			})
		})

		g.Describe("Advance", func() {
			var char rune

			g.JustBeforeEach(func() { sut.Advance(char) })

			g.Describe("with a regular character", func() {
				g.BeforeEach(func() { char = 'a' })

				g.It("should advance the column", func() {
					Expect(sut.Index).To(Equal(1))
					Expect(sut.String()).To(Equal("line 1, column 2"))
				})
			})

			g.Describe("with a line break", func() {
				g.BeforeEach(func() { char = '\n' })

				g.It("should advance the line", func() {
					Expect(sut.Index).To(Equal(1))
					Expect(sut.String()).To(Equal("line 2, column 1"))
				})
			})
		})

		g.Describe("Nested", func() {
			var nested source.Location

			g.BeforeEach(func() {
				sut = &source.Location{Index: 1, Line: 2, Column: 3}
			})

			g.JustBeforeEach(func() { nested = source.Location{Index: 1, Line: 2, Column: 3} })

			g.It("should return a new location that is nested within the current location", func() {
				Expect(sut.Nested(nested)).To(Equal(source.Location{Index: 2, Line: 3, Column: 5}))
			})
		})
	})
}
