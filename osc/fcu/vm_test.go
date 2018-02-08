package fcu

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/terraform-providers/terraform-provider-outscale/osc"
)

var (
	mux *http.ServeMux

	ctx = context.TODO()

	client *Client

	server *httptest.Server
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	url, _ := url.Parse(server.URL)

	config := osc.Config{
		BaseURL: url,
		Credentials: &osc.Credentials{
			AccessKey: "AKID",
			SecretKey: "SecretKey",
			Region:    "region",
		},
	}

	client, _ = NewFCUClient(config)

}

func teardown() {
	server.Close()
}
