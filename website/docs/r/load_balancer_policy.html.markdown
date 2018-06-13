---
layout: "outscale"
page_title: "OUTSCALE: outscale_load_balancer_policy"
sidebar_current: "docs-outscale-datasource-load-balancer-policy"
description: |-
  Replaces the current set of policies for a load balancer with another specified one.
If the PolicyNames.member.N parameter is empty, all current policies are disabled.
---

# outscale_load_balancer_policy

Replaces the current set of policies for a load balancer with another specified one.
If the PolicyNames.member.N parameter is empty, all current policies are disabled.

## Example Usage

```hcl
resource "outscale_load_balancer" "lb" {
    availability_zones = ["eu-west-2a"]
    load_balancer_name = "foobar-terraform-lbu-1"
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

resource "outscale_load_balancer_policy" "outscale_load_balancer_policy" {
    load_balancer_name = "${outscale_load_balancer.outscale_load_balancer.load_balancer_name}"

    load_balancer_port = "${outscale_load_balancer.outscale_load_balancer.listeners.0.load_balancer_port}"
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_name` The name of the load balancer.
* `policy_names.N` - The list of policies names (must contain all the policies to be enabled).
* `load_balancer_port` - The external port of the load balancer.

## Attributes Reference

The following attributes are exported:

* `load_balancer_name` The name of the load balancer.
* `policy_names.N` - The list of policies names (must contain all the policies to be enabled).
* `load_balancer_port` - The external port of the load balancer.
* `request_id` - The ID of the request.

[See detailed description](http://docs.outscale.com/api_lbu/operations/Action_SetLoadBalancerPoliciesOfListener_get.html#_api_lbu-action_setloadbalancerpoliciesoflistener_get)
