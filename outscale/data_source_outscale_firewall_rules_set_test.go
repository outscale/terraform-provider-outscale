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
					// resource.TestCheckResourceAttr(
					// 	"data.outscale_firewall_rules_set.by_id", "security_group_info.#", "3"),
					resource.TestCheckResourceAttr(
						"data.outscale_firewall_rules_set.by_filter_public", "security_group_info.#", "1"),
					resource.TestCheckResourceAttr(
						"data.outscale_firewall_rules_set.by_filter_public", "security_group_info.0.ip_permissions_egress.#", "1"),
				),
			},
		},
	})
}

func testAccDataSourceOutscaleSecurityGroupConfig_vpc(rInt int) string {
	return fmt.Sprintf(`
		resource "outscale_outbound_rule" "outscale_outbound_rule1" {
	ip_permissions = {
		from_port = 22
		to_port = 22
		ip_protocol = "tcp"
		ip_ranges = ["46.231.147.8/32"]
	}

	group_id = "${outscale_firewall_rules_set.outscale_firewall_rules_set.id}"
}

resource "outscale_inbound_rule" "outscale_inbound_rule1" {
	ip_permissions = {
		from_port = 22
		to_port = 22
		ip_protocol = "tcp"
		ip_ranges = ["46.231.147.8/32"]
	}

	group_id = "${outscale_firewall_rules_set.outscale_firewall_rules_set.id}"
}

resource "outscale_inbound_rule" "outscale_inbound_rule2" {
	 ip_permissions = {
		from_port = 443
		to_port = 443
		ip_protocol = "tcp"
		ip_ranges = ["46.231.147.8/32"]
	}

	group_id = "${outscale_firewall_rules_set.outscale_firewall_rules_set.id}"
}

resource "outscale_firewall_rules_set" "outscale_firewall_rules_set" {
		group_description = "Used in the terraform acceptance tests"
		group_name = "test-%d"
		vpc_id = "vpc-e9d09d63"
		tag = {
			Name = "tf-acctest"
			Seed = "%d"
		}
	}
	data "outscale_firewall_rules_set" "by_filter_public" {
		filter {
			name = "group-name"
			values = ["${outscale_firewall_rules_set.outscale_firewall_rules_set.group_name}"]
		}
	}`, rInt, rInt)
}
