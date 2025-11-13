---
layout: "outscale"
page_title: "OUTSCALE: outscale_accounts"
subcategory: "Account"
sidebar_current: "outscale-accounts"
description: |-
  [Provides information about accounts.]
---

# outscale_accounts Data Source

Provides information about accounts.

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Your-Account.html).  
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-account).

## Example Usage

```hcl
data "outscale_accounts" "all_accounts" {
  
}
```

## Argument Reference

No argument is supported.

## Attribute Reference

The following attributes are exported:

* `accounts` - The list of the accounts.
    * `account_id` - The ID of the account.
    * `additional_emails` - One or more additional email addresses for the account. These addresses are used for notifications only.
    * `city` - The city of the account owner.
    * `company_name` - The name of the company for the account.
    * `country` - The country of the account owner.
    * `customer_id` - The ID of the customer.
    * `email` - The main email address for the account. This address is used for your credentials and for notifications.
    * `first_name` - The first name of the account owner.
    * `job_title` - The job title of the account owner.
    * `last_name` - The last name of the account owner.
    * `mobile_number` - The mobile phone number of the account owner.
    * `phone_number` - The landline phone number of the account owner.
    * `state_province` - The state/province of the account.
    * `vat_number` - The value added tax (VAT) number for the account.
    * `zip_code` - The ZIP code of the city.
