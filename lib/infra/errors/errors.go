package errors

import (
	e "errors"
	"fmt"
)

// CodeType is a string that contains error's code description
type CodeType string

// KindType is a string that contains error's kind description
type KindType string

// CustomError is a structure that encodes useful information about a given error.
// It's supposed to flow within the application in detriment of the the default golang error,
// since its Kind and Code attributes are the keys to express its semantic and uniqueness, respectively.
// It should be generated once by the peace of code that found the error (because it's where we have more context about the error),
// and be by passed to the upper layers of the application.
type CustomError struct {
	kind    KindType
	code    CodeType
	message string
}

const (
	// DefaultKind is the default kind returned by the error library, whenever the method Kind(err error) KindType is called
	// and the error argument isnt a CustomError
	DefaultKind KindType = "DEFAULT_ERROR_KIND"
	//DefaultCode is the default code returned by the error library, whenever the method Code(err error) CodeType is called
	// and the error argument isnt a CustomError
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

// New returns a new instance of CustomError with the given message
func New(message string, args ...interface{}) CustomError {
	return CustomError{
		kind:    DefaultKind,
		code:    DefaultCode,
		message: fmt.Sprintf(message, args...),
	}
}

// WithKind return a copy of the CustomError with the given KindType filled
func (ce CustomError) WithKind(kind KindType) CustomError {
	ce.kind = kind
	return ce
}

// WithCode return a copy of the CustomError with the given CodeType filled
func (ce CustomError) WithCode(code CodeType) CustomError {
	ce.code = code
	return ce
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
