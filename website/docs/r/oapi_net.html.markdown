---
layout: "outscale"
page_title: "OUTSCALE: outscale_net"
sidebar_current: "docs-outscale-resource-net"
description: |-
  Creates a Net with a specified IP range.
---

# outscale_net

Creates a Net with a specified IP range.

The IP range (network range) of your Net must be between a /28 netmask (16 IP addresses) and a /16 netmask (65 536 IP addresses).

## Example Usage

```hcl
resource "outscale_net" "net" {
  ip_range = "10.0.0.0/16"
}
```

## Argument Reference

The following arguments are supported:

* `ip_range` - (Required) The IP range for the Net, in CIDR notation (for example, 10.0.0.0/16).
* `tenancy` - (Optional) The tenancy options for the VMs (`default` if a VM created in a Net can be launched with any tenancy, `dedicated` if it can be launched with dedicated tenancy VMs running on single-tenant hardware).

See detailed information in [Outscale VMs](http://docs.outscale.com/api_fcu/operations/Action_CreateVpc_get.html#_api_fcu-action_createvpc_get).

## Attributes Reference

The following attributes are exported:

* `ip_range` - The IP range for the Net, in CIDR notation (for example, 10.0.0.0/16).
* `tenancy` - The VM tenancy in a Net.
* `dhcp_options_set_id` - The ID of the DHCP options set (or `default` if you want to associate the default one).
* `state` - The state of the Net (`pending` | `available`)
* `tags` - One or more tags associated with the Net.
* `net_id` - The ID of the Net.
* `request_id` - The ID of the request.

See detailed information in [CreateNet](http://docs.outscale.com/api_fcu/operations/Action_CreateVpc_get.html#_api_fcu-action_createvpc_get).
