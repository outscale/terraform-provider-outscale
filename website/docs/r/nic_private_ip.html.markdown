---
layout: "outscale"
page_title: "3DS OUTSCALE: outscale_nic_private_ip"
sidebar_current: "outscale-nic-private-ip"
description: |-
  [Manages a NIC private IP.]
---

# outscale_nic_private_ip Resource

Manages a NIC private IP.
For more information on this resource, see the [User Guide](https://wiki.outscale.net/display/EN/About+FNIs#AboutFNIs-FNIsAttributes).
For more information on this resource actions, see the [API documentation](https://docs-beta.outscale.com/#linkprivateips).

## Example Usage

```hcl

#resource "outscale_net" "net01" {
#	ip_range = "10.0.0.0/16"
#}

#resource "outscale_subnet" "subnet01" {
#	subregion_name = "${var.region}a"
#	ip_range       = "10.0.0.0/16"
#	net_id         = outscale_net.net01.net_id
#}

#resource "outscale_nic" "nic01" {
#	subnet_id = outscale_subnet.subnet01.subnet_id
#}

resource "outscale_nic_private_ip" "nic_private_ip01" {
	nic_id      = outscale_nic.nic01.nic_id
	private_ips = ["10.0.12.34"]
}
resource "outscale_nic_private_ip" "nic_private_ip02" {
	nic_id                     = outscale_nic.nic01.nic_id
	secondary_private_ip_count = 2
}


```

## Argument Reference

The following arguments are supported:

* `allow_relink` - (Optional) If `true`, allows an IP address that is already assigned to another NIC in the same Subnet to be assigned to the NIC you specified.
* `nic_id` - (Required) The ID of the NIC.
* `private_ips` - (Optional) The secondary private IP address or addresses you want to assign to the NIC within the IP address range of the Subnet.
* `secondary_private_ip_count` - (Optional) The number of secondary private IP addresses to assign to the NIC.

## Attribute Reference

No attribute is exported.