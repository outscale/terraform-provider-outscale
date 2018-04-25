---
layout: "outscale"
page_title: "OUTSCALE: outscale_client_endpoint"
sidebar_current: "docs-outscale-datasource-client-endpoint"
description: |-
Describes client endpoint


---

# outscale_client_endpoint

Describes client endpoint


## Example Usage

```hcl
data "outscale_firewall_rules_sets" "by_filter" {
	filter {
		name = "group-name"
		values = ["${outscale_firewall_rules_sets.outscale_firewall_rules_sets.group_name}"]
	}
}`, rInt, rInt)
```

## Argument Reference

The following arguments are supported:

* `customer_gateway_id.N` (Optional)One or more customer gateways IDs.
* `Fulter.N` (Optional)One or more customer gateways IDs.



## Filters

Use the Filter.N parameter to filter the described instances on the following properties:

* `state`  The state of the customer gateway (pending | available | deleting | deleted).
* `customer-gateway-id`  The ID of the customer gateway.
* `ip-address	`  The public IPv4 address of the customer gateway.
* `ip-permission.cidr	` The ASN number.
* `bgp-asn`  	The type of communication tunnel to the gateway.
* `type`  The type of communication tunnel to the gateway.
* `tag-key`  The key of a tag assigned to the resource. This filter is independent of the tag-value filter.
* `tag-value`  The value of a tag assigned to the resource. This filter is independent of the tag-key filter.



## Attributes Reference

The following attributes are exported:

* `bgp_asn` - An unsigned 32-bits ASN (Autonomous System Number) used by the BGP (Border Gateway Protocol) to find out the path to the customer gateway through the Internet.
* `customer_gateway_id` - The ID of the customer gateway.
* `ip_address` - The public IPv4 address of the customer gateway (must be a fixed address into a NATed network).
* `state` - The state of the customer gateway (pending | available | deleting | deleted).
* `tag_set.N` - One or more tags associated with the customer gateway.
* `type` - The type of communication tunnel used by the customer gateway (only ipsec.1 is supported).
