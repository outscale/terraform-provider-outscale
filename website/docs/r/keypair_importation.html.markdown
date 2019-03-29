---
layout: "outscale"
page_title: "OUTSCALE: outscale_keypair_importation"
sidebar_current: "docs-outscale-resource-keypair-importation"
description: |-
  Imports a provided public key and creates a keypair.
---

# outscale_keypair

Imports a provided public key and creates a keypair. [provisioning](/docs/provisioners/index.html).

## Example Usage

```hcl
resource "outscale_keypair_importation" "outscale_keypair_importation" {

    key_name            = "keyname_test"
    public_key_material = "${file("keypair_public_test")}"
}
```

## Argument Reference

The following arguments are supported:

* `keypair_name` - (Required)	A unique name for the keypair, with a maximum length of 255 ASCII characters.
* `publicKey` - (Required)	The public key. If you are not using command line tools, it must be encoded in base64.



## Attributes Reference

The following attributes are exported:

* `keypair_fingerprint` -	The MD5 public key fingerprint as specified in section 4 of RFC 4716.	false	string
* `keypair_name` -	The keypair name you specified.	false	string
* `request_id` -	The ID of the request	false	string

See detailed information in [Describe KeyPair Importtion](http://docs.outscale.com/api_fcu/operations/Action_ImportKeyPair_get.html#_api_fcu-action_importkeypair_get).
