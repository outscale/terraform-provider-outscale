---
layout: "outscale"
page_title: "OUTSCALE: outscale_cas"
subcategory: "Identity Access Management (IAM)"
sidebar_current: "outscale-cas"
description: |-
  [Provides information about Certificate Authorities (CAs).]
---

# outscale_cas Data Source

Provides information about Certificate Authorities (CAs).

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-API-Access-Rules.html).  
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-ca).

## Example Usage

```hcl
data "outscale_cas" "cas01" {
    filter {
        name   = "ca_ids"
        values = ["ca-12345678", "ca-87654321"]
    }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) A combination of a filter name and one or more filter values. You can specify this argument for as many filter names as you need. The filter name can be any of the following:
    * `ca_fingerprints` - (Optional) The fingerprints of the CAs.
    * `ca_ids` - (Optional) The IDs of the CAs.
    * `descriptions` - (Optional) The descriptions of the CAs.

## Attribute Reference

The following attributes are exported:

* `cas` - Information about one or more CAs.
    * `ca_fingerprint` - The fingerprint of the CA.
    * `ca_id` - The ID of the CA.
    * `description` - The description of the CA.
