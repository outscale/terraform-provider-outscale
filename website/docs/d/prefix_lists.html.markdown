---
layout: "outscale"
page_title: "OUTSCALE: outscale_prefix_lists"
sidebar_current: "docs-outscale-datasource-prefix-lists"
description: |-
  Describes Prefix List through Terraform.
---

# outscale_prefix_lists

Describes Prefix List through Terraform.

## Example Usage

```hcl
data "outscale_prefix_lists" "s3_by_id" {
  prefix_list_id = ["pl-a14a8cdc"]
}
```

## Argument Reference

The following arguments are supported:

* `filter` -One or more filters. 
* `prefix_list_id` -One or more prefix list IDs

See detailed information in [Outscale Prefix List](https://wiki.outscale.net/display/DOCU/Getting+Information+About+Your+Instances).

## Filters

Use the Filter.N parameter to filter the described Prefix List on the following properties:

* `prefix-list-id` - The ID of a prefix list.
* `prefix-list-name` - The name of a prefix list.

## Attributes Reference

The following attributes are exported:

* `prefix_list_set.N`	- Information about one or more prefix lists, each containing the following attributes:
  - `cidr_et`	- The list of network prefixes used by the service, in CIDR notation.
  - `prefix_list_id` - The ID of the prefix list.	
  - `prefix_list_name` - The name of the prefix list, which identifies the Outscale service it is associated with.
  - `request_id`-	The ID of the request	false	string

See detailed information in [Describe Prefix List](http://docs.outscale.com/api_fcu/operations/Action_DescribePrefixLists_get.html#_api_fcu-action_describeprefixlists_get).