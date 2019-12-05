---
layout: "outscale"
page_title: "Provider: 3DS OUTSCALE"
sidebar_current: "docs-outscale-index"
description: |-
  The 3DS OUTSCALE Services provider is used to interact with the many resources supported by 3DS OUTSCALE. The provider needs to be configured with the proper credentials before it can be used.
---

# 3DS OUTSCALE Provider

The 3DS OUTSCALE provider is used to interact with the
many resources supported by Outscale. The provider needs to be configured
with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
provider "outscale" {
  region     = "eu-west-2"
  access_key_id = "anaccesskey"
  secret_key_id = "asecretkey"
}
```
Available regions are: eu-west-2, us-east-2, us-west-1, cn-southeast-1


## Authentication

The 3DS OUTSCALE provider offers a flexible means of providing credentials for
authentication. The following methods are supported, in this order, and
explained below:

- access_key_id
- secret_key_id

### Static credentials ###

Static credentials can be provided by adding an `access_key_id` and `secret_key_id` in-line in the
Outscale provider block:

Usage:

```hcl
provider "outscale" {
  region     = "eu-west-2"
  access_key_id = "anaccesskey"
  secret_key_id = "asecretkey"
}
```

### Environment variables

You can provide your credentials via the `ACCESS_KEY_ID` and
`SECRET_KEY_ID`, environment variables, representing your 3DS OUTSCALE
Access Key and 3DS OUTSCALE Secret Key, respectively. 

```hcl
provider "outscale" {}
```

Usage:

```hcl
$ export ACCESS_KEY_ID="anaccesskey"
$ export SECRET_ACCESS_KEY="asecretkey"
$ export DEFAULT_REGION="eu-west-2"
$ terraform plan
```

### Shared Credentials file

You can use an 3DS OUTSCALE credentials file to specify your credentials.

Usage:

```hcl
provider "outscale" {
  region                  = "eu-west-2"
  shared_credentials_file = "/Users/tf_user/.outscale/creds"
  profile                 = "customprofile"
}
```

### 3DS OUTSCALE Role

The default deadline for the FCU metadata API endpoint is 100 milliseconds,
which can be overidden by setting the `OUTSCALE_METADATA_TIMEOUT` environment
variable. The variable expects a positive golang Time.Duration string, which is
a sequence of decimal numbers and a unit suffix; valid suffixes are `ns`
(nanoseconds), `us` (microseconds), `ms` (milliseconds), `s` (seconds), `m`
(minutes), and `h` (hours). Examples of valid inputs: `100ms`, `250ms`, `1s`,
`2.5s`, `2.5m`, `1m30s`.
