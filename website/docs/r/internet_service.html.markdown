---
layout: "outscale"
page_title: "3DS OUTSCALE: outscale_internet_service"
sidebar_current: "outscale-internet-service"
description: |-
  [Manages an Internet service.]
---

# outscale_internet_service Resource

Manages an Internet service.
For more information on this resource, see the [User Guide](https://wiki.outscale.net/display/EN/About+Internet+Gateways).
For more information on this resource actions, see the [API documentation](https://docs-beta.outscale.com/#3ds-outscale-api-internetservice).

## Example Usage

```hcl

resource "outscale_internet_service" "internet_service01" {	
}


```

## Argument Reference

The following arguments are supported:

* `tags` - One or more tags to add to this resource.
    * `key` - The key of the tag, with a minimum of 1 character.
    * `value` - The value of the tag, between 0 and 255 characters.

## Attribute Reference

The following attributes are exported:

* `internet_service` - Information about the Internet service.
  * `internet_service_id` - The ID of the Internet service.
  * `net_id` - The ID of the Net attached to the Internet service.
  * `state` - The state of the attachment of the Net to the Internet service (always `available`).
  * `tags` - One or more tags associated with the Internet service.
    * `key` - The key of the tag, with a minimum of 1 character.
    * `value` - The value of the tag, between 0 and 255 characters.
