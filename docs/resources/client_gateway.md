---
layout: "outscale"
page_title: "OUTSCALE: outscale_client_gateway"
subcategory: "Client Gateway"
sidebar_current: "outscale-client-gateway"
description: |-
  [Manages a client gateway.]
---

# outscale_client_gateway Resource

Manages a client gateway.

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Client-Gateways.html).  
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-clientgateway).

## Example Usage

```hcl
resource "outscale_client_gateway" "client_gateway01" {
    bgp_asn         = 65000
    public_ip       = "111.11.11.111"
    connection_type = "ipsec.1"
    tags {
        key   = "Name"
        value = "client_gateway_01"
    }
}
```

## Argument Reference

The following arguments are supported:

* `bgp_asn` - (Required) The Autonomous System Number (ASN) used by the Border Gateway Protocol (BGP) to find the path to your client gateway through the Internet. <br/>
This number must be between `1` and `4294967295`, except `50624`, `53306`, and `132418`. <br/>
If you do not have an ASN, you can choose one between `64512` and `65534` (both included), or between `4200000000` and `4294967295` (both included).
* `connection_type` - (Required) The communication protocol used to establish tunnel with your client gateway (always `ipsec.1`).
* `public_ip` - (Required) The public fixed IPv4 address of your client gateway.
* `tags` - (Optional) A tag to add to this resource. You can specify this argument several times.
    * `key` - (Required) The key of the tag, with a minimum of 1 character.
    * `value` - (Required) The value of the tag, between 0 and 255 characters.

## Attribute Reference

The following attributes are exported:

* `bgp_asn` - The Autonomous System Number (ASN) used by the Border Gateway Protocol (BGP) to find the path to your client gateway through the Internet.
* `client_gateway_id` - The ID of the client gateway.
* `connection_type` - The type of communication tunnel used by the client gateway (always `ipsec.1`).
* `public_ip` - The public IPv4 address of the client gateway (must be a fixed address into a NATed network).
* `state` - The state of the client gateway (`pending` \| `available` \| `deleting` \| `deleted`).
* `tags` - One or more tags associated with the client gateway.
    * `key` - The key of the tag, with a minimum of 1 character.
    * `value` - The value of the tag, between 0 and 255 characters.

## Import

A client gateway can be imported using its ID. For example:

```console

$ terraform import outscale_client_gateway.ImportedClientGateway cgw-12345678

```