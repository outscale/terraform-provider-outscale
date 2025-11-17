---
layout: "outscale"
page_title: "OUTSCALE: outscale_entities_linked_to_policy"
subcategory: "OUTSCALE API"
sidebar_current: "outscale-entities-linked-to-policy"
description: |-
  [Provides information about  entities (account, users, or user groups) linked to a specific managed policy.]
---

# outscale_entities_linked_to_policy Data Source

Provides information about  entities (account, users, or user groups) linked to a specific managed policy.

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Policies.html).  
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api.html#readentitieslinkedtopolicy).

## Example Usage

```hcl
data "outscale_entities_linked_to_policy" "entities_linked_policy01" {
    policy_orn    = "orn:ows:idauth::012345678910:policy/example/example-policy"
    entities_type = ["USER","GROUP","ACCOUNT"]
}
```

## Argument Reference

The following arguments are supported:

* `entities_type` - (Optional) The type of entity linked to the policy (`ACCOUNT` \| `USER` \| `GROUP`) you want to get information about.
* `first_item` - (Optional) The item starting the list of entities requested.
* `policy_orn` - (Optional) The OUTSCALE Resource Name (ORN) of the policy. For more information, see [Resource Identifiers](https://docs.outscale.com/en/userguide/Resource-Identifiers.html).
* `results_per_page` - (Optional) The maximum number of items that can be returned in a single response (by default, 100).

## Attribute Reference

The following attributes are exported:

* `accounts` - The accounts linked to the specified policy.
    * `id` - The ID of the entity.
    * `name` - The name of the entity.
    * `orn` - The OUTSCALE Resource Name (ORN) of the entity. For more information, see [Resource Identifiers](https://docs.outscale.com/en/userguide/Resource-Identifiers.html).
* `groups` - The groups linked to the specified policy.
    * `id` - The ID of the entity.
    * `name` - The name of the entity.
    * `orn` - The OUTSCALE Resource Name (ORN) of the entity. For more information, see [Resource Identifiers](https://docs.outscale.com/en/userguide/Resource-Identifiers.html).
* `has_more_items` - If true, there are more items to return using the `first_item` parameter in a new request.
* `items_count` - The number of entities the specified policy is linked to.
* `max_results_limit` - Indicates maximum results defined for the operation.
* `max_results_truncated` - If true, indicates whether requested page size is more than allowed.
* `users` - The users linked to the specified policy.
    * `id` - The ID of the entity.
    * `name` - The name of the entity.
    * `orn` - The OUTSCALE Resource Name (ORN) of the entity. For more information, see [Resource Identifiers](https://docs.outscale.com/en/userguide/Resource-Identifiers.html).
