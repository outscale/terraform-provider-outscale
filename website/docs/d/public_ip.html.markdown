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

* `allocation_id` - (Optional) One allocation IDs.
* `filter` - (Optional) One or more filters.
* `public_ip` - (Optional) One External IP address.


See detailed information in [Outscale Public IP](http://docs.outscale.com/api_fcu/operations/Action_DescribeAddresses_get.html#_api_fcu-action_describeaddresses_get).

## Filters

Use the Filter.N parameter to filter the described instances on the following properties:

* `domain` - Whether the EIP is for use in the public Cloud or in a VPC.
* `instance-id` - The ID of the instance the address is associated with (if any).
* `public-ip` - The EIP.
* `allocation-id` - The allocation ID for the EIP.
* `association-id` - The association ID for the EIP.
* `network-interface-id` - The architecture of the instance (i386 | x86_64).
* `network-interface-owner-id` - The account ID of the owner.
* `private-ip-address` - The private IP address associated with the EIP.

## Attributes Reference

The following attributes are exported:
* `allocation_id` The allocation ID for the EIP.
* `association_id` The association ID for the EIP.
* `domain` Whether the EIP is for use in the public Cloud or in a VPC.
* `instance_id` The ID of the instance the address is associated with (if any).
* `network_interface_id` The ID of the network interface the address is associated with (if any).
* `network_interface_owner_id` The account ID of the owner.
* `private_ip_address` The private IP address associated with the EIP.
* `public_ip` The EIP.
* `request_id` - The ID of the request.

See detailed information in [Describe Public IP](http://docs.outscale.com/api_fcu/operations/Action_DescribeAddresses_get.html#_api_fcu-action_describeaddresses_get).
