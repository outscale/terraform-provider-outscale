---
layout: "outscale"
page_title: "Outscale: outscale_directlink"
sidebar_current: "docs-outscale-resource-directlink"
description: |-
  Provides a Connection of Direct Connect.
---
# outscale_directlink

Provides a Connection of Direct Connect.

## Example Usage

```hcl
data "outscale_sites" "test" {}

resource "outscale_directlink" "hoge" {
  bandwidth = "1Gbps"
  connection_name = "test-directlink-%d"
  location = "${data.outscale_sites.test.locations.0.location_code}"
}

data "outscale_directlink" "test" {
  connection_id = "${outscale_directlink.hoge.id}"
}
```

## Argument Reference

The following arguments are supported:

* `connection_id` - (Required) The name of the connection.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `connections.N` - Information about one or more DirectLink connections.
  * `bandwidth` - The physical link bandwidth (either 1 GiB/s or 10 GiB/s).
  * `connection_id` - The ID of the connection (for example, dcx-xxxxxxxx).
  * `connection_name` - The name of the connection.
  * `connection_state` - The state of the connection. Connection states are: `requested`: The connection is requested but the request has not been validated yet. `pending`: The connection request has been validated. It remains in the pending state until you establish the physical link. available: The physical link is established and the connection is ready to use. `deleting`: The deletion process is in progress. `deleted`: The connection is deleted.
  * `location` - The datacenter where the connection is located.
  * `owner_account` - The account ID of the owner of the connection.
  * `region` - The Region in which the connection has been created.
* `request_id` - The ID of the request.

[See detailed description](http://docs.outscale.com/api_directlink/operations/Action_DescribeConnections_get.html#_api_directlink-action_describeconnections_get)