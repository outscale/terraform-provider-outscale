---
layout: "outscale"
page_title: "OUTSCALE: outscale_ca"
subcategory: "CA (Client Certificate Authority)"
sidebar_current: "outscale-ca"
description: |-
  [Provides information about a Certificate Authority (CA).]
---

# outscale_ca Data Source

Provides information about a Certificate Authority (CA).

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-API-Access-Rules.html).  
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-ca).

## Example Usage

```hcl
data "outscale_ca" "ca01" { 
    filter {
        name   = "ca_ids"
        values = ["ca-12345678"]
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

* `ca_fingerprint` - The fingerprint of the CA.
* `ca_id` - The ID of the CA.
* `description` - The description of the CA.
