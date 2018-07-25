package outscale

import (
	"crypto/tls"
	"net/http"
	"os"

	"github.com/terraform-providers/terraform-provider-outscale/osc"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"
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

	u := os.Getenv("OUTSCALE_OAPI")

	oapicfg := &oapi.Config{
		AccessKey: c.AccessKeyID,
		SecretKey: c.SecretKeyID,
		Region:    c.Region,
		Service:   "oapi",
		URL:       u,
	}

	skipClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	oapiClient := oapi.NewClient(oapicfg, skipClient)

	client := &OutscaleClient{
		FCU:  fcu,
		OAPI: oapiClient,
	}

	return client, nil
}

// OutscaleClient client
type OutscaleClient struct {
	FCU  *fcu.Client
	OAPI *oapi.Client
}
