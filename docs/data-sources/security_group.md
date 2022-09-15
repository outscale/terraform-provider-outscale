---
layout: "outscale"
page_title: "OUTSCALE: outscale_security_group"
sidebar_current: "outscale-security-group"
description: |-
  [Provides information about a specific security group.]
---

# outscale_security_group Data Source

Provides information about a specific security group.
For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Security-Groups.html).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-securitygroup).

## Example Usage

```hcl
data "outscale_security_group" "security_group01" {
  filter {
    name   = "security_group_ids"
    values = ["sg-12345678"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) A combination of a filter name and one or more filter values. You can specify this argument for as many filter names as you need. The filter name can be any of the following:
    * `account_ids` - (Optional) The account IDs of the owners of the security groups.
    * `descriptions` - (Optional) The descriptions of the security groups.
    * `inbound_rule_account_ids` - (Optional) The account IDs that have been granted permissions.
    * `inbound_rule_from_port_ranges` - (Optional) The beginnings of the port ranges for the TCP and UDP protocols, or the ICMP type numbers.
    * `inbound_rule_ip_ranges` - (Optional) The IP ranges that have been granted permissions, in CIDR notation (for example, 10.0.0.0/24).
    * `inbound_rule_protocols` - (Optional) The IP protocols for the permissions (`tcp` \| `udp` \| `icmp`, or a protocol number, or `-1` for all protocols).
    * `inbound_rule_security_group_ids` - (Optional) The IDs of the security groups that have been granted permissions.
    * `inbound_rule_security_group_names` - (Optional) The names of the security groups that have been granted permissions.
    * `inbound_rule_to_port_ranges` - (Optional) The ends of the port ranges for the TCP and UDP protocols, or the ICMP codes.
    * `net_ids` - (Optional) The IDs of the Nets specified when the security groups were created.
    * `outbound_rule_account_ids` - (Optional) The account IDs that have been granted permissions.
    * `outbound_rule_from_port_ranges` - (Optional) The beginnings of the port ranges for the TCP and UDP protocols, or the ICMP type numbers.
    * `outbound_rule_ip_ranges` - (Optional) The IP ranges that have been granted permissions, in CIDR notation (for example, 10.0.0.0/24).
    * `outbound_rule_protocols` - (Optional) The IP protocols for the permissions (`tcp` \| `udp` \| `icmp`, or a protocol number, or `-1` for all protocols).
    * `outbound_rule_security_group_ids` - (Optional) The IDs of the security groups that have been granted permissions.
    * `outbound_rule_security_group_names` - (Optional) The names of the security groups that have been granted permissions.
    * `outbound_rule_to_port_ranges` - (Optional) The ends of the port ranges for the TCP and UDP protocols, or the ICMP codes.
    * `security_group_ids` - (Optional) The IDs of the security groups.
    * `security_group_names` - (Optional) The names of the security groups.
    * `tag_keys` - (Optional) The keys of the tags associated with the security groups.
    * `tag_values` - (Optional) The values of the tags associated with the security groups.
    * `tags` - (Optional) The key/value combination of the tags associated with the security groups, in the following format: &quot;Filters&quot;:{&quot;Tags&quot;:[&quot;TAGKEY=TAGVALUE&quot;]}.

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
