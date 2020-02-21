---
layout: "outscale"
page_title: "3DS OUTSCALE: outscale_vm_state"
sidebar_current: "outscale-vm-state"
description: |-
  [Provides information about VM states.]
---

# outscale_vm_state Data Source

Provides information about VM states.
For more information on this resource, see the [User Guide](https://wiki.outscale.net/display/EN/About+Instance+Lifecycle).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#readvmsstate).

## Example Usage

```hcl

data "outscale_vm_states" "vm_states01" {
  filter {
    name   = "vm_ids"
    values = ["i-12345678", "i-12345679"]
  }
}


```

## Argument Reference

The following arguments are supported:

* `all_vms` - (Optional) If `true`, includes the status of all VMs. By default or if set to `false`, only includes the status of running VMs.
* `filter` - One or more filters.
  * `subregion_names` - (Optional) The names of the Subregions of the VMs.
  * `vm_ids` - (Optional) One or more IDs of VMs.
  * `vm_states` - (Optional) The states of the VMs (`pending` \| `running` \| `shutting-down` \| `terminated` \| `stopping` \| `stopped`).

## Attribute Reference

The following attributes are exported:

* `vm_states` - Information about one or more VM states.
  * `maintenance_events` - One or more scheduled events associated with the VM.
    * `code` - The code of the event (`system-reboot` \| `system-maintenance`).
    * `description` - The description of the event.
    * `not_after` - The latest scheduled end time for the event.
    * `not_before` - The earliest scheduled start time for the event.
  * `subregion_name` - The name of the Subregion of the VM.
  * `vm_id` - The ID of the VM.
  * `vm_state` - The state of the VM (`pending` \| `running` \| `shutting-down` \| `terminated` \| `stopping` \| `stopped`).
