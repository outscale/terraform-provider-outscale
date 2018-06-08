---
layout: "outscale"
page_title: "OUTSCALE: outscale_lin_api_access"
sidebar_current: "docs-outscale-datasource-lin-api-access"
description: |-
  Creates a Virtual Private Cloud (VPC) endpoint to access an Outscale service from this VPC without using the Internet and External IP addresses.
  You specify the service using its prefix list name.
---

# outscale_lin_api_access

Creates a Virtual Private Cloud (VPC) endpoint to access an Outscale service from this VPC without using the Internet and External IP addresses.
You specify the service using its prefix list name.

## Example Usage

```hcl
resource "outscale_lin" "foo" {
    cidr_block = "10.1.0.0/16"
}

resource "outscale_route_table" "foo" {
    vpc_id = "${outscale_lin.foo.id}"
}

resource "outscale_lin_api_access" "link" {
    vpc_id = "${outscale_lin.foo.id}"
    route_table_id = [
        "${outscale_route_table.foo.id}"
    ]
    service_name = "com.outscale.eu-west-2.osu"
}

data "outscale_lin_api_access" "test" {
    vpc_endpoint_id = "${outscale_lin_api_access.link.id}"
}
```

## Argument Reference

The following arguments are supported:

* `vpc_enpoint_id` - (Optional) The ID of the VPC Endpoint.
* `filter.N` - (Optional) One or more filters.

## Attributes

* `service_name`: - The prefix list name corresponding to the service (for example, com.outscale.eu-west-2.osu for OSU).
* `route_table_id` - One or more IDs of route tables to use for the connection.
* `prefix_list_id` - The Prefix ID List for the service given.
* `cidr_blocks` - The CIDR Blocks for Prefix Ids
* `request_id` - The ID of the request.
* `state` - The state of the VPC endpoint (pending| available| deleting | deleted).

See detailed information in [VPC Endpoints](http://docs.outscale.com/api_fcu/index.html#_vpc_endpoints).
