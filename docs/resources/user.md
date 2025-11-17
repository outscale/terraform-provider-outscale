---
layout: "outscale"
page_title: "OUTSCALE: outscale_user"
subcategory: "OUTSCALE API"
sidebar_current: "outscale-user"
description: |-
  [Manages a user.]
---

# outscale_user Resource

Manages a user.

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-EIM-Users.html).  
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api.html#createuser).

## Example Usage

### Creating a user

```hcl
resource "outscale_user" "user-1"  {
    user_name  = "User-TF-1"
    user_email = "test-TF1@test2.fr"
    path       = "/terraform/"
}
```

### Linking a policy to a user

```hcl
resource "outscale_user" "user-1"  {
    user_name = "User-TF-1"
    policy {
        policy_orn         = outscale_policy.policy-1.orn
        default_version_id = "V1"
    }
}
```

## Argument Reference

The following arguments are supported:

* `default_version_id` - The ID of a policy version that you want to make the default one (the active one).
* `path` - (Optional) The path to the EIM user you want to create (by default, `/`). This path name must begin and end with a slash (`/`), and contain between 1 and 512 alphanumeric characters and/or slashes (`/`), or underscores (`_`).
* `user_email` - (Optional) The email address of the EIM user.
* `user_name` - (Required) The name of the EIM user. This user name must contain between 1 and 64 alphanumeric characters and/or pluses (`+`), equals (`=`), commas (`,`), periods (`.`), at signs (`@`), dashes (`-`), or underscores (`_`).

## Attribute Reference

The following attributes are exported:

* `creation_date` - The date and time (UTC) of creation of the EIM user.
* `last_modification_date` - The date and time (UTC) of the last modification of the EIM user.
* `path` - The path to the EIM user.
* `user_email` - The email address of the EIM user.
* `user_id` - The ID of the EIM user.
* `user_name` - The name of the EIM user.

## Import

A user can be imported using its ID. For example:

```console

$ terraform import outscale_user.user1 user_id

```