// Package osc ...
package osc

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/signer/v4"
)

const (
	libraryVersion = "1.0"
	// DefaultBaseURL ...
	DefaultBaseURL = "https://%s.%s.outscale.com"
	opaqueBaseURL  = "/%s.%s.outscale.com/%s"
	// UserAgent ...
	UserAgent           = "osc/" + libraryVersion
	mediaTypeJSON       = "application/json"
	mediaTypeWSDL       = "application/wsdl+xml"
	mediaTypeURLEncoded = "application/x-www-form-urlencoded"
	signatureVersion    = "4"
)

// BuildRequestHandler creates a new request and marshals the body depending on the implementation
type BuildRequestHandler func(v interface{}, method, url string) (*http.Request, io.ReadSeeker, error)

// MarshalHander marshals the incoming body to a desired format
type MarshalHander func(v interface{}, action, version string, isLBU bool) (string, error)

// UnmarshalHandler unmarshals the body request depending on different implementations
type UnmarshalHandler func(v interface{}, req *http.Response) error

// UnmarshalLBUXML ...
type UnmarshalLBUXML func(v interface{}, req *http.Response, operation string) error

// UnmarshalErrorHandler unmarshals the errors coming from an http respose
type UnmarshalErrorHandler func(r *http.Response) error

// Client manages the communication between the Outscale API's
type Client struct {
	Config Config
	Signer *v4.Signer

	// Handlers
	MarshalHander         MarshalHander
	BuildRequestHandler   BuildRequestHandler
	UnmarshalHandler      UnmarshalHandler
	UnmarshalLBUXML       UnmarshalLBUXML
	UnmarshalErrorHandler UnmarshalErrorHandler
}

// Config Configuration of the client
type Config struct {
	Target      string
	Credentials *Credentials

	// HTTP client used to communicate with the Outscale API.
	Client *http.Client

	// Base URL for API requests.
	BaseURL *url.URL

	// User agent for client
	UserAgent string

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
	return c.Signer.Sign(req, body, c.Config.Target, c.Config.Credentials.Region, timestamp)
}

// NewRequest creates a request and signs it
func (c *Client) NewRequest(ctx context.Context, operation, method, urlStr string, body interface{}) (*http.Request, error) {
	rel, errp := url.Parse(urlStr)
	if errp != nil {
		return nil, errp
	}

	var b interface{}
	var err error

	isLBU := (strings.Contains(operation, "LoadBalancer") || strings.Contains(operation, "ConfigureHealthCheck"))

	// method for FCU API
	if method != http.MethodPost {
		b, err = c.MarshalHander(body, operation, "2018-05-14", !isLBU)
		if err != nil {
			return nil, err
		}
	} else if method == http.MethodPost {
		v := struct {
			Action               string `json:"Action"`
			Version              string `json:"Version"`
			AuthenticationMethod string `json:"AuthenticationMethod"`
		}{operation, "2017-12-15", "accesskey"}

		var m map[string]string

		ja, _ := json.Marshal(v)
		json.Unmarshal(ja, &m)
		jb, _ := json.Marshal(body)
		json.Unmarshal(jb, &m)

		jm, _ := json.Marshal(m)

		b = string(jm)
	}

	u := c.Config.BaseURL.ResolveReference(rel)

	req, reader, err := c.BuildRequestHandler(b, method, u.String())
	if err != nil {
		return nil, err
	}

	fmt.Println(rel.Opaque)

	if strings.Contains(operation, "LoadBalancer") || strings.Contains(operation, "ConfigureHealthCheck") {
		c.SetHeaders(req, "lbu_20180514", operation)
	} else if strings.Contains(operation, "AccessKey") {
		c.SetHeaders(req, "TinaIcuService", operation)
	}

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
	req.Header.Add("Content-Type", mediaTypeURLEncoded)
}

// Do sends the request to the API's
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) error {

	req = req.WithContext(ctx)

	resp, err := c.Config.Client.Do(req)
	if err != nil {
		return err
	}

	err = c.checkResponse(resp)
	if err != nil {
		return err
	}

	if req.Method == "POST" {
		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(v)
		if err != nil {
			return err
		}

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
