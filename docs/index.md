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
  access_key_id = var.access_key_id
  secret_key_id = var.secret_key_id
  api {
    endpoint       = "https://api.eu-west-2.outscale.com"
    region         = "eu-west-2"
    x509_cert_path = "/path/to/cert.pem"
    x509_key_path  = "/path/to/key.pem"
  }
  oks {
    endpoint = "https://api.eu-west-2.oks.outscale.com/api/v2"
    region   = "eu-west-2"
  }
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
  access_key_id = "myaccesskey"
  secret_key_id = "mysecretkey"
  api {
    region = "eu-west-2"
  }
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

## Configuration

### Set a profile using a configuration file

You can set a named profile by specifying the `profile` attribute in the provider block.

The location of the shared configuration and credential file can be specified using the `config_file` attribute:

Example:

```hcl
provider "outscale" {
  config_file = "./.osc/config.json"
  profile     = "default"
}
```

### Set a profile using environment variables

You can also set a named profile by specifying the `OSC_PROFILE` environment variable.

The locations of the shared configuration and credential file can be specified using the `OSC_CONFIG_FILE` environment variable:

```hcl
# For Linux and macOS
export OSC_CONFIG_FILE="$HOME/.osc/config.json"
export OSC_PROFILE="default"
 
# For Windows
export OSC_CONFIG_FILE="%USERPROFILE%.osc\config.json"
 ```


## Arguments Reference

In addition to [generic provider arguments](https://www.terraform.io/docs/configuration/providers.html), the following arguments are supported in the OUTSCALE provider block:

* `config_file` - (Optional) The path to an OSC config file. It can also be sourced from the `OSC_CONFIG_FILE` [environment variable](#set-a-profile-using-environment-variables).
* `profile` - (Optional) The named profile you want to use in the OSC config file. It can also be sourced from the `OSC_PROFILE` [environment variable](#set-a-profile-using-environment-variables).
* `access_key_id` - (Optional) The ID of the OUTSCALE access key. It must be provided, but it can also be sourced from the `OUTSCALE_ACCESSKEYID` [environment variable](#environment-variables).
* `secret_key_id` - (Optional) The OUTSCALE secret key. It must be provided, but it can also be sourced from the `OUTSCALE_SECRETKEYID` [environment variable](#environment-variables).
* `api` - (Optional) Configuration elements for OUTSCALE API operations.
    * `endpoint` - (Optional) The endpoint to use for OUTSCALE API operations. For more information on available endpoints, see [API Endpoints Reference > OUTSCALE API](https://docs.outscale.com/en/userguide/API-Endpoints-Reference.html#_outscale_api).
    * `region` - (Optional) The Region to use for OUTSCALE API operations. It can also be sourced from the `OUTSCALE_REGION` [environment variable](#environment-variables). For more information on available Regions, see [About Regions and Subregions](https://docs.outscale.com/en/userguide/About-Regions-and-Subregions.html).
    * `x509_cert_path` - (Optional) The path to the x509 Client Certificate. It can also be sourced from the `OUTSCALE_X509CERT` [environment variable](#environment-variables). For more information on the use of those certificates, see [About API Access Rules](https://docs.outscale.com/en/userguide/About-API-Access-Rules.html).
    * `x509_key_path` - (Optional) The path to the private key of the x509 Client Certificate. It can also be sourced from the `OUTSCALE_X509KEY` [environment variable](#environment-variables). For more information on the use of those certificates, see [About API Access Rules](https://docs.outscale.com/en/userguide/About-API-Access-Rules.html).
    * `insecure` - (Optional) Enables TLS insecure connection.
* `oks` - (Optional) Configuration elements for OKS API operations.
    * `endpoint` - (Optional) The endpoint to use for OKS API operations. For more information on available endpoints, see [API Endpoints Reference > OUTSCALE Kubernetes as a Service (OKS)](https://docs.outscale.com/en/userguide/API-Endpoints-Reference.html#_outscale_kubernetes_as_a_service_oks).
    * `region` - (Optional) The Region to use for OKS API operations. It can also be sourced from the `OUTSCALE_REGION` [environment variable](#environment-variables). For more information on available Regions, see [About Regions and Subregions](https://docs.outscale.com/en/userguide/About-Regions-and-Subregions.html).

The following top-level arguments are deprecated but still supported as a fallback:
* `endpoints` - (Optional, deprecated) The endpoints to use for OUTSCALE API and OKS API operations. For more information on available endpoints, see [Regions, Endpoints and Availability Zones Reference](https://docs.outscale.com/en/userguide/Regions-Endpoints-and-Availability-Zones-Reference.html).
    * `api` - (Optional, deprecated) For OUTSCALE API.
    * `oks` - (Optional, deprecated) For OKS API.
* `region` - (Optional, deprecated) The Region to use for OUTSCALE API and OKS API operations. It can also be sourced from the `OUTSCALE_REGION` [environment variable](#environment-variables). For more information on available Regions, see [About Regions and Subregions](https://docs.outscale.com/en/userguide/About-Regions-and-Subregions.html).
* `x509_cert_path` - (Optional, deprecated) The path to the x509 Client Certificate. It can also be sourced from the `OUTSCALE_X509CERT` [environment variable](#environment-variables). For more information on the use of those certificates, see [About API Access Rules](https://docs.outscale.com/en/userguide/About-API-Access-Rules.html).
* `x509_key_path` - (Optional, deprecated) The path to the private key of the x509 Client Certificate. It can also be sourced from the `OUTSCALE_X509KEY` [environment variable](#environment-variables). For more information on the use of those certificates, see [About API Access Rules](https://docs.outscale.com/en/userguide/About-API-Access-Rules.html).
* `insecure` - (Optional, deprecated) Enables TLS insecure connection.
