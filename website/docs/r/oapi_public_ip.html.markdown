---
layout: "outscale"
page_title: "OUTSCALE: outscale_public_ip"
sidebar_current: "docs-outscale-resource-public_ip"
description: |-
  Acquires an External IP address (EIP) for your account.
  An EIP is a static IP address designed for dynamic Cloud computing. It can be used for virtual machines (VMs) in the public Cloud (standard) or in a Net, for a network interface card (NIC), or for a NAT service.
---

# outscale_public_ip

Acquires an External IP address (EIP) for your account.
An EIP is a static IP address designed for dynamic Cloud computing. It can be used for virtual machines (VMs) in the public Cloud (standard) or in a Net, for a network interface card (NIC), or for a NAT service.

## Example Usage

```hcl
resource "outscale_public_ip" "bar" {}
```

## Argument Reference

No arguments are supported.

## Attributes Reference

* `link_public_ip_id` - The ID representing the association of the EIP with the VM or the NIC.
* `nic_account_id` - The account ID of the owner of the NIC.
* `nic_id` - The ID of the NIC the EIP is associated with (if any).
* `private_ip` - The private IP address associated with the EIP.
* `public_ip` - The External IP address (EIP) associated with the NAT service.
* `public_ip_id` - The allocation ID of the EIP associated with the NAT service.
* `vm_id` - The ID of the VM the External IP (EIP) is associated with (if any).
* `request_id` - The ID of the request.

See detailed information in [OAPI PublicIp](http://docs.outscale.com/oapi/index.html#tocspublicip).
