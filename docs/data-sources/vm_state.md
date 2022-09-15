---
layout: "outscale"
page_title: "OUTSCALE: outscale_vm_state"
sidebar_current: "outscale-vm-state"
description: |-
  [Provides information about a specific VM state.]
---

# outscale_vm_state Data Source

Provides information about a specific VM state.
For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Instance-Lifecycle.html).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#readvmsstate).

## Example Usage

```hcl
data "outscale_vm_state" "vm_state01" {
  filter {
    name   = "vm_ids"
    values = ["i-12345678"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `all_vms` - (Optional) If true, includes the status of all VMs. By default or if set to false, only includes the status of running VMs.
* `filter` - (Optional) A combination of a filter name and one or more filter values. You can specify this argument for as many filter names as you need. The filter name can be any of the following:
    * `maintenance_event_codes` - (Optional) The code for the scheduled event (`system-reboot` \| `system-maintenance`).
    * `maintenance_event_descriptions` - (Optional) The description of the scheduled event.
    * `maintenance_events_not_after` - (Optional) The latest time the event can end.
    * `maintenance_events_not_before` - (Optional) The earliest time the event can start.
    * `subregion_names` - (Optional) The names of the Subregions of the VMs.
    * `vm_ids` - (Optional) One or more IDs of VMs.
    * `vm_states` - (Optional) The states of the VMs (`pending` \| `running` \| `stopping` \| `stopped` \| `shutting-down` \| `terminated` \| `quarantine`).

## Attribute Reference

The following attributes are exported:

* `maintenance_events` - One or more scheduled events associated with the VM.
    * `code` - The code of the event (`system-reboot` \| `system-maintenance`).
    * `description` - The description of the event.
    * `not_after` - The latest scheduled end time for the event.
    * `not_before` - The earliest scheduled start time for the event.
* `subregion_name` - The name of the Subregion of the VM.
* `vm_id` - The ID of the VM.
* `vm_state` - The state of the VM (`pending` \| `running` \| `stopping` \| `stopped` \| `shutting-down` \| `terminated` \| `quarantine`).
