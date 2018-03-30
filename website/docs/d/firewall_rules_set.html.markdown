---
layout: "outscale"
page_title: "OUTSCALE: outscale_firewall_rules_set"
sidebar_current: "docs-outscale-datasource-firewall_rules_set"
description: |-
  Describes one or more security groups.
---

# outscale_firewall_rules_set

	Describes one or more security groups. You can specify either the name of the security groups or their IDs.

## Example Usage

```hcl
	data "outscale_firewall_rules_set" "by_id" {
		group_id = ["${outscale_firewall_rules_set.test.id}", "${outscale_firewall_rules_set.test2.id}", "${outscale_firewall_rules_set.test3.id}"]
	}
	data "outscale_firewall_rules_set" "by_filter" {
		filter {
			name = "group-name"
			values = ["${outscale_firewall_rules_set.test.group_name}"]
		}
	}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Optional) The ID of one or more security groups..
* `group_name` - (Optional) The name of one or more security groups. This parameter only matches security groups in the public Cloud. To match security groups in a VPC, use the group-name filter instead.

## Filters

You can use the Filter.N parameter to filter the security groups on the following properties:

* `description` -  The description of the security group.
* `group-id` -  The ID of the security group.
* `group-name` -  The name of the security group.
* `ip-permission.cidr` -  A CIDR range that has been granted permission.
* `ip-permission.from-port` -  The beginning of the port range for the TCP and UDP protocols, or an ICMP type number.
* `ip-permission.group-id` -  The ID of a security group that has been granted permission.
* `ip-permission.group-name` -  The name of a security group that has been granted permission.
* `ip-permission.protocol` -  The IP protocol for the permission (tcp | udp | icmp, or a protocol number).
* `ip-permission.to-port` -  The end of the port range for the TCP and UDP protocols, or an ICMP code.
* `ip-permission.user-id` -  The account ID of a user that has been granted permission.
* `owner-id` -  The account ID of the owner of the security group.
* `tag` -  The key/value combination of a tag that is assigned to the resource; in the following format` -  key=value.
* `tag-key` -  The key of a tag associated with the security group.
* `tag-value` -  The value of a tag associated with the security group.
* `vpc-id` -  The ID of the VPC specified when the security group was created.

## Attributes Reference

* `groupDescription` - 	A description of the security group.	false	string
* `groupId` - 	The ID of the security group.	false	string
* `groupName` - 	The name of the security group.	false	string
* `ipPermissions` - 	The inbound rules associated with the security group.	false	IpPermission
* `ipPermissionsEgress` - 	The outbound rules associated with the security group.	false	IpPermission
* `ownerId` - 	The account ID of the owner of the security group.	false	string
* `tagSet` - 	One or more tags associated with the security group.	false	Tag
* `vpcId` - 	The ID of the VPC for the security group.	false	string
