---
layout: "outscale"
page_title: "OUTSCALE: outscale_internet_service_link"
sidebar_current: "outscale-internet-service-link"
description: |-
  [Manages an Internet service link.]
---

# outscale_internet_service_link Resource

Manages an Internet service link.
For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Internet-Gateways.html).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-internetservice).

## Example Usage

### Required resources

```hcl
resource "outscale_net" "net01" {
	ip_range = "10.0.0.0/18"
}

resource "outscale_internet_service" "internet_service01" {
}
```


### Link an Internet service to a Net

```hcl
resource "outscale_internet_service_link" "internet_service_link01" {
	internet_service_id = outscale_internet_service.internet_service01.internet_service_id
	net_id              = outscale_net.net01.net_id
}
```

## Argument Reference

The following arguments are supported:

* `internet_service_id` - (Required) The ID of the Internet service you want to attach.
* `net_id` - (Required) The ID of the Net to which you want to attach the Internet service.

## Attribute Reference

The following attributes are exported:

* `internet_service_id` - The ID of the Internet service.
* `state` - The state of the attachment of the Internet service to the Net (always `available`).
* `tags` - One or more tags associated with the Internet service.
    * `key` - The key of the tag, with a minimum of 1 character.
    * `value` - The value of the tag, between 0 and 255 characters.

## Import

An internet service link can be imported using the internet service ID. For example:

```console

$ terraform import outscale_internet_service_link.ImportedInternetServiceLink igw-87654321

```