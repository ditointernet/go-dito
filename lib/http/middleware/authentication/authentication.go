package authentication

import (
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/ditointernet/go-dito/lib/errors"
	"github.com/gin-gonic/gin"
)

// AccountAuthenticator structure responsible for handling request authentication.
type AccountAuthenticator struct {
	jwks JWKSClient
}

// NewAccountAuthenticator creates a new instance of the AccountAuthenticator structure.
func NewAccountAuthenticator(jwks JWKSClient) (AccountAuthenticator, error) {
	if jwks == nil {
		return AccountAuthenticator{}, errors.NewMissingRequiredDependency("jwks")
	}

	return AccountAuthenticator{
		jwks: jwks,
	}, nil
}

// MustNewAccountAuthenticator creates a new instance of the AccountAuthenticator structure.
// It panics if any error is found.
func MustNewAccountAuthenticator(jwks JWKSClient) AccountAuthenticator {
	auth, err := NewAccountAuthenticator(jwks)
	if err != nil {
		panic(err)
	}

	return auth
}

// Authenticate is responsible for verify if the request is authenticated.
//
// It tries to authenticate the token with the certifications on memory,
// if it fails, the certifications are renewed and the authentication is
// run again.
func (ua AccountAuthenticator) Authenticate(ctx *gin.Context) {
	authHeader := string(ctx.GetHeader("Authorization"))
	if len(authHeader) < 7 || strings.ToLower(authHeader[:7]) != "bearer " || authHeader[7:] == "" {
		ctx.Error(ErrMissingOrInvalidAuthenticationToken)
		ctx.Abort()
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
		ctx.Error(ErrRenewCertificates)
		ctx.Abort()
		return
	}

	certs = ua.jwks.Certs()
	parsedToken, err = jwt.Parse(token, verifyJWTSignature(certs))
	if err != nil {
		ctx.Error(ErrInvalidJWTToken)
		ctx.Abort()
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
