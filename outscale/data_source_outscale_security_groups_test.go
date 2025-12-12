package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNet_WithSecurityGroups_DataSource(t *testing.T) {
	rInt := acctest.RandInt()
	resouceName1 := "data.outscale_security_groups.by_id"
	resouceName2 := "data.outscale_security_groups.by_filter"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleSecurityGroupConfigVPC(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resouceName1, "security_groups.#", "3"),
					resource.TestCheckResourceAttr(resouceName2, "security_groups.#", "1"),
				),
			},
		},
	})
}

func testAccDataSourceOutscaleSecurityGroupConfigVPC(rInt int) string {
	return fmt.Sprintf(`
		resource "outscale_net" "outscale_net" {
			ip_range = "10.0.0.0/16"
		}

		resource "outscale_security_group" "test" {
			net_id = "${outscale_net.outscale_net.id}"
			description = "Used in the terraform acceptance tests"
			security_group_name = "test-1-%[1]d"
			tags {
				key = "tf-acctest"
				value = "%[1]d"
			}
		}

		resource "outscale_security_group" "test2" {
			net_id = "${outscale_net.outscale_net.id}"
			description = "Used in the terraform acceptance tests"
			security_group_name = "test-2-%[1]d"
			tags {
				key = "tf-acctest"
				value = "%[1]d"
			}
		}

		resource "outscale_security_group" "test3" {
			net_id = "${outscale_net.outscale_net.id}"
			description = "Used in the terraform acceptance tests"
			security_group_name = "test-3-%[1]d"
			tags {
				key = "tf-acctest"
				value = "%[1]d"
			}
		}

		data "outscale_security_groups" "by_id" {
			security_group_ids = [outscale_security_group.test.id, outscale_security_group.test2.id, outscale_security_group.test3.id]
		}

		data "outscale_security_groups" "by_filter" {
			filter {
				name = "security_group_names"
				values = [outscale_security_group.test.security_group_name]
			}
		}
	`, rInt)
}
