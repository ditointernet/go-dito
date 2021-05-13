package http

import (
	"context"
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

type logger interface {
	Debug(ctx context.Context, msg string, args ...interface{})
	Info(ctx context.Context, msg string, args ...interface{})
	Warning(ctx context.Context, msg string, args ...interface{})
	Error(ctx context.Context, err error)
	Critical(ctx context.Context, err error)
}

// ContextKeyRequestIPAddress is the key of RequestIP information injected into the request context
const ContextKeyRequestIPAddress = "request_ip"

// ServerInput encapsulates the necessary Inputs to initialize a Server
type ServerInput struct {
	Port           int
	AllowedOrigins []string
	Handler        func(*routing.Router)
	Logger         logger
}

// Server is an HTTPServer object that can serve HTTP requests
// It is based on httpfast implementation
type Server struct {
	port           int
	allowedOrigins []string
	router         *routing.Router
	logger         logger
}

// NewServer creates a new instance of Server
func NewServer(in ServerInput) Server {
	router := routing.New()

	server := Server{
		port:           in.Port,
		allowedOrigins: in.AllowedOrigins,
		router:         router,
		logger:         in.Logger,
	}

	router.Use(slash.Remover(http.StatusMovedPermanently))
	router.Use(content.TypeNegotiator(content.JSON))
	router.Use(fault.ErrorHandler(nil, customErrorContentType))

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
	return fasthttp.ListenAndServe(fmt.Sprintf(":%d", s.port), s.router.HandleRequest)
}

func (s Server) addCorsMiddleware() {
	if len(s.allowedOrigins) == 0 {
		return
	}

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: s.allowedOrigins,
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodDelete,
			http.MethodPut,
			http.MethodPatch,
		},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
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

func customErrorContentType(ctx *routing.Context, err error) error {
	_, ok := err.(ErrorResponse)
	if ok {
		ctx.Response.Header.SetContentType(content.JSON)
	}

	_, ok = err.(ErrorListResponse)
	if ok {
		ctx.Response.Header.SetContentType(content.JSON)
	}

	return err
}
