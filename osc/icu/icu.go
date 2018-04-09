package icu

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/terraform-providers/terraform-provider-outscale/osc"
	"github.com/terraform-providers/terraform-provider-outscale/osc/handler"
)

//FCU the name of the api for url building
const ICU = "icu"

//Client manages the FCU API
type Client struct {
	client *osc.Client
	API    ICUService
}

// NewFCUClient return a client to operate FCU resources
func NewICUClient(config osc.Config) (*Client, error) {

	s := &v4.Signer{
		Credentials: credentials.NewStaticCredentials(config.Credentials.AccessKey,
			config.Credentials.SecretKey, ""),
	}

	u, err := url.Parse(fmt.Sprintf(osc.DefaultBaseURL, ICU, config.Credentials.Region))
	if err != nil {
		return nil, err
	}

	config.Target = ICU
	config.BaseURL = u
	config.UserAgent = osc.UserAgent
	config.Client = &http.Client{}

	fmt.Printf("\n\n[DEBUG] CONFIG => %v\n\n", config)

	c := osc.Client{
		Config:                config,
		Signer:                s,
		MarshalHander:         handler.URLEncodeMarshalHander,
		BuildRequestHandler:   handler.BuildURLEncodedRequest,
		UnmarshalHandler:      handler.UnmarshalXML,
		UnmarshalErrorHandler: handler.UnmarshalErrorHandler,
	}

	f := &Client{client: &c,
		API: ICUOperations{client: &c},
	}
	return f, nil
}
