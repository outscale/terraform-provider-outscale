---
layout: "outscale"
page_title: "3DS OUTSCALE: outscale_security_group_rule"
sidebar_current: "outscale-security-group-rule"
description: |-
  [Manages a security group rule.]
---

# outscale_security_group_rule Resource

Manages a security group rule.
For more information on this resource, see the [User Guide](https://wiki.outscale.net/display/EN/About+Security+Group+Rules).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-securitygrouprule).

## Example Usage

```hcl

# Set rule from IP range

resource "outscale_security_group" "security_group01" {
  description         = "Terraform security group for sg rule"
  security_group_name = "terraform-security-group-test-01"
}

resource "outscale_security_group_rule" "security_group_rule01" {
  flow              = "Inbound"
  security_group_id = outscale_security_group.security_group01.security_group_id
  from_port_range   = "80"
  to_port_range     = "80"
  ip_protocol       = "tcp"
  ip_range          = "10.0.0.0/16"
}

# Set rule from another security group

resource "outscale_security_group" "security_group02" {
    description         = "Terraform security group for sg rule"
    security_group_name = "terraform-security-group-test-02"
}

resource "outscale_security_group" "security_group03" {
    description         = "Terraform security group for sg rule"
    security_group_name = "terraform-security-group-test-03"
}

resource "outscale_security_group_rule" "security_group_rule02" 
    flow              = "Inbound"
    security_group_id = outscale_security_group.security_group02.security_group_id
    rules {
     from_port_range   = "22"
     to_port_range     = "22"
     ip_protocol       = "tcp"
     security_groups_members {
        account_id          =  "012345678910"
        security_group_name = outscale_security_group.security_group03.security_group_name
       }
     }


```

## Argument Reference

The following arguments are supported:

* `flow` - (Required) The direction of the flow: `Inbound` or `Outbound`. You can specify `Outbound` for Nets only.
* `from_port_range` - (Optional) The beginning of the port range for the TCP and UDP protocols, or an ICMP type number.
* `ip_protocol` - (Optional) The IP protocol name (`tcp`, `udp`, `icmp`) or protocol number. By default, `-1`, which means all protocols.
* `ip_range` - (Optional) The IP range for the security group rule, in CIDR notation (for example, 10.0.0.0/16).
* `rules` - (Optional) Information about the security group rule to create.
  * `from_port_range` - (Optional) The beginning of the port range for the TCP and UDP protocols, or an ICMP type number.
  * `ip_protocol` - (Optional) The IP protocol name (`tcp`, `udp`, `icmp`) or protocol number. By default, `-1`, which means all protocols.
  * `ip_ranges` - (Optional) One or more IP ranges for the security group rules, in CIDR notation (for example, 10.0.0.0/16).
  * `security_groups_members` - (Optional) Information about one or more members of a security group.
     * `account_id` - (Optional) The account ID of a user.
     * `security_group_id` - (Required) The ID of the security group.
     * `security_group_name` - (Optional) The name of the security group.
  * `service_ids` - (Optional) One or more service IDs to allow traffic from a Net to access the corresponding OUTSCALE services. For more information, see [ReadNetAccessPointServices](https://docs.outscale.com/api#readnetaccesspointservices).
  * `to_port_range` - (Optional) The end of the port range for the TCP and UDP protocols, or an ICMP type number.
* `security_group_account_id_to_link` - (Optional) The account ID of the owner of the security group for which you want to create a rule.
* `security_group_id` - (Required) The ID of the security group for which you want to create a rule.
* `security_group_name_to_link` - (Optional) The ID of the source security group. If you are in the Public Cloud, you can also specify the name of the source security group.
* `to_port_range` - (Optional) The end of the port range for the TCP and UDP protocols, or an ICMP type number.

## Attribute Reference

The following attributes are exported:

* `security_group` - Information about the security group.
  * `account_id` - The account ID of a user that has been granted permission.
  * `description` - The description of the security group.
  * `inbound_rules` - The inbound rules associated with the security group.
     * `from_port_range` - The beginning of the port range for the TCP and UDP protocols, or an ICMP type number.
     * `ip_protocol` - The IP protocol name (`tcp`, `udp`, `icmp`) or protocol number. By default, `-1`, which means all protocols.
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
     * `ip_protocol` - The IP protocol name (`tcp`, `udp`, `icmp`) or protocol number. By default, `-1`, which means all protocols.
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
