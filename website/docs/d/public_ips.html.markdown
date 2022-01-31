---
layout: "outscale"
page_title: "OUTSCALE: outscale_public_ips"
sidebar_current: "outscale-public-ips"
description: |-
  [Provides information about public IPs.]
---

# outscale_public_ips Data Source

Provides information about public IPs.
For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-EIPs.html).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-publicip).

## Example Usage

```hcl
data "outscale_public_ips" "public_ips01" {
  filter {
    name   = "public_ips"
    values = ["111.11.111.1", "222.22.222.2"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) A combination of a filter name and one or more filter values. You can specify this argument for as many filter names as you need. The filter name can be any of the following:
    * `link_public_ip_ids` - (Optional) The IDs representing the associations of public IPs with VMs or NICs.
    * `nic_account_ids` - (Optional) The account IDs of the owners of the NICs.
    * `nic_ids` - (Optional) The IDs of the NICs.
    * `placements` - (Optional) Whether the public IPs are for use in the public Cloud or in a Net.
    * `private_ips` - (Optional) The private IP addresses associated with the public IPs.
    * `public_ip_ids` - (Optional) The IDs of the public IPs.
    * `public_ips` - (Optional) The public IPs.
    * `tag_keys` - (Optional) The keys of the tags associated with the public IPs.
    * `tag_values` - (Optional) The values of the tags associated with the public IPs.
    * `tags` - (Optional) The key/value combination of the tags associated with the public IPs, in the following format: &quot;Filters&quot;:{&quot;Tags&quot;:[&quot;TAGKEY=TAGVALUE&quot;]}.
    * `vm_ids` - (Optional) The IDs of the VMs.

## Attribute Reference

The following attributes are exported:

* `public_ips` - Information about one or more public IPs.
    * `link_public_ip_id` - (Required in a Net) The ID representing the association of the public IP with the VM or the NIC.
    * `nic_account_id` - The account ID of the owner of the NIC.
    * `nic_id` - The ID of the NIC the public IP is associated with (if any).
    * `private_ip` - The private IP address associated with the public IP.
    * `public_ip` - The public IP.
    * `public_ip_id` - The allocation ID of the public IP.
    * `tags` - One or more tags associated with the public IP.
        * `key` - The key of the tag, with a minimum of 1 character.
        * `value` - The value of the tag, between 0 and 255 characters.
    * `vm_id` - The ID of the VM the public IP is associated with (if any).
