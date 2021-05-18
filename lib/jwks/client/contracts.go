package client

import (
	"context"

	"github.com/ditointernet/go-dito/lib/http"
)

// HttpClient is the public interface of the http client lib
type HttpClient interface {
	Patch(ctx context.Context, request http.HttpRequest) (rst http.HttpResult, err error)
	Put(ctx context.Context, request http.HttpRequest) (rst http.HttpResult, err error)
	Post(ctx context.Context, request http.HttpRequest) (rst http.HttpResult, err error)
	Delete(ctx context.Context, request http.HttpRequest) (rst http.HttpResult, err error)
	Get(ctx context.Context, request http.HttpRequest) (rst http.HttpResult, err error)
	PostForm(ctx context.Context, request http.HttpRequest) (rst http.HttpResult, err error)
}
