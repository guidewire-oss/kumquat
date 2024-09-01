package renderer

import "fmt"

type LookupError struct {
	rendererName string
}

func (e *LookupError) Error() string {
	return fmt.Sprintf("unknown renderer '%s'", e.rendererName)
}

type Error struct {
	err    error
	line   int
	column int
}

func NewError(err error, line, column int) *Error {
	if line == 0 {
		column = 0
	}

	return &Error{err, line, column}
}

// Error returns the error message.
func (e *Error) Error() string {
	if e.line == 0 {
		return fmt.Sprintf("[line ?, column ?] %s", e.err.Error())
	}

	if e.column == 0 {
		return fmt.Sprintf("[line %d, column ?] %s", e.line, e.err.Error())
	}

	return fmt.Sprintf("[line %d, column %d] %s", e.line, e.column, e.err.Error())
}

// Unwrap returns the wrapped error.
func (e *Error) Unwrap() error {
	return e.err
}

// Line returns the line number where the error occurred; 0 if unknown.
func (e *Error) Line() int {
	return e.line
}

// Column returns the column number where the error occurred; 0 if unknown.
func (e *Error) Column() int {
	return e.column
}
