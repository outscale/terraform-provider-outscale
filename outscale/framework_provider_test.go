package outscale

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	vers "github.com/terraform-providers/terraform-provider-outscale/version"
)

func TestFwProvider_impl(t *testing.T) {
	var _ provider.Provider = New(vers.GetVersion())
}

func TestAccFwPreCheck(t *testing.T) {
	if os.Getenv("OUTSCALE_ACCESSKEYID") == "" ||
		os.Getenv("OUTSCALE_REGION") == "" ||
		os.Getenv("OUTSCALE_SECRETKEYID") == "" ||
		os.Getenv("OUTSCALE_IMAGEID") == "" ||
		os.Getenv("OUTSCALE_ACCOUNT") == "" {
		t.Fatal("`OUTSCALE_ACCESSKEYID`, `OUTSCALE_SECRETKEYID`, `OUTSCALE_REGION`, `OUTSCALE_ACCOUNT` and `OUTSCALE_IMAGEID` must be set for acceptance testing")
	}
}

func TestMuxServer(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: map[string]func() (tfprotov5.ProviderServer, error){
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
		},
		Steps: []resource.TestStep{
			{
				Config: fwtestAccDataSourceOutscaleOAPIQuotaConfig,
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
				Config: fwtestAccDataSourceOutscaleOAPIQuotaConfig,
			},
			{
				ProtoV5ProviderFactories: map[string]func() (tfprotov5.ProviderServer, error){
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
				},
				Config:   fwtestAccDataSourceOutscaleOAPIQuotaConfig,
				PlanOnly: true,
			},
		},
	})
}

const fwtestAccDataSourceOutscaleOAPIQuotaConfig = `
	data "outscale_quota" "lbuquota1" { 
  filter {
        name     = "quota_names"
        values   = ["lb_listeners_limit"]
    }
}
`
