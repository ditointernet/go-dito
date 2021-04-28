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
	DefaultKind KindType = "DEFAULT_ERROR_KIND"
	DefaultCode CodeType = "DEFAULT_ERROR_CODE"
	// KindInternal are errors caused by some internal fail like failed IO calls or invalid memory states
	KindInternal KindType = "INTERNAL"
	// KindInvalidInput are errors caused by some invalid values on the input
	KindInvalidInput KindType = "INVALID_INPUT"
	// KindNotFound are errors caused by any required resources that not exists on the data repository
	KindNotFound KindType = "NOT_FOUND"
	// KindAuthentication are errors caused by an unauthenticated call
	KindAuthentication KindType = "AUTHENTICATION"
	// KindAuthorization are errors caused by an unauthorized call
	KindAuthorization KindType = "AUTHORIZATION"
)

// GetService returns a new instance of CustomError with message, kind and code
func New(message string, kind KindType, code CodeType) CustomError {
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
