---
layout: "outscale"
page_title: "OUTSCALE: outscale_internet_service_link"
sidebar_current: "docs-outscale-resource-internet_service-link"
description: |-
  Attaches an Internet service to a Net.
---

# outscale_internet_service_link

To enable the connection between the Internet and a Net, you must attach an Internet service to this Net.

## Example Usage

```hcl
resource "outscale_internet_service" "service" {}

resource "outscale_net" "net" {
  ip_range = "10.0.0.0/16"
}

resource "outscale_internet_service_link" "link" {
  net_id = "${outscale_net.net.id}"
  internet_service_id = "${outscale_internet_service.service.id}"
}
```

## Argument Reference

The following arguments are supported:

* `internet_service_id` - The ID of the Internet service you want to attach.
* `net_id` - The ID of the Net to which you want to attach the Internet service.

See detailed information in [Outscale Link Internet Service](http://docs.outscale.com/api_fcu/operations/Action_AttachInternetGateway_get.html#_api_fcu-action_attachinternetservice_get).
More info here [Link Internet Service](https://docs-beta.outscale.com/oapi#linkinternetservice)

## Attributes Reference

The following attributes are exported:

* `request_id` - The ID of the request.
* `state` - The state of the attachment of the Net to the Internet service (always `available`).
* `tags` - One or more tags associated with the Internet service.

See detailed information in [Outscale Link Internet Service](http://docs.outscale.com/api_fcu/operations/Action_AttachInternetGateway_get.html#_api_fcu-action_attachinternetservice_get).
More info here [Link Internet Service](https://docs-beta.outscale.com/oapi#linkinternetservice)
