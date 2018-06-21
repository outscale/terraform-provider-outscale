---
layout: "outscale"
page_title: "OUTSCALE: outscale_load_balancer_ssl_certificate"
sidebar_current: "docs-outscale-resource-load-balancer-ssl-certificate"
description: |-
  Creates a load balancer ssl certificate.
---

# outscale_load_balancer_ssl_certificate

Sets a new SSL certificate to an SSL or HTTPS listener of a load balancer.
This certificate replaces any certificate used on the same load balancer and port. To do so, you first need to upload a certificate using the [UploadServerCertificate](http://docs.outscale.com/api_eim/index.html#_certificates) method in the Elastic Identity Management (EIM) API.

## Example Usage

```hcl
resource "outscale_load_balancer" "bar" {
  availability_zones = ["eu-west-2a"]
  load_balancer_name = "foobar-terraform-lbu-test"
  listeners {
    instance_port = 8000
    instance_protocol = "HTTP"
    load_balancer_port = 80
    protocol = "HTTP"
  }

    tag {
        bar = "baz"
    }

}

resource "outscale_server_certificate" "test_cert" {
  server_certificate_name = "terraform-test-cert"
  certificate_body = "${file("filepath/cert.pem")}"
  private_key =  <<EOF
-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQDKdH6BU9Q0xBVPfeX5NjCC/B2Pm3WsFGnTtRw4abkD+r4to9wD
eYUgjH2yPCyonNOA8mNiCQgDTtaLfbA8LjBYoodt7rgaTO7C0ugRtmTNK96DmYxm
f8Gs5ZS6eC3yeaFv58d1w2mow7tv0+DRk8uXwzVfaaMxoalsCtlLznmZHwIDAQAB
AoGABZj69nBu6ZaSUERW23EYHkcCOjo+Iqfd1TCouxaROv7vyytApgfyGlhIEWmA
gpjzcBlDji5Zvl2rqOesu707MOuJavZvluo+JHy/VIuU+yGUrWuO/QVCu6Jn3yns
vS7g48ConuZ962cTzRPcpPDspONBVOAhVCF33Y8PsnxV0wECQQD5RqeoqxEUupsy
QhrDui0KkYXLdT0uhrEQ69n9rvAiQoHPsiX0MswfEKnj/g9N3VwGLdgWytT0TvcI
8fDPRB4/AkEAz+qF3taX77gB69XRPQwCGWqE1fHIFMwX7QeYdEsk3iRZ0EKVcdp6
vIPCB2Cq4a4eXcaFa/bXen4yeYgyTbeNIQJBAO92dWctdoowPRiJskZmGhC1/Q6X
gH+qenyj5VSy8hInS6anH5i4F6icDGhtzmvhgx6YeaZjkTFkjiG0sb2aVWcCQQDD
WL7UwtzX/xPXB/ril5C1Xo5WESgC2ks0ielkgmGuUYsNEDInWbXtvwGjOuDyz0x6
oRYkfTSxQzabVyqkOGvhAkBtbjUxOD8wgBIjb4T6mAMokQo6PeEAZGUTyPifjJNo
detWVr2WRvgNgQvcRnNPECwfq1RtMJJpavaI3kgeaSxg
-----END RSA PRIVATE KEY-----
EOF
}

resource "outscale_load_balancer_ssl_certificate" "test" {
    load_balancer_name = "${outscale_load_balancer.bar.id}"
    load_balancer_port = "${outscale_load_balancer.bar.listeners.0.load_balancer_port}"
    ssl_certificate_id = "${outscale_server_certificate.test_cert.id}"
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_name` -  (Required) The name of the load balancer.
* `load_balancer_port` - (Required) The port using SSL certificate.
* `ssl_certificate_id` - (Required) The Outscale Resource Name (ORN) of the SSL .certificate.

## Attributes

* `load_balancer_name` - The name of the load balancer.
* `load_balancer_port` - The port using SSL certificate.
* `ssl_certificate_id` - The Outscale Resource Name (ORN) of the SSL.
* `request_id` - The ID of the request.


[See detailed information.](http://docs.outscale.com/api_lbu/operations/Action_CreateLoadBalancer_get.html#_api_lbu-action_createloadbalancer_get)
