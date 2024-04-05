---
layout: "outscale"
page_title: "OUTSCALE: outscale_keypair"
sidebar_current: "outscale-keypair"
description: |-
  [Provides information about a keypair.]
---

# outscale_keypair Data Source

Provides information about a keypair.

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Keypairs.html).  
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-keypair).

## Example Usage

```hcl
data "outscale_keypair" "keypair01" {
	filter {
		name   = "keypair_names"
		values = ["terraform-keypair-01"]
	}
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) A combination of a filter name and one or more filter values. You can specify this argument for as many filter names as you need. The filter name can be any of the following:
    * `keypair_fingerprints` - (Optional) The fingerprints of the keypairs.
    * `keypair_names` - (Optional) The names of the keypairs.
    * `keypair_types` - (Optional) The types of the keypairs (`ssh-rsa`, `ssh-ed25519`, `ecdsa-sha2-nistp256`, `ecdsa-sha2-nistp384`, or `ecdsa-sha2-nistp521`).

## Attribute Reference

The following attributes are exported:

* `keypair_fingerprint` - The MD5 public key fingerprint as specified in section 4 of RFC 4716.
* `keypair_name` - The name of the keypair.
* `keypair_type` - The type of the keypair (`ssh-rsa`, `ssh-ed25519`, `ecdsa-sha2-nistp256`, `ecdsa-sha2-nistp384`, or `ecdsa-sha2-nistp521`).
