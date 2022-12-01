package celplate_test

import (
	"errors"
	"testing"

	"github.com/franela/goblin"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"github.com/spacelift-io/celplate"
)

type mockEvaluator struct {
	mock.Mock
}

func (m *mockEvaluator) Evaluate(expression string) (string, error) {
	args := m.Called(expression)
	return args.String(0), args.Error(1)
}

func TestScanner(t *testing.T) {
	g := goblin.Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("Scanner", func() {
		var sut *celplate.Scanner
		var mockCaller *mockEvaluator

		g.BeforeEach(func() {
			mockCaller = new(mockEvaluator)
			sut = celplate.NewScanner(mockCaller)
		})

		g.Describe("Transform", func() {
			var input []byte
			var err error
			var output []byte

			g.JustBeforeEach(func() { output, err = sut.Transform(input) })

			g.Describe("when the input is valid", func() {
				g.BeforeEach(func() { input = []byte("Hello, world!") })

				g.It("should succeed", func() {
					Expect(err).NotTo(HaveOccurred())

					Expect(err).NotTo(HaveOccurred())
					Expect(output).To(Equal(input))
				})
			})

			g.Describe("when the input is invalid", func() {
				g.BeforeEach(func() { input = []byte("Hello, ${{ world }!") })

				g.It("should fail with the location", func() {
					Expect(err).To(MatchError("line 1, column 19: unexpected character '!', expected '}'; line 1, column 20: unexpected end of input"))
					Expect(output).To(BeNil())
				})
			})

			g.Describe("when the input contains an expression", func() {
				var evaluateCall *mock.Call

				g.BeforeEach(func() {
					input = []byte("Hello, ${{ world }}!")
					evaluateCall = mockCaller.On("Evaluate", " world ")
				})

				g.Describe("when the evaluation fails", func() {
					g.BeforeEach(func() { evaluateCall.Return("", errors.New("error")) })

					g.It("should fail with the location", func() {
						Expect(err).To(MatchError("error; line 1, column 20: unexpected character '!', expected '}'; line 1, column 21: unexpected end of input"))
						Expect(output).To(BeNil())
					})
				})

				g.Describe("when the evaluation succeeds", func() {
					g.BeforeEach(func() { evaluateCall.Return("world", nil) })

					g.It("should succeed", func() {
						Expect(err).NotTo(HaveOccurred())
						Expect(output).To(Equal([]byte("Hello, world!")))
					})
				})
			})
		})
	})
}
