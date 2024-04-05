---
layout: "outscale"
page_title: "OUTSCALE: outscale_virtual_gateways"
sidebar_current: "outscale-virtual-gateways"
description: |-
  [Provides information about virtual gateways.]
---

# outscale_virtual_gateways Data Source

Provides information about virtual gateways.

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Virtual-Private-Gateways.html).  
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-virtualgateway).

## Example Usage

```hcl
data "outscale_virtual_gateways" "virtual_gateways01" {
    filter {
        name   = "states"
        values = ["available"]
    }
    filter {
        name   = "link_states"
        values = ["attached", "detached"]
    }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) A combination of a filter name and one or more filter values. You can specify this argument for as many filter names as you need. The filter name can be any of the following:
    * `connection_types` - (Optional) The types of the virtual gateways (only `ipsec.1` is supported).
    * `link_net_ids` - (Optional) The IDs of the Nets the virtual gateways are attached to.
    * `link_states` - (Optional) The current states of the attachments between the virtual gateways and the Nets (`attaching` \| `attached` \| `detaching` \| `detached`).
    * `states` - (Optional) The states of the virtual gateways (`pending` \| `available` \| `deleting` \| `deleted`).
    * `tag_keys` - (Optional) The keys of the tags associated with the virtual gateways.
    * `tag_values` - (Optional) The values of the tags associated with the virtual gateways.
    * `tags` - (Optional) The key/value combinations of the tags associated with the virtual gateways, in the following format: `TAGKEY=TAGVALUE`.
    * `virtual_gateway_ids` - (Optional) The IDs of the virtual gateways.
* `next_page_token` - (Optional) The token to request the next page of results. Each token refers to a specific page.
* `results_per_page` - (Optional) The maximum number of logs returned in a single response (between `1`and `1000`, both included). By default, `100`.

## Attribute Reference

The following attributes are exported:

* `next_page_token` - The token to request the next page of results. Each token refers to a specific page.
* `virtual_gateways` - Information about one or more virtual gateways.
    * `connection_type` - The type of VPN connection supported by the virtual gateway (only `ipsec.1` is supported).
    * `net_to_virtual_gateway_links` - The Net to which the virtual gateway is attached.
        * `net_id` - The ID of the Net to which the virtual gateway is attached.
        * `state` - The state of the attachment (`attaching` \| `attached` \| `detaching` \| `detached`).
    * `state` - The state of the virtual gateway (`pending` \| `available` \| `deleting` \| `deleted`).
    * `tags` - One or more tags associated with the virtual gateway.
        * `key` - The key of the tag, with a minimum of 1 character.
        * `value` - The value of the tag, between 0 and 255 characters.
    * `virtual_gateway_id` - The ID of the virtual gateway.
