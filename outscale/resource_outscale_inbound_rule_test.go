package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func TestAccOutscaleInboundRule_Ingress_VPC(t *testing.T) {
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
		CheckDestroy: testAccCheckOutscaleSGRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleInboundRuleIngressConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleSecurityGroupRuleExists("outscale_firewall_rules_set.web", &group),
					testAccCheckOutscaleSecurityGroupRuleAttributes("outscale_inbound_rule.ingress_1", &group, nil, "ingress"),
					resource.TestCheckResourceAttr(
						"outscale_inbound_rule.ingress_1", "ip_permissions.0.from_port", "80"),
					testRuleCount,
				),
			},
		},
	})
}

func TestAccOutscaleInboundRule_Ingress_Protocol(t *testing.T) {
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
		CheckDestroy: testAccCheckOutscaleSGRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleInboundRule_protocolConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleSecurityGroupRuleExists("outscale_firewall_rules_set.web", &group),
					testAccCheckOutscaleSecurityGroupRuleAttributes("outscale_inbound_rule.ingress_1", &group, nil, "ingress"),
					resource.TestCheckResourceAttr(
						"outscale_inbound_rule.ingress_1", "ip_permissions.0.ip_protocol", "6"),
					testRuleCount,
				),
			},
		},
	})
}

func TestAccOutscaleInboundRule_Ingress_Classic(t *testing.T) {
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
		CheckDestroy: testAccCheckOutscaleSGRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleInboundRuleIngressClassicConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleSecurityGroupRuleExists("outscale_firewall_rules_set.web", &group),
					testAccCheckOutscaleSecurityGroupRuleAttributes("outscale_inbound_rule.ingress_1", &group, nil, "ingress"),
					resource.TestCheckResourceAttr(
						"outscale_inbound_rule.ingress_1", "ip_permissions.0.from_port", "80"),
					testRuleCount,
				),
			},
		},
	})
}

func TestAccOutscaleInboundRule_MultiIngress(t *testing.T) {
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
		CheckDestroy: testAccCheckOutscaleSGRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleInboundRuleConfigMultiIngress,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleSecurityGroupRuleExists("outscale_firewall_rules_set.web", &group),
					testMultiRuleCount,
				),
			},
		},
	})
}

func testAccOutscaleInboundRuleIngressConfig(rInt int) string {
	return fmt.Sprintf(`
	resource "outscale_firewall_rules_set" "web" {
		group_name = "terraform_test_%d"
		group_description = "Used in the terraform acceptance tests"
		vpc_id = "vpc-e9d09d63"
					tags {
									Name = "tf-acc-test"
					}
	}
	resource "outscale_inbound_rule" "ingress_1" {
		ip_permissions = {
			ip_protocol = "tcp"
			from_port = 80
			to_port = 8000
			ip_ranges = {
				cidr_ip = "10.0.0.0/8"
			}
		}
		group_id = "${outscale_firewall_rules_set.web.id}"
		source_security_group_name = "${outscale_firewall_rules_set.web.group_name}"
	}`, rInt)
}

func testAccOutscaleInboundRule_protocolConfig(rInt int) string {
	return fmt.Sprintf(`
resource "outscale_firewall_rules_set" "web" {
		group_name = "terraform_test_%d"
		group_description = "Used in the terraform acceptance tests"
		vpc_id = "vpc-e9d09d63"
					tags {
									Name = "tf-acc-test"
					}
	}
resource "outscale_inbound_rule" "ingress_1" {
	ip_permissions = {
			ip_protocol    = "6"
			from_port   = 80
			to_port     = 8000
			ip_ranges = {
				cidr_ip = "10.0.0.0/8"
			}
		}
		group_id = "${outscale_firewall_rules_set.web.id}"
		source_security_group_name = "${outscale_firewall_rules_set.web.group_name}"
}
`, rInt)
}

func testAccOutscaleInboundRuleIngressClassicConfig(rInt int) string {
	return fmt.Sprintf(`
	resource "outscale_firewall_rules_set" "web" {
		group_name = "terraform_test_%d"
		group_description = "Used in the terraform acceptance tests"
		vpc_id = "vpc-e9d09d63"
					tags {
									Name = "tf-acc-test"
					}
	}
	resource "outscale_inbound_rule" "ingress_1" {
		ip_permissions = {
			ip_protocol = "tcp"
			from_port = 80
			to_port = 8000
			ip_ranges = {
				cidr_ip = "10.0.0.0/8"
			}
		}
		group_id = "${outscale_firewall_rules_set.web.id}"
		source_security_group_name = "${outscale_firewall_rules_set.web.group_name}"
	}`, rInt)
}

const testAccOutscaleInboundRuleConfigMultiIngress = `
resource "outscale_firewall_rules_set" "web" {
  group_name = "terraform_acceptance_test_example_2"
	group_description = "Used in the terraform acceptance tests"
	vpc_id = "vpc-e9d09d63"
}

resource "outscale_inbound_rule" "ingress_1" {
	ip_permissions = {
			ip_protocol = "tcp"
			from_port = 22
			to_port = 22
			ip_ranges = {
					cidr_ip = "10.0.0.0/8"
			}
	}
  group_id = "${outscale_firewall_rules_set.web.id}"
}
resource "outscale_inbound_rule" "ingress_2" {
  
	ip_permissions = {
			ip_protocol = "tcp"
  		from_port = 80
			to_port = 8000
			ip_ranges = {
					cidr_ip = "10.0.0.0/8"
			}
		}
  group_id = "${outscale_firewall_rules_set.web.id}"
}
`
