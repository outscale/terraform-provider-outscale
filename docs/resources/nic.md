---
layout: "outscale"
page_title: "OUTSCALE: outscale_nic"
subcategory: "OUTSCALE API"
sidebar_current: "outscale-nic"
description: |-
  [Manages a network interface card (NIC).]
---

# outscale_nic Resource

Manages a network interface card (NIC).

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-NICs.html).  
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-nic).

## Example Usage

### Required resources

```hcl
resource "outscale_net" "net01" {
    ip_range = "10.0.0.0/16"
}

resource "outscale_subnet" "subnet01" {
    subregion_name = "eu-west-2a"
    ip_range       = "10.0.0.0/18"
    net_id         = outscale_net.net01.net_id
}

resource "outscale_security_group" "security_group01" {
    description         = "Terraform security group for nic with private IPs"
    security_group_name = "terraform-security-group-nic-ips"
    net_id              = outscale_net.net01.net_id
}
```

### Create a NIC

```hcl
resource "outscale_nic" "nic01" {
    subnet_id = outscale_subnet.subnet01.subnet_id
    security_group_ids = [outscale_security_group.security_group01.security_group_id]
}

```

### Create a NIC with private IP addresses

```hcl
resource "outscale_nic" "nic02" {
    description       = "Terraform nic with private IPs"
    subnet_id         = outscale_subnet.subnet01.subnet_id
    security_group_ids = [outscale_security_group.security_group01.security_group_id]
    private_ips {
        is_primary = true
        private_ip = "10.0.0.4"
    }
    private_ips {
        is_primary = false
        private_ip = "10.0.0.5"
    }
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) A description for the NIC.
* `private_ips` - (Optional) Information about the private IP or IPs of the NIC.
    * `is_primary` - (Optional) If true, the IP is the primary private IP of the NIC.
    * `private_ip` - (Optional) A private IP for the NIC. This IP must be within the IP range of the Subnet that you specify with the `subnet_id` attribute. However, it cannot be one of the first four IPs (ending in `.0`, `.1`, `.2`, `.3`) or the last IP (ending in `.255`) of the Subnet, as these are reserved by 3DS OUTSCALE. For more information, see [About Nets](https://docs.outscale.com/en/userguide/About-Nets.html).<br />
If you do not specify this argument, a random private IP is selected within the IP range of the Subnet.
* `security_group_ids` - (Optional) One or more IDs of security groups for the NIC.
* `subnet_id` - (Required) The ID of the Subnet in which you want to create the NIC.
* `tags` - (Optional) A tag to add to this resource. You can specify this argument several times.
    * `key` - (Required) The key of the tag, with a minimum of 1 character.
    * `value` - (Required) The value of the tag, between 0 and 255 characters.

## Attribute Reference

The following attributes are exported:

* `account_id` - The account ID of the owner of the NIC.
* `description` - The description of the NIC.
* `is_source_dest_checked` - (Net only) If true, the source/destination check is enabled. If false, it is disabled.
* `link_nic` - Information about the NIC attachment.
    * `delete_on_vm_deletion` - If true, the NIC is deleted when the VM is terminated.
    * `device_number` - The device index for the NIC attachment (between `1` and `7`, both included).
    * `link_nic_id` - The ID of the NIC to attach.
    * `state` - The state of the attachment (`attaching` \| `attached` \| `detaching` \| `detached`).
    * `vm_account_id` - The account ID of the owner of the VM.
    * `vm_id` - The ID of the VM.
* `link_public_ip` - Information about the public IP association.
    * `link_public_ip_id` - (Required in a Net) The ID representing the association of the public IP with the VM or the NIC.
    * `public_dns_name` - The name of the public DNS.
    * `public_ip` - The public IP associated with the NIC.
    * `public_ip_account_id` - The account ID of the owner of the public IP.
    * `public_ip_id` - The allocation ID of the public IP.
* `mac_address` - The Media Access Control (MAC) address of the NIC.
* `net_id` - The ID of the Net for the NIC.
* `nic_id` - The ID of the NIC.
* `private_dns_name` - The name of the private DNS.
* `private_ips` - The private IPs of the NIC.
    * `is_primary` - If true, the IP is the primary private IP of the NIC.
    * `link_public_ip` - Information about the public IP association.
        * `link_public_ip_id` - (Required in a Net) The ID representing the association of the public IP with the VM or the NIC.
        * `public_dns_name` - The name of the public DNS.
        * `public_ip` - The public IP associated with the NIC.
        * `public_ip_account_id` - The account ID of the owner of the public IP.
        * `public_ip_id` - The allocation ID of the public IP.
    * `private_dns_name` - The name of the private DNS.
    * `private_ip` - The private IP of the NIC.
* `security_groups` - One or more IDs of security groups for the NIC.
    * `security_group_id` - The ID of the security group.
    * `security_group_name` - The name of the security group.
* `state` - The state of the NIC (`available` \| `attaching` \| `in-use` \| `detaching`).
* `subnet_id` - The ID of the Subnet.
* `subregion_name` - The Subregion in which the NIC is located.
* `tags` - One or more tags associated with the NIC.
    * `key` - The key of the tag, with a minimum of 1 character.
    * `value` - The value of the tag, between 0 and 255 characters.

## Timeouts

The `timeouts` block enables you to configure [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - Defaults to 10 minutes.
* `read` - Defaults to 5 minutes.
* `update` - Defaults to 10 minutes.
* `delete` - Defaults to 5 minutes.

## Import

A NIC can be imported using its ID. For example:

```console

$ terraform import outscale_nic.ImportedNic eni-12345678

```