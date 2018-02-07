package outscale

import "github.com/terraform-providers/terraform-provider-outscale/osc"

// Config ...
type Config struct {
	AccessKeyID string
	SecretKeyID string
	TokenID     string
	OApi        bool
}

// Client ...
func (c *Config) Client() (*osc.Client, error) {
	return nil, nil
}
