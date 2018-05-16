---
layout: "outscale"
page_title: "OUTSCALE: outscale_nics"
sidebar_current: "docs-outscale-datasource-nics"
description: |-
  Describes one or more Network Interfaces.

---

# outscale_nics

Describes one or more Network Interfaces.

## Example Usage

```hcl
data "outscale_nics" "nic" {
		network_interface_id = "NICID"
		subnet_id = "1"
}
```

## Argument Reference

The following arguments are supported:

* `nic_id` - (Optional)	One or more network interface IDs.

See detailed information in [Outscale Network Interfaces](https://wiki.outscale.net/display/DOCU/Getting+Information+About+Your+Instances).

## Filters

Use the Filter.N parameter to filter the described instances on the following properties:

* `addresses.private-ip-address`: The private IP addresses associated with the network interface.
* `addresses.primary`: Whether the private IP address is the primary IP address associated with the network interface.
* `addresses.association.public-ip`: The association ID returned when the network interface was associated with the External IP address.
* `addresses.association.owner-id`: The account ID of the owner of the addresses associated with the network interface.
* `association.association-id`: The association ID returned when the network interface was associated with an IP address.
* `association.allocation-id`: The allocation ID returned when you allocated the External IP address for your network interface.
* `association.ip-owner-id`: The account ID of the owner of the External IP address associated with the network interface.
* `association.public-ip`: The External IP address associated with the network interface.
* `attachment.attachment-id`: The ID of the interface attachment.
* `attachment.instance-id`: The ID of the instance the network interface is attached to.
* `attachment.instance-owner-id`: The account ID of the owner of the instance the network interface is attached to.
* `attachment.device-index`: The device index the network interface is attached to.
* `attachment.status`: The state of the attachment (attaching | attached | detaching | detached).
* `attachment.attach.time`: The time at which the network interface was attached to an instance.
* `attachment.delete-on-termination`: Indicates whether the attachment is deleted when the instance is terminated.
* `availability-zone`: The Availability Zone of the network interface.
* `description`: The description of the network interface.
* `group-id`: The ID of a security group associated with the network interface.
* `group-name`: The name of a security group associated with the network interface.
* `mac-address`: The MAC address of the network interface.
* `network-interface-id`: The ID of the network interface.
* `owner-id`: The Outscale account ID of the network interface owner.
* `private-ip-address`: The private IP address or addresses of the network interface.
* `private-dns-name`: The private DNS name of the network interface.
* `requester-id`: The ID of the entity that launched the instance.
* `requester-managed`: Indicates whether the network interface is managed by an Outscale service.
* `source-dest-check`: If true, the source/destination checking is enabled. If false, it is disabled. This value must be false for a NAT instance to perform NAT (network address translation) in a VPC.
* `status`: The state of the network interface.
* `subnet-id`: The ID of the subnet for the network interface.
* `vpc-id`: The ID of the VPC for the network interface.


## Attributes Reference

The following attributes are exported:

* `association` - 	The association information for an External IP associated with the network interface.	false	NetworkInterfaceAssociation


* `nic`	Information about the network interfaces.	
* `lin_id`	The ID of the VPC.

See detailed information in [Describe Outscale Network Interfaces](http://docs.outscale.com/api_fcu/operations/Action_DescribeNetworkInterfaces_get.html#_api_fcu-action_describenetworkinterfaces_get).
