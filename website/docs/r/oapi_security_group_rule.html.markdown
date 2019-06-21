---
layout: "outscale"
page_title: "OUTSCALE: outscale_security_group"
sidebar_current: "docs-outscale-resource-security-group-rule"
description: |-
  Creates a security group rule.
---

# outscale_security_group_rule

Configures the rules for a security group.
The modifications are effective at virtual machine (VM) level as quickly as possible, but a small delay may occur.

You can add one or more egress rules to a security group for use with a Net.
It allows VMs to send traffic to either one or more destination IP address ranges or destination security groups for the same Net.
We recommend using a set of IP permissions to authorize outbound access to a destination security group. We also recommended this method to create a rule with a specific IP protocol and a specific port range. In a set of IP permissions, we recommend to specify the the protocol.

You can also add one or more ingress rules to a security group.
In the public Cloud, this action allows one or more IP address ranges to access a security group for your account, or allows one or more security groups (source groups) to access a security group for your own Outscale account or another one.
In a Net, this action allows one or more IP address ranges to access a security group for your Net, or allows one or more other security groups (source groups) to access a security group for your Net. All the security groups must be for the same Net.


## Example Usage

```hcl
resource "outscale_security_group_rule" "outscale_security_group_rule" {
	flow              = "Inbound"
	security_group_id = "${outscale_security_group.outscale_security_group.security_group_id}"

	from_port_range = "0"
	to_port_range = "0"
	ip_protocol = "tcp"
	ip_range = "0.0.0.0/0"
}

resource "outscale_security_group_rule" "outscale_security_group_rule_https" {
	flow = "Inbound"
	from_port_range = 443
	to_port_range = 443
	ip_protocol = "tcp"
	ip_range = "46.231.147.8/32"
	security_group_id = "${outscale_security_group.outscale_security_group.security_group_id}"
	}

resource "outscale_security_group" "outscale_security_group" {
	description         = "test group"
	security_group_name = "sg1-test-group_test_%d"
}
```

## Argument Reference

The following arguments are supported:

* `flow` - (Required) The direction of the flow: Inbound or Outbound. You can specify Outbound for Nets only.
* `ip_range` - (Optional) The IP range for the security group rule, in CIDR notation (for example, 10.0.0.0/16).
* `from_port_range` - (Optional) The beginning of the port range for the TCP and UDP protocols, or an ICMP type number.
* `security_group_id` - (Optional) The ID of the security group for which you want to create a rule.
* `security_group_name` - (Optional) The name of the security group for which you want to create a rule. 
* `ip_protocol` - (Optional) The IP protocol name (tcp, udp, icmp) or protocol number. By default, -1, which means all protocols.
* `security_group_name_to_link` - (Optional) The name of the security group for which you want to create a rule.
* `security_group_account_id_to_link` - (Optional) 	The account ID of the owner of the security group for which you want to create a rule.
* `to_port_range` - (Optional) The end of the port range for the TCP and UDP protocols, or an ICMP type number.
* `rules` - (Optional) Information about the security group rule to create.
* `inbound_rules`- (Optional) Information about the security group inbound rule to create.
* `outbound_rules` - (Optional) Information about the security group outbound rule to create.



## Attributes Reference

The following attributes are supported:

* `request_id` - The ID of the request.
* `net_id` - The VPC if the Security Group

Se detailed information: [CreateSecurityGroupRule](https://docs-beta.outscale.com/#createsecuritygrouprule) method.