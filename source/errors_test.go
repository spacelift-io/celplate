package source_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/spacelift-io/celplate/source"
)

func TestErrors(t *testing.T) {
	sut := &source.Errors{}

	require.NoError(t, sut.ErrorOrNil())

	serr1 := &source.Error{Location: source.Location{Line: 1, Column: 2}, Message: "foo"}
	serr2 := &source.Error{Location: source.Location{Line: 3, Column: 4}, Message: "bar"}
	sut.Push(serr1)
	sut.Push(serr2)
	sut.Push(errors.New("internal error"))

	require.Error(t, sut.ErrorOrNil())
	assert.Equal(t, "line 1, column 2: foo; line 3, column 4: bar; internal error", sut.Error())

	result := source.GetErrors(sut)
	require.Len(t, result, 2)
	assert.Equal(t, serr1, result[0])
	assert.Equal(t, serr2, result[1])
}
