---
layout: "outscale"
page_title: "OUTSCALE: outscale_server_certificate"
sidebar_current: "outscale-server-certificate"
description: |-
  [Manages a server certificate.]
---

# outscale_server_certificate Resource

Manages a server certificate.
For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Server-Certificates-in-EIM.html).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-servercertificate).

## Example Usage

```hcl
resource "outscale_server_certificate" "server_certificate_01" { 
    name        =  "terraform-server-certificate"
    body        =  file("<PATH>")
    chain       =  file("<PATH>")
    private_key =  file("<PATH>")
    path        =  "<PATH>"
}
```


## Argument Reference

The following arguments are supported:

* `body` - (Required) The PEM-encoded X509 certificate.
* `chain` - (Optional) The PEM-encoded intermediate certification authorities.
* `name` - (Required) A unique name for the certificate. Constraints: 1-128 alphanumeric characters, pluses (+), equals (=), commas (,), periods (.), at signs (@), minuses (-), or underscores (_).
* `path` - (Optional) The path to the server certificate, set to a slash (/) if not specified.
* `private_key` - (Required) The PEM-encoded private key matching the certificate.

## Attribute Reference

The following attributes are exported:

* `expiration_date` - The date at which the server certificate expires.
* `id` - The ID of the server certificate.
* `name` - The name of the server certificate.
* `path` - The path to the server certificate.
* `upload_date` - The date at which the server certificate has been uploaded.

## Import

A server certificate can be imported using its ID. For example:

```console

$ terraform import outscale_server_certificate.ImportedServerCertificate 0123456789

```