package outscale

import (
	"crypto/tls"
	"net/http"

	"github.com/outscale/osc-go/oapi"

	"github.com/terraform-providers/terraform-provider-outscale/osc"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"

	"github.com/hashicorp/terraform/helper/logging"
)

// Config ...
type Config struct {
	AccessKeyID string
	SecretKeyID string
	Region      string
	TokenID     string
	OApi        bool
}

//OutscaleClient client
type OutscaleClient struct {
	FCU  *fcu.Client
	OAPI *oapi.Client
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

	oapiClient := oapi.NewClient(oapicfg, skipClient)

	client := &OutscaleClient{
		FCU:  fcu,
		ICU:  icu,
		OAPI: oapiClient,
	}

	return client, nil
}
