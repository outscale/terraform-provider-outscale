---
layout: "outscale"
page_title: "OUTSCALE: outscale_server_certificates"
sidebar_current: "docs-outscale-data-source-server-certificates"
description: |-
    Lists your server certificates. 
---

# outscale_server_certificates

Lists your server certificates.

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
  path = "/test_path/"
  private_key             = "${file("test-key.pem")}"
}

data "outscale_server_certificates" "test"{
    path_prefix = "${outscale_server_certificate.test_cert.path}"
}

```

## Argument Reference

The following arguments are supported:

* `path_prefix` - (Optional) The path prefix of the server certificates, set to a slash (/) if not specified.

## Attributes Reference

* `certificate_body` – The contents of the public key certificate in
  PEM-encoded format.
* `certificate_chain` – The PEM-encoded intermediate certification authorities.
* `server_certificate_metadata_list` - Information about one or more server certificates.

### Server Certificate Metadata List

The `server_certificate_metadata_list` has the following attributes:

* `server_certificate_id` - The unique Server Certificate name
* `server_certificate_name` - The name of the Server Certificate
* `arn` - The unique identifier of the server certificate (between 20 and 2048 characters), which can be used by EIM policies.
* `path` - The path to the server certificate.

[1]: http://docs.outscale.com/api_eim/index.html#_certificates
