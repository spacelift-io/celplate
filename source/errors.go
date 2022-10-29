package source

import "strings"

// Errors represents a collection of source code errors.
type Errors []Error

func (e Errors) Error() string {
	parts := make([]string, len(e))

	for i, err := range e {
		parts[i] = err.Error()
	}

	return strings.Join(parts, "; ")
}
