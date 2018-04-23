---
layout: "outscale"
page_title: "OUTSCALE: outscale_client_endpoint"
sidebar_current: "docs-outscale-resource-client-endpoint"
description: |-
  Creates an OMI from an existing instance which is either running or stopped.
---

# outscale_client_endpoint

Provides information about your customer gateway.

## Example Usage

```hcl
		resource "outscale_client_endpoint" "foo" {
			bgp_asn = %d
			ip_address = "172.10.10.1"
			type = "ipsec.1"
			tag {
				Name = "foo-gateway-%d"
				Another = "tag"
			}
		}
```

## Argument Reference

The following arguments are supported:


* `bgp_asn` - (Required)	An unsigned 32-bits Autonomous System Number (ASN) used by the Border Gateway Protocol (BGP) to find out the path to your customer gateway through the Internet network. The integer must be within the [0;4294967295] range.
* `ip_address` -	(Required) The public fixed IPv4 address of your customer gateway.	
* `type` - (Required)	The communication protocol used to establish tunnel with your customer gateway (only ipsec.1 is supported).	


## Attributes

* `bgp_asn`	An unsigned 32-bits ASN (Autonomous System Number) used by the BGP (Border Gateway Protocol) to find out the path to the customer gateway through the Internet.	false	integer
* `customer_gateway_id`	The ID of the customer gateway.	false	string
* `ip_address`	The public IPv4 address of the customer gateway (must be a fixed address into a NATed network).	false	string
* `state`	The state of the customer gateway (pending | available | deleting | deleted).	false	string
* `tag_set`	One or more tags associated with the customer gateway.	false	Tag
* `type`	The type of communication tunnel used by the customer gateway (only ipsec.1 is supported).	false	string
* `request_id`	The ID of the request


See detailed information in [Describe Client Endpoint](http://docs.outscale.com/api_fcu/operations/Action_DescribeCustomerGateways_get.html#_api_fcu-action_describecustomergateways_get).

