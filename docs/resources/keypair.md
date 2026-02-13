---
layout: "outscale"
page_title: "OUTSCALE: outscale_keypair"
subcategory: "OUTSCALE API"
sidebar_current: "outscale-keypair"
description: |-
  [Manages a keypair.]
---

# outscale_keypair Resource

Manages a keypair.

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Keypairs.html).  
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-keypair).

## Example Usage

### Create a keypair

```hcl
resource "outscale_keypair" "keypair01" {
	keypair_name = "terraform-keypair-create"
}
```

### Import keypairs

```hcl
resource "outscale_keypair" "keypair02" {
	keypair_name = "terraform-keypair-import-file"
	public_key   = file("<PATH>")
}

resource "outscale_keypair" "keypair03" {
	keypair_name = "terraform-keypair-import-text"
	public_key   = "UFVCTElDIEtFWQ=="
}
```

## Argument Reference

The following arguments are supported:

* `keypair_name` - (Required) A unique name for the keypair, with a maximum length of 255 [ASCII printable characters](https://en.wikipedia.org/wiki/ASCII#Printable_characters).
* `public_key` - (Optional) The public key to import in your account, if you are importing an existing keypair. This value must be Base64-encoded.
* `tags` - (Optional) A tag to add to this resource. You can specify this argument several times.
    * `key` - (Required) The key of the tag, with a minimum of 1 character.
    * `value` - (Required) The value of the tag, between 0 and 255 characters.

## Attribute Reference

The following attributes are exported:

* `keypair_fingerprint` - The MD5 public key fingerprint, as specified in section 4 of RFC 4716.
* `keypair_id` - The ID of the keypair.
* `keypair_name` - The name of the keypair.
* `keypair_type` - The type of the keypair (`ssh-rsa`, `ssh-ed25519`, `ecdsa-sha2-nistp256`, `ecdsa-sha2-nistp384`, or `ecdsa-sha2-nistp521`).
* `private_key` - The private key, returned only if you are creating a keypair (not if you are importing). When you save this private key in a .rsa file, make sure you replace the `\n` escape sequences with real line breaks.
* `tags` - One or more tags associated with the keypair.
    * `key` - The key of the tag, with a minimum of 1 character.
    * `value` - The value of the tag, between 0 and 255 characters.

## Timeouts

The `timeouts` block enables you to configure [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - Defaults to 10 minutes.
* `read` - Defaults to 5 minutes.
* `update` - Defaults to 10 minutes.
* `delete` - Defaults to 5 minutes.

## Import

A keypair can be imported using its name. For example:

```console

$ terraform import outscale_keypair.ImportedKeypair keypair_id

```