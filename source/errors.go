package source

import (
	"errors"
	"strings"
)

// Errors represents a collection of source code errors.
type Errors struct {
	Errs []error
}

func (e *Errors) Push(err error) {
	if err == nil {
		return
	}

	e.Errs = append(e.Errs, err)
}

func (e *Errors) ErrorOrNil() error {
	if e == nil || len(e.Errs) == 0 {
		return nil
	}

	return e
}

func (e *Errors) Error() string {
	if len(e.Errs) == 1 {
		return e.Errs[0].Error()
	}

	parts := make([]string, len(e.Errs))
	for i, err := range e.Errs {
		parts[i] = err.Error()
	}

	return strings.Join(parts, "; ")
}

func GetErrors(err error) []*Error {
	result := []*Error{}

	var errs *Errors
	if errors.As(err, &errs) {
		for _, e := range errs.Errs {
			var src *Error
			if errors.As(e, &src) {
				result = append(result, src)
			}
		}
	}

	return result
}
