---
layout: "outscale"
page_title: "3DS OUTSCALE: outscale_nat_service"
sidebar_current: "outscale-nat-service"
description: |-
  [Manages a NAT service.]
---

# outscale_nat_service Resource

Manages a NAT service.
For more information on this resource, see the [User Guide](https://wiki.outscale.net/display/EN/About+NAT+Devices).
For more information on this resource actions, see the [API documentation](https://docs-beta.outscale.com/#3ds-outscale-api-natservice).

## Example Usage

```hcl
[exemple de code]
```

## Argument Reference

The following arguments are supported:

* `public_ip_id` - (Required) The allocation ID of the EIP to associate with the NAT service.<br />
If the EIP is already associated with another resource, you must first disassociate it.
* `subnet_id` - (Required) The ID of the Subnet in which you want to create the NAT service.

## Attribute Reference

The following attributes are exported:

* `nat_service` - Information about the NAT service.
  * `nat_service_id` - The ID of the NAT service.
  * `net_id` - The ID of the Net in which the NAT service is.
  * `public_ips` - Information about the External IP address or addresses (EIPs) associated with the NAT service.
    * `public_ip` - The External IP address (EIP) associated with the NAT service.
    * `public_ip_id` - The allocation ID of the EIP associated with the NAT service.
  * `state` - The state of the NAT service (`pending` \| `available` \| `deleting` \| `deleted`).
  * `subnet_id` - The ID of the Subnet in which the NAT service is.
  * `tags` - One or more tags associated with the NAT service.
    * `key` - The key of the tag, with a minimum of 1 character.
    * `value` - The value of the tag, between 0 and 255 characters.
