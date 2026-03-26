---
layout: "outscale"
page_title: "OUTSCALE: outscale_oks_quotas"
subcategory: "OKS API"
sidebar_current: "outscale-oks-quotas"
description: |-
  [Provides information about OKS quotas.]
---

# outscale_oks_quotas Data Source

Provides information about OKS quotas.

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/Getting-Information-About-the-Quotas-of-a-Profile.html).  
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/oks.html#getquotas).

## Example Usage

```hcl
data "outscale_oks_quotas" "oks_quotas" {
}
```

## Argument Reference

No argument is supported.

## Attribute Reference

The following attributes are exported:

* `clusters_per_project` - The maximum allowed number of clusters per project.
* `cp_subregions` - The list of available Subregions.
* `kube_versions` - The list of available Kubernetes versions.
* `projects` - The maximum allowed number of projects.
