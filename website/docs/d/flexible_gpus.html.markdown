---
layout: "outscale"
page_title: "OUTSCALE: outscale_flexible_gpus"
sidebar_current: "outscale-flexible-gpus"
description: |-
  [Provides information about flexible GPUs.]
---

# outscale_flexible_gpus Data Source

Provides information about flexible GPUs.
For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Flexible-GPUs.html).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-flexiblegpu).

## Example Usage

```hcl
data "outscale_flexible_gpus" "flexible_gpus01" {
  filter {
    name   = "fgpu_ids"
    values = ["fgpu-12345678", "fgpu-12345679"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) A combination of a filter name and one or more filter values. You can specify this argument for as many filter names as you need. The filter name can be any of the following:
    * `delete_on_vm_deletion` - (Optional) Indicates whether the fGPU is deleted when terminating the VM.
    * `flexible_gpu_ids` - (Optional) One or more IDs of fGPUs.
    * `generations` - (Optional) The processor generations that the fGPUs are compatible with.
    * `model_names` - (Optional) One or more models of fGPUs. For more information, see [About Flexible GPUs](https://wiki.outscale.net/display/EN/About+Flexible+GPUs).
    * `states` - (Optional) The states of the fGPUs (`allocated` \| `attaching` \| `attached` \| `detaching`).
    * `subregion_names` - (Optional) The Subregions where the fGPUs are located.
    * `vm_ids` - (Optional) One or more IDs of VMs.

## Attribute Reference

The following attributes are exported:

* `flexible_gpus` - Information about one or more fGPUs.
    * `delete_on_vm_deletion` - If true, the fGPU is deleted when the VM is terminated.
    * `flexible_gpu_id` - The ID of the fGPU.
    * `generation` - The compatible processor generation.
    * `model_name` - The model of fGPU. For more information, see [About Flexible GPUs](https://wiki.outscale.net/display/EN/About+Flexible+GPUs).
    * `state` - The state of the fGPU (`allocated` \| `attaching` \| `attached` \| `detaching`).
    * `subregion_name` - The Subregion where the fGPU is located.
    * `vm_id` - The ID of the VM the fGPU is attached to, if any.
