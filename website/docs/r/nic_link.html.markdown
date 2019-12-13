---
layout: "outscale"
page_title: "3DS OUTSCALE: outscale_nic_link"
sidebar_current: "outscale-nic-link"
description: |-
  [Manages a NIC link.]
---

# outscale_nic_link Resource

Manages a NIC link.
For more information on this resource, see the [User Guide](https://wiki.outscale.net/display/EN/About+FNIs#AboutFNIs-FNIAttachmentFNIsAttachmenttoInstances).
For more information on this resource actions, see the [API documentation](https://docs-beta.outscale.com/#linknic).

## Example Usage

```hcl
[exemple de code]
```

## Argument Reference

The following arguments are supported:

* `device_number` - (Required) The index of the VM device for the NIC attachment (between 1 and 7, both included).
* `nic_id` - (Required) The ID of the NIC you want to attach.
* `vm_id` - (Required) The ID of the VM to which you want to attach the NIC.

## Attribute Reference

The following attributes are exported:

* `link_nic_id` - The ID of the NIC attachment.