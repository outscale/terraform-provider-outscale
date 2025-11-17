---
layout: "outscale"
page_title: "OUTSCALE: outscale_api_access_policy"
subcategory: "OUTSCALE API"
sidebar_current: "outscale-api-access-policy"
description: |-
  [Manages the API access policy.]
---

# outscale_api_access_policy Resource

Manages the API access policy.

To activate a trusted session, first you must:
* Set expiration dates to all your access keys.
* Specify a Certificate Authority (CA) in all your API access rules.

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Your-API-Access-Policy.html).  
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-apiaccesspolicy).

## Example Usage

### Require expiration dates for your access keys

```hcl
resource "outscale_api_access_policy" "unique" {
    max_access_key_expiration_seconds = 31536000 # 1 year
    require_trusted_env               = false
}
```

### Activate a trusted session

```hcl
resource "outscale_api_access_policy" "unique" {
    max_access_key_expiration_seconds = 3153600000 # 100 years
    require_trusted_env               = true
}
```

### Deactivate a trusted session

```hcl
resource "outscale_api_access_policy" "unique" {
    max_access_key_expiration_seconds = 0
    require_trusted_env               = false
}
```

## Argument Reference

The following arguments are supported:

* `max_access_key_expiration_seconds` - (Required) The maximum possible lifetime for your access keys, in seconds (between `0` and `3153600000`, both included). If set to `O`, your access keys can have unlimited lifetimes, but a trusted session cannot be activated. Otherwise, all your access keys must have an expiration date. This value must be greater than the remaining lifetime of each access key of your account.
* `require_trusted_env` - (Required) If true, a trusted session is activated, provided that you specify the `max_access_key_expiration_seconds` parameter with a value greater than `0`.<br />
Enabling this will require you and all your users to log in to Cockpit v2 using the WebAuthn method for multi-factor authentication. For more information, see [About Authentication > Multi-Factor Authentication](https://docs.outscale.com/en/userguide/About-Authentication.html#_multi_factor_authentication).

## Attribute Reference

The following attributes are exported:

* `max_access_key_expiration_seconds` - The maximum possible lifetime for your access keys, in seconds. If `0`, your access keys can have unlimited lifetimes.
* `require_trusted_env` - If true, a trusted session is activated, allowing you to bypass Certificate Authorities (CAs) enforcement. For more information, see [About Your API Access Policy](https://docs.outscale.com/en/userguide/About-Your-API-Access-Policy.html).<br />
If this is enabled, it is required that you and all your users log in to Cockpit v2 using the WebAuthn method for multi-factor authentication. For more information, see [About Authentication > Multi-Factor Authentication](https://docs.outscale.com/en/userguide/About-Authentication.html#_multi_factor_authentication).

