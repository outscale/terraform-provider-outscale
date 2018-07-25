---
layout: "outscale"
page_title: "OUTSCALE: outscale_public_ip"
sidebar_current: "docs-outscale-datasource-public-ip"
description: |-
  Describes one External IP address allocated to your account.
---

# outscale_public_ip

By default, this action returns information about your EIP: available, associated with an instance or network interface, or used for a NAT gateway.
## Example Usage

```hcl

data "outscale_public_ip" "by_allocation_id" {
  allocation_id = "${outscale_public_ip.test.allocation_id}"
}
data "outscale_public_ip" "by_public_ip" {
  public_ip = "${outscale_public_ip.test.public_ip}"
}
```

## Argument Reference

The following arguments are supported:

* `reservation_id` - (Optional) One allocation IDs.
* `filter` - (Optional) One or more filters.
* `public_ip` - (Optional) One External IP address.

## Filters

Use the Filter.N parameter to filter the described instances on the following properties:

* `domain` Whether the EIP is for use in the public Cloud or in a VPC.
* `instance-id` The ID of the instance the address is associated with (if any).
* `public-ip` The EIP.
* `allocation-id` The allocation ID for the EIP.
* `association-id` The association ID for the EIP.
* `network-interface-id` The ID of the network interface the address is associated with (if any).
* `network-interface-owner-id` The account ID of the owner.
* `private-ip-address` The private IP address associated with the EIP.

## Attributes Reference

The following attributes are exported:

* `reservation_id` - The ID of the allocation.
* `link_id` - The association ID for the EIP.
* `placement` - Whether the EIP is for use in the public Cloud or in a VPC.
* `nic_id` - The ID of the network interface the address is associated with (if any).
* `nic_account_id` - The account ID of the owner.
* `private_ip` - The private IP address associated with the EIP.
* `public_ip` - The EIP.
* `request_id` - The ID of the request.