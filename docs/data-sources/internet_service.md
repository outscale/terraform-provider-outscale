---
layout: "outscale"
page_title: "OUTSCALE: outscale_internet_service"
sidebar_current: "outscale-internet-service"
description: |-
  [Provides information about a specific Internet service.]
---

# outscale_internet_service Data Source

Provides information about a specific Internet service.
For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Internet-Gateways.html).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-internetservice).

## Example Usage

```hcl
data "outscale_internet_service" "internet_service01" {
  filter {
    name   = "internet_service_ids"
    values = ["igw-12345678"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) A combination of a filter name and one or more filter values. You can specify this argument for as many filter names as you need. The filter name can be any of the following:
    * `internet_service_ids` - (Optional) The IDs of the Internet services.
    * `link_net_ids` - (Optional) The IDs of the Nets the Internet services are attached to.
    * `link_states` - (Optional) The current states of the attachments between the Internet services and the Nets (only `available`, if the Internet gateway is attached to a VPC).
    * `tag_keys` - (Optional) The keys of the tags associated with the Internet services.
    * `tag_values` - (Optional) The values of the tags associated with the Internet services.
    * `tags` - (Optional) The key/value combination of the tags associated with the Internet services, in the following format: &quot;Filters&quot;:{&quot;Tags&quot;:[&quot;TAGKEY=TAGVALUE&quot;]}.

## Attribute Reference

The following attributes are exported:

* `internet_service_id` - The ID of the Internet service.
* `net_id` - The ID of the Net attached to the Internet service.
* `state` - The state of the attachment of the Internet service to the Net (always `available`).
* `tags` - One or more tags associated with the Internet service.
    * `key` - The key of the tag, with a minimum of 1 character.
    * `value` - The value of the tag, between 0 and 255 characters.
