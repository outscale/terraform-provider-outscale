---
layout: "outscale"
page_title: "OUTSCALE: outscale_load_balancer_access_logs"
sidebar_current: "docs-outscale-datasource-load-balancer-access-logs"
description: |-
  Describes the attributes of the specified load balancer.
---

# outscale_load_balancer_access_logs

Describes the attributes of the specified load balancer.

## Example Usage

```hcl
resource "outscale_load_balancer" "bar" {
    availability_zones = ["eu-west-2a"]
    load_balancer_name               = "foobar-terraform-elb-ds"
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

resource "outscale_load_balancer_attributes" "bar2" {
    enabled = "false"
    s3_bucket_name = "donustestbucket"
    load_balancer_name = "${outscale_load_balancer.bar.id}"
}

data "outscale_load_balancer_access_logs" "test" {
    load_balancer_name = "${outscale_load_balancer.bar.id}"
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_name` - The name of the load balancer.

## Attributes Reference

The following attributes are exported:

* `emit_interval` - (Optional) The time interval for access logs publication into the OSU bucket, that can be either 5 or 60 minutes (by default, 60 minutes).
* `enabled` - If true, the access logs are enabled for your load balancer. If false, they are not.
* `s3_bucket_name` - (Optional) The name of the Object Storage Unit (OSU) bucket for the access logs.
* `s3_bucket_prefix` - (Optional) The path to the folder in your OSU bucket for your access logs information. If not specified, the access logs are published at the root level of your bucket.

[See detailed description](http://docs.outscale.com/api_lbu/operations/Action_DescribeLoadBalancerAttributes_get.html#_api_lbu-action_describeloadbalancerattributes_get)
