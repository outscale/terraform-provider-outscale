package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNet_WithSecurityGroupDataSource_basic(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "outscale_security_group.netSGtest"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleSecurityGroupConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "net_id"),
				),
			},
		},
	})
}

func TestAccOthers_WithSecurityGroupPublic(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "outscale_security_group.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleSecurityGroupPublicConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "security_group_name"),
				),
			},
		},
	})
}

func testAccDataSourceOutscaleSecurityGroupConfig(rInt int) string {
	return fmt.Sprintf(`
		resource "outscale_net" "vpc" {
				ip_range = "10.0.0.0/16"
				tags {
					key = "Name"
					value = "testacc-sec-group-ds"
				}
		}

		resource "outscale_security_group" "netSGtest" {
			net_id = outscale_net.vpc.id
			description = "Used in the terraform acceptance tests"
			security_group_name = "netSGtest-%d"
		}

		data "outscale_security_group" "by_filter" {
			filter {
				name = "security_group_ids"
				values = [outscale_security_group.netSGtest.security_group_id]
			}
		}`, rInt)
}

func testAccDataSourceOutscaleSecurityGroupPublicConfig(rInt int) string {
	return fmt.Sprintf(`
		resource "outscale_security_group" "test" {
			description = "Used in the terraform acceptance tests"
			security_group_name = "test-%d"
			tags {
				key = "tf-acctest"
				value = "%d"
			}
		}

		data "outscale_security_group" "by_filter_public" {
			filter {
		        name = "security_group_names"
				values = [outscale_security_group.test.security_group_name]
			}
		}`, rInt, rInt)
}
