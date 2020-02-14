---
layout: "outscale"
page_title: "3DS OUTSCALE: outscale_snapshot_attributes"
sidebar_current: "outscale-snapshot-attributes"
description: |-
  [Manages snapshot attributes.]
---

# outscale_snapshot_attributes Resource

Manages snapshot attributes.
For more information on this resource, see the [User Guide](https://wiki.outscale.net/display/EN/About+Snapshots#AboutSnapshots-SnapshotPermissionsandCopy).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#updatesnapshot).

## Example Usage

```hcl

#resource "outscale_volume" "volume01" {
#  subregion_name = "eu-west-2a"
#  size           = 40
#}

#resource "outscale_snapshot" "snapshot01" {
#  volume_id = outscale_volume.volume01.volume_id
#  tags {
#    key   = "name"
#    value = "terraform-snapshot-test"
#  }
#}

# Add permissions

resource "outscale_snapshot_attributes" "snapshot_attributes01" {
  snapshot_id = outscale_snapshot.snapshot01.snapshot_id
  permissions_to_create_volume_additions {
    account_ids = ["012345678910"]
  }
}

# Remove permissions

resource "outscale_snapshot_attributes" "snapshot_attributes02" {
  snapshot_id = outscale_snapshot.snapshot01.snapshot_id
  permissions_to_create_volume_removals {
    account_ids = ["012345678910"]
  }
}


```

## Argument Reference

The following arguments are supported:

* `snapshot_id` - (Required) The ID of the snapshot.
* `permissions_to_create_volume_additions` - (Optional) Information about the users you want to give permissions for the resource.
  * `account_ids` - (Optional) The account ID of one or more users you want to give permissions to.
  * `global_permission` - (Optional) If `true`, the resource is public. If `false`, the resource is private.
* `permissions_to_create_volume_removals` - (Optional) Information about the users you want to remove permissions for the resource.
  * `account_ids` - (Optional) The account ID of one or more users you want to remove permissions from.
  * `global_permission` - (Optional) If `true`, the resource is public. If `false`, the resource is private.

## Attribute Reference

The following attributes are exported:

* `snapshot_id` - The ID of the snapshot.
* `permissions_to_create_volume_additions` - Information about the permissions for the resource.
  * `account_ids` - The account ID of one or more users who have permissions for the resource.
  * `global_permission` - If `true`, the resource is public. If `false`, the resource is private.
