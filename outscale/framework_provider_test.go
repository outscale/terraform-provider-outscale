package outscale

import (
	"context"
	"testing"

	sdkresource "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	vers "github.com/outscale/terraform-provider-outscale/version"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/outscale/terraform-provider-outscale/utils"
	"github.com/outscale/terraform-provider-outscale/version"
)

func TestFwProvider_impl(t *testing.T) {
	var _ provider.Provider = New(vers.GetVersion())
}

func TestAccFwPreCheck(t *testing.T) {
	if !utils.IsEnvVariableSet([]string{"OUTSCALE_ACCESSKEYID", "OUTSCALE_SECRETKEYID", "OUTSCALE_REGION", "OUTSCALE_ACCOUNT", "OUTSCALE_IMAGEID"}) {
		t.Fatal("`OUTSCALE_ACCESSKEYID`, `OUTSCALE_SECRETKEYID`, `OUTSCALE_REGION`, `OUTSCALE_ACCOUNT` and `OUTSCALE_IMAGEID` must be set for acceptance testing")
	}
}

func TestMuxServer(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
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
				ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
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

var testAccConfiguredClient *OutscaleClientFW

func DefineTestProviderFactoriesV6() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"outscale": func() (tfprotov6.ProviderServer, error) {
			ctx := context.Background()
			upgradedSdkServer, err := tf5to6server.UpgradeServer(
				ctx,
				Provider().GRPCProvider,
			)

			if err != nil {
				return nil, err
			}

			p := NewWithConfigure(version.GetVersion(), func(client *OutscaleClientFW) {
				testAccConfiguredClient = client
			}).(*frameworkProvider)

			providers := []func() tfprotov6.ProviderServer{
				providerserver.NewProtocol6(p),
				func() tfprotov6.ProviderServer {
					return upgradedSdkServer
				},
			}

			muxServer, err := tf6muxserver.NewMuxServer(ctx, providers...)

			if err != nil {
				return nil, err
			}

			return muxServer.ProviderServer(), nil
		},
	}
}

func FrameworkMigrationTestSteps(sdkVersion string, config string) []sdkresource.TestStep {
	return []sdkresource.TestStep{
		{
			ExternalProviders: map[string]sdkresource.ExternalProvider{
				"outscale": {
					VersionConstraint: sdkVersion,
					Source:            "outscale/outscale",
				},
			},
			Config: config,
		},
		{
			ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
			Config:                   config,
			PlanOnly:                 true,
			ExpectNonEmptyPlan:       false,
		},
	}
}
