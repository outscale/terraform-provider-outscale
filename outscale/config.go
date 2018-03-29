package outscale

import (
	"github.com/terraform-providers/terraform-provider-outscale/osc"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
	"github.com/terraform-providers/terraform-provider-outscale/osc/icu"
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
	icu, err := icu.NewICUClient(config)
	if err != nil {
		return nil, err
	}
	client := &OutscaleClient{
		FCU: fcu,
		ICU: icu,
	}

	return client, nil
}

// OutscaleClient client
type OutscaleClient struct {
	FCU *fcu.Client
	ICU *icu.Client
}
