//Outscale Client
package osc

import (
	"net/http"
	"net/url"
)

const (
	libraryVersion = "1.0"
	defaultBaseURL = "https://%s.%s.outscale.com"
	userAgent      = "osc/" + libraryVersion
	mediaTypeJSON  = "application/json"
	mediaTypeWSDL  = "application/wsdl+xml"
)

// Client manages the communication between the Outscale API's
type Client struct {
	Config Config
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
