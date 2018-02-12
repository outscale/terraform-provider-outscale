---
layout: "outscale"
page_title: "Provider: OUTSCALE"
sidebar_current: "docs-outscale-index"
description: |-
  The Outscale Services provider is used to interact with the many resources supported by Outscale. The provider needs to be configured with the proper credentials before it can be used.
---

# Outscale Provider

The Outscale provider is used to interact with the
many resources supported by Outscale. The provider needs to be configured
with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
# Configure the Outscale Provider
provider "outscale" {
  access_key = "${var.outscale_access_key}"
  secret_key = "${var.outscale_secret_key}"
  region     = "us-east-1"
}

# Create a web server
resource "outscale_vm" "web" {
  # ...
}
```

## Authentication

The Outscale provider offers a flexible means of providing credentials for
authentication. The following methods are supported, in this order, and
explained below:

- access_key_id
- secret_key_id

### Static credentials ###

Static credentials can be provided by adding an `access_key` and `secret_key` in-line in the
Outscale provider block:

Usage:

```hcl
provider "outscale" {
  region     = "us-west-2"
  access_key_id = "anaccesskey"
  secret_key_id = "asecretkey"
}
```

### Environment variables

You can provide your credentials via the `ACCESS_KEY_ID` and
`SECRET_KEY_ID`, environment variables, representing your AWS
Access Key and AWS Secret Key, respectively.  Note that setting your
AWS credentials using either these (or legacy) environment variables
will override the use of `AWS_SHARED_CREDENTIALS_FILE` and `AWS_PROFILE`.
The `AWS_DEFAULT_REGION` and `AWS_SESSION_TOKEN` environment variables
are also used, if applicable:

```hcl
provider "aws" {}
```

Usage:

```hcl
$ export ACCESS_KEY_ID="anaccesskey"
$ export SECRET_ACCESS_KEY="asecretkey"
$ export DEFAULT_REGION="us-west-2"
$ terraform plan
```

### Shared Credentials file

You can use an Outscale credentials file to specify your credentials.

Usage:

```hcl
provider "aws" {
  region                  = "us-west-2"
  shared_credentials_file = "/Users/tf_user/.aws/creds"
  profile                 = "customprofile"
}
```

### Outscale Role

The default deadline for the EC2 metadata API endpoint is 100 milliseconds,
which can be overidden by setting the `OUTSCALE_METADATA_TIMEOUT` environment
variable. The variable expects a positive golang Time.Duration string, which is
a sequence of decimal numbers and a unit suffix; valid suffixes are `ns`
(nanoseconds), `us` (microseconds), `ms` (milliseconds), `s` (seconds), `m`
(minutes), and `h` (hours). Examples of valid inputs: `100ms`, `250ms`, `1s`,
`2.5s`, `2.5m`, `1m30s`.


## Argument Reference

The following arguments are supported:

* `block_device_mapping` - (Optional) The block device mapping of the instance.
* `client_token` - (Optional) A unique identifier which enables you to manage the idempotency.
* `disable_api_termination` - (Optional) If true, you cannot terminate the instance using Cockpit, the CLI or the API. If false, you can.
* `dry_run` - (Optional) If true, checks whether you have the required permissions to perform the action.
* `ebs_optimized` - (Optional) If true, the instance is created with optimized BSU I/O. All Outscale instances have optimized BSU I/O.
* `image_id` - (Required) The ID of the OMI. You can find the list of OMIs by calling the DescribeImages method.
* `instance_initiated_shutdown_behavior` - (Optional) The instance behavior when you stop or terminate it. By default or if set to stop, the instance stops. If set to restart, the instance stops then automatically restarts. If set to terminate, the instance stops and is terminated.
* `instance_type` - (Optional) The type of instance. For more information, see Instance Types.
* `key_name` - (Optional) The name of the keypair.
* `max_count` - (Required) The maximum number of instances you want to launch. If all the instances cannot be created, the largest possible number of instances above MinCount are created and launched.
* `min_count` - (Required) The minimum number of instances you want to launch. If this number of instances cannot be created, FCU does not create and launch any instance.
* `network_interface` - (Optional) One or more network interfaces.
* `placement` - (Optional) A specific placement where you want to create the instances (for example, Availability Zone, dedicated host, affinity criteria and so on).
* `private_ip_address` - (Optional) In a VPC, the unique primary IP address. The IP address must come from the IP address range of the subnet.
* `private_ip_addresses` - (Optional) In a VPC, the list of primary IP addresses when you create several instances. The IP addresses must come from the IP address range of the subnet.
* `security_group` - (Optional) One or more security group names.
* `security_group_id` - (Optional) One or more security group IDs.
* `subnet_id` - (Optional) In a VPC, the ID of the subnet in which you want to launch the instance.
* `user_data` - (Optional) Data or a script used to add a specific configuration to the instance when launching it. If you are not using a command line tool, this must be base64-encoded.
