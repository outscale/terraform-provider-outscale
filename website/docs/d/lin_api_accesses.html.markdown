---
layout: "outscale"
page_title: "OUTSCALE: outscale_lin_api_access"
sidebar_current: "docs-outscale-datasource-lin-api-accesses"
description: |-
Describes one or more Virtual Private Cloud (VPC) endpoints.
---

# outscale_lin_api_accesses

Describes one or more Virtual Private Cloud (VPC) endpoints.

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

data "outscale_lin_api_accesses" "test" {
    filter {
        name = "service-name"
        values = ["${outscale_lin_api_access.link.service_name}"]
    }
}
```

## Argument Reference

The following arguments are supported:

* `vpc_enpoint_id` - (Optional) The IDs of the VPC Endpoint.
* `filter` - (Optional) One or more filters.

#### Filters

Use the filter (Filter.N) parameter to filter the described instances on the following properties:

* `service-name` - The name of the prefix list corresponding to the service. For more information, see DescribePrefixLists.
* `vpc-id` - The ID of the VPC.
* `vpc-endpoint-id` The ID of the VPC endpoint.
* `vpc-endpoint-state` The state of the VPC endpoint (pending | available | deleting | deleted).


## Attributes

* `vpc_endpoint_set`: - Information about one or more VPC endpoints.
* `request_id` - The ID of the request.

### Vpc Endpoint Set

The vpc_endpoint_set element has the following fields:

* `service_name`: - The prefix list name corresponding to the service (for example, com.outscale.eu-west-2.osu for OSU).
* `route_table_id` - One or more IDs of route tables to use for the connection.
* `prefix_list_id` - The Prefix ID List for the service given.
* `cidr_blocks` - The CIDR Blocks for Prefix Ids
* `state` - The state of the VPC endpoint (pending| available| deleting | deleted).

See detailed information in [VPC Endpoints](http://docs.outscale.com/api_fcu/index.html#_vpc_endpoints).
