---
layout: "outscale"
page_title: "3DS OUTSCALE: outscale_load_balancer"
sidebar_current: "outscale-load-balancer"
description: |-
  [Manages a load balancer.]
---

# outscale_load_balancer Resource

Manages a load balancer.
For more information on this resource, see the [User Guide](https://wiki.outscale.net/display/EN/About+Load+Balancers).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-loadbalancer).

## Example Usage

```hcl
# Create a load balancer in the public Cloud

resource "outscale_load_balancer" "load_balancer01" {
    load_balancer_name = "terraform-public-load-balancer"
    subregion_names    = ["${var.region}a"]
    listeners {
        backend_port           = 8080
        backend_protocol       = "HTTP"
        load_balancer_protocol = "HTTP"
        load_balancer_port     = 8080
      }
    tags {
        key   = "name"
        value = "terraform-public-load-balancer"
      }
}

# Create a load balancer in a Net

#resource "outscale_net" "net01" {
#    ip_range = "10.0.0.0/16"
#}

#resource "outscale_subnet" "subnet01" {
#    net_id   = outscale_net.net01.net_id
#    ip_range = "10.0.0.0/24"
#    tags {
#        key   = "Name"
#        value = "terraform-subnet-for-internal-load-balancer"
#    }
#}

#resource "outscale_security_group" "security_group01" {
#    description         = "Terraform security group for internal load balancer"
#    security_group_name = "terraform-security-group-for-internal-load-balancer"
#    net_id              = outscale_net.net01.net_id
#     tags {
#         key   = "Name"
#         value = "terraform-security-group-for-internal-load-balancer"
#    }
#}

resource "outscale_load_balancer" "load_balancer02" {
    load_balancer_name = "terraform-private-load-balancer"
    listeners {
        backend_port           = 80
        backend_protocol       = "TCP"
        load_balancer_protocol = "TCP"
        load_balancer_port     = 80
      }
    subnets            = [outscale_subnet.subnet01.subnet_id]
    security_groups    = [outscale_security_group.security_group01.security_group_id]
    load_balancer_type = "internal"
    tags {
        key   = "name"
        value = "terraform-private-load-balancer"
      }
}

# Create an internet-facing load balancer in a Net

#resource "outscale_net" "net02" {
#    ip_range = "10.0.0.0/16"
#}

#resource "outscale_subnet" "subnet02" {
#    net_id   = outscale_net.net02.net_id
#    ip_range = "10.0.0.0/24"
#    tags {
#        key   = "Name"
#        value = "terraform-security-group-for-load-balancer"
#    }
#}

#resource "outscale_internet_service" "internet_service01" {
#    depends_on = "outscale_net.net02"
#}

#resource "outscale_internet_service_link" "internet_service_link01" {
#    internet_service_id = outscale_internet_service.internet_service01.internet_service_id
#    net_id              = outscale_net.net02.net_id
#}

#resource "outscale_route" "route01" {
#    gateway_id           = outscale_internet_service.internet_service01.id
#    destination_ip_range = "10.0.0.0/0"
#    route_table_id       = outscale_route_table.route_table01.route_table_id
#}

#resource "outscale_route_table" "route_table01" {
#    net_id = outscale_net.net02.net_id
#    tags {
#        key   = "name"
#        value = "terraform-route-table-for-load-balancer"
#    }
#}

#resource "outscale_route_table_link" "route_table_link01" {
#    route_table_id  = outscale_route_table.route_table01.route_table_id
#    subnet_id       = outscale_subnet.subnet02.subnet_id
#}

resource "outscale_load_balancer" "load_balancer03" {
    load_balancer_name = "terraform-internet-facing-private-load-balancer"
    listeners {
      backend_port           = 80
      backend_protocol       = "TCP"
      load_balancer_protocol = "TCP"
      load_balancer_port     = 80
     }
    listeners {
     backend_port            = 8080
     backend_protocol        = "HTTP"
     load_balancer_protocol  = "HTTP"
     load_balancer_port      = 8080
     }
    subnets            = [outscale_subnet.subnet02.subnet_id]
    load_balancer_type = "internet-facing"
    tags {
        key   = "name"
        value = "terraform-internet-facing-private-load-balancer"
     }
    depends_on = [outscale_route.route01,outscale_route_table_link.route_table_link01]
}
```

## Argument Reference

The following arguments are supported:

* `listeners` - (Required) One or more listeners to create.
  * `backend_port` - (Optional) The port on which the back-end VM is listening (between `1` and `65535`, both included).
  * `backend_protocol` - (Optional) The protocol for routing traffic to back-end VMs (`HTTP` \| `HTTPS` \| `TCP` \| `SSL` \| `UDP`).
  * `load_balancer_port` - (Optional) The port on which the load balancer is listening (between `1` and `65535`, both included).
  * `load_balancer_protocol` - (Optional) The routing protocol (`HTTP` \| `HTTPS` \| `TCP` \| `SSL` \| `UDP`).
  * `server_certificate_id` - (Optional) The ID of the server certificate.
* `load_balancer_name` - (Required) The unique name of the load balancer (32 alphanumeric or hyphen characters maximum, but cannot start or end with a hyphen).
* `load_balancer_type` - (Optional) The type of load balancer: `internet-facing` or `internal`. Use this parameter only for load balancers in a Net.
* `security_groups` - (Optional) (Net only) One or more IDs of security groups you want to assign to the load balancer. If not specified, the default security group of the Net is assigned to the load balancer.
* `subnets` - (Optional) One or more IDs of Subnets in your Net that you want to attach to the load balancer.
* `subregion_names` - (Optional) One or more names of Subregions (currently, only one Subregion is supported). This parameter is not required if you create a load balancer in a Net. To create an internal load balancer, use the `LoadBalancerType` parameter.
* `tags` - (Optional) One or more tags assigned to the load balancer.
  * `key` - (Optional) The key of the tag, with a minimum of 1 character.
  * `value` - (Optional) The value of the tag, between 0 and 255 characters.

## Attribute Reference

The following attributes are exported:

* `load_balancer` - Information about the load balancer.
  * `access_log` - Information about access logs.
    * `is_enabled` - If `true`, access logs are enabled for your load balancer. If `false`, they are not. If you set this to `true` in your request, the `OsuBucketName` parameter is required.
    * `osu_bucket_name` - The name of the Object Storage Unit (OSU) bucket for the access logs.
    * `osu_bucket_prefix` - The path to the folder of the access logs in your Object Storage Unit (OSU) bucket (by default, the `root` level of your bucket).
    * `publication_interval` - The time interval for the publication of access logs in the Object Storage Unit (OSU) bucket, in minutes. This value can be either 5 or 60 (by default, 60).
  * `application_sticky_cookie_policies` - The stickiness policies defined for the load balancer.
    * `cookie_name` - The name of the application cookie used for stickiness.
    * `policy_name` - The mnemonic name for the policy being created. The name must be unique within a set of policies for this load balancer.
  * `backend_vm_ids` - One or more IDs of back-end VMs for the load balancer.
  * `dns_name` - The DNS name of the load balancer.
  * `health_check` - Information about the health check configuration.
    * `check_interval` - The number of seconds between two pings (between `5` and `600` both included).
    * `healthy_threshold` - The number of consecutive successful pings before considering the VM as healthy (between `2` and `10` both included).
    * `path` - The path for HTTP or HTTPS requests.
    * `port` - The port number (between `1` and `65535`, both included).
    * `protocol` - The protocol for the URL of the VM (`HTTP` \| `HTTPS` \| `TCP` \| `SSL` \| `UDP`).
    * `timeout` - The maximum waiting time for a response before considering the VM as unhealthy, in seconds (between `2` and `60` both included).
    * `unhealthy_threshold` - The number of consecutive failed pings before considering the VM as unhealthy (between `2` and `10` both included).
  * `listeners` - The listeners for the load balancer.
    * `backend_port` - The port on which the back-end VM is listening (between `1` and `65535`, both included).
    * `backend_protocol` - The protocol for routing traffic to back-end VMs (`HTTP` \| `HTTPS` \| `TCP` \| `SSL` \| `UDP`).
    * `load_balancer_port` - The port on which the load balancer is listening (between 1 and `65535`, both included).
    * `load_balancer_protocol` - The routing protocol (`HTTP` \| `HTTPS` \| `TCP` \| `SSL` \| `UDP`).
    * `policy_names` - The names of the policies. If there are no policies enabled, the list is empty.
    * `server_certificate_id` - The ID of the server certificate.
  * `load_balancer_name` - The name of the load balancer.
  * `load_balancer_sticky_cookie_policies` - The policies defined for the load balancer.
    * `policy_name` - The name of the stickiness policy.
  * `load_balancer_type` - The type of load balancer. Valid only for load balancers in a Net.<br />
If `LoadBalancerType` is `internet-facing`, the load balancer has a public DNS name that resolves to a public IP address.<br />
If `LoadBalancerType` is `internal`, the load balancer has a public DNS name that resolves to a private IP address.
  * `net_id` - The ID of the Net for the load balancer.
  * `security_groups` - One or more IDs of security groups for the load balancers. Valid only for load balancers in a Net.
  * `source_security_group` - Information about the source security group of the load balancer, which you can use as part of your inbound rules for your registered VMs.<br />
To only allow traffic from load balancers, add a security group rule that specifies this source security group as the inbound source.
    * `security_group_account_id` - The account ID of the owner of the security group.
    * `security_group_name` - The name of the security group.
  * `subnets` - The IDs of the Subnets for the load balancer.
  * `subregion_names` - One or more names of Subregions for the load balancer.
  * `tags` - One or more tags associated with the load balancer.
    * `key` - The key of the tag, with a minimum of 1 character.
    * `value` - The value of the tag, between 0 and 255 characters.

## Import

A load balancer can be imported using its name. For example:

```

$ terraform import outscale_load_balancer.ImportedLbu Name-of-the-Lbu

```