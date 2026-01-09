package testacc

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
	"github.com/outscale/terraform-provider-outscale/provider"
	"github.com/outscale/terraform-provider-outscale/version"
	"github.com/samber/lo"
)

var (
	// ConfiguredClient is the client configured during provider setup in tests
	ConfiguredClient *client.OutscaleClient

	// SDKProvider is the SDK v2 provider instance for tests
	SDKProvider  = provider.Provider()
	SDKProviders = map[string]*schema.Provider{
		"outscale": SDKProvider,
	}
)

func PreCheck(t *testing.T) {
	if utils.GetEnvVariableValue([]string{"OSC_ACCESS_KEY", "OUTSCALE_ACCESSKEYID"}) == "" ||
		utils.GetEnvVariableValue([]string{"OSC_SECRET_KEY", "OUTSCALE_SECRETKEYID"}) == "" ||
		utils.GetEnvVariableValue([]string{"OSC_REGION", "OUTSCALE_REGION"}) == "" ||
		utils.GetEnvVariableValue([]string{"OSC_ACCOUNT_ID", "OUTSCALE_ACCOUNT"}) == "" ||
		utils.GetEnvVariableValue([]string{"OSC_IMAGE_ID", "OUTSCALE_IMAGEID"}) == "" {
		t.Fatal("`OSC_ACCESS_KEY`, `OSC_SECRET_KEY`, `OSC_REGION`, `OSC_ACCOUNT_ID` and `OUTSCALE_IMAGEID` must be set for acceptance testing")
	}
}

// Returns a map of provider factories for testing with protocol v6
// This includes both the SDK v2 provider (upgraded to v6) and the Framework provider
func ProtoV6ProviderFactories() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"outscale": func() (tfprotov6.ProviderServer, error) {
			ctx := context.Background()

			upgradedSdkServer, err := tf5to6server.UpgradeServer(
				ctx,
				provider.Provider().GRPCProvider,
			)
			if err != nil {
				return nil, err
			}

			// Create Framework provider with configure callback to capture client
			frameworkProvider := provider.NewWithConfigure(version.GetVersion(), func(c *client.OutscaleClient) {
				ConfiguredClient = c
			})

			providers := []func() tfprotov6.ProviderServer{
				providerserver.NewProtocol6(frameworkProvider),
				func() tfprotov6.ProviderServer {
					return upgradedSdkServer
				},
			}

			// Mux both providers together
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

// Creates migration test steps which expects an empty plan after apply
// comparing plan of current and specified version
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

// Creates migration test steps that aims to reproduce
// a refresh + empty plan. This uses the implicit legacy empty plan check by leaving ConfigPlanChecks empty.
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
				ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
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
			ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
			Config:                   c.Config,
			ConfigPlanChecks:         c.ConfigPlanChecks,
		})
		return steps
	})
}

func ImportStepWithStateIdFunc(resourceName string, importStateIdFunc resource.ImportStateIdFunc, ignore ...string) resource.TestStep {
	return importStep(resourceName, importStateIdFunc, ignore...)
}

func ImportStep(resourceName string, ignore ...string) resource.TestStep {
	idFunc := func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}
		return rs.Primary.ID, nil
	}
	return importStep(resourceName, idFunc, ignore...)
}

func importStep(resourceName string, importStateIdFunc resource.ImportStateIdFunc, ignore ...string) resource.TestStep {
	step := resource.TestStep{
		ResourceName:      resourceName,
		ImportState:       true,
		ImportStateVerify: true,
		ImportStateIdFunc: importStateIdFunc,
	}

	if len(ignore) > 0 {
		step.ImportStateVerifyIgnore = ignore
	}

	return step
}

func DefaultIgnores() []string {
	return []string{
		"request_id",
	}
}
