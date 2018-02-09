---
layout: "outscale"
page_title: "OUTSCALE: outscale_vms"
sidebar_current: "docs-outscale-resource-vms"
description: |-
  Provides an Outscale instance resource. This allows instances to be created, updated, and deleted. Instances also support provisioning.
---

# outscale_instance

Provides an Outscale instance resource. This allows instances to be created, updated,
and deleted. Instances also support [provisioning](/docs/provisioners/index.html).

## Example Usage

```hcl
data "outscale_image" "centos73" { 
  most_recent = true 
  executable_by = ["self"] 

  filter {
    name = "owner" 
    values = ["Outscale"] 
  }

  filter {
  name = "description" values = ["Centos 7.3*"] 
  } 
} 
/* instance creation */
resource "outscale_vm" "web" { 
  ami = "${data.outscale_image.centos73.image_id}"
  instance_type = "t2.micro" 
}
```

## Argument Reference

The following arguments are supported:

* `BlockDeviceMapping.N` - (Optional) The block device mapping of the instance.
* `ClientToken` - (Optional) A unique identifier which enables you to manage the idempotency.
* `DisableApiTermination` - (Optional) If true, you cannot terminate the instance using Cockpit, the CLI or the API. If false, you can.
* `DryRun` - (Optional) If true, checks whether you have the required permissions to perform the action.
* `EbsOptimized` - (Optional) If true, the instance is created with optimized BSU I/O. All Outscale instances have optimized BSU I/O.
* `ImageId` - (Required) The ID of the OMI. You can find the list of OMIs by calling the DescribeImages method.
* `InstanceInitiatedShutdownBehavior` - (Optional) The instance behavior when you stop or terminate it. By default or if set to stop, the instance stops. If set to restart, the instance stops then automatically restarts. If set to terminate, the instance stops and is terminated.
* `InstanceType` - (Optional) The type of instance. For more information, see Instance Types.
* `KeyName` - (Optional) The name of the keypair.
* `MaxCount` - (Required) The maximum number of instances you want to launch. If all the instances cannot be created, the largest possible number of instances above MinCount are created and launched.
* `MinCount` - (Required) The minimum number of instances you want to launch. If this number of instances cannot be created, FCU does not create and launch any instance.
* `NetworkInterface.N` - (Optional) One or more network interfaces.
* `Placement` - (Optional) A specific placement where you want to create the instances (for example, Availability Zone, dedicated host, affinity criteria and so on).
* `PrivateIpAddress` - (Optional) In a VPC, the unique primary IP address. The IP address must come from the IP address range of the subnet.
* `PrivateIpAddresses` - (Optional) In a VPC, the list of primary IP addresses when you create several instances. The IP addresses must come from the IP address range of the subnet.
* `SecurityGroup.N` - (Optional) One or more security group names.
* `SecurityGroupId.N` - (Optional) One or more security group IDs.
* `SubnetId` - (Optional) In a VPC, the ID of the subnet in which you want to launch the instance.
* `UserData` - (Optional) Data or a script used to add a specific configuration to the instance when launching it. If you are not using a command line tool, this must be base64-encoded.


## Attributes Reference

The following attributes are exported:

* `groupSet.N` - (Optional) One or more security groups.
* `instancesSet.N` - (Optional) One or more instances.
* `ownerId` - (Optional) The ID of the account which has reserved the instances.
* `requesterId` - (Optional) The ID of the requester.
* `reservationId` - (Optional) Zero or more reservations, giving you information about your request.
