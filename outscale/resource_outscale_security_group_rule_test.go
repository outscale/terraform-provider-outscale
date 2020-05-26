package outscale

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccOutscaleOAPIOutboundRule(t *testing.T) {
	var group oscgo.SecurityGroup
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
					testAccCheckOutscaleOAPIRuleAttributes("outscale_security_group_rule.outscale_security_group_rule_https", &group, nil, "Inbound"),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPISecurityGroupRuleDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_security_group_rule" {
			continue
		}

		_, resp, err := readSecurityGroups(conn, rs.Primary.ID)
		if err == nil || len(resp.GetSecurityGroups()) > 0 {
			return fmt.Errorf("Outscale Security Group Rule(%s) still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckOutscaleOAPIRuleAttributes(n string, group *oscgo.SecurityGroup, p *oscgo.SecurityGroupRule, ruleType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Security Group Rule Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Security Group Rule is set")
		}

		if p == nil {
			fromPortRange := int64(443)
			toPortRange := int64(443)
			ipProtocol := "tcp"
			p = &oscgo.SecurityGroupRule{
				FromPortRange: &fromPortRange,
				ToPortRange:   &toPortRange,
				IpProtocol:    &ipProtocol,
				IpRanges:      &[]string{"46.231.147.8/32"},
			}
		}

		var matchingRule *oscgo.SecurityGroupRule
		var rules []oscgo.SecurityGroupRule
		if ruleType == "Inbound" {
			rules = group.GetInboundRules()
		} else {
			rules = group.GetOutboundRules()
		}

		if len(rules) == 0 {
			return fmt.Errorf("No IPPerms")
		}

		for _, r := range rules {
			if p.GetToPortRange() != r.GetToPortRange() {
				continue
			}

			if p.GetFromPortRange() != r.GetFromPortRange() {
				continue
			}

			if p.GetIpProtocol() != r.GetIpProtocol() {
				continue
			}

			remaining := len(p.GetIpRanges())
			for _, ip := range p.GetIpRanges() {
				for _, rip := range r.GetIpRanges() {
					if ip == rip {
						remaining--
					}
				}
			}

			if remaining > 0 {
				continue
			}

			remaining = len(p.GetSecurityGroupsMembers())
			for _, ip := range p.GetSecurityGroupsMembers() {
				for _, rip := range r.GetSecurityGroupsMembers() {
					if ip.GetSecurityGroupId() == rip.GetSecurityGroupId() {
						remaining--
					}
				}
			}

			if remaining > 0 {
				continue
			}

			remaining = len(p.GetServiceIds())
			for _, pip := range p.GetServiceIds() {
				for _, rpip := range r.GetServiceIds() {
					if pip == rpip {
						remaining--
					}
				}
			}

			if remaining > 0 {
				continue
			}

			matchingRule = &r
		}

		if matchingRule != nil {
			return nil
		}

		return fmt.Errorf("Security Rules: looking for %+v, wasn't found in %+v", p, rules)
	}
}

func testAccCheckOutscaleOAPIRuleExists(n string, group *oscgo.SecurityGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Security Group is set")
		}

		conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI
		req := oscgo.ReadSecurityGroupsRequest{
			Filters: &oscgo.FiltersSecurityGroup{
				SecurityGroupIds: &[]string{rs.Primary.ID},
			},
		}

		var resp oscgo.ReadSecurityGroupsResponse
		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, _, err = conn.SecurityGroupApi.ReadSecurityGroups(context.Background(), &oscgo.ReadSecurityGroupsOpts{ReadSecurityGroupsRequest: optional.NewInterface(req)})

			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded") {
					fmt.Printf("\n\n[INFO] Request limit exceeded")
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}

			return nil
		})

		if err != nil {
			return err
		}

		if len(resp.GetSecurityGroups()) > 0 && resp.GetSecurityGroups()[0].GetSecurityGroupId() == rs.Primary.ID {
			*group = resp.GetSecurityGroups()[0]
			return nil
		}

		return fmt.Errorf("Security Group not found")
	}
}

func testAccOutscaleOAPISecurityGroupRuleEgressConfig(rInt int) string {
	return fmt.Sprintf(`
		resource "outscale_security_group_rule" "outscale_security_group_rule" {
			flow              = "Inbound"
			security_group_id = "${outscale_security_group.outscale_security_group.security_group_id}"

			to_port_range   = "0"
			ip_protocol     = "tcp"
			ip_range        = "0.0.0.0/0"
		}

		resource "outscale_security_group_rule" "outscale_security_group_rule_https" {
			flow              = "Inbound"
			from_port_range   = 443
			to_port_range     = 443
			ip_protocol       = "tcp"
			ip_range          = "46.231.147.8/32"
			security_group_id = "${outscale_security_group.outscale_security_group.security_group_id}"
		}

		resource "outscale_security_group" "outscale_security_group" {
			description         = "test group"
			security_group_name = "sg1-test-group_test_%d"
		}
	`, rInt)
}
