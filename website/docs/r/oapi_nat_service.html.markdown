---
layout: "outscale"
page_title: "OUTSCALE: outscale_nat_service"
sidebar_current: "docs-outscale-resource-nat-service"
description: |-
  Provides an Outscale Nat Gateway resource. This allows instances to be created, described, and deleted. Nat Gateway also support provisioning.
---

# outscale_nat_service

  Provides an Outscale Nat Gateway resource. This allows instances to be created, described, and deleted. Nat Gateway also support [provisioning](/docs/provisioners/index.html).

## Example Usage

```hcl
resource "outscale_nat_service" "gateway" {
    public_ip_id = "eipalloc-32e506e8"
    subnet_id = "subnet-861fbecc"
}
```

## Argument Reference

The following arguments are supported:

* `public_ip_id` - (Required) The allocation ID of the EIP to associate with the NAT service.
* `subnet_id` - (Required) The ID of the Subnet in which you want to create the NAT service.
* `token` - (Optional) A unique identifier which enables you to manage the idempotency.

## Attributes Reference

The following attributes are exported:

* `public_ips` - Information about the External IP address or addresses (EIPs) associated with the NAT service (List).
* `nat_service_id` - The ID of the NAT service.
* `state` - The state of the NAT gateway (pending | available| deleting | deleted).
* `subnet-id` - The ID of the subnet in which the NAT gateway is.
* `net_id` - The ID of the Net in which the NAT service is.
* `request_id` - The ID of the request.
