---
layout: "outscale"
page_title: "OUTSCALE: outscale_directlink_interface"
sidebar_current: "docs-outscale-directlink-interface"
description: |-
  Creates a private virtual interface.
---

# outscale_directlink_interface

Creates a private virtual interface.

Private interfaces enable you to reach one of your Virtual Private Clouds (VPCs) through a virtual private gateway.

## Example Usage

```hcl
resource "outscale_vpn_gateway" "foo" {
    tag {
        Name = "terraform-testacc-dxvif-1"
    }
}

resource "outscale_directlink_interface" "foo" {
    connection_id = "TBD"

    new_private_virtual_interface {
        virtual_gateway_id = "${outscale_vpn_gateway.foo.id}"
        virtual_interface_name = "terraform-testacc-dxvif-1"
        vlan = 4094
        asn  = 3
    }
}
```

## Argument Reference

The following arguments are supported:

* `connection_id` - The ID of the connection over which you want to create the private virtual interface.
* `new_private_virtual_interface` - Detailed information about the configuration parameters of the private virtual interface.
  * `amazon_address` - If provided, the IP address to set on Outscale’s side of the virtual interface that must include a network prefix (for example 172.16.0.1/30).
  * `asn` - The BGP (Border Gateway Protocol) AS (autonomous system) number on the customer’s side of the virtual interface.
  * `auth_key` - The BGP authentication key.
  * `customer_address` - The IP address on the customer’s side of the virtual interface (must be provided if the AmazonAddress is provided, and be in the same network subnet).
  * `virtual_gateway_id` - The target virtual private gateway.
  * `virtual_interface_name` - The name of the virtual interface.
  * `vlan` - The unique VLAN ID for the virtual interface.

## Attributes Reference

The following attributes are exported:

* `amazon_address` - If provided, the IP address to set on Outscale’s side of the virtual interface that must include a network prefix (for example 172.16.0.1/30).
* `asn` - The BGP (Border Gateway Protocol) AS (autonomous system) number on the customer’s side of the virtual interface.
* `auth_key` - The BGP authentication key.
* `connection_id` - The ID of the DirectLink connection.
* `customer_address` - The IP address on the customer’s side of the virtual interface (must be provided if the AmazonAddress is provided, and be in the same network subnet).
* `location` - The datacenter where the virtual interface is located.
* `owner_account` - The account ID of the owner of the virtual interface.
* `virtual_gateway_id` - The target virtual private gateway.
* `virtual_interface_id` - The ID of the virtual interface.
* `virtual_interface_name` - The name of the virtual interface.
* `virtual_interface_state` - The state of the virtual interface (pending | available | deleting | deleted).
* `virtual_interface_type` - The type of the virtual interface (always private).
* `vlan` - The unique VLAN ID for the virtual interface.
* `request_id` - The ID of the request

[See detailed description](http://docs.outscale.com/api_directlink/operations/Action_CreatePrivateVirtualInterface_get.html#_api_directlink-action_createprivatevirtualinterface_get)
