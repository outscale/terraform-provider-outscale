package osc

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// ErrMsg Error lists available
var ErrMsg = map[string]string{
	"SerializationError": "unable to unmarshal EC2 metadata error respose",
	"HTTP":               "HTTP Error",
}

// UnmarshalErrorHandler for HTTP Response
func UnmarshalErrorHandler(r *http.Response) error {
	defer r.Body.Close()
	b := &bytes.Buffer{}
	if _, err := io.Copy(b, r.Body); err != nil {
		return SendError(ErrMsg["SerializationError"], err)

	}

	// Response body format is not consistent between metadata endpoints.
	// Grab the error message as a string and include that as the source error
	return SendError(ErrMsg["HTTP"], errors.New(b.String()))

}

// SendError method which receives the message and the error
func SendError(msg string, err error) error {
	return errors.New(msg + " - " + fmt.Sprint(err))
}
