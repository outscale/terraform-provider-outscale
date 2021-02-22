---
layout: "outscale"
page_title: "Provider: 3DS OUTSCALE"
description: |-
  The 3DS OUTSCALE provider is used to manage 3DS OUTSCALE Cloud resources. The provider needs to be configured with the proper credentials before it can be used.
---

# 3DS OUTSCALE Provider

The 3DS OUTSCALE provider is used to manage 3DS OUTSCALE Cloud resources.  
Use the navigation to the left to read about the available resources. For more information on our resources, see the [User Guide](https://wiki.outscale.net/display/EN#).

The provider is based on our 3DS OUTSCALE API. For more information, see [APIs Reference](https://wiki.outscale.net/display/EN/3DS+OUTSCALE+APIs+Reference) and the [API Documentation](https://docs.outscale.com/api#3ds-outscale-api).  

The provider needs to be configured with the proper credentials before it can be used.  

-> **Note:** Since the release of Terraform 0.13, provider declaration has changed. For more information, see our [README](https://github.com/outscale-dev/terraform-provider-outscale#using-the-provider) and the [Terraform documentation](https://www.terraform.io/docs/configuration/provider-requirements.html).


## Example

```hcl
provider "outscale" {
  access_key_id  = var.access_key_id
  secret_key_id  = var.secret_key_id
  region         = "eu-west-2"
  endpoints {
    api  = "api.eu-west-2.outscale.com"
    }
  x509_cert_path = "/tmp/client-certificate.pem"
  x509_key_path  = "/tmp/key.pem"
}
```

## Authentication

3DS OUTSCALE authentication is based on access keys composed of an **access key ID** and a **secret key**.
For more information on access keys, see [About Access Keys](https://wiki.outscale.net/display/EN/About+Access+Keys).
To retrieve your access keys, see [Getting Information About Your Access Keys](https://wiki.outscale.net/display/EN/Getting+Information+About+Your+Access+Keys).

The 3DS OUTSCALE provider supports different ways of providing credentials for authentication. The following methods are supported:

- [3DS OUTSCALE Provider](#3ds-outscale-provider)
  - [Example](#example)
  - [Authentication](#authentication)
    - [Static credentials](#static-credentials)
    - [Environment variables](#environment-variables)
  - [Arguments Reference](#arguments-reference)

### Static credentials

!> **Warning:** Hard-coding credentials into any Terraform configuration is not recommended, and risks secret leakage should this file ever be committed to a public version control system.

You can provide your credentials by specifying the `access_key_id` and `secret_key_id` attributes in the provider block:

Example:

```hcl
provider "outscale" {
  access_key_id   = "myaccesskey"
  secret_key_id   = "mysecretkey"
  region          = "eu-west-2"
}
```

### Environment variables

You can provide your credentials with the `OUTSCALE_ACCESSKEYID` and `OUTSCALE_SECRETKEYID` environment variables:

Example:

```hcl
provider "outscale" {}
```

Usage:

```bash
$ export OUTSCALE_ACCESSKEYID="myaccesskey"
$ export OUTSCALE_SECRETKEYID="mysecretkey"
$ export OUTSCALE_REGION="cloudgouv-eu-west-1"
$ export OUTSCALE_X509CERT="~/certificate/certificate.crt"
$ export OUTSCALE_X509KEY="~/certificate/certificate.key"

$ terraform plan
```

## Arguments Reference

In addition to [generic provider arguments](https://www.terraform.io/docs/configuration/providers.html), the following arguments are supported in the 3DS OUTSCALE provider block:

* `access_key_id` - (Optional) The ID of the 3DS OUTSCALE access key. It must be provided, but it can also be sourced from the `OUTSCALE_ACCESSKEYID` [environment variable](#environment-variables).

* `secret_key_id` - (Optional) The 3DS OUTSCALE secret key. It must be provided, but it can also be sourced from the `OUTSCALE_SECRETKEYID` [environment variable](#environment-variables).

* `region` - (Optional) The Region that will be used as default value for all resources. It can also be sourced from the `OUTSCALE_REGION` [environment variable](#environment-variables). For more information on available Regions, see [Regions Reference](https://wiki.outscale.net/display/EN/Regions%2C+Endpoints+and+Availability+Zones+Reference).

* `endpoints` - (Optional) The shortened custom endpoint that will be used as default value for all resources. For more information on available endpoints, see [Endpoints Reference](https://wiki.outscale.net/display/EN/Regions%2C+Endpoints+and+Availability+Zones+Reference).

* `x509_cert_path` - (Optional) The path to the x509 Client Certificate. It can also be sourced from the `OUTSCALE_X509CERT` [environment variable](#environment-variables). For more information on the use of those certificates, see [About API Access Rules](https://wiki.outscale.net/display/EN/About+API+Access+Rules).

* `x509_key_path` - (Optional) The path to the private key of the x509 Client Certificate. It can also be sourced from the `OUTSCALE_X509KEY` [environment variable](#environment-variables). For more information on the use of those certificates, see [About API Access Rules](https://wiki.outscale.net/display/EN/About+API+Access+Rules).