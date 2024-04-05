---
layout: "outscale"
page_title: "OUTSCALE: outscale_vm_types"
sidebar_current: "outscale-vm-types"
description: |-
  [Provides information about VM types.]
---

# outscale_vm_types Data Source

Provides information about VM types.

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/Instance-Types.html).  
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#readvmtypes).

## Example Usage

### All types of VMs
```hcl
data "outscale_vm_types" "all_vm_types" {
}
```

### VMs optimized for Block Storage Unit (BSU)
```hcl
data "outscale_vm_types" "vm_types01" {
    filter {
        name   = "bsu_optimized"
        values = [true]
    }
}
```

### Specific VM type
```hcl
data "outscale_vm_types" "vm_types02" {
    filter {
        name   = "vm_type_names"
        values = ["m3.large"]
    }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) A combination of a filter name and one or more filter values. You can specify this argument for as many filter names as you need. The filter name can be any of the following:
    * `bsu_optimized` - (Optional) This parameter is not available. It is present in our API for the sake of historical compatibility with AWS.
    * `ephemerals_types` - (Optional) The types of ephemeral storage disk.
    * `eths` - (Optional) The number of Ethernet interfaces available.
    * `gpus` - (Optional) The number of GPUs available.
    * `memory_sizes` - (Optional) The amounts of memory, in gibibytes (GiB).
    * `vcore_counts` - (Optional) The numbers of vCores.
    * `vm_type_names` - (Optional) The names of the VM types. For more information, see [VM Types](https://docs.outscale.com/en/userguide/VM-Types.html).
    * `volume_counts` - (Optional) The maximum number of ephemeral storage disks.
    * `volume_sizes` - (Optional) The size of one ephemeral storage disk, in gibibytes (GiB).

## Attribute Reference

The following attributes are exported:

* `vm_types` - Information about one or more VM types.
    * `bsu_optimized` - This parameter is not available. It is present in our API for the sake of historical compatibility with AWS.
    * `ephemerals_type` - The type of ephemeral storage disk.
    * `eth` - The number of Ethernet interface available.
    * `gpu` - The number of GPU available.
    * `max_private_ips` - The maximum number of private IPs per network interface card (NIC).
    * `memory_size` - The amount of memory, in gibibytes.
    * `vcore_count` - The number of vCores.
    * `vm_type_name` - The name of the VM type.
    * `volume_count` - The maximum number of ephemeral storage disks.
    * `volume_size` - The size of one ephemeral storage disk, in gibibytes (GiB).
