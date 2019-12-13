---
layout: "outscale"
page_title: "3DS OUTSCALE: outscale_net"
sidebar_current: "outscale-net"
description: |-
  [Manages a Net.]
---

# outscale_net Resource

Manages a Net.
For more information on this resource, see the [User Guide](https://wiki.outscale.net/display/EN/About+VPCs).
For more information on this resource actions, see the [API documentation](https://docs-beta.outscale.com/#3ds-outscale-api-net).

## Example Usage

```hcl

# Create a Net with a Subnet

resource "outscale_net" "net01" {
	ip_range = "10.0.0.0/16"
}

resource "outscale_subnet" "subnet01" {
	net_id   = outscale_net.net01.net_id
	ip_range = "10.0.0.0/18"
}

# Create a Net with a security group

resource "outscale_net" "net02" {
	ip_range = "10.0.0.0/16"
}

resource "outscale_security_group" "security_group01" {
	description         = "Terraform security group for Net"
	security_group_name = "terraform-security-group-test-01"
	net_id              = outscale_net.net02.net_id
}

# Create a Net with a VM

resource "outscale_net" "net03" {
	ip_range = "10.0.0.0/16"
}

resource "outscale_subnet" "subnet02" {
	subregion_name = "${var.region}a"
	ip_range       = "10.0.0.0/16"
	net_id         = outscale_net.net03.net_id
}

resource "outscale_security_group" "security_group02" {
	description         = "Terraform security group for Net with VM"
	security_group_name = "terraform-security-group-test-02"
	net_id              = outscale_net.net03.net_id
}

resource "outscale_vm" "outscale_vm01" {
	image_id                 = var.image_id
	vm_type                  = var.vm_type
	keypair_name             = var.keypair_name
	security_group_ids       = [outscale_security_group.security_group02[0].security_group_id]
	placement_subregion_name = "${var.region}a"
	placement_tenancy        = "default"
	subnet_id                = outscale_subnet.subnet02.subnet_id
}

# Create a Net with a network

resource "outscale_net" "net04" {
	ip_range = "10.0.0.0/16"
}

resource "outscale_subnet" "subnet03" {
	net_id   = outscale_net.net04.net_id
	ip_range = "10.0.0.0/18"
}

resource "outscale_public_ip" "public_ip01" {
}

resource "outscale_nat_service" "nat_service01" {
	subnet_id    = outscale_subnet.subnet03.subnet_id
	public_ip_id = outscale_public_ip.public_ip01.public_ip_id
}

resource "outscale_route_table" "route_table01" {
	net_id = outscale_net.net04.net_id
}

resource "outscale_route" "route01" {
	destination_ip_range = "0.0.0.0/0"
	gateway_id           = outscale_internet_service.internet_service01.internet_service_id
	route_table_id       = outscale_route_table.route_table01.route_table_id
}

resource "outscale_route_table_link" "route_table_link01" {
	subnet_id      = outscale_subnet.subnet03.subnet_id
	route_table_id = outscale_route_table.route_table01.id
}

resource "outscale_internet_service" "internet_service01" {
}

resource "outscale_internet_service_link" "internet_service_link01" {
	net_id              = outscale_net.net04.net_id
	internet_service_id = outscale_internet_service.internet_service01.id
}


```

## Argument Reference

The following arguments are supported:

* `ip_range` - (Required) The IP range for the Net, in CIDR notation (for example, 10.0.0.0/16).
* `tenancy` - (Optional) The tenancy options for the VMs (`default` if a VM created in a Net can be launched with any tenancy, `dedicated` if it can be launched with dedicated tenancy VMs running on single-tenant hardware).

## Attribute Reference

The following attributes are exported:

* `net` - Information about the Net.
  * `dhcp_options_set_id` - The ID of the DHCP options set (or `default` if you want to associate the default one).
  * `ip_range` - The IP range for the Net, in CIDR notation (for example, 10.0.0.0/16).
  * `net_id` - The ID of the Net.
  * `state` - The state of the Net (`pending` \| `available`).
  * `tags` - One or more tags associated with the Net.
    * `key` - The key of the tag, with a minimum of 1 character.
    * `value` - The value of the tag, between 0 and 255 characters.
  * `tenancy` - The VM tenancy in a Net.
