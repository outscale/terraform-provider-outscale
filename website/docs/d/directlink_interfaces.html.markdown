---
layout: "outscale"
page_title: "OUTSCALE: outscale_directlink_interfaces"
sidebar_current: "docs-outscale-datasource-directlink-interfaces"
description: |-
  Describes one or more of your virtual interfaces.
---

# outscale_directlink_interfaces

Describes one or more of your virtual interfaces.

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
        asn = 64512
    }
}

data "outscale_directlink_interfaces" "outscale_directlink_interfaces" {
    connection_id = "TBD"
}
```

## Argument Reference

The following arguments are supported:

* `connection_id` - The ID of a DirectLink connection.

## Filters

None.

## Attributes Reference

The following attributes are exported:

* `virtual_interfaces.N` - Information about one or more virtual interfaces.
  * `amazon_address` - The IP address on Outscale’s side of the virtual interface.
  * `asn` - The BGP (Border Gateway Protocol) AS (Autonomous System) number on the customer’s side of the virtual interface.
  * `auth_key` - The BGP authentication key.
  * `connection_id` - The ID of the DirectLink connection.
  * `customer_address` - The IP address on the customer’s side of the virtual interface.
  * `location` - The datacenter where the virtual interface is located.
  * `owner_account` - The account ID of the owner of the virtual interface.
  * `virtual_gateway_id` - The target virtual private gateway.
  * `virtual_interface_id` - The ID of the virtual interface.
  * `virtual_interface_name` - The name of the virtual interface.
  * `virtual_interface_state` - The state of the virtual interface (pending | available | deleting | deleted).
  * `virtual_interface_type` - The type of the virtual interface (always private).
  * `vlan` - The VLAN number associated with the virtual interface.
* `request_id` - The ID of the request

[See detailed description](http://docs.outscale.com/api_directlink/operations/Action_DescribeVirtualInterfaces_get.html#_api_directlink-action_describevirtualinterfaces_get)
