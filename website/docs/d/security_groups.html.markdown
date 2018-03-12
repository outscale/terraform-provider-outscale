---
layout: "outscale"
page_title: "OUTSCALE: outscale_security_groups"
sidebar_current: "docs-outscale-datasource-security-groups"
description: |-
  Describes one or more security groups.
---

# outscale_security_groups

You can specify either the name of the security groups or their IDs.

## Example Usage

```hcl
data "outscale_firewall_rules_sets" "outscale_firewall_rules_sets" {

    filter {

        owner-id ="339215505907"

    }

}


```

## Argument Reference

The following arguments are supported:

* `group_id` - (Optional) The ID of one or more security groups.
* `filter` - (Optional) One or more filters
* `group_name` - (Optional) The name of one or more security groups. This parameter only matches security groups in the public Cloud. To match security groups in a VPC, use the group-name filter instead.


See detailed information in [Outscale Instances](http://docs.outscale.com/api_fcu/operations/Action_DescribeSecurityGroups_get.html#_api_fcu-action_describesecuritygroups_get).

## Filters

Use the Filter.N parameter to filter the described instances on the following properties:

* `description` The description of the security group.
* `group-id` The ID of the security group.
* `group-name` The name of the security group.
* `ip-permission.cidr` A CIDR range that has been granted permission.
* `ip-permission.from-port` The beginning of the port range for the TCP and UDP protocols, or an ICMP type number.
* `ip-permission.group-id` The ID of a security group that has been granted permission.
* `ip-permission.group-name` The name of a security group that has been granted permission.
* `ip-permission.protocol` The IP protocol for the permission (tcp | udp | icmp, or a protocol number).
* `ip-permission.to-port` The end of the port range for the TCP and UDP protocols, or an ICMP code.
* `ip-permission.user-id` The account ID of a user that has been granted permission.
* `owner-id` The account ID of the owner of the security group.
* `tag` The key/value combination of a tag that is assigned to the resource; in the following format: key=value.
* `tag-key` The key of a tag associated with the security group.
* `tag-value` The value of a tag associated with the security group.
* `vpc-id` The ID of the VPC specified when the security group was created.

## Attributes Reference

The following attributes are exported:

* `request_id` - The ID of the request.
* `security_groups_info` - Information about one or more security group.





See detailed information in [Describe Instances](http://docs.outscale.com/api_fcu/operations/Action_DescribeSecurityGroups_get.html#_api_fcu-action_describesecuritygroups_get).
