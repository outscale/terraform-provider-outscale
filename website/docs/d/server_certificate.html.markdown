---
layout: "outscale"
page_title: "OUTSCALE: outscale_server_certificate"
sidebar_current: "docs-outscale-data-source-server-certificate"
description: |-
  Retrieves an EIM Server Certificate
---

# outscale_server_certificate

Retrieves an EIM Server Certificate resource to upload Server Certificates.

These elements can be used with other services (for example, to configure SSL termination on load balancers).
You can also specify the chain of intermediate certification authorities if your certificate is not directly signed by a root one. You can specify multiple intermediate certification authorities in the `CertificateChain` parameter. To do so, concatenate all certificates in the correct order (the first certificate must be the authority of your certificate, the second must the the authority of the first one, and so on).

The private key must be a RSA key in PKCS1 form. To check this, open the PEM file and ensure its header reads as follows: **BEGIN RSA PRIVATE KEY**.

See detailed information in: [OUTSCALE Certificates Documentation](1)

## Example Usage

**Using certs on file:**

```hcl
resource "outscale_server_certificate" "test_cert" {
  server_certificate_name = "some_test_cert"
  certificate_body        = "${file("self-ca-cert.pem")}"
  private_key             = "${file("test-key.pem")}"
}

data "outscale_server_certificate" "test"{
    server_certificate_name = "${outscale_server_certificate.test_cert.server_certificate_name}"
}

```

## Argument Reference

The following arguments are supported:

* `server_certificate_name` - (Required) The name of the certificate, which must be unique. Do not include the
  path in this value.

## Attributes Reference

* `certificate_body` – The contents of the public key certificate in
  PEM-encoded format.
* `certificate_chain` – The PEM-encoded intermediate certification authorities.
* `server_certificate_metadata` - The metadata of the server certificate.

### Server Certificate Metadata

The `server_certificate_metadata` has the following attributes:

* `server_certificate_id` - The unique Server Certificate name
* `server_certificate_name` - The name of the Server Certificate
* `arn` - The unique identifier of the server certificate (between 20 and 2048 characters), which can be used by EIM policies.
* `path` - The path to the server certificate.


[1]: http://docs.outscale.com/api_eim/index.html#_certificates
