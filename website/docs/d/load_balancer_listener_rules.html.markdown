---
layout: "outscale"
page_title: "OUTSCALE: outscale_load_balancer_listener_rules"
sidebar_current: "outscale-load-balancer-listener-rules"
description: |-
  [Provides information about load balancer listener rules.]
---

# outscale_load_balancer_listener_rules Data Source

Provides information about load balancer listener rules.
For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Load-Balancers.html).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-listener).

## Example Usage

```hcl
data "outscale_load_balancer_listener_rules" "rules01" {
  filter {
    name   = "listener_rule_names"
    values = ["terraform-listener-rule02","terraform-listener-rule01"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) A combination of a filter name and one or more filter values. You can specify this argument for as many filter names as you need. The filter name can be any of the following:
    * `listener_rule_names` - (Optional) The names of the listener rules.

## Attribute Reference

The following attributes are exported:

* `listener_rules` - The list of the rules to describe.
    * `action` - The type of action for the rule (always `forward`).
    * `host_name_pattern` - A host-name pattern for the rule, with a maximum length of 128 characters. This host-name pattern supports maximum three wildcards, and must not contain any special characters except [-.?].
    * `listener_id` - The ID of the listener.
    * `listener_rule_id` - The ID of the listener rule.
    * `listener_rule_name` - A human-readable name for the listener rule.
    * `path_pattern` - A path pattern for the rule, with a maximum length of 128 characters. This path pattern supports maximum three wildcards, and must not contain any special characters except [_-.$/~&quot;'@:+?].
    * `priority` - The priority level of the listener rule, between `1` and `19999` both included. Each rule must have a unique priority level. Otherwise, an error is returned.
    * `vm_ids` - The IDs of the backend VMs.
