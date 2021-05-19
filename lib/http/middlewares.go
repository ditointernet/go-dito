package http

import (
	"context"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/ditointernet/go-dito/lib/errors"
	routing "github.com/jackwhelpton/fasthttp-routing/v2"
)

// ContextKeyAccountID is the key used to retrieve and save accountId into the context
const ContextKeyAccountID string = "account-id"

type jwksClient interface {
	GetCerts(ctx context.Context) error
	RenewCerts(ctx context.Context) error
	Certs() map[string]string
}

// UserAuthenticator structure responsible for handling request authentication
type UserAuthenticator struct {
	logger logger
	jwks   jwksClient
}

// NewUserAuthenticator creates a new instance of the UserAuthenticator structure
func NewUserAuthenticator(logger logger, jwks jwksClient) UserAuthenticator {
	return UserAuthenticator{
		logger: logger,
		jwks:   jwks,
	}
}

// Authenticate is responsible for verify if the request is authenticated
//
// It tries to authenticate the token with the certifications on memory,
// if it fails, the certifications are renewed and the authentication is
// run again.
func (ua UserAuthenticator) Authenticate(ctx *routing.Context) error {
	authHeader := string(ctx.Request.Header.Peek("Authorization"))
	if len(authHeader) < 7 || strings.ToLower(authHeader[:7]) != "bearer " || authHeader[7:] == "" {
		err := errors.New("unauthenticated request").WithKind(errors.KindUnauthenticated)
		return NewErrorResponse(ctx, err)
	}
	token := authHeader[7:]

	certs := ua.jwks.Certs()

	if parsedToken, err := jwt.Parse(token, verifyJWTSignature(certs)); err == nil {
		setAccountID(ctx, parsedToken)
		return nil
	}

	if err := ua.jwks.RenewCerts(ctx); err != nil {
		err = errors.New("error on renewing the certificates").WithKind(errors.KindInternal)
		ua.logger.Error(ctx, err)
		return NewErrorResponse(ctx, err)
	}
	certs = ua.jwks.Certs()

	parsedToken, err := jwt.Parse(token, verifyJWTSignature(certs))
	if err != nil {
		return NewErrorResponse(ctx, err)
	}

	setAccountID(ctx, parsedToken)
	return nil
}

func setAccountID(ctx *routing.Context, token *jwt.Token) {
	claims, _ := token.Claims.(jwt.MapClaims)
	sub, _ := claims["sub"].(string)
	// removes auth provider prefix 'auth0|' to get only the user identifier.
	accountID := strings.Split(sub, "|")[1]
	ctx.SetUserValue(ContextKeyAccountID, accountID)
}

func verifyJWTSignature(certs map[string]string) func(token *jwt.Token) (interface{}, error) {
	return func(token *jwt.Token) (interface{}, error) {
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("token's kid header not found").WithKind(errors.KindUnauthenticated)
		}

		cert, ok := certs[kid]
		if !ok {
			return nil, errors.New("cert key not found").WithKind(errors.KindUnauthenticated)
		}

		result, err := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
		if err != nil {
			return nil, errors.New("error trying to validate signature").WithKind(errors.KindInternal)
		}

		return result, nil
	}
}