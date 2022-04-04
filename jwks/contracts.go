package jwks

import (
	"context"

	"github.com/ditointernet/go-dito/lib/http"
)

// HTTPClient is the public interface of the http client lib
type HTTPClient interface {
	Patch(ctx context.Context, request http.HTTPRequest) (rst http.HTTPResult, err error)
	Put(ctx context.Context, request http.HTTPRequest) (rst http.HTTPResult, err error)
	Post(ctx context.Context, request http.HTTPRequest) (rst http.HTTPResult, err error)
	Delete(ctx context.Context, request http.HTTPRequest) (rst http.HTTPResult, err error)
	Get(ctx context.Context, request http.HTTPRequest) (rst http.HTTPResult, err error)
	PostForm(ctx context.Context, request http.HTTPRequest) (rst http.HTTPResult, err error)
}
