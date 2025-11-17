---
layout: "outscale"
page_title: "OUTSCALE: outscale_ca"
subcategory: "OUTSCALE API"
sidebar_current: "outscale-ca"
description: |-
  [Manages a Certificate Authority (CA).]
---

# outscale_ca Resource

Manages a Certificate Authority (CA).

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-API-Access-Rules.html).  
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-ca).

## Example Usage

```hcl
resource "outscale_ca" "ca01" {
    ca_pem      = file("<PATH>")
    description = "Terraform certificate authority"
}
```

## Argument Reference

The following arguments are supported:

* `ca_pem` - (Required) The CA in PEM format.
* `description` - (Optional) The description of the CA.

## Attribute Reference

The following attributes are exported:

* `ca_fingerprint` - The fingerprint of the CA.
* `ca_id` - The ID of the CA.
* `description` - The description of the CA.

## Import

A CA can be imported using its ID. For example:

```console

$ terraform import outscale_ca.ImportedCa ca-12345678

```