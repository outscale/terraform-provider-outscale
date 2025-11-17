---
layout: "outscale"
page_title: "OUTSCALE: outscale_api_access_rules"
subcategory: "Identity Access Management (IAM)"
sidebar_current: "outscale-api-access-rules"
description: |-
  [Provides information about API access rules.]
---

# outscale_api_access_rules Data Source

Provides information about API access rules.

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-API-Access-Rules.html).  
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-apiaccessrule).

## Example Usage

```hcl
data "outscale_api_access_rules" "api_access_rules01" {
    filter {
        name   = "ca_ids"
        values = ["ca-12345678", "ca-87654321"]
    }
    filter {
        name   = "ip_ranges"
        values = ["192.0.2.0/16"]
    }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) A combination of a filter name and one or more filter values. You can specify this argument for as many filter names as you need. The filter name can be any of the following:
    * `api_access_rule_ids` - (Optional) One or more IDs of API access rules.
    * `ca_ids` - (Optional) One or more IDs of Client Certificate Authorities (CAs).
    * `cns` - (Optional) One or more Client Certificate Common Names (CNs).
    * `descriptions` - (Optional) One or more descriptions of API access rules.
    * `ip_ranges` - (Optional) One or more IPs or CIDR blocks (for example, `192.0.2.0/16`).

## Attribute Reference

The following attributes are exported:

* `api_access_rules` - A list of API access rules.
    * `api_access_rule_id` - The ID of the API access rule.
    * `ca_ids` - One or more IDs of Client Certificate Authorities (CAs) used for the API access rule.
    * `cns` - One or more Client Certificate Common Names (CNs).
    * `description` - The description of the API access rule.
    * `ip_ranges` - One or more IP ranges used for the API access rule, in CIDR notation (for example, `192.0.2.0/16`).
