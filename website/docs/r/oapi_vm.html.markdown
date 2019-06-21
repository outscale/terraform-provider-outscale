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
resource "outscale_vm" "basic" {
	image_id = "ami-8a6a0120"
	vm_type = "t2.micro"
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
* `vm_type` - (Optional) The type of instance. For more information, see Instance Types.
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

* `group_set.N` - One or more security groups.
  * `group_id` - The ID of the security group.
  * `group_name` - The name of the security group.
* `instances_set.N` - One or more instances.
  * `ami_launch_index` - The launch index of the OMI.
  * `architecture` - The architecture of the image.
  * `block_device_mapping.N` - One or more entries of block device mapping.
    * `device_name` - The name of the instance device name.
    * `ebs` - One or more parameters used to automatically set up volumes when the instance is launched.
      * `delete_on_termination` - If true, the volume is deleted when the instance is terminated.
      * `status` - The attachment state (attaching | attached | detaching | detached).   
      * `volume_id` - The ID of the volume.
  * `client_token` - A unique identifier which enables you to manage the idempotency. 
  * `dns_name` - The name of the public DNS assigned to the instance. 
  * `ebs_optimized` - If true, the instance is created with optimized BSU I/O. 
  * `group_set.N` - One and more security groups for the instance.
    * `group_id` - The ID of the security group.
    * `group_name` - The name of the security group.
  * `hypervisor` - The hypervisor type of the instance.
  * `iam_instance_profile` - The EIM instance profile associated with the instance.
    * `arn` - The unique identifier of the ressource (between 20 and 2048 characters).
    * `id` - The ID of the instance profile.
  * `image_id` - The ID of the OMI.
  * `instance_id` - The ID of the instance.
  * `instance_lifecycle` - Indicates whether it is a spot instance. 
  * `instance_state` - The current state of the instance.
    * `code` - The code of the state of the instance (0 pending | 16 running | 32 shutting­down | 48 terminated | 64 stopping | 80 stopped)
    * `name` - The state of the instance (pending | running | shutting­down | terminated | stopping | stopped). 
  * `vm_type` - The type of instance.
  * `ip_address` - The public IP address assigned to the instance.
  * `kernel_id` - The ID of the associated kernel.
  * `key_name` - The name of the keypair.
  * `monitoring` - The monitoring information for the instance.
    * `state` - The state of detail monitoring (enabled | disabled | disabling | pending)
  * `network_interface_set.N` - In a VPC, one or more network interfaces for the instance.
    * `association` - Information about an External IP associated with the interface.
      * `ip_owner_id` - The account ID of the owner of the External IP address.
      * `public_dns_name` - The name of the public DNS.
      * `public_ip` - The public IP address or the External IP address associated with the network interface.
    * `attachment` - The attachment of the network interface.
      * `attachment_id` - The ID of the network interface attachment.
      * `delete_on_termination` - If true, the network interface is deleted when the instance is terminated. 
      * `device_index` - The index of the instance device for the attachment.
      * `attachment.status` - The state of the netowrk interface (attaching | attached | detaching | detached).
    * `description` - The description of the network interface. 
    * `group_set.N` - One and more security groups for the instance.
      * `group_id` - The ID of the security group.
      * `group_name` - The name of the security group.
    * `mac_address` - The MAC address of the interface.
    * `network_interface_id` - The ID of the network interface.
    * `owner_id` - The account ID of the owner of the instance reservation.
    * `private_dns_name` - The name of the private DNS assigned to the instance. 
    * `private_ip_address` - The private IP address assigned to the instance. 
    * `private_ip_addresses_set.N` - The private IP addresses assigned to the network interface.
      * `association` - Information about an External IP address associated with the interface.
        * `ip_owner_id` - The account ID of the owner of the External IP address.
        * `public_dns_name` - The name of the public DNS.
        * `public_ip` - The public IP address or the External IP address associated with the network interface.
      * `primary` - If true, the IP address is the primary address of the network interface. 
      * `private_dns_name` - The name of the private DNS.
      * `private_ip_address` - The private IP address of the network interface.
    * `source_dest_check` - If true in a VPC, the check to perform NAT is enabled. 
    * `status` - The state of the interface.
    * `subnet_id` - In a VPC, the ID of the subnet in which to launch the instance. 
    * `vpc_id` - The ID of the VPC in which the instance is launched.
  * `placement` - A specific placement where you want to create the instances. 
    * `tenancy` - The tenancy of the instance (default|dedicated).
  * `platform` - Indicates whether it is a Windows instance.
  * `private_dns_name` - The name of the private DNS assigned to the instance. 
  * `private_ip_address` - The private IP address assigned to the instance. 
  * `product_codes.N` - The code of the product attached to the instance.
    * `product_code` - The code of the product. (001 Linux/Unix | 002 Windows | 003 MapR | 004 Linux/oracle | 005 Windows 10)
    * `type` - The type of product code.
  * `ramdisk_id` - The ID of the associated RAM disk. 
  * `reason` - Information about the latest state change. 
  * `root_device_name` - The name of the root device.
  * `root_device_type` - The type of root device used by the OMI. 
  * `source_dest_check` - If true in a VPC, the check to perform NAT is enabled. 
  * `spot_instance_request_id` - The ID of the spot instance request. 
  * `sriov_net_support` - If true, the enhanced networking is enabled. 
  * `state_reason` - Information about the latest state change.
    * `code` - The code of the change of state.
    * `message` - The message explaining the change of state.
  * `subnet_id` - In a VPC, the ID of the subnet in which you want to launch the instance. 
  * `tag_set.N` - One or more tags associated with the instance.
    * `key` - The key of the tag.
    * `value` - The value of the tag.
  * `virtualization_type` - The virtualization type.
  * `vpc_id` - The ID of the VPC in which the instance is launched.
* `owner_id` - The ID of the account which has reserved the instances.
* `password_data` - Password for windows environments.
* `requester_id` - The ID of the requester.
* `reservation_id` - Zero or more reservations, giving you information about your request.

See detailed information in [Instances](http://docs.outscale.com/api_fcu/index.html#_instances).
