---
layout: "outscale"
page_title: "OUTSCALE: outscale_net_access_points"
sidebar_current: "outscale-net-access-points"
description: |-
  [Provides information about Net access points.]
---

# outscale_net_access_points Data Source

Provides information about Net access points.
For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-VPC-Endpoints.html).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-netaccesspoint).

## Example Usage

### List Net access points

```hcl
data "outscale_net_access_points" "net_access_points01" {
    filter {
        name   = "net_access_point_ids"
        values = ["vpce-12345678","vpce-12345679"]
    }
}
```

### List Net access points according to their Net and state

```hcl
data "outscale_net_access_points" "net_access_points02" {
    filter {
        name   = "net_ids"
        values = ["vpc-12345678"]
    }
    filter {
        name   = "states"
        values = ["available"]
    }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) A combination of a filter name and one or more filter values. You can specify this argument for as many filter names as you need. The filter name can be any of the following:
    * `net_access_point_ids` - (Optional) The IDs of the Net access points.
    * `net_ids` - (Optional) The IDs of the Nets.
    * `service_names` - (Optional) The names of the services. For more information, see [ReadNetAccessPointServices](https://docs.outscale.com/api#readnetaccesspointservices).
    * `states` - (Optional) The states of the Net access points (`pending` \| `available` \| `deleting` \| `deleted`).
    * `tag_keys` - (Optional) The keys of the tags associated with the Net access points.
    * `tag_values` - (Optional) The values of the tags associated with the Net access points.
    * `tags` - (Optional) The key/value combination of the tags associated with the Net access points, in the following format: &quot;Filters&quot;:{&quot;Tags&quot;:[&quot;TAGKEY=TAGVALUE&quot;]}.

## Attribute Reference

The following attributes are exported:

* `net_access_points` - One or more Net access points.
    * `net_access_point_id` - The ID of the Net access point.
    * `net_id` - The ID of the Net with which the Net access point is associated.
    * `route_table_ids` - The ID of the route tables associated with the Net access point.
    * `service_name` - The name of the service with which the Net access point is associated.
    * `state` - The state of the Net access point (`pending` \| `available` \| `deleting` \| `deleted`).
    * `tags` - One or more tags associated with the Net access point.
        * `key` - The key of the tag, with a minimum of 1 character.
        * `value` - The value of the tag, between 0 and 255 characters.
