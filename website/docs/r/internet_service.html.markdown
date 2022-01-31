---
layout: "outscale"
page_title: "OUTSCALE: outscale_internet_service"
sidebar_current: "outscale-internet-service"
description: |-
  [Manages an Internet service.]
---

# outscale_internet_service Resource

Manages an Internet service.
For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Internet-Gateways.html).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-internetservice).

## Example Usage

```hcl
resource "outscale_internet_service" "internet_service01" {	
}
```

## Argument Reference

The following arguments are supported:

* `tags` - (Optional) A tag to add to this resource. You can specify this argument several times.
    * `key` - (Required) The key of the tag, with a minimum of 1 character.
    * `value` - (Required) The value of the tag, between 0 and 255 characters.

## Attribute Reference

The following attributes are exported:

* `internet_service_id` - The ID of the Internet service.
* `net_id` - The ID of the Net attached to the Internet service.
* `state` - The state of the attachment of the Internet service to the Net (always `available`).
* `tags` - One or more tags associated with the Internet service.
    * `key` - The key of the tag, with a minimum of 1 character.
    * `value` - The value of the tag, between 0 and 255 characters.

## Import

An internet service can be imported using its ID. For example:

```console

$ terraform import outscale_internet_service.ImportedInternetService igw-12345678

```