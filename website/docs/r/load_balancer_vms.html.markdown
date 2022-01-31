---
layout: "outscale"
page_title: "OUTSCALE: outscale_load_balancer_vms"
sidebar_current: "outscale-load-balancer-vms"
description: |-
  [Manages load balancer VMs.]
---

# outscale_load_balancer_vms Resource

Manages load balancer VMs.
For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Load-Balancers.html).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-loadbalancer).

## Example Usage

### Required resources

```hcl
resource "outscale_vm" "outscale_vm01" {
    image_id     = "ami-12345678"
    vm_type      = "t2.small"
    keypair_name = var.keypair_name
}

resource "outscale_vm" "outscale_vm_02" {
    image_id     = var.image_id
    vm_type      = var.vm_type
    keypair_name = var.keypair_name
}

resource "outscale_load_balancer" "load_balancer01" {
    load_balancer_name = "load-balancer-for-backend-vms"
    subregion_names    = ["${var.region}a"]
    listeners {
       backend_port           = 80
       backend_protocol       = "TCP"
       load_balancer_protocol = "TCP"
       load_balancer_port     = 80
    }
    tags {
       key   = "name"
       value = "outscale_load_balancer01"
    }
}
```

### Register VMs with a load balancer

```hcl
resource "outscale_load_balancer_vms" "outscale_load_balancer_vms01" {
    load_balancer_name = load-balancer-for-backend-vms
    backend_vm_ids     = [outscale_vm.outscale_vm01.vm_id,outscale_vm.outscale_vm_02.vm_id]
}
```

## Argument Reference

The following arguments are supported:

* `backend_vm_ids` - (Required) One or more IDs of back-end VMs.<br />
Specifying the same ID several times has no effect as each back-end VM has equal weight.
* `load_balancer_name` - (Required) The name of the load balancer.

## Attribute Reference

No attribute is exported.

