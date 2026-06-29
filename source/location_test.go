package source_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/spacelift-io/celplate/source"
)

func TestLocation_Start(t *testing.T) {
	sut := source.Start()
	assert.Zero(t, sut.Index)
	assert.Equal(t, "line 1, column 1", sut.String())
}

func TestLocation_Advance_RegularCharacter(t *testing.T) {
	sut := source.Start()
	sut.Advance('a')
	assert.Equal(t, 1, sut.Index)
	assert.Equal(t, "line 1, column 2", sut.String())
}

func TestLocation_Advance_LineBreak(t *testing.T) {
	sut := source.Start()
	sut.Advance('\n')
	assert.Equal(t, 1, sut.Index)
	assert.Equal(t, "line 2, column 1", sut.String())
}

func TestLocation_Nested(t *testing.T) {
	sut := &source.Location{Index: 1, Line: 2, Column: 3}
	nested := source.Location{Index: 1, Line: 2, Column: 3}
	assert.Equal(t, source.Location{Index: 2, Line: 3, Column: 5}, sut.Nested(nested))
}
