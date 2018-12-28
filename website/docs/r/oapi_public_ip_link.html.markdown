---
layout: "outscale"
page_title: "OUTSCALE: outscale_public_ip_link"
sidebar_current: "docs-outscale-resource-public_ip_link"
description: |-
  Associates an External IP address (EIP) with a virtual machine (VM) or a network interface card (NIC).
---

# outscale_public_ip_link

Associates an External IP address (EIP) with a virtual machine (VM) or a network interface card (NIC), in the public Cloud or in a Net. You can associate an EIP with only one VM or network interface at a time.

NOTE: You can associate an EIP with a network address translation (NAT) service only when creating the NAT service. To modify its EIP, you need to delete the NAT service and re-create it with the new EIP. For more information, see the CreateNatService method.

## Example Usage

```hcl
resource "outscale_vm" "basic" {
  image_id = "ami-8a6a0120"
  instance_type = "t2.micro"
  key_name = "terraform-basic"
  subnet_id = "subnet-861fbecc"
}

resource "outscale_public_ip" "bar" {}

resource "outscale_public_ip_link" "by_public_ip" {
  public_ip = "${outscale_public_ip.bar.public_ip}"
  vm_id = "${outscale_vm.basic.id}"
}
```

## Argument Reference

The following arguments are supported:

* `allow_relink` - If true, allows the EIP to be associated with the VM or NIC that you specify even if it is already associated with another VM or NIC. If false, prevents the EIP from being associated with the VM or NIC that you specify if it is already associated with another VM or NIC.
* `nic_id` - (Optional) (Net only) The ID of the NIC. This parameter is required if the VM has more than one NIC attached. Otherwise, you need to specify the VmId parameter instead. You cannot specify both parameters at the same time.
* `private_ip` - (Optional) (Net only) The primary or secondary private IP address of the specified NIC. By default, the primary private IP address.
* `public_ip_id` - (Optional) The allocation ID of the EIP. In a Net, this parameter is required.
* `vm_id` - (Optional) The ID of the VM.
  * In the public Cloud, this parameter is required.
  * In a Net, this parameter is required if the VM has only one NIC. Otherwise, you need to specify the NicId parameter instead. You cannot specify both parameters at the same time.
* `public_ip` - (Optional) The External IP address.

## Attributes Reference

* `link_public_ip_id` - The ID representing the association of the EIP with the VM or the NIC.
* `nic_account_id` - The account ID of the owner of the NIC.
* `nic_id` - The ID of the NIC the EIP is associated with (if any).
* `private_ip` - The private IP address associated with the EIP.
* `public_ip` - The External IP address (EIP) associated with the NAT service.
* `public_ip_id` - The allocation ID of the EIP associated with the NAT service.
* `vm_id` - The ID of the VM the External IP (EIP) is associated with (if any).
* `request_id` - The ID of the request.
