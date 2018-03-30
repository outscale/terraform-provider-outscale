---
layout: "outscale"
page_title: "OUTSCALE: outscale_volume"
sidebar_current: "docs-outscale-datasource-volume"
description: |-
  Describes one or more specified Block Storage Unit (BSU) volume..
---

# outscale_volume

  Describes one or more specified Block Storage Unit (BSU) volume.

## Example Usage

```hcl
resource "outscale_volume" "external1" {
    availability_zone = "eu-west-2a"
    volume_type = "gp2"
    size = 10
    tags {
        Name = "External Volume 1"
    }
}
data "outscale_volumes" "ebs_volume" {
    filter {
	name = "size"
	values = ["${outscale_volume.external1.size}"]
    }
    filter {
	name = "volume-type"
	values = ["${outscale_volume.external1.volume_type}"]
    }
}
```

## Argument Reference

The following arguments are supported:

* `volume_id` - The ID of the volume.

See detailed information in [Outscale Instances](https://wiki.outscale.net/display/DOCU/Getting+Information+About+Your+Instances).

## Filters

Use the Filter.N parameter to filter the described volume on the following properties:

* `attachment.attach-time`: - The time at which the attachment was initiated.
* `attachment.delete-on-termination`: - Whether the volume is deleted when terminating the instance.
* `attachment.device`: - The device to which the volume is plugged in.
* `attachment.instance-id`: - The ID of the instance the volume is attached to.
* `attachment.status`: - The attachment state (attaching | attached | detaching | detached).
* `availability-zone`: - The Availability Zone in which the volume was created.
* `create-time`: - The time at which the volume was created.
* `tag`: - The key/value combination of a tag that is assigned to the resource, in the following format`: - key=value.
* `tag-key`: - The key of a tag associated with the resource.
* `tag-value`: - The value of a tag associated with the resource.
* `volume-id`: - The ID of the volume.
* `volume-type`: - The type of the volume (standard | gp2 | io1 | sc1| st1).
* `snapshot-id`: - The snapshot from which the volume was created.
* `size`: - The size of the volume, in Gibibytes (GiB).
* `status`: - The status of the volume (creating | available | in-use | deleting | error).


## Attributes Reference

The following attributes are exported:

* `attachment_set` - Information about your volume attachment.
* `availability_zone` - The Availability Zone where the volume is.
* `iops` - The number of I/O operations per second (only for io1 and gp2 volumes).
* `size` - The size of the volume, in Gibibytes (GiB).
* `snapshot_id` - The ID of the snapshot from which the volume was created.
* `status` - The state of the volume (creating| available| in-use| deleting| error).
* `tag_set` - One or more tags associated with the volume.
* `volume_id` - The ID of the volume.
* `volume_type` - The type of the volume (standard | gp2 | io1 | sc1 | st1).
* `request_id` - The ID of the request.

See detailed information in [Describe Instances]http://docs.outscale.com/api_fcu/operations/Action_DescribeVolumes_get.html#_api_fcu-action_describevolumes_get.