package outscale

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"

	"github.com/aws/aws-sdk-go/aws"
)

func TestIpPermissionIDHash(t *testing.T) {
	simple := &fcu.IpPermission{
		IpProtocol: aws.String("tcp"),
		FromPort:   aws.Int64(int64(80)),
		ToPort:     aws.Int64(int64(8000)),
		IpRanges: []*fcu.IpRange{
			{
				CidrIp: aws.String("10.0.0.0/8"),
			},
		},
	}

	egress := &fcu.IpPermission{
		IpProtocol: aws.String("tcp"),
		FromPort:   aws.Int64(int64(80)),
		ToPort:     aws.Int64(int64(8000)),
		IpRanges: []*fcu.IpRange{
			{
				CidrIp: aws.String("10.0.0.0/8"),
			},
		},
	}

	egress_all := &fcu.IpPermission{
		IpProtocol: aws.String("-1"),
		IpRanges: []*fcu.IpRange{
			{
				CidrIp: aws.String("10.0.0.0/8"),
			},
		},
	}

	vpc_security_group_source := &fcu.IpPermission{
		IpProtocol: aws.String("tcp"),
		FromPort:   aws.Int64(int64(80)),
		ToPort:     aws.Int64(int64(8000)),
		UserIdGroupPairs: []*fcu.UserIdGroupPair{
			{
				UserId:  aws.String("987654321"),
				GroupId: aws.String("sg-12345678"),
			},
			{
				UserId:  aws.String("123456789"),
				GroupId: aws.String("sg-987654321"),
			},
			{
				UserId:  aws.String("123456789"),
				GroupId: aws.String("sg-12345678"),
			},
		},
	}

	security_group_source := &fcu.IpPermission{
		IpProtocol: aws.String("tcp"),
		FromPort:   aws.Int64(int64(80)),
		ToPort:     aws.Int64(int64(8000)),
		UserIdGroupPairs: []*fcu.UserIdGroupPair{
			{
				UserId:    aws.String("987654321"),
				GroupName: aws.String("my-security-group"),
			},
			{
				UserId:    aws.String("123456789"),
				GroupName: aws.String("my-security-group"),
			},
			{
				UserId:    aws.String("123456789"),
				GroupName: aws.String("my-other-security-group"),
			},
		},
	}

	// hardcoded hashes, to detect future change
	cases := []struct {
		Input  *fcu.IpPermission
		Type   string
		Output string
	}{
		{simple, "ingress", "sgrule-3403497314"},
		{egress, "egress", "sgrule-1173186295"},
		{egress_all, "egress", "sgrule-766323498"},
		{vpc_security_group_source, "egress", "sgrule-351225364"},
		{security_group_source, "egress", "sgrule-2198807188"},
	}

	for _, tc := range cases {
		actual := ipPermissionIDHash("sg-12345", tc.Type, tc.Input)
		if actual != tc.Output {
			t.Errorf("input: %s - %s\noutput: %s", tc.Type, tc.Input, actual)
		}
	}
}

func TestAccOutscaleOutboundRule(t *testing.T) {
	var group fcu.SecurityGroup
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOutboundRuleConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleSecurityGroupRuleExists("outscale_firewall_rules_set.web", &group),
					testAccCheckOutscaleSecurityGroupRuleAttributes("outscale_outbound_rule.egress_1", &group, nil, "egress"),
				),
			},
		},
	})
}

func testAccCheckOutscaleSecurityGroupRuleAttributes(n string, group *fcu.SecurityGroup, p *fcu.IpPermission, ruleType string) resource.TestCheckFunc {
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
		if ruleType == "ingress" {
			rules = group.IpPermissions
		} else {
			rules = group.IpPermissionsEgress
		}

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
			log.Printf("[DEBUG] Matching rule found : %s", matchingRule)
			return nil
		}

		return fmt.Errorf("Error here\n\tlooking for %s, wasn't found in %s", p, rules)
	}
}

func testAccOutscaleOutboundRuleConfig(rInt int) string {
	return fmt.Sprintf(`
	resource "outscale_firewall_rules_set" "web" {
		group_name = "terraform_test_%d"
		group_description = "Used in the terraform acceptance tests"
		tags {
			Name = "tf-acc-test"
		}
	}
	resource "outscale_outbound_rule" "egress_1" {
		ip_permissions = {
			protocol = "tcp"
			from_port = 80
			to_port = 8000
			ip_ranges = ["10.0.0.0/8"]
		}
		group_id = "${outscale_firewall_rules_set.web.id}"
	}`, rInt)
}
