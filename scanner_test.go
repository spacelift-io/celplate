package celplate_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/spacelift-io/celplate"
)

type mockEvaluator struct {
	mock.Mock
}

func (m *mockEvaluator) Evaluate(expression string) (string, error) {
	args := m.Called(expression)
	return args.String(0), args.Error(1)
}

func TestScanner_Transform_PlainText(t *testing.T) {
	sut := celplate.NewScanner(new(mockEvaluator))
	input := []byte("Hello, world!")

	output, err := sut.Transform(input)

	require.NoError(t, err)
	assert.Equal(t, input, output)
}

func TestScanner_Transform_InvalidExpression(t *testing.T) {
	sut := celplate.NewScanner(new(mockEvaluator))

	output, err := sut.Transform([]byte("Hello, ${{ world }!"))

	assert.EqualError(t, err, "line 1, column 19: unexpected character '!', expected '}'")
	assert.Nil(t, output)
}

func TestScanner_Transform_BareDollarAtEndOfLine(t *testing.T) {
	sut := celplate.NewScanner(new(mockEvaluator))
	input := []byte("pattern: ^[a-z]+$\nrequired: true")

	output, err := sut.Transform(input)

	require.NoError(t, err)
	assert.Equal(t, input, output)
}

func TestScanner_Transform_ExpressionEvaluationError(t *testing.T) {
	ev := new(mockEvaluator)
	ev.On("Evaluate", " world ").Return("", errors.New("error"))
	sut := celplate.NewScanner(ev)

	output, err := sut.Transform([]byte("Hello, ${{ world }}!"))

	assert.EqualError(t, err, "line 1, column 19: error")
	assert.Nil(t, output)
}

func TestScanner_Transform_ExpressionEvaluationSuccess(t *testing.T) {
	ev := new(mockEvaluator)
	ev.On("Evaluate", " world ").Return("world", nil)
	sut := celplate.NewScanner(ev)

	output, err := sut.Transform([]byte("Hello, ${{ world }}!"))

	require.NoError(t, err)
	assert.Equal(t, []byte("Hello, world!"), output)
}
