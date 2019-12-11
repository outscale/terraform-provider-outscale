---
layout: "outscale"
page_title: "3DS OUTSCALE: outscale_tag"
sidebar_current: "outscale-tag"
description: |-
  [Provides information about tags.]
---

# outscale_tag Data Source

Provides information about tags.
For more information on this resource, see the [User Guide](https://wiki.outscale.net/display/EN/About+Tags).
For more information on this resource actions, see the [API documentation](https://docs-beta.outscale.com/#3ds-outscale-api-tag).

## Example Usage

```hcl
[exemple de code]
```

## Argument Reference

The following arguments are supported:

* `filters` - One or more filters.
  * `keys` - (Optional) The keys of the tags that are assigned to the resources. You can use this filter alongside the `Values` filter. In that case, you filter the resources corresponding to each tag, regardless of the other filter.
  * `resource_ids` - (Optional) The IDs of the resources with which the tags are associated.
  * `resource_types` - (Optional) The resource type (`instance` \| `image` \| `volume` \| `snapshot` \| `public-ip` \| `security-group` \| `route-table` \| `network-interface` \| `vpc` \| `subnet` \| `network-link` \| `vpc-endpoint` \| `nat-gateway` \| `internet-gateway` \| `customer-gateway` \| `vpn-gateway` \| `vpn-connection` \| `dhcp-options` \| `task`).
  * `values` - (Optional) The values of the tags that are assigned to the resources. You can use this filter alongside the `TagKeys` filter. In that case, you filter the resources corresponding to each tag, regardless of the other filter.

## Attribute Reference

The following attributes are exported:

* `tags` - Information about one or more tags.
  * `key` - The key of the tag, with a minimum of 1 character.
  * `resource_id` - The ID of the resource.
  * `resource_type` - The type of the resource.
  * `value` - The value of the tag, between 0 and 255 characters.
