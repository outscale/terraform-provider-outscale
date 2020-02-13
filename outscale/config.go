package outscale

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/logging"
	"github.com/outscale/osc-go/oapi"

	oscgo "github.com/marinsalinas/osc-sdk-go"
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
	OAPI   *oapi.Client
	OSCAPI *oscgo.APIClient
}

// Client ...
func (c *Config) Client() (*OutscaleClient, error) {
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

	basePath := fmt.Sprintf("https://api.%s.outscale.com/oapi/latest", c.Region)

	if endpoint, ok := c.Endpoints["api"]; ok {
		basePath = endpoint.(string)
	}

	oscConfig := &oscgo.Configuration{
		BasePath:      basePath,
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
