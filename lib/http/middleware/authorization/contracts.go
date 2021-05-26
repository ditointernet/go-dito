package authorization

import (
	"context"

	"github.com/ditointernet/go-dito/lib/opa"
)

type authorizatorClient interface {
	DecideIfAllowed(ctx context.Context, regoQuery string, method, path, brandID, userID string) (bool, error)
	ExecuteQuery(ctx context.Context, query string, input map[string]interface{}) (opa.AuthorizationResult, error)
}

type logger interface {
	Debug(ctx context.Context, msg string, args ...interface{})
	Info(ctx context.Context, msg string, args ...interface{})
	Warning(ctx context.Context, msg string, args ...interface{})
	Error(ctx context.Context, err error)
	Critical(ctx context.Context, err error)
}
