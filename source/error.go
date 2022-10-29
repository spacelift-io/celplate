package source

import "fmt"

// Error represents an error in the source code at a given location.
type Error struct {
	Location Location
	Message  string
}

func (l *Error) Error() string {
	return fmt.Sprintf("%s: %s", l.Location.String(), l.Message)
}
