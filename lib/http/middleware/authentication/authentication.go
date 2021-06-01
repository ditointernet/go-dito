package authentication

import (
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/ditointernet/go-dito/lib/errors"
	"github.com/ditointernet/go-dito/lib/http"
	routing "github.com/jackwhelpton/fasthttp-routing/v2"
)

// ContextKeyAccountID is the key used to retrieve and save accountId into the context
const ContextKeyAccountID string = "account-id"

const (
	// CodeTypeMissingBearerToken indicates that the bearer token was not provided
	CodeTypeMissingBearerToken errors.CodeType = "MISSING_BEARER_TOKEN"
	// CodeTypeErrorOnRenewingCerts indicates that the application couldnt renew the JWKS certificates
	CodeTypeErrorOnRenewingCerts errors.CodeType = "COULD_NOT_RENEW_CERTS"
	// CodeTypeErrorOnParsingJWTToken indicates that the application couldn't parse the JWT token
	CodeTypeErrorOnParsingJWTToken errors.CodeType = "COULD_NOT_HANDLE_TOKEN"
)

// AccountAuthenticator structure responsible for handling request authentication
type AccountAuthenticator struct {
	logger logger
	jwks   jwksClient
}

// NewAccountAuthenticator creates a new instance of the AccountAuthenticator structure
func NewAccountAuthenticator(logger logger, jwks jwksClient) (AccountAuthenticator, error) {
	if logger == nil {
		return AccountAuthenticator{}, errors.NewMissingRequiredDependency("logger")
	}

	if jwks == nil {
		return AccountAuthenticator{}, errors.NewMissingRequiredDependency("jwks")
	}

	return AccountAuthenticator{
		logger: logger,
		jwks:   jwks,
	}, nil
}

// NewAccountAuthenticator creates a new instance of the AccountAuthenticator structure.
// It panics if any error is found.
func MustNewAccountAuthenticator(logger logger, jwks jwksClient) AccountAuthenticator {
	auth, err := NewAccountAuthenticator(logger, jwks)
	if err != nil {
		panic(err)
	}

	return auth
}

// Authenticate is responsible for verify if the request is authenticated
//
// It tries to authenticate the token with the certifications on memory,
// if it fails, the certifications are renewed and the authentication is
// run again.
func (ua AccountAuthenticator) Authenticate(ctx *routing.Context) error {
	authHeader := string(ctx.Request.Header.Peek("Authorization"))
	if len(authHeader) < 7 || strings.ToLower(authHeader[:7]) != "bearer " || authHeader[7:] == "" {
		err := errors.New("missing or invalid authentication token").WithKind(errors.KindUnauthenticated).WithCode(CodeTypeMissingBearerToken)
		ua.logger.Error(ctx, err)
		return http.NewErrorResponse(ctx, err)
	}
	token := authHeader[7:]

	certs := ua.jwks.Certs()
	parsedToken, err := jwt.Parse(token, verifyJWTSignature(certs))
	if err == nil {
		setAccountID(ctx, parsedToken)
		return nil
	}

	if err := ua.jwks.RenewCerts(ctx); err != nil {
		err = errors.New("error on renewing the certificates").WithKind(errors.KindInternal).WithCode(CodeTypeErrorOnRenewingCerts)
		ua.logger.Error(ctx, err)
		return http.NewErrorResponse(ctx, err)
	}
	certs = ua.jwks.Certs()

	parsedToken, err = jwt.Parse(token, verifyJWTSignature(certs))
	if err != nil {
		ua.logger.Error(ctx, err)
		err = errors.New(err.Error()).WithKind(errors.KindUnauthenticated).WithCode(CodeTypeErrorOnParsingJWTToken)
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
	ctx.SetUserValue(ContextKeyAccountID, accountID)
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
			return nil, errors.New("error trying to validate JWT signature")
		}

		return result, nil
	}
}
