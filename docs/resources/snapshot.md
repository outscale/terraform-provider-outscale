---
layout: "outscale"
page_title: "OUTSCALE: outscale_snapshot"
subcategory: "OUTSCALE API"
sidebar_current: "outscale-snapshot"
description: |-
  [Manages a snapshot.]
---

# outscale_snapshot Resource

Manages a snapshot.

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Snapshots.html).  
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-snapshot).

## Example Usage

### Required resource

```hcl
resource "outscale_volume" "volume01" {
	subregion_name = "${var.region}a"
	size           = 40
}
```

### Create a snapshot

```hcl
resource "outscale_snapshot" "snapshot01" {
	volume_id = outscale_volume.volume01.volume_id
}
```

### Copy a snapshot

```hcl
resource "outscale_snapshot" "snapshot02" {
	description        = "Terraform snapshot copy"
	source_snapshot_id = "snap-12345678"
	source_region_name = "eu-west-2"
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) A description for the snapshot.
* `file_location` - (Optional) **(when importing from a bucket)** The pre-signed URL of the snapshot you want to import. For more information, see [Creating a Pre-signed URL](https://docs.outscale.com/en/userguide/Creating-a-Pre-Signed-URL.html).
* `snapshot_size` - (Optional) **(when importing from a bucket)** The size of the snapshot you want to create in your account, in bytes. This size must be greater than or equal to the size of the original, uncompressed snapshot.
* `source_region_name` - (Optional) **(when copying a snapshot)** The name of the source Region, which must be the same as the Region of your account.
* `source_snapshot_id` - (Optional) **(when copying a snapshot)** The ID of the snapshot you want to copy.
* `tags` - (Optional) A tag to add to this resource. You can specify this argument several times.
    * `key` - (Required) The key of the tag, with a minimum of 1 character.
    * `value` - (Required) The value of the tag, between 0 and 255 characters.
* `volume_id` - (Optional) **(when creating from a volume)** The ID of the volume you want to create a snapshot of.

## Attribute Reference

The following attributes are exported:

* `account_alias` - The account alias of the owner of the snapshot.
* `account_id` - The account ID of the owner of the snapshot.
* `creation_date` - The date and time (UTC) at which the snapshot was created.
* `description` - The description of the snapshot.
* `permissions_to_create_volume` - Permissions for the resource.
    * `account_ids` - One or more account IDs that the permission is associated with.
    * `global_permission` - A global permission for all accounts.<br />
(Request) Set this parameter to true to make the resource public (if the parent parameter is `Additions`) or to make the resource private (if the parent parameter is `Removals`).<br />
(Response) If true, the resource is public. If false, the resource is private.
* `progress` - The progress of the snapshot, as a percentage.
* `snapshot_id` - The ID of the snapshot.
* `state` - The state of the snapshot (`in-queue` \| `pending` \| `completed` \| `error` \| `deleting`).
* `tags` - One or more tags associated with the snapshot.
    * `key` - The key of the tag, with a minimum of 1 character.
    * `value` - The value of the tag, between 0 and 255 characters.
* `volume_id` - The ID of the volume used to create the snapshot.
* `volume_size` - The size of the volume used to create the snapshot, in gibibytes (GiB).

## Timeouts

The `timeouts` block enables you to configure [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - Defaults to 40 minutes.
* `read` - Defaults to 5 minutes.
* `update` - Defaults to 10 minutes.
* `delete` - Defaults to 5 minutes.

## Import

A snapshot can be imported using its ID. For example:

```console

$ terraform import terraform import outscale_snapshot.ImportedSnapshot snap-12345678

```