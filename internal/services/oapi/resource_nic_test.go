package oapi_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func TestAccNet_WithNic_basic(t *testing.T) {
	subregion := utils.GetRegion()
	resourceName := "outscale_nic.outscale_nic"
	sgName := acctest.RandomWithPrefix("testacc-sg")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),

		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleENIConfig(subregion, sgName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "private_ips.#", "2"),
				),
			},
			{
				Config: testAccOutscaleENIConfigUpdate(subregion, sgName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "subregion_name", fmt.Sprintf("%sa", subregion)),
					resource.TestCheckResourceAttr(resourceName, "private_ips.#", "3"),
				),
			},
		},
	})
}

func TestAccNet_Nic_Migration(t *testing.T) {
	subregion := utils.GetRegion()
	sgName := acctest.RandomWithPrefix("testacc-sg")

	resource.Test(t, resource.TestCase{
		Steps: testacc.FrameworkMigrationTestSteps("1.5.0",
			testAccOutscaleENIConfig(subregion, sgName),
			testAccOutscaleENIConfigUpdate(subregion, sgName),
		),
	})
}

func TestAccNet_WithNic_privateIPsValidation(t *testing.T) {
	subregion := utils.GetRegion()
	sgName := acctest.RandomWithPrefix("testacc-sg")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),

		Steps: []resource.TestStep{
			{
				Config:      testAccOutscaleENIConfigInvalidPrimary(subregion, sgName),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`At least one private_ips block must set is_primary = true`),
			},
		},
	})
}

func TestAccNet_WithNic_multiplePrimaryValidation(t *testing.T) {
	subregion := utils.GetRegion()
	sgName := acctest.RandomWithPrefix("testacc-sg")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),

		Steps: []resource.TestStep{
			{
				Config:      testAccOutscaleENIConfigMultiplePrimary(subregion, sgName),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`Only one private_ips block can set is_primary = true`),
			},
		},
	})
}

func TestAccNet_WithNic_primaryIPChangeRequiresReplace(t *testing.T) {
	subregion := utils.GetRegion()
	sgName := acctest.RandomWithPrefix("testacc-sg")
	resourceName := "outscale_nic.outscale_nic"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),

		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleENIConfigPrimaryOnly(subregion, sgName, "10.0.0.6"),
			},
			{
				Config: testAccOutscaleENIConfigPrimaryOnly(subregion, sgName, "10.0.0.7"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionReplace),
					},
				},
			},
		},
	})
}

func TestAccNet_WithNic_addSecondaryDoesNotReplace(t *testing.T) {
	subregion := utils.GetRegion()
	sgName := acctest.RandomWithPrefix("testacc-sg")
	resourceName := "outscale_nic.outscale_nic"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),

		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleENIConfigPrimaryOnly(subregion, sgName, "10.0.0.6"),
			},
			{
				Config: testAccOutscaleENIConfigPrimaryAndSecondary(subregion, sgName, "10.0.0.6", "10.0.0.7"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
			},
		},
	})
}

func TestAccNet_WithNic_addSecondaryDoesNotDiffAfterApply(t *testing.T) {
	subregion := utils.GetRegion()
	sgName := acctest.RandomWithPrefix("testacc-sg")
	resourceName := "outscale_nic.outscale_nic"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),

		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleENIConfigPrimaryOnly(subregion, sgName, "10.0.67.91"),
			},
			{
				Config: testAccOutscaleENIConfigPrimaryAndSecondary(subregion, sgName, "10.0.67.91", "10.0.67.92"),
			},
			{
				Config: testAccOutscaleENIConfigPrimaryAndSecondary(subregion, sgName, "10.0.67.91", "10.0.67.92"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionNoop),
					},
				},
			},
		},
	})
}

func TestAccNet_WithNic_omitPrivateIPsDoesNotDiff(t *testing.T) {
	subregion := utils.GetRegion()
	sgName := acctest.RandomWithPrefix("testacc-sg")
	resourceName := "outscale_nic.outscale_nic"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),

		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleENIConfigPrimaryOnly(subregion, sgName, "10.0.0.6"),
			},
			{
				Config: testAccOutscaleENIConfigWithoutPrivateIPs(subregion, sgName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionNoop),
					},
				},
			},
		},
	})
}

func TestAccNet_WithNic_omitSecurityGroupIDsDoesNotDiff(t *testing.T) {
	subregion := utils.GetRegion()
	sgName := acctest.RandomWithPrefix("testacc-sg")
	resourceName := "outscale_nic.outscale_nic"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),

		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleENIConfigWithSecurityGroupIDs(subregion, sgName),
			},
			{
				Config: testAccOutscaleENIConfigWithoutSecurityGroupIDs(subregion, sgName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionNoop),
					},
				},
			},
		},
	})
}

func TestAccNet_WithNic_privateIpTopLevel(t *testing.T) {
	subregion := utils.GetRegion()
	sgName := acctest.RandomWithPrefix("testacc-sg")
	resourceName := "outscale_nic.outscale_nic"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),

		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleENIConfigPrimaryOnly(subregion, sgName, "10.0.0.6"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "private_ip", "10.0.0.6"),
				),
			},
			{
				Config: testAccOutscaleENIConfigPrimaryOnly(subregion, sgName, "10.0.0.6"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionNoop),
					},
				},
			},
		},
	})
}

func TestAccNet_WithNic_securityGroupIdsUpdate(t *testing.T) {
	subregion := utils.GetRegion()
	sgName := acctest.RandomWithPrefix("testacc-sg-upd")
	resourceName := "outscale_nic.outscale_nic"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),

		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleENIConfigWithoutSecurityGroupIDs(subregion, sgName),
			},
			{
				Config: testAccOutscaleENIConfigWithSecurityGroupIDs(subregion, sgName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				Config: testAccOutscaleENIConfigWithSecurityGroupIDs(subregion, sgName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionNoop),
					},
				},
			},
		},
	})
}

func TestAccNet_Nic_MigrationWithSecurityGroupIds(t *testing.T) {
	subregion := utils.GetRegion()
	sgName := acctest.RandomWithPrefix("testacc-sg")

	resource.Test(t, resource.TestCase{
		Steps: testacc.FrameworkMigrationTestSteps("1.5.0",
			testAccOutscaleENIConfigWithSecurityGroupIDs(subregion, sgName),
			testAccOutscaleENIConfigWithoutSecurityGroupIDs(subregion, sgName),
		),
	})
}

func testAccOutscaleENIConfig(subregion, sgName string) string {
	return fmt.Sprintf(`
		resource "outscale_net" "outscale_net" {
			ip_range = "10.0.0.0/16"
			tags {
				key = "Name"
				value = "testacc-nic-rs"
			}
		}

		resource "outscale_subnet" "outscale_subnet" {
			subregion_name = "%sa"
			ip_range       = "10.0.0.0/16"
			net_id         = outscale_net.outscale_net.net_id
		}

		resource "outscale_security_group" "outscale_sg" {
			description         = "sg for terraform tests"
			security_group_name = "%[2]s"
			net_id              = outscale_net.outscale_net.net_id
		}

		resource "outscale_nic" "outscale_nic" {
			subnet_id          = outscale_subnet.outscale_subnet.subnet_id
			security_group_ids = [outscale_security_group.outscale_sg.security_group_id]

			private_ips {
				is_primary = true
				private_ip = "10.0.0.23"
			}

			private_ips	{
				is_primary = false
				private_ip = "10.0.0.46"
			}
		}
	`, subregion, sgName)
}

func testAccOutscaleENIConfigUpdate(subregion, sgName string) string {
	return fmt.Sprintf(`
		resource "outscale_net" "outscale_net" {
			ip_range = "10.0.0.0/16"
			tags {
				key = "Name"
				value = "testacc-nic-rs"
			}
		}

		resource "outscale_subnet" "outscale_subnet" {
			subregion_name = "%sa"
			ip_range       = "10.0.0.0/16"
			net_id         = outscale_net.outscale_net.net_id
		}

		resource "outscale_security_group" "outscale_sg" {
			description         = "sg for terraform tests"
			security_group_name = "%[2]s"
			net_id              = outscale_net.outscale_net.net_id
		}

		resource "outscale_nic" "outscale_nic" {
			subnet_id          = outscale_subnet.outscale_subnet.subnet_id
			security_group_ids = [outscale_security_group.outscale_sg.security_group_id]

			private_ips {
				is_primary = true
				private_ip = "10.0.0.23"
			}

			private_ips {
				is_primary = false
				private_ip = "10.0.0.46"
			}

			private_ips {
				is_primary = false
				private_ip = "10.0.0.69"
			}
		}
	`, subregion, sgName)
}

func testAccOutscaleENIConfigPrimaryOnly(subregion, sgName, privateIP string) string {
	return fmt.Sprintf(`
		resource "outscale_net" "outscale_net" {
			ip_range = "10.0.0.0/16"
			tags {
				key = "Name"
				value = "testacc-nic-rs"
			}
		}

		resource "outscale_subnet" "outscale_subnet" {
			subregion_name = "%sa"
			ip_range       = "10.0.0.0/16"
			net_id         = outscale_net.outscale_net.net_id
		}

		resource "outscale_security_group" "outscale_sg" {
			description         = "sg for terraform tests"
			security_group_name = "%[2]s"
			net_id              = outscale_net.outscale_net.net_id
		}

		resource "outscale_nic" "outscale_nic" {
			subnet_id          = outscale_subnet.outscale_subnet.subnet_id
			security_group_ids = [outscale_security_group.outscale_sg.security_group_id]

			private_ips {
				is_primary = true
				private_ip = "%[3]s"
			}
		}
	`, subregion, sgName, privateIP)
}

func testAccOutscaleENIConfigPrimaryAndSecondary(subregion, sgName, primaryIP, secondaryIP string) string {
	return fmt.Sprintf(`
		resource "outscale_net" "outscale_net" {
			ip_range = "10.0.0.0/16"
			tags {
				key = "Name"
				value = "testacc-nic-rs"
			}
		}

		resource "outscale_subnet" "outscale_subnet" {
			subregion_name = "%sa"
			ip_range       = "10.0.0.0/16"
			net_id         = outscale_net.outscale_net.net_id
		}

		resource "outscale_security_group" "outscale_sg" {
			description         = "sg for terraform tests"
			security_group_name = "%[2]s"
			net_id              = outscale_net.outscale_net.net_id
		}

		resource "outscale_nic" "outscale_nic" {
			subnet_id          = outscale_subnet.outscale_subnet.subnet_id
			security_group_ids = [outscale_security_group.outscale_sg.security_group_id]

			private_ips {
				is_primary = true
				private_ip = "%[3]s"
			}

			private_ips {
				private_ip = "%[4]s"
			}
		}
	`, subregion, sgName, primaryIP, secondaryIP)
}

func testAccOutscaleENIConfigInvalidPrimary(subregion, sgName string) string {
	return fmt.Sprintf(`
		resource "outscale_net" "outscale_net" {
			ip_range = "10.0.0.0/16"
			tags {
				key = "Name"
				value = "testacc-nic-rs"
			}
		}

		resource "outscale_subnet" "outscale_subnet" {
			subregion_name = "%sa"
			ip_range       = "10.0.0.0/16"
			net_id         = outscale_net.outscale_net.net_id
		}

		resource "outscale_security_group" "outscale_sg" {
			description         = "sg for terraform tests"
			security_group_name = "%[2]s"
			net_id              = outscale_net.outscale_net.net_id
		}

		resource "outscale_nic" "outscale_nic" {
			subnet_id          = outscale_subnet.outscale_subnet.subnet_id
			security_group_ids = [outscale_security_group.outscale_sg.security_group_id]

			private_ips {
				is_primary = false
				private_ip = "10.0.0.6"
			}

			private_ips {
				is_primary = false
				private_ip = "10.0.0.7"
			}
		}
	`, subregion, sgName)
}

func testAccOutscaleENIConfigMultiplePrimary(subregion, sgName string) string {
	return fmt.Sprintf(`
		resource "outscale_net" "outscale_net" {
			ip_range = "10.0.0.0/16"
			tags {
				key = "Name"
				value = "testacc-nic-rs"
			}
		}

		resource "outscale_subnet" "outscale_subnet" {
			subregion_name = "%sa"
			ip_range       = "10.0.0.0/16"
			net_id         = outscale_net.outscale_net.net_id
		}

		resource "outscale_security_group" "outscale_sg" {
			description         = "sg for terraform tests"
			security_group_name = "%[2]s"
			net_id              = outscale_net.outscale_net.net_id
		}

		resource "outscale_nic" "outscale_nic" {
			subnet_id          = outscale_subnet.outscale_subnet.subnet_id
			security_group_ids = [outscale_security_group.outscale_sg.security_group_id]

			private_ips {
				is_primary = true
				private_ip = "10.0.0.6"
			}

			private_ips {
				is_primary = true
				private_ip = "10.0.0.7"
			}
		}
	`, subregion, sgName)
}

func testAccOutscaleENIConfigWithoutPrivateIPs(subregion, sgName string) string {
	return fmt.Sprintf(`
		resource "outscale_net" "outscale_net" {
			ip_range = "10.0.0.0/16"
			tags {
				key = "Name"
				value = "testacc-nic-rs"
			}
		}

		resource "outscale_subnet" "outscale_subnet" {
			subregion_name = "%sa"
			ip_range       = "10.0.0.0/16"
			net_id         = outscale_net.outscale_net.net_id
		}

		resource "outscale_security_group" "outscale_sg" {
			description         = "sg for terraform tests"
			security_group_name = "%[2]s"
			net_id              = outscale_net.outscale_net.net_id
		}

		resource "outscale_nic" "outscale_nic" {
			subnet_id          = outscale_subnet.outscale_subnet.subnet_id
			security_group_ids = [outscale_security_group.outscale_sg.security_group_id]
		}
	`, subregion, sgName)
}

func testAccOutscaleENIConfigWithSecurityGroupIDs(subregion, sgName string) string {
	return fmt.Sprintf(`
		resource "outscale_net" "outscale_net" {
			ip_range = "10.0.0.0/16"
			tags {
				key = "Name"
				value = "testacc-nic-rs"
			}
		}

		resource "outscale_subnet" "outscale_subnet" {
			subregion_name = "%sa"
			ip_range       = "10.0.0.0/16"
			net_id         = outscale_net.outscale_net.net_id
		}

		resource "outscale_security_group" "outscale_sg" {
			description         = "sg for terraform tests"
			security_group_name = "%[2]s"
			net_id              = outscale_net.outscale_net.net_id
		}

		resource "outscale_nic" "outscale_nic" {
			subnet_id          = outscale_subnet.outscale_subnet.subnet_id
			security_group_ids = [outscale_security_group.outscale_sg.security_group_id]
		}
	`, subregion, sgName)
}

func testAccOutscaleENIConfigWithoutSecurityGroupIDs(subregion, sgName string) string {
	return fmt.Sprintf(`
		resource "outscale_net" "outscale_net" {
			ip_range = "10.0.0.0/16"
			tags {
				key = "Name"
				value = "testacc-nic-rs"
			}
		}

		resource "outscale_subnet" "outscale_subnet" {
			subregion_name = "%sa"
			ip_range       = "10.0.0.0/16"
			net_id         = outscale_net.outscale_net.net_id
		}

		resource "outscale_security_group" "outscale_sg" {
			description         = "sg for terraform tests"
			security_group_name = "%[2]s"
			net_id              = outscale_net.outscale_net.net_id
		}

		resource "outscale_nic" "outscale_nic" {
			subnet_id = outscale_subnet.outscale_subnet.subnet_id
		}
	`, subregion, sgName)
}
