package outscale

import "github.com/terraform-providers/terraform-provider-outscale/osc"

type Config struct {
	AccessKeyId string
	SecretKeyId string
	OApi        bool
}

func (c *Config) Client() (*osc.Client, error) {
	return nil, nil
}
