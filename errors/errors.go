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
	message error
}

const (
	// CodeUnknown is the default code returned when the application doesn't attach any code into the error
	CodeUnknown CodeType = "UNKNOWN"
	// KindUnexpected is the default kind returned  when the application doesn't attach any kind into the error
	KindUnexpected KindType = "UNEXPECTED"
	// KindConflict are errors caused by requests with data that conflicts with the current state of the system
	KindConflict KindType = "CONFLICT"
	// KindInternal are errors caused by some internal fail like failed IO calls or invalid memory states
	KindInternal KindType = "INTERNAL"
	// KindInvalidInput are errors caused by some invalid values on the input
	KindInvalidInput KindType = "INVALID_INPUT"
	// KindNotFound are errors caused by any required resources that not exists on the data repository
	KindNotFound KindType = "NOT_FOUND"
	// KindUnauthentication are errors caused by an unauthenticated call
	KindUnauthenticated KindType = "UNAUTHENTICATED"
	// KindUnauthorized are errors caused by an unauthorized call
	KindUnauthorized KindType = "UNAUTHORIZED"
)

// New returns a new instance of CustomError with the given message
func New(message string, args ...interface{}) CustomError {
	return CustomError{
		kind:    KindUnexpected,
		code:    CodeUnknown,
		message: fmt.Errorf(message, args...),
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
	return ce.message.Error()
}

func (ce CustomError) Unwrap() error {
	return ce.message
}

// NewMissingRequiredDependency creates a new error that indicates a missing required dependency.
// It should be producing at struct constructors.
func NewMissingRequiredDependency(name string) error {
	return New("missing required dependency: %s", name).WithCode("MISSING_REQUIRED_DEPENDENCY")
}

// Kind this method receives an error, then compares its interface type with the CustomError interface
// if the interfaces types matches, returns its kind
func Kind(err error) KindType {
	var customError CustomError
	if e.As(err, &customError) {
		return customError.kind
	}
	return KindUnexpected
}

// Kind this method receives an error, then compares its interface type with the CustomError interface
// if the interfaces types matches, returns its Code
func Code(err error) CodeType {
	var customError CustomError
	if e.As(err, &customError) {
		return customError.code
	}
	return CodeUnknown
}
