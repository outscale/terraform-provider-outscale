---
layout: "outscale"
page_title: "OUTSCALE: outscale_public_ip"
sidebar_current: "docs-outscale-resource-public_ip"
description: |-
  Provides an Outscale Public IP Association as a top level resource, to associate and disassociate Public IPs from Outscale VMs and Network Interfaces.
---

# outscale_public_ip

NOTE: outscale_public_ip is useful in scenarios where Public IPs are either pre-existing or distributed to customers or users and therefore cannot be changed.

## Example Usage

```hcl
resource "outscale_public_ip" "bar" {}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Optional) The type of platform in which you want to use the EIP (standard | vpc).

## Attributes Reference

* `domain` - (Optional) The type of platform in which you can use the EIP.

## Attribute reference

* `allocation_id` - The ID that represents the allocation of the EIP for use with instances in a VPC.
* `domain` - The type of platform in which you can use the EIP.
* `public_ip` - The External IP address.
* `request_id` - The ID of the request.


See detailed information in [FCU Address](http://docs.outscale.com/api_fcu/definitions/Address.html#_api_fcu-address).
