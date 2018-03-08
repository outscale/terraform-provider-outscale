---
layout: "outscale"
page_title: "OUTSCALE: outscale_inbound_rule"
sidebar_current: "docs-outscale-resource-inbound_rule"
description: |-
Adds one or more ingress rules to a security group.
The modifications are effective at instances level as quickly as possible, but a small delay may occur.
In the public Cloud, this action allows one or more CIDR IP address ranges to access a security group for your account, or allows one or more security groups (source groups) to access a security group for your own Outscale account or another one.
In a VPC, this action allows one or more CIDR IP address ranges to access a security group for your VPC, or allows one or more other security groups (source groups) to access a security group for your VPC. All the security groups must be for the same VPC.
To create a rule with a specific IP protocol and a specific port range, we recommand to use a set of IP permissions. We also recommand to specify the protocol in a set of IP permissions.
---

* NOTE - 
By default, traffic between two security groups is allowed through both public and private IP addresses. To restrict it to private IP addresses only, contact our Support team: support@outscale.com.

## Example Usage

```hcl
resource "outscale_firewall_rules_set" "web" {
		group_name = "terraform_test_%d"
		group_description = "Used in the terraform acceptance tests"
					tag = {
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
	}
```

## Argument Reference

The following arguments are supported:

* `cidr_ip` - The CIDR IP address range.

* `from_port` - The beginning of the port range for the TCP and UDP protocols, or an ICMP type number.

* `group_id` - The ID of the security group (mandatory for a non-default VPC).

* `group_name` - The name of the security group.

* `ip_permissions` - Describes a security group rule.

* `ip_protocol` - The IP protocol name or number.

* `source_security_group_name` - The name of the source security group (cannot be combined with the FromPort, ToPort, CidrIp and IpProtocol parameters).

* `source_security_group_owner_id` - The Outscale account ID of the owner of the source security group, creating rules that grant full ICMP, UDP, and TCP access (cannot be combined with the FromPort, ToPort, CidrIp and IpProtocol parameters).

* `to_port` - The end of port range for the TCP and UDP protocols, or an ICMP type number.


The IP Permissions block has the following attributes:

* `from_port` - The beginning of the port range for the TCP and UDP protocols, or an ICMP type number.

* `groups` - One or more security groups and account ID pairs.

* `ip_protocol` - The IP protocol name or number.

* `ip_ranges` - One or more IP ranges.

* `prefix_list_ids` - One or more prefix list IDs to allow traffic from a VPC to access the corresponding Outscale services. For more information, see DescribePrefixLists

* `to_port` - The end of port range for the TCP and UDP protocols, or an ICMP type number.



See detailed information in [Authorize Security Group Ingress](http://docs.outscale.com/api_fcu/operations/Action_AuthorizeSecurityGroupIngress_get.html#_api_fcu-action_authorizesecuritygroupingress_get).

See detailed information in [Revoke Security Group Ingress](http://docs.outscale.com/api_fcu/operations/Action_RevokeSecurityGroupIngress_get.html#_api_fcu-action_revokesecuritygroupingress_get).