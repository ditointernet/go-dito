package errors

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"example/infra"
)

// New ...
func New(args ...interface{}) *infra.Error {
	err := infra.Error{
		Metadata: infra.Metadata{},
	}

	for _, arg := range args {
		switch arg := arg.(type) {
		case context.Context:
			err.Ctx = arg
		case error:
			err.Err = arg
		case infra.OpName:
			err.OpName = arg
		case infra.ErrorKind:
			err.Kind = arg
		case infra.Severity:
			err.Severity = arg
		case infra.Metadata:
			err.Metadata = err.Metadata.Merge(&arg)
		case string:
			err.Err = errors.New(arg)
		}
	}

	return &err
}

// Context ...
func Context(err error) context.Context {
	e, ok := err.(*infra.Error)
	if !ok {
		return context.Background()
	}

	if e.Ctx != nil {
		return e.Ctx
	}

	return Context(e.Err)
}

// Trace ...
func Trace(err *infra.Error) []infra.OpName {
	trace := []infra.OpName{err.OpName}

	nextError, ok := err.Err.(*infra.Error)
	if !ok {
		return trace
	}

	trace = append(trace, Trace(nextError)...)

	return trace
}

// Kind ...
func Kind(err error) infra.ErrorKind {
	e, ok := err.(*infra.Error)
	if !ok {
		return infra.KindUnexpected
	}

	if e.Kind != 0 {
		return e.Kind
	}

	return Kind(e.Err)
}

// Severity ...
func Severity(err error) infra.Severity {
	e, ok := err.(*infra.Error)
	if !ok {
		return infra.SeverityError
	}

	if e.Severity != "" {
		return e.Severity
	}

	return Severity(e.Err)
}

// Metadata ...
func Metadata(err *infra.Error) infra.Metadata {
	nextErr, ok := err.Err.(*infra.Error)
	if !ok {
		return err.Metadata
	}

	return Metadata(nextErr).Merge(&err.Metadata)
}

// OpName ...
func OpName(err *infra.Error) infra.OpName {
	nextError, ok := err.Err.(*infra.Error)
	if !ok {
		return err.OpName
	}

	return OpName(nextError)
}

// Error ...
func Error(err *infra.Error) error {
	nextError, ok := err.Err.(*infra.Error)
	if !ok {
		return err.Err
	}

	return Error(nextError)
}

// Log ...
func Log(log infra.LogProvider, err *infra.Error) {
	method := fmt.Sprintf("%sMetadata", Severity(err))

	values := []reflect.Value{
		reflect.ValueOf(Context(err)),
		reflect.ValueOf(OpName(err)),
		reflect.ValueOf(err.Err.Error()),
		reflect.ValueOf(Metadata(err).Merge(&infra.Metadata{
			"trace": Trace(err),
			"kind":  Kind(err),
		})),
	}

	reflect.ValueOf(log).MethodByName(method).Call(values)
}
