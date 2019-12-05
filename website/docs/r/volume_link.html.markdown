---
layout: "outscale"
page_title: "3DS OUTSCALE: outscale_volume_link"
sidebar_current: "docs-outscale-resource-volume-link"
description: |-
  [Manages a volume link.]
---

# outscale_volume_link Resource

Manages a volume link.
For more information on this resource, see the [User Guide](https://wiki.outscale.net/display/EN/About+Volumes).
For more information on this resource actions, see the [API documentation](https://docs-beta.outscale.com/#linkvolume).

## Example Usage

```hcl
[exemple de code]
```

## Argument Reference

The following arguments are supported:

* `device_name` - (Required) The name of the device.
* `vm_id` - (Required) The ID of the VM you want to attach the volume to.
* `volume_id` - (Required) The ID of the volume you want to attach.

## Attribute Reference

The following attributes are exported:

* `delete_on_vm_deletion` - Indicates whether the volume is deleted when terminating the instance
* `device_name` - The name of the device.
* `vm_id` - The ID of the VM you want to attach the volume to.
* `state` - The attachment state of the volume (`attaching` | `detaching` | `attached` | `detached`).
* `volume_id` - The ID of the volume you want to attach.
