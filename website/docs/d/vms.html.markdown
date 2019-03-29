---
layout: "outscale"
page_title: "OUTSCALE: outscale_vms"
sidebar_current: "docs-outscale-datasource-vms"
description: |-
  Provides Outscale instances resource attributes. It can be used to recover attributes of instances not managed in the current configuration file.
---

# outscale_vms

  Provides Outscale instances resource attributes. It can be used to recover attributes of instances not managed in the current configuration file.

## Example Usage

```hcl
resource "outscale_vm" "basic" {
  image_id = "ami-8a6a0120"
	instance_type = "t2.micro"
	key_name = "terraform-basic"
	security_group = ["sg-6ed31f3e"]
}

resource "outscale_vm" "web" {
  image_id = "ami-8a6a0120"
	instance_type = "t2.micro"
	key_name = "terraform-basic"
	security_group = ["sg-6ed31f3e"]
}

data "outscale_vms" "basic_web" {
	instance_id = ["${outscale_vm.basic.id}", "${outscale_vm.web.id}"]
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) One or more filters.
* `instance_id` - (Optional) One or more instance IDs.

See detailed information in [Outscale Instances](https://wiki.outscale.net/display/DOCU/Getting+Information+About+Your+Instances).

## Filters

Use the Filter.N parameter to filter the described instances on the following properties:

* `architecture` - The architecture of the instance (i386 | x86_64).
* `availability-zone` - The Availability Zone of the instance.
* `block-device-mapping.attach-time` - The attach time for a BSU volume mapped to the instance (for example, 2016-01-23T18` -45` -30.000Z).
* `block-device-mapping.delete-on-termination` - Indicates whether the BSU volume is deleted when terminating the instance.
* `block-device-mapping.device-name` - The device name for a BSU volume (for example, /dev/sdh or xvdh).
* `block-device-mapping.status` - The status for the BSU volume (attaching | attached| detaching | detached).
* `block-device-mapping.volume-id` - The volume ID of a BSU volume.
* `client-token` - The idempotency token provided when launching the instance.
* `dns-name` - The public DNS name of the instance.
* `group-id` - The ID of the security group for the instance (only in the public Cloud).
* `group-name` - The name of the security group for the instance (only in the public Cloud).
* `hypervisor` - The hypervisor type of the instance (ovm | xen).
* `image-id` - The ID of the image used to launch the instance.
* `instance-id` - The ID of the instance.
* `instance-lifecycle` - Indicates whether the instance is a Spot Instance (spot).
* `instance-state-code` - The state of the instance (a 16-bit unsigned integer). The high byte is an opaque internal value you should ignore. The low byte is set based on the state represented. The valid values are 0 (pending), 16 (running), 32 (shutting-down), 48 (terminated), 64 (stopping), and 80 (stopped).
* `instance-state-name` - The state of the instance (pending | running | shutting-down | terminated | stopping | stopped).
* `instance-type` - The instance type (for example, t2.micro).
* `instance.group-id` - The ID of the security group for the instance.
* `instance.group-name` - The name of the security group for the instance.
* `ip-address` - The public IP address of the instance.
* `kernel-id` - The ID of the kernel.
* `key-name` - The name of the keypair used when launching the instance.
* `launch-index` - The index for the instance when launching a group of several instances (for example, 0, 1, 2, and so on).
* `launch-time` - The time when the instance was launched.
* `monitoring-state` - Indicates whether monitoring is enabled for the instance (disabled | enabled).
* `owner-id` - The Outscale account ID of the instance owner.
* `placement-group-name` - The name of the placement group for the instance.
* `platform` - The platform. Use windows if you have Windows instances. Otherwise, leave this filter blank.
* `private-dns-name` - The private DNS name of the instance.
* `private-ip-address` - The private IP address of the instance.
* `product-code` - The product code associated with the OMI used to launch the instance.
* `ramdisk-id` - The ID of the RAM disk.
* `reason` - The reason explaining the current state of the instance. This filter is like the state-reason-code one.
* `requester-id` - The ID of the entity that launched the instance (for example, Cockpit or Auto Scaling).
* `reservation-id` - The ID of the reservation of the instance, created every time you launch an instance. This reservation ID can can be associated with several instances when you lauch a group of instances using the same launch request.
* `root-device-name` - The name of the root device for the instance (for example, /dev/vda1).
* `root-device-type` - The root device type used by the instance (always ebs).
* `source-dest-check` - If true, the source/destination checking is enabled. If false, it is disabled. This value must be false for a NAT instance to perform NAT (network address translation) in a VPC.
* `spot-instance-request-id` - The ID of the Spot Instance request.
* `state-reason-code` - The reason code for the state change.
* `state-reason-message` - A message describing the state change.
* `subnet-id` - The ID of the subnet for the instance.
* `tag` -key=value` - The key/value combination of a tag that is assigned to the resource.
* `tag-key` - The key of a tag that is assigned to the resource. You can use this filter alongside the tag-value filter. In that case, you filter the resources corresponding to each tag, regardless of the other filter.
* `tag-value` - The value of a tag that is assigned to the resource. You can use this filter alongside the tag-key filter. In that case, you filter the resources corresponding to each tag, regardless of the other filter.
* `tenancy` - The tenancy of an instance (dedicated | default | host).
* `virtualization-type` - The virtualization type of the instance (always hvm).
* `vpc-id` - The ID of the VPC in which the instance is running.
* `network-interface.description` - The description of the network interface.
* `network-interface.subnet-id` - The ID of the subnet for the network interface.
* `network-interface.vpc-id` - The ID of the VPC for the network interface.
* `network-interface.network-interface-id` - The ID of the network interface.
* `network-interface.owner-id` - The ID of the owner of the network interface.
* `network-interface.availability-zone` - The Availability Zone for the network interface.
* `network-interface.requester-id` - The requester ID for the network interface.
* `network-interface.requester-managed` - Indicates whether the network interface is managed by Outscale.
* `network-interface.status` - The status of the network interface (available | in-use).
* `network-interface.mac-address` - The MAC address of the network interface.
* `network-interface.private-dns-name` - The private DNS name of the network interface.
* `network-interface.source-dest-check` - If true, the source/destination checking of the network interface is enabled. If false, it is disabled. The value must be false for the network interface to perform NAT (network address translation) in a VPC.
* `network-interface.group-id` - The ID of a security group associated with the network interface.
* `network-interface.group-name` - The name of a security group associated with the network interface.
* `network-interface.attachment.attachment-id` - The ID of the interface attachment.
* `network-interface.attachment.instance-id` - The ID of the instance the network interface is attached to.
* `network-interface.attachment.instance-owner-id` - The account ID of the ID of the instance the network interface is attached to.
* `network-interface.addresses.private-ip-address` - The private IP address associated with the network interface.
* `network-interface.attachment.device-index` - The device index the network interface is attached to.
* `network-interface.attachment.status` - The status of the attachment (attaching | attached | detaching | detached).
* `network-interface.attachment.attach-time` - The time when the network interface was attached to an instance.
* `network-interface.attachment.delete-on-termination` - Indicates whether the attachment is deleted when terminating an instance.
* `network-interface.addresses.primary` - Indicates whether the IP address of the network interface is the primary private IP address.
* `network-interface.addresses.association.public-ip` - The ID of the association of an External IP address with a network interface.
* `network-interface.addresses.association.ip-owner-id` - The account ID of the owner of the private IP address associated with the network interface.
* `network-interface.association.public-ip` - The External IP address associated with the network interface.
* `network-interface.association.ip-owner-id` - The account ID of the owner of the External IP address associated with the network interface.
* `network-interface.association.allocation-id` - The allocation ID. This ID is returned when you allocate the External IP address for your network interface.
* `network-interface.association.association-id` - The association ID. This ID is returned when the network interface is associated with an IP address.

## Attributes Reference

The following attributes are exported:

* `reservation_set` - Zero or more reservations.
* `request_id`- The ID of the request.



See detailed information in [Describe Instances](http://docs.outscale.com/api_fcu/operations/Action_DescribeInstances_get.html#_api_fcu-action_describeinstances_get).
