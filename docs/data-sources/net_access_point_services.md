---
layout: "outscale"
page_title: "OUTSCALE: outscale_net_access_point_services"
sidebar_current: "outscale-net-access-point-services"
description: |-
  [Provides information about Net access point services.]
---

# outscale_net_access_point_services Data Source

Provides information about Net access point services.
For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-VPC-Endpoints.html).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-netaccesspoint).

## Example Usage

### List all services available to create Net access points

```hcl
data "outscale_net_access_point_services" "all-services" { 
}
```

### List one or more services according to their service IDs

```hcl
data "outscale_net_access_point_services" "services01" {
  filter {
    name   = "service_ids"
    values = ["pl-12345678","pl-12345679"]
  }
}
```

### List one or more services according to their service names

```hcl
data "outscale_net_access_point_services" "services02" {
  filter {
    name   = "service_names"
    values = ["com.outscale.eu-west-2.api"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) A combination of a filter name and one or more filter values. You can specify this argument for as many filter names as you need. The filter name can be any of the following:
    * `service_ids` - (Optional) The IDs of the services.
    * `service_names` - (Optional) The names of the services.

## Attribute Reference

The following attributes are exported:

* `services` - The names of the services you can use for Net access points.
    * `ip_ranges` - The list of network prefixes used by the service, in CIDR notation.
    * `service_id` - The ID of the service.
    * `service_name` - The name of the service.
