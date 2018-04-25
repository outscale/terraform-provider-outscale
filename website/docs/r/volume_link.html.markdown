---
layout: "outscale"
page_title: "OUTSCALE: outscale_volume_link"
sidebar_current: "docs-outscale-resource-volume-link"
description: |-
  Provides an Outscale Volume Link resource. This allows volumes link to be created, deleted, described and imported.
---

# outscale_vm

  Provides an Outscale Volume Link resource. This allows volumes to be created, deleted, described and imported. Instances also support [provisioning](/docs/provisioners/index.html).

## Example Usage

```hcl
data "outscale_volume" "outscale_volume" {
  most_recent = true

  filter {
    name   = "volume-type"
    values = ["gp2"]
  }

  filter {
    name   = "tag:Name"
    values = ["Example"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `device` - The instance device name.
* `instance_id` - The ID of the instance you want to attach the volume to.
* `volume_id` - The ID of the volume you want to attach.

See detailed information in [Outscale Instances](https://wiki.outscale.net/display/DOCU/Getting+Information+About+Your+Instances).


## Attributes Reference

The following attributes are exported:

* `delete_on_termination` - Indicates whether the volume is deleted when terminating the instance
* `device` - The instance device name.
* `instance_id` -	The ID of the instance the volume is attached to.
* `status` - The attachment state of the volume (attaching | detaching | attached | detached).
* `volume_id` - The ID of the volume.


See detailed information in [Describe Instances](http://docs.outscale.com/api_fcu/definitions/VolumeAttachment.html#_api_fcu-volumeattachment).