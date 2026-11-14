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

* `cluster_id` - (Required) The ID of the cluster on which you want to apply the manifest.
* `manifest` - (Required) The Kubernetes YAML manifest.
* `force_conflicts` - (Optional) If set to true, server-side apply takes ownership of the fields that another field manager (for example `kubectl`, another Terraform provider or a controller) owns with a different value, instead of failing with a conflict. This also applies to the dry-run validation during plan. By default, false.
* `skip_delete` - (Optional) If set to true, Terraform removes only the manifest resource from the state during delete, without deleting the Kubernetes object from the cluster. By default, false.
* `wait` - (Optional): If set to true, Terraform waits for the Kubernetes object to be deleted during destroy. By default, false.
* `wait_for` - (Optional) Wait until the fields in the Kubernetes object match the expected values after apply:
    * `fields` - (Required) Maps of key/value pairs in the `"{.field_path}" = "expected_value"` format.<br />
Each key must be a [JSONPath](https://kubernetes.io/docs/reference/kubectl/jsonpath/) field path, but the enclosing characters (`{` `}`) and the first `.` are optional, and each value must be a regex pattern. All the configured fields must match for the wait to complete.<br />
Examples: `"{.status.progress.ready}" = "1"`, `"status.progress.ready" = "1"`, `"status.state.name" = "idle|reconciliation"`.
    * `timeout` - (Optional) A custom timeout for the `wait_for` checks. If not specified, falls back to the CRUD operation default timeout.

## Attribute Reference

The following attribute is exported:

* `object` - The applied Kubernetes object returned as YAML.

## Timeouts

The `timeouts` block enables you to configure [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - Defaults to 15 minutes.
* `read` - Defaults to 2 minutes.
* `update` - Defaults to 10 minutes.
* `delete` - Defaults to 10 minutes.

## Import

A Kubernetes object that already exists on an OKS cluster can be imported using an ID made of the cluster ID and the `apiVersion`, `kind` and `name` of the object, plus its `namespace` for a namespaced object. A namespaced object imported without a namespace is looked up in the `default` namespace, and the namespace of a cluster-scoped object is ignored. For example:

```console

$ terraform import outscale_oks_manifest.nodepool cluster_id=cluster_id,apiVersion=oks.dev/v1beta2,kind=NodePool,name=pool-1
$ terraform import outscale_oks_manifest.config cluster_id=cluster_id,apiVersion=v1,kind=ConfigMap,namespace=kube-system,name=example

```

With an `import` block:

```hcl
import {
  to = outscale_oks_manifest.nodepool
  id = "cluster_id=${outscale_oks_cluster.cluster.id},apiVersion=oks.dev/v1beta2,kind=NodePool,name=pool-1"
}
```

The imported `manifest` holds only the `apiVersion`, `kind` and `metadata` of the object, so the next apply updates the manifest in place, and server-side apply takes ownership of the configured fields. If another field manager owns some of these fields with a different value, the plan fails with a conflict unless `force_conflicts` is set to true. Once that apply is done, `force_conflicts` can be removed, so that a later conflict, with a controller for example, is reported instead of overridden.
