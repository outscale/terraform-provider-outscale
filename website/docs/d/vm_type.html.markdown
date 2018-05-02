---
layout: "outscale"
page_title: "OUTSCALE: outscale_vm_type"
sidebar_current: "docs-outscale-datasource-vm-type"
description: |-
  Provides an Outscale resource to get Instances Types.
---

# outscale_vm_types

Provides an Outscale resource to get Instances Types.

## Example Usage

```hcl
data "outscale_vm_type" "outscale_vm_type" {
    filter {
        name = "name"
        values = ["c4.large"]
    }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) One or more filters.

See detailed information in [Outscale VM Types](https://wiki.outscale.net/display/DOCU/Getting+Information+About+Your+Instances).

## Filters

Use the Filter.N parameter to filter the described VM Types on the following properties:

* `name`: - The name of the instance type.
* `vcpu`: - The number of vCPU.
* `memory`: - The amount of memory.
* `storage-size`: - The size of the ephemeral storage.
* `storage-count`: - The number of ephemeral storage.
* `ebs-optimized-available`: -  Whether optimized storage bandwidth is available.



## Attributes Reference

The following attributes are exported:

* `instance_type_set` - Information about zero or more instance types.
* `request_id` - 	The ID of the request

See detailed information in [Describe VM Types](http://docs.outscale.com/api_fcu/operations/Action_DescribeInstanceTypes_get.html#_api_fcu-action_describeinstancetypes_get).
