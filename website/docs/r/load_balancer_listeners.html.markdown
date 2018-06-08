---
layout: "outscale"
page_title: "OUTSCALE: outscale_load_balancer_listeners"
sidebar_current: "docs-outscale-datasource-load-balancer-listeners"
description: |-
  Creates one or more listeners for a specified load balancer.
---

# outscale_load_balancer_listeners

Creates one or more listeners for a specified load balancer.

## Example Usage

```hcl
resource "outscale_load_balancer" "lb" {
    availability_zones = ["eu-west-2a"]
    load_balancer_name               = "foobar-terraform-lbu-1"
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

resource "outscale_load_balancer_listeners" "bar" {
    load_balancer_name               = "${outscale_load_balancer.lb.id}"
    listeners {
        instance_port = 9000
        instance_protocol = "HTTP"
        load_balancer_port = 9000
        protocol = "HTTP"
    }
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_name` The name of the load balancer.
* `listeners.N` - One or more listeners, each with the following attributes:
  - `instance_port` - The port on which the instance is listening (between 1 and 65535 both included).
  - `instance_protocol` - (Optional) The protocol for routing traffic to back-end instances (HTTP | TCP).
  - `load_balancer_port` - The port on which the load balancer is listening (25, 80, 443, 465, 587, or between 1024 and 65535 both included).
  - `protocol` - TThe routing protocol (HTTP | TCP).
  - `ssl_certificate_id` - (Optional) The ID of the server certificate.

## Attributes Reference

The following attributes are exported:

`load_balancer_name` The name of the load balancer.
* `listeners.N` - One or more listeners, each with the following attributes:
  - `instance_port` - The port on which the instance is listening (between 1 and 65535 both included).
  - `instance_protocol` - (Optional) The protocol for routing traffic to back-end instances (HTTP | TCP).
  - `load_balancer_port` - The port on which the load balancer is listening (25, 80, 443, 465, 587, or between 1024 and 65535 both included).
  - `protocol` - TThe routing protocol (HTTP | TCP).
  - `ssl_certificate_id` - (Optional) The ID of the server certificate.

[See detailed description](http://docs.outscale.com/api_lbu/operations/Action_CreateLoadBalancerListeners_get.html#_api_lbu-action_createloadbalancerlisteners_get)
