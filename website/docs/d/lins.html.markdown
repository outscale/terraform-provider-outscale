---
layout: "outscale"
page_title: "OUTSCALE: outscale_lins"
sidebar_current: "docs-outscale-datasource-lins"
description: |-
    Describes one or more Virtual Private Clouds (VPCs)

---

# outscale_lin

Describes one or more Virtual Private Clouds (VPCs).
You can use the Filter.N parameter to filter the VPCs on the following properties:

## Example Usage

```hcl
data "outscale_lins" "by_id" {
  vpc_id = ["${outscale_lin.test.id}"]
}`, cidr, tag)
```

## Argument Reference

The following arguments are supported:

* `filter.N` - (Optional) One or more filters.
* `vpc_id.N` - (Optional) One or more VPC IDs.





See detailed information in [Outscale lins](http://docs.outscale.com/api_fcu/operations/Action_DescribeVpcs_get.html#_api_fcu-action_describevpcs_get).

## Filters

Use the Filter.N parameter to filter the described instances on the following properties:

* `vpc-id` - (Optional) The ID of the VPC.
* `state` - (Optional) The state of the VPC (pending | available).
* `cidr` - (Optional) The exact CIDR block of the VPC.
* `cidr-block` - (Optional) The exact CIDR block of the VPC (similar to cdr and cidrBlock).
* `cidrBlock` - (Optional) The exact CIDR block of the VPC (similar to cidr and cidr-block).
* `dhcp-options-id` - (Optional) The ID of a DHCP options.
* `is-default` - (Optional) Indicates Whether the VPC is the default one.
* `isDefault` - (Optional) Alias to is-default filter.
* `tag` - (Optional) The key/value combination of a tag associated with the resource.
* `tag-key` - (Optional) The key of a tag associated with the resource.
* `tag-value` - (Optional) The value of a tag associated with the resource.


## Attributes Reference

The following attributes are exported:

* `vpcSet` - Information about the specified and described VPCs.

See detailed information in [Describe lins](http://docs.outscale.com/api_fcu/operations/Action_DescribeVpcs_get.html#_api_fcu-action_describevpcs_get).
