---
layout: "outscale"
page_title: "OUTSCALE: outscale_flexible_gpu_catalog"
sidebar_current: "outscale-flexible-gpu-catalog"
description: |-
  [Provides information about the flexible GPU catalog.]
---

# outscale_flexible_gpu_catalog Data Source

Provides information about the flexible GPU catalog.
For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Flexible-GPUs.html).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#readflexiblegpucatalog).

## Example Usage

```hcl
data "outscale_flexible_gpu_catalog" "flexible_gpu_catalog01" {
}
```

## Argument Reference

No argument is supported.

## Attribute Reference

The following attributes are exported:

* `flexible_gpu_catalog` - Information about one or more fGPUs available in the public catalog.
    * `generations` - The generations of VMs that the fGPU is compatible with.
    * `max_cpu` - The maximum number of VM vCores that the fGPU is compatible with.
    * `max_ram` - The maximum amount of VM memory that the fGPU is compatible with.
    * `model_name` - The model of fGPU.
    * `v_ram` - The amount of video RAM (VRAM) of the fGPU.
