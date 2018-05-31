---
layout: "outscale"
page_title: "OUTSCALE: outscale_api_key"
sidebar_current: "docs-outscale-datasource-api-key"
description: |-
Describes client endpoint
---

# outscale_api_key

Describes the api key

## Example Usage

```hcl
data "outscale_api_key" "test" {}
```

## Argument Reference

No arguments are supported

## Attributes Reference

The following attributes are exported:

* `access_key.N` - A list of access keys and their metadata.
  - `access_key_id` - The ID of the access key.
  - `owner_id` - The account name of the user the access key is associated with.
  - `status` - The state of the access key (active if the key is valid for API calls, or inactive if not).

[See detailed description](http://docs.outscale.com/api_eim/operations/Action_ListAccessKeys_get.html#_api_eim-action_listaccesskeys_get)
