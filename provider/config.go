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
	IAASEndpoint string
	IAASRegion   string
	IAASX509Cert string
	IAASX509Key  string
	IAASInsecure bool
	OKSEndpoint  string
	OKSRegion    string
}

func (c *Config) Client() (*client.OutscaleClient, error) {
	oscClient, err := client.NewOAPIClient(client.Config{
		AccessKeyID:  c.AccessKeyID,
		SecretKeyID:  c.SecretKeyID,
		Region:       c.IAASRegion,
		APIEndpoint:  c.IAASEndpoint,
		X509CertPath: c.IAASX509Cert,
		X509KeyPath:  c.IAASX509Key,
		Insecure:     c.IAASInsecure,
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
