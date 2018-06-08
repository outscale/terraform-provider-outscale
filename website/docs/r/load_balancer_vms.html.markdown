---
layout: "outscale"
page_title: "OUTSCALE: outscale_load_balancer_vms"
sidebar_current: "docs-outscale-datasource-load-balancer-vms"
description: |-
  Registers one or more instances with a specified load balancer.
---

# outscale_load_balancer_vms

Registers one or more instances with a specified load balancer.
The instances must be running in the same network as the load balancer (in the public Cloud or in the same VPC). It may take a little time for an instance to be registered with the load balancer. Once the instance is registered with a load balancer, it receives traffic and requests from this load balancer and is called a back-end instance.

## Example Usage

```hcl
resource "outscale_load_balancer" "bar" {
    load_balancer_name = "load-test"

    availability_zones = ["eu-west-2a"]
        listeners {
        instance_port = 8000
        instance_protocol = "HTTP"
        load_balancer_port = 80
        protocol = "HTTP"
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
```

## Argument Reference

The following arguments are supported:

* `load_balancer_name` The name of the load balancer.
* `instances.N` - One or more instance IDs, each with following attributes:
  - `instance_id` - If true, the access logs are enabled for your load balancer. If false, they are not.

## Attributes Reference

No attributes are exported.

[See detailed description](http://docs.outscale.com/api_lbu/operations/Action_RegisterInstancesWithLoadBalancer_get.html#_api_lbu-action_registerinstanceswithloadbalancer_get)
