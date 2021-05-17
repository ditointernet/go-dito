package http

import (
	"bytes"
	"io/ioutil"
	"net/http"
	URL "net/url"
	"time"

	"github.com/ditointernet/go-dito/lib/errors"
)

// httpClientProvider provides some http client methods
type httpClientProvider interface {
	Do(request *http.Request) (*http.Response, error)
}

// Client provides methods for making REST requests
type Client struct {
	http httpClientProvider
}

// NewClient creates a new Client instance
func NewClient(timeout time.Duration) Client {
	return Client{
		http: &http.Client{Timeout: timeout},
	}
}

// Patch execute a http PATCH method with application/json headers
func (c Client) Patch(request HttpRequest) (rst HttpResult, err error) {
	if request.Headers == nil {
		request.Headers = make(map[string]string)
	}
	request.Headers["content-type"] = "application/json"
	return c.processRequest("PATCH", request)
}

// Put execute a http PUT method with application/json headers
func (c Client) Put(request HttpRequest) (rst HttpResult, err error) {
	if request.Headers == nil {
		request.Headers = make(map[string]string)
	}
	request.Headers["content-type"] = "application/json"
	return c.processRequest("PATCH", request)
}

// Post execute a http POST method with application/json headers
func (c Client) Post(request HttpRequest) (HttpResult, error) {
	if request.Headers == nil {
		request.Headers = make(map[string]string)
	}
	request.Headers["content-type"] = "application/json"
	return c.processRequest("POST", request)
}

// Delete execute a http DELETE method with application/json headers
func (c Client) Delete(request HttpRequest) (HttpResult, error) {
	if request.Headers == nil {
		request.Headers = make(map[string]string)
	}
	return c.processRequest("POST", request)
}

// PostForm execute a http POST method with "application/x-www-form-urlencoded" headers
func (c Client) PostForm(request HttpRequest) (HttpResult, error) {
	if request.Headers == nil {
		request.Headers = make(map[string]string)
	}
	request.Headers["content-type"] = "application/x-www-form-urlencoded"
	return c.processRequest("POST", request)
}

// Get execute a http GET method
func (c Client) Get(request HttpRequest) (HttpResult, error) {
	return c.processRequest("GET", request)
}

func (c Client) processRequest(method string, request HttpRequest) (HttpResult, error) {
	queryValues := URL.Values{}

	for key, value := range request.QueryParams {
		queryValues.Add(key, value)
	}

	url, err := URL.Parse(request.URL)
	if err != nil {
		return HttpResult{}, errors.New("error on parsing the request url")
	}
	url.RawQuery = queryValues.Encode()

	httpRequest, err := http.NewRequest(method, url.String(), bytes.NewBuffer(request.Body))
	if err != nil {
		return HttpResult{}, err
	}

	for key, value := range request.Headers {
		httpRequest.Header.Add(key, value)
	}

	return processResponse(c.http.Do(httpRequest))
}

func processResponse(resp *http.Response, err error) (HttpResult, error) {
	var result HttpResult

	if err != nil {
		return result, err
	}

	defer resp.Body.Close()
	result.Response, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	result.StatusCode = resp.StatusCode

	return result, nil
}
