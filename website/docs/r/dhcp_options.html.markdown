---
layout: "outscale"
page_title: "OUTSCALE: outscale_dhcp_options"
sidebar_current: "docs-outscale-resource-dhcp-options"
description: |-
  Provides an Outscale DHCP Options resource. Creates, Delete, Describe and Import a new set of DHCP options that you can then associate to a Virtual Private Cloud (VPC).
---

# outscale_dhcp_options

  Provides an Outscale DHCP Options resource. Creates, Delete, Describe and Import a new set of DHCP options that you can then associate to a Virtual Private Cloud (VPC). Instances also support [provisioning](/docs/provisioners/index.html).

## Example Usage

```hcl
data "outscale_dhcp_options" "outscale_dhcp_options" {
  most_recent = true

  filter {
    name   = "volume-type"
    values = ["gp2"]
  }

  filter {
    name   = "tag:Name"
    values = ["Example"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `dhcp_configuration` - (Optional) A DHCP configuration option.

See detailed information in [Outscale Instances](https://wiki.outscale.net/display/DOCU/Getting+Information+About+Your+Instances).


## Attributes Reference

The following attributes are exported:

* `dhcp_configuration_set` - One or more DHCP options in the set.
* `dhcp_options_id` - The ID of the DHCP options set.
* `tag_set` - One or more tags associated with the DHCP options set.
* `request_id` - The ID of the request


See detailed information in [Describe DHCP Options](http://docs.outscale.com/api_fcu/definitions/Volume.html#_api_fcu-volume).