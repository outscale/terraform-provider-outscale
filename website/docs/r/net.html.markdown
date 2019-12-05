---
layout: "outscale"
page_title: "3DS OUTSCALE: outscale_net"
sidebar_current: "docs-outscale-resource-net"
description: |-
  [Manages a Net.]
---

# outscale_net Resource

Manages a Net.
For more information on this resource, see the [User Guide](https://wiki.outscale.net/display/EN/About+VPCs).
For more information on this resource actions, see the [API documentation](https://docs-beta.outscale.com/#3ds-outscale-api-net).

## Example Usage

```hcl
[exemple de code]
```

## Argument Reference

The following arguments are supported:

* `ip_range` - (Required) The IP range for the Net, in CIDR notation (for example, 10.0.0.0/16).
* `tenancy` - (Optional) The tenancy options for the VMs (`default` if a VM created in a Net can be launched with any tenancy, `dedicated` if it can be launched with dedicated tenancy VMs running on single-tenant hardware).

## Attribute Reference

The following attributes are exported:

* `net` - Information about the Net.
  * `dhcp_options_set_id` - The ID of the DHCP options set (or `default` if you want to associate the default one).
  * `ip_range` - The IP range for the Net, in CIDR notation (for example, 10.0.0.0/16).
  * `net_id` - The ID of the Net.
  * `state` - The state of the Net (`pending` \| `available`).
  * `tags` - One or more tags associated with the Net.
    * `key` - The key of the tag, with a minimum of 1 character.
    * `value` - The value of the tag, between 0 and 255 characters.
  * `tenancy` - The VM tenancy in a Net.
