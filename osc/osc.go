// Package osc ...
package osc

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go/aws/signer/v4"
)

const (
	libraryVersion      = "1.0"
	defaultBaseURL      = "https://%s.%s.outscale.com"
	opaqueBaseURL       = "/%s.%s.outscale.com/%s"
	userAgent           = "osc/" + libraryVersion
	mediaTypeJSON       = "application/json"
	mediaTypeWSDL       = "application/wsdl+xml"
	mediaTypeURLEncoded = "application/x-www-form-urlencoded"
	signatureVersion    = "4"
)

// BuildRequestHandler creates a new request and marshals the body depending on the implementation
type BuildRequestHandler func(v interface{}, method, url string) (*http.Request, io.ReadSeeker, error)

// MarshalHander marshals the incoming body to a desired format
type MarshalHander func(v interface{}, action, version string) (string, error)

// UnmarshalHandler unmarshals the body request depending on different implementations
type UnmarshalHandler func(v interface{}, req *http.Response) error

// UnmarshalErrorHandler unmarshals the errors coming from an http respose
type UnmarshalErrorHandler func(r *http.Response) error

// Client manages the communication between the Outscale API's
type Client struct {
	Config Config
	signer *v4.Signer

	// Handlers
	MarshalHander         MarshalHander
	BuildRequestHandler   BuildRequestHandler
	UnmarshalHandler      UnmarshalHandler
	UnmarshalErrorHandler UnmarshalErrorHandler
}

// Config Configuration of the client
type Config struct {
	Endpoint    string
	Target      string
	Credentials *Credentials

	// HTTP client used to communicate with the Outscale API.
	client *http.Client

	// Base URL for API requests.
	BaseURL *url.URL

	// User agent for client
	UserAgent string

	// Services used for communicating with the API
	// To be implemented

	// Optional function called after every successful request made to the DO APIs
	onRequestCompleted RequestCompletionCallback
}

// Credentials needed access key, secret key and region
type Credentials struct {
	AccessKey string
	SecretKey string
	Region    string
}

// RequestCompletionCallback defines the type of the request callback function.
type RequestCompletionCallback func(*http.Request, *http.Response)

// Sign HTTP Request for authentication
func (c Client) Sign(req *http.Request, body io.ReadSeeker, timestamp time.Time, service string) (http.Header, error) {
	return c.signer.Sign(req, body, c.Config.Target, c.Config.Credentials.Region, timestamp)
}

// NewRequest creates a request and signs it
func (c *Client) NewRequest(ctx context.Context, operation, method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	b, err := c.MarshalHander(body, operation, "")
	if err != nil {
		return nil, err
	}

	u := c.Config.BaseURL.ResolveReference(rel)

	req, reader, err := c.BuildRequestHandler(b, method, u.String())
	if err != nil {
		return nil, err
	}

	fmt.Println(rel.Opaque)

	_, err = c.Sign(req, reader, time.Now(), c.Config.Target)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// SetHeaders sets the headers for the request
func (c Client) SetHeaders(req *http.Request, target, operation string) {
	req.Header.Add("User-Agent", c.Config.UserAgent)
	req.Header.Add("X-Amz-Target", fmt.Sprintf("%s.%s", target, operation))
}

// Do sends the request to the API's
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) error {

	req = req.WithContext(ctx)

	resp, err := c.Config.client.Do(req)
	if err != nil {
		return err
	}

	err = c.checkResponse(resp)
	if err != nil {
		return err
	}

	return c.UnmarshalHandler(v, resp)
}

func (c Client) checkResponse(r *http.Response) error {
	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil
	}

	return c.UnmarshalErrorHandler(r)

}
