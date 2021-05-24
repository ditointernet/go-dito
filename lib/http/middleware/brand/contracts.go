package brand

import "context"

type logger interface {
	Debug(ctx context.Context, msg string, args ...interface{})
	Info(ctx context.Context, msg string, args ...interface{})
	Warning(ctx context.Context, msg string, args ...interface{})
	Error(ctx context.Context, err error)
	Critical(ctx context.Context, err error)
}
