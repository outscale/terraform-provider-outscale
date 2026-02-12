---
layout: "outscale"
page_title: "OUTSCALE: outscale_internet_service"
subcategory: "OUTSCALE API"
sidebar_current: "outscale-internet-service"
description: |-
  [Manages an Internet service.]
---

# outscale_internet_service Resource

Manages an Internet service.

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Internet-Services.html).  
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

* `internet_service_id` - The ID of the internet service.
* `net_id` - The ID of the Net attached to the internet service.
* `state` - The state of the attachment of the internet service to the Net (always `available`).
* `tags` - One or more tags associated with the internet service.
    * `key` - The key of the tag, with a minimum of 1 character.
    * `value` - The value of the tag, between 0 and 255 characters.

## Timeouts

The `timeouts` block enables you to configure [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - Defaults to 10 minutes.
* `read` - Defaults to 5 minutes.
* `update` - Defaults to 10 minutes.
* `delete` - Defaults to 5 minutes.

## Import

An internet service can be imported using its ID. For example:

```console

$ terraform import outscale_internet_service.ImportedInternetService igw-12345678

```