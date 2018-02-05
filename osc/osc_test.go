// Package osc ...
package osc

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
)

func buildRequest(service, region, body string) (*http.Request, io.ReadSeeker) {
	reader := strings.NewReader(body)
	endpoint := fmt.Sprintf("https://%s.%s.outscale.com", service, region)
	req, _ := http.NewRequest("POST", endpoint, reader)
	req.URL.Opaque = "//example.org/bucket/key-._~,!@#$%^&*()"
	req.Header.Add("X-Amz-Target", "prefix.Operation")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", string(len(body)))
	return req, reader
}

func buildSigner() *v4.Signer {
	return &v4.Signer{
		Credentials: credentials.NewStaticCredentials("AKID", "SECRET", ""),
	}
}

func buildClient() *Client {

	return &Client{
		signer: buildSigner(),
		Config: Config{
			Credentials: &Credentials{
				Region: "eu-west-1",
			},
		},
	}
}

func TestSignRequest(t *testing.T) {
	req, body := buildRequest("fcu", "eu-west-1", "{}")
	client := buildClient()
	cor, _ := client.Sign(req, body, time.Unix(0, 0), "fcu")
	fmt.Println(cor)

	expectedDate := "19700101T000000Z"
	expectedSig := "AWS4-HMAC-SHA256 Credential=AKID/19700101/eu-west-1/fcu/aws4_request, SignedHeaders=content-length;content-type;host;x-amz-date;x-amz-target, Signature=70b1a1fc5246f4f9444daa40d7530264e7255b4e966d73a5bb4747f339e862b2"

	q := req.Header

	if e, a := expectedDate, q.Get("X-Amz-Date"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

	if e, a := expectedSig, q.Get("Authorization"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

}
