package outscale

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func TestAccOutscaleOutboundRule(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi == false {
		t.Skip()
	}
	var group fcu.SecurityGroup
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleSecurityGroupRuleEgressConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleRuleExists("outscale_firewall_rules_set.web", &group),
					testAccCheckOutscaleRuleAttributes("outscale_outbound_rule.egress_1", &group, nil, "egress"),
				),
			},
		},
	})
}

func testAccOutscaleSecurityGroupRuleEgressConfig(rInt int) string {
	return fmt.Sprintf(`
	resource "outscale_firewall_rules_set" "web" {
		group_name = "terraform_test_%d"
		group_description = "Used in the terraform acceptance tests"
		vpc_id = "vpc-e9d09d63"
		tag = {
						Name = "tf-acc-test"
		}
	}
	resource "outscale_outbound_rule" "egress_1" {
			ip_permissions = {
				ip_protocol = "tcp"
				from_port = 80
				to_port = 8000
				ip_ranges = ["10.0.0.0/8"]
		}
		group_id = "${outscale_firewall_rules_set.web.id}"
	}`, rInt)
}
