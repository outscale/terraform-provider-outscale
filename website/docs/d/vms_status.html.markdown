---
layout: "outscale"
page_title: "OUTSCALE: outscale_vms_state"
sidebar_current: "docs-outscale-datasource-vms-state"
description: |-
  Describes the status of one or more instances.
---

# outscale_vms_state

Describes the status of one or more instances.

## Example Usage

```hcl
data "outscale_vms_state" "state" {
  instance_id = ["i-5adcfa0f", "i-31f8427d"]
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) One or more filters.
* `instance_id` - (Optional) The IDs of the instances.

See detailed information in [Outscale VM Status](http://docs.outscale.com/api_fcu/operations/Action_DescribeInstanceStatus_get.html#_api_fcu-action_describeinstancestatus_get).

## Filters

Use the Filter.N parameter to filter the described instances on the following properties:

* `availability-zone` - The Availability Zone of the instance.
* `event.code` - The code for the scheduled event (`system-reboot` | `system-maintenance`).
* `event.description` - Indicates whether the BSU volume is deleted when terminating the instance.
* `event.not-after` - The latest end time for a scheduled event (for example, `2016-01-23T18:45:30.000Z`).
* `event.not-before` - The earliest start time for a scheduled event (for example, `2016-01-23T18:45:30.000Z`).
* `instance-state-code` - The state of the instance (a 16-bit unsigned integer). The high byte is an internal value you should ignore. The low byte represents the state of the instance: `0` (pending), `16` (running), `32` (shutting-down), `48` (terminated), `64` (stopping), or `80` (stopped).
* `client-token` - The idempotency token provided when launching the instance.
* `instance-state-name` - The state of the instance (`pending` | `running` | `shutting-down` | `terminated` | `stopping` | `stopped`).

## Attributes Reference

The following attributes are exported:

* `instance_status_set` - One or more instance status descriptions.
* `request_id` - The ID of the request.

See detailed information in [Instance Status Set](http://docs.outscale.com/api_fcu/definitions/InstanceStatus.html#_api_fcu-instancestatus).
