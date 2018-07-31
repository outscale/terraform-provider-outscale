package outscale

import (
	"crypto/tls"
	"net/http"
	"os"

	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"

	"github.com/terraform-providers/terraform-provider-outscale/osc"
	"github.com/terraform-providers/terraform-provider-outscale/osc/dl"
	"github.com/terraform-providers/terraform-provider-outscale/osc/eim"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
	"github.com/terraform-providers/terraform-provider-outscale/osc/icu"
	"github.com/terraform-providers/terraform-provider-outscale/osc/lbu"
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
	ICU  *icu.Client
	LBU  *lbu.Client
	EIM  *eim.Client
	DL   *dl.Client
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
	icu, err := icu.NewICUClient(config)
	if err != nil {
		return nil, err
	}
	lbu, err := lbu.NewLBUClient(config)
	if err != nil {
		return nil, err
	}
	eim, err := eim.NewEIMClient(config)
	if err != nil {
		return nil, err
	}
	dl, err := dl.NewDLClient(config)
	if err != nil {
		return nil, err
	}

	u := os.Getenv("OUTSCALE_OAPI_URL")

	oapicfg := &oapi.Config{
		AccessKey: c.AccessKeyID,
		SecretKey: c.SecretKeyID,
		Region:    c.Region,
		Service:   "oapi-gtw",
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
		ICU:  icu,
		LBU:  lbu,
		EIM:  eim,
		DL:   dl,
		OAPI: oapiClient,
	}

	return client, nil
}
