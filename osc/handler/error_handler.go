package handler

import (
	"encoding/json"
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

//JSONICUError ...
type JSONICUError struct {
	Msj     string `json:"message"`
	Type    string `json:"__type"`
	Message string `json:"Message"`
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
	return fmt.Errorf("%s: %s", v.Errors[0].Code, v.Errors[0].Message)
}

// UnmarshalJSONErrorHandler for HTTP Response
func UnmarshalJSONErrorHandler(r *http.Response) error {
	defer r.Body.Close()
	v := JSONICUError{}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("Read body: %v", err)
	}

	err = json.Unmarshal(data, &v)
	if err != nil {
		return fmt.Errorf("error unmarshalling response %v", err)
	}

	// Response body format is not consistent between metadata endpoints.
	// Grab the error message as a string and include that as the source error
	return fmt.Errorf("%s: %s", v.Type, v.Message)
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
