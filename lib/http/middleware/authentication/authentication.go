package authentication

import (
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/ditointernet/go-dito/lib/errors"
	ditohttp "github.com/ditointernet/go-dito/lib/http"
	"github.com/ditointernet/go-dito/lib/http/infra"
	"github.com/gin-gonic/gin"
)

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
	jwks infra.JWKSClient
}

// NewAccountAuthenticator creates a new instance of the AccountAuthenticator structure
func NewAccountAuthenticator(jwks infra.JWKSClient) (AccountAuthenticator, error) {
	if jwks == nil {
		return AccountAuthenticator{}, errors.NewMissingRequiredDependency("jwks")
	}

	return AccountAuthenticator{
		jwks: jwks,
	}, nil
}

// MustNewAccountAuthenticator creates a new instance of the AccountAuthenticator structure.
// It panics if any error is found.
func MustNewAccountAuthenticator(jwks infra.JWKSClient) AccountAuthenticator {
	auth, err := NewAccountAuthenticator(jwks)
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
func (ua AccountAuthenticator) Authenticate(ctx *gin.Context) {
	authHeader := string(ctx.GetHeader("Authorization"))
	if len(authHeader) < 7 || strings.ToLower(authHeader[:7]) != "bearer " || authHeader[7:] == "" {
		err := errors.
			New("missing or invalid authentication token").
			WithKind(errors.KindUnauthenticated).
			WithCode(CodeTypeMissingBearerToken)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, ditohttp.NewErrorResponse(ctx, err))
		return
	}
	token := authHeader[7:]

	certs := ua.jwks.Certs()
	parsedToken, err := jwt.Parse(token, verifyJWTSignature(certs))
	if err == nil {
		injectClaims(ctx, parsedToken)
		return
	}

	if err := ua.jwks.RenewCerts(ctx); err != nil {
		err = errors.
			New("error on renewing the certificates").
			WithKind(errors.KindInternal).
			WithCode(CodeTypeErrorOnRenewingCerts)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, ditohttp.NewErrorResponse(ctx, err))
		return
	}
	certs = ua.jwks.Certs()

	parsedToken, err = jwt.Parse(token, verifyJWTSignature(certs))
	if err != nil {
		err = errors.New(err.Error()).WithKind(errors.KindUnauthenticated).WithCode(CodeTypeErrorOnParsingJWTToken)
		ditohttp.NewErrorResponse(ctx, err)
		return
	}

	injectClaims(ctx, parsedToken)
	return
}

func injectClaims(ctx *gin.Context, token *jwt.Token) {
	claims, _ := token.Claims.(jwt.MapClaims)
	ctx.Set("claims", claims)
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
