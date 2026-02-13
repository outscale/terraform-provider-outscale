---
layout: "outscale"
page_title: "OUTSCALE: outscale_api_access_rule"
subcategory: "OUTSCALE API"
sidebar_current: "outscale-api-access-rule"
description: |-
  [Manages an API access rule.]
---

# outscale_api_access_rule Resource

Manages an API access rule.

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-API-Access-Rules.html).  
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-apiaccessrule).

## Example Usage

### Create an API access rule based on IPs

```hcl
resource "outscale_api_access_rule" "api_access_rule01" {
    ip_ranges   = ["192.0.2.0", "192.0.2.0/16"]
    description = "Basic API Access Rule from Terraform"
}
```

### Create an API access rule based on IPs and Certificate Authority (CA)

```hcl
resource "outscale_ca" "ca01" {
    ca_pem      = file("<PATH>")
    description = "Terraform CA"
}

resource "outscale_api_access_rule" "api_access_rule02" {
    ip_ranges   = ["192.0.2.0", "192.0.2.0/16"]
    ca_ids      = [outscale_ca.ca01.ca_id]
    description = "API Access Rule with CA from Terraform"
}
```

## Argument Reference

The following arguments are supported:

* `ca_ids` - (Optional) One or more IDs of Client Certificate Authorities (CAs).
* `cns` - (Optional) One or more Client Certificate Common Names (CNs). If this parameter is specified, you must also specify the `ca_ids` parameter.
* `description` - (Optional) A description for the API access rule.
* `ip_ranges` - (Optional) One or more IPs or CIDR blocks (for example, `192.0.2.0/16`).

## Attribute Reference

The following attributes are exported:

* `api_access_rule_id` - The ID of the API access rule.
* `ca_ids` - One or more IDs of Client Certificate Authorities (CAs) used for the API access rule.
* `cns` - One or more Client Certificate Common Names (CNs).
* `description` - The description of the API access rule.
* `ip_ranges` - One or more IP ranges used for the API access rule, in CIDR notation (for example, `192.0.2.0/16`).

## Timeouts

The `timeouts` block enables you to configure [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - Defaults to 10 minutes.
* `read` - Defaults to 5 minutes.
* `update` - Defaults to 10 minutes.
* `delete` - Defaults to 5 minutes.

## Import

An API access rule can be imported using its ID. For example:

```console

$ terraform import outscale_api_access_rule.ImportedAPIAccessRule "aar-12345678"

```