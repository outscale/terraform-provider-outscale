package outscale

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func TestAccOutscaleSecurityGroupRule_Ingress_VPC(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi {
		t.Skip()
	}

	var group fcu.SecurityGroup
	rInt := acctest.RandInt()

	testRuleCount := func(*terraform.State) error {
		if len(group.IpPermissions) != 1 {
			return fmt.Errorf("Wrong Security Group rule count, expected %d, got %d",
				1, len(group.IpPermissions))
		}

		rule := group.IpPermissions[0]
		if *rule.FromPort != int64(80) {
			return fmt.Errorf("Wrong Security Group port setting, expected %d, got %d",
				80, int(*rule.FromPort))
		}

		return nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleSecurityGroupRuleIngressConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleSecurityGroupRuleExists("outscale_firewall_rules_set.web", &group),
					testAccCheckOutscaleSecurityGroupInboundRuleAttributes("outscale_inbound_rule.ingress_1", &group, nil, "ingress"),
					resource.TestCheckResourceAttr(
						"outscale_inbound_rule.ingress_1", "ip_permissions.0.from_port", "80"),
					testRuleCount,
				),
			},
		},
	})
}

func TestAccOutscaleSecurityGroupRule_MultiIngress(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi {
		t.Skip()
	}

	var group fcu.SecurityGroup

	testMultiRuleCount := func(*terraform.State) error {
		if len(group.IpPermissions) != 2 {
			return fmt.Errorf("Wrong Security Group rule count, expected %d, got %d",
				2, len(group.IpPermissions))
		}

		var rule *fcu.IpPermission
		for _, r := range group.IpPermissions {
			if *r.FromPort == int64(80) {
				rule = r
			}
		}

		if *rule.ToPort != int64(8000) {
			return fmt.Errorf("Wrong Security Group port 2 setting, expected %d, got %d",
				8000, int(*rule.ToPort))
		}

		return nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleSecurityGroupRuleConfigMultiIngress,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleSecurityGroupRuleExists("outscale_firewall_rules_set.web", &group),
					testMultiRuleCount,
				),
			},
		},
	})
}

func testAccCheckOutscaleSecurityGroupInboundRuleAttributes(n string, group *fcu.SecurityGroup, p *fcu.IpPermission, ruleType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Security Group Rule Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Security Group Rule is set")
		}

		if p == nil {
			p = &fcu.IpPermission{
				FromPort:   aws.Int64(80),
				ToPort:     aws.Int64(8000),
				IpProtocol: aws.String("tcp"),
				IpRanges:   []*fcu.IpRange{{CidrIp: aws.String("10.0.0.0/8")}},
			}
		}

		var matchingRule *fcu.IpPermission
		var rules []*fcu.IpPermission
		rules = group.IpPermissions

		if len(rules) == 0 {
			return fmt.Errorf("No IPPerms")
		}

		for _, r := range rules {
			if r.ToPort != nil && *p.ToPort != *r.ToPort {
				continue
			}

			if r.FromPort != nil && *p.FromPort != *r.FromPort {
				continue
			}

			if r.IpProtocol != nil && *p.IpProtocol != *r.IpProtocol {
				continue
			}

			remaining := len(p.IpRanges)
			for _, ip := range p.IpRanges {
				for _, rip := range r.IpRanges {
					if *ip.CidrIp == *rip.CidrIp {
						remaining--
					}
				}
			}

			if remaining > 0 {
				continue
			}

			remaining = len(p.UserIdGroupPairs)
			for _, ip := range p.UserIdGroupPairs {
				for _, rip := range r.UserIdGroupPairs {
					if *ip.GroupId == *rip.GroupId {
						remaining--
					}
				}
			}

			if remaining > 0 {
				continue
			}

			remaining = len(p.PrefixListIds)
			for _, pip := range p.PrefixListIds {
				for _, rpip := range r.PrefixListIds {
					if *pip.PrefixListId == *rpip.PrefixListId {
						remaining--
					}
				}
			}

			if remaining > 0 {
				continue
			}

			matchingRule = r
		}

		if matchingRule != nil {
			return nil
		}

		return fmt.Errorf("Error here\n\tlooking for %+v, wasn't found in %+v", p, rules)
	}
}

func testAccOutscaleSecurityGroupRuleIngressConfig(rInt int) string {
	return fmt.Sprintf(`
	resource "outscale_firewall_rules_set" "web" {
		group_name = "terraform_test_%d"
		group_description = "Used in the terraform acceptance tests"
					tag {
									Name = "tf-acc-test"
					}
	}
	resource "outscale_inbound_rule" "ingress_1" {
		ip_permissions = {
			ip_protocol = "tcp"
		from_port = 80
		to_port = 8000
		ip_ranges = ["10.0.0.0/8"]
		}
		group_id = "${outscale_firewall_rules_set.web.id}"
	}`, rInt)
}

const testAccOutscaleSecurityGroupRuleConfigMultiIngress = `
resource "outscale_firewall_rules_set" "web" {
  group_name = "terraform_acceptance_test_example_2"
  group_description = "Used in the terraform acceptance tests"
}
resource "outscale_firewall_rules_set" "worker" {
  group_name = "terraform_acceptance_test_example_worker"
  group_description = "Used in the terraform acceptance tests"
}
resource "outscale_inbound_rule" "ingress_1" {
  ip_permissions = {
		ip_protocol = "tcp"
  from_port = 22
  to_port = 22
  ip_ranges = ["10.0.0.0/8"]
	}
  group_id = "${outscale_firewall_rules_set.web.id}"
}
resource "outscale_inbound_rule" "ingress_2" {
  ip_permissions = {
		ip_protocol = "tcp"
  from_port = 80
  to_port = 8000
	}
  group_id = "${outscale_firewall_rules_set.web.id}"
}
`
