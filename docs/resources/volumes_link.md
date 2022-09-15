---
layout: "outscale"
page_title: "OUTSCALE: outscale_volumes_link"
sidebar_current: "outscale-volumes-link"
description: |-
  [Manages a volume link.]
---

# outscale_volumes_link Resource

Manages a volume link.
For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Volumes.html).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-volume).

## Example Usage

### Required resources

```hcl
resource "outscale_volume" "volume01" {
	subregion_name = "${var.region}a"
	size           = 40
}

resource "outscale_vm" "vm01" {
	image_id           = var.image_id
	vm_type            = var.vm_type
	keypair_name       = var.keypair_name
	security_group_ids = [var.security_group_id]
}
```

### Link a volume to a VM

```hcl
resource "outscale_volumes_link" "volumes_link01" {
	device_name = "/dev/xvdc"
	volume_id   = outscale_volume.volume01.id
	vm_id       = outscale_vm.vm01.id
}
```

## Argument Reference

The following arguments are supported:

* `device_name` - (Required) The name of the device.
* `vm_id` - (Required) The ID of the VM you want to attach the volume to.
* `volume_id` - (Required) The ID of the volume you want to attach.

## Attribute Reference

The following attributes are exported:

* `delete_on_vm_deletion` - If true, the volume is deleted when terminating the VM. If false, the volume is not deleted when terminating the VM.
* `device_name` - The name of the device.
* `state` - The state of the attachment of the volume (`attaching` | `detaching` | `attached` | `detached`).
* `vm_id` - The ID of the VM.
* `volume_id` - The ID of the volume.

## Import

A volume link can be imported using a volume ID. For example:

```console

$ terraform import outscale_volumes_link.ImportedVolumeLink vol-12345678

```