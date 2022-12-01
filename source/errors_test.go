package source_test

import (
	"errors"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/spacelift-io/celplate/source"
)

func TestErrors(t *testing.T) {
	RegisterTestingT(t)

	sut := &source.Errors{}

	// No errors is nil
	Expect(sut.ErrorOrNil()).To(BeNil())

	// Push 3 errors, two of type source
	serr1 := &source.Error{Location: source.Location{Line: 1, Column: 2}, Message: "foo"}
	serr2 := &source.Error{Location: source.Location{Line: 3, Column: 4}, Message: "bar"}
	sut.Push(serr1)
	sut.Push(serr2)
	sut.Push(errors.New("internal error"))

	// Check for the erors
	Expect(sut.ErrorOrNil()).ToNot(BeNil())
	Expect(sut.Error()).To(Equal("line 1, column 2: foo; line 3, column 4: bar; internal error"))

	// Get only the source errors
	result := source.GetErrors(sut)
	Expect(len(result)).To(Equal(2))
	Expect(result[0]).To(Equal(serr1))
	Expect(result[1]).To(Equal(serr2))
}
