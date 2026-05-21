---
layout: "outscale"
page_title: "OUTSCALE: outscale_oks_kubeconfig"
subcategory: "OKS API"
sidebar_current: "outscale-oks-kubeconfig"
description: |-
  [Provides information about a kubeconfig file.]
---

# outscale_oks_kubeconfig Data Source

Provides information about a kubeconfig file.

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/Accessing-a-Cluster.html).  
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/oks.html#getkubeconfigwithpubkeynacl).

## Example Usage

```hcl
data "outscale_oks_kubeconfig" "config" {
  cluster_id = "00000000-0000-4000-8000-000000000000"
} 
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required) The ID of the cluster.
* `user` - (Optional) The user of the kubeconfig file.
* `group` - (Optional) The group of the kubeconfig file.
* `ttl` - (Optional) The time to live (TTL) of the kubeconfig file.
* `x-encrypt-nacl` - (Optional) The header to encrypt the kubeconfig file.

## Attribute Reference

The following attributes are exported:

* `kubeconfig`  - (Sensitive value) A file containing access configuration to the cluster.
* `kubeconfig_attributes` - (Sensitive value) Access configuration to the cluster.
    * `client_certificate` - The client certificate for accessing the cluster.
    * `client_key` - The private key matching the client certificate.
    * `cluster_ca_certificate` - The Certificate Authority (CA) associated with the cluster.
    * `host` - The URL of the Kubernetes API server of the cluster.
