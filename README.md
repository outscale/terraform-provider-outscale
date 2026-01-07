# 3DS OUTSCALE Terraform Provider

[![Project Graduated](https://docs.outscale.com/fr/userguide/_images/Project-Graduated-green.svg)](https://docs.outscale.com/en/userguide/Open-Source-Projects.html)

<p align="center">
  <img alt="Terraform" src="https://www.datocms-assets.com/2885/1731373310-terraform_white.svg" width="200px">
</p>

## 🌐 Links
- Terraform website: https://www.terraform.io
- [Gitter chat](https://gitter.im/hashicorp-terraform/Lobby)
- [Mailing list](http://groups.google.com/group/terraform-tool)

---

## 📄 Table of Contents
- [Requirements](#-requirements)
- [Migration to v1](#-migration-to-v1)
- [Breaking Changes](#-breaking-changes)
- [Using the Provider](#-using-the-provider)
  - [With Terraform](#with-terraform)
  - [With OpenTofu](#with-opentofu)
- [Proxy Configuration](#-proxy-configuration)
- [x509 Authentication](#-x509-authentication)
- [Building the Provider](#-building-the-provider)
- [Using the Built Provider](#-using-the-built-provider)
- [Contributing](#-contributing)
- [Building the Documentation](#-building-the-documentation)

---

## ✅ Requirements
- [Terraform 1.10.x or latest](https://www.terraform.io/downloads.html)
- [Go 1.24.2](https://golang.org/doc/install) (to build the provider)

---

## 🚀 Migration to v1

> ⚠️ **Warning:** Always backup your state file before migrating!

See [MIGRATION GUIDE](./MIGRATION.md) for full instructions.

<details>
<summary>Migration Steps</summary>

### Step 1: Upgrade provider version
```hcl
terraform {
  required_providers {
    outscale = {
      source  = "outscale/outscale"
      version = "1.3.1"
    }
  }
}

provider "outscale" {
  # Configuration
}
```
```sh
terraform init -upgrade
```

### Step 2: Clean up state & configuration

**Linux**
```sh
sed -i 's/outscale_volumes_link/outscale_volume_link/g' terraform.tfstate
# + Other sed commands
```

**macOS**
```sh
sed -i='' 's/outscale_volumes_link/outscale_volume_link/g' terraform.tfstate
# + Other sed commands
```

### Step 3: Refresh
```sh
terraform refresh
```
</details>

---

## 💥 Breaking Changes

> ⚠️ **Important:**
There is a breaking change when creating an `access_key` without expiration date in versions `< v0.9.0`.
See [Issue #342](https://github.com/outscale/terraform-provider-outscale/issues/342).

---

## 🚀 Using the Provider

### With Terraform

```hcl
terraform {
  required_providers {
    outscale = {
      source  = "outscale/outscale"
      version = "1.3.1"
    }
  }
}

provider "outscale" {
  # Configuration options
}
```
```sh
terraform init
terraform plan
```

### With OpenTofu
```hcl
terraform {
  required_providers {
    outscale = {
      source  = "outscale/outscale"
      version = "1.3.1"
    }
  }
}

provider "outscale" {
  # Configuration options
}
```
```sh
tofu init
tofu plan
```

👉 See [OpenTofu migration guide](https://opentofu.org/docs/intro/migration/).

---

## 🌍 Proxy Configuration

**Linux/macOS**
```sh
export HTTPS_PROXY=http://192.168.1.24:3128
```

**Windows**
```cmd
set HTTPS_PROXY=http://192.168.1.24:3128
```

---

## 🔐 x509 Authentication

Add to your provider config:
```hcl
provider "outscale" {
  x509_cert_path = "/myrepository/certificate/client_ca.crt"
  x509_key_path  = "/myrepository/certificate/client_ca.key"
}
```
Or set environment variables:
```sh
export OUTSCALE_X509CERT=/myrepository/certificate/client_ca.crt
export OUTSCALE_X509KEY=/myrepository/certificate/client_ca.key
```

---

## 🛠 Building the Provider

Clone and build:
```sh
git clone --branch v1.3.1 https://github.com/outscale/terraform-provider-outscale
cd terraform-provider-outscale
go build -o terraform-provider-outscale_v1.3.1
```

---

## 📦 Using the Built Provider

After building the provider manually, install it locally depending on your platform and tooling:

### For Terraform

<details>
<summary>On Linux</summary>
1. Download and install [Terraform](https://www.terraform.io/downloads.html).

2. Move the plugin to the repository:
```sh
mkdir -p terraform.d/plugins/registry.terraform.io/outscale/outscale/1.3.1/linux_amd64
mv terraform-provider-outscale_v1.3.1 terraform.d/plugins/registry.terraform.io/outscale/outscale/1.3.1/linux_amd64/
```

3. Initialize Terraform:
```sh
terraform init
```

4. Plan your Terraform configuration:
```sh
terraform plan
```
</details>

<details>
<summary>On macOS</summary>
1. Download and install [Terraform](https://www.terraform.io/downloads.html).

2. Move the plugin to the repository:
```sh
mkdir -p terraform.d/plugins/registry.terraform.io/outscale/outscale/1.3.1/darwin_arm64
mv terraform-provider-outscale_v1.3.1 terraform.d/plugins/registry.terraform.io/outscale/outscale/1.3.1/darwin_arm64/
```

3. Initialize Terraform:
```sh
terraform init
```

4. Plan your Terraform configuration:
```sh
terraform plan
```
</details>

---

### For OpenTofu

<details>
<summary>On Linux</summary>
1. Download and install [OpenTofu](https://opentofu.org/docs/intro/install/deb/).

2. Move the plugin to the repository:
```sh
mkdir -p terraform.d/plugins/registry.opentofu.org/outscale/outscale/1.3.1/linux_amd64
mv terraform-provider-outscale_v1.3.1 terraform.d/plugins/registry.opentofu.org/outscale/outscale/1.3.1/linux_amd64/
```

3. Initialize OpenTofu:
```sh
tofu init
```

4. Plan your configuration:
```sh
tofu plan
```
</details>

<details>
<summary>On macOS</summary>
1. Download and install [OpenTofu](https://opentofu.org/docs/intro/install/homebrew/).

2. Move the plugin to the repository:
```sh
mkdir -p terraform.d/plugins/registry.opentofu.org/outscale/outscale/1.3.1/darwin_arm64
mv terraform-provider-outscale_v1.3.1 terraform.d/plugins/registry.opentofu.org/outscale/outscale/1.3.1/darwin_arm64/
```

3. Initialize OpenTofu:
```sh
tofu init
```

4. Plan your configuration:
```sh
tofu plan
```
</details>

---

## 🤝 Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md).

---

## 📝 Building the Documentation

Requirements:
- `make`
- `python3`
- `python-venv`

```sh
make doc
```
