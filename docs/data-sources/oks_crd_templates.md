---
layout: "outscale"
page_title: "OUTSCALE: outscale_oks_crd_templates"
subcategory: "OKS API"
sidebar_current: "outscale-oks-crd-templates"
description: |-
  [Provides information about YAML manifest templates provided by the OKS API to create CRDs.]
---

# outscale_oks_crd_templates Data Source

Provides information about the YAML manifest templates provided by the OKS API to create Custom Resource Definitions (CRDs).

For more information on the templates, see the [API documentation](https://docs.outscale.com/oks.html#oks-api-templates).

## Example Usage

```hcl
data "outscale_oks_crd_templates" "templates" {
}
```

## Argument Reference

No argument is supported.

## Attribute Reference

The following attribute is exported:

* `manifests` - YAML manifest templates provided by the OKS API to create CRDs. For now, the following templates are exported: [`NodePool`](https://docs.outscale.com/oks.html#getnodepooltemplate), [`NetPeeringRequest`](https://docs.outscale.com/oks.html#getnetpeeringrequesttemplate), [`NetPeeringAcceptance`](https://docs.outscale.com/oks.html#getnetpeeringacceptancetemplate).
