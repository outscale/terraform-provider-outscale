---
layout: "outscale"
page_title: "3DS OUTSCALE: outscale_net_attributes"
sidebar_current: "outscale-net-attributes"
description: |-
  [Provides information about Net attributes.]
---

# outscale_net_attributes Data Source

Provides information about Net attributes.
For more information on this resource, see the [User Guide](https://wiki.outscale.net/display/EN/About+VPCs).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#updatenet).

## Example Usage

```hcl

data "outscale_net_attributes" "net_attributes01" {
  net_id = "vpc-12345678"
}


```

## Argument Reference

The following arguments are supported:

* `net_id` - (Required) The ID of the Net.

## Attribute Reference

The following attributes are exported:

* `dhcp_options_set_id` - The ID of the DHCP options set (or `default` if you want to associate the default one).
* `net_id` - The ID of the Net.
