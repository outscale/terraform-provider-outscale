---
layout: "outscale"
page_title: "OUTSCALE: outscale_tags"
sidebar_current: "docs-outscale-datasource-tags"
description: |-
  Describes one or more tags for your resources.
  You can use the Filter.N parameter to filter the tags on the following properties
---

# outscale_tags

Describes one or more tags for your resources.

## Example Usage

```hcl
resource "outscale_vm" "basic" {
  image_id = "ami-8a6a0120"
	instance_type = "m1.small"
	tags = {
		foo = "bar"
	}
}
resource "outscale_vm" "basic2" {
  image_id = "ami-8a6a0120"
	instance_type = "m1.small"
	tags = {
		foo = "baz"
	}
}

data "outscale_tags" "web" {
	filter {
    name = "resource-type"
    value = ["instance"]
	}
}
```

## Argument Reference

The following arguments are supported:

* `Filter.N` - (Optional) One or more filters.

See detailed information in [Outscale Tags](http://docs.outscale.com/api_fcu/operations/Action_DescribeTags_get.html#_api_fcu-action_describetags_get).

## Filters

Use the Filter.N parameter to filter the described instances on the following properties:

* `key` The key of the tag.
* `value` The value of the tag.
* `resource-id` The ID of the tag.
* `resource-type` The resource type (volume | snapshot | image | instance | vpc | internet-gateway | subnet | route-table | network-interface | vpn-gateway | customer-gateway | vpn-connection).

## Attributes Reference

The following attributes are exported:

* `request_id` - The ID of the request.
* `tagSet.N` - Information about attribute list.

See detailed information in [Describe Tags](http://docs.outscale.com/api_fcu/operations/Action_DescribeTags_get.html#_body_parameter).
