---
layout: "outscale"
page_title: "OUTSCALE: outscale_security_group"
subcategory: "OUTSCALE API"
sidebar_current: "outscale-security-group"
description: |-
  [Manages a security group.]
---

# outscale_security_group Resource

Manages a security group.

Security groups you create to use in a Net contain a default outbound rule that allows all outbound flows.

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
	description         = "Terraform security group"
	security_group_name = "terraform-security-group"
	net_id              = outscale_net.net01.net_id
}
```

### Create a security group for a Net without the default outbound rule 

```hcl
resource "outscale_security_group" "security_group02" {
    remove_default_outbound_rule = true
    description                  = "Terraform security group without outbound rule"
    security_group_name          = "terraform-security-group-empty"
    net_id                       = outscale_net.net01.net_id
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Required) A description for the security group.<br />
This description can contain between 1 and 255 characters. Allowed characters are `a-z`, `A-Z`, `0-9`, accented letters, spaces, and `_.-:/()#,@[]+=&;{}!$*`.
* `net_id` - (Optional) The ID of the Net for the security group.
* `remove_default_outbound_rule` - (Optional) (Net only) By default or if set to false, the security group is created with a default outbound rule allowing all outbound flows. If set to true, the security group is created without a default outbound rule. For an existing security group, setting this parameter to true deletes the security group and creates a new one.
* `security_group_name` - (Optional) A name for the security group.<br />
This name must be unique and contain between 1 and 255 characters. It must not start with `sg-`. Allowed characters are `a-z`, `A-Z`, `0-9`, spaces, and `_.-:/()#,@[]+=&;{}!$*`.<br />
If not specified, the security group name is randomly generated.
* `tags` - (Optional) A tag to add to this resource. You can specify this argument several times.
    * `key` - (Required) The key of the tag, with a minimum of 1 character.
    * `value` - (Required) The value of the tag, between 0 and 255 characters.

## Attribute Reference

The following attributes are exported:

* `account_id` - The account ID that has been granted permission.
* `description` - The description of the security group.
* `inbound_rules` - The inbound rules associated with the security group.
    * `from_port_range` - The beginning of the port range for the TCP and UDP protocols, or an ICMP type number.
    * `ip_protocol` - The IP protocol name (`tcp`, `udp`, `icmp`, or `-1` for all protocols). By default, `-1`. In a Net, this can also be an IP protocol number. For more information, see the [IANA.org website](https://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml).
    * `ip_ranges` - One or more IP ranges for the security group rules, in CIDR notation (for example, `10.0.0.0/16`).
    * `security_groups_members` - Information about one or more source or destination security groups.
        * `account_id` - The account ID that owns the source or destination security group.
        * `security_group_id` - The ID of a source or destination security group that you want to link to the security group of the rule.
        * `security_group_name` - (Public Cloud only) The name of a source or destination security group that you want to link to the security group of the rule.
    * `service_ids` - One or more service IDs to allow traffic from a Net to access the corresponding OUTSCALE services. For more information, see [ReadNetAccessPointServices](https://docs.outscale.com/api#readnetaccesspointservices).
    * `to_port_range` - The end of the port range for the TCP and UDP protocols, or an ICMP code number.
* `net_id` - The ID of the Net for the security group.
* `outbound_rules` - The outbound rules associated with the security group.
    * `from_port_range` - The beginning of the port range for the TCP and UDP protocols, or an ICMP type number.
    * `ip_protocol` - The IP protocol name (`tcp`, `udp`, `icmp`, or `-1` for all protocols). By default, `-1`. In a Net, this can also be an IP protocol number. For more information, see the [IANA.org website](https://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml).
    * `ip_ranges` - One or more IP ranges for the security group rules, in CIDR notation (for example, `10.0.0.0/16`).
    * `security_groups_members` - Information about one or more source or destination security groups.
        * `account_id` - The account ID that owns the source or destination security group.
        * `security_group_id` - The ID of a source or destination security group that you want to link to the security group of the rule.
        * `security_group_name` - (Public Cloud only) The name of a source or destination security group that you want to link to the security group of the rule.
    * `service_ids` - One or more service IDs to allow traffic from a Net to access the corresponding OUTSCALE services. For more information, see [ReadNetAccessPointServices](https://docs.outscale.com/api#readnetaccesspointservices).
    * `to_port_range` - The end of the port range for the TCP and UDP protocols, or an ICMP code number.
* `remove_default_outbound_rule` - If false, the security group is created with a default outbound rule allowing all outbound flows. If true, the security group is created without a default outbound rule.
* `security_group_id` - The ID of the security group.
* `security_group_name` - The name of the security group.
* `tags` - One or more tags associated with the security group.
    * `key` - The key of the tag, with a minimum of 1 character.
    * `value` - The value of the tag, between 0 and 255 characters.

## Timeouts

The `timeouts` block enables you to configure [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - Defaults to 10 minutes.
* `read` - Defaults to 5 minutes.
* `update` - Defaults to 10 minutes.
* `delete` - Defaults to 5 minutes.

## Import

A security group can be imported using its ID. For example:

```console

$ terraform import outscale_security_group.ImportedSecurityGroup sg-87654321

```