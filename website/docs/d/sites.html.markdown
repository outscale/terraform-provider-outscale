---
layout: "outscale"
page_title: "OUTSCALE: outscale_sites"
sidebar_current: "docs-outscale-datasource-sites"
description: |-
  Describes the locations, corresponding to datacenters, where you can set up a DirectLink connection.
---

# outscale_sites

Describes the locations, corresponding to datacenters, where you can set up a DirectLink connection.

## Example Usage

```hcl
data "outscale_sites" "test" {}
```

## Argument Reference

No arguments are supported

## Attributes Reference

The following attributes are exported:

* `locations.N` - Information about one or more locations.
  * `location_code` - The location code, to be set as the Location parameter of the CreateConnection method when creating a connection.
  * `location_name` - The name and description of the location, corresponding to a datacenter.
* `request_id` - The ID of the request.

[See detailed description](http://docs.outscale.com/api_directlink/operations/Action_DescribeLocations_get.html#_api_directlink-action_describelocations_get)
