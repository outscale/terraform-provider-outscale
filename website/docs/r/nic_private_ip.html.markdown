---
layout: "outscale"
page_title: "3DS OUTSCALE: outscale_nic_private_ip"
sidebar_current: "outscale-nic-private-ip"
description: |-
  [Manages a NIC private IP.]
---

# outscale_nic_private_ip Resource

Manages a NIC private IP.
For more information on this resource, see the [User Guide](https://wiki.outscale.net/display/EN/About+FNIs#AboutFNIs-FNIsAttributes).
For more information on this resource actions, see the [API documentation](https://docs-beta.outscale.com/#linkprivateips).

## Example Usage

```hcl
[exemple de code]
```

## Argument Reference

The following arguments are supported:

* `allow_relink` - (Optional) If `true`, allows an IP address that is already assigned to another NIC in the same Subnet to be assigned to the NIC you specified.
* `nic_id` - (Required) The ID of the NIC.
* `private_ips` - (Optional) The secondary private IP address or addresses you want to assign to the NIC within the IP address range of the Subnet.
* `secondary_private_ip_count` - (Optional) The number of secondary private IP addresses to assign to the NIC.

## Attribute Reference

No attribute is exported.