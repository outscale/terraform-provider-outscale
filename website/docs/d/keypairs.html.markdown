---
layout: "outscale"
page_title: "OUTSCALE: outscale_keypairs"
sidebar_current: "docs-outscale-datasource-keypairs"
description: |-
Describes your keypairs.
---

# outscale_keypair

Describes one or more of your keypairs.

## Example Usage

```hcl
resource "outscale_key_pair" "outscale_keypair" {
    count = 1

    key_name = "keyname_test_"
}

resource "outscale_key_pair" "outscale_keypair2" {
    count = 1

    key_name = "keyname_test_"
}




```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) One or more filters
* `key_names` - (Optional) One or more keypairs.


See detailed information in [Outscale Instances](http://docs.outscale.com/api_fcu/operations/Action_DescribeKeyPairs_get.html#_api_fcu-action_describekeypairs_get).

## Filters

Use the Filter.N parameter to filter the described instances on the following properties:

* `fingerprint` The fingerprint of the keypair.
* `key-name` The name of the keypair


## Attributes Reference

The following attributes are exported:

* `request_id` - The ID of the request.
* `key_sets` - Information about one or more keypairs.







See detailed information in [Describe Instances](http://docs.outscale.com/api_fcu/operations/Action_DescribeKeyPairs_get.html#_api_fcu-action_describekeypairs_get).
