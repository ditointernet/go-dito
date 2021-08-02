package infra

import (
	"context"
)

type Logger interface {
	Debug(ctx context.Context, msg string, args ...interface{})
	Info(ctx context.Context, msg string, args ...interface{})
	Warning(ctx context.Context, msg string, args ...interface{})
	Error(ctx context.Context, err error)
	Critical(ctx context.Context, err error)
}

type AuthorizatorClient interface {
	DecideIfAllowed(ctx context.Context, regoQuery string, method, path, brandID, userID string) (bool, error)
	ExecuteQuery(ctx context.Context, query string, input map[string]interface{}) ([]map[string]interface{}, error)
}
