---
layout: "outscale"
page_title: "OUTSCALE: outscale_flexible_gpu"
subcategory: "OUTSCALE API"
sidebar_current: "outscale-flexible-gpu"
description: |-
  [Manages a flexible GPU.]
---

# outscale_flexible_gpu Resource

Manages a flexible GPU.

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Flexible-GPUs.html).  
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-flexiblegpu).

## Example Usage

### Create a flexible GPU

```hcl
resource "outscale_flexible_gpu" "flexible_gpu01" {
    model_name             =  var.model_name
    generation             =  "v4"
    subregion_name         =  "${var.region}a"
    delete_on_vm_deletion  =  true
}
```

## Argument Reference

The following arguments are supported:

* `delete_on_vm_deletion` - (Optional) If true, the fGPU is deleted when the VM is terminated.
* `generation` - (Optional) The processor generation that the fGPU must be compatible with. If not specified, the oldest possible processor generation is selected (as provided by [ReadFlexibleGpuCatalog](https://docs.outscale.com/api#readflexiblegpucatalog) for the specified model of fGPU).
* `model_name` - (Required) The model of fGPU you want to allocate. For more information, see [About Flexible GPUs](https://docs.outscale.com/en/userguide/About-Flexible-GPUs.html).
* `subregion_name` - (Required) The Subregion in which you want to create the fGPU.

## Attribute Reference

The following attributes are exported:

* `delete_on_vm_deletion` - If true, the fGPU is deleted when the VM is terminated.
* `flexible_gpu_id` - The ID of the fGPU.
* `generation` - The compatible processor generation.
* `model_name` - The model of fGPU. For more information, see [About Flexible GPUs](https://docs.outscale.com/en/userguide/About-Flexible-GPUs.html).
* `state` - The state of the fGPU (`allocated` \| `attaching` \| `attached` \| `detaching`).
* `subregion_name` - The Subregion where the fGPU is located.
* `vm_id` - The ID of the VM the fGPU is attached to, if any.

## Timeouts

The `timeouts` block enables you to configure [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - Defaults to 10 minutes.
* `read` - Defaults to 5 minutes.
* `update` - Defaults to 10 minutes.
* `delete` - Defaults to 5 minutes.

## Import

A flexible GPU can be imported using its ID. For example:

```console

$ terraform import outscale_flexible_gpu.imported_fgpu fgpu-12345678

```