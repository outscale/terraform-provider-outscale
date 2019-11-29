---
layout: "outscale"
page_title: "3DS OUTSCALE: outscale_load_balancer"
sidebar_current: "docs-outscale-resource-load-balancer"
description: |-
  [Manages a load balancer.]
---

# outscale_load_balancer

Manages a load balancer.
For more information on this resource, see the [User Guide](https://wiki.outscale.net/display/EN/About+Load+Balancers).
For more information on this resource actions, see the [API documentation](https://docs-beta.outscale.com/#3ds-outscale-api-loadbalancer).

## Example Usage

```hcl
[exemple de code]
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
* `security_groups` - (Optional) One or more IDs of security groups you want to assign to the load balancer.<br />
In a Net, this attribute is required. In the public Cloud, it is optional and default security groups can be applied.
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
    * `security_group_name` - (Public Cloud only) The name of the security group.
  * `subnets` - The IDs of the Subnets for the load balancer.
  * `subregion_names` - One or more names of Subregions for the load balancer.
  * `tags` - One or more tags associated with the load balancer.
    * `key` - The key of the tag, with a minimum of 1 character.
    * `value` - The value of the tag, between 0 and 255 characters.
