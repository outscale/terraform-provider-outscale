package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const mediaTypeURLEncoded = "application/x-www-form-urlencoded"

// BuildURLEncodedRequest the request with a body, if it's post then adds it to the body of the request,
// otherwise adds it to the url query
func BuildURLEncodedRequest(body interface{}, method, url string) (*http.Request, io.ReadSeeker, error) {

	isLBU := strings.Contains(url, "lbu")

	if method == http.MethodPost && !isLBU {
		reader := strings.NewReader(body.(string))
		req, err := http.NewRequest(method, url, reader)
		if err != nil {
			return nil, nil, err
		}
		return req, reader, nil
	}

	if isLBU {
		value := body.(string)
		i := strings.Index(value, "&")
		substring := value[7:i]

		fmt.Println("substring =>", substring)

		v := struct {
			Action  string `json:"Action"`
			Version string `json:"Version"`
		}{substring, "2017-12-15"}

		ja, _ := json.Marshal(v)
		b := string(ja)

		fmt.Println("READER =>", b)

		reader := strings.NewReader(b)
		req, err := http.NewRequest(method, url, reader)
		if err != nil {
			return nil, nil, err
		}
		req.URL.RawQuery = body.(string)
		return req, reader, nil
	}

	if method == http.MethodGet {
		req, err := http.NewRequest(method, url, nil)
		if err != nil {
			return nil, nil, err
		}

		req.URL.RawQuery = body.(string)
		return req, nil, nil

	}
	return nil, nil, fmt.Errorf("Method %s not supported", method)
}
