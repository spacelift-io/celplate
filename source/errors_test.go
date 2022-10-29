package source_test

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/spacelift-io/celplate/source"
)

func TestErrors(t *testing.T) {
	RegisterTestingT(t)

	sut := source.Errors{
		{Location: source.Location{Line: 1, Column: 2}, Message: "foo"},
		{Location: source.Location{Line: 3, Column: 4}, Message: "bar"},
	}

	Expect(sut.Error()).To(Equal("line 1, column 2: foo; line 3, column 4: bar"))
}
