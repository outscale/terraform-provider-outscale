---
layout: "outscale"
page_title: "3DS OUTSCALE: outscale_internet_service"
sidebar_current: "outscale-internet-service"
description: |-
  [Provides information about a specific Internet service.]
---

# outscale_internet_service Data Source

Provides information about a specific Internet service.
For more information on this resource, see the [User Guide](https://wiki.outscale.net/display/EN/About+Internet+Gateways).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-internetservice).

## Example Usage

```hcl

data "outscale_internet_service" "internet_service01" {
  filter {
    name   = "internet_service_ids"
    values = ["igw-12345678"]
  }
}


```

## Argument Reference

The following arguments are supported:

* `filter` - One or more filters.
  * `internet_service_ids` - (Optional) The IDs of the Internet services.
  * `tag_keys` - (Optional) The keys of the tags associated with the Internet services.
  * `tag_values` - (Optional) The values of the tags associated with the Internet services.
  * `tags` - (Optional) The key/value combination of the tags associated with the Internet services, in the following format: "Filters":{"Tags":["TAGKEY=TAGVALUE"]}.

## Attribute Reference

The following attributes are exported:

* `internet_services` - Information about one or more Internet services.
  * `internet_service_id` - The ID of the Internet service.
  * `net_id` - The ID of the Net attached to the Internet service.
  * `state` - The state of the attachment of the Net to the Internet service (always `available`).
  * `tags` - One or more tags associated with the Internet service.
    * `key` - The key of the tag, with a minimum of 1 character.
    * `value` - The value of the tag, between 0 and 255 characters.
