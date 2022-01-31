---
layout: "outscale"
page_title: "OUTSCALE: outscale_dhcp_option"
sidebar_current: "outscale-dhcp-option"
description: |-
  [Manages a DHCP option.]
---

# outscale_dhcp_option Resource

Manages a DHCP option.
For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-DHCP-Options.html).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-dhcpoption).

## Example Usage

### Create a basic DHCP options set

```hcl
resource "outscale_dhcp_option" "dhcp_option_01" {
	domain_name = "MyCompany.com"
}
```

### Create a complete DHCP options set

```hcl
resource "outscale_dhcp_option" "dhcp_option_02" {
	domain_name         = "MyCompany.com"
	domain_name_servers = ["111.111.11.111","222.222.22.222"]
	ntp_servers         = ["111.1.1.1","222.2.2.2"]
	tags {
		key = "Name"
		value = "DHCP01"
	}
}
```

## Argument Reference

The following arguments are supported:

* `domain_name_servers` - (Optional) The IP addresses of domain name servers. If no IP addresses are specified, the `OutscaleProvidedDNS` value is set by default.
* `domain_name` - (Optional) Specify a domain name (for example, MyCompany.com). You can specify only one domain name.
* `ntp_servers` - (Optional) The IP addresses of the Network Time Protocol (NTP) servers.
* `tags` - (Optional) A tag to add to this resource. You can specify this argument several times.
    * `key` - (Required) The key of the tag, with a minimum of 1 character.
    * `value` - (Required) The value of the tag, between 0 and 255 characters.

## Attribute Reference

The following attributes are exported:

* `default` - If true, the DHCP options set is a default one. If false, it is not.
* `dhcp_options_set_id` - The ID of the DHCP options set.
* `domain_name_servers` - One or more IP addresses for the domain name servers.
* `domain_name` - The domain name.
* `ntp_servers` - One or more IP addresses for the NTP servers.
* `tags` - One or more tags associated with the DHCP options set.
    * `key` - The key of the tag, with a minimum of 1 character.
    * `value` - The value of the tag, between 0 and 255 characters.

## Import

DHCP options can be imported using the DHCP option ID. For example:

```console

$ terraform import outscale_dhcp_option.ImportedDhcpSet dopt-87654321

```