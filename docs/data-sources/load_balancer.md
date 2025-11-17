---
layout: "outscale"
page_title: "OUTSCALE: outscale_load_balancer"
subcategory: "OUTSCALE API"
sidebar_current: "outscale-load-balancer"
description: |-
  [Provides information about a load balancer.]
---

# outscale_load_balancer Data Source

Provides information about a load balancer.

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Load-Balancers.html).  
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-loadbalancer).

## Example Usage

```hcl
data "outscale_load_balancer" "load_balancer01" {
    filter {
        name   = "load_balancer_names"
        values = ["load_balancer01"]
    }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) A combination of a filter name and one or more filter values. You can specify this argument for as many filter names as you need. The filter name can be any of the following:
    * `load_balancer_names` - (Optional) The names of the load balancers.

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
