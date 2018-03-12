---
layout: "outscale"
page_title: "OUTSCALE: outscale_images"
sidebar_current: "docs-outscale-datasource-images"
description: |-
    Describes one or more OMIs you can use.


---

# outscale_images

Describes one or more OMIs you can use.
You can filter the described OMIs using the ImageId.N, the Owner.N and the ExecutableBy.N parameters.
You can also use the Filter.N parameter to filter the OMIs on the following properties:

## Example Usage

```hcl
data "outscale_images" "nat_ami" {
    filter {
        name = "architecture"
        values = ["x86_64"]
    }
    filter {
        name = "virtualization-type"
        values = ["hvm"]
    }
    filter {
        name = "root-device-type"
        values = ["ebs"]
    }
    filter {
        name = "block-device-mapping.volume-type"
        values = ["standard"]
    }
}
```

## Argument Reference

The following arguments are supported:

* `dryRun` - (Optional) If set to true, checks whether you have the required permissions to perform the action.
* `executable_by` - (Optional) One or more instance IDs.
* `filter` - (Optional) One or more filters.
* `image_id` - (Optional) One or more OMI IDs.
* `owner_id` - (Optional) The user ID of one or more owners of OMIs. By default, all the OMIs for which you have launch permissions are described.


See detailed information in [Outscale Instances](https://wiki.outscale.net/display/DOCU/Getting+Information+About+Your+Instances).

## Filters

Use the Filter.N parameter to filter the described instances on the following properties:

* `architecture` - The architecture of the instance (i386 | x86_64).
* `block-device-mapping.delete-on-termination` Whether the volume is deleted when terminating the instance.
* `block-device-mapping.device-name` The device name for the volume.
* `block-device-mapping.snapshot-id` The ID of the snapshot used to create the volume.
* `block-device-mapping.volume-size` The size of the volume, in Gibibytes (GiB).
* `block-device-mapping.volume-type` The type of volume (standard | gp2 | io1 | sc1 | st1).
* `description` The description of the OMI, provided when it was created.
* `hypervisor` The hypervisor type (always xen).
* `image-id` The ID of the OMI.
* `image-type` The type of OMI (always machine for official OMIs).
* `is-public` Whether the OMI has public launch permissions.
* `kernel-id` The ID of kernel.
* `manisfest-location` The location of the OMI manifest.
* `name` The name of the OMI, provided when it was created.
* `owner-alias` The account alias of the owner of the OMI.
* `owner-id` The account ID of the owner of the OMI.
* `platform` The platform.
* `product-code` The product code associated with the OMI.
* `ramdisk-id` The ID of the RAM disk.
* `root-device-name` The device name of the root device (for example, /dev/sda1).
* `root-device-type` The type of root device used by the OMI (always ebs).
* `state` The current state of the OMI.
* `tag` The key/value combination of a tag associated with the OMI, in the following format: key=value.
* `tag-key` The key of a tag associated with the OMI, independent of tag-value.
* `tag-value` The value of a tag associated with the OMI, independent of tag-key.
* `virtualization-type` The virtualization type (always hvm).


## Attributes Reference

The following attributes are exported:

* `image_set  ` - Information about one or more OMIs.
* `requester_id` - The ID of the request.

See detailed information in [Describe Instances](http://docs.outscale.com/api_fcu/operations/Action_DescribeImages_get.html#_body_parameter).
