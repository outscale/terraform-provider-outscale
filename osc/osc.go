// Package osc ...
package osc

import (
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go/aws/signer/v4"
)

const (
	libraryVersion   = "1.0"
	defaultBaseURL   = "https://%s.%s.outscale.com"
	userAgent        = "osc/" + libraryVersion
	mediaTypeJSON    = "application/json"
	mediaTypeWSDL    = "application/wsdl+xml"
	signatureVersion = "4&SnapshotId"
)

// Client manages the communication between the Outscale API's
type Client struct {
	Config Config
	signer *v4.Signer
}

// Config Configuration of the client
type Config struct {
	Endpoint    string
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

// Sign ...
func (c Client) Sign(req *http.Request, body io.ReadSeeker, timestamp time.Time, service string) (http.Header, error) {
	return c.signer.Sign(req, body, service, c.Config.Credentials.Region, timestamp)
}
