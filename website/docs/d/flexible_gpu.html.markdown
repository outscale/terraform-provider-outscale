---
layout: "outscale"
page_title: "3DS OUTSCALE: outscale_flexible_gpu"
sidebar_current: "outscale-flexible-gpu"
description: |-
  [Provides information about a specific flexible GPU.]
---

# outscale_flexible_gpu Data Source

Provides information about a specific flexible GPU.
For more information on this resource, see the [User Guide](https://wiki.outscale.net/display/EN/About+Flexible+GPUs).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-flexiblegpu).

## Example Usage

```
data "outscale_flexible_gpu" "flexible_gpu01" {
  filter {
    name   = "flexible_gpu_ids"
    values = ["fgpu-12345678"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - One or more filters.
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
