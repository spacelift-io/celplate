package celplate

import (
	"bytes"
	"errors"
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

func NewScanner(evaluator Evaluator) *Scanner {
	return &Scanner{
		currentExpression: bytes.NewBuffer(nil),
		output:            bytes.NewBuffer(nil),
		state:             ssDefault,
		evaluator:         evaluator,
		location:          source.Start(),
	}
}

func (s *Scanner) Transform(input []byte) ([]byte, error) {
	for _, char := range string(input) {
		if err := s.consumeWithError(char); err != nil {
			return nil, err
		}
	}

	if s.state != ssDefault {
		return nil, &source.Error{
			Location: *s.location,
			Message:  "unexpected end of input",
		}
	}

	return s.output.Bytes(), nil
}

func (s *Scanner) consumeWithError(char rune) (err error) {
	defer func() { s.location.Advance(char) }()

	if err = s.consume(char); err == nil {
		return nil
	}

	var sourceErrors source.Errors
	if errors.As(err, &sourceErrors) {
		// The location of each expression error is relative to the expression
		// start location, so we need to do some math to get the absolute
		// location.
		for i := range sourceErrors {
			sourceErrors[i].Location = s.currentExpressionStart.Nested(sourceErrors[i].Location)
		}

		return sourceErrors
	}

	return &source.Error{
		Location: *s.location,
		Message:  err.Error(),
	}
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
		return fmt.Errorf("unexpected character %q, expected %q", char, closeChar)
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
