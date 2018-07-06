---
layout: "outscale"
page_title: "OUTSCALE: outscale_account"
sidebar_current: "docs-outscale-datasource-account"
description: |-
  Gets information about the account that sent the request.
---

# outscale_account

Gets information about the account that sent the request.

## Example Usage

```hcl
data "outscale_account" "account" {}
```

## Argument Reference

No arguments are supported

## Attributes Reference

The following attributes are exported:

* `account_pid` - The personal identifier (PID) of the account.
* `city` - The city of the account owner.
* `company_name` - The name of the company for the account.
* `country` - The country of the account owner.
* `customer_id` - The ID of the customer.
* `email` - The email address for the account.
* `first_name` - The first name of the account owner.
* `job_title` - The job title of the account owner.
* `last_name` - The last name of the account owner.
* `mobile_number` - The mobile phone number of the account owner.
* `phone_number` - The landline phone number of the account owner.
* `state` - The state of the account owner.
* `vat_number` - The VAT number of the account.
* `zipcode` - The zip code of the city.
* `request_id` - The ID of the request.

[See detailed description](http://docs.outscale.com/api_icu/operations/Action_GetAccount_get.html#_api_icu-action_getaccount_get)
