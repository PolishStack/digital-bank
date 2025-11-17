package errs

import (
	"errors"
	"fmt"
)

// Code is a stable machine-readable code for each error kind.
type Code string

const (
	CodeInternal   Code = "internal_error"
	CodeNotFound   Code = "not_found"
	CodeConflict   Code = "conflict"
	CodeBadRequest Code = "bad_request"
	CodeUnauth     Code = "unauthenticated"
)

// CodedError is an error that carries a machine code and a safe message.
type CodedError struct {
	CodeVal Code   // machine code
	Msg     string // developer message - safe to log
	Err     error  // wrapped underlying error (may be nil)
}

func (e *CodedError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Msg, e.Err)
	}
	return e.Msg
}

func (e *CodedError) Unwrap() error { return e.Err }

// Code returns the machine code.
func (e *CodedError) Code() Code { return e.CodeVal }

// Helpers
func New(code Code, msg string) *CodedError {
	return &CodedError{CodeVal: code, Msg: msg}
}

func Wrap(err error, code Code, msg string) *CodedError {
	if err == nil {
		return New(code, msg)
	}
	// If err is already CodedError, keep its code if it's the same; otherwise wrap
	var ce *CodedError
	if errors.As(err, &ce) {
		// If caller provides same code, just return original wrapped with extra context
		if ce.CodeVal == code {
			return &CodedError{CodeVal: code, Msg: msg, Err: err}
		}
	}
	return &CodedError{CodeVal: code, Msg: msg, Err: err}
}

// Predicates
func IsCode(err error, code Code) bool {
	var ce *CodedError
	if errors.As(err, &ce) {
		return ce.CodeVal == code
	}
	return false
}

func AsCoded(err error) (*CodedError, bool) {
	var ce *CodedError
	if errors.As(err, &ce) {
		return ce, true
	}
	return nil, false
}
