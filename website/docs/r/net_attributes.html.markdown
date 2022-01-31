---
layout: "outscale"
page_title: "OUTSCALE: outscale_net_attributes"
sidebar_current: "outscale-net-attributes"
description: |-
  [Manages Net attributes.]
---

# outscale_net_attributes Resource

Manages Net attributes.
For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-DHCP-Options.html).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#updatenet).

## Example Usage

### Required resource

```hcl
resource "outscale_net" "net01" {
	ip_range = "10.0.0.0/16"
}
```

### Associate a DHCP option set to a Net

```hcl
resource "outscale_net_attributes" "net_attributes01" {
	net_id              = outscale_net.net01.net_id
	dhcp_options_set_id = var.dhcp_options_set_id
}
```

## Argument Reference

The following arguments are supported:

* `dhcp_options_set_id` - (Required) The ID of the DHCP options set (or `default` if you want to associate the default one).
* `net_id` - (Required) The ID of the Net.

## Attribute Reference

The following attributes are exported:

* `dhcp_options_set_id` - The ID of the DHCP options set (or `default` if you want to associate the default one).
* `ip_range` - The IP range for the Net, in CIDR notation (for example, 10.0.0.0/16).
* `net_id` - The ID of the Net.
* `state` - The state of the Net (`pending` \| `available`).
* `tags` - One or more tags associated with the Net.
    * `key` - The key of the tag, with a minimum of 1 character.
    * `value` - The value of the tag, between 0 and 255 characters.
* `tenancy` - The VM tenancy in a Net.

## Import

A Net attribute can be imported using the Net ID. For example:

```console

$ terraform import outscale_net_attributes.ImportedNet vpc-12345678

```