---
layout: "outscale"
page_title: "OUTSCALE: outscale_load_balancer_listener_descriptions"
sidebar_current: "docs-outscale-datasource-load-balancer-listener-descriptions"
description: |-
  Describes the listener of the specified load balancer.
---

# outscale_load_balancer_listener_descriptions

Describes the listener of the specified load balancer.

## Example Usage

```hcl
resource "outscale_load_balancer" "bar" {
    availability_zones = ["eu-west-2a"]
    load_balancer_name               = "foobar-terraform-elb"
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

data "outscale_load_balancer_listener_descriptions" "test" {
    load_balancer_names = ["${outscale_load_balancer.bar.id}"]
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_names.N` - The names of the load balancers.

## Attributes Reference

The following attributes are exported:

* `listener_descriptions` - The listeners for the load balancer each with the following attributes:
  - `listener` - the listener with the following attributes:
    - `instance_port` - The port on which the instance is listening (between 1 and 65535 both included).
    - `instance_protocol` - (Optional) The protocol for routing traffic to back-end instances (HTTP | TCP).
    - `load_balancer_port` - The port on which the load balancer is listening (25, 80, 443, 465, 587, or between 1024 and 65535 both included).
    - `protocol` - TThe routing protocol (HTTP | TCP).
    - `ssl_certificate_id` - (Optional) The ID of the server certificate.
  - `policy_names` - the policy names.

[See detailed description](http://docs.outscale.com/api_lbu/operations/Action_DescribeLoadBalancers_get.html#_api_lbu-action_describeloadbalancers_get)
