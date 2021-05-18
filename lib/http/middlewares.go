package http

import (
	"context"
	"fmt"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/ditointernet/go-dito/lib/errors"
	routing "github.com/jackwhelpton/fasthttp-routing/v2"
)

type JWKSCLient interface {
	GetCerts(ctx context.Context) error
	RenewCerts(ctx context.Context) error
	Certs() map[string]string
}

// UserAuthenticator ...
type UserAuthenticator struct {
	logger logger
	jwks   JWKSCLient
}

// NewUserAuthenticator ...
func NewUserAuthenticator(logger logger, jwks JWKSCLient) UserAuthenticator {
	return UserAuthenticator{
		logger: logger,
		jwks:   jwks,
	}
}

// Authenticate ...
//
// It tries to authenticate the token with the certifications on memory,
// if it fails, the certifications are renewed and the authentication is
// run again.
func (ua UserAuthenticator) Authenticate(ctx *routing.Context) error {
	authHeader := string(ctx.Request.Header.Peek("Authorization"))
	if len(authHeader) < 7 || strings.ToLower(authHeader[:7]) != "bearer " || authHeader[7:] == "" {
		// error
		return nil
	}
	token := authHeader[7:]

	certs := ua.jwks.Certs()

	if parsedToken, err := jwt.Parse(token, verifyJWTSignature(certs)); err == nil {
		setUserID(ctx, parsedToken)
		return nil
	}

	if err := ua.jwks.RenewCerts(ctx); err != nil {
		// error on renewing token
		return nil
	}
	certs = ua.jwks.Certs()

	parsedToken, err := jwt.Parse(token, verifyJWTSignature(certs))
	if err == nil {
		setUserID(ctx, parsedToken)
		return nil
	}

	// error invalid token
	return nil
}

func setUserID(ctx *routing.Context, token *jwt.Token) {
	claims, _ := token.Claims.(jwt.MapClaims)
	sub, _ := claims["sub"].(string)
	// removes auth provider prefix 'auth0|' to get only the user identifier.
	userID := strings.Split(sub, "|")[1]
	ctx.SetUserValue("userID", userID)
}

func verifyJWTSignature(certs map[string]string) func(token *jwt.Token) (interface{}, error) {
	return func(token *jwt.Token) (interface{}, error) {
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("token's kid header not found")
		}

		cert, ok := certs[kid]
		if !ok {
			return nil, errors.New("cert key not found")
		}

		result, err := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
		if err != nil {
			return nil, errors.New("error trying to validate signature")
		}

		return result, nil
	}
}

func handleAuthenticationError(ctx *routing.Context, logger, err error) {
	rID, _ := ctx.UserValue("request-id").(string)
	var errResponse string
	fmt.Println(rID)
	// if e, ok := err.(domain.APIError); ok {
	// 	ctx.SetStatusCode(http.StatusUnauthorized)
	// 	errResponse = fmt.Sprintf(`{"errors":[{"error":"%s","code":"%s"}]}`, err.Error(), e.Code())
	// 	// logger.Errorf(infra.LogOptions{TraceID: rID},
	// 	// 	"token-authentication-error", "error: %s | code: %s", err.Error(), e.Code())
	// } else {
	// 	ctx.SetStatusCode(http.StatusInternalServerError)
	// 	errResponse = `{"errors":[{"error":"internal"}]}`
	// 	logger.Errorf(infra.LogOptions{TraceID: rID},
	// 		"token-authentication-error", "internal: %s", err.Error())
	// }

	ctx.SetContentTypeBytes([]byte("application/json"))
	ctx.SetBodyString(errResponse)
	ctx.Abort()
}
