---
layout: "outscale"
page_title: "3DS OUTSCALE: outscale_tag"
sidebar_current: "outscale-tag"
description: |-
  [Manages a tag.]
---

# outscale_tag Resource

Manages a tag.
For more information on this resource, see the [User Guide](https://wiki.outscale.net/display/EN/About+Tags).
For more information on this resource actions, see the [API documentation](https://docs-beta.outscale.com/#3ds-outscale-api-tag).

## Example Usage

```hcl

resource "outscale_tag" "tag01" {
	resource_ids = [var.vm_id]
	tag {
		key = "name"
		value = "terraform-vm-with-tag"
	}
}


```

## Argument Reference

The following arguments are supported:

* `resource_ids` - (Required) One or more resource IDs.
* `tag` - (Required) One or more tags to add to the specified resources.
  * `key` - (Optional) The key of the tag, with a minimum of 1 character.
  * `value` - (Optional) The value of the tag, between 0 and 255 characters.

## Attribute Reference

No attribute is exported.

