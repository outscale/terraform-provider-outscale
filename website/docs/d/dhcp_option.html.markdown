---
layout: "outscale"
page_title: "OUTSCALE: outscale_dhcp_option"
sidebar_current: "outscale-dhcp-option"
description: |-
  [Provides information about a specific DHCP option.]
---

# outscale_dhcp_option Data Source

Provides information about a specific DHCP option.
For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-DHCP-Options.html).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-dhcpoption).

## Example Usage

```hcl
data "outscale_dhcp_option" "data_dhcp_option" {
	filter {
		name   = "dhcp_options_set_id"
		values = ["dopt-12345678"]
	}
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) A combination of a filter name and one or more filter values. You can specify this argument for as many filter names as you need. The filter name can be any of the following:
    * `default` - (Optional) If true, lists all default DHCP options set. If false, lists all non-default DHCP options set.
    * `dhcp_options_set_ids` - (Optional) The IDs of the DHCP options sets.
    * `domain_name_servers` - (Optional) The domain name servers used for the DHCP options sets.
    * `domain_names` - (Optional) The domain names used for the DHCP options sets.
    * `ntp_servers` - (Optional) The Network Time Protocol (NTP) servers used for the DHCP options sets.
    * `tag_keys` - (Optional) The keys of the tags associated with the DHCP options sets.
    * `tag_values` - (Optional) The values of the tags associated with the DHCP options sets.
    * `tags` - (Optional) The key/value combination of the tags associated with the DHCP options sets, in the following format: &quot;Filters&quot;:{&quot;Tags&quot;:[&quot;TAGKEY=TAGVALUE&quot;]}.

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
