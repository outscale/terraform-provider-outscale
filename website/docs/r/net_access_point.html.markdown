---
layout: "outscale"
page_title: "OUTSCALE: outscale_net_access_point"
sidebar_current: "outscale-net-access-point"
description: |-
  [Manages a Net access point.]
---

# outscale_net_access_point Resource

Manages a Net access point.
For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-VPC-Endpoints.html).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-netaccesspoint).

## Example Usage

### Required resources

```hcl
resource "outscale_net" "net01" { 
  ip_range = "10.0.0.0/16"
}

resource "outscale_route_table" "route_table01" {
  net_id = outscale_net.net01.net_id
}
```

### Create a Net access point

```hcl
resource "outscale_net_access_point" "net_access_point01" {
   net_id          = outscale_net.net01.net_id
   route_table_ids = [outscale_route_table.route_table01.route_table_id]
   service_name    = "com.outscale.eu-west-2.api"
   tags {
      key   = "name"
      value = "terraform-net-access-point"
   }
}
```

## Argument Reference

The following arguments are supported:

* `net_id` - (Required) The ID of the Net.
* `route_table_ids` - (Optional) One or more IDs of route tables to use for the connection.
* `service_name` - (Required) The name of the service (in the format `com.outscale.region.service`).
* `tags` - (Optional) A tag to add to this resource. You can specify this argument several times.
    * `key` - (Required) The key of the tag, with a minimum of 1 character.
    * `value` - (Required) The value of the tag, between 0 and 255 characters.

## Attribute Reference

The following attributes are exported:

* `net_access_point_id` - The ID of the Net access point.
* `net_id` - The ID of the Net with which the Net access point is associated.
* `route_table_ids` - The ID of the route tables associated with the Net access point.
* `service_name` - The name of the service with which the Net access point is associated.
* `state` - The state of the Net access point (`pending` \| `available` \| `deleting` \| `deleted`).
* `tags` - One or more tags associated with the Net access point.
    * `key` - The key of the tag, with a minimum of 1 character.
    * `value` - The value of the tag, between 0 and 255 characters.

## Import

A Net access point can be imported using its ID. For example:

```console

$ terraform import outscale_net_access_point.ImportedNetAccessPoint vpce-87654321

```