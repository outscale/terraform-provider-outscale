---
layout: "outscale"
page_title: "OUTSCALE: outscale_directlink"
sidebar_current: "docs-outscale-resource-directlink"
description: |-
  Creates a new DirectLink connection between a customer network and a specified DirectLink location.
---

# outscale_directlink

Creates a new DirectLink connection between a customer network and a specified DirectLink location.

## Example Usage

```hcl
resource "outscale_directlink" "hoge" {
    bandwidth = "1Gbps"
    connection_name = "test-directlink-tf-dx-2"
    location = "PAR1"
}
```

## Argument Reference

The following arguments are supported:

* `bandwidth` - The bandwidth of the connection (1G GiB/s | 10 GiB/s).
* `connection_name` - The name of the connection.
* `location` - The code of the requested location for the connection, returned by the DescribeLocation method.

## Attributes Reference

The following attributes are exported:

* `bandwidth` - The bandwidth of the connection (1G GiB/s | 10 GiB/s).
* `connection_id` - The ID of the connection (for example, dcx-xxxxxxxx).
* `connection_name` - The name of the connection.
* `connection_state` - The state of the connection.
* `location` - The datacenter where the connection is located.
* `owner_account` - The account ID of the owner of the connection.
* `region` - The Region in which the connection has been created.
* `request_id` - The ID of the request

[See detailed description](http://docs.outscale.com/api_directlink/operations/Action_CreateConnection_get.html#_api_directlink-action_createconnection_get)
