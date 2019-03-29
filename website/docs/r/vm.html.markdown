---
layout: "outscale"
page_title: "OUTSCALE: outscale_vm"
sidebar_current: "docs-outscale-resource-vm"
description: |-
  Provides an Outscale instance resource. This allows instances to be created, updated, and deleted. Instances also support provisioning.
---

# outscale_vm

Provides an Outscale instance resource. This allows instances to be created, updated,
and deleted. Instances also support [provisioning](/docs/provisioners/index.html).

## Example Usage

```hcl
data "outscale_image" "centos73" {
  most_recent = true
  executable_by = ["self"]

  filter {
    name = "owner"
    value = ["Outscale"]
  }

  filter {
    name = "description"
    value = ["Centos 7.3*"]
  }
}
/* instance creation */
resource "outscale_vm" "web" {
  image_id = "${data.outscale_image.centos73.image_id}"
  instance_type = "t2.micro"
}
```

## Argument Reference

The following arguments are supported:

* `block_device_mapping` - (Optional) The block device mapping of the instance.
* `client_token` - (Optional) A unique identifier which enables you to manage the idempotency.
* `disable_api_termination` - (Optional) If true, you cannot terminate the instance using Cockpit, the CLI or the API. If false, you can.
* `ebs_optimized` - (Optional) If true, the instance is created with optimized BSU I/O. All Outscale instances have optimized BSU I/O.
* `image_id` - (Required) The ID of the OMI. You can find the list of OMIs by calling the DescribeImages method.
* `instance_initiated_shutdown_behavior` - (Optional) The instance behavior when you stop or terminate it. By default or if set to stop, the instance stops. If set to restart, the instance stops then automatically restarts. If set to terminate, the instance stops and is terminated.
* `instance_type` - (Optional) The type of instance. For more information, see Instance Types.
* `key_name` - (Optional) The name of the keypair.
* `network_interface` - (Optional) One or more network interfaces.
* `placement` - (Optional) A specific placement where you want to create the instances (for example, Availability Zone, dedicated host, affinity criteria and so on).
* `private_ip_address` - (Optional) In a VPC, the unique primary IP address. The IP address must come from the IP address range of the subnet.
* `private_ip_addresses` - (Optional) In a VPC, the list of primary IP addresses when you create several instances. The IP addresses must come from the IP address range of the subnet.
* `security_group` - (Optional) One or more security group names.
* `security_group_id` - (Optional) One or more security group IDs.
* `subnet_id` - (Optional) In a VPC, the ID of the subnet in which you want to launch the instance.
* `user_data` - (Optional) Data or a script used to add a specific configuration to the instance when launching it. If you are not using a command line tool, this must be base64-encoded.

See detailed information in [Outscale Instances](https://wiki.outscale.net/display/DOCU/Getting+Information+About+Your+Instances).


## Attributes Reference

The following attributes are exported:

* `group_set` - One or more security groups.
* `instances_set` - One or more instances.
* `owner_id` - The ID of the account which has reserved the instances.
* `password_data` - Password for windows environments.
* `requester_id` - The ID of the requester.
* `request_id` - The ID of the request.
* `reservation_id` - Zero or more reservations, giving you information about your request.

See detailed information in [Instances](http://docs.outscale.com/api_fcu/index.html#_instances).
