---
layout: "outscale"
page_title: "OUTSCALE: outscale_load_balancers"
sidebar_current: "docs-outscale-datasource-load-balancers"
description: |-
  Describes one or more load balancers
---

# outscale_load_balancers

Describes one or more load balancers

## Example Usage

```hcl
resource "outscale_load_balancer" "bar" {
    availability_zones_member = ["eu-west-2a"]
    load_balancer_name               = "foobar-terraform-elb"
    listeners_member {
        instance_port = 8000
        instance_protocol = "HTTP"
        load_balancer_port = 80
        // Protocol should be case insensitive
        protocol = "HTTP"
    }

    tag {
        bar = "baz"
    }
}

data "outscale_load_balancers" "test" {
    load_balancer_name = ["${outscale_load_balancer.bar.id}"]
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_name.N` (Optional). The name of one or more load balancers.

## Attributes Reference

The following attributes are exported:

* `load_balancer_descriptions_member.N` - Information about one or more load balancers, each containing the following data:
  - `security_groups_member.N` - The security groups for the load balancer. Valid only for load balancers in a VPC.
  - `subnets_member.N` - The IDs of the subnets for the load balancer.
  - `listener_descriptions_member.N` - The listeners for the load balancer. 
  - `policies` - The policies defined for the load balancer.
  - `health_check` - Information about the health checks conducted on the load balancer.
  - `instances_member.N` - The IDs of the instances for the load balancer.
  - `availability_zones_member.N` - The Availability Zones for the load balancer.
  - `scheme` - The type of load balancer. Valid only for load balancers in a VPC.
  If Scheme is internet-facing, the load balancer has a public DNS name that resolves to a public IP address.
  If Scheme is internal, the load balancer has a public DNS name that resolves to a private IP address.
  - `source_security_group` - The security group for the load balancer, which you can use as part of your inbound rules for your registered instances.
  To only allow traffic from load balancers, add a security group rule that specifies this source security group as the inbound source.
  - `vpc_id` - The ID of the VPC for the load balancer.
  - `dns_name` - The DNS name of the load balancer.

[See detailed description](http://docs.outscale.com/api_lbu/operations/Action_DescribeLoadBalancers_get.html#_api_lbu-action_describeloadbalancers_get)
