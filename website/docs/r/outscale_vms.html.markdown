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
resource "outscale_vms" "web" {
  ami = "${data.outscale_image.centos73.image_id}"
  instance_type = "t2.micro"
}
```

## Argument Reference

The following arguments are supported:

* `amiLaunchIndex` - (Optional)The launch index of the OMI.
* `architecture` - (Optional)The architecture of the image.
* `blockDeviceMapping.N` - (Optional)One or more entries of block device mapping.
* `clientToken` - (Optional)A unique identifier which enables you to manage the idempotency.
* `dnsName` - (Optional)The name of the public DNS assigned to the instance.
* `ebsOptimized` - (Optional)If true, the instance is created with optimized BSU I/O. All Outscale instances have optimized BSU I/O.
* `groupSet.N` - (Optional)One and more security groups for the instance.
* `hypervisor` - (Optional)The hypervisor type of the instance.
* `iamInstanceProfile` - (Optional)The EIM instance profile associated with the instance.
* `imageId` - (Optional)The ID of the OMI.
* `instanceId` - (Optional)The ID of the instance.
* `instanceLifecycle` - (Optional)Indicates whether it is a spot instance.
* `instanceState` - (Optional)The current state of the instance.
* `instanceType` - (Optional)The type of instance.
* `ipAddress` - (Optional)The public IP address assigned to the instance.
* `kernelId` - (Optional)The ID of the associated kernel.
* `keyName` - (Optional)The name of the keypair.
* `monitoring` - (Optional)The monitoring information for the instance.
* `networkInterfaceSet.N` - (Optional)In a VPC, one or more network interfaces for the instance.
* `placement` - (Optional)A specific placement where you want to create the instances.
* `platform` - (Optional)Indicates whether it is a Windows instance.
* `privateDnsName` - (Optional)The name of the private DNS assigned to the instance.
* `privateIpAddress` - (Optional)The private IP address assigned to the instance.
* `productCodes.N` - (Optional)The code of the product attached to the instance.
* `ramdiskId` - (Optional)The ID of the associated RAM disk.
* `reason` - (Optional)Information about the latest state change.
* `rootDeviceName` - (Optional)The name of the root device.
* `rootDeviceType` - (Optional)The type of root device used by the OMI.
* `sourceDestCheck` - (Optional)If true in a VPC, the check to perform NAT is enabled.
* `spotInstanceRequestId` - (Optional)The ID of the spot instance request.
* `sriovNetSupport` - (Optional)If true, the enhanced networking is enabled.
* `stateReason` - (Optional)Information about the latest state change.
* `subnetId` - (Optional)In a VPC, the ID of the subnet in which you want to launch the instance.
* `tagSet.N` - (Optional)One or more tags associated with the instance.
* `virtualizationType` - (Optional)The virtualization type.
* `vpcId` - (Optional)The ID of the VPC in which the instance is launched.


## Attributes Reference

The following attributes are exported:


