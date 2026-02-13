package provider

import (
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/version"
)

var UserAgent = "terraform-provider-outscale/" + version.GetVersion()

type Config struct {
	AccessKeyID string
	SecretKeyID string
	Region      string
	TokenID     string
	ConfigFile  string
	Profile     string

	// deprecated fields
	Endpoints    map[string]string
	X509CertPath string
	X509KeyPath  string
	Insecure     bool

	// per-service fields
	APIEndpoint string
	APIRegion   string
	APIX509Cert string
	APIX509Key  string
	APIInsecure bool
	OKSEndpoint string
	OKSRegion   string
}

func (c *Config) Client() (*client.OutscaleClient, error) {
	oscClient, err := client.NewOAPIClient(client.Config{
		AccessKeyID:  c.AccessKeyID,
		SecretKeyID:  c.SecretKeyID,
		Region:       c.APIRegion,
		APIEndpoint:  c.APIEndpoint,
		X509CertPath: c.APIX509Cert,
		X509KeyPath:  c.APIX509Key,
		Insecure:     c.APIInsecure,
		UserAgent:    UserAgent,
	})
	if err != nil {
		return nil, err
	}

	oksClient, err := client.NewOKSClient(client.Config{
		AccessKeyID: c.AccessKeyID,
		SecretKeyID: c.SecretKeyID,
		Region:      c.OKSRegion,
		OKSEndpoint: c.OKSEndpoint,
		UserAgent:   UserAgent,
	})
	if err != nil {
		return nil, err
	}

	return &client.OutscaleClient{
		OSCAPI: oscClient,
		OKS:    oksClient,
	}, nil
}
