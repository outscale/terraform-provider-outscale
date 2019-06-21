---
layout: "outscale"
page_title: "OUTSCALE: outscale_security_group"
sidebar_current: "docs-outscale-resource-security-group"
description: |-
  Creates a security group.
---

# outscale_security_group

Creates a security group.
This action creates a security group either in the public Cloud or in a specified Net. By default, a default security group for use in the public Cloud and a default security group for use in a Net are created.
When launching a virtual machine (VM), if no security group is explicitly specified, the appropriate default security group is assigned to the VM. Default security groups include a default rule granting VMs network access to each other.
When creating a security group, you specify a name. Two security groups for use in the public Cloud or for use in a Net cannot have the same name.
You can have up to 500 security groups in the public Cloud. You can create up to 500 security groups per Net.
To add or remove rules, use the [CreateSecurityGroupRule](https://docs-beta.outscale.com/#createsecuritygrouprule) method.

## Example Usage

```hcl
resource "outscale_net" "net" {
  ip_range = "10.0.0.0/16"
}


resource "outscale_security_group" "web" {
  security_group_name = "terraform_test_%d"
  description = "Used in the terraform acceptance tests"
  tags = {
    key= "Name"
    value = "tf-acc-test"
  }
  net_id = "${outscale_net.net.id}"
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) A description for the security group, with a maximum length of 255 [ASCII printable characters](https://en.wikipedia.org/wiki/ASCII#Printable_characters).
* `security_group_name` - (Required) The name of the security group. This name must be unique and contain between 1 and 255 ASCII characters. Accented letters are not allowed.
* `net_id` - (Optional) The ID of the Net for the security group.

## Attributes Reference

* `security_group_id` - The ID of the security group.
* `security_group_name` - The name of the security group.
* `inbound_rules` - The inbound rules associated with the security group.
* `outbound_rules` - The outbound rules associated with the security group.
* `account_id` - The account ID of a user that has been granted permission.
* `request_id` - The ID of the request.
* `tags` - One or more tags associated with the security group.

Se detailed information: [Security Group](https://docs-beta.outscale.com/#createsecuritygroup)