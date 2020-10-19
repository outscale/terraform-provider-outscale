package outscale

import (
	"crypto/tls"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/logging"

	oscgo "github.com/outscale/osc-sdk-go/osc"
)

// Config ...
type Config struct {
	AccessKeyID string
	SecretKeyID string
	Region      string
	TokenID     string
	Endpoints   map[string]interface{}
}

//OutscaleClient client
type OutscaleClient struct {
	OSCAPI *oscgo.APIClient
}

// Client ...
func (c *Config) Client() (*OutscaleClient, error) {
	skipClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy:           http.ProxyFromEnvironment,
		},
	}

	skipClient.Transport = logging.NewTransport("Outscale", skipClient.Transport)

	skipClient.Transport = NewTransport(c.AccessKeyID, c.SecretKeyID, c.Region, skipClient.Transport)

	oscConfig := oscgo.NewConfiguration()
	oscConfig.Debug = true
	oscConfig.HTTPClient = skipClient

	oscClient := oscgo.NewAPIClient(oscConfig)

	client := &OutscaleClient{
		OSCAPI: oscClient,
	}

	return client, nil
}
