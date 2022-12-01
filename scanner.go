package celplate

import (
	"bytes"
	"fmt"

	"github.com/spacelift-io/celplate/source"
)

type Scanner struct {
	currentExpression      *bytes.Buffer
	currentExpressionStart source.Location
	output                 *bytes.Buffer

	state    scannerState
	location *source.Location

	evaluator Evaluator
}

// Evaluator evaluates expressions nested inside supported blocks (${{ ... }}).
type Evaluator interface {
	// Evaluate evaluates the given expression and returns its result, and an
	// error, if any.
	Evaluate(expression string) (string, error)
}

type scannerState int

const (
	dollarChar = '$'
	openChar   = '{'
	closeChar  = '}'

	ssDefault scannerState = iota
	ssDollar
	ssOpenChar
	ssExpression
	ssCloseChar
)

// NewScanner returns a new generic `Scanner` object.
func NewScanner(evaluator Evaluator) *Scanner {
	return &Scanner{
		currentExpression: bytes.NewBuffer(nil),
		output:            bytes.NewBuffer(nil),
		state:             ssDefault,
		evaluator:         evaluator,
		location:          source.Start(),
	}
}

// Transform will transform a given byte slice by using the evaluator.
// It will continue even if it encounters an error gathering all
// errors and returning at the end of input.
func (s *Scanner) Transform(input []byte) ([]byte, error) {
	errs := &source.Errors{}
	for _, char := range string(input) {
		if err := s.consumeWithError(char); err != nil {
			errs.Push(err)
		}
	}

	if s.state != ssDefault {
		errs.Push(&source.Error{
			Location: *s.location,
			Message:  "unexpected end of input",
		})
	}

	if err := errs.ErrorOrNil(); err != nil {
		return nil, err
	}
	return s.output.Bytes(), nil
}

func (s *Scanner) consumeWithError(char rune) error {
	defer func() { s.location.Advance(char) }()

	return s.consume(char)
}

func (s *Scanner) consume(char rune) error {
	switch s.state {
	case ssDefault:
		return s.onDefault(char)
	case ssDollar:
		return s.onDollar(char)
	case ssOpenChar:
		return s.onWaitOpen(char)
	case ssExpression:
		return s.onExpression(char)
	case ssCloseChar:
		return s.onWaitClose(char)
	}

	return fmt.Errorf("impossible to handle state %q", s.state)
}

func (s *Scanner) onDefault(char rune) (err error) {
	if char == dollarChar {
		s.state = ssDollar
		s.currentExpressionStart = *s.location
		return
	}

	_, err = s.output.WriteRune(char)
	return
}

func (s *Scanner) onDollar(char rune) (err error) {
	if char == openChar {
		s.state = ssOpenChar
		return
	}

	s.state = ssDefault
	_, err = s.output.Write([]byte{dollarChar, byte(char)})

	return
}

func (s *Scanner) onWaitOpen(char rune) (err error) {
	if char == openChar {
		s.state = ssExpression
		return
	}

	s.state = ssDefault
	_, err = s.output.Write([]byte{dollarChar, openChar, byte(char)})

	return
}

func (s *Scanner) onExpression(char rune) (err error) {
	if char == closeChar {
		s.state = ssCloseChar
		return
	}

	_, err = s.currentExpression.WriteRune(char)

	return
}

func (s *Scanner) onWaitClose(char rune) (err error) {
	if char != closeChar {
		return &source.Error{
			Location: *s.location,
			Message:  fmt.Sprintf("unexpected character %q, expected %q", char, closeChar),
		}
	}

	var out string
	if out, err = s.evaluator.Evaluate(s.currentExpression.String()); err != nil {
		return err
	}

	s.state = ssDefault
	s.currentExpression.Reset()

	_, err = s.output.WriteString(out)

	return
}
