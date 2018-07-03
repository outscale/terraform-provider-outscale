---
layout: "outscale"
page_title: "OUTSCALE: outscale_load_balancer_tags"
sidebar_current: "docs-outscale-datasource-load-balancer-tags"
description: |-
  Describes the tags associated with one or more specified load balancers.
---

# outscale_load_balancer_tags

Describes the tags associated with one or more specified load balancers.

## Example Usage

```hcl
resource "outscale_load_balancer" "bar" {
    availability_zones = ["eu-west-2a"]
    load_balancer_name = "foobar-terraform-elb-aaaaa"
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

resource "outscale_load_balancer_tags" "tags" {
    load_balancer_names = ["${outscale_load_balancer.bar.id}"]
    tags = [{
        key = "bar2"
        value = "baz2"
    }]
}

data "outscale_load_balancer_tags" "testds" {
    load_balancer_names = ["${outscale_load_balancer.bar.id}"]
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_name.N` - The name of one or more load balancers.

## Attributes Reference

The following attributes are exported:

* `tag_descriptions.N` - Information about the tags associated with the specified load balancers.
    - `load_balancer_name` - (Optional) The name of the load balancer.
    - `tags.N` - One or more tags associated with the load balancer.

[See detailed description](http://docs.outscale.com/api_lbu/operations/Action_DescribeTags_get.html#_api_lbu-action_describetags_get)
