package main

import (
	routing "github.com/jackwhelpton/fasthttp-routing/v2"

	"github.com/ditointernet/go-dito/errors"
	"github.com/ditointernet/go-dito/http"
	"github.com/ditointernet/go-dito/log"
)

func main() {
	logger := log.NewLogger(log.LoggerInput{Level: "DEBUG"})

	server := http.NewServer(http.ServerInput{
		Port:   3000,
		Logger: logger,
		Handler: func(r *routing.Router) {
			r.Get("/route1", func(ctx *routing.Context) error {
				return http.NewErrorResponse(ctx, errors.New("random error").
					WithCode("RANDOM_ERROR").
					WithKind(errors.KindInvalidInput))
			})

			r.Get("/route2", func(ctx *routing.Context) error {
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

			r.Get("/route3", func(ctx *routing.Context) error {
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

			r.Get("/route4", func(ctx *routing.Context) error {
				ctx.Write(http.NewMessageResponse("something was done"))
				ctx.SetStatusCode(200)

				return nil
			})

			r.Get("/route5", func(ctx *routing.Context) error {
				ctx.Write(http.NewResourceCreatedResponse("random-id"))
				ctx.SetStatusCode(201)

				return nil
			})

			r.Get("/route6", func(ctx *routing.Context) error {
				ctx.Write(http.NewResourceCreatedResponse("random-id").WithMessage("Good job, %s!", "buddy"))
				ctx.SetStatusCode(201)

				return nil
			})
		},
	})

	server.Run()
}
