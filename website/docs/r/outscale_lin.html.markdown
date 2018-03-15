---
layout: "outscale"
page_title: "OUTSCALE: outscale_lin"
sidebar_current: "docs-outscale-resource-lin"
description: |-
  Creates a Virtual Private Cloud (VPC) with a specified CIDR block.
---

# outscale_lin

Creates a Virtual Private Cloud (VPC) with a specified CIDR block.
The CIDR block (network range) of your VPC must be between a /28 netmask (16 IP addresses) and a /16 netmask (65 536 IP addresses).

## Example Usage

```hcl
resource "outscale_lin" "vpc" {
	cidr_block = "10.0.0.0/16"
}
```

## Argument Reference

The following arguments are supported:

* `cidr_block` - The CIDR block for the VPC (for example, 10.0.0.0/16).
* `instance_tenancy` - The tenancy options of the instances (`default` if an instance created in a VPC can be lauched with any tenancy, `dedicated` if it can be launched with dedicated tenancy instances running on single-tenant hardware).

See detailed information in [Outscale Instances](http://docs.outscale.com/api_fcu/operations/Action_CreateVpc_get.html#_api_fcu-action_createvpc_get).


## Attributes Reference

The following attributes are exported:

* `cidr_block` - The CIDR block of the VPC, in the [16;28] range (for example, 10.84.7.0/24).
* `dhcp_options_id` - The ID of the DHCP options set associated with the VPC.
* `instance_tenancy` -	The tenancy of the instances (`default` if an instance created in a VPC can be lauched with any tenancy, `dedicated` if it can be launched with dedicated tenancy instances running on single-tenant hardware).
* `state` - The state of the VPC (`pending` | `available`)
* `tag_set` - One or more tags associated with the VPC.
* `vpc_id` - The ID of the Virtual Private Cloud (VPC).


See detailed information in [Create VPC](http://docs.outscale.com/api_fcu/operations/Action_CreateVpc_get.html#_api_fcu-action_createvpc_get).
