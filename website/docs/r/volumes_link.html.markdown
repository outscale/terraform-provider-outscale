---
layout: "outscale"
page_title: "OUTSCALE: outscale_volume_link"
sidebar_current: "docs-outscale-resource-volume-link"
description: |-
  Provides an Outscale Volume Link resource. This allows volumes link to be created, deleted, described and imported.
---

# outscale_volume_link

  Provides an Outscale Volume Link resource. This allows volumes to be created, deleted, described and imported. Instances also support [provisioning](/docs/provisioners/index.html).

## Example Usage

```hcl
resource "outscale_vm" "web" {
	image_id = "ami-8a6a0120"
	instance_type = "t1.micro"
	tags {
		Name = "HelloWorld"
	}
}
resource "outscale_volume" "example" {
  availability_zone = "eu-west-2a"
	size = 1
}
resource "outscale_volume_link" "ebs_att" {
  device = "/dev/sdh"
	volume_id = "${outscale_volume.example.id}"
	instance_id = "${outscale_vm.web.id}"
}

```

## Argument Reference

The following arguments are supported:

* `device` - The instance device name.
* `instance_id` - The ID of the instance you want to attach the volume to.
* `volume_id` - The ID of the volume you want to attach.

See detailed information in [Outscale Volume](http://docs.outscale.com/api_fcu/operations/Action_AttachVolume_get.html#_api_fcu-action_attachvolume_get).


## Attributes Reference

The following attributes are exported:

* `delete_on_termination` - Indicates whether the volume is deleted when terminating the instance
* `device` - The instance device name.
* `instance_id` -	The ID of the instance the volume is attached to.
* `status` - The attachment state of the volume (attaching | detaching | attached | detached).
* `volume_id` - The ID of the volume.


See detailed information in [Volume Attachment](http://docs.outscale.com/api_fcu/definitions/VolumeAttachment.html#_api_fcu-volumeattachment).
