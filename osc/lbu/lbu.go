package lbu

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/terraform-providers/terraform-provider-outscale/osc"
	"github.com/terraform-providers/terraform-provider-outscale/osc/handler"
)

//LBU the name of the api for url building
const LBU = "lbu"

//Client manages the FCU API
type Client struct {
	client *osc.Client
	API    Service
}

// NewLBUClient return a client to operate FCU resources
func NewLBUClient(config osc.Config) (*Client, error) {

	s := &v4.Signer{
		Credentials: credentials.NewStaticCredentials(config.Credentials.AccessKey,
			config.Credentials.SecretKey, ""),
	}

	u, err := url.Parse(fmt.Sprintf(osc.DefaultBaseURL, LBU, config.Credentials.Region))
	if err != nil {
		return nil, err
	}

	config.Target = LBU
	config.BaseURL = u
	config.UserAgent = osc.UserAgent
	config.Client = &http.Client{}

	c := osc.Client{
		Config:                config,
		Signer:                s,
		MarshalHander:         handler.URLLBUEncodeMarshalHander,
		BuildRequestHandler:   handler.BuildURLEncodedRequest,
		UnmarshalHandler:      handler.UnmarshalLBUXML,
		UnmarshalErrorHandler: handler.UnmarshalLBUErrorHandler,
	}

	f := &Client{client: &c,
		API: Operations{client: &c},
	}
	return f, nil
}
