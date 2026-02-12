---
layout: "outscale"
page_title: "OUTSCALE: outscale_load_balancer"
subcategory: "OUTSCALE API"
sidebar_current: "outscale-load-balancer"
description: |-
  [Manages a load balancer.]
---

# outscale_load_balancer Resource

Manages a load balancer.

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Load-Balancers.html).  
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-loadbalancer).

## Example Usage

### Create a load balancer in the public Cloud

```hcl
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
```

### Create a load balancer in a Net

```hcl
resource "outscale_net" "net01" {
    ip_range = "10.0.0.0/16"
}

resource "outscale_subnet" "subnet01" {
    net_id   = outscale_net.net01.net_id
    ip_range = "10.0.0.0/24"
    tags {
        key   = "Name"
        value = "terraform-subnet-for-internal-load-balancer"
    }
}

resource "outscale_security_group" "security_group01" {
    description         = "Terraform security group for internal load balancer"
    security_group_name = "terraform-security-group-for-internal-load-balancer"
    net_id              = outscale_net.net01.net_id
    tags {
        key   = "Name"
        value = "terraform-security-group-for-internal-load-balancer"
    }
}

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
```

### Create an internet-facing load balancer in a Net

```hcl
resource "outscale_net" "net02" {
    ip_range = "10.0.0.0/16"
}

resource "outscale_subnet" "subnet02" {
    net_id   = outscale_net.net02.net_id
    ip_range = "10.0.0.0/24"
    tags {
        key   = "Name"
        value = "terraform-security-group-for-load-balancer"
    }
}

resource "outscale_internet_service" "internet_service01" {
    depends_on = [outscale_net.net02]
}

resource "outscale_internet_service_link" "internet_service_link01" {
    internet_service_id = outscale_internet_service.internet_service01.internet_service_id
    net_id              = outscale_net.net02.net_id
}

resource "outscale_route_table" "route_table01" {
    net_id = outscale_net.net02.net_id
    tags {
        key   = "name"
        value = "terraform-route-table-for-load-balancer"
    }
}

resource "outscale_route" "route01" {
    gateway_id           = outscale_internet_service.internet_service01.id
    destination_ip_range = "0.0.0.0/0"
    route_table_id       = outscale_route_table.route_table01.route_table_id
}

resource "outscale_route_table_link" "route_table_link01" {
    route_table_id = outscale_route_table.route_table01.route_table_id
    subnet_id      = outscale_subnet.subnet02.subnet_id
}

resource "outscale_load_balancer" "load_balancer03" {
    load_balancer_name = "terraform-internet-private-lb"
    listeners {
        backend_port           = 80
        backend_protocol       = "TCP"
        load_balancer_protocol = "TCP"
        load_balancer_port     = 80
    }
    listeners {
        backend_port           = 8080
        backend_protocol       = "HTTP"
        load_balancer_protocol = "HTTP"
        load_balancer_port     = 8080
    }
    subnets            = [outscale_subnet.subnet02.subnet_id]
    load_balancer_type = "internet-facing"
    public_ip          = "192.0.2.0"
    tags {
        key   = "name"
        value = "terraform-internet-private-lb"
    }
    depends_on = [outscale_route.route01,outscale_route_table_link.route_table_link01]
}
```


## Argument Reference

The following arguments are supported:

* `listeners` - (Required) One or more listeners to create.
    * `backend_port` - (Optional) The port on which the backend VM is listening (between `1` and `65535`, both included).
    * `backend_protocol` - (Optional) The protocol for routing traffic to backend VMs (`HTTP` \| `HTTPS` \| `TCP` \| `SSL`).
    * `load_balancer_port` - (Optional) The port on which the load balancer is listening (between `1` and `65535`, both included).
    * `load_balancer_protocol` - (Optional) The routing protocol (`HTTP` \| `HTTPS` \| `TCP` \| `SSL`).
    * `server_certificate_id` - (Optional) The OUTSCALE Resource Name (ORN) of the server certificate. For more information, see [Resource Identifiers > OUTSCALE Resource Names (ORNs)](https://docs.outscale.com/en/userguide/Resource-Identifiers.html#_outscale_resource_names_orns).<br/>
This parameter is required for `HTTPS` and `SSL` protocols.
* `load_balancer_name` - (Required) The unique name of the load balancer, with a maximum length of 32 alphanumeric characters and dashes (`-`). This name must not start or end with a dash.
* `load_balancer_type` - (Optional) The type of load balancer: `internet-facing` or `internal`. Use this parameter only for load balancers in a Net.
* `public_ip` - (Optional) (internet-facing only) The public IP you want to associate with the load balancer. If not specified, a public IP owned by 3DS OUTSCALE is associated.
* `security_groups` - (Optional) (Net only) One or more IDs of security groups you want to assign to the load balancer. If not specified, the default security group of the Net is assigned to the load balancer.
* `subnets` - (Optional) (Net only) The ID of the Subnet in which you want to create the load balancer. Regardless of this Subnet, the load balancer can distribute traffic to all Subnets. This parameter is required in a Net.
* `subregion_names` - (Optional) (public Cloud only) The Subregion in which you want to create the load balancer. Regardless of this Subregion, the load balancer can distribute traffic to all Subregions. This parameter is required in the public Cloud.
* `tags` - (Optional) A tag to add to this resource. You can specify this argument several times.
    * `key` - (Required) The key of the tag, with a minimum of 1 character.
    * `value` - (Required) The value of the tag, between 0 and 255 characters.

## Attribute Reference

The following attributes are exported:

* `access_log` - Information about access logs.
    * `is_enabled` - If true, access logs are enabled for your load balancer. If false, they are not. If you set this to true in your request, the `osu_bucket_name` parameter is required.
    * `osu_bucket_name` - The name of the OOS bucket for the access logs.
    * `osu_bucket_prefix` - The path to the folder of the access logs in your OOS bucket (by default, the `root` level of your bucket).
    * `publication_interval` - The time interval for the publication of access logs in the OOS bucket, in minutes. This value can be either `5` or `60` (by default, `60`).
* `application_sticky_cookie_policies` - The stickiness policies defined for the load balancer.
    * `cookie_name` - The name of the application cookie used for stickiness.
    * `policy_name` - The mnemonic name for the policy being created. The name must be unique within a set of policies for this load balancer.
* `backend_ips` - One or more public IPs of backend VMs.
* `backend_vm_ids` - One or more IDs of backend VMs for the load balancer.
* `dns_name` - The DNS name of the load balancer.
* `health_check` - Information about the health check configuration.
    * `check_interval` - The number of seconds between two requests (between `5` and `600` both included).
    * `healthy_threshold` - The number of consecutive successful requests before considering the VM as healthy (between `2` and `10` both included).
    * `path` - If you use the HTTP or HTTPS protocols, the request URL path. Always starts with a slash (`/`).
    * `port` - The port number (between `1` and `65535`, both included).
    * `protocol` - The protocol for the URL of the VM (`HTTP` \| `HTTPS` \| `TCP` \| `SSL`).
    * `timeout` - The maximum waiting time for a response before considering the VM as unhealthy, in seconds (between `2` and `60` both included).
    * `unhealthy_threshold` - The number of consecutive failed requests before considering the VM as unhealthy (between `2` and `10` both included).
* `listeners` - The listeners for the load balancer.
    * `backend_port` - The port on which the backend VM is listening (between `1` and `65535`, both included).
    * `backend_protocol` - The protocol for routing traffic to backend VMs (`HTTP` \| `HTTPS` \| `TCP` \| `SSL`).
    * `load_balancer_port` - The port on which the load balancer is listening (between `1` and `65535`, both included).
    * `load_balancer_protocol` - The routing protocol (`HTTP` \| `HTTPS` \| `TCP` \| `SSL`).
    * `policy_names` - The names of the policies. If there are no policies enabled, the list is empty.
    * `server_certificate_id` - The OUTSCALE Resource Name (ORN) of the server certificate. For more information, see [Resource Identifiers > OUTSCALE Resource Names (ORNs)](https://docs.outscale.com/en/userguide/Resource-Identifiers.html#_outscale_resource_names_orns).
* `load_balancer_name` - The name of the load balancer.
* `load_balancer_sticky_cookie_policies` - The policies defined for the load balancer.
    * `policy_name` - The name of the stickiness policy.
* `load_balancer_type` - The type of load balancer. Valid only for load balancers in a Net.<br />
If `load_balancer_type` is `internet-facing`, the load balancer has a public DNS name that resolves to a public IP.<br />
If `load_balancer_type` is `internal`, the load balancer has a public DNS name that resolves to a private IP.
* `net_id` - The ID of the Net for the load balancer.
* `public_ip` - (internet-facing only) The public IP associated with the load balancer.
* `secured_cookies` - Whether secure cookies are enabled for the load balancer.
* `security_groups` - One or more IDs of security groups for the load balancers. Valid only for load balancers in a Net.
* `source_security_group` - Information about the source security group of the load balancer, which you can use as part of your inbound rules for your registered VMs.<br />
To only allow traffic from load balancers, add a security group rule that specifies this source security group as the inbound source.
    * `security_group_account_id` - The account ID of the owner of the security group.
    * `security_group_name` - The name of the security group.
* `subnets` - The ID of the Subnet in which the load balancer was created.
* `subregion_names` - The ID of the Subregion in which the load balancer was created.
* `tags` - One or more tags associated with the load balancer.
    * `key` - The key of the tag, with a minimum of 1 character.
    * `value` - The value of the tag, between 0 and 255 characters.

## Timeouts

The `timeouts` block enables you to configure [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - Defaults to 10 minutes.
* `read` - Defaults to 5 minutes.
* `update` - Defaults to 10 minutes.
* `delete` - Defaults to 5 minutes.

## Import

A load balancer can be imported using its name. For example:

```console

$ terraform import outscale_load_balancer.ImportedLbu Name-of-the-Lbu

```