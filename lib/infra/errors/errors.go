package errors

import (
	e "errors"
)

type CodeType string
type KindType string

type CustomError struct {
	kind    KindType
	code    CodeType
	message string
}

const (
	DefaultKind KindType = ""
	DefaultCode CodeType = ""
)

// GetService returns a new instance of CustomError with kind, code and message
func New(kind KindType, message string, code CodeType) CustomError {
	return CustomError{
		kind:    kind,
		code:    code,
		message: message,
	}
}

// Error returns CustomError message
func (ce CustomError) Error() string {
	return ce.message
}

// Kind this method receives an error, then compares its interface type with the CustomError interface
// if the interfaces types matches, returns its kind
func Kind(err error) KindType {
	var customError CustomError
	if e.As(err, &customError) {
		return customError.kind
	}
	return DefaultKind
}

// Kind this method receives an error, then compares its interface type with the CustomError interface
// if the interfaces types matches, returns its Code
func Code(err error) CodeType {
	var customError CustomError
	if e.As(err, &customError) {
		return customError.code
	}
	return DefaultCode
}
