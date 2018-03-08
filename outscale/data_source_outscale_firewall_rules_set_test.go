package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceOutscaleSecurityGroups_vpc(t *testing.T) {
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleSecurityGroupConfig_vpc(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.outscale_firewall_rules_set.by_id", "security_group_info.#", "3"),
					resource.TestCheckResourceAttr(
						"data.outscale_firewall_rules_set.by_filter", "security_group_info.#", "1"),
				),
			},
		},
	})
}

func testAccDataSourceOutscaleSecurityGroupConfig_vpc(rInt int) string {
	return fmt.Sprintf(`
	resource "outscale_firewall_rules_set" "test" {
		vpc_id = "vpc-e9d09d63"
		group_description = "Used in the terraform acceptance tests"
		group_name = "test-1--%d"
		tag = {
			Name = "tf-acctest"
			Seed = "%d"
		}
	}
	resource "outscale_firewall_rules_set" "test2" {
		vpc_id = "vpc-e9d09d63"
		group_description = "Used in the terraform acceptance tests"
		group_name = "test-2--%d"
		tag = {
			Name = "tf-acctest"
			Seed = "%d"
		}
	}
	resource "outscale_firewall_rules_set" "test3" {
		vpc_id = "vpc-e9d09d63"
		group_description = "Used in the terraform acceptance tests"
		group_name = "test-3--%d"
		tag = {
			Name = "tf-acctest"
			Seed = "%d"
		}
	}
	data "outscale_firewall_rules_set" "by_id" {
		group_id = ["${outscale_firewall_rules_set.test.id}", "${outscale_firewall_rules_set.test2.id}", "${outscale_firewall_rules_set.test3.id}"]
	}
	data "outscale_firewall_rules_set" "by_filter" {
		filter {
			name = "group-name"
			values = ["${outscale_firewall_rules_set.test.group_name}"]
		}
	}
	`, rInt, rInt, rInt, rInt, rInt, rInt)
}
