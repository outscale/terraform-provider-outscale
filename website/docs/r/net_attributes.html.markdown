---
layout: "outscale"
page_title: "3DS OUTSCALE: outscale_net_attributes"
sidebar_current: "outscale-net-attributes"
description: |-
  [Manages a Net attribute.]
---

# outscale_net_attributes Resource

Manages a Net attribute.
For more information on this resource, see the [User Guide](https://wiki.outscale.net/display/EN/About+VPCs).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#updatenet).

## Example Usage

```hcl

#resource "outscale_net" "net01" {
#	ip_range = "10.0.0.0/16"
#}

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
* `state` - The state of the Net (`pending` | `available`).
* `tenancy` - The VM tenancy in a Net (`default` | `dedicated`).

## Import

A Net attribute can be imported using the Net ID. For example:

```

$ terraform import outscale_net_attributes.ImportedNet vpc-12345678

```