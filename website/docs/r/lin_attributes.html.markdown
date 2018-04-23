---
layout: "outscale"
page_title: "OUTSCALE: outscale_lin_attributes"
sidebar_current: "docs-outscale-resource-lin-attributes"
description: |-
  Modifies a specified attribute of a VPC. You can modify one attribute only at a time.
---

# outscale_lin_attibutes

Modifies a specified attribute of a VPC. You can modify one attribute only at a time.

## Example Usage

```hcl
resource "outscale_lin" "vpc" {
	cidr_block = "10.0.0.0/16"
}

resource "outscale_lin_attributes" "outscale_lin_attributes" {
	enable_dns_hostnames = true
	vpc_id = "${outscale_lin.vpc.id}"
	attribute = "enableDnsSupport"
}
```

## Argument Reference

The following arguments are supported:

* `enable_dns_hostnames` (Optional)	If set to true, the instances launched in the VPC get DNS hostnames.
* `enable_dns_support`	(Optional) If set to true, the DNS resolution is supported for the VPC.	
* `vpc_id`	(Required) The ID of the VPC.

## Attributes

* `enable_dns_hostnames` - If set to true, the instances launched in the VPC get DNS hostnames.
* `enable_dns_support` - If set to true, the DNS resolution is supported for the VPC.	
* `vpc_id` - The ID of the VPC.


See detailed information in [Modify Lin Attributes](http://docs.outscale.com/api_fcu/operations/Action_ModifyVpcAttribute_get.html#_api_fcu-action_modifyvpcattribute_get).
