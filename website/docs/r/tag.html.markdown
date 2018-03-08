---
layout: "outscale"
page_title: "OUTSCALE: outscale_tag"
sidebar_current: "docs-outscale-resource-tag"
description: |-
  Adds or replaces one or more tags for one or more specified resources.
---

#outscale_tag

Adds or replaces one or more tags for one or more specified resources.
A tag consists of a key and a value. This combination must be unique for each resource.
Tags allow associating user data with resources. You can tag the following resources:

Volumes (vol-xxxxxxxx)

Snapshots (snap-xxxxxxxx)

OMIs (ami-xxxxxxxx)

Instances (i-xxxxxxxx)

Internet gateways (igw-xxxxxxxx)

Network interfaces (eni-xxxxxxxx)

Route tables (rtb-xxxxxxxx)

Subnets (subnet-xxxxxxxx)

VPCs (vpc-xxxxxxxx)

Customer gateways (cgw-xxxxxxxx)

VPN gateways (vgw-xxxxxxxx)

VPN connections (vpn-xxxxxxxx)

## Example Usage

```hcl
resource "outscale_vm" "foo" {
	image_id = "ami-8a6a0120"
	instance_type = "m1.small"
	tags {
		foo = "bar"
	}
}

resource "outscale_tag" "foo" {
	resource_ids = ["${outscale_vm.foo.id}"]
	tags {
		faz = "baz"
	}
}
```

## Argument Reference

The following arguments are supported:

* `dry_run` - If set to true, checks whether you have the required permissions to perform the action.
* `resource_ids` - One or more resource IDs.
* `filter` - One or more filters.
* `tags` - A list of tags to add to the specified resources.
* `max_results` - The maximum number of results that can be returned in a single page. You can use the NextToken attribute to request the next results pages. This value is between 5 and 1000. If you provide a value larger than 1000, only 1000 results are returned.
* `next_token` - The token to request the next results page.
* `tag_set` - Information about one or more tags.

See detailed information in [FCU Address](http://docs.outscale.com/api_fcu/definitions/Address.html#_api_fcu-address).
