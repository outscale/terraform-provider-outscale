package client

import (
	oscgo "github.com/outscale/osc-sdk-go/v2"
	sdkv3_oks "github.com/outscale/osc-sdk-go/v3/pkg/oks"
)

type OutscaleClient struct {
	OSCAPI *oscgo.APIClient
	OKS    *sdkv3_oks.Client
}
