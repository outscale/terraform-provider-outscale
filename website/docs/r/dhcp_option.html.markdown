---
layout: "outscale"
page_title: "OUTSCALE: outscale_dhcp_option"
sidebar_current: "docs-outscale-resource-dhcp-option"
description: |-
	Creates a new set of DHCP options that you can then associate to a Virtual Private Cloud (VPC).
---

# outscale_dhcp_option

Creates a new set of DHCP options that you can then associate to a Virtual Private Cloud (VPC).

## Example Usage

```hcl
resource "outscale_dhcp_option" "foo" {}
```

## Argument Reference

The following arguments are supported:

* `DhcpConfiguration` - (Optional)	A DHCP configuration option.

## Attributes

* `dhcp_configurationSet.N` -	One or more DHCP options in the set.
* `dhcp_optionsId` - The ID of the DHCP options set.
* `tag_set` - One or more tags associated with the DHCP options set.
* `request_id` - The ID of the request.


See detailed information in [Create DHCP Option](http://docs.outscale.com/api_fcu/operations/Action_DescribeDhcpOptions_get.html#_api_fcu-action_describedhcpoptions_get).

