---
layout: "outscale"
page_title: "Provider: 3DS OUTSCALE"
description: |-
  The 3DS OUTSCALE provider is used to manage 3DS OUTSCALE Cloud resources. The provider needs to be configured with the proper credentials before it can be used.
---

# 3DS OUTSCALE Provider

The 3DS OUTSCALE provider is used to manage 3DS OUTSCALE Cloud resources.
Use the navigation to the left to read about the available resources.
For more information on our resources, see the [User Guide](https://wiki.outscale.net/display/EN#).

The provider is based on our 3DS OUTSCALE API. For more information on our APIs, see [APIs Reference](https://wiki.outscale.net/display/EN/3DS+OUTSCALE+APIs+Reference).

The provider needs to be configured with the proper credentials before it can be used.

## Example

```hcl
provider "outscale" {
  access_key_id = "AZERTY123456QSDF7890"
  secret_key_id = "123456AZERTY7890QSDFAZERTY123456QSDF7890"
  region        = "eu-west-2"
}
```

## Authentication

3DS OUTSCALE authentication is based on access keys composed of an **access key ID** and a **secret key**.
For more information on access keys, see [About Access Keys](https://wiki.outscale.net/display/EN/About+Access+Keys).
To retrieve your access keys, see [Getting Information About Your Access Keys](https://wiki.outscale.net/display/EN/Getting+Information+About+Your+Access+Keys).

To provide your credentials to Terraform, you need to specify the `access_key_id` and `secret_key_id` attributes in your configuration file.
The 3DS OUTSCALE provider offers several ways to specify these attributes. The following methods are supported:

1. [Static credentials](#static-credentials)
2. [Environment variables](#environment-variables)

### Static credentials

!> **Warning**: Hard-coding credentials into any Terraform configuration is not recommended, and risks secret leakage should this file ever be committed to a public version control system.

In the provider block of your configuration file, you can provide your credentials with raw values:

Example:

```hcl
provider "outscale" {
  access_key_id   = "myaccesskey"
  secret_key_id   = "mysecretkey"
  region          = "regionname"
}
```

### Environment variables

In the provider block of your configuration file, you can provide your credentials with the `OSC_ACCESS_KEY_ID`and `OSC_SECRET_ACCESS_KEY` environment variables:

Example:

```hcl
provider "outscale" {
	access_key_id   = "var.access_key_id"
  secret_key_id = "var.secret_key_id"
  region        = "var.region"
}
```

Usage:

```bash
$ export OSC_ACCESS_KEY_ID="myaccesskey"
$ export OSC_SECRET_ACCESS_KEY="mysecretkey"
$ export OSC_DEFAULT_REGION="regionname"
$ terraform plan
```

## Arguments Reference

In addition to [generic provider arguments](https://www.terraform.io/docs/configuration/providers.html), the following arguments are supported in the 3DS OUTSCALE provider block:

- `access_key_id` - (Optional) The ID of the 3DS OUTSCALE access key. It must be provided, but it can also be sourced from the `OSC_ACCESS_KEY_ID` [environment variable](#environment-variables).

- `secret_key_id` - (Optional) The 3DS OUTSCALE secret key. It must be provided, but it can also be sourced from the `OSC_SECRET_ACCESS_KEY` [environment variable](#environment-variables).

- `region` - (Optional) The Region that will be used as default value for all resources. It can also be sourced from the `OSC_DEFAULT_REGION` [environment variable](#environment-variables). For more information on available Regions, see [Regions Reference](https://wiki.outscale.net/display/EN/Regions%2C+Endpoints+and+Availability+Zones+Reference