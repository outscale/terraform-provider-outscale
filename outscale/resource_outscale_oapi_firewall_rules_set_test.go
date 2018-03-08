package outscale

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func TestAccOutscaleOAPIFirewallRulesSet(t *testing.T) {
	var group fcu.SecurityGroup
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleOAPISGRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPIFirewallRulesSetConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPISecurityGroupRuleExists("outscale_firewall_rules_set.web", &group),
					resource.TestCheckResourceAttr(
						"outscale_firewall_rules_set.web", "firewall_rules_set_name", fmt.Sprintf("terraform_test_%d", rInt)),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPISGRuleDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient).FCU

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_firewall_rules_set" {
			continue
		}

		// Retrieve our group
		req := &fcu.DescribeSecurityGroupsInput{
			GroupIds: []*string{aws.String(rs.Primary.ID)},
		}

		var resp *fcu.DescribeSecurityGroupsOutput
		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err = conn.VM.DescribeSecurityGroups(req)

			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}

			return nil
		})

		if err == nil {
			if len(resp.SecurityGroups) > 0 && *resp.SecurityGroups[0].GroupId == rs.Primary.ID {
				return fmt.Errorf("Security Group (%s) still exists.", rs.Primary.ID)
			}

			return nil
		}

		if resp == nil {
			return nil
		}

		if strings.Contains(err.Error(), "InvalidGroup.NotFound") {
			return err
		}

		return err
	}

	return nil
}

func testAccCheckOutscaleOAPISecurityGroupRuleExists(n string, group *fcu.SecurityGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Security Group is set")
		}

		conn := testAccProvider.Meta().(*OutscaleClient).FCU
		req := &fcu.DescribeSecurityGroupsInput{
			GroupIds: []*string{aws.String(rs.Primary.ID)},
		}

		var resp *fcu.DescribeSecurityGroupsOutput
		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err = conn.VM.DescribeSecurityGroups(req)

			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}

			return nil
		})

		if err != nil {
			return err
		}

		if len(resp.SecurityGroups) > 0 && *resp.SecurityGroups[0].GroupId == rs.Primary.ID {
			*group = *resp.SecurityGroups[0]
			return nil
		}

		return fmt.Errorf("Security Group not found")
	}
}

func testAccOutscaleOAPIFirewallRulesSetConfig(rInt int) string {
	return fmt.Sprintf(`
	resource "outscale_firewall_rules_set" "web" {
		firewall_rules_set_name = "terraform_test_%d"
		description = "Used in the terraform acceptance tests"
		tag = {
						Name = "tf-acc-test"
		}
		lin_id = "vpc-e9d09d63"
	}`, rInt)
}
