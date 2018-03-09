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
    values = ["instance"]
	}
}
```

## Argument Reference

The following arguments are supported:

* `dryRun` - (Optional) If set to true, checks whether you have the required permissions to perform the action.
* `Filter.N` - (Optional) One or more filters.
* `MaxResults` - (Optional) The maximum number of results that can be returned in a single page. You can use the NextToken attribute to request the next results pages. This value is between 5 and 1000. If you provide a value larger than 1000, only 1000 results are returned.
* `NextToken` - (Optional) The token to request the next results page.



See detailed information in [Outscale Instances](http://docs.outscale.com/api_fcu/operations/Action_DescribeTags_get.html#_api_fcu-action_describetags_get).

## Filters

Use the Filter.N parameter to filter the described instances on the following properties:

* `key` The key of the tag.
* `value` The value of the tag.
* `resource-id` The ID of the tag.
* `resource-type` The resource type (volume | snapshot | image | instance | vpc | internet-gateway | subnet | route-table | network-interface | vpn-gateway | customer-gateway | vpn-connection).

## Attributes Reference

The following attributes are exported:

* `request_id` - The ID of the request.
* `tagSet.N` - Information about one or more tags.

See detailed information in [Describe Instances](http://docs.outscale.com/api_fcu/operations/Action_DescribeTags_get.html#_body_parameter).
