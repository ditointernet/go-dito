package main

import (
	"context"

	routing "github.com/jackwhelpton/fasthttp-routing/v2"

	"github.com/ditointernet/go-dito/lib/errors"
	"github.com/ditointernet/go-dito/lib/http"
	"github.com/ditointernet/go-dito/lib/log"
)

func main() {
	ctx := context.Background()

	logger := log.NewLogger(log.LoggerInput{Level: "DEBUG"})

	server := http.NewServer(http.ServerInput{
		Port:   3000,
		Logger: logger,
		Handler: func(r *routing.Router) {
			r.Get("/route1", func(c *routing.Context) error {
				return http.NewErrorResponse(ctx, errors.New("random error").
					WithCode("RANDOM_ERROR").
					WithKind(errors.KindInvalidInput))
			})

			r.Get("/route2", func(c *routing.Context) error {
				errs := []error{
					errors.New("random error 1").
						WithCode("RANDOM_ERROR_1").
						WithKind(errors.KindUnauthenticated),
					errors.New("random error 2").
						WithCode("RANDOM_ERROR_2").
						WithKind(errors.KindUnauthorized),
				}

				return http.NewErrorListResponse(ctx, errs...)
			})

			r.Get("/route3", func(c *routing.Context) error {
				errs := []error{
					errors.New("random error 1").
						WithCode("RANDOM_ERROR_1").
						WithKind(errors.KindUnauthenticated),
					errors.New("random error 2").
						WithCode("RANDOM_ERROR_2").
						WithKind(errors.KindUnauthorized),
				}

				return http.NewErrorListResponse(ctx, errs...).WithStatusCode(409)
			})

			r.Get("/route4", func(c *routing.Context) error {
				c.WriteString(http.NewResponseMessage("something was created"))
				c.SetStatusCode(201)

				return nil
			})
		},
	})

	server.Run()
}
