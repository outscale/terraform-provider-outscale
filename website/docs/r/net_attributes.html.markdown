---
layout: "outscale"
page_title: "3DS OUTSCALE: outscale_net_attribute"
sidebar_current: "outscale-net-attribute"
description: |-
  [Manages a Net attribute.]
---

# outscale_net_attribute Resource

Manages a Net attribute.
For more information on this resource, see the [User Guide](https://wiki.outscale.net/display/EN/About+VPCs).
For more information on this resource actions, see the [API documentation](https://docs-beta.outscale.com/#updatenet).

## Example Usage

```hcl
[exemple de code]
```

## Argument Reference

The following arguments are supported:

* `dhcp_options_set_id` - (Required) The ID of the DHCP options set (or `default` if you want to associate the default one).
* `net_id` - (Required) The ID of the Net.

## Attribute Reference

The following attributes are exported:

* `dhcp_options_set_id` - The ID of the DHCP options set (or `default` if you want to associate the default one).
* `net_id` - The ID of the Net.
