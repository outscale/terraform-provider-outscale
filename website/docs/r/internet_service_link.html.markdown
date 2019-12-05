---
layout: "outscale"
page_title: "3DS OUTSCALE: outscale_internet_service_link"
sidebar_current: "docs-outscale-resource-internet-service-link"
description: |-
  [Manages an internet service link.]
---

# outscale_internet_service_link Resource

Manages an internet service link.
For more information on this resource, see the [User Guide](https://wiki.outscale.net/display/EN/About+Internet+Gateways).
For more information on this resource actions, see the [API documentation](https://docs-beta.outscale.com/#linkinternetservice).

## Example Usage

```hcl
[exemple de code]
```

## Argument Reference

The following arguments are supported:

* `internet_service_id` - (Required) The ID of the Internet service you want to attach.
* `net_id` - (Required) The ID of the Net to which you want to attach the Internet service.

## Attribute Reference

The following attributes are exported:

* `internet_service_id` - The ID of the Internet service you want to attach.
* `state` - The state of the attachment of the Net to the Internet service (always `available`).
* `tags` - One or more tags associated with the Internet service.
