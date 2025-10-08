---
layout: "outscale"
page_title: "OUTSCALE: outscale_keypairs"
subcategory: "Keypair"
sidebar_current: "outscale-keypairs"
description: |-
  [Provides information about keypairs.]
---

# outscale_keypairs Data Source

Provides information about keypairs.

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Keypairs.html).  
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

* `filter` - (Optional) A combination of a filter name and one or more filter values. You can specify this argument for as many filter names as you need. The filter name can be any of the following:
    * `keypair_fingerprints` - (Optional) The fingerprints of the keypairs.
    * `keypair_ids` - (Optional) The IDs of the keypairs.
    * `keypair_names` - (Optional) The names of the keypairs.
    * `keypair_types` - (Optional) The types of the keypairs (`ssh-rsa`, `ssh-ed25519`, `ecdsa-sha2-nistp256`, `ecdsa-sha2-nistp384`, or `ecdsa-sha2-nistp521`).
    * `tag_keys` - (Optional) The keys of the tags associated with the keypairs.
    * `tag_values` - (Optional) The values of the tags associated with the keypairs.
    * `tags` - (Optional) The key/value combination of the tags associated with the keypairs, in the following format: &quot;Filters&quot;:{&quot;Tags&quot;:[&quot;TAGKEY=TAGVALUE&quot;]}.

## Attribute Reference

The following attributes are exported:

* `keypairs` - Information about one or more keypairs.
    * `keypair_fingerprint` - The MD5 public key fingerprint as specified in section 4 of RFC 4716.
    * `keypair_id` - The ID of the keypair.
    * `keypair_name` - The name of the keypair.
    * `keypair_type` - The type of the keypair (`ssh-rsa`, `ssh-ed25519`, `ecdsa-sha2-nistp256`, `ecdsa-sha2-nistp384`, or `ecdsa-sha2-nistp521`).
    * `tags` - One or more tags associated with the keypair.
        * `key` - The key of the tag, with a minimum of 1 character.
        * `value` - The value of the tag, between 0 and 255 characters.
