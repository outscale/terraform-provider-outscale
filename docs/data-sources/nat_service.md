---
layout: "outscale"
page_title: "OUTSCALE: outscale_nat_service"
sidebar_current: "outscale-nat-service"
description: |-
  [Provides information about a specific NAT service.]
---

# outscale_nat_service Data Source

Provides information about a specific NAT service.
For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-NAT-Gateways.html).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-natservice).

## Example Usage

```hcl
data "outscale_nat_service" "nat_service01" {
  filter {
    name   = "nat_service_ids"
    values = ["nat-12345678"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) A combination of a filter name and one or more filter values. You can specify this argument for as many filter names as you need. The filter name can be any of the following:
    * `nat_service_ids` - (Optional) The IDs of the NAT services.
    * `net_ids` - (Optional) The IDs of the Nets in which the NAT services are.
    * `states` - (Optional) The states of the NAT services (`pending` \| `available` \| `deleting` \| `deleted`).
    * `subnet_ids` - (Optional) The IDs of the Subnets in which the NAT services are.
    * `tag_keys` - (Optional) The keys of the tags associated with the NAT services.
    * `tag_values` - (Optional) The values of the tags associated with the NAT services.
    * `tags` - (Optional) The key/value combination of the tags associated with the NAT services, in the following format: &quot;Filters&quot;:{&quot;Tags&quot;:[&quot;TAGKEY=TAGVALUE&quot;]}.

## Attribute Reference

The following attributes are exported:

* `nat_service_id` - The ID of the NAT service.
* `net_id` - The ID of the Net in which the NAT service is.
* `public_ips` - Information about the public IP or IPs associated with the NAT service.
    * `public_ip` - The public IP associated with the NAT service.
    * `public_ip_id` - The allocation ID of the public IP associated with the NAT service.
* `state` - The state of the NAT service (`pending` \| `available` \| `deleting` \| `deleted`).
* `subnet_id` - The ID of the Subnet in which the NAT service is.
* `tags` - One or more tags associated with the NAT service.
    * `key` - The key of the tag, with a minimum of 1 character.
    * `value` - The value of the tag, between 0 and 255 characters.
