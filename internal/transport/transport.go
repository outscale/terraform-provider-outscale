package transport

import (
	"crypto/sha256"
	"io"
	"net/http"

	"github.com/aws/smithy-go/aws-http-auth/credentials"
	"github.com/aws/smithy-go/aws-http-auth/sigv4"
)

type transport struct {
	inner  http.RoundTripper
	signer *sigv4.Signer
	*securityProviderAWSv4
}

func NewSecurityProviderAWSv4(
	accessKey, secretKey, sessionToken, service, region string,
) *securityProviderAWSv4 {
	return &securityProviderAWSv4{
		accessKey:    accessKey,
		secretKey:    secretKey,
		sessionToken: sessionToken,
		service:      service,
		region:       region,
	}
}

type securityProviderAWSv4 struct {
	accessKey    string
	secretKey    string
	sessionToken string
	service      string
	region       string
}

func (t *transport) sign(req *http.Request) error {
	sigReq := sigv4.SignRequestInput{
		Request: req,
		Credentials: credentials.Credentials{
			AccessKeyID:     t.accessKey,
			SecretAccessKey: t.secretKey,
			SessionToken:    t.sessionToken,
		},
		Service: t.service,
		Region:  t.region,
	}
	if req.GetBody != nil {
		h := sha256.New()

		bodyreader, err := req.GetBody()
		if err != nil {
			return err
		}

		if _, err := io.Copy(h, bodyreader); err != nil {
			return err
		}

		sigReq.PayloadHash = h.Sum(nil)
	}

	err := t.signer.SignRequest(&sigReq)
	if err != nil {
		return err
	}

	return nil
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if err := t.sign(req); err != nil {
		return nil, err
	}

	resp, err := t.inner.RoundTrip(req)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func NewTransport(t http.RoundTripper, provider *securityProviderAWSv4) *transport {
	return &transport{t, sigv4.New(), provider}
}
