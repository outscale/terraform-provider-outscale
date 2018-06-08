---
layout: "outscale"
page_title: "OUTSCALE: outscale_load_balancer_tags"
sidebar_current: "docs-outscale-datasource-load-balancer-tags"
description: |-
  Adds one or more tags to a specified load balancer.
---

# outscale_load_balancer_tags

Adds one or more tags to a specified load balancer.
You can add up to 10 tags per load balancer.

## Example Usage

```hcl
resource "outscale_load_balancer" "bar" {
    availability_zones = ["eu-west-2a"]
    load_balancer_name = "foobar-terraform-elb-1"
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
```

## Argument Reference

The following arguments are supported:

* `load_balancer_names.N` The name of the load balancer.
* `tags.N` - One or more tags, each with the following attributes:
  - `key` - The key of the tag.
  - `value` - (Optional) The value of the tag.

## Attributes Reference

No attributes are exported.

[See detailed description](http://docs.outscale.com/api_lbu/operations/Action_AddTags_get.html#_api_lbu-action_addtags_get)
