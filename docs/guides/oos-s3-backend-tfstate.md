---
page_title: "Terraform Remote State on OOS (S3 Backend)"
description: "How to use OUTSCALE Object Storage (OOS) as an S3 backend for Terraform state (tfstate)."
---

# Terraform Remote State on OOS (S3 Backend)

This guide explains how to store Terraform state (`tfstate`) in [**OUTSCALE Object Storage (OOS)**](https://docs.outscale.com/en/userguide/Introduction-to-OOS.html) using Terraform's built-in [**S3 backend**](https://developer.hashicorp.com/terraform/language/backend/s3).

## Prerequisites

- Terraform installed (1.x recommended)
- An existing OOS bucket dedicated to Terraform state (see [Step 1](#1.-create-the-bucket-in-OOS-(one-time)))
- OOS credentials with read/write permissions on that bucket
- Network access to the OOS endpoint

## Conventions / placeholders

Replace these values with your own:

- `<OOS_REGION>`: e.g. `eu-west-2`
- `<OOS_ENDPOINT>`: `https://oos.<OOS_REGION>.outscale.com`  
  e.g.: `https://oos.eu-west-2.outscale.com`
- `<TFSTATE_BUCKET>`: e.g. `myproject-tfstate`
- `<TFSTATE_KEY>`: e.g. `terraform.tfstate` or `projects/<project>/terraform.tfstate`

## 1. Create the bucket in OOS (one-time)

The bucket **must exist** before Terraform can use it as a backend.

### Option A: Create the bucket using [OCTL](https://github.com/outscale/octl)

```bash
octl storage bucket create --bucket <TFSTATE_BUCKET>
```

### Option B: Create the bucket using `s3cmd`

See OUTSCALE documentation: https://docs.outscale.com/en/userguide/s3cmd.html

1. Configure `s3cmd`:

```bash
s3cmd --configure
```

2. Create the bucket:

```bash
s3cmd mb s3://<TFSTATE_BUCKET>
```


## 2. Provide credentials for Terraform backend

Terraform's `s3` backend requires S3 credentials. Two common approaches are supported.

### Option A: Environment variables (recommended for CI)

```bash
export AWS_ACCESS_KEY_ID="<OOS_ACCESS_KEY>"
export AWS_SECRET_ACCESS_KEY="<OOS_SECRET_KEY>"
```

### Option B: OOS shared credentials file (recommended for local dev)

Create or update `~/.osc/credential`:

```ini
[<PROFILE>]
aws_access_key_id = <OOS_ACCESS_KEY>
aws_secret_access_key = <OOS_SECRET_KEY>
```

## 3. Required: Set AWS checksum compatibility env vars

To ensure compatibility with OOS, you must set these environment variables before running Terraform:

```bash
export AWS_REQUEST_CHECKSUM_CALCULATION=WHEN_REQUIRED
export AWS_RESPONSE_CHECKSUM_VALIDATION=WHEN_REQUIRED
```

Reference (OUTSCALE documentation):
[AWS SDK and CLI Compatibility Warning](https://docs.outscale.com/en/userguide/AWS-SDK-and-CLI-Compatibility-Warning.html)

## 4. Configure the Terraform S3 backend for OOS

Create a dedicated file (recommended) like `backend.tf`:

```hcl
terraform {
  backend "s3" {
    bucket = "<TFSTATE_BUCKET>"
    key    = "<TFSTATE_KEY>"
    region = "<OOS_REGION>"

    endpoints = {
      s3 = "<OOS_ENDPOINT>"
    }

    # Credential source:
    # - keep it if you rely on ~/.osc/credential
    # - remove it if you use AWS_ACCESS_KEY_ID / AWS_SECRET_ACCESS_KEY
    profile                  = "<PROFILE>"
    shared_credentials_files = ["~/.osc/credential"]

    # S3-compatible / non-AWS settings
    skip_credentials_validation = true
    skip_region_validation      = true
    skip_requesting_account_id  = true
    skip_s3_checksum            = true

    # Optional: enable lock file (recommended if your workflow supports it)
    # use_lockfile = true
  }
}
```

### Working example (as-is)

```hcl
terraform {
  required_providers {
    outscale = {
      source  = "outscale/outscale"
      version = "1.4.0"
    }
  }

  backend "s3" {
    bucket = "myproject-tfstate"
    key     = "terraform.tfstate"
    region  = "eu-west-2"

    endpoints = {
      s3 = "https://oos.eu-west-2.outscale.com"
    }
    
    profile = "default"
    shared_credentials_files = ["./osc/credential"]
    
    skip_credentials_validation = true
    skip_region_validation      = true
    skip_requesting_account_id  = true
    skip_s3_checksum            = true

    use_lockfile = true
  }
}
```

## 5. Initialize Terraform (and migrate state if needed)

Run:

```bash
terraform init
```

If a local state already exists and no remote state is found in OOS, Terraform may prompt:

```text
Initializing the backend...
Do you want to copy existing state to the new backend?
  Pre-existing state was found while migrating the previous "local" backend to the
  newly configured "s3" backend. No existing state was found in the newly
  configured "s3" backend. Do you want to copy this state to the new "s3"
  backend? Enter "yes" to copy and "no" to start with an empty state.

  Enter a value: yes
```

After completion you should see:

```text
Successfully configured the backend "s3"! Terraform will automatically
use this backend unless the backend configuration changes.
```

If you change backend configuration later (bucket/key/endpoint/etc.), re-run init with:

```bash
terraform init -reconfigure
```

## 6. Verify the state is stored in OOS

After deploying your Terraform resources, you can see the object corresponding to `<TFSTATE_KEY>` in your bucket.

### Option A: Using OCTL

```bash
octl storage object list --bucket <TFSTATE_BUCKET>
```

### Option B: Using `s3cmd`
```bash
s3cmd ls s3://<TFSTATE_BUCKET>
```

## 7. Recommendations (security and recoverability)

Terraform state can contain sensitive values. Treat the bucket as sensitive:

* Keep the bucket private
* Restrict access to only Terraform users/CI runners that need it
* Enable bucket versioning to help recover from accidental overwrites or deletions

### Enable versioning (recommended)

```bash
octl storage api PutBucketVersioning \
  --Bucket <TFSTATE_BUCKET> \
  --VersioningConfiguration.Status Enabled \
```

## 8. Troubleshooting

### Error: checksum / "trailing checksum is not supported" / BadRequest on PutObject

Make sure these are exported (local shell and CI environment):

```bash
export AWS_REQUEST_CHECKSUM_CALCULATION=WHEN_REQUIRED
export AWS_RESPONSE_CHECKSUM_VALIDATION=WHEN_REQUIRED
```

Reference:
[AWS SDK and CLI Compatibility Warning](https://docs.outscale.com/en/userguide/AWS-SDK-and-CLI-Compatibility-Warning.html)

### Error: "No valid credential sources found"

* Ensure `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` are set, or
* Ensure `~/.osc/credential` has the profile you configured in the backend block, and `profile = "..."` matches.

### Error: region/account/STS validation failures

Keep the following in the backend configuration:

* `skip_credentials_validation = true`
* `skip_region_validation = true`
* `skip_requesting_account_id = true`

### Error: Unable to create resource / dial tcp: lookup api..outscale.com: no such host

In addition to the S3 credentials for OOS, make sure you configure [your credentials for the OUTSCALE provider](../index.md) itself too.