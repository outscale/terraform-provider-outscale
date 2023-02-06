---
page_title: "Provider: OUTSCALE"
---

# OUTSCALE Provider

The OUTSCALE provider is used to manage OUTSCALE Cloud resources.  
Use the navigation to the left to read about the available resources. For more information on our resources, see the [User Guide](https://docs.outscale.com/en/userguide/Home.html).

The provider is based on our OUTSCALE API. For more information, see [APIs Reference](https://docs.outscale.com/en/userguide/OUTSCALE-APIs-Reference.html) and the [API Documentation](https://docs.outscale.com/api).  

The provider needs to be configured with the proper credentials before it can be used.  

-> **Note:** 
To configure the provider, see our [README](https://github.com/outscale/terraform-provider-outscale#using-the-provider) and the [Terraform documentation](https://www.terraform.io/docs/configuration/provider-requirements.html). <br />
To configure a proxy, see our [README](https://github.com/outscale/terraform-provider-outscale#configuring-the-proxy-if-any).

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

OUTSCALE authentication is based on access keys composed of an **access key ID** and a **secret key**.
For more information on access keys, see [About Access Keys](https://docs.outscale.com/en/userguide/About-Access-Keys.html).
To retrieve your access keys, see [Getting Information About Your Access Keys](https://docs.outscale.com/en/userguide/Getting-Information-About-Your-Access-Keys.html).

The OUTSCALE provider supports different ways of providing credentials for authentication. The following methods are supported:

1. [Static credentials](#static-credentials)
2. [Environment variables](#environment-variables)

### Static credentials

!> Hard-coding credentials into any Terraform configuration is not recommended, and risks secret leakage should this file ever be committed to a public version control system.

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

```console
$ export OUTSCALE_ACCESSKEYID="myaccesskey"
$ export OUTSCALE_SECRETKEYID="mysecretkey"
$ export OUTSCALE_REGION="cloudgouv-eu-west-1"
$ export OUTSCALE_X509CERT="~/certificate/certificate.crt"
$ export OUTSCALE_X509KEY="~/certificate/certificate.key"

$ terraform plan
```

## Arguments Reference

In addition to [generic provider arguments](https://www.terraform.io/docs/configuration/providers.html), the following arguments are supported in the OUTSCALE provider block:

* `access_key_id` - (Optional) The ID of the OUTSCALE access key. It must be provided, but it can also be sourced from the `OUTSCALE_ACCESSKEYID` [environment variable](#environment-variables).

* `secret_key_id` - (Optional) The OUTSCALE secret key. It must be provided, but it can also be sourced from the `OUTSCALE_SECRETKEYID` [environment variable](#environment-variables).

* `region` - (Optional) The Region that will be used as default value for all resources. It can also be sourced from the `OUTSCALE_REGION` [environment variable](#environment-variables). For more information on available Regions, see [Regions, Endpoints and Availability Zones Reference](https://docs.outscale.com/en/userguide/Regions-Endpoints-and-Availability-Zones-Reference.html).

* `endpoints` - (Optional) The shortened custom endpoint that will be used as default value for all resources. For more information on available endpoints, see [Regions, Endpoints and Availability Zones Reference](https://docs.outscale.com/en/userguide/Regions-Endpoints-and-Availability-Zones-Reference.html).

* `x509_cert_path` - (Optional) The path to the x509 Client Certificate. It can also be sourced from the `OUTSCALE_X509CERT` [environment variable](#environment-variables). For more information on the use of those certificates, see [About API Access Rules](https://docs.outscale.com/en/userguide/About-API-Access-Rules.html).

* `x509_key_path` - (Optional) The path to the private key of the x509 Client Certificate. It can also be sourced from the `OUTSCALE_X509KEY` [environment variable](#environment-variables). For more information on the use of those certificates, see [About API Access Rules](https://docs.outscale.com/en/userguide/About-API-Access-Rules.html).