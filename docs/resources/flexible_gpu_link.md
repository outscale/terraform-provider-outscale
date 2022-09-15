---
layout: "outscale"
page_title: "OUTSCALE: outscale_flexible_gpu_link"
sidebar_current: "outscale-flexible-gpu-link"
description: |-
  [Manages a flexible GPU link.]
---

# outscale_flexible_gpu_link Resource

Manages a flexible GPU link.
For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Flexible-GPUs.html).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-flexiblegpu).

## Example Usage

### Required resources

```hcl
resource "outscale_vm" "vm01" {
  image_id                 = ami-12345678
  vm_type                  = t2.small
  keypair_name             = var.keypair_name
  placement_subregion_name = "eu-west-2a"
}
resource "outscale_flexible_gpu" "flexible_gpu01" {
  model_name            = var.model_name
  generation            = "v4"
  subregion_name        = "eu-west-2a"
  delete_on_vm_deletion = true
}
```

### Create a flexible GPU link

```hcl
resource "outscale_flexible_gpu_link" "link_fgpu01" {
  flexible_gpu_id = outscale_flexible_gpu.flexible_gpu01.flexible_gpu_id
  vm_id           = outscale_vm.vm01.vm_id
}
```

## Argument Reference

The following arguments are supported:

* `flexible_gpu_id` - (Required) The ID of the fGPU you want to attach.
* `vm_id` - (Required) The ID of the VM you want to attach the fGPU to.

## Attribute Reference

No attribute is exported.

## Import

A flexible GPU link can be imported using the flexible GPU ID. For example:

```console

$ terraform import outscale_flexible_gpu_link.imported_link_fgpu fgpu-12345678

```