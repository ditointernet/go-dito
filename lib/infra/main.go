// package main

// import (
// 	"context"

// 	"github.com/ditointernet/go-dito/lib/infra/log"
// )

// func main() {
// 	ctx := context.Background()

// 	logger := log.NewLogger(log.LoggerInput{Level: "INFO", Attributes: log.LogAttributeSet{
// 		"id1": true,
// 		"id2": true,
// 	}})

// 	ctx = context.WithValue(ctx, "id1", "value1")
// 	ctx = context.WithValue(ctx, "id2", "value2")
// 	// fmt.Println(ctx.Value("id1"))

// 	logger.Info(ctx, "Hello World")
// }
