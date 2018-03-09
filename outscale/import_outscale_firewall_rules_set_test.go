package outscale

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOutscaleFirewallRulesSet_importBasic(t *testing.T) {
	resourceName := "outscale_firewall_rules_set.outscale_firewall_rules_set"

	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccAWSSecurityGroupRuleIngressConfig(rInt),
			},

			resource.TestStep{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"associate_public_ip_address", "user_data", "security_group", "tag"},
			},
		},
	})
}
