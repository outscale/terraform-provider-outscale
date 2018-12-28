---
layout: "outscale"
page_title: "OUTSCALE: outscale_internet_service"
sidebar_current: "docs-outscale-resource-internet-service"
description: |-
  Creates an Internet service you can use with a Net.
---

# outscale_internet_service

An Internet service enables your virtual machines (VMs) launched in a Net to connect to the Internet. By default, a Net includes an Internet service, and each Subnet is public. Every VM launched within a default Subnet has a private and a public IP addresses.

## Example Usage

```hcl
resource "outscale_internet_service" "iservice" {}
```

## Attributes Reference

The following attributes are exported:

* `internet_service_id` - The ID of the Internet service.
* `net_id` - The ID of the Net attached to the Internet service.
* `state` - The state of the attachment of the Net to the Internet service (always `available`).
* `tags` -One or more tags associated with the Internet service.
* `request_id` - The ID of the request.

See detailed information in [Create Internet Service](http://docs.outscale.com/api_fcu/operations/Action_CreateInternetGateway_get.html#_api_fcu-action_createinternetservice_get).
