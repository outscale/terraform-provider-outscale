package transport

import (
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
)

type transport struct {
	transport http.RoundTripper
	signer    *v4.Signer
	region    string
}

func (t *transport) sign(req *http.Request, body []byte) error {
	reader := strings.NewReader(string(body))
	timestamp := time.Now()
	_, err := t.signer.Sign(req, reader, "oapi", t.region, timestamp)
	return err
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	if err := t.sign(req, body); err != nil {
		return nil, err
	}

	resp, err := t.transport.RoundTrip(req)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func NewTransport(accessKey, accessSecret, region string, t http.RoundTripper) *transport {
	s := &v4.Signer{
		Credentials: credentials.NewStaticCredentials(accessKey,
			accessSecret, ""),
	}
	return &transport{t, s, region}
}
