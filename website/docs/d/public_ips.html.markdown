---
layout: "outscale"
page_title: "OUTSCALE: outscale_public_ips"
sidebar_current: "docs-outscale-datasource-public-ips"
description: |-
  Describes one or more External IP addresses (EIPs) allocated to your account.
---

# outscale_public_ips

By default, this action returns information about all your EIPs: available, associated with an instance or network interface, or used for a NAT gateway.
## Example Usage

```hcl
  data "outscale_public_ips" "by_public_ips" {
    public_ips = ["${outscale_public_ip.test.public_ip}", "${outscale_public_ip.test1.public_ip}", "${outscale_public_ip.test2.public_ip}"]
  }
```

## Argument Reference

The following arguments are supported:

* `AllocationId` - (Optional) One or more allocation IDs.
* `Filter.N` - (Optional) One or more filters.
* `PublicIp` - (Optional) One or more External IP address.


See detailed information in [Outscale Public IPs](http://docs.outscale.com/api_fcu/operations/Action_DescribeAddresses_get.html#_api_fcu-action_describeaddresses_get).

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

* `request_id` - The ID of the request.
* `addresses_set` - Information about one or more External IP addresses.



See detailed information in [Describe Public IPs](http://docs.outscale.com/api_fcu/operations/Action_DescribeAddresses_get.html#_api_fcu-action_describeaddresses_get).
