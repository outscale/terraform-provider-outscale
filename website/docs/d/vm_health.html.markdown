---
layout: "outscale"
page_title: "OUTSCALE: outscale_vm_health"
sidebar_current: "docs-outscale-datasource-vm-health"
description: |-
  Describes the state of one or more back-end instances registered with a specified load balancer.
---

# outscale_vm_health

Describes the state of one or more back-end instances registered with a specified load balancer.

## Example Usage

```hcl
resource "outscale_load_balancer" "bar" {
    availability_zones = ["eu-west-2a"]
    load_balancer_name = "foobar-terraform-elb"
    listeners {
        instance_port = 8000
        instance_protocol = "HTTP"
        load_balancer_port = 80
        protocol = "HTTP"
    }

    tag {
        bar = "baz"
    }
}

resource "outscale_vm" "foo1" {
    image_id = "ami-880caa66"
    instance_type = "t2.micro"
}

resource "outscale_load_balancer_vms" "foo1" {
    load_balancer_name      = "${outscale_load_balancer.bar.id}"
    instances = [{
        instance_id = "${outscale_vm.foo1.id}"
    }]
}

data "outscale_vm_health" "web" {
    load_balancer_name = "${outscale_load_balancer.bar.id}"
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_name` - The name of the load balancer.
* `instances.N` - One or more instance IDs, each with following attributes:
  - `instance_id` - If true, the access logs are enabled for your load balancer. If false, they are not.

## Attributes Reference

The following attributes are exported:

* `instance_states.N` - (Optional) Information about the health of one or more back-end instances, each containing the following attributes:
  - `description` - (Optional) A description of the instance state.
  - `instance_id` - (Optional) The ID of the instance.
  - `reason_code` - (Optional) Information about the cause of OutOfService instances.\
  Specifically, whether the cause is Elastic Load Balancing or the instance (ELB | Instance | N/A).
  - `state` - (Optional) The current state of the instance (InService | OutOfService | Unknown).

[See detailed description](http://docs.outscale.com/api_lbu/operations/Action_DescribeInstanceHealth_get.html#_api_lbu-action_describeinstancehealth_get)
