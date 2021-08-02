package authentication

// import (
// 	"fmt"
// 	"net"
// 	"testing"

// 	"github.com/ditointernet/go-dito/lib/errors"
// 	"github.com/ditointernet/go-dito/lib/http/mocks"
// 	"github.com/golang/mock/gomock"
// 	routing "github.com/jackwhelpton/fasthttp-routing/v2"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/valyala/fasthttp"
// )

// func newCtx() *routing.Context {
// 	req := fasthttp.AcquireRequest()
// 	reqCtx := &fasthttp.RequestCtx{}
// 	reqCtx.Init(req, &net.IPAddr{IP: net.IPv4(127, 0, 0, 1)}, nil)
// 	return &routing.Context{
// 		RequestCtx: reqCtx,
// 	}
// }

// func newCtxWithHeaders(headers map[string]string) *routing.Context {
// 	ctx := newCtx()

// 	for key, value := range headers {
// 		ctx.Request.Header.Set(key, value)
// 	}

// 	return ctx
// }
// func TestAccountAuthenticator_Authenticate(t *testing.T) {
// 	var logger *mocks.MockLogger
// 	var jwks *mocks.MockJWKSClient

// 	withMock := func(runner func(t *testing.T, ua AccountAuthenticator)) func(t *testing.T) {
// 		return func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			logger = mocks.NewMockLogger(ctrl)
// 			jwks = mocks.NewMockJWKSClient(ctrl)
// 			middleware, _ := NewAccountAuthenticator(logger, jwks)
// 			runner(t, middleware)
// 		}
// 	}

// 	t.Run("should return Unauthorized APIError when the Authorization Header is not given",
// 		withMock(func(t *testing.T, ua AccountAuthenticator) {
// 			logger.EXPECT().Error(gomock.Any(), gomock.Any())
// 			ctx := newCtx()

// 			err := ua.Authenticate(ctx)

// 			expected := `{"error":{"code":"MISSING_BEARER_TOKEN","message":"missing or invalid authentication token"}}`
// 			assert.Equal(t, expected, err.Error())
// 		}))

// 	t.Run("should return Unauthorized APIError when the Authorization Header given is not bearer",
// 		withMock(func(t *testing.T, ua AccountAuthenticator) {
// 			logger.EXPECT().Error(gomock.Any(), gomock.Any())
// 			ctx := newCtxWithHeaders(map[string]string{"Authorization": "basic64 "})
// 			err := ua.Authenticate(ctx)

// 			expected := `{"error":{"code":"MISSING_BEARER_TOKEN","message":"missing or invalid authentication token"}}`
// 			assert.Equal(t, expected, err.Error())
// 		}))

// 	t.Run("should return Unauthorized APIError when the Authorization Header given does not contain the bearer token",
// 		withMock(func(t *testing.T, ua AccountAuthenticator) {
// 			logger.EXPECT().Error(gomock.Any(), gomock.Any())

// 			ctx := newCtxWithHeaders(map[string]string{"Authorization": "bearer "})
// 			err := ua.Authenticate(ctx)

// 			expected := `{"error":{"code":"MISSING_BEARER_TOKEN","message":"missing or invalid authentication token"}}`
// 			assert.Equal(t, expected, err.Error())
// 		}))

// 	t.Run("should return internal server error when an error is found trying to renew the jwks certificates",
// 		withMock(func(t *testing.T, ua AccountAuthenticator) {
// 			logger.EXPECT().
// 				Error(gomock.Any(), gomock.Any())
// 			ctx := newCtxWithHeaders(map[string]string{"Authorization": "bearer t"})

// 			jwks.EXPECT().
// 				Certs().
// 				Return(map[string]string{})

// 			jwks.EXPECT().
// 				RenewCerts(ctx).
// 				Return(errors.New("testErr"))

// 			err := ua.Authenticate(ctx)

// 			expected := `{"error":{"code":"COULD_NOT_RENEW_CERTS","message":"error on renewing the certificates"}}`
// 			assert.Equal(t, expected, err.Error())
// 		}))

// 	t.Run("should return Unauthorized APIError when the kid header is not found on the token",
// 		withMock(func(t *testing.T, ua AccountAuthenticator) {
// 			logger.EXPECT().
// 				Error(gomock.Any(), gomock.Any())

// 			token := `eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwczovL3Rlc3QuYXV0aDAuY29tLyIsInN1YiI6ImF1dGgwfHRlc3QiLCJhdWQiOiJodHRwOi8vb3BlbmFwaS5kaXRvLmNvbS5ici8iLCJpYXQiOjEwNDEzNzkyMDAwLCJleHAiOjEwNDEzNzkyMDAwLCJhenAiOiJBSDRhMHZJenRFQjBZT0ViRXlkV2p2d1BvV2lrMEQ4biIsInNjb3BlIjoib2ZmbGluZV9hY2Nlc3MiLCJndHkiOiJwYXNzd29yZCJ9.EWAJ-SKeE2ck0oXevxED1HDdGsDuhMJ3o2r8m80swg8gvbGQDJ2bkfoM9Idi-x0ztmyW5H_2SP9AnZrPcqZqwfJn2W-PSLQjWwNKY3kLSKVOevlstGelX4zu2LGapBlsszm1P5LzYpLl85U-47beo3ezwq0fM4b5K4bl9GhMmDmndsqZ_eyVtfHOkHfBKQzo72_trsyEk86Oa6jrmxNjPjwqIfCxCFwjPGK9vkOGE505eZCJ9at6ileRP4i7aO_KGxzfME2NCfvCTXdpeW5tTSdeXuQ-m-u_k3dwhQviKhY8k6D2bW-bns5u-5fl01v0-nFYE3-LiJCPTb6Hj5dRb171eNdoCG2L092Vvz68nWvQLNuzSU8xINnCaZ-uoAEDmsKWJtvRS5Xi_CAto7awLbHHb_S9zzXNe1QRTJAg2Tdec1zwzciskeYvhVgVHMrYfDVdMxTLyVGNRIva2SIbkNZBxATXqBfLSZhAvXG3h58e8QSP7Q7WaK4Gqzxrqa5OLiQRGYPCgOLDP0zIN3R8i5gWuqh3YsCmwV8B0k2bHJ4OcL_23-5LrPHRs2GvZ0_o1eP8ZhONJTuKLmVk4AIZ0zqvLyMEFCOj_HpPExm4ranhd67wWR3HViL-3fmJJGI53LfpxEUW6CcSxnl3IrTJGxsE7qnD8be9oPG7X6b-xFw`
// 			ctx := newCtxWithHeaders(map[string]string{"Authorization": "bearer " + token})

// 			jwks.EXPECT().
// 				Certs().
// 				Return(map[string]string{}).MinTimes(2)

// 			jwks.EXPECT().
// 				RenewCerts(ctx).
// 				Return(nil)

// 			err := ua.Authenticate(ctx)

// 			expected := `{"error":{"code":"COULD_NOT_HANDLE_TOKEN","message":"token's kid header not found"}}`
// 			assert.Equal(t, expected, err.Error())
// 		}))

// 	t.Run("should return Unauthorized APIError when the certificate for the token is not found",
// 		withMock(func(t *testing.T, ua AccountAuthenticator) {
// 			logger.EXPECT().
// 				Error(gomock.Any(), gomock.Any())
// 			token := `eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3QifQ.eyJpc3MiOiJodHRwczovL3Rlc3QuYXV0aDAuY29tLyIsInN1YiI6ImF1dGgwfHRlc3QiLCJhdWQiOiJodHRwOi8vb3BlbmFwaS5kaXRvLmNvbS5ici8iLCJpYXQiOjEwNDEzNzkyMDAwLCJleHAiOjEwNDEzNzkyMDAwLCJhenAiOiJBSDRhMHZJenRFQjBZT0ViRXlkV2p2d1BvV2lrMEQ4biIsInNjb3BlIjoib2ZmbGluZV9hY2Nlc3MiLCJndHkiOiJwYXNzd29yZCJ9.eM-Z_PRPUYYYIhOxvZkHZw3bQvl_DL8WZ9bt4yT0iDMuwcunb7uEY5dnXDAqfYe6KYndYX-dqeMzvBBTnJ2vwLY_EPAN6caIZNhgtNO6Nk5GBkQHEQ9d-MkEvCLZHCiiKjlIx_Dz-bpIqNqvhH_t-LPkRwSiwLw0JcfzPOpw5_WlB0ay4Y95NAJI_ielTBc0cyiO-GHNYAqOim4ES2XMS-6fGElhasy6q0MpCQpyiYFJcPPW6HZ9gdYu_dOpY9FdZCR1Vb13JnEjft-oll0YTZW9YYzVEe47dNjlqhaUeM12UrazeZWdk2x9ToI7E3OVfleOV7jlLzTatEuYdYRQzlFPTndpQl8-_GKdK3oziy3UDcpUxFhaxyhuv-80V5bMgdOWYHxm8Ykf6RgoOu-5yw9BBXIHaHVVrtImAZtenhsQTikE8EDtm2OqSxfv14zja7MswKy5opPeHralb5OMFHrc_7lYowqJNq31b77niDpme3-KwJUAHZvHLoMfF5fW0deKBmahdySGwfm4Kdo0xY1IWtPwIoWZayusP1I2iCl0fzeECn8JGI2Ml8BYyKdzjMy7eEVw6RWgdzhO8Q5l4wG7VqSH08Nr_z25Cq_9FupMyhLzcPU2aRBEIzFOiXZMw-aAn5x1X2XJ0Sr6CaIsRNMvG83EojhufSG07BjDl8Q`
// 			ctx := newCtxWithHeaders(map[string]string{"Authorization": "bearer " + token})

// 			jwks.EXPECT().
// 				Certs().
// 				Return(map[string]string{}).MinTimes(2)

// 			jwks.EXPECT().
// 				RenewCerts(ctx).
// 				Return(nil)

// 			err := ua.Authenticate(ctx)

// 			expected := `{"error":{"code":"COULD_NOT_HANDLE_TOKEN","message":"cert key not found"}}`
// 			assert.Equal(t, expected, err.Error())
// 		}))

// 	t.Run("should return Unauthorized APIError when the token signature does not match the certificate",
// 		withMock(func(t *testing.T, ua AccountAuthenticator) {
// 			logger.EXPECT().
// 				Error(gomock.Any(), gomock.Any())

// 			token := `eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3QifQ.eyJpc3MiOiJodHRwczovL3Rlc3QuYXV0aDAuY29tLyIsInN1YiI6ImF1dGgwfHRlc3QiLCJhdWQiOiJodHRwOi8vb3BlbmFwaS5kaXRvLmNvbS5ici8iLCJpYXQiOjEwNDEzNzkyMDAwLCJleHAiOjEwNDEzNzkyMDAwLCJhenAiOiJBSDRhMHZJenRFQjBZT0ViRXlkV2p2d1BvV2lrMEQ4biIsInNjb3BlIjoib2ZmbGluZV9hY2Nlc3MiLCJndHkiOiJwYXNzd29yZCJ9.eM-Z_PRPUYYYIhOxvZkHZw3bQvl_DL8WZ9bt4yT0iDMuwcunb7uEY5dnXDAqfYe6KYndYX-dqeMzvBBTnJ2vwLY_EPAN6caIZNhgtNO6Nk5GBkQHEQ9d-MkEvCLZHCiiKjlIx_Dz-bpIqNqvhH_t-LPkRwSiwLw0JcfzPOpw5_WlB0ay4Y95NAJI_ielTBc0cyiO-GHNYAqOim4ES2XMS-6fGElhasy6q0MpCQpyiYFJcPPW6HZ9gdYu_dOpY9FdZCR1Vb13JnEjft-oll0YTZW9YYzVEe47dNjlqhaUeM12UrazeZWdk2x9ToI7E3OVfleOV7jlLzTatEuYdYRQzlFPTndpQl8-_GKdK3oziy3UDcpUxFhaxyhuv-80V5bMgdOWYHxm8Ykf6RgoOu-5yw9BBXIHaHVVrtImAZtenhsQTikE8EDtm2OqSxfv14zja7MswKy5opPeHralb5OMFHrc_7lYowqJNq31b77niDpme3-KwJUAHZvHLoMfF5fW0deKBmahdySGwfm4Kdo0xY1IWtPwIoWZayusP1I2iCl0fzeECn8JGI2Ml8BYyKdzjMy7eEVw6RWgdzhO8Q5l4wG7VqSH08Nr_z25Cq_9FupMyhLzcPU2aRBEIzFOiXZMw-aAn5x1X2XJ0Sr6CaIsRNMvG83EojhufSG07BjDl8Q`
// 			ctx := newCtxWithHeaders(map[string]string{"Authorization": "bearer " + token})
// 			jwks.EXPECT().
// 				Certs().
// 				Return(map[string]string{"test": ""}).MinTimes(2)

// 			jwks.EXPECT().
// 				RenewCerts(ctx).
// 				Return(nil)

// 			err := ua.Authenticate(ctx)

// 			expected := `{"error":{"code":"COULD_NOT_HANDLE_TOKEN","message":"error trying to validate JWT signature"}}`
// 			assert.Equal(t, expected, err.Error())
// 		}))

// 	t.Run("should not interfere on the request result when the token is authenticated",
// 		withMock(func(t *testing.T, ua AccountAuthenticator) {
// 			token := `eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3QifQ.eyJpc3MiOiJodHRwczovL3Rlc3QuYXV0aDAuY29tLyIsInN1YiI6ImF1dGgwfHRlc3QiLCJhdWQiOiJodHRwOi8vb3BlbmFwaS5kaXRvLmNvbS5ici8iLCJleHAiOjEwNDEzNzkyMDAwLCJhenAiOiJBSDRhMHZJenRFQjBZT0ViRXlkV2p2d1BvV2lrMEQ4biIsInNjb3BlIjoib2ZmbGluZV9hY2Nlc3MiLCJndHkiOiJwYXNzd29yZCJ9.ZEkaYsciuCRvS_76SvGT8g8vywdNG3dG2e1LZ0PkQsvtuFSnPdLmU5rQofBGjQGn7wC_aZRpr4Rl-zEhp6BQR6iYsTj5Iw2HUwAvimxtUa-Ztzreiz_seCaFvTdePOmCr9T7pwKc_WnN4rpZBcxx4AD4Y1K2SfIO7eWGvcUSquxdcKVjqUz8c7gokwMvHWv3D6qc7kAz3cm-ug7uj-Mzh1jWm1v_ssFqPQICT28_0aGP7lbAepd_tkI1lKeo7RnerfkNBu1Jh3fiOrjcU_pSWGJ36jZ9ifx5zP0bHfjzUF0r0ijzti-0ufu_fGmrXr73OBzCT54NDdbkUfBYd4l5Uahi3ajNNp6hnNVoBNueYHa7k3J5Svh9Xszij1F5ssSr1L5H-e-6JEtiAqkGGd1rem-leA2qf9HCf0C6RKbyn6j4Pr93Xlu-6F5FWTDcWRSoNXd2sLzLkXwj6B-zQDri5tS2BgBKs8ajHQ038yrW4uw8zS2CXoQ-f4tB63mhW0ESYCqCM0N5w83_P4tC9LWlPXCUZqi_nLYvHC84F5eJRpPlwDREkfzD1P_CE9tatgLNPnZ8FU2AZW5GNOEdtXPPMVoEzdikBVDc3oRdR1HPKTXEjkXidPh44KabsYcNKxP5FRO6fTv3Hwd9_m1nmSD2WAsBE6CqsZv_1NySla5Ewas`

// 			cert := `-----BEGIN PUBLIC KEY-----
// MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAyDrJZeDYxOWYhHGVJqw5
// jD/HucDR1euWZoUleSp/ugPi6SscAXHzXoyZu301n5P/x/yVPpDAdI/GRuP2R7/H
// UC3hCt1rRAKGFD+UNcG7PQItFteMXCLnuFhY3PTWRw2Y4NLBk+eNp0RCvJvKyp8S
// kDN8Rf2t+9yWbn4PUIx4lRyIrUNg1oK7oQq0RDEl71ZHiqRAuqf7+htPKG+o067K
// Yal2VWtgIwoY1wSn6L6RPw58NghRVJVvZIjSdJzsTcWWnxgHK0yTpjMVd0fL2shR
// 0Vb7DBMjJp8tnPPizPgtGpRL36MDy8sWLcFx2WYWW3Ga4rpQHnsXXy1vdRyAHtqn
// QBZVkmIyxvfQhq8QnW/TLRWIszliQ+brljA/zwVf5lyW8o0AHDhnnY91L3+A//eo
// 40A8wgLeHRaJIz6UTVnjXzpP18d6bjOHnQc7aedCSpp6gHiKGHii2GEJqvNiXcC0
// bwkj44kPHTji8fSfxLYjbRktALxqnTQNKlvCh/4ddFEJbGiBJVFo4p0StQ8SAmOW
// 7CuEMuLCKvGIomC1NS3gVgY9oi5wiRTaU6Xf+VnVMhq5IEqwC0oJS9yuRKPVEE3x
// P1z+rQl1jiFIGpLjZkDLCMDzcZj3zyrVBqFXMKOFnntlw3/2JVKswDUtrc/bI2ES
// CfIuR66YJ+RDPXtZXAaiX9UCAwEAAQ==
// -----END PUBLIC KEY-----`

// 			jwks.EXPECT().
// 				Certs().
// 				Return(map[string]string{"test": cert})

// 			ctx := newCtxWithHeaders(map[string]string{"Authorization": "bearer " + token})
// 			err := ua.Authenticate(ctx)
// 			assert.Equal(t, err, nil)
// 		}))

// 	t.Run("should not interfere on the request result when the token is authenticated after certificates renew",
// 		withMock(func(t *testing.T, ua AccountAuthenticator) {
// 			token := `eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3QifQ.eyJpc3MiOiJodHRwczovL3Rlc3QuYXV0aDAuY29tLyIsInN1YiI6ImF1dGgwfHRlc3QiLCJhdWQiOiJodHRwOi8vb3BlbmFwaS5kaXRvLmNvbS5ici8iLCJleHAiOjEwNDEzNzkyMDAwLCJhenAiOiJBSDRhMHZJenRFQjBZT0ViRXlkV2p2d1BvV2lrMEQ4biIsInNjb3BlIjoib2ZmbGluZV9hY2Nlc3MiLCJndHkiOiJwYXNzd29yZCJ9.ZEkaYsciuCRvS_76SvGT8g8vywdNG3dG2e1LZ0PkQsvtuFSnPdLmU5rQofBGjQGn7wC_aZRpr4Rl-zEhp6BQR6iYsTj5Iw2HUwAvimxtUa-Ztzreiz_seCaFvTdePOmCr9T7pwKc_WnN4rpZBcxx4AD4Y1K2SfIO7eWGvcUSquxdcKVjqUz8c7gokwMvHWv3D6qc7kAz3cm-ug7uj-Mzh1jWm1v_ssFqPQICT28_0aGP7lbAepd_tkI1lKeo7RnerfkNBu1Jh3fiOrjcU_pSWGJ36jZ9ifx5zP0bHfjzUF0r0ijzti-0ufu_fGmrXr73OBzCT54NDdbkUfBYd4l5Uahi3ajNNp6hnNVoBNueYHa7k3J5Svh9Xszij1F5ssSr1L5H-e-6JEtiAqkGGd1rem-leA2qf9HCf0C6RKbyn6j4Pr93Xlu-6F5FWTDcWRSoNXd2sLzLkXwj6B-zQDri5tS2BgBKs8ajHQ038yrW4uw8zS2CXoQ-f4tB63mhW0ESYCqCM0N5w83_P4tC9LWlPXCUZqi_nLYvHC84F5eJRpPlwDREkfzD1P_CE9tatgLNPnZ8FU2AZW5GNOEdtXPPMVoEzdikBVDc3oRdR1HPKTXEjkXidPh44KabsYcNKxP5FRO6fTv3Hwd9_m1nmSD2WAsBE6CqsZv_1NySla5Ewas`

// 			cert := `-----BEGIN PUBLIC KEY-----
// MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAyDrJZeDYxOWYhHGVJqw5
// jD/HucDR1euWZoUleSp/ugPi6SscAXHzXoyZu301n5P/x/yVPpDAdI/GRuP2R7/H
// UC3hCt1rRAKGFD+UNcG7PQItFteMXCLnuFhY3PTWRw2Y4NLBk+eNp0RCvJvKyp8S
// kDN8Rf2t+9yWbn4PUIx4lRyIrUNg1oK7oQq0RDEl71ZHiqRAuqf7+htPKG+o067K
// Yal2VWtgIwoY1wSn6L6RPw58NghRVJVvZIjSdJzsTcWWnxgHK0yTpjMVd0fL2shR
// 0Vb7DBMjJp8tnPPizPgtGpRL36MDy8sWLcFx2WYWW3Ga4rpQHnsXXy1vdRyAHtqn
// QBZVkmIyxvfQhq8QnW/TLRWIszliQ+brljA/zwVf5lyW8o0AHDhnnY91L3+A//eo
// 40A8wgLeHRaJIz6UTVnjXzpP18d6bjOHnQc7aedCSpp6gHiKGHii2GEJqvNiXcC0
// bwkj44kPHTji8fSfxLYjbRktALxqnTQNKlvCh/4ddFEJbGiBJVFo4p0StQ8SAmOW
// 7CuEMuLCKvGIomC1NS3gVgY9oi5wiRTaU6Xf+VnVMhq5IEqwC0oJS9yuRKPVEE3x
// P1z+rQl1jiFIGpLjZkDLCMDzcZj3zyrVBqFXMKOFnntlw3/2JVKswDUtrc/bI2ES
// CfIuR66YJ+RDPXtZXAaiX9UCAwEAAQ==
// -----END PUBLIC KEY-----`

// 			ctx := newCtxWithHeaders(map[string]string{"Authorization": "bearer " + token})
// 			gomock.InOrder(
// 				jwks.EXPECT().
// 					Certs().
// 					Return(map[string]string{}),

// 				jwks.EXPECT().
// 					RenewCerts(ctx).
// 					Return(nil),

// 				jwks.EXPECT().
// 					Certs().
// 					Return(map[string]string{"test": cert}),
// 			)

// 			err := ua.Authenticate(ctx)
// 			assert.Equal(t, err, nil)
// 		}))

// 	t.Run("should set the account id on the context when the given token is authenticated",
// 		withMock(func(t *testing.T, ua AccountAuthenticator) {
// 			token := `eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3QifQ.eyJpc3MiOiJodHRwczovL3Rlc3QuYXV0aDAuY29tLyIsInN1YiI6ImF1dGgwfHRlc3QiLCJhdWQiOiJodHRwOi8vb3BlbmFwaS5kaXRvLmNvbS5ici8iLCJleHAiOjEwNDEzNzkyMDAwLCJhenAiOiJBSDRhMHZJenRFQjBZT0ViRXlkV2p2d1BvV2lrMEQ4biIsInNjb3BlIjoib2ZmbGluZV9hY2Nlc3MiLCJndHkiOiJwYXNzd29yZCJ9.ZEkaYsciuCRvS_76SvGT8g8vywdNG3dG2e1LZ0PkQsvtuFSnPdLmU5rQofBGjQGn7wC_aZRpr4Rl-zEhp6BQR6iYsTj5Iw2HUwAvimxtUa-Ztzreiz_seCaFvTdePOmCr9T7pwKc_WnN4rpZBcxx4AD4Y1K2SfIO7eWGvcUSquxdcKVjqUz8c7gokwMvHWv3D6qc7kAz3cm-ug7uj-Mzh1jWm1v_ssFqPQICT28_0aGP7lbAepd_tkI1lKeo7RnerfkNBu1Jh3fiOrjcU_pSWGJ36jZ9ifx5zP0bHfjzUF0r0ijzti-0ufu_fGmrXr73OBzCT54NDdbkUfBYd4l5Uahi3ajNNp6hnNVoBNueYHa7k3J5Svh9Xszij1F5ssSr1L5H-e-6JEtiAqkGGd1rem-leA2qf9HCf0C6RKbyn6j4Pr93Xlu-6F5FWTDcWRSoNXd2sLzLkXwj6B-zQDri5tS2BgBKs8ajHQ038yrW4uw8zS2CXoQ-f4tB63mhW0ESYCqCM0N5w83_P4tC9LWlPXCUZqi_nLYvHC84F5eJRpPlwDREkfzD1P_CE9tatgLNPnZ8FU2AZW5GNOEdtXPPMVoEzdikBVDc3oRdR1HPKTXEjkXidPh44KabsYcNKxP5FRO6fTv3Hwd9_m1nmSD2WAsBE6CqsZv_1NySla5Ewas`
// 			ctx := newCtxWithHeaders(map[string]string{"Authorization": "bearer " + token})
// 			cert := `-----BEGIN PUBLIC KEY-----
// MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAyDrJZeDYxOWYhHGVJqw5
// jD/HucDR1euWZoUleSp/ugPi6SscAXHzXoyZu301n5P/x/yVPpDAdI/GRuP2R7/H
// UC3hCt1rRAKGFD+UNcG7PQItFteMXCLnuFhY3PTWRw2Y4NLBk+eNp0RCvJvKyp8S
// kDN8Rf2t+9yWbn4PUIx4lRyIrUNg1oK7oQq0RDEl71ZHiqRAuqf7+htPKG+o067K
// Yal2VWtgIwoY1wSn6L6RPw58NghRVJVvZIjSdJzsTcWWnxgHK0yTpjMVd0fL2shR
// 0Vb7DBMjJp8tnPPizPgtGpRL36MDy8sWLcFx2WYWW3Ga4rpQHnsXXy1vdRyAHtqn
// QBZVkmIyxvfQhq8QnW/TLRWIszliQ+brljA/zwVf5lyW8o0AHDhnnY91L3+A//eo
// 40A8wgLeHRaJIz6UTVnjXzpP18d6bjOHnQc7aedCSpp6gHiKGHii2GEJqvNiXcC0
// bwkj44kPHTji8fSfxLYjbRktALxqnTQNKlvCh/4ddFEJbGiBJVFo4p0StQ8SAmOW
// 7CuEMuLCKvGIomC1NS3gVgY9oi5wiRTaU6Xf+VnVMhq5IEqwC0oJS9yuRKPVEE3x
// P1z+rQl1jiFIGpLjZkDLCMDzcZj3zyrVBqFXMKOFnntlw3/2JVKswDUtrc/bI2ES
// CfIuR66YJ+RDPXtZXAaiX9UCAwEAAQ==
// -----END PUBLIC KEY-----`

// 			jwks.EXPECT().
// 				Certs().
// 				Return(map[string]string{"test": cert})

// 			ua.Authenticate(ctx)

// 			accountID, ok := ctx.Value(ContextKeyAccountID).(string)
// 			fmt.Println(accountID, ok)
// 			assert.Equal(t, "test", accountID)
// 		}))

// 	t.Run("should set the contextId into the context when the given token is authenticated after certificates renew",
// 		withMock(func(t *testing.T, ua AccountAuthenticator) {
// 			token := `eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3QifQ.eyJpc3MiOiJodHRwczovL3Rlc3QuYXV0aDAuY29tLyIsInN1YiI6ImF1dGgwfHRlc3QiLCJhdWQiOiJodHRwOi8vb3BlbmFwaS5kaXRvLmNvbS5ici8iLCJleHAiOjEwNDEzNzkyMDAwLCJhenAiOiJBSDRhMHZJenRFQjBZT0ViRXlkV2p2d1BvV2lrMEQ4biIsInNjb3BlIjoib2ZmbGluZV9hY2Nlc3MiLCJndHkiOiJwYXNzd29yZCJ9.ZEkaYsciuCRvS_76SvGT8g8vywdNG3dG2e1LZ0PkQsvtuFSnPdLmU5rQofBGjQGn7wC_aZRpr4Rl-zEhp6BQR6iYsTj5Iw2HUwAvimxtUa-Ztzreiz_seCaFvTdePOmCr9T7pwKc_WnN4rpZBcxx4AD4Y1K2SfIO7eWGvcUSquxdcKVjqUz8c7gokwMvHWv3D6qc7kAz3cm-ug7uj-Mzh1jWm1v_ssFqPQICT28_0aGP7lbAepd_tkI1lKeo7RnerfkNBu1Jh3fiOrjcU_pSWGJ36jZ9ifx5zP0bHfjzUF0r0ijzti-0ufu_fGmrXr73OBzCT54NDdbkUfBYd4l5Uahi3ajNNp6hnNVoBNueYHa7k3J5Svh9Xszij1F5ssSr1L5H-e-6JEtiAqkGGd1rem-leA2qf9HCf0C6RKbyn6j4Pr93Xlu-6F5FWTDcWRSoNXd2sLzLkXwj6B-zQDri5tS2BgBKs8ajHQ038yrW4uw8zS2CXoQ-f4tB63mhW0ESYCqCM0N5w83_P4tC9LWlPXCUZqi_nLYvHC84F5eJRpPlwDREkfzD1P_CE9tatgLNPnZ8FU2AZW5GNOEdtXPPMVoEzdikBVDc3oRdR1HPKTXEjkXidPh44KabsYcNKxP5FRO6fTv3Hwd9_m1nmSD2WAsBE6CqsZv_1NySla5Ewas`
// 			ctx := newCtxWithHeaders(map[string]string{"Authorization": "bearer " + token})
// 			cert := `-----BEGIN PUBLIC KEY-----
// MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAyDrJZeDYxOWYhHGVJqw5
// jD/HucDR1euWZoUleSp/ugPi6SscAXHzXoyZu301n5P/x/yVPpDAdI/GRuP2R7/H
// UC3hCt1rRAKGFD+UNcG7PQItFteMXCLnuFhY3PTWRw2Y4NLBk+eNp0RCvJvKyp8S
// kDN8Rf2t+9yWbn4PUIx4lRyIrUNg1oK7oQq0RDEl71ZHiqRAuqf7+htPKG+o067K
// Yal2VWtgIwoY1wSn6L6RPw58NghRVJVvZIjSdJzsTcWWnxgHK0yTpjMVd0fL2shR
// 0Vb7DBMjJp8tnPPizPgtGpRL36MDy8sWLcFx2WYWW3Ga4rpQHnsXXy1vdRyAHtqn
// QBZVkmIyxvfQhq8QnW/TLRWIszliQ+brljA/zwVf5lyW8o0AHDhnnY91L3+A//eo
// 40A8wgLeHRaJIz6UTVnjXzpP18d6bjOHnQc7aedCSpp6gHiKGHii2GEJqvNiXcC0
// bwkj44kPHTji8fSfxLYjbRktALxqnTQNKlvCh/4ddFEJbGiBJVFo4p0StQ8SAmOW
// 7CuEMuLCKvGIomC1NS3gVgY9oi5wiRTaU6Xf+VnVMhq5IEqwC0oJS9yuRKPVEE3x
// P1z+rQl1jiFIGpLjZkDLCMDzcZj3zyrVBqFXMKOFnntlw3/2JVKswDUtrc/bI2ES
// CfIuR66YJ+RDPXtZXAaiX9UCAwEAAQ==
// -----END PUBLIC KEY-----`

// 			gomock.InOrder(
// 				jwks.EXPECT().
// 					Certs().
// 					Return(map[string]string{}),

// 				jwks.EXPECT().
// 					RenewCerts(ctx).
// 					Return(nil),

// 				jwks.EXPECT().
// 					Certs().
// 					Return(map[string]string{"test": cert}),
// 			)

// 			ua.Authenticate(ctx)

// 			accountID, _ := ctx.UserValue(ContextKeyAccountID).(string)
// 			fmt.Println(accountID)
// 			assert.Equal(t, "test", accountID)
// 		}))
// }
