---
layout: "outscale"
page_title: "OUTSCALE: outscale_firewall_rules_set"
sidebar_current: "docs-outscale-resource-firewall_rules_set"
description: |-
Creates a security group.
This action creates a security group either in the public Cloud or in a specified VPC. By default, a default security group for use in the public Cloud and a default security group for use in a VPC are created.
When launching an instance, if no security group is explicitly specified, the appropriate default security group is assigned to the instance. Default security groups include a default rule granting instances network access to each other.
When creating a security group, you specify a name. Two security groups for use in the public Cloud or for use in a VPC cannot have the same name.
You can have up to 500 security groups in the public Cloud. You can create up to 500 security groups per VPC.
To add or remove rules, use the (AuthorizeSecurityGroupIngress)[http://docs.outscale.com/api_fcu/operations/Action_AuthorizeSecurityGroupIngress_post.html#_api_fcu-action_authorizesecuritygroupingress_post], (AuthorizeSecurityGroupEgress)[http://docs.outscale.com/api_fcu/operations/Action_AuthorizeSecurityGroupEgress_post.html#_api_fcu-action_authorizesecuritygroupegress_post], (RevokeSecurityGroupIngress)[http://docs.outscale.com/api_fcu/operations/Action_RevokeSecurityGroupIngress_post.html#_api_fcu-action_revokesecuritygroupingress_post] or (RevokeSecurityGroupEgress)[http://docs.outscale.com/api_fcu/operations/Action_RevokeSecurityGroupEgress_post.html#_api_fcu-action_revokesecuritygroupegress_post] methods.
---

## Example Usage

```hcl
resource "aws_security_group" "allow_all" {
  name        = "allow_all"
  description = "Allow all inbound traffic"
  vpc_id      = "${aws_vpc.main.id}"

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port       = 0
    to_port         = 0
    protocol        = "-1"
    cidr_blocks     = ["0.0.0.0/0"]
    prefix_list_ids = ["pl-12c4e678"]
  }
}
```

Basic usage with tags:

```hcl
resource "aws_security_group" "allow_all" {
  name        = "allow_all"
  description = "Allow all inbound traffic"

  ingress {
    from_port   = 0
    to_port     = 65535
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    Name = "allow_all"
  }
}
```

## Argument Reference

The following arguments are supported:

* `group_name` - (Required) The name of the security group (between 1 and 255 of the following characters: ASCII for the public Cloud, and a-z, A-Z, 0-9, spaces or ._-:/()#,@[]+=&;{}!$* for a VPC).
* `group_description` - (Required) A description for the security group (between 1 and 255 of the following characters: ASCII for the public Cloud, and a-z, A-Z, 0-9, spaces or ._-:/()#,@[]+=&;{}!$* for a VPC).
* `ip_permissions` - The inbound rules associated with the security group.
* `ip_permissions_egress` - The outbound rules associated with the security group.
* `vpc_id` - The ID of the VPC.
* `tags` - (Optional) A mapping of tags to assign to the resource.
* `tag_set` - One or more tags associated with the security group.
* `ownerId` - The ID of the VPC for the security group.


See detailed information in [Authorize Security Group Egress](http://docs.outscale.com/api_fcu/operations/Action_AuthorizeSecurityGroupEgress_get.html#_api_fcu-action_authorizesecuritygroupegress_get).
See detailed information in [Authorize Security Group Ingress](http://docs.outscale.com/api_fcu/operations/Action_AuthorizeSecurityGroupIngress_get.html#_api_fcu-action_authorizesecuritygroupingress_get).
See detailed information in [Create Security Group](http://docs.outscale.com/api_fcu/operations/Action_CreateSecurityGroup_get.html#_api_fcu-action_createsecuritygroup_get).
See detailed information in [Delete Security Group](http://docs.outscale.com/api_fcu/operations/Action_DeleteSecurityGroup_get.html#_api_fcu-action_deletesecuritygroup_get).
See detailed information in [Describe Security Groups](http://docs.outscale.com/api_fcu/operations/Action_DescribeSecurityGroups_get.html#_api_fcu-action_describesecuritygroups_get).
See detailed information in [Revoke Security Group Egress](http://docs.outscale.com/api_fcu/operations/Action_RevokeSecurityGroupEgress_get.html#_api_fcu-action_revokesecuritygroupegress_get).
See detailed information in [Revoke Security Group Ingress](http://docs.outscale.com/api_fcu/operations/Action_RevokeSecurityGroupIngress_get.html#_api_fcu-action_revokesecuritygroupingress_get).