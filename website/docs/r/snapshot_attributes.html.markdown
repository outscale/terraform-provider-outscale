---
layout: "outscale"
page_title: "OUTSCALE: outscale_snapshot_attributes"
sidebar_current: "docs-outscale-resource-snapshot-attributes"
description: |-
  Modifies the permissions for a specified snapshot.
---

# outscale_snapshot_attributes

Modifies the permissions for a specified snapshot.
You can add or remove permissions for specified user IDs or groups. You can share a snapshot with a user that is in the same Region. The user can create a copy of the snapshot you shared, obtaining all the rights for the copy of the snapshot. 

## Example Usage

```hcl
resource "outscale_volume" "description_test" {
    availability_zone = "eu-west-2a"
    size = 1
}

resource "outscale_snapshot" "test" {
    volume_id = "${outscale_volume.description_test.id}"
    description = "Snapshot Acceptance Test"
}
```

## Argument Reference

The following arguments are supported:

* `snapshot_id` - The ID of the snapshot.
* `create_volume_permission_add` - (Optional) Enables you to modify the permissions to create a volume for user IDs or groups.
  * `user_id` (Optional) The account ID of the user.
  * `group_id` (Optional) The name of the group (all if public).

## Attributes Reference

The following attributes are exported:

* `account_id` The account ID.

[See detailed description](http://docs.outscale.com/api_fcu/operations/Action_ModifySnapshotAttribute_get.html#_api_fcu-action_modifysnapshotattribute_get)
