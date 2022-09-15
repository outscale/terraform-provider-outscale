---
layout: "outscale"
page_title: "OUTSCALE: outscale_net"
sidebar_current: "outscale-net"
description: |-
  [Manages a Net.]
---

# outscale_net Resource

Manages a Net.
For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-VPCs.html).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-net).

## Example Usage

### Create a Net

```hcl
resource "outscale_net" "net01" {
	ip_range = "10.10.0.0/16"
	tenancy  = "default"
}
```

### Create a Net with a network

```hcl
resource "outscale_net" "net02" {
	ip_range = "10.0.0.0/16"
}

resource "outscale_subnet" "subnet01" {
	net_id   = outscale_net.net02.net_id
	ip_range = "10.0.0.0/18"
}

resource "outscale_public_ip" "public_ip01" {
}

resource "outscale_nat_service" "nat_service01" {
	subnet_id    = outscale_subnet.subnet01.subnet_id
	public_ip_id = outscale_public_ip.public_ip01.public_ip_id
}

resource "outscale_route_table" "route_table01" {
	net_id = outscale_net.net02.net_id
}

resource "outscale_route" "route01" {
	destination_ip_range = "0.0.0.0/0"
	gateway_id           = outscale_internet_service.internet_service01.internet_service_id
	route_table_id       = outscale_route_table.route_table01.route_table_id
}

resource "outscale_route_table_link" "route_table_link01" {
	subnet_id      = outscale_subnet.subnet01.subnet_id
	route_table_id = outscale_route_table.route_table01.id
}

resource "outscale_internet_service" "internet_service01" {
}

resource "outscale_internet_service_link" "internet_service_link01" {
	net_id              = outscale_net.net02.net_id
	internet_service_id = outscale_internet_service.internet_service01.id
}
```

## Argument Reference

The following arguments are supported:

* `ip_range` - (Required) The IP range for the Net, in CIDR notation (for example, 10.0.0.0/16).
* `tags` - (Optional) A tag to add to this resource. You can specify this argument several times.
    * `key` - (Required) The key of the tag, with a minimum of 1 character.
    * `value` - (Required) The value of the tag, between 0 and 255 characters.
* `tenancy` - (Optional) The tenancy options for the VMs (`default` if a VM created in a Net can be launched with any tenancy, `dedicated` if it can be launched with dedicated tenancy VMs running on single-tenant hardware).

## Attribute Reference

The following attributes are exported:

* `dhcp_options_set_id` - The ID of the DHCP options set (or `default` if you want to associate the default one).
* `ip_range` - The IP range for the Net, in CIDR notation (for example, 10.0.0.0/16).
* `net_id` - The ID of the Net.
* `state` - The state of the Net (`pending` \| `available`).
* `tags` - One or more tags associated with the Net.
    * `key` - The key of the tag, with a minimum of 1 character.
    * `value` - The value of the tag, between 0 and 255 characters.
* `tenancy` - The VM tenancy in a Net.

## Import

A Net can be imported using its ID. For example:

```console

$ terraform import outscale_net.ImportedNet vpc-87654321

```