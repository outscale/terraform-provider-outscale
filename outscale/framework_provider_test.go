package outscale

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/outscale/terraform-provider-outscale/utils"
	"github.com/outscale/terraform-provider-outscale/version"
	"github.com/samber/lo"
)

func TestFwProvider_impl(t *testing.T) {
	var _ provider.Provider = New(version.GetVersion())
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

type MigrationTestConfig struct {
	Config                  string
	ExpectUpdateActionsAddr []string
}

func FrameworkMigrationTestSteps(sdkVersion string, configs ...string) []resource.TestStep {
	migrationConfigs := lo.Map(configs, func(config string, _ int) MigrationTestConfig {
		return MigrationTestConfig{Config: config}
	})
	return frameworkMigrationTestStepsWithOptions(sdkVersion, migrationConfigs...)
}

// Creates migration test steps with expected update actions (without resource replacement)
func FrameworkMigrationTestStepsWithUpdate(sdkVersion string, configs ...MigrationTestConfig) []resource.TestStep {
	return frameworkMigrationTestStepsWithOptions(sdkVersion, configs...)
}

func frameworkMigrationTestStepsWithOptions(sdkVersion string, configs ...MigrationTestConfig) []resource.TestStep {
	return lo.FlatMap(configs, func(c MigrationTestConfig, i int) []resource.TestStep {
		var steps []resource.TestStep

		// If not the first config, destroy the previous one first to avoid provider init conflict
		if i > 0 {
			steps = append(steps, resource.TestStep{
				ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
				Config:                   configs[i-1].Config,
				Destroy:                  true,
			})
		}

		var planChecks []plancheck.PlanCheck
		if len(c.ExpectUpdateActionsAddr) > 0 {
			for _, addr := range c.ExpectUpdateActionsAddr {
				planChecks = append(planChecks,
					plancheck.ExpectResourceAction(addr, plancheck.ResourceActionUpdate),
				)
			}
		} else {
			planChecks = []plancheck.PlanCheck{
				plancheck.ExpectEmptyPlan(),
			}
		}

		return append(steps,
			resource.TestStep{
				ExternalProviders: map[string]resource.ExternalProvider{
					"outscale": {
						VersionConstraint: sdkVersion,
						Source:            "outscale/outscale",
					},
				},
				Config: c.Config,
			},
			resource.TestStep{
				ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
				Config:                   c.Config,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: planChecks,
				},
			},
		)
	})
}
