package outscale

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/outscale/terraform-provider-outscale/utils"
	vers "github.com/outscale/terraform-provider-outscale/version"
)

func TestFwProvider_impl(t *testing.T) {
	var _ provider.Provider = New(vers.GetVersion())
}

func TestAccFwPreCheck(t *testing.T) {
	if !utils.IsEnvVariableSet([]string{"OUTSCALE_ACCESSKEYID", "OUTSCALE_SECRETKEYID", "OUTSCALE_REGION", "OUTSCALE_ACCOUNT", "OUTSCALE_IMAGEID"}) {
		t.Fatal("`OUTSCALE_ACCESSKEYID`, `OUTSCALE_SECRETKEYID`, `OUTSCALE_REGION`, `OUTSCALE_ACCOUNT` and `OUTSCALE_IMAGEID` must be set for acceptance testing")
	}
}

func protoV5ProviderFactories() map[string]func() (tfprotov5.ProviderServer, error) {
	return map[string]func() (tfprotov5.ProviderServer, error){
		"outscale": providerserver.NewProtocol5WithError(New(vers.GetVersion())),
	}
}

func TestMuxServer(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: defineTestProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fwtestAccDataSourceOutscaleQuotaConfig,
			},
		},
	})
}

func TestDataSource_UpgradeFromVersion(t *testing.T) {
	resource.Test(t, resource.TestCase{

		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"outscale": {
						VersionConstraint: "0.10.0",
						Source:            "outscale/outscale",
					},
				},
				Config: fwtestAccDataSourceOutscaleQuotaConfig,
			},
			{
				ProtoV5ProviderFactories: protoV5ProviderFactories(),
				Config:                   fwtestAccDataSourceOutscaleQuotaConfig,
				PlanOnly:                 true,
			},
		},
	})
}

const fwtestAccDataSourceOutscaleQuotaConfig = `
	data "outscale_quota" "lbuquota1" {
        filter {
            name     = "quota_names"
            values   = ["lb_listeners_limit"]
    }
}
`

func defineTestProviderFactories() map[string]func() (tfprotov5.ProviderServer, error) {
	return map[string]func() (tfprotov5.ProviderServer, error){
		"outscale": func() (tfprotov5.ProviderServer, error) {
			ctx := context.Background()
			providers := []func() tfprotov5.ProviderServer{
				providerserver.NewProtocol5(New(vers.GetVersion())), // Example terraform-plugin-framework provider
				Provider().GRPCProvider,                             // Example terraform-plugin-sdk provider
			}
			muxServer, err := tf5muxserver.NewMuxServer(ctx, providers...)
			if err != nil {
				return nil, err
			}
			return muxServer.ProviderServer(), nil
		},
	}
}
