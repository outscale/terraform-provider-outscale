package outscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
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
