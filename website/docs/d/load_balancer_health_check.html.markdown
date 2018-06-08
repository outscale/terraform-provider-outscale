---
layout: "outscale"
page_title: "OUTSCALE: outscale_load_balancer_health_check"
sidebar_current: "docs-outscale-datasource-load-balancer-health-check"
description: |-
  Describes one or more load balancers's health.
---

# outscale_load_balancer_health_check

Describes one or more load balancers's health.

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

## Attributes Reference

The following attributes are exported:

* `healthy_threshold` - The number of consecutive successful pings before considering the instance as healthy (between 2 and 10 both included).
* `interval` - The number of seconds between two pings (between 5 and 600 both included).
* `target` - The URL of the checked instance.
* `unhealthy_threshold` - The number of consecutive failed pings before considering the instance as unhealthy (between 2 and 10 both included).

[See detailed description](http://docs.outscale.com/api_lbu/operations/Action_DescribeLoadBalancers_get.html#_api_lbu-action_describeloadbalancers_get)
