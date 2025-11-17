---
layout: "outscale"
page_title: "OUTSCALE: outscale_api_access_policy"
subcategory: "Identity Access Management (IAM)"
sidebar_current: "outscale-api-access-policy"
description: |-
  [Provides information about the API access policy.]
---

# outscale_api_access_policy Data Source

Provides information about the API access policy.

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Your-API-Access-Policy.html).  
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-apiaccesspolicy).

## Example Usage

```hcl
data "outscale_api_access_policy" "unique" {
}
```

## Argument Reference

No argument is supported.

## Attribute Reference

The following attributes are exported:

* `max_access_key_expiration_seconds` - The maximum possible lifetime for your access keys, in seconds. If `0`, your access keys can have unlimited lifetimes.
* `require_trusted_env` - If true, a trusted session is activated, allowing you to bypass Certificate Authorities (CAs) enforcement. For more information, see [About Your API Access Policy](https://docs.outscale.com/en/userguide/About-Your-API-Access-Policy.html).<br />
If this is enabled, it is required that you and all your users log in to Cockpit v2 using the WebAuthn method for multi-factor authentication. For more information, see [About Authentication > Multi-Factor Authentication](https://docs.outscale.com/en/userguide/About-Authentication.html#_multi_factor_authentication).
