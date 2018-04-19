---
layout: "outscale"
page_title: "OUTSCALE: outscale_firewall_rules_set"
sidebar_current: "docs-outscale-datasource-firewall-rules-set"
description: |-
Describes one or more security groups.


---

# outscale_firewall_rules_set

Describes one or more security groups.
You can specify either the name of the security groups or their IDs.


## Example Usage

```hcl
data "outscale_firewall_rules_set" "by_id" {
	group_id = "${outscale_firewall_rules_set.test.id}"
}
data "outscale_firewall_rules_set" "by_filter" {
	filter {
		name = "group-name"
		values = ["${outscale_firewall_rules_set.test.group_name}"]
	}
}`, rInt, rInt)
```

## Argument Reference

The following arguments are supported:

* `group_id.N` (Optional) The ID of one or more security groups.
* `group_name.N` - (Optional) The name of one or more security groups. This parameter only matches security groups in the public Cloud. To match security groups in a VPC, use the group-name filter instead.
* `filter.N` - (Optional) One or more filters.


See detailed information in [Outscale firewallRulesSets](http://docs.outscale.com/api_fcu/operations/Action_DescribeSecurityGroups_get.html#_api_fcu-action_describesecuritygroups_get).

## Filters

Use the Filter.N parameter to filter the described instances on the following properties:

* `description`  The description of the security group.
* `group-id`  The ID of the security group.
* `group-name`  The name of the security group.
* `ip-permission.cidr	` A CIDR range that has been granted permission. hostnames.
* `ip-permission.group-id`  The ID of a security group that has been granted permission.
* `ip-permission.group-name`  The name of a security group that has been granted permission.
* `ip-permission.protocol`  The IP protocol for the permission (tcp | udp | icmp, or a protocol number).
* `ip-permission-to-port`  The end of the port range for the TCP and UDP protocols, or an ICMP code.
* `ip-permission.user-id`  The account ID of a user that has been granted permission.
* `owner-id`  The account ID of the owner of the security group.
* `tag:key=value`  The key/value combination of a tag that is assigned to the resource.
* `tag-key`  The key of a tag associated with the security group.
* `tag-value` The value of a tag associated with the security group.
* `vpc-id`  The ID of the VPC specified when the security group was created.



## Attributes Reference

The following attributes are exported:

* `group_description` - A description of the security group.
* `group_id` - The ID of the security group.
* `group_name` - The name of the security group.
* `ip_permissions.N` - The inbound rules associated with the security group.
* `ip_permissions_egress.N` - The outbound rules associated with the security group
* `owner_id` - The account ID of the owner of the security group.
* `tag_set.N` - One or more tags associated with the security group.
* `vpc_id` - The ID of the VPC for the security group.

See detailed information in [Describe firewallRulesSets](http://docs.outscale.com/api_fcu/operations/Action_DescribeSecurityGroups_get.html#_api_fcu-action_describesecuritygroups_get).
