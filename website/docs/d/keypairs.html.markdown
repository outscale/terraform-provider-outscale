---
layout: "outscale"
page_title: "3DS OUTSCALE: outscale_keypairs"
sidebar_current: "outscale-keypairs"
description: |-
  [Provides information about keypairs.]
---

# outscale_keypairs Data Source

Provides information about keypairs.
For more information on this resource, see the [User Guide](https://wiki.outscale.net/display/EN/About+Keypairs).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-keypair).

## Example Usage

```hcl

data "outscale_keypairs" "keypairs01" {
	filter {
		name   = "keypair_names"
		values = ["terraform-keypair-01", "terraform-keypair-02"]
	}
}


```

## Argument Reference

The following arguments are supported:

* `filter` - One or more filters.
  * `keypair_fingerprints` - (Optional) The fingerprints of the keypairs.
  * `keypair_names` - (Optional) The names of the keypairs.

## Attribute Reference

The following attributes are exported:

* `keypairs` - Information about one or more keypairs.
  * `keypair_fingerprint` - If you create a keypair, the SHA-1 digest of the DER encoded private key.<br />
If you import a keypair, the MD5 public key fingerprint as specified in section 4 of RFC 4716.
  * `keypair_name` - The name of the keypair.
