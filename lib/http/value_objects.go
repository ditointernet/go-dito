package http

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ditointernet/go-dito/lib/errors"
)

type errorPayload struct {
	Code    errors.CodeType `json:"code,omitempty"`
	Message string          `json:"message"`
}

// ErrorResponse defines how errors should be presented to clients of HTTPServer.
// It implements fasthttp's HTTPError contract (https://github.com/jackwhelpton/fasthttp-routing/blob/master/error.go#L12),
// so it keeps its error handling behaviour.
type ErrorResponse struct {
	status int

	TraceID string       `json:"trace_id,omitempty"`
	Err     errorPayload `json:"error"`
}

// NewErrorResponse creates a new ErrorResponse object.
func NewErrorResponse(ctx context.Context, err error) ErrorResponse {
	return ErrorResponse{
		// TraceID will be collected in the future, when tracing is properly implemented by GoDito
		status: kindToHTTPStatusCode(errors.Kind(err)),
		Err: errorPayload{
			Code:    errors.Code(err),
			Message: err.Error(),
		},
	}
}

// Error returns the error message.
func (e ErrorResponse) Error() string {
	msg, _ := json.Marshal(e)
	return string(msg)
}

// StatusCode returns the HTTP status code.
func (e ErrorResponse) StatusCode() int {
	return e.status
}

// ErrorListResponse defines how list of errors should be presented to clients of HTTPServer.
// It implements fasthttp's HTTPError contract (https://github.com/jackwhelpton/fasthttp-routing/blob/master/error.go#L12),
// so it keeps its error handling behaviour.
type ErrorListResponse struct {
	status int

	TraceID string         `json:"trace_id,omitempty"`
	Errs    []errorPayload `json:"errors"`
}

// NewErrorListResponse creates a new ErrorListResponse object.
func NewErrorListResponse(ctx context.Context, errs ...error) ErrorListResponse {
	if len(errs) == 0 {
		return ErrorListResponse{
			status: 500,
		}
	}

	errsPayload := []errorPayload{}
	for _, err := range errs {
		errsPayload = append(errsPayload, errorPayload{
			Code:    errors.Code(err),
			Message: err.Error(),
		})
	}

	return ErrorListResponse{
		// TraceID will be collected in the future, when tracing is properly implemented by GoDito
		status: kindToHTTPStatusCode(errors.Kind(errs[0])),
		Errs:   errsPayload,
	}
}

// Error returns the error message.
func (e ErrorListResponse) Error() string {
	msg, _ := json.Marshal(e)
	return string(msg)
}

// StatusCode returns the HTTP status code.
func (e ErrorListResponse) StatusCode() int {
	return e.status
}

// WithStatusCode is a way to explicitly define the desired StatusCode to be sent to the client
// with an ErrorListResponse.
func (e ErrorListResponse) WithStatusCode(status int) ErrorListResponse {
	e.status = status
	return e
}

func kindToHTTPStatusCode(kind errors.KindType) int {
	switch kind {
	case errors.KindInvalidInput:
		return 400
	case errors.KindUnauthenticated:
		return 401
	case errors.KindUnauthorized:
		return 403
	case errors.KindNotFound:
		return 404
	case errors.KindConflict:
		return 409
	case errors.KindUnexpected:
		return 500
	case errors.KindInternal:
		return 500
	default:
		return 500
	}
}

// MessageResponse is a generic message that should be sent to a client of HTTP Server.
type MessageResponse struct {
	Message string `json:"message"`
}

// NewMessageResponse creates a MessageResponse
func NewMessageResponse(msg string, params ...interface{}) MessageResponse {
	return MessageResponse{Message: fmt.Sprintf(msg, params...)}
}

// ResourceCreatedResponse is is the default response sent when a new resource is created in the system.
type ResourceCreatedResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

// NewResourceCreatedResponse creates a new ResourceCreatedResponse.
func NewResourceCreatedResponse(id string) ResourceCreatedResponse {
	return ResourceCreatedResponse{ID: id, Message: "Resource created succesfuly!"}
}

// WithMessage overrides the default message.
func (r ResourceCreatedResponse) WithMessage(msg string, params ...interface{}) ResourceCreatedResponse {
	r.Message = fmt.Sprintf(msg, params...)
	return r
}
