---
layout: "outscale"
page_title: "OUTSCALE: outscale_load_balancer_cookiepolicy"
sidebar_current: "docs-outscale-datasource-load-balancer-cookiepolicy"
description: |-
  Creates a stickiness policy with sticky session lifetimes following the one of an application-generated cookie.
---

# outscale_load_balancer_cookiepolicy

Creates a stickiness policy with sticky session lifetimes following the one of an application-generated cookie.
The created policy can be used only with HTTP listeners.

## Example Usage

```hcl
resource "outscale_load_balancer" "lb" {
    load_balancer_name = "tf-test-lb-abc"
    availability_zones = ["eu-west-2a"]
    listeners {
        instance_port = 8000
        instance_protocol = "HTTP"
        load_balancer_port = 80
        protocol = "HTTP"
    }
}

resource "outscale_load_balancer_cookiepolicy" "foo" {
    policy_name = "foo-policy"
    load_balancer_name = "${outscale_load_balancer.lb.id}"
    cookie_name = "MyAppCookie"
}
```

## Argument Reference

The following arguments are supported:

* `cookie_name` The name of the application cookie used for stickiness.
* `policy_name` The unique name of the policy (alphanumeric characters and dashes (-)).
* `load_balancer_name` The name of the load balancer.

## Attributes Reference

The following attributes are exported:

* `cookie_name` The name of the application cookie used for stickiness.
* `policy_name` The unique name of the policy (alphanumeric characters and dashes (-)).
* `load_balancer_name` The name of the load balancer.

[See detailed description](http://docs.outscale.com/api_lbu/operations/Action_CreateAppCookieStickinessPolicy_get.html#_api_lbu-action_createappcookiestickinesspolicy_get)
