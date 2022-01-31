---
layout: "outscale"
page_title: "OUTSCALE: outscale_nic_private_ip"
sidebar_current: "outscale-nic-private-ip"
description: |-
  [Manages a NIC's private IPs.]
---

# outscale_nic_private_ip Resource

Manages a NIC's private IPs.
For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-FNIs.html).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-nic).

## Example Usage

### Required resources

```hcl
resource "outscale_net" "net01" {
	ip_range = "10.0.0.0/16"
}

resource "outscale_subnet" "subnet01" {
	subregion_name = "${var.region}a"
	ip_range       = "10.0.0.0/16"
	net_id         = outscale_net.net01.net_id
}

resource "outscale_nic" "nic01" {
	subnet_id = outscale_subnet.subnet01.subnet_id
}
```

### Link a specific secondary private IP address to a NIC

```hcl
resource "outscale_nic_private_ip" "nic_private_ip01" {
	nic_id      = outscale_nic.nic01.nic_id
	private_ips = ["10.0.12.34", "10.0.12.35"]
}
```

### Link several automatic secondary private IP addresses to a NIC

```hcl
resource "outscale_nic_private_ip" "nic_private_ip02" {
	nic_id                     = outscale_nic.nic01.nic_id
	secondary_private_ip_count = 2
}
```

## Argument Reference

The following arguments are supported:

* `allow_relink` - (Optional) If true, allows an IP address that is already assigned to another NIC in the same Subnet to be assigned to the NIC you specified.
* `nic_id` - (Required) The ID of the NIC.
* `private_ips` - (Optional) The secondary private IP address or addresses you want to assign to the NIC within the IP address range of the Subnet.
* `secondary_private_ip_count` - (Optional) The number of secondary private IP addresses to assign to the NIC.

## Attribute Reference

No attribute is exported.

