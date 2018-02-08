package outscale

import (
	"github.com/terraform-providers/terraform-provider-outscale/osc"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

// Config ...
type Config struct {
	AccessKeyID string
	SecretKeyID string
	Region      string
	TokenID     string
	OApi        bool
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
	fcu, err := fcu.NewFCUClient(config)
	if err != nil {
		return nil, err
	}

	client := &OutscaleClient{
		FCU: fcu,
	}

	return client, nil
}

// OutscaleClient client
type OutscaleClient struct {
	FCU *fcu.Client
}
