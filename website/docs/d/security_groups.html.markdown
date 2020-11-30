---
layout: "outscale"
page_title: "3DS OUTSCALE: outscale_security_group"
sidebar_current: "outscale-security-group"
description: |-
  [Provides information about security groups.]
---

# outscale_security_group Data Source

Provides information about security groups.
For more information on this resource, see the [User Guide](https://wiki.outscale.net/display/EN/About+Security+Groups).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-securitygroup).

## Example Usage

```hcl

data "outscale_security_groups" "security_groups01" {
  filter {
    name   = "security_group_ids"
    values = ["sg-12345678", "sg-12345679"]
  }
}


```

## Argument Reference

The following arguments are supported:

* `filter` - One or more filters.
  * `account_ids` - (Optional) The account IDs of the owners of the security groups.
  * `security_group_ids` - (Optional) The IDs of the security groups.
  * `security_group_names` - (Optional) The names of the security groups.
  * `tag_keys` - (Optional) The keys of the tags associated with the security groups.
  * `tag_values` - (Optional) The values of the tags associated with the security groups.
  * `tags` - (Optional) The key/value combination of the tags associated with the security groups, in the following format: "Filters":{"Tags":["TAGKEY=TAGVALUE"]}.

## Attribute Reference

The following attributes are exported:

* `security_groups` - Information about one or more security groups.
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
      * `service_ids` - One or more service IDs to allow traffic from a Net to access the corresponding 3DS OUTSCALE services. For more information, see [ReadNetAccessPointServices](https://docs.outscale.com/api#readnetaccesspointservices).
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
      * `service_ids` - One or more service IDs to allow traffic from a Net to access the corresponding 3DS OUTSCALE services. For more information, see [ReadNetAccessPointServices](https://docs.outscale.com/api#readnetaccesspointservices).
      * `to_port_range` - The end of the port range for the TCP and UDP protocols, or an ICMP type number.
  * `security_group_id` - The ID of the security group.
  * `security_group_name` - The name of the security group.
  * `tags` - One or more tags associated with the security group.
      * `key` - The key of the tag, with a minimum of 1 character.
      * `value` - The value of the tag, between 0 and 255 characters.
