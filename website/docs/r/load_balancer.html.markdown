---
layout: "outscale"
page_title: "OUTSCALE: outscale_load_balancer"
sidebar_current: "docs-outscale-resource-load-balancer"
description: |-
  Creates a load balancer.
---

# outscale_load_balancer

Creates a load balancer.
The load balancer is created with a unique Domain Name Service (DNS) name. It receives the incoming traffic and routes it to its registered instances.
By default, this action creates an Internet-facing load balancer, resolving to public IP addresses. To create an internal load balancer in a Virtual Private Cloud (VPC), resolving to private IP addresses, use the Scheme parameter.

## Example Usage

```hcl
resource "outscale_load_balancer" "bar" {
    availability_zones = ["eu-west-2a"]
    load_balancer_name               = "foobar-terraform-elb-1"
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
```

## Argument Reference

The following arguments are supported:

* `availability_zones.N` - The name of the Availability Zone (currently, only one AZ is supported). This parameter is not required if you create a load balancer in a VPC. To create an internal load balancer, use the Scheme parameter.
* `listeners.N` - One or more listeners.
* `load_balancer_name` - The unique name of the load balancer (32 alphanumeric or hyphen characters maximum, but cannot start or end with a hyphen).
* `scheme` - The type of load balancer. Use this parameter only for load balancers in a VPC. To create an internal load balancer, set this parameter to internal.
* `security_groups.N` - (Optional) One or more security groups you want to assign to the load balancer.
  In a VPC, this attribute is required. In the public Cloud, it is optional and default security groups can be applied.
* `subnets.N` - (Optional) One or more subnet IDs in your VPC to attach to the load balancer.
* `tag.N` - (Optional) One or more tags assigned to the load balancer.

## Attributes

* `security_groups_member.N` - The security groups for the load balancer. Valid only for load balancers in a VPC.
* `subnets_member.N` - The IDs of the subnets for the load balancer.
* `listener_descriptions_member.N` - The listeners for the load balancer.
* `policies` - The policies defined for the load balancer.
* `health_check` - Information about the health checks conducted on the load balancer.
* `instances_member.N` - The IDs of the instances for the load balancer.
* `availability_zones_member.N` - The Availability Zones for the load balancer.
* `scheme` - The type of load balancer. Valid only for load balancers in a VPC.
  If Scheme is internet-facing, the load balancer has a public DNS name that resolves to a public IP address.
  If Scheme is internal, the load balancer has a public DNS name that resolves to a private IP address.
* `source_security_group` - The security group for the load balancer, which you can use as part of your inbound rules for your registered instances.
  To only allow traffic from load balancers, add a security group rule that specifies this source security group as the inbound source.
* `vpc_id` - The ID of the VPC for the load balancer.
* `dns_name` - The DNS name of the load balancer.
* `load_balancer_name` - The unique name of the load balancer.

[See detailed information.](http://docs.outscale.com/api_lbu/operations/Action_CreateLoadBalancer_get.html#_api_lbu-action_createloadbalancer_get)
