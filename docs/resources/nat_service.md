---
layout: "outscale"
page_title: "OUTSCALE: outscale_nat_service"
sidebar_current: "outscale-nat-service"
description: |-
  [Manages a NAT service.]
---

# outscale_nat_service Resource

Manages a NAT service.
For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-NAT-Gateways.html).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-natservice).

## Example Usage

### Required resources

```hcl
resource "outscale_net" "net01" {
    ip_range = "10.0.0.0/16"
}

resource "outscale_subnet" "subnet01" {
    net_id   = outscale_net.net01.net_id
    ip_range = "10.0.0.0/18"
}

resource "outscale_route_table" "route_table01" {
    net_id = outscale_net.net01.net_id
}

resource "outscale_route_table_link" "outscale_route_table_link01" {
    subnet_id      = outscale_subnet.subnet01.subnet_id
    route_table_id = outscale_route_table.route_table01.route_table_id
}

resource "outscale_internet_service" "internet_service01" {
}

resource "outscale_internet_service_link" "internet_service_link01" {
    net_id              = outscale_net.net01.net_id
    internet_service_id = outscale_internet_service.internet_service01.internet_service_id
}

resource "outscale_route" "route01" {
    destination_ip_range = "0.0.0.0/0"
    gateway_id           = outscale_internet_service.internet_service01.internet_service_id
    route_table_id       = outscale_route_table.route_table01.route_table_id
    depends_on           = [outscale_internet_service_link.internet_service_link01]
}

resource "outscale_public_ip" "public_ip01" {
}
```

### Create a NAT service

```hcl
resource "outscale_nat_service" "nat_service01" {
    subnet_id    = outscale_subnet.subnet01.subnet_id
    public_ip_id = outscale_public_ip.public_ip01.public_ip_id
    depends_on   = [outscale_route.route01]
}
```

## Argument Reference

The following arguments are supported:

* `public_ip_id` - (Required) The allocation ID of the public IP to associate with the NAT service.<br />
If the public IP is already associated with another resource, you must first disassociate it.
* `subnet_id` - (Required) The ID of the Subnet in which you want to create the NAT service.
* `tags` - (Optional) A tag to add to this resource. You can specify this argument several times.
    * `key` - (Required) The key of the tag, with a minimum of 1 character.
    * `value` - (Required) The value of the tag, between 0 and 255 characters.

## Attribute Reference

The following attributes are exported:

* `nat_service_id` - The ID of the NAT service.
* `net_id` - The ID of the Net in which the NAT service is.
* `public_ips` - Information about the public IP or IPs associated with the NAT service.
    * `public_ip` - The public IP associated with the NAT service.
    * `public_ip_id` - The allocation ID of the public IP associated with the NAT service.
* `state` - The state of the NAT service (`pending` \| `available` \| `deleting` \| `deleted`).
* `subnet_id` - The ID of the Subnet in which the NAT service is.
* `tags` - One or more tags associated with the NAT service.
    * `key` - The key of the tag, with a minimum of 1 character.
    * `value` - The value of the tag, between 0 and 255 characters.

## Import

A NAT service can be imported using its ID. For example:

```console

$ terraform import outscale_nat_service.ImportedNatService nat-87654321

```