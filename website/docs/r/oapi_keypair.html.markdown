---
layout: "outscale"
page_title: "OUTSCALE: outscale_keypair"
sidebar_current: "docs-outscale-resource-keypair"
description: |-
  Creates a 2048-bit RSA keypair with a specified name. This action returns the private key that you need to save. The public key is stored by Outscale.
---

# outscale_keypair

Creates a 2048-bit RSA keypair with a specified name.
This action returns the private key that you need to save. The public key is stored by Outscale.

You can also use this method to import a provided public key and create a keypair.

This action imports the public key of a keypair created by a third-party tool and uses it to create a new keypair. The private key is never provided to Outscale.

## Example Usage

```hcl
resource "outscale_keypair" "a_key_pair" {
 keypair_name   = "tf-acc-key-pair"
}
```

## Argument Reference

The following arguments are supported:

* `keypair_name` - (Required) A unique name for the keypair, with a maximum length of 255 [ASCII printable characters](https://en.wikipedia.org/wiki/ASCII#Printable_characters).
* `public_key` - (Optional) The public key. If you are not using command line tools, it must be encoded in Base64.

See detailed information in [Outscale CreateKeypair](http://docs.outscale.com/api_fcu/operations/Action_CreateKeyPair_get.html#_api_fcu-action_createkeypair_get).

## Attributes Reference

The following attributes are exported:

* `keypair_fingerprint` - If you created the keypair, the SHA-1 digest of the DER encoded private key. If you imported the keypair, the MD5 public key fingerprint as specified in section 4 of RFC4716.
* `private_key` - The private key.
* `keypair_name` - The name of the keypair.
* `request_id` - The ID of the request.

See detailed information in [ReadKeyPairs](http://docs.outscale.com/api_fcu/definitions/KeyPairInfo.html#_api_fcu-keypairinfo).
