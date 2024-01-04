package outscale

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/version"
)

// Config ...
type Config struct {
	AccessKeyID string
	SecretKeyID string
	Region      string
	TokenID     string
	Endpoints   map[string]interface{}
	X509cert    string
	X509key     string
	Insecure    bool
}

// OutscaleClient client
type OutscaleClient struct {
	OSCAPI *oscgo.APIClient
}

// Client ...
func (c *Config) Client() (*OutscaleClient, error) {
	tlsconfig := &tls.Config{InsecureSkipVerify: c.Insecure}
	cert, err := tls.LoadX509KeyPair(c.X509cert, c.X509key)
	if err == nil {
		tlsconfig = &tls.Config{
			InsecureSkipVerify: false,
			Certificates:       []tls.Certificate{cert},
		}
	}

	skipClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsconfig,
			Proxy:           http.ProxyFromEnvironment,
		},
	}

	skipClient.Transport = logging.NewTransport("Outscale", skipClient.Transport)

	skipClient.Transport = NewTransport(c.AccessKeyID, c.SecretKeyID, c.Region, skipClient.Transport)

	basePath := fmt.Sprintf("api.%s.outscale.com", c.Region)
	if endpoint, ok := c.Endpoints["api"]; ok {
		basePath = endpoint.(string)
	}

	oscConfig := oscgo.NewConfiguration()
	oscConfig.Debug = true
	oscConfig.HTTPClient = skipClient
	oscConfig.Host = basePath
	oscConfig.UserAgent = fmt.Sprintf("terraform-provider-outscale/%s", version.GetVersion())

	oscClient := oscgo.NewAPIClient(oscConfig)

	client := &OutscaleClient{
		OSCAPI: oscClient,
	}
	return client, nil
}
