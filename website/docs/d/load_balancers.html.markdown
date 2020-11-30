---
layout: "outscale"
page_title: "3DS OUTSCALE: outscale_load_balancer"
sidebar_current: "outscale-load-balancer"
description: |-
  [Provides information about load balancers.]
---

# outscale_load_balancer Data Source

Provides information about load balancers.
For more information on this resource, see the [User Guide](https://wiki.outscale.net/display/EN/About+Load+Balancers).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-loadbalancer).

## Example Usage

```hcl
data "outscale_load_balancers" "load_balancers01" {
filter {
        name   = "load_balancer_names"
        values = ["load_balancer01","load_balancer02"]
      }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - A combination of a filter name and one or more filter values. You can specify this argument for as many filter names as you need. The filter name can be any of the following:
  * `load_balancer_names` - (Optional) The names of the load balancers.

## Attribute Reference

The following attributes are exported:

* `load_balancers` - Information about one or more load balancers.
  * `access_log` - Information about access logs.
      * `is_enabled` - If `true`, access logs are enabled for your load balancer. If `false`, they are not. If you set this to `true` in your request, the `osu_bucket_name` parameter is required.
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

