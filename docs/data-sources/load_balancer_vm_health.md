---
layout: "outscale"
page_title: "OUTSCALE: outscale_load_balancer_vm_health"
sidebar_current: "outscale-load-balancer-vm-health"
description: |-
  [Provides information about the health of one or more back-end VMs registered with a specific load balancer.]
---

# outscale_load_balancer_vm_health Data Source

Provides information about the health of one or more back-end VMs registered with a specific load balancer.
For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Load-Balancers.html).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#readvmshealth).

## Example Usage

```hcl
data "outscale_load_balancer_vm_health" "load_balancer_vm_health01" {
    load_balancer_name = "load_balancer01"
    backend_vm_ids     = ["i-12345678","i-87654321"]
}
```

## Argument Reference

The following arguments are supported:

* `backend_vm_ids` - (Optional) One or more IDs of back-end VMs.
* `load_balancer_name` - (Required) The name of the load balancer.

## Attribute Reference

The following attributes are exported:

* `backend_vm_health` - Information about the health of one or more back-end VMs.
    * `description` - The description of the state of the back-end VM.
    * `state` - The state of the back-end VM (`InService` \| `OutOfService` \| `Unknown`).
    * `state_reason` - Information about the cause of `OutOfService` VMs.<br />
Specifically, whether the cause is Elastic Load Balancing or the VM (`ELB` \| `Instance` \| `N/A`).
    * `vm_id` - The ID of the back-end VM.
