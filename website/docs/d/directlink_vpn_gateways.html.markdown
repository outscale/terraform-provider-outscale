---
layout: "outscale"
page_title: "OUTSCALE: outscale_directlink_vpn_gateways"
sidebar_current: "docs-outscale-datasource-sites"
description: |-
 Returns a list of your virtual gateways that can be used as a target by a private virtual interface.
---

# outscale_directlink_vpn_gateways

Returns a list of your virtual gateways that can be used as a target by a private virtual interface.

## Example Usage

```hcl
data "outscale_directlink_vpn_gateways" "test" {}
```

## Argument Reference

No arguments are supported

## Attributes Reference

The following attributes are exported:

* `virtual_gateways.N` - Information about one or more virtual gateways.
  * `virtual_gateway_id` - The ID of the virtual gateway.
  * `virtual_gateway_state` - The state of the virtual gateway (pending | available | deleting | deleted).
* `request_id` - The ID of the request.

[See detailed description](http://docs.outscale.com/api_directlink/operations/Action_DescribeVirtualGateways_get.html#_api_directlink-action_describevirtualgateways_get)
