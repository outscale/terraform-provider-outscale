---
layout: "outscale"
page_title: "OUTSCALE: outscale_dhcp_options"
sidebar_current: "outscale-dhcp-options"
description: |-
  [Provides information about DHCP options.]
---

# outscale_dhcp_options Data Source

Provides information about DHCP options.
For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-DHCP-Options.html).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-dhcpoption).

## Example Usage

```hcl
data "outscale_dhcp_option" "data_dhcp_options" {
	filter {
		name   = "dhcp_options_set_ids"
		values = ["dopt-12345678","dopt-87654321"]
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

* `dhcp_options_sets` - Information about one or more DHCP options sets.
    * `default` - If true, the DHCP options set is a default one. If false, it is not.
    * `dhcp_options_set_id` - The ID of the DHCP options set.
    * `domain_name` - The domain name.
    * `domain_name_servers` - One or more IP addresses for the domain name servers.
    * `ntp_servers` - One or more IP addresses for the NTP servers.
    * `tags` - One or more tags associated with the DHCP options set.
        * `key` - The key of the tag, with a minimum of 1 character.
        * `value` - The value of the tag, between 0 and 255 characters.
