---
layout: "outscale"
page_title: "OUTSCALE: outscale_tag"
sidebar_current: "docs-outscale-datasource-tag"
description: |-
  Describes one tag for your resources.
---

# outscale_tag

Describes one tag for your resources.

## Example Usage

```hcl
resource "outscale_vm" "basic" {
  image_id = "ami-8a6a0120"
	instance_type = "m1.small"
	tags = {
		foo = "bar"
	}
}

data "outscale_tag" "web" {
	filter {
    name = "resource-id"
    values = ["${outscale_vm.basic.id}"]
	}
}
```

## Argument Reference

The following arguments are supported:

* `dryRun` - (Optional) If set to true, checks whether you have the required permissions to perform the action.
* `filter` - (Optional) One or more filters.

See detailed information in [Outscale Instances](http://docs.outscale.com/api_fcu/operations/Action_DescribeTags_get.html#_api_fcu-action_describetags_get).

## Filters

Use the Filter.N parameter to filter the described instances on the following properties:

* `key` The key of the tag.
* `value` The value of the tag.
* `resource-id` The ID of the tag.
* `resource-type` The resource type (volume | snapshot | image | instance | vpc | internet-gateway | subnet | route-table | network-interface | vpn-gateway | customer-gateway | vpn-connection).

## Attributes Reference

The following attributes are exported:

* `tag_set` - Information about one.
* `request_id` - The ID of the request.

See detailed information in [Describe Instances](http://docs.outscale.com/api_fcu/operations/Action_DescribeTags_get.html#_body_parameter).
