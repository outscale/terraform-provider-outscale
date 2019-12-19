---
layout: "outscale"
page_title: "3DS OUTSCALE: outscale_nat_service"
sidebar_current: "outscale-nat-service"
description: |-
  [Provides information about NAT services.]
---

# outscale_nat_service Data Source

Provides information about NAT services.
For more information on this resource, see the [User Guide](https://wiki.outscale.net/display/EN/About+NAT+Devices).
For more information on this resource actions, see the [API documentation](https://docs-beta.outscale.com/#3ds-outscale-api-natservice).

## Example Usage

```hcl

data "outscale_nat_services" "nat_services01" {
  filter {
    name   = "nat_service_ids"
    values = ["nat-12345678", "nat-12345679"]
  }
}


```

## Argument Reference

The following arguments are supported:

* `filters` - One or more filters.
  * `nat_service_ids` - (Optional) The IDs of the NAT services.
  * `net_ids` - (Optional) The IDs of the Nets in which the NAT services are.
  * `states` - (Optional) The states of the NAT services (`pending` \| `available` \| `deleting` \| `deleted`).
  * `subnet_ids` - (Optional) The IDs of the Subnets in which the NAT services are.
  * `tag_keys` - (Optional) The keys of the tags associated with the NAT services.
  * `tag_values` - (Optional) The values of the tags associated with the NAT services.
  * `tags` - (Optional) The key/value combination of the tags associated with the NAT services, in the following format: "Filters":{"Tags":["TAGKEY=TAGVALUE"]}.

## Attribute Reference

The following attributes are exported:

* `nat_services` - Information about one or more NAT services.
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
