---
layout: "outscale"
page_title: "OUTSCALE: outscale_public_ip"
sidebar_current: "docs-outscale-resource-public_ip"
description: |-
  Provides an Outscale Public IP Association as a top level resource, to associate and disassociate Public IPs from Outscale VMs and Network Interfaces.
---

# outscale_public_ip

NOTE: outscale_public_ip is useful in scenarios where Public IPs are either pre-existing or distributed to customers or users and therefore cannot be changed.

## Example Usage

```hcl
resource "outscale_public_ip" "bar" {}
```

## Argument Reference

The following arguments are supported:

* `placement` - (Optional) Indicates whether the External IP address is for use with instances in the public Cloud or in a VPC.

## Attribute reference

* `reservation_id` - The ID of the allocation.
* `link_id` - The association ID for the EIP.
* `placement` - Whether the EIP is for use in the public Cloud or in a VPC.
* `nic_id` - The ID of the network interface the address is associated with (if any).
* `nic_account_id` - The account ID of the owner.
* `private_ip` - The private IP address associated with the EIP.
* `public_ip` - The EIP.
* `request_id` - The ID of the request.