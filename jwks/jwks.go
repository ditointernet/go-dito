package jwks

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/ditointernet/go-dito/errors"
	client "github.com/ditointernet/go-dito/http"
)

// Client is the structure responsible for handling JWKS certificates
type Client struct {
	jwksURI              string
	http                 HTTPClient
	certs                map[string]string
	lastRenewTime        time.Time
	renewMinuteThreshold int
	mux                  sync.Mutex
}

// NewClient constructs a new JWKS instance
func NewClient(jwksURI string, http HTTPClient, renewMinuteThreshold int) (*Client, error) {
	if jwksURI == "" {
		return nil, errors.NewMissingRequiredDependency("jwksURI")
	}

	if http == nil {
		return nil, errors.NewMissingRequiredDependency("http")
	}

	if renewMinuteThreshold < 0 {
		renewMinuteThreshold = 5
	}

	return &Client{
		jwksURI:              jwksURI,
		http:                 http,
		renewMinuteThreshold: renewMinuteThreshold,
	}, nil
}

// MustNewClient constructs a new JWKS instance.
// It panics if any error is found.
func MustNewClient(jwksURI string, http HTTPClient, renewMinuteThreshold int) *Client {
	cli, err := NewClient(jwksURI, http, renewMinuteThreshold)
	if err != nil {
		panic(err)
	}

	return cli
}

// GetCerts makes a http request to jwksURI and retrieves a list of valid certtificates
func (c *Client) GetCerts(ctx context.Context) error {
	certs, err := c.fetchCerts(ctx)
	if err != nil {
		return err
	}

	c.certs = certs

	return nil
}

// RenewCerts compare if the certs are valid if not retrieves a new list of valid certificates
func (c *Client) RenewCerts(ctx context.Context) error {
	c.mux.Lock()
	defer c.mux.Unlock()
	if time.Since(c.lastRenewTime).Minutes() > float64(c.renewMinuteThreshold) {
		certs, err := c.fetchCerts(ctx)
		if err != nil {
			return err
		}
		c.certs = certs
	}
	return nil
}

// Certs return a list of valid certs
func (c *Client) Certs() map[string]string {
	return c.certs
}

type jwk struct {
	KeyID           string   `json:"kid"`
	X509Certificate []string `json:"x5c"`
}

type jwks struct {
	Keys []jwk `json:"keys"`
}

func (c *Client) fetchCerts(ctx context.Context) (map[string]string, error) {
	resp, err := c.http.Get(ctx, client.HTTPRequest{URL: c.jwksURI})
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("request was not successful. Received status: %d", resp.StatusCode).WithKind(errors.KindInternal)
	}

	jwks := jwks{}
	err = json.Unmarshal(resp.Response, &jwks)
	if err != nil {
		return nil, err
	}

	certs := make(map[string]string)
	c.lastRenewTime = time.Now()

	for _, k := range jwks.Keys {
		certs[k.KeyID] = "-----BEGIN CERTIFICATE-----\n" + k.X509Certificate[0] + "\n-----END CERTIFICATE-----"
	}

	return certs, nil
}
