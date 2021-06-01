package authentication

import (
	"fmt"
	"net"
	"testing"

	"github.com/ditointernet/go-dito/lib/errors"
	"github.com/ditointernet/go-dito/lib/http/middleware/authentication/mocks"
	"github.com/golang/mock/gomock"
	routing "github.com/jackwhelpton/fasthttp-routing/v2"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func newCtx() *routing.Context {
	req := fasthttp.AcquireRequest()
	reqCtx := &fasthttp.RequestCtx{}
	reqCtx.Init(req, &net.IPAddr{IP: net.IPv4(127, 0, 0, 1)}, nil)
	return &routing.Context{
		RequestCtx: reqCtx,
	}
}

func newCtxWithHeaders(headers map[string]string) *routing.Context {
	ctx := newCtx()

	for key, value := range headers {
		ctx.Request.Header.Set(key, value)
	}

	return ctx
}
func TestAccountAuthenticator_Authenticate(t *testing.T) {
	var logger *mocks.Mocklogger
	var jwks *mocks.MockjwksClient

	withMock := func(runner func(t *testing.T, ua AccountAuthenticator)) func(t *testing.T) {
		return func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			logger = mocks.NewMocklogger(ctrl)
			jwks = mocks.NewMockjwksClient(ctrl)
			middleware, _ := NewAccountAuthenticator(logger, jwks)
			runner(t, middleware)
		}
	}

	t.Run("should return Unauthorized APIError when the Authorization Header is not given",
		withMock(func(t *testing.T, ua AccountAuthenticator) {
			logger.EXPECT().Error(gomock.Any(), gomock.Any())
			ctx := newCtx()

			err := ua.Authenticate(ctx)

			expected := `{"error":{"code":"MISSING_BEARER_TOKEN","message":"missing or invalid authentication token"}}`
			assert.Equal(t, expected, err.Error())
		}))

	t.Run("should return Unauthorized APIError when the Authorization Header given is not bearer",
		withMock(func(t *testing.T, ua AccountAuthenticator) {
			logger.EXPECT().Error(gomock.Any(), gomock.Any())
			ctx := newCtxWithHeaders(map[string]string{"Authorization": "basic64 "})
			err := ua.Authenticate(ctx)

			expected := `{"error":{"code":"MISSING_BEARER_TOKEN","message":"missing or invalid authentication token"}}`
			assert.Equal(t, expected, err.Error())
		}))

	t.Run("should return Unauthorized APIError when the Authorization Header given does not contain the bearer token",
		withMock(func(t *testing.T, ua AccountAuthenticator) {
			logger.EXPECT().Error(gomock.Any(), gomock.Any())

			ctx := newCtxWithHeaders(map[string]string{"Authorization": "bearer "})
			err := ua.Authenticate(ctx)

			expected := `{"error":{"code":"MISSING_BEARER_TOKEN","message":"missing or invalid authentication token"}}`
			assert.Equal(t, expected, err.Error())
		}))

	t.Run("should return internal server error when an error is found trying to renew the jwks certificates",
		withMock(func(t *testing.T, ua AccountAuthenticator) {
			logger.EXPECT().
				Error(gomock.Any(), gomock.Any())
			ctx := newCtxWithHeaders(map[string]string{"Authorization": "bearer t"})

			jwks.EXPECT().
				Certs().
				Return(map[string]string{})

			jwks.EXPECT().
				RenewCerts(ctx).
				Return(errors.New("testErr"))

			err := ua.Authenticate(ctx)

			expected := `{"error":{"code":"COULD_NOT_RENEW_CERTS","message":"error on renewing the certificates"}}`
			assert.Equal(t, expected, err.Error())
		}))

	t.Run("should return Unauthorized APIError when the kid header is not found on the token",
		withMock(func(t *testing.T, ua AccountAuthenticator) {
			logger.EXPECT().
				Error(gomock.Any(), gomock.Any())

			token := `mock-token`
			ctx := newCtxWithHeaders(map[string]string{"Authorization": "bearer " + token})

			jwks.EXPECT().
				Certs().
				Return(map[string]string{}).MinTimes(2)

			jwks.EXPECT().
				RenewCerts(ctx).
				Return(nil)

			err := ua.Authenticate(ctx)

			expected := `{"error":{"code":"ERROR_ON_PARSING_JWT","message":"token's kid header not found"}}`
			assert.Equal(t, expected, err.Error())
		}))

	t.Run("should return Unauthorized APIError when the certificate for the token is not found",
		withMock(func(t *testing.T, ua AccountAuthenticator) {
			logger.EXPECT().
				Error(gomock.Any(), gomock.Any())
			token := `mock-token`
			ctx := newCtxWithHeaders(map[string]string{"Authorization": "bearer " + token})

			jwks.EXPECT().
				Certs().
				Return(map[string]string{}).MinTimes(2)

			jwks.EXPECT().
				RenewCerts(ctx).
				Return(nil)

			err := ua.Authenticate(ctx)

			expected := `{"error":{"code":"ERROR_ON_PARSING_JWT","message":"cert key not found"}}`
			assert.Equal(t, expected, err.Error())
		}))

	t.Run("should return Unauthorized APIError when the token signature does not match the certificate",
		withMock(func(t *testing.T, ua AccountAuthenticator) {
			logger.EXPECT().
				Error(gomock.Any(), gomock.Any())

			token := `mock-token`
			ctx := newCtxWithHeaders(map[string]string{"Authorization": "bearer " + token})
			jwks.EXPECT().
				Certs().
				Return(map[string]string{"test": ""}).MinTimes(2)

			jwks.EXPECT().
				RenewCerts(ctx).
				Return(nil)

			err := ua.Authenticate(ctx)

			expected := `{"error":{"code":"ERROR_ON_PARSING_JWT","message":"error trying to validate JWT signature"}}`
			assert.Equal(t, expected, err.Error())
		}))

	t.Run("should not interfere on the request result when the token is authenticated",
		withMock(func(t *testing.T, ua AccountAuthenticator) {
			token := `mock-token`

			cert := `-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAyDrJZeDYxOWYhHGVJqw5
jD/HucDR1euWZoUleSp/ugPi6SscAXHzXoyZu301n5P/x/yVPpDAdI/GRuP2R7/H
UC3hCt1rRAKGFD+UNcG7PQItFteMXCLnuFhY3PTWRw2Y4NLBk+eNp0RCvJvKyp8S
kDN8Rf2t+9yWbn4PUIx4lRyIrUNg1oK7oQq0RDEl71ZHiqRAuqf7+htPKG+o067K
Yal2VWtgIwoY1wSn6L6RPw58NghRVJVvZIjSdJzsTcWWnxgHK0yTpjMVd0fL2shR
0Vb7DBMjJp8tnPPizPgtGpRL36MDy8sWLcFx2WYWW3Ga4rpQHnsXXy1vdRyAHtqn
QBZVkmIyxvfQhq8QnW/TLRWIszliQ+brljA/zwVf5lyW8o0AHDhnnY91L3+A//eo
40A8wgLeHRaJIz6UTVnjXzpP18d6bjOHnQc7aedCSpp6gHiKGHii2GEJqvNiXcC0
bwkj44kPHTji8fSfxLYjbRktALxqnTQNKlvCh/4ddFEJbGiBJVFo4p0StQ8SAmOW
7CuEMuLCKvGIomC1NS3gVgY9oi5wiRTaU6Xf+VnVMhq5IEqwC0oJS9yuRKPVEE3x
P1z+rQl1jiFIGpLjZkDLCMDzcZj3zyrVBqFXMKOFnntlw3/2JVKswDUtrc/bI2ES
CfIuR66YJ+RDPXtZXAaiX9UCAwEAAQ==
-----END PUBLIC KEY-----`

			jwks.EXPECT().
				Certs().
				Return(map[string]string{"test": cert})

			ctx := newCtxWithHeaders(map[string]string{"Authorization": "bearer " + token})
			err := ua.Authenticate(ctx)
			assert.Equal(t, err, nil)
		}))

	t.Run("should not interfere on the request result when the token is authenticated after certificates renew",
		withMock(func(t *testing.T, ua AccountAuthenticator) {
			token := `mock-token`

			cert := `-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAyDrJZeDYxOWYhHGVJqw5
jD/HucDR1euWZoUleSp/ugPi6SscAXHzXoyZu301n5P/x/yVPpDAdI/GRuP2R7/H
UC3hCt1rRAKGFD+UNcG7PQItFteMXCLnuFhY3PTWRw2Y4NLBk+eNp0RCvJvKyp8S
kDN8Rf2t+9yWbn4PUIx4lRyIrUNg1oK7oQq0RDEl71ZHiqRAuqf7+htPKG+o067K
Yal2VWtgIwoY1wSn6L6RPw58NghRVJVvZIjSdJzsTcWWnxgHK0yTpjMVd0fL2shR
0Vb7DBMjJp8tnPPizPgtGpRL36MDy8sWLcFx2WYWW3Ga4rpQHnsXXy1vdRyAHtqn
QBZVkmIyxvfQhq8QnW/TLRWIszliQ+brljA/zwVf5lyW8o0AHDhnnY91L3+A//eo
40A8wgLeHRaJIz6UTVnjXzpP18d6bjOHnQc7aedCSpp6gHiKGHii2GEJqvNiXcC0
bwkj44kPHTji8fSfxLYjbRktALxqnTQNKlvCh/4ddFEJbGiBJVFo4p0StQ8SAmOW
7CuEMuLCKvGIomC1NS3gVgY9oi5wiRTaU6Xf+VnVMhq5IEqwC0oJS9yuRKPVEE3x
P1z+rQl1jiFIGpLjZkDLCMDzcZj3zyrVBqFXMKOFnntlw3/2JVKswDUtrc/bI2ES
CfIuR66YJ+RDPXtZXAaiX9UCAwEAAQ==
-----END PUBLIC KEY-----`

			ctx := newCtxWithHeaders(map[string]string{"Authorization": "bearer " + token})
			gomock.InOrder(
				jwks.EXPECT().
					Certs().
					Return(map[string]string{}),

				jwks.EXPECT().
					RenewCerts(ctx).
					Return(nil),

				jwks.EXPECT().
					Certs().
					Return(map[string]string{"test": cert}),
			)

			err := ua.Authenticate(ctx)
			assert.Equal(t, err, nil)
		}))

	t.Run("should set the account id on the context when the given token is authenticated",
		withMock(func(t *testing.T, ua AccountAuthenticator) {
			token := `mock-token`
			ctx := newCtxWithHeaders(map[string]string{"Authorization": "bearer " + token})
			cert := `-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAyDrJZeDYxOWYhHGVJqw5
jD/HucDR1euWZoUleSp/ugPi6SscAXHzXoyZu301n5P/x/yVPpDAdI/GRuP2R7/H
UC3hCt1rRAKGFD+UNcG7PQItFteMXCLnuFhY3PTWRw2Y4NLBk+eNp0RCvJvKyp8S
kDN8Rf2t+9yWbn4PUIx4lRyIrUNg1oK7oQq0RDEl71ZHiqRAuqf7+htPKG+o067K
Yal2VWtgIwoY1wSn6L6RPw58NghRVJVvZIjSdJzsTcWWnxgHK0yTpjMVd0fL2shR
0Vb7DBMjJp8tnPPizPgtGpRL36MDy8sWLcFx2WYWW3Ga4rpQHnsXXy1vdRyAHtqn
QBZVkmIyxvfQhq8QnW/TLRWIszliQ+brljA/zwVf5lyW8o0AHDhnnY91L3+A//eo
40A8wgLeHRaJIz6UTVnjXzpP18d6bjOHnQc7aedCSpp6gHiKGHii2GEJqvNiXcC0
bwkj44kPHTji8fSfxLYjbRktALxqnTQNKlvCh/4ddFEJbGiBJVFo4p0StQ8SAmOW
7CuEMuLCKvGIomC1NS3gVgY9oi5wiRTaU6Xf+VnVMhq5IEqwC0oJS9yuRKPVEE3x
P1z+rQl1jiFIGpLjZkDLCMDzcZj3zyrVBqFXMKOFnntlw3/2JVKswDUtrc/bI2ES
CfIuR66YJ+RDPXtZXAaiX9UCAwEAAQ==
-----END PUBLIC KEY-----`

			jwks.EXPECT().
				Certs().
				Return(map[string]string{"test": cert})

			ua.Authenticate(ctx)

			accountID, ok := ctx.Value(ContextKeyAccountID).(string)
			fmt.Println(accountID, ok)
			assert.Equal(t, "test", accountID)
		}))

	t.Run("should set the contextId into the context when the given token is authenticated after certificates renew",
		withMock(func(t *testing.T, ua AccountAuthenticator) {
			token := `mock-token`
			ctx := newCtxWithHeaders(map[string]string{"Authorization": "bearer " + token})
			cert := `-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAyDrJZeDYxOWYhHGVJqw5
jD/HucDR1euWZoUleSp/ugPi6SscAXHzXoyZu301n5P/x/yVPpDAdI/GRuP2R7/H
UC3hCt1rRAKGFD+UNcG7PQItFteMXCLnuFhY3PTWRw2Y4NLBk+eNp0RCvJvKyp8S
kDN8Rf2t+9yWbn4PUIx4lRyIrUNg1oK7oQq0RDEl71ZHiqRAuqf7+htPKG+o067K
Yal2VWtgIwoY1wSn6L6RPw58NghRVJVvZIjSdJzsTcWWnxgHK0yTpjMVd0fL2shR
0Vb7DBMjJp8tnPPizPgtGpRL36MDy8sWLcFx2WYWW3Ga4rpQHnsXXy1vdRyAHtqn
QBZVkmIyxvfQhq8QnW/TLRWIszliQ+brljA/zwVf5lyW8o0AHDhnnY91L3+A//eo
40A8wgLeHRaJIz6UTVnjXzpP18d6bjOHnQc7aedCSpp6gHiKGHii2GEJqvNiXcC0
bwkj44kPHTji8fSfxLYjbRktALxqnTQNKlvCh/4ddFEJbGiBJVFo4p0StQ8SAmOW
7CuEMuLCKvGIomC1NS3gVgY9oi5wiRTaU6Xf+VnVMhq5IEqwC0oJS9yuRKPVEE3x
P1z+rQl1jiFIGpLjZkDLCMDzcZj3zyrVBqFXMKOFnntlw3/2JVKswDUtrc/bI2ES
CfIuR66YJ+RDPXtZXAaiX9UCAwEAAQ==
-----END PUBLIC KEY-----`

			gomock.InOrder(
				jwks.EXPECT().
					Certs().
					Return(map[string]string{}),

				jwks.EXPECT().
					RenewCerts(ctx).
					Return(nil),

				jwks.EXPECT().
					Certs().
					Return(map[string]string{"test": cert}),
			)

			ua.Authenticate(ctx)

			accountID, _ := ctx.UserValue(ContextKeyAccountID).(string)
			fmt.Println(accountID)
			assert.Equal(t, "test", accountID)
		}))
}
