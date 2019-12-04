---
layout: "outscale"
page_title: "3DS OUTSCALE: outscale_vpn_connection_route"
sidebar_current: "docs-outscale-resource-vpn-connection-route"
description: |-
  [Manages a VPN connection route.]
---

# outscale_vpn_connection_route Resource

Manages a VPN connection route.
For more information on this resource, see the [User Guide](https://wiki.outscale.net/display/EN/About+Routing+Configuration+for+VPN+Connections).
For more information on this resource actions, see the [API documentation](https://docs-beta.outscale.com/#3ds-outscale-api-vpnconnection).

## Example Usage

```hcl
[exemple de code]
```

## Argument Reference

The following arguments are supported:

* `destination_ip_range` - (Required) The network prefix of the route, in CIDR notation (for example, 10.12.0.0/16).
* `vpn_connection_id` - (Required) The ID of the target VPN connection of the static route.

## Attribute Reference

The following attributes are exported:

