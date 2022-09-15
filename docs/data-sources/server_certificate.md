---
layout: "outscale"
page_title: "OUTSCALE: outscale_server_certificate"
sidebar_current: "outscale-server-certificate"
description: |-
  [Provides information about a specific server certificate.]
---

# outscale_server_certificate Data Source

Provides information about a specific server certificate.
For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Server-Certificates-in-EIM.html).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-servercertificate).

## Example Usage

```hcl
data "outscale_server_certificate" "server_certificate01" {
  filter {
    name   = "paths"
    values = "<PATH>"
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) A combination of a filter name and one or more filter values. You can specify this argument for as many filter names as you need. The filter name can be any of the following:
    * `paths` - (Optional) The paths to the server certificates.

## Attribute Reference

The following attributes are exported:

* `expiration_date` - The date at which the server certificate expires.
* `id` - The ID of the server certificate.
* `name` - The name of the server certificate.
* `path` - The path to the server certificate.
* `upload_date` - The date at which the server certificate has been uploaded.
