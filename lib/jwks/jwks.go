package jwks

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	httpCLient "github.com/ditointernet/go-dito/lib/http"
	"github.com/ditointernet/go-dito/lib/jwks/client"
)

// Client is the structure responsible for handling the JWKS certificates
type Client struct {
	jwksURI              string
	http                 client.HttpClient
	certs                map[string]string
	lastRenewTime        time.Time
	renewMinuteThreshold int
	mux                  sync.Mutex
}

// NewClient constructs a new JWKS instance
func NewClient(jwksURI string, http client.HttpClient, renewMinuteThreshold int) (*Client, error) {
	return &Client{
		jwksURI:              jwksURI,
		http:                 http,
		lastRenewTime:        time.Now(),
		renewMinuteThreshold: renewMinuteThreshold,
	}, nil
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
	if time.Since(c.lastRenewTime).Minutes() > float64(c.renewMinuteThreshold) {
		certs, err := c.fetchCerts(ctx)
		if err != nil {
			return err
		}
		c.certs = certs
		c.lastRenewTime = time.Now()
	}
	c.mux.Unlock()
	return nil
}

// Certs return a list of valid certs
func (c Client) Certs() map[string]string {
	return c.certs
}

type jwk struct {
	KeyID           string   `json:"kid"`
	X509Certificate []string `json:"x5c"`
}

type jwks struct {
	Keys []jwk `json:"keys"`
}

func (c Client) fetchCerts(ctx context.Context) (map[string]string, error) {
	resp, err := c.http.Get(ctx, httpCLient.HTTPRequest{URL: c.jwksURI})
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request was not successful. Received status: %d", resp.StatusCode)
	}

	jwks := jwks{}
	err = json.Unmarshal(resp.Response, &jwks)
	if err != nil {
		return nil, err
	}

	certs := make(map[string]string)
	for _, k := range jwks.Keys {
		certs[k.KeyID] = "-----BEGIN CERTIFICATE-----\n" + k.X509Certificate[0] + "\n-----END CERTIFICATE-----"
	}

	return certs, nil
}
