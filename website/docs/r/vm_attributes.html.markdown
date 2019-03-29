---
layout: "outscale"
page_title: "OUTSCALE: outscale_vm_attributes"
sidebar_current: "docs-outscale-vm-attributes"
description: |-
  Creates an OMI from an existing instance which is either running or stopped.
---

# outscale_image

Creates an OMI from an existing instance which is either running or stopped.

## Example Usage

```hcl
resource "outscale_vm" "basic" {
	image_id = "ami-8a6a0120"
	instance_type = "t2.micro"
	key_name = "terraform-basic"
	security_group = ["sg-6ed31f3e"]
}

resource "outscale_image" "foo" {
	name = "tf-outscale-image-name"
	instance_id = "${outscale_vm.basic.id}"
}
```

## Argument Reference

The following arguments are supported:

* `Attribute`	(Optional) The name of instance attribute (userData | rootDeviceName | instanceType | groupSet | ebsOptimized | sourceDestCheck | deleteOnTermination | disableApiTermination | instanceInitiatedShutdownBehavior).
* `BlockDeviceMapping` (Optional)	The block device mapping of the instance.
* `DisableApiTermination`	(Optional) If true, you cannot terminate the instance using Cockpit, the CLI or the API. If false, you can.
* `EbsOptimized` (Optional)	If true, the instance is optimized for BSU I/O. All Outscale instances have optimized BSU I/O.
* `GroupId` (Optional)	A list of security groups IDs associated with the instance.	false	string	outscale_firewall_rules_set
* `InstanceId` (Required)	The ID of the instance.	true	string	outscale_vm 
* `InstanceInitiatedShutdownBehavior` (Optional)	The instance behavior when you stop or terminate it.
* `InstanceType` (Optional)	The type of instance. For more information, see Instance Types. 
* `SourceDestCheck` (Optional)	If true, the source/destination checking is enabled. If false, it is disabled. This value must be false for a NAT instance to perform NAT (network address translation).
* `UserData` (Optional)	The base64-encoded MIME user data.	false	BlobAttributeValue	- 
* `Value`	(Optional) The new value for the instance attribute.



# Attributes


* `block_device_mapping` - The block device mapping of the instance.
* `disable_api_termination` - If true, you cannot terminate the instance using Cockpit, the CLI or the API. If false, you can.
* `ebs_optimized` - Indicates whether the instance is optimized for BSU I/O.
* `group_set` - The security groups associated with the instance.
* `instance_id` - The ID of the instance.
* `instance_initiated_shutdown_behavior` - Indicates whether the instance stops, terminates or restarts when you stop or terminate it. 
* `instance_type` - The type of instance.
* `kernel` -
* `product_codes` -
* `ramdisk` - The ID of the RAM disk.
* `root_device_name` - The name of the root device.
* `source_dest_check` - (VPC only) If true, the source/destination checking is enabled. If false, it is disabled. This value must be false for a NAT instance to perform Network Address Translation (NAT) in a VPC.
* `sriov_net_support` -
* `user_data` - The Base64-encoded MIME user data.
* `request_id` -


See detailed information in [Describe VM Attibutes](http://docs.outscale.com/api_fcu/operations/Action_DescribeInstanceAttribute_get.html#_api_fcu-action_describeinstanceattribute_get).

