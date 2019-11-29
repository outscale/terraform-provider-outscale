---
layout: "outscale"
page_title: "3DS OUTSCALE: outscale_vm_type"
sidebar_current: "docs-outscale-datasource-vm-type"
description: |-
  [Provides information about a specific VM type.]
---

# outscale_vm_type Data Source

Provides information about a specific VM type.
For more information on this resource, see the [User Guide](https://wiki.outscale.net/display/EN/Instance+Types).
For more information on this resource actions, see the [API documentation](https://docs-beta.outscale.com/#readvmtypes).

## Example Usage

```hcl
[exemple de code]
```

## Argument Reference

The following arguments are supported:

* `filters` - One or more filters.
  * `bsu_optimized` - (Optional) Indicates whether the VM is optimized for BSU I/O.
  * `memory_sizes` - (Optional) The amounts of memory, in bytes.
  * `vcore_counts` - (Optional) The numbers of vCores.
  * `vm_type_names` - (Optional) The names of the VM types. For more information, see [Instance Types](https://wiki.outscale.net/display/EN/Instance+Types).
  * `volume_counts` - (Optional) The maximum number of ephemeral storage disks.
  * `volume_sizes` - (Optional) The size of one ephemeral storage disk, in gibibytes (GiB).

## Attribute Reference

The following attributes are exported:

* `vm_types` - Information about one or more VM types.
  * `bsu_optimized` - Indicates whether the VM is optimized for BSU I/O.
  * `max_private_ips` - The maximum number of private IP addresses per network interface card (NIC).
  * `memory_size` - The amount of memory, in gibibytes.
  * `vcore_count` - The number of vCores.
  * `vm_type_name` - The name of the VM type.
  * `volume_count` - The maximum number of ephemeral storage disks.
  * `volume_size` - The size of one ephemeral storage disk, in gibibytes (GiB).
