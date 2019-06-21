---
layout: "outscale"
page_title: "OUTSCALE: outscale_volumes_link"
sidebar_current: "docs-outscale-resource-volumes-link"
description: |-
  Provides an Outscale Volume Link resource. This allows volumes link to be created, deleted, described and imported.
---

# outscale_volumes_link

  Provides an Outscale Volume Link resource. This allows volumes to be created, deleted, described and imported. Instances also support [provisioning](/docs/provisioners/index.html).

## Example Usage

```hcl
resource "outscale_vm" "web" {
	image_id = "ami-8a6a0120"
	vm_type  = "t1.micro"
}

resource "outscale_volume" "example" {
  subregion_name = "eu-west-2a"
  size = 1
}
resource "outscale_volume_link" "ebs_att" {
  device_name = "/dev/sdh"
	volume_id = "${outscale_volume.example.id}"
	vm_id = "${outscale_vm.web.id}"
}
```

## Argument Reference

The following arguments are supported:

* `device_name` - The instance device name.
* `vm_id` - The ID of the instance you want to attach the volume to.
* `volume_id` - The ID of the volume you want to attach.

## Attributes Reference

The following attributes are exported:

* `delete_on_vm_deletion` - Indicates whether the volume is deleted when terminating the instance
* `device_name` - The instance device name.
* `vm_id` -	The ID of the instance the volume is attached to.
* `state` - The attachment state of the volume (attaching | detaching | attached | detached).
* `volume_id` - The ID of the volume.
* `request_id` - The ID of the request.
