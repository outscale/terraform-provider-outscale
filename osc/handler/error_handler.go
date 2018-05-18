package handler

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

// ErrMsg Error lists available
var ErrMsg = map[string]string{
	"SerializationError": "unable to unmarshal EC2 metadata error respose",
	"HTTP":               "HTTP Error",
}

// Error ...
type Error struct {
	Code    string `xml:"Code"`
	Message string `xml:"Message"`
}

// XMLError ...
type XMLError struct {
	XMLName   xml.Name `xml:"Response"`
	Errors    []Error  `xml:"Errors>Error"`
	RequestID string   `xml:"RequestID"`
}

// XMLLBUError ...
type XMLLBUError struct {
	XMLName   xml.Name `xml:"ErrorResponse"`
	Errors    Error    `xml:"Error"`
	RequestID string   `xml:"RequestID"`
}

// UnmarshalErrorHandler for HTTP Response
func UnmarshalErrorHandler(r *http.Response) error {
	defer r.Body.Close()
	v := XMLError{}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("Read body: %v", err)
	}

	err = xml.Unmarshal(data, &v)
	if err != nil {
		return fmt.Errorf("error unmarshalling response %v", err)
	}

	// Response body format is not consistent between metadata endpoints.
	// Grab the error message as a string and include that as the source error
	return SendError(v)

}

// SendError method which receives the message and the error
func SendError(msg XMLError) error {
	return fmt.Errorf("%s: %s", msg.Errors[0].Code, msg.Errors[0].Message)
}

// UnmarshalLBUErrorHandler for HTTP Response
func UnmarshalLBUErrorHandler(r *http.Response) error {
	defer r.Body.Close()
	v := XMLLBUError{}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("Read body: %v", err)
	}

	err = xml.Unmarshal(data, &v)
	if err != nil {
		return fmt.Errorf("error unmarshalling response %v", err)
	}

	return fmt.Errorf("%s: %s", v.Errors.Code, v.Errors.Message)

}
