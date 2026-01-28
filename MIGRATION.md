Of course! Here's a clean, structured, user-friendly **MIGRATION.md** file you can directly use in your repo:

---

# 🚀 Migration Guide to v1.x

This guide will help you safely migrate your Terraform configuration and state to version **v1.x** of the **3DS OUTSCALE Terraform Provider**.

---

## ⚠️ Backup Recommendation

Before starting, **always back up your Terraform state and configuration files.**

- If you use a **local state**, copy your `terraform.tfstate` and `terraform.tfstate.backup` files.
- If you use a **remote backend** (like S3), follow your backend's backup procedure.
- Optionally, back up or version your Terraform `.tf` configuration files.

---

## 💥 Breaking Changes

- **Access Key Resource Behavior:**
  In versions **< v0.9.0**, creating an `access_key` **without an expiration date** is no longer supported.
  See [Issue #342](https://github.com/outscale/terraform-provider-outscale/issues/342) for more details.

---

## 📝 Migration Steps

### Step 1: Refresh Terraform State

To make sure the state is up-to-date with the current provider version, run:

```sh
terraform refresh
```

### Step 2: Update Provider Version

In your Terraform configuration:

```hcl
terraform {
  required_providers {
    outscale = {
      source  = "outscale/outscale"
      version = "1.3.2"
    }
  }
}

provider "outscale" {
  # Configuration options
}
```

Then run:

```sh
terraform init -upgrade
```

---

### Step 3: Clean Up Terraform State & Configuration

Some resource attributes and block names have changed in v1.x.
You must clean and adapt your state file and Terraform configuration.

#### On Linux

```sh
terraform fmt
sed -i 's/outscale_volumes_link/outscale_volume_link/g' terraform.tfstate
sed -i '/"block_device_mappings_created": \[/, /\],/d' terraform.tfstate
sed -i '/"source_security_group": {/, /},/d' terraform.tfstate
sed -i '/"flexible_gpu_id": "/, /",/d' terraform.tfstate
sed -i '/"link_public_ip": {/, /},/d' terraform.tfstate
sed -i '/"accepter_net": {/, /},/d' terraform.tfstate
sed -i '/"health_check": {/, /},/d' terraform.tfstate
sed -i '/"access_log": {/, /},/d' terraform.tfstate
sed -i '/"source_net": {/, /},/d' terraform.tfstate
sed -i '/"link_nic": {/, /},/d' terraform.tfstate
sed -i '/"state": {/, /},/d' terraform.tfstate
sed -i 's/outscale_volumes_link/outscale_volume_link/g' *.tf
sed -i 's/flexible_gpu_id /flexible_gpu_ids /g' *.tf
sed -i '/outscale_flexible_gpu\./s/$/ ]/' *.tf
sed -i '/flexible_gpu_ids /s/= /= [/' *.tf
terraform fmt
```

#### On macOS

```sh
terraform fmt
sed -i='' 's/outscale_volumes_link/outscale_volume_link/g' terraform.tfstate
sed -i='' '/"block_device_mappings_created": \[/, /\],/d' terraform.tfstate
sed -i='' '/"source_security_group": {/, /},/d' terraform.tfstate
sed -i='' '/"flexible_gpu_id": "/, /",/d' terraform.tfstate
sed -i='' '/"link_public_ip": {/, /},/d' terraform.tfstate
sed -i='' '/"accepter_net": {/, /},/d' terraform.tfstate
sed -i='' '/"health_check": {/, /},/d' terraform.tfstate
sed -i='' '/"access_log": {/, /},/d' terraform.tfstate
sed -i='' '/"source_net": {/, /},/d' terraform.tfstate
sed -i='' '/"link_nic": {/, /},/d' terraform.tfstate
sed -i='' '/"state": {/, /},/d' terraform.tfstate
sed -i='' 's/outscale_volumes_link/outscale_volume_link/g' *.tf
sed -i='' 's/flexible_gpu_id /flexible_gpu_ids /g' *.tf
sed -i='' '/outscale_flexible_gpu\./s/$/ ]/' *.tf
sed -i='' '/flexible_gpu_ids /s/= /= [/' *.tf
terraform fmt
```

---

### Step 3: Refresh Terraform State

Once the state and configuration files are updated:

```sh
terraform refresh
```

---

### ✅ Verify Configuration

You can now validate and plan:

```sh
terraform validate
terraform plan
```

---

## ℹ️ Need Help?

If you encounter issues during migration, check the open migration issues:
👉 [Migration Issues](https://github.com/outscale/terraform-provider-outscale/issues?q=is%3Aissue)

---

**Happy provisioning!** 🎉
