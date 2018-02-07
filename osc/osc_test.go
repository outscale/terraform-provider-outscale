// Package osc ...
package osc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
)

var (
	mux *http.ServeMux

	ctx = context.TODO()

	client *Client

	server *httptest.Server
)

func init() {
	client = buildClient()
}

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	client = buildClient()
	url, _ := url.Parse(server.URL)
	client.Config.BaseURL = url
}

func teardown() {
	server.Close()
}

func buildRequest(service, region, body string) (*http.Request, io.ReadSeeker) {
	reader := strings.NewReader(body)
	endpoint := fmt.Sprintf("https://%s.%s.outscale.com", service, region)
	req, _ := http.NewRequest("POST", endpoint, reader)
	req.URL.Opaque = "//example.org/bucket/key-._~,!@#$%^&*()"
	req.Header.Add("X-Amz-Target", "prefix.Operation")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", fmt.Sprintf("%d", len(body)))
	return req, reader
}

func buildSigner() *v4.Signer {
	return &v4.Signer{
		Credentials: credentials.NewStaticCredentials("AKID", "SECRET", ""),
	}
}

func buildClient() *Client {

	baseURL, _ := url.Parse(fmt.Sprintf(DefaultBaseURL, "fcu", "eu-west-2"))
	fmt.Println(baseURL.Opaque)

	return &Client{
		Signer:                buildSigner(),
		BuildRequestHandler:   buildTestHandler,
		MarshalHander:         testBuildRequestHandler,
		UnmarshalHandler:      unmarshalTestHandler,
		UnmarshalErrorHandler: testUnmarshalErrorHandler,
		Config: Config{
			UserAgent: "test",
			Target:    "fcu",
			BaseURL:   baseURL,
			Client:    &http.Client{},
			Credentials: &Credentials{
				Region: "eu-west-1",
			},
		},
	}
}

func buildTestHandler(v interface{}, method, url string) (*http.Request, io.ReadSeeker, error) {
	reader := strings.NewReader("{}")
	req, _ := http.NewRequest(method, url, reader)

	req.Header.Add("Content-Type", mediaTypeURLEncoded)

	return req, reader, nil
}

func testBuildRequestHandler(v interface{}, action, version string) (string, error) {
	return "{}", nil
}

func unmarshalTestHandler(v interface{}, req *http.Response) error {
	return nil
}

func testUnmarshalErrorHandler(r *http.Response) error {
	return errors.New("This is an error")
}

func TestSign(t *testing.T) {
	req, body := buildRequest("fcu", "eu-west-1", "{}")
	client.Sign(req, body, time.Unix(0, 0), "fcu")

	expectedDate := "19700101T000000Z"
	expectedSig := "AWS4-HMAC-SHA256 Credential=AKID/19700101/eu-west-1/fcu/aws4_request, SignedHeaders=content-length;content-type;host;x-amz-date;x-amz-target, Signature=3e6c29372bb5d7c5ce2c605a29bd774b1be3a8d794ea31e033c50af2c5e27302"

	q := req.Header

	if e, a := expectedDate, q.Get("X-Amz-Date"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

	if e, a := expectedSig, q.Get("Authorization"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

}

func TestSetHeaders(t *testing.T) {
	c := buildClient()

	req, _ := http.NewRequest(http.MethodGet, "http//:example.org/", nil)
	c.SetHeaders(req, "fcu", "DescribeInstances")

	q := req.Header
	targetExpected := "fcu.DescribeInstances"
	agentExpected := "test"

	if e, a := agentExpected, q.Get("User-Agent"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := targetExpected, q.Get("X-Amz-Target"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
}

func TestNewRequest(t *testing.T) {
	c := buildClient()

	inURL, outURL := "foo", fmt.Sprintf(DefaultBaseURL+"/foo", "fcu", "eu-west-2")
	inBody, outBody := "{}", "{}"
	req, _ := c.NewRequest(context.TODO(), "operation", http.MethodGet, inURL, inBody)
	fmt.Println(req.URL.Opaque)

	// test relative URL was expanded
	if req.URL.String() != outURL {
		t.Errorf("NewRequest(%v) URL = %v, expected %v", inURL, req.URL, outURL)
	}

	// test body was JSON encoded
	body, _ := ioutil.ReadAll(req.Body)
	if string(body) != outBody {
		t.Errorf("NewRequest(%v)Body = %v, expected %v", inBody, string(body), outBody)
	}
}

func TestDo(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if m := http.MethodGet; m != r.Method {
			t.Errorf("Request method = %v, expected %v", r.Method, m)
		}
		fmt.Fprint(w, `{}`)
	})

	inURL := "/"
	inBody := "{}"

	req, _ := client.NewRequest(context.TODO(), "operation", http.MethodGet, inURL, inBody)
	err := client.Do(context.Background(), req, nil)
	if err != nil {
		t.Fatalf("Do(): %v", err)
	}
}

func TestDo_ErrorResponse(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if m := http.MethodGet; m != r.Method {
			t.Errorf("Request method = %v, expected %v", r.Method, m)
		}
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, `{}`)
	})

	inURL := "/"
	inBody := "{}"

	req, _ := client.NewRequest(context.TODO(), "operation", http.MethodGet, inURL, inBody)
	err := client.Do(context.Background(), req, nil)
	if err == nil {
		t.Fatalf("Do(): %v", err)
	}
}
