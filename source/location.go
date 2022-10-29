package source

import (
	"fmt"
)

const lineBreak = '\n'

// Location represents a location in the source code.
type Location struct {
	Index  int
	Line   int
	Column int
}

// Start creates a new starting location.
func Start() *Location {
	return &Location{Line: 1, Column: 1}
}

// Advance advances the location based on a single character.
func (l *Location) Advance(char rune) {
	l.Index++

	if char != lineBreak {
		l.Column++
		return
	}

	l.Line++
	l.Column = 1
}

// Nested returns a new location that is nested within the current location.
func (l *Location) Nested(nested Location) Location {
	return Location{
		Index:  l.Index + nested.Index,
		Line:   l.Line + nested.Line - 1,
		Column: l.Column + nested.Column - 1,
	}
}

// String returns a string representation of the location.
func (l *Location) String() string {
	return fmt.Sprintf("line %v, column %v", l.Line, l.Column)
}
