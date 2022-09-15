---
layout: "outscale"
page_title: "OUTSCALE: outscale_security_group"
sidebar_current: "outscale-security-group"
description: |-
  [Manages a security group.]
---

# outscale_security_group Resource

Manages a security group.
For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Security-Groups.html).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-securitygroup).

## Example Usage

### Optional resource

```hcl
resource "outscale_net" "net01" {
	ip_range = "10.0.0.0/16"
}
```

### Create a security group for a Net

```hcl
resource "outscale_security_group" "security_group01" {
	description         = "Terraform security group test"
	security_group_name = "terraform-security-group"
	net_id              = outscale_net.net01.net_id
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Required) A description for the security group, with a maximum length of 255 [ASCII printable characters](https://en.wikipedia.org/wiki/ASCII#Printable_characters).
* `net_id` - (Optional) The ID of the Net for the security group.
* `security_group_name` - (Required) The name of the security group.<br />
This name must not start with `sg-`.</br>
This name must be unique and contain between 1 and 255 ASCII characters. Accented letters are not allowed.
* `tags` - (Optional) A tag to add to this resource. You can specify this argument several times.
    * `key` - (Required) The key of the tag, with a minimum of 1 character.
    * `value` - (Required) The value of the tag, between 0 and 255 characters.

## Attribute Reference

The following attributes are exported:

* `account_id` - The account ID of a user that has been granted permission.
* `description` - The description of the security group.
* `inbound_rules` - The inbound rules associated with the security group.
    * `from_port_range` - The beginning of the port range for the TCP and UDP protocols, or an ICMP type number.
    * `ip_protocol` - The IP protocol name (`tcp`, `udp`, `icmp`, or `-1` for all protocols). By default, `-1`. In a Net, this can also be an IP protocol number. For more information, see the [IANA.org website](https://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml).
    * `ip_ranges` - One or more IP ranges for the security group rules, in CIDR notation (for example, 10.0.0.0/16).
    * `security_groups_members` - Information about one or more members of a security group.
        * `account_id` - The account ID of a user.
        * `security_group_id` - The ID of the security group.
        * `security_group_name` - The name of the security group.
    * `service_ids` - One or more service IDs to allow traffic from a Net to access the corresponding OUTSCALE services. For more information, see [ReadNetAccessPointServices](https://docs.outscale.com/api#readnetaccesspointservices).
    * `to_port_range` - The end of the port range for the TCP and UDP protocols, or an ICMP type number.
* `net_id` - The ID of the Net for the security group.
* `outbound_rules` - The outbound rules associated with the security group.
    * `from_port_range` - The beginning of the port range for the TCP and UDP protocols, or an ICMP type number.
    * `ip_protocol` - The IP protocol name (`tcp`, `udp`, `icmp`, or `-1` for all protocols). By default, `-1`. In a Net, this can also be an IP protocol number. For more information, see the [IANA.org website](https://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml).
    * `ip_ranges` - One or more IP ranges for the security group rules, in CIDR notation (for example, 10.0.0.0/16).
    * `security_groups_members` - Information about one or more members of a security group.
        * `account_id` - The account ID of a user.
        * `security_group_id` - The ID of the security group.
        * `security_group_name` - The name of the security group.
    * `service_ids` - One or more service IDs to allow traffic from a Net to access the corresponding OUTSCALE services. For more information, see [ReadNetAccessPointServices](https://docs.outscale.com/api#readnetaccesspointservices).
    * `to_port_range` - The end of the port range for the TCP and UDP protocols, or an ICMP type number.
* `security_group_id` - The ID of the security group.
* `security_group_name` - The name of the security group.
* `tags` - One or more tags associated with the security group.
    * `key` - The key of the tag, with a minimum of 1 character.
    * `value` - The value of the tag, between 0 and 255 characters.

## Import

A security group can be imported using its ID. For example:

```console

$ terraform import outscale_security_group.ImportedSecurityGroup sg-87654321

```