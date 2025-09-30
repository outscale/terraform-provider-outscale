---
layout: "outscale"
page_title: "OUTSCALE: outscale_oks_cluster"
sidebar_current: "outscale-oks-cluster"
description: |-
  [Manages a cluster.]
---

# outscale_oks_cluster Resource

Manages a cluster.

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-OKS.html#_clusters).  
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/oks.html#oks-api-clusters).

## Example Usage

```hcl
resource "outscale_oks_project" "project01" {
  name   = "project01"
  cidr   = "10.50.0.0/18"
  region = "eu-west-2"
} 

resource "outscale_oks_cluster" "cluster01" {
  project_id      = outscale_oks_project.project01.id
  admin_whitelist = ["0.0.0.0/0"]
  cidr_pods       = "10.91.0.0/16"
  cidr_service    = "10.92.0.0/16"
  version         = "1.32"
  name            = "cluster01"
  control_planes  = "cp.mono.master"
  tags            = {
    tagkey = "tagvalue"
  }
}
```

## Argument Reference

The following arguments are supported:

* `admin_whitelist` - (Required) The list of CIDR blocks or IPs allowed to access the cluster via the Kubernetes API.
* `cidr_pods` - (Required) The CIDR block for Kubernetes pods' network.
* `cidr_service` - (Required) The CIDR block for the Kubernetes services' network.
* `name` - (Required) A unique name for the cluster within the project. Between 1 and 40 characters, this name must start with a letter and contain only lowercase letters, numbers, or hyphens.
* `project_id` - (Required) The ID of the project in which you want to create a cluster.
* `version` - (Required) The Kubernetes version to be deployed for the cluster. For more information, see [GetKubernetesVersions](https://docs.outscale.com/oks.html#getkubenetesversions).
* `admin_lbu` - (Optional) If true, load balancer administration is enabled for cluster management. If false, it is disabled. By default, false.
* `admission_flags` - (Optional) The configuration for Kubernetes admission controllers.
    * `disable_admission_plugins` - The list of Kubernetes admission plugins to disable.
    * `enable_admission_plugins` - The list of Kubernetes admission plugins to enable.
* `auto_maintenances` - (Optional) The configurations for automated maintenance windows.
    * `minor_upgrade_maintenance` - The maintenance window configuration for minor Kubernetes upgrades.
        * `duration_hours` - The duration of the maintenance window, in hours. By default, `0`.
        * `enabled` - If true, a maintenance window is enabled. By default, true.
        * `start_hour` - The starting time of the maintenance window, in hours. By default, `12`.
        * `tz` - The timezone for the maintenance window. By default, `UTC`.
        * `week_day` - The weekday on which the maintenance window begins (`Mon` \| `Tue` \| `Wed` \| `Thu` \| `Fri` \| `Sat` \| `Sun`). By default, `Tue`.
    * `patch_upgrade_maintenance`- The maintenance window configuration for patch Kubernetes upgrades.
        * `duration_hours` - The duration of the maintenance window, in hours. By default, `0`.
        * `enabled` - If true, a maintenance window is enabled. By default, true.
        * `start_hour` - The starting time of the maintenance window, in hours. By default, `12`.
        * `tz` - The timezone for the maintenance window. By default, `UTC`.
        * `week_day` - The weekday on which the maintenance window begins (`Mon` \| `Tue` \| `Wed` \| `Thu` \| `Fri` \| `Sat` \| `Sun`). By default, `Tue`.
* `cluster_dns` - (Optional) The IP for the cluster's DNS service.
* `control_planes` - (Optional) The size of control plane deployment for the cluster. For more information, see [About OKS > Control Planes](https://docs.outscale.com/en/userguide/About-OKS.html#_control_planes). By default, `cp.3.masters.small`.
* `cp_multi_az` - (Optional) If true, multi-Subregion deployment is enabled for the control plane. If false, it is disabled. By default, false.
* `cp_subregions` - (Optional) The list of Subregions where control plane components are deployed.
* `description` - (Optional) A description for the project.
* `disable_api_termination` - (Optional) If true, cluster deletion through the API is disabled. If false, it is enabled. By default, false.
* `quirks` - (Optional) The list of special configurations or behaviors for the cluster.
* `tags` - (Optional) 	The key/value combinations of the tags associated with the cluster's metadata.

## Attribute Reference

The following attributes are exported:

* `cni` - The Container Network Interface (CNI) used in the cluster.
* `id` - The Universally Unique Identifier (UUID) of the cluster.
* `admission_flags` - The configuration for Kubernetes admission controllers.
    * `applied_admission_plugins` - The list of admission plugins that are currently applied to the cluster.
    * `disable_admission_plugins_actual` - The list of Kubernetes admission plugins that are disabled.
    * `enable_admission_plugins_actual` - The list of Kubernetes admission plugins that are enabled.
* `kubeconfig` - A file containing access configuration to the cluster.
* `request_id` - The ID of the API request.
* `statuses` - The status information of the cluster.
    * `available_upgrade` - Any available version of Kubernetes for upgrade (if applicable). For more information, see [GetKubernetesVersions](https://docs.outscale.com/oks.html#getkubenetesversions).
    * `created_at` - The timestamp when the project was created (date-time).
    * `deleted_at` - The timestamp when the project was deleted (if applicable) (date-time).
    * `status` - The status of the project.
    * `updated_at` - The timestamp when the project was last updated (date-time).
* `auto_maintenances` - The configurations for automated maintenance windows.
    * `minor_upgrade_maintenance_actual` - The maintenance window configuration for minor Kubernetes upgrades.
    * `patch_upgrade_maintenance_actual` - The maintenance window configuration for minor Kubernetes upgrades.
* `disable_api_termination` - If true, cluster deletion through the API is disabled. If false, it is enabled.


## Import

An OKS project can be imported using its ID. For example:

```console

$ terraform import outscale_oks_cluster.cluster id

```