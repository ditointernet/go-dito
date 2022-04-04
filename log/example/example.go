package main

import (
	"context"
	"errors"

	"github.com/ditointernet/go-dito/lib/log"
)

func main() {
	ctx := context.WithValue(context.Background(), "RandomContextKey", "RandomContextValue")

	logger := log.NewLogger(log.LoggerInput{
		Level: "DEBUG",
		Attributes: log.LogAttributeSet{
			log.LogAttribute("RandomContextKey"): true,
		},
	})

	logger.Debug(ctx, "hello world")
	logger.Info(ctx, "hello world")
	logger.Warning(ctx, "hello world")
	logger.Error(ctx, errors.New("random error"))
	logger.Critical(ctx, errors.New("random error"))
}
