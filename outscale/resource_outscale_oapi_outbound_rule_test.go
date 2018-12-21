package outscale

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"
)

func TestAccOutscaleOAPIOutboundRule(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapiFlag, err := strconv.ParseBool(o)
	if err != nil {
		oapiFlag = false
	}

	if !oapiFlag {
		t.Skip()
	}
	var group oapi.SecurityGroup
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleOAPISecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPISecurityGroupRuleEgressConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIRuleExists("outscale_security_group.outscale_security_group", &group),
					testAccCheckOutscaleOAPIRuleAttributes("outscale_security_group_rule.outscale_security_group_rule", &group, nil, "Inbound"),
				),
			},
		},
	})
}

func testAccOutscaleOAPISecurityGroupRuleEgressConfig(rInt int) string {
	return fmt.Sprintf(`
resource "outscale_security_group_rule" "outscale_security_group_rule" {
	flow              = "Inbound"
	rules {
		from_port_range   = 22
		to_port_range     = 22
		ip_protocol       = "tcp"
		ip_ranges         = ["46.231.147.8/32"]
		security_groups_members {
			security_group_id = "${outscale_security_group.outscale_security_group.security_group_id}"	
		}
	}
}

resource "outscale_security_group" "outscale_security_group" {
	description         = "test group"
	security_group_name = "sg1-test-group_test_%d"
}
`, rInt)
}
