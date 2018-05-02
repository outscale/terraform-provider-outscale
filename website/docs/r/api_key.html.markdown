---
layout: "outscale"
page_title: "OUTSCALE: outscale_api_key"
sidebar_current: "docs-outscale-resource-api-key"
description: |-
  Creates a new secret access key and the corresponding access key ID for the account that sends the request.
---

# outscale_api_key

Creates a new secret access key and the corresponding access key ID for the account that sends the request. Instances also support [provisioning](/docs/provisioners/index.html).

## Example Usage

```hcl
resource "outscale_api_key" "outscale_api_key" {
    
    tag = {
        Name = "api_key_test"
    }
    
}
```

## Argument Reference

The following arguments are supported:

* `api_key_id` - (Optional)	The ID of the access key. If not provided, it will be created automatically.
* `secret_key` - (Optional)	The secret access key that enables you to send requests. If not provided, it will be created automatically.
* `tag` - (Optional)	A list of tags to add to the specified resources.



## Attributes Reference

The following attributes are exported:

* `api_keyId` - 	The ID of the access key.	
* `account_id` - 	The account ID of the owner of the access key.
* `secret_key` - 	The secret access key that enables you to send requests.
* `state` - 	The state of the access key (active if the key is valid for API calls, or inactive if not).
* `tag` - 	One or more tags associated with the access key.

See detailed information in [Describe API Key](http://docs.outscale.com/api_icu/operations/Action_CreateAccessKey_get.html#_api_icu-action_createaccesskey_get).
