---
layout: "outscale"
page_title: "3DS OUTSCALE: outscale_tag"
sidebar_current: "docs-outscale-resource-tag"
description: |-
  [Manages a tag.]
---

# outscale_tag Resource

Manages a tag.
For more information on this resource, see the [User Guide](https://wiki.outscale.net/display/EN/About+Tags).
For more information on this resource actions, see the [API documentation](https://docs-beta.outscale.com/#3ds-outscale-api-tag).

## Example Usage

```hcl
[exemple de code]
```

## Argument Reference

The following arguments are supported:

* `resource_ids` - (Required) One or more resource IDs.
* `tags` - (Required) One or more tags to add to the specified resources.
  * `key` - (Optional) The key of the tag, with a minimum of 1 character.
  * `value` - (Optional) The value of the tag, between 0 and 255 characters.

## Attribute Reference

The following attributes are exported:

