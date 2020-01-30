package outscale

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/outscale/osc-go/oapi"

	"github.com/hashicorp/terraform/helper/logging"

	oscgo "github.com/marinsalinas/osc-sdk-go"
)

// Config ...
type Config struct {
	AccessKeyID string
	SecretKeyID string
	Region      string
	TokenID     string
	OApi        bool
}

//OutscaleClient client
type OutscaleClient struct {
	OAPI   *oapi.Client
	OSCAPI *oscgo.APIClient
}

// Client ...
func (c *Config) Client() (*OutscaleClient, error) {
	config := osc.Config{
		Credentials: &osc.Credentials{
			AccessKey: c.AccessKeyID,
			SecretKey: c.SecretKeyID,
			Region:    c.Region,
		},
	}
	if err != nil {
		return nil, err
	}

	oapicfg := &oapi.Config{
		AccessKey: c.AccessKeyID,
		SecretKey: c.SecretKeyID,
		Region:    c.Region,
		Service:   "api",
		URL:       "outscale.com/oapi/latest",
	}

	skipClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	skipClient.Transport = logging.NewTransport("Outscale", skipClient.Transport)

	skipClient.Transport = oscgo.NewTransport(c.AccessKeyID, c.SecretKeyID, c.Region, skipClient.Transport)

	oscConfig := &oscgo.Configuration{
		BasePath:      fmt.Sprintf("https://api.%s.outscale.com/oapi/latest", c.Region),
		DefaultHeader: make(map[string]string),
		UserAgent:     "terraform-provider-outscale-dev",
		HTTPClient:    skipClient,
	}

	oscClient := oscgo.NewAPIClient(oscConfig)

	oapiClient := oapi.NewClient(oapicfg, skipClient)

	client := &OutscaleClient{
		OAPI:   oapiClient,
		OSCAPI: oscClient,
	}

	return client, nil
}
