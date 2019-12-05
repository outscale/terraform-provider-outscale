---
layout: "outscale"
page_title: "3DS OUTSCALE: outscale_internet_service"
sidebar_current: "docs-outscale-resource-internet-service"
description: |-
  [Manages an internet service.]
---

# outscale_internet_service

Manages an internet service.
For more information on this resource, see the [User Guide](https://wiki.outscale.net/display/EN/About+Internet+Gateways).
For more information on this resource actions, see the [API documentation](https://docs-beta.outscale.com/#3ds-outscale-api-internetservice).

## Example Usage

```hcl
[exemple de code]
```

## Argument Reference

The following arguments are supported:


## Attribute Reference

The following attributes are exported:

* `internet_service` - Information about the Internet service.
  * `internet_service_id` - The ID of the Internet service.
  * `net_id` - The ID of the Net attached to the Internet service.
  * `state` - The state of the attachment of the Net to the Internet service (always `available`).
  * `tags` - One or more tags associated with the Internet service.
    * `key` - The key of the tag, with a minimum of 1 character.
    * `value` - The value of the tag, between 0 and 255 characters.
