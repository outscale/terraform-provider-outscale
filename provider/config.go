package provider

import (
	"errors"

	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/version"
)

var UserAgent = "terraform-provider-outscale/" + version.GetVersion()

// Config ...
type Config struct {
	AccessKeyID  string
	SecretKeyID  string
	Region       string
	TokenID      string
	Endpoints    map[string]string
	X509CertPath string
	X509KeyPath  string
	Insecure     bool
	ConfigFile   string
	Profile      string
}

// Client ...
func (c *Config) Client() (*client.OutscaleClient, error) {
	endpoint := c.Endpoints["api"]
	if c.Region == "" && endpoint == "" {
		return nil, errors.New("'region' or 'endpoints' must be set for provider configuration")
	}

	oscClient, err := client.NewOAPIClient(client.Config{
		AccessKeyID:  c.AccessKeyID,
		SecretKeyID:  c.SecretKeyID,
		Region:       c.Region,
		APIEndpoint:  endpoint,
		X509CertPath: c.X509CertPath,
		X509KeyPath:  c.X509KeyPath,
		Insecure:     c.Insecure,
		UserAgent:    UserAgent,
	})
	if err != nil {
		return nil, err
	}

	outscaleClient := &client.OutscaleClient{
		OSCAPI: oscClient,
	}
	return outscaleClient, nil
}
