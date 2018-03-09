---
layout: "outscale"
page_title: "OUTSCALE: outscale_keypair"
sidebar_current: "docs-outscale-resource-keypair"
description: |-
  Provides an Outscale keypair resource. This allows keypairs to be created, deleted, described and imported.
---

# outscale_vm

Provides an Outscale keypair resource. This allows keypairs to be created, deleted,
described and imported. Instances also support [provisioning](/docs/provisioners/index.html).

## Example Usage

```hcl
resource "outscale_keypair" "a_key_pair" {
	key_name   = "tf-acc-key-pair"
}
```

## Argument Reference

The following arguments are supported:

* `key_name` - (Required) A unique name for the keypair, with a maximum length of 255 ASCII characters.

See detailed information in [Outscale Instances](https://wiki.outscale.net/display/DOCU/Getting+Information+About+Your+Instances).


## Attributes Reference

The following attributes are exported:

* `key_fingerprint` - The SHA-1 digest of the DER encoded private key.
* `key_material` - The private key.
* `key_name` - A unique name for the keypair.

See detailed information in [Describe Instances](http://docs.outscale.com/api_fcu/definitions/KeyPairInfo.html#_api_fcu-keypairinfo).
