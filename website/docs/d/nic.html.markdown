---
layout: "outscale"
page_title: "OUTSCALE: outscale_nic"
sidebar_current: "docs-outscale-datasource-nic"
description: |-
  Describes one or more Network Interfaces through Terraform.

---

# outscale_lin

Describes one or more Network Interfaces through Terraform.

## Example Usage

```hcl
data "outscale_nic" "nic" {
		network_interface_id = "NICID"
		subnet_id = "1"
}
```

## Argument Reference

The following arguments are supported:

* `network_interface_id` - (Optional)	One or more network interface IDs.

See detailed information in [Outscale Nic](http://docs.outscale.com/api_fcu/operations/Action_DescribeNetworkInterfaces_get.html#_api_fcu-action_describenetworkinterfaces_get).

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
* `attachment` - 	The network interface attachment.	false	NetworkInterfaceAttachment
* `availability_zone` - 	The Availability Zone in which the network interface is located.	false	string
* `description` - 	A description of the network interface.	false	string
* `group_set` - 	One or more security groups for the network interface.	false	GroupIdentifier
* `mac_address` - 	The MAC address.	false	string
* `network_interface_id` - 	The ID of the network interface.	false	string
* `owner_id` - 	The account ID of the owner of the network interface.	false	string
* `private_dns_name` - 	The name of the private DNS.	false	string
* `private_ip_address` - 	The private IP addresses assigned to the network interface, in the CIDR of its subnet.	false	string
* `private_ip_addresses_set` - 	Information about one or more private IP addresses assigned to the network interface.	false	NetworkInterfacePrivateIpAddress
* `requester_id` - 	The ID of the requester that launched the instances on your behalf.	false	string
* `requester_managed` - 	If true, the network interface is being managed by Outscale.	false	boolean
* `source_dest_check` - 	If true, the traffic to or from the instance is validated.	false	boolean
* `status` - 	The state of the network interface (available | attaching | in-use | detaching).	false	string
* `subnet_id` - 	The ID of the subnet.	false	string
* `tag_set` - 	One or more tags associated with the network interface.	false	Tag
* `vpc_idd` - 	The ID of the VPC.	false	string

See detailed information in [Describe Outscale Nic](http://docs.outscale.com/api_fcu/operations/Action_DescribeNetworkInterfaces_get.html#_api_fcu-action_describenetworkinterfaces_get).
