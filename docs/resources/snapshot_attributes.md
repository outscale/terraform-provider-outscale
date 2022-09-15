---
layout: "outscale"
page_title: "OUTSCALE: outscale_snapshot_attributes"
sidebar_current: "outscale-snapshot-attributes"
description: |-
  [Manages snapshot attributes.]
---

# outscale_snapshot_attributes Resource

Manages snapshot attributes.
For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Snapshots.html).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#updatesnapshot).

## Example Usage

### Required resources

```hcl
resource "outscale_volume" "volume01" {
	subregion_name = "eu-west-2a"
	size           = 40
}

resource "outscale_snapshot" "snapshot01" {
	volume_id = outscale_volume.volume01.volume_id
	tags {
		key   = "name"
		value = "terraform-snapshot-test"
	}
}
```

### Add permissions

```hcl
resource "outscale_snapshot_attributes" "snapshot_attributes01" {
	snapshot_id = outscale_snapshot.snapshot01.snapshot_id
	permissions_to_create_volume_additions {
		account_ids = ["012345678910"]
	}
}
```

### Remove permissions

```hcl
resource "outscale_snapshot_attributes" "snapshot_attributes02" {
	snapshot_id = outscale_snapshot.snapshot01.snapshot_id
	permissions_to_create_volume_removals {
		account_ids = ["012345678910"]
	}
}
```

## Argument Reference

The following arguments are supported:

* `permissions_to_create_volume_additions` - (Optional) Information about the users to whom you want to give permissions for the resource.
    * `account_ids` - (Optional) The account ID of one or more users to whom you want to give permissions.
    * `global_permission` - (Optional) If true, the resource is public. If false, the resource is private.
* `permissions_to_create_volume_removals` - (Optional) Information about the users from whom you want to remove permissions for the resource.
    * `account_ids` - (Optional) The account ID of one or more users from whom you want to remove permissions.
    * `global_permission` - (Optional) If true, the resource is public. If false, the resource is private.
* `snapshot_id` - (Required) The ID of the snapshot.

## Attribute Reference

The following attributes are exported:

* `account_id` - The account ID of the owner of the snapshot.
* `snapshot_id` - The ID of the snapshot.

