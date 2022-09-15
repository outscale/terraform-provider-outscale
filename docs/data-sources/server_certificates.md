---
layout: "outscale"
page_title: "OUTSCALE: outscale_server_certificates"
sidebar_current: "outscale-server-certificates"
description: |-
  [Provides information about server certificates.]
---

# outscale_server_certificates Data Source

Provides information about server certificates.
For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Server-Certificates-in-EIM.html).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-servercertificate).

## Example Usage

### Read specific server certificates

```hcl
data "outscale_server_certificates" "server_certificates01" {
  filter {
    name   = "paths"
    values = ["<PATH01>", "<PATH02>"]
  }
}
```

### Read all server certificates

```hcl
data "outscale_server_certificates" "all_server_certificates" {
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) A combination of a filter name and one or more filter values. You can specify this argument for as many filter names as you need. The filter name can be any of the following:
    * `paths` - (Optional) The paths to the server certificates.

## Attribute Reference

The following attributes are exported:

* `server_certificates` - Information about one or more server certificates.
    * `expiration_date` - The date at which the server certificate expires.
    * `id` - The ID of the server certificate.
    * `name` - The name of the server certificate.
    * `path` - The path to the server certificate.
    * `upload_date` - The date at which the server certificate has been uploaded.
