package handler

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

const mediaTypeURLEncoded = "application/x-www-form-urlencoded"

// BuildURLEncodedRequest the request with a body, if it's post then adds it to the body of the request,
// otherwise adds it to the url query
func BuildURLEncodedRequest(body interface{}, method, url string) (*http.Request, io.ReadSeeker, error) {

	if method == http.MethodPost {
		reader := strings.NewReader(body.(string))
		req, err := http.NewRequest(method, url, reader)
		if err != nil {
			return nil, nil, err
		}
		return req, reader, nil
	}

	if method == http.MethodGet {
		req, err := http.NewRequest(method, url, nil)
		if err != nil {
			return nil, nil, err
		}

		req.URL.RawQuery = body.(string)

		if strings.Contains(url, "lbu") {
			replacer := strings.NewReplacer("Member", "member")
			output := replacer.Replace(body.(string))

			req.URL.RawQuery = output
		}

		return req, nil, nil

	}
	return nil, nil, fmt.Errorf("Method %s not supported", method)
}
