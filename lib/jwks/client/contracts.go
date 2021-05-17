package client

import "context"

// HTTPHeaders is a map'containing the relation key=value of the heades used on the http rest request.
type HTTPHeaders map[string]string

// HTTPQueryParams is a map'containing the relation key=value of the query params used on the http rest request
type HTTPQueryParams map[string]string

// HttpRequest are the params used to build a new http rest request

type HttpRequest struct {
	URL         string
	Body        []byte
	Headers     HTTPHeaders
	QueryParams HTTPQueryParams
}

type HttpResult struct {
	StatusCode int
	Response   []byte
}

// HttpClient is the public interface of the http client lib
type HttpClient interface {
	Patch(ctx context.Context, request HttpRequest) (rst HttpResult, err error)
	Put(ctx context.Context, request HttpRequest) (rst HttpResult, err error)
	Post(ctx context.Context, request HttpRequest) (rst HttpResult, err error)
	Delete(ctx context.Context, request HttpRequest) (rst HttpResult, err error)
	Get(ctx context.Context, request HttpRequest) (rst HttpResult, err error)
	PostForm(ctx context.Context, request HttpRequest) (rst HttpResult, err error)
}
