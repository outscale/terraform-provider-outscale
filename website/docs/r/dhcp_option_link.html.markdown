---
layout: "outscale"
page_title: "OUTSCALE: outscale_dhcp_option_link"
sidebar_current: "docs-outscale-resource-dhcp-option-link"
description: |-
	Associates a DHCP options set with a specified Virtual Private Cloud (VPC).
---

# outscale_dhcp_option_link

Associates a DHCP options set with a specified Virtual Private Cloud (VPC).

## Example Usage

```hcl
resource "outscale_lin" "foo" {
	cidr_block = "10.1.0.0/16"
}

resource "outscale_dhcp_option" "foo" {}

resource "outscale_dhcp_option_link" "foo" {
	vpc_id = "${outscale_lin.foo.id}"
	dhcp_options_id = "${outscale_dhcp_option.foo.id}"
}
```

## Argument Reference

The following arguments are supported:

* `dhcp_options_id`	- (Required) The ID of the DHCP options set, or default if you do not want to associate any DHCP options with the VPC.
* `vpc_id`- (Required)	The ID of the VPC.

## Attributes

* `dhcp_options_id`	- The ID of the DHCP options set, or default if you do not want to associate any DHCP options with the VPC.
* `vpc_id`- The ID of the VPC.
* `request_id`	The ID of the request.


See detailed information in [Associates DHCP Option](http://docs.outscale.com/api_fcu/operations/Action_AssociateDhcpOptions_get.html#_api_fcu-action_associatedhcpoptions_get).

