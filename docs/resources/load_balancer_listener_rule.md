---
layout: "outscale"
page_title: "OUTSCALE: outscale_load_balancer_listener_rule"
sidebar_current: "outscale-load-balancer-listener-rule"
description: |-
  [Manages a load balancer listener rule.]
---

# outscale_load_balancer_listener_rule Resource

Manages a load balancer listener rule.
For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Load-Balancers.html).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-listener).

## Example Usage

### Required resources

```hcl
resource "outscale_vm" "vm01" {
  image_id     = var.image_id
  vm_type      = var.vm_type
  keypair_name = var.keypair_name
}

resource "outscale_load_balancer" "load_balancer01" {
  load_balancer_name = "terraform-public-load-balancer"
  subregion_names    = ["${var.region}a"]
  listeners {
    backend_port           = 80
    backend_protocol       = "TCP"
    load_balancer_protocol = "TCP"
    load_balancer_port     = 80
  }
  tags {
    key   = "name"
    value = "terraform-public-load-balancer"
  }
}

resource "outscale_load_balancer_vms" "outscale_load_balancer_vms01" {
  load_balancer_name = outscale_load_balancer.load_balancer01.id
  backend_vm_ids     = [outscale_vm.vm01.vm_id]
}
```

### Create a listener rule based on path pattern

```hcl
resource "outscale_load_balancer_listener_rule" "rule01" {
  listener {
    load_balancer_name = outscale_load_balancer.load_balancer01.id
    load_balancer_port = 80
  }
  listener_rule {
    action             = "forward"
    listener_rule_name = "terraform-listener-rule01"
    path_pattern       = "*.abc.*.abc.*.com"
    priority           = 10
  }
  vm_ids = [outscale_vm.vm01.vm_id]
}
```

### Create a listener rule based on host pattern

```hcl
resource "outscale_load_balancer_listener_rule" "rule02" {
  listener  {
    load_balancer_name = outscale_load_balancer.load_balancer01.id
    load_balancer_port = 80
  }
  listener_rule {
    action             = "forward"
    listener_rule_name = "terraform-listener-rule02"
    host_name_pattern  = "*.abc.-.abc.*.com"
    priority           = 1
  }
  vm_ids = [outscale_vm.vm01.vm_id]
}
```

## Argument Reference

The following arguments are supported:

* `listener_rule` - Information about the listener rule.
    * `action` - (Optional) The type of action for the rule (always `forward`).
    * `host_name_pattern` - (Optional) A host-name pattern for the rule, with a maximum length of 128 characters. This host-name pattern supports maximum three wildcards, and must not contain any special characters except [-.?]. 
    * `listener_rule_name` - (Optional) A human-readable name for the listener rule.
    * `path_pattern` - (Optional) A path pattern for the rule, with a maximum length of 128 characters. This path pattern supports maximum three wildcards, and must not contain any special characters except [_-.$/~&quot;'@:+?].
    * `priority` - (Optional) The priority level of the listener rule, between `1` and `19999` both included. Each rule must have a unique priority level. Otherwise, an error is returned.
* `listener` - Information about the load balancer.
    * `load_balancer_name` - (Optional) The name of the load balancer to which the listener is attached.
    * `load_balancer_port` - (Optional) The port of load balancer on which the load balancer is listening (between `1` and `65535` both included).
* `vm_ids` - (Required) The IDs of the backend VMs.

## Attribute Reference

The following attributes are exported:

* `action` - The type of action for the rule (always `forward`).
* `host_name_pattern` - A host-name pattern for the rule, with a maximum length of 128 characters. This host-name pattern supports maximum three wildcards, and must not contain any special characters except [-.?].
* `listener_id` - The ID of the listener.
* `listener_rule_id` - The ID of the listener rule.
* `listener_rule_name` - A human-readable name for the listener rule.
* `path_pattern` - A path pattern for the rule, with a maximum length of 128 characters. This path pattern supports maximum three wildcards, and must not contain any special characters except [_-.$/~&quot;'@:+?].
* `priority` - The priority level of the listener rule, between `1` and `19999` both included. Each rule must have a unique priority level. Otherwise, an error is returned.
* `vm_ids` - The IDs of the backend VMs.

