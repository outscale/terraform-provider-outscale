---
layout: "outscale"
page_title: "OUTSCALE: outscale_client_gateways"
subcategory: "OUTSCALE API"
sidebar_current: "outscale-client-gateways"
description: |-
  [Provides information about client gateways.]
---

# outscale_client_gateways Data Source

Provides information about client gateways.

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Client-Gateways.html).  
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-clientgateway).

## Example Usage

```hcl
data "outscale_client_gateways" "client_gateways01" {
    filter {
        name   = "bgp_asns"
        values = ["65000"]
    }
    filter {
        name   = "public_ips"
        values = ["111.11.111.1", "222.22.222.2"]
    }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) A combination of a filter name and one or more filter values. You can specify this argument for as many filter names as you need. The filter name can be any of the following:
    * `bgp_asns` - (Optional) The Border Gateway Protocol (BGP) Autonomous System Numbers (ASNs) of the connections.
    * `client_gateway_ids` - (Optional) The IDs of the client gateways.
    * `connection_types` - (Optional) The types of communication tunnels used by the client gateways (always `ipsec.1`).
    * `public_ips` - (Optional) The public IPv4 addresses of the client gateways.
    * `states` - (Optional) The states of the client gateways (`pending` \| `available` \| `deleting` \| `deleted`).
    * `tag_keys` - (Optional) The keys of the tags associated with the client gateways.
    * `tag_values` - (Optional) The values of the tags associated with the client gateways.
    * `tags` - (Optional) The key/value combinations of the tags associated with the client gateways, in the following format: `TAGKEY=TAGVALUE`.

## Attribute Reference

The following attributes are exported:

* `client_gateways` - Information about one or more client gateways.
    * `bgp_asn` - The Autonomous System Number (ASN) used by the Border Gateway Protocol (BGP) to find the path to your client gateway through the Internet.
    * `client_gateway_id` - The ID of the client gateway.
    * `connection_type` - The type of communication tunnel used by the client gateway (always `ipsec.1`).
    * `public_ip` - The public IPv4 address of the client gateway (must be a fixed address into a NATed network).
    * `state` - The state of the client gateway (`pending` \| `available` \| `deleting` \| `deleted`).
    * `tags` - One or more tags associated with the client gateway.
        * `key` - The key of the tag, with a minimum of 1 character.
        * `value` - The value of the tag, between 0 and 255 characters.
