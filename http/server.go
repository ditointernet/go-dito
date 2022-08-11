package http

import (
	"errors"
	"fmt"
	"net/http"

	routing "github.com/jackwhelpton/fasthttp-routing/v2"
	"github.com/jackwhelpton/fasthttp-routing/v2/access"
	"github.com/jackwhelpton/fasthttp-routing/v2/content"
	"github.com/jackwhelpton/fasthttp-routing/v2/fault"
	"github.com/jackwhelpton/fasthttp-routing/v2/slash"
	"github.com/rs/cors"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

// ServerInput encapsulates the necessary Inputs to initialize a Server
type ServerInput struct {
	Port           int
	ReadBufferSize int // Bytes, e.g. 4096 == 4kB
	AllowedOrigins []string
	AllowedHeaders []string
	ExposedHeaders []string
	Handler        func(*routing.Router)
	Logger         logger
}

// Server is an HTTPServer object that can serve HTTP requests
// It is based on httpfast implementation
type Server struct {
	port           int
	readBufferSize int
	allowedOrigins []string
	allowedHeaders []string
	exposedHeaders []string
	router         *routing.Router
	logger         logger
}

var defaultAllowedHeaders = []string{
	"Accept",
	"Authorization",
	"Brand",
	"Content-Type",
	"X-CSRF-Token",
}

var defaultExposedHeaders = []string{
	"X-TOTAL-COUNT",
}

// NewServer creates a new instance of Server
func NewServer(in ServerInput) Server {
	router := routing.New()

	if len(in.AllowedHeaders) == 0 {
		in.AllowedHeaders = defaultAllowedHeaders
	}

	if len(in.ExposedHeaders) == 0 {
		in.ExposedHeaders = defaultExposedHeaders
	}

	server := Server{
		port:           in.Port,
		readBufferSize: in.ReadBufferSize,
		allowedOrigins: in.AllowedOrigins,
		allowedHeaders: in.AllowedHeaders,
		exposedHeaders: in.ExposedHeaders,
		router:         router,
		logger:         in.Logger,
	}

	router.Use(
		slash.Remover(http.StatusMovedPermanently),
		content.TypeNegotiator(content.JSON),
		fault.ErrorHandler(nil, customErrorHandler),
	)

	server.addCorsMiddleware()
	server.addRequestIPIntoContext()
	if in.Logger != nil {
		server.addRequestLogger()
	}

	in.Handler(router)

	return server
}

// Run listen and serves HTTP requests at the Port specified in Server construction
func (s Server) Run() error {
	server := fasthttp.Server{
		Handler:        s.router.HandleRequest,
		ReadBufferSize: s.readBufferSize,
	}

	return server.ListenAndServe(fmt.Sprintf(":%d", s.port))
}

func (s Server) addCorsMiddleware() {
	if len(s.allowedOrigins) == 0 {
		return
	}

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: s.allowedOrigins,
		ExposedHeaders: s.exposedHeaders,
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodDelete,
			http.MethodPut,
			http.MethodPatch,
		},
		AllowedHeaders:   s.allowedHeaders,
		AllowCredentials: true,
		MaxAge:           300,
	})

	s.router.Use(routing.RequestHandlerFunc(fasthttpadaptor.NewFastHTTPHandlerFunc(corsMiddleware.HandlerFunc)))
}

func (s Server) addRequestIPIntoContext() {
	s.router.Use(func(ctx *routing.Context) error {
		ctx.SetUserValue(ContextKeyRequestIPAddress, access.GetClientIP(ctx.RequestCtx))
		return nil
	})
}

func (s Server) addRequestLogger() {
	s.router.Use(access.CustomLogger(func(ctx *fasthttp.RequestCtx, elapsed float64) {
		ip := ctx.UserValue(ContextKeyRequestIPAddress)
		req := fmt.Sprintf("%s %s %s", string(ctx.Request.Header.Method()), string(ctx.RequestURI()), string(ctx.Request.URI().Scheme()))
		s.logger.Info(ctx, `[%v] [%.3fms] %s %d %d`, ip, elapsed, req, ctx.Response.StatusCode(), len(ctx.Response.Body()))
	}))
}

func customErrorHandler(ctx *routing.Context, err error) error {
	var errorResponse ErrorResponse
	var errorListResponse ErrorListResponse
	if errors.As(err, &errorResponse) || errors.As(err, &errorListResponse) {
		return err
	}

	return NewErrorResponse(ctx, err)
}
