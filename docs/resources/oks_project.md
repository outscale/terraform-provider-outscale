---
layout: "outscale"
page_title: "OUTSCALE: outscale_oks_project"
subcategory: "OKS"
sidebar_current: "outscale-oks-project"
description: |-
  [Manages a project.]
---

# outscale_oks_project Resource

Manages a project.

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-OKS.html#_projects).  
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/oks.html#oks-api-projects).

## Example Usage

```hcl
resource "outscale_oks_project" "project01" {
  name   = "project01"
  cidr   = "10.50.0.0/18"
  region = "eu-west-2"
  tags   = {
    tagkey = "tagvalue"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A unique name for the project, between 1 and 40 characters. This name must start with a letter and contain only lowercase letters, numbers, or hyphens.
* `description` - (Optional) A description for the project.
* `cidr` - (Required) The CIDR block to associate with the Net of the project.
* `region` - (Required) The Region on which the project is deployed.
* `tags` - (Optional) The key/value combinations of the tags associated with the resource.
* `quirks` - (Optional) A list of special configurations or behaviors for the project.
* `disable_api_termination` - (Optional) If true, project deletion through the API is disabled. If false, it is enabled. By default, false.

## Attribute Reference

The following attributes are exported:

* `id` - The ID of the project.
* `name` - The name of the project.
* `description` - The description of the project.
* `cidr` - The CIDR block associated with the Net of the project.
* `region` - The Region on which the project is deployed.
* `status` - The status of the project.
* `tags` - The key/value combinations of the tags associated with the resource.
* `disable_api_termination` - If true, project deletion through the API is disabled. If false, it is enabled.
* `created_at` - The timestamp when the project was created (date-time).
* `updated_at` - The timestamp when the project was last updated (date-time).
* `deleted_at` - The timestamp when the project was deleted (if applicable) (date-time).

## Import

An OKS project can be imported using its ID. For example:

```console

$ terraform import outscale_oks_project.project id

```

