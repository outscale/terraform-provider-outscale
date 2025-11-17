---
layout: "outscale"
page_title: "OUTSCALE: outscale_user_group"
subcategory: "OUTSCALE API"
sidebar_current: "outscale-user-group"
description: |-
  [Manages a user group.]
---

# outscale_user_group Resource

Manages a user group.

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-EIM-Groups.html).  
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api.html#createusergroup).

## Example Usage

### Create a user group

```hcl
resource "outscale_user_group" "group-1" {
    user_group_name = "Group-TF-test-1"
    path            = "/terraform/"
}
```

### Link a policy to a user group

```hcl
resource "outscale_user_group" "group-1" {
    user_group_name = "Group-TF-test-1"
    policy {
        policy_orn         = outscale_policy.policy-2.orn
        default_version_id = "V2"
    }
}
```

### Add a user to a user group

```hcl
resource "outscale_user_group" "group-1" {
    user_group_name = "Group-TF-test-1"
    user {
        user_name = "user-name-1"
        path      = "/terraform/"
    }
    user {
        user_name = "user-name-2"
    }
}
```

### Create a user group, and add a user and a policy to it

```hcl
resource "outscale_user_group" "group-1" {
  user_group_name = "Group-TF-test-1"
    user {
        user_name = "user-name-1"
        path      = "/terraform/"
    }
    user {
        user_name = "user-name-2"
    }
    policy {
        policy_orn = outscale_policy.policy-2.orn
        version_id = "V2"
    }
}
```


## Argument Reference

The following arguments are supported:

* `default_version_id` - The ID of a policy version that you want to make the default one (the active one).
* `path` - (Optional) The path to the group. If not specified, it is set to a slash (`/`).
* `user_group_name` - (Required) The name of the group.

## Attribute Reference

The following attributes are exported:

* `creation_date` - The date and time (UTC) of creation of the user group.
* `last_modification_date` - The date and time (UTC) of the last modification of the user group.
* `name` - The name of the user group.
* `orn` - The Outscale Resource Name (ORN) of the user group. For more information, see [Resource Identifiers](https://docs.outscale.com/en/userguide/Resource-Identifiers.html).
* `path` - The path to the user group.
* `user_group_id` - The ID of the user group.

## Import

A user group can be imported using its group ID. For example:

```console

$ terraform import outscale_user_group.group1 user_group_id

```