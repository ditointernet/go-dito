package authentication

import (
	"context"
	"fmt"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/ditointernet/go-dito/lib/errors"
	"github.com/ditointernet/go-dito/lib/http"
	routing "github.com/jackwhelpton/fasthttp-routing/v2"
)

// ContextKeyAccountID is the key used to retrieve and save accountId into the context
const ContextKeyAccountID string = "account-id"

// AccountAuthenticator structure responsible for handling request authentication
type AccountAuthenticator struct {
	logger logger
	jwks   jwksClient
}

// NewAccountAuthenticator creates a new instance of the AccountAuthenticator structure
func NewAccountAuthenticator(logger logger, jwks jwksClient) (AccountAuthenticator, error) {
	if logger == nil {
		return AccountAuthenticator{}, errors.New("missing logger dependency").WithKind(errors.KindInternal)
	}
	if jwks == nil {
		return AccountAuthenticator{}, errors.New("missing jkwks client dependency").WithKind(errors.KindInternal)
	}
	return AccountAuthenticator{
		logger: logger,
		jwks:   jwks,
	}, nil
}

// Authenticate is responsible for verify if the request is authenticated
//
// It tries to authenticate the token with the certifications on memory,
// if it fails, the certifications are renewed and the authentication is
// run again.
func (ua AccountAuthenticator) Authenticate(ctx *routing.Context) error {
	authHeader := string(ctx.Request.Header.Peek("Authorization"))
	if len(authHeader) < 7 || strings.ToLower(authHeader[:7]) != "bearer " || authHeader[7:] == "" {
		err := errors.New("unauthenticated request").WithKind(errors.KindUnauthenticated)
		ua.logger.Error(ctx, err)
		return http.NewErrorResponse(ctx, err)
	}
	token := authHeader[7:]

	certs := ua.jwks.Certs()
	parsedToken, err := jwt.Parse(token, ua.verifyJWTSignature(ctx, certs))
	if err == nil {
		setAccountID(ctx, parsedToken)
		return nil
	}

	if err := ua.jwks.RenewCerts(ctx); err != nil {
		err = errors.New("error on renewing the certificates").WithKind(errors.KindInternal)
		ua.logger.Error(ctx, err)
		return http.NewErrorResponse(ctx, err)
	}
	certs = ua.jwks.Certs()

	parsedToken, err = jwt.Parse(token, ua.verifyJWTSignature(ctx, certs))
	if err != nil {
		ua.logger.Error(ctx, err)
		return http.NewErrorResponse(ctx, err)
	}

	setAccountID(ctx, parsedToken)
	return nil
}

func setAccountID(ctx *routing.Context, token *jwt.Token) {
	claims, _ := token.Claims.(jwt.MapClaims)
	sub, _ := claims["sub"].(string)
	// removes auth provider prefix 'auth0|' to get only the user identifier.
	accountID := strings.Split(sub, "|")[1]
	fmt.Println(accountID)
	ctx.SetUserValue(ContextKeyAccountID, accountID)
}

func (ua AccountAuthenticator) verifyJWTSignature(ctx context.Context, certs map[string]string) func(token *jwt.Token) (interface{}, error) {
	return func(token *jwt.Token) (interface{}, error) {
		kid, ok := token.Header["kid"].(string)
		if !ok {
			err := errors.New("token's kid header not found").WithKind(errors.KindUnauthenticated)
			return nil, err
		}
		cert, ok := certs[kid]
		if !ok {
			err := errors.New("cert key not found").WithKind(errors.KindUnauthenticated)
			return nil, err
		}

		result, err := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
		if err != nil {
			err := errors.New("error trying to validate signature").WithKind(errors.KindInternal)
			return nil, err
		}

		return result, nil
	}
}
