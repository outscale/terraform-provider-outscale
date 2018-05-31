---
layout: "outscale"
page_title: "OUTSCALE: outscale_quotas"
sidebar_current: "docs-outscale-datasource-quotas"
description: |-
  Describes one or more of your quotas.
---

# outscale_quotas

Describes one or more of your quotas.

## Example Usage

```hcl
data "outscale_quotas" "s3_by_id" {
  quota_name = ["vm_limit"]
}
```

## Argument Reference

The following arguments are supported:
	 
* `quota_name`-	(Optional) One or more names of quota.	

See detailed information in [Outscale Quota](https://wiki.outscale.net/display/DOCU/Getting+Information+About+Your+Instances).

## Filters

Use the Filter.N parameter to filter the described quotas on the following properties:

* `reference`: -The reference of the quota.
* `quota.display-name`: -	The display name of the quota.
* `quota.group-name`: -	The group name of the quota.

## Attributes Reference

The following attributes are exported:

* `quota_set`	- One or more quotas associated with the user.
* `reference`	- The resource ID if it is a resource-specific quota, global if it is not.
* `request_id` -The ID of the request

See detailed information in [Describe Quota](http://docs.outscale.com/api_fcu/operations/Action_DescribeQuotas_get.html#_api_fcu-action_describequotas_get).