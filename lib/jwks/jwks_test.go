package jwks

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/ditointernet/go-dito/lib/jwks/client"
	"github.com/ditointernet/go-dito/lib/jwks/mocks"
	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
)

func TestClient_Certs(t *testing.T) {
	withMock := func(runner func(t *testing.T, c Client)) func(t *testing.T) {
		return func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			runner(t, Client{
				jwksURI: "jwksURI",
				http:    mocks.NewMockHttpClient(ctrl),
				certs:   map[string]string{"1": "-----BEGIN CERTIFICATE-----\ntest\n-----END CERTIFICATE-----"},
			})
		}
	}

	t.Run("should return the certificates when the method is called",
		withMock(func(t *testing.T, c Client) {
			assert.Equal(t, c.Certs(), c.certs)
		}))
}

func TestClient_RenewCerts(t *testing.T) {
	var httpMock *mocks.MockHttpClient

	withMock := func(runner func(t *testing.T, c Client)) func(t *testing.T) {
		return func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			httpMock = mocks.NewMockHttpClient(ctrl)
			runner(t, Client{
				jwksURI:       "jwksURI",
				http:          httpMock,
				certs:         map[string]string{"1": "-----BEGIN CERTIFICATE-----\nt\n-----END CERTIFICATE-----"},
				lastRenewTime: time.Now(),
				mux:           sync.Mutex{},
			})
		}
	}

	t.Run("should fail when the certificates fetch return error",
		withMock(func(t *testing.T, c Client) {
			httpMock.EXPECT().Get(context.TODO(), gomock.Any()).Return(client.HttpResult{}, errors.New("error"))

			assert.EqualError(t, c.RenewCerts(context.TODO()), "error")
		}))

	t.Run("should update the certificates when the renew process is successful",
		withMock(func(t *testing.T, c Client) {
			httpMock.EXPECT().Get(context.TODO(), gomock.Any()).Return(client.HttpResult{
				StatusCode: 200,
				Response:   getBody(),
			}, nil)

			c.RenewCerts(context.TODO())

			expected := map[string]string{
				"1": "-----BEGIN CERTIFICATE-----\ntest\n-----END CERTIFICATE-----",
				"2": "-----BEGIN CERTIFICATE-----\ntest2\n-----END CERTIFICATE-----",
			}

			assert.Equal(t, expected, c.certs)
		}))

	t.Run("should not update the certificates when the renew process is called in the renew threshold",
		withMock(func(t *testing.T, c Client) {
			httpMock.EXPECT().Get(context.TODO(), gomock.Any()).Do(func(request client.HttpRequest) (client.HttpResult, error) {
				t.FailNow()
				return client.HttpResult{}, nil
			}).AnyTimes()

			c.renewMinuteThreshold = 5

			go c.RenewCerts(context.TODO())
			go c.RenewCerts(context.TODO())
		}))
}

func TestClient_NewClient(t *testing.T) {
	var httpMock *mocks.MockHttpClient

	withMock := func(runner func(t *testing.T)) func(t *testing.T) {
		return func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			httpMock = mocks.NewMockHttpClient(ctrl)
			runner(t)
		}
	}

	t.Run("should fail when the certificates fetch don't receives a Ok http status",
		withMock(func(t *testing.T) {
			httpMock.EXPECT().Get(context.TODO(), gomock.Any()).Return(client.HttpResult{StatusCode: 500}, nil)

			clientJwks, _ := NewClient("jwksURI", httpMock, 5)
			err := clientJwks.GetCerts(context.TODO())

			assert.EqualError(t, err, "request was not successful. Received status: 500")
		}))

	t.Run("should fail when the received certificates couldn't be parsed",
		withMock(func(t *testing.T) {
			httpMock.EXPECT().Get(context.TODO(), gomock.Any()).Return(client.HttpResult{
				StatusCode: 200,
				Response:   []byte("0"),
			}, nil)

			clientJwks, _ := NewClient("jwksURI", httpMock, 5)

			err := clientJwks.GetCerts(context.TODO())

			assert.EqualError(t, err, "json: cannot unmarshal number into Go value of type jwks.jwks")
		}))

	t.Run("should return the certificates when they are received and parsed with success",
		withMock(func(t *testing.T) {
			httpMock.EXPECT().Get(context.TODO(), gomock.Any()).Return(client.HttpResult{
				StatusCode: 200,
				Response:   getBody(),
			}, nil)

			clientJwks, _ := NewClient("jwksURI", httpMock, 5)

			clientJwks.GetCerts(context.TODO())

			expected := map[string]string{
				"1": "-----BEGIN CERTIFICATE-----\ntest\n-----END CERTIFICATE-----",
				"2": "-----BEGIN CERTIFICATE-----\ntest2\n-----END CERTIFICATE-----",
			}

			assert.Equal(t, expected, clientJwks.certs)
		}))
}

func getBody() []byte {
	body, _ := json.Marshal(jwks{Keys: []jwk{
		{
			KeyID:           "1",
			X509Certificate: []string{"test"},
		},
		{
			KeyID:           "2",
			X509Certificate: []string{"test2"},
		},
	}})

	return body
}
