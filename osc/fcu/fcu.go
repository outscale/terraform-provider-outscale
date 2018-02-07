package fcu

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/terraform-providers/terraform-provider-outscale/osc"
	"github.com/terraform-providers/terraform-provider-outscale/osc/handler"
)

//FCU the name of the api for url building
const FCU = "fcu"

//Client manages the FCU API
type Client struct {
	client *osc.Client
	VM     VMService
}

// NewFCUClient return a client to operate FCU resources
func NewFCUClient(config osc.Config) (*Client, error) {

	s := &v4.Signer{
		Credentials: credentials.NewStaticCredentials(config.Credentials.AccessKey,
			config.Credentials.SecretKey, ""),
	}

	u, err := url.Parse(fmt.Sprintf(osc.DefaultBaseURL, FCU, config.Credentials.Region))
	if err != nil {
		return nil, err
	}

	config.Target = FCU
	config.BaseURL = u
	config.UserAgent = osc.UserAgent
	config.Client = &http.Client{}

	c := osc.Client{
		Config:                config,
		Signer:                s,
		MarshalHander:         handler.URLEncodeMarshalHander,
		BuildRequestHandler:   handler.BuildURLEncodedRequest,
		UnmarshalHandler:      handler.UnmarshalXML,
		UnmarshalErrorHandler: handler.UnmarshalErrorHandler,
	}

	f := &Client{client: &c,
		VM: VMOperations{client: &c},
	}
	return f, nil
}

// DescribeInstances method
func (c *Client) DescribeInstances(input *DescribeInstancesInput) (*DescribeInstancesOutput, error) {
	inURL := "/"
	endpoint := "DescribeInstances"
	output := &DescribeInstancesOutput{}

	if input == nil {
		input = &DescribeInstancesInput{}
	}

	req, err := c.client.NewRequest(context.TODO(), endpoint, http.MethodGet, inURL, input)

	if err != nil {
		return nil, err
	}

	if req != nil {
		return nil, nil
	}

	return output, nil
}
