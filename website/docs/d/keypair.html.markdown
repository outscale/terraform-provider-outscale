---
layout: "outscale"
page_title: "3DS OUTSCALE: outscale_keypair"
sidebar_current: "docs-outscale-datasource-keypair"
description: |-
  [Provides information about a specific keypair.]
---

# outscale_keypair Data Source

Provides information about a specific keypair.
For more information on this resource, see the [User Guide](https://wiki.outscale.net/display/EN/About+Keypairs).
For more information on this resource actions, see the [API documentation](https://docs-beta.outscale.com/#3ds-outscale-api-keypair).

## Example Usage

```hcl
[exemple de code]
```

## Argument Reference

The following arguments are supported:

* `filters` - One or more filters.
  * `keypair_fingerprints` - (Optional) The fingerprints of the keypairs.
  * `keypair_names` - (Optional) The names of the keypairs.

## Attribute Reference

The following attributes are exported:

* `keypairs` - Information about one or more keypairs.
  * `keypair_fingerprint` - If you create a keypair, the SHA-1 digest of the DER encoded private key.<br />
If you import a keypair, the MD5 public key fingerprint as specified in section 4 of RFC 4716.
  * `keypair_name` - The name of the keypair.
