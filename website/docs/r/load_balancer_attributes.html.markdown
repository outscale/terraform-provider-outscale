---
layout: "outscale"
page_title: "OUTSCALE: outscale_load_balancer_attributes"
sidebar_current: "docs-outscale-datasource-load-balancer-cookiepolicy"
description: |-
  Modifies the specified attributes of a load balancer.
---

# outscale_load_balancer_attributes

Modifies the specified attributes of a load balancer.
You can modify the load balancer AccessLogs attribute, that you can either enable or disable.

## Example Usage

```hcl
resource "outscale_load_balancer" "bar" {
    availability_zones = ["eu-west-2a"]
    load_balancer_name               = "foobar-terraform-elb-%d"
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
```

## Argument Reference

The following arguments are supported:

* `load_balancer_name` The name of the load balancer.
* `emit_interval` - (Optional) The time interval for access logs publication into the OSU bucket, that can be either 5 or 60 minutes (by default, 60 minutes).
* `enabled` - If true, the access logs are enabled for your load balancer. If false, they are not.
* `s3_bucket_name` - (Optional) The name of the Object Storage Unit (OSU) bucket for the access logs.
* `s3_bucket_prefix` - (Optional) The path to the folder in your OSU bucket for your access logs information. If not specified, the access logs are published at the root level of your bucket.

## Attributes Reference

The following attributes are exported:

* `load_balancer_name` The name of the load balancer.
* `load_balancer_attributes` - IInformation about the load balancer attributes.:
  - `access_log` - If enabled, information about requests that the load balancer are written into the specified Object Storage Unit (OSU) bucket.
    - `emit_interval` - (Optional) The time interval for access logs publication into the OSU bucket, that can be either 5 or 60 minutes (by default, 60 minutes).
    - `enabled` - If true, the access logs are enabled for your load balancer. If false, they are not.
    - `s3_bucket_name` - (Optional) The name of the Object Storage Unit (OSU) bucket for the access logs.
    - `s3_bucket_prefix` - (Optional) The path to the folder in your OSU bucket for your access logs information. If not specified, the access logs are published at the root level of your bucket.

[See detailed description](http://docs.outscale.com/api_lbu/operations/Action_ModifyLoadBalancerAttributes_get.html#_api_lbu-action_modifyloadbalancerattributes_get)
