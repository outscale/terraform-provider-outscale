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
	Config           string
	ConfigPlanChecks resource.ConfigPlanChecks
}

// Creates migration test steps which expects an empty plan after apply comparing plan of current and specified version
func FrameworkMigrationTestSteps(sdkVersion string, configs ...string) []resource.TestStep {
	migrationConfigs := lo.Map(configs, func(config string, _ int) MigrationTestConfig {
		return MigrationTestConfig{
			Config: config,
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectEmptyPlan(),
				},
			},
		}
	})
	return frameworkMigrationTestStepsWithConfigs(sdkVersion, migrationConfigs...)
}

// Creates migration test steps that aims to reproduce a refresh + empty plan.
// This uses the implicit legacy empty plan check by leaving ConfigPlanChecks empty.
func FrameworkMigrationTestStepsWithExpectNonEmptyPlan(sdkVersion string, configs ...string) []resource.TestStep {
	migrationConfigs := lo.Map(configs, func(config string, _ int) MigrationTestConfig {
		return MigrationTestConfig{
			Config: config,
			// ConfigPlanChecks is left nil/empty to trigger the implicit default check
		}
	})
	return frameworkMigrationTestStepsWithConfigs(sdkVersion, migrationConfigs...)
}

// Creates migration test steps with configurable plan checks
func FrameworkMigrationTestStepsWithConfigs(sdkVersion string, configs ...MigrationTestConfig) []resource.TestStep {
	return frameworkMigrationTestStepsWithConfigs(sdkVersion, configs...)
}

func frameworkMigrationTestStepsWithConfigs(sdkVersion string, configs ...MigrationTestConfig) []resource.TestStep {
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
		steps = append(steps, resource.TestStep{
			ExternalProviders: map[string]resource.ExternalProvider{
				"outscale": {
					VersionConstraint: sdkVersion,
					Source:            "outscale/outscale",
				},
			},
			Config: c.Config,
		})

		// If c.ConfigPlanChecks is set, it runs those checks.
		// If c.ConfigPlanChecks is empty, it runs the default "Expect Empty Plan" check.
		steps = append(steps, resource.TestStep{
			ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
			Config:                   c.Config,
			ConfigPlanChecks:         c.ConfigPlanChecks,
		})
		return steps
	})
}
