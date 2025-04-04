---
layout: "outscale"
page_title: "OUTSCALE: outscale_internet_services"
sidebar_current: "outscale-internet-services"
description: |-
  [Provides information about Internet services.]
---

# outscale_internet_services Data Source

Provides information about Internet services.

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Internet-Services.html).  
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-internetservice).

## Example Usage

```hcl
data "outscale_internet_services" "internet_services01" {
    filter {
        name   = "tag_keys"
        values = ["env"]
    }
    filter {
        name   = "tag_values"
        values = ["prod", "test"]
    }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) A combination of a filter name and one or more filter values. You can specify this argument for as many filter names as you need. The filter name can be any of the following:
    * `internet_service_ids` - (Optional) The IDs of the internet services.
    * `link_net_ids` - (Optional) The IDs of the Nets the internet services are attached to.
    * `link_states` - (Optional) The current states of the attachments between the internet services and the Nets (only `available`, if the internet gateway is attached to a Net).
    * `tag_keys` - (Optional) The keys of the tags associated with the internet services.
    * `tag_values` - (Optional) The values of the tags associated with the internet services.
    * `tags` - (Optional) The key/value combinations of the tags associated with the Internet services, in the following format: `TAGKEY=TAGVALUE`.

## Attribute Reference

The following attributes are exported:

* `internet_services` - Information about one or more internet services.
    * `internet_service_id` - The ID of the internet service.
    * `net_id` - The ID of the Net attached to the internet service.
    * `state` - The state of the attachment of the internet service to the Net (always `available`).
    * `tags` - One or more tags associated with the internet service.
        * `key` - The key of the tag, with a minimum of 1 character.
        * `value` - The value of the tag, between 0 and 255 characters.
