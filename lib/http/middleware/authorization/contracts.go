package authorization

import "context"

// AuthQueryResult ...
type AuthQueryResult []map[string]interface{}

type AuthorizatorClient interface {
	DecideIfAllowed(ctx context.Context, regoQuery string, method, path, brandID, userID string) (bool, error)
	ExecuteQuery(ctx context.Context, query string, input map[string]interface{}) (AuthQueryResult, error)
}

type Logger interface {
	Debug(ctx context.Context, msg string, args ...interface{})
	Info(ctx context.Context, msg string, args ...interface{})
	Warning(ctx context.Context, msg string, args ...interface{})
	Error(ctx context.Context, err error)
	Critical(ctx context.Context, err error)
}
