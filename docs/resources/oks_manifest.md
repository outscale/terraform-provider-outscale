---
layout: "outscale"
page_title: "OUTSCALE: outscale_oks_manifest"
subcategory: "OKS API"
sidebar_current: "outscale-oks-manifest"
description: |-
  [Manages a manifest.]
---

# outscale_oks_manifest Resource

Applies a generic Kubernetes manifest on an OKS cluster using the cluster's kubeconfig.

The manifest can be written from a Custom Resource Definition (CRD) retrieved from the Kubernetes API, or from the templates returned by the [outscale_oks_crd_templates](https://registry.terraform.io/providers/outscale/outscale/latest/docs/data-sources/oks_crd_templates) data source. Both cluster-scoped and namescape-scoped CRDs are supported.

This uses server-side apply (SSA) for each Kubernetes apply/validation call. Create and update operations use patch SSA.

-> **Note:** This resource is a smaller implementation of what the [kubernetes_manifest](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/manifest) and the [kubectl_manifest](https://registry.terraform.io/providers/hashicorp-oss/kubectl/latest/docs/resources/manifest) resources offer, that integrates natively with an OKS cluster without needing to initialize another Terraform provider.

## Example Usage

### Required resources

```hcl
resource "outscale_oks_project" "project" {
  name            = "project01"
  cidr            = "10.50.0.0/18"
  region          = "eu-west-2"
}

resource "outscale_oks_cluster" "cluster" {
  project_id      = outscale_oks_project.project.id
  admin_whitelist = ["0.0.0.0/0"]
  cidr_pods       = "10.91.0.0/16"
  cidr_service    = "10.92.0.0/16"
  version         = "1.35"
  name            = "cluster01"
  control_planes  = "cp.mono.master"
}
```

### Create an OKS node pool

```hcl
resource "outscale_oks_manifest" "nodepool" {
  cluster_id = outscale_oks_cluster.cluster.id
  manifest   = <<-YAML
apiVersion: oks.dev/v1beta2
kind: NodePool
metadata:
  name: pool-1
spec:
  autoHealing: true
  desiredNodes: 1
  nodeType: tinav7.c1r1p1
  upgradeStrategy:
    autoUpgradeEnabled: true
    autoUpgradeMaintenance:
      durationHours: 1
      startHour: 12
      weekDay: Tue
    maxSurge: 0
    maxUnavailable: 1
  volumes:
  - device: root
    dir: /
    size: 100
    type: gp2
  zones: [eu-west-2a]
YAML

  wait_for = {
    fields = {
      "status.progress.ready" = "1"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required) The cluster ID on which the manifest is applied.
* `manifest` - (Required) The Kubernetes YAML manifest.
* `skip_delete` (Optional, default: false) If set to true, Terraform removes only the manifest resource from the state during delete, without deleting the Kubernetes object from the cluster.
* `wait` - (Optional, default: false): If set to true, Terraform waits for the Kubernetes object to be deleted during destroy.
* `wait_for` - (Optional) Wait until the fields in the Kubernetes object match the expected values after apply:
    * `fields` - (Required) Maps of key-value pairs (field path => expected pattern).
        * Each key must be a [JSONPath](https://kubernetes.io/docs/reference/kubectl/jsonpath/) field path, but a simple field path such as `status.progress.ready` is accepted and converted internally to a JSONPath expression (that is, the enclosing `{` `}` and the first `.` can be omitted).
        * Each value is a regex pattern.
        * All configured fields must match for the wait to complete.
        * Examples: `"status.progress.ready" = "1"`, `"{.status.progress.ready}" = "1"`, `"{.status.state.name}" = "idle|reconciliation"`.
    * `timeout` - (Optional) Custom timeout for the `wait_for` checks. If not specified, falls back to the CRUD operation default timeout.

## Attribute Reference

The following attribute is exported:

* `object` - The applied Kubernetes object returned as YAML.

## Timeouts

The `timeouts` block enables you to configure [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - Defaults to 15 minutes.
* `read` - Defaults to 2 minutes.
* `update` - Defaults to 10 minutes.
* `delete` - Defaults to 10 minutes.
