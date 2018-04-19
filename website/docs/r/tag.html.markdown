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
	tag {
		faz = "baz"
	}
}
```

## Argument Reference

The following arguments are supported:

* `resource_ids` - One or more resource IDs.
* `filter` - One or more filters.
* `tag` - A list of tags to add to the specified resources.

See detailed information in [FCU Address](http://docs.outscale.com/api_fcu/definitions/Address.html#_api_fcu-address).
