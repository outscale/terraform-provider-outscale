# 3DS OUTSCALE Terraform Provider
[![Project Graduated](https://docs.outscale.com/fr/userguide/_images/Project-Graduated-green.svg)](https://docs.outscale.com/en/userguide/Open-Source-Projects.html)

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)
<img  alt="Terraform"  src="https://camo.githubusercontent.com/6d6ec94bb2909d75122df9cf17e1940b522a805587c890a2e37a57eba61f7eb1/68747470733a2f2f7777772e6461746f636d732d6173736574732e636f6d2f323838352f313632393934313234322d6c6f676f2d7465727261666f726d2d6d61696e2e737667"  width="200px">

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 1.x.x

- [Go](https://golang.org/doc/install) 1.22.0 (to build the provider plugin)


## Migration to v1

> [!WARNING]
>
> Before you begin using the ```v1``` binary on your Terraform code, make sure to back up your state file!
>
> If you are using a local state file, make copy of your terraform.tfstate file in your project directory.
>
> If you are using a remote backend such as an S3 bucket, make sure that you follow the backup procedures and that you exercise the restore procedure at least once.
>
> Additionally, make sure you backup or version your code as migration will require some code changes (on Flexible_gpu resource).

### Step 1: Upgrade provider version

```sh
terraform {
  required_providers {
    outscale = {
      source = "outscale/outscale"
      version = "1.0.1"
    }
  }
}

provider "outscale" {
  # Configuration options
}
```

```sh
terraform init -upgrade
```

### Step 2: Edit terraform state and configuration files

Some block types changed in terraform state, the following script will delete those blocks.

Then ``` terraform refresh``` or ``` terraform apply ``` will set the right block type.

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
sed -i '/outscale_flexible_gpu\./s/$/ \]/' *.tf
sed -i '/flexible_gpu_ids /s/= /= \[/' *.tf
terraform fmt
```

#### On MacOS
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
sed -i='' '/outscale_flexible_gpu\./s/$/\]/' *.tf
sed -i='' '/flexible_gpu_ids /s/= /= \[/' *.tf
terraform fmt
```
### Step 3: Refresh configuration to update terraform state

```sh
terraform refresh
```

## Breaking change

> **Warning**
>
> We have a broken change on our api when creating access_key without expiration date for all version less then v0.9.0. ([GH-issues](https://github.com/outscale/terraform-provider-outscale/issues/342))
>
> We recommend to upgrade on the latest ([v1.0.1](https://registry.terraform.io/providers/outscale/outscale/latest))

## Using the Provider with Terraform

> **Warning**
>
> Our provider terraform has been moved from [outscale-dev](https://registry.terraform.io/providers/outscale-dev/outscale/latest) to [outscale](https://registry.terraform.io/providers/outscale/outscale/latest) organisation on terraform registry
>
> The next releases will be only publish under [outscale organization on terraform registry](https://registry.terraform.io/providers/outscale/outscale/latest)

Add the following lines in the Terraform configuration to permit to get the provider from the Terrafom registry:

```sh
terraform {
  required_providers {
    outscale = {
      source = "outscale/outscale"
      version = "1.0.1"
    }
  }
}

provider "outscale" {
  # Configuration options
}
```
1. Execute `terraform init`

2. Execute `terraform plan`

## Using the Provider with OpenTofu

```sh
terraform {
  required_providers {
    outscale = {
      source = "outscale/outscale"
      version = "1.0.1"
    }
  }
}

provider "outscale" {
  # Configuration options
}
```
1. Execute `tofu init`

2. Execute `tofu plan`

## Migrating to OpenTofu from Terraform
Follow [migration link](https://opentofu.org/docs/intro/migration/)

## Configuring the proxy, if any
### on Linux/macOS
```sh
export HTTPS_PROXY=http://192.168.1.24:3128
```
### on Windows
 ```sh
set HTTPS_PROXY=http://192.168.1.24:3128
```

## x509 client authentication, if any
Add the following lines in the Terraform configuration to define certificate location:
```sh
terraform {
  required_providers {
    outscale = {
      source = "outscale/outscale"
      version = "1.0.1"
    }
  }
}

provider "outscale" {
  access_key_id = var.access_key_id
  secret_key_id = var.secret_key_id
  region = var.region
  x509_cert_path = "/myrepository/certificate/client_ca.crt"
  x509_key_path = "/myrepository/certificate/client_ca.key"
}
```
or set the following environment variables:

```sh
export OUTSCALE_X509CERT=/myrepository/certificate/client_ca.crt
export OUTSCALE_X509KEY=/myrepository/certificate/client_ca.key
```
## Building The Provider
Clone repository to: `$GOPATH/src/github.com/outscale/terraform-provider-outscale`
```sh
mkdir -p $GOPATH/src/github.com/terraform-providers
cd  $GOPATH/src/github.com/terraform-providers
git clone --branch v1.0.1 https://github.com/outscale/terraform-provider-outscale
```
Enter the provider directory and build the provider
```sh
cd  $GOPATH/src/github.com/terraform-providers/terraform-provider-outscale
go build -o terraform-provider-outscale_v1.0.1
```
## Using the provider built
### For Terraform
#### On Linux

1. Download and install [Terraform](https://www.terraform.io/downloads.html)

2. Move the plugin to the repository ~/.terraform.d/plugins/registry.terraform.io/outscale/outscale/1.0.1/linux_amd64/.
```shell
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/outscale/outscale/1.0.1/linux_amd64
mv terraform-provider-outscale_v1.0.1 ~/.terraform.d/plugins/registry.terraform.io/outscale/outscale/1.0.1/linux_amd64
```
3. Execute `terraform init`

4. Execute `terraform plan`

#### On macOS
1. Download and install [Terraform](https://www.terraform.io/downloads.html)

2. Move the plugin to the repository ~/.terraform.d/plugins/registry.terraform.io/outscale/outscale/1.0.1/darwin_arm64
```shell
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/outscale/outscale/1.0.1/darwin_arm64
mv terraform-provider-outscale_v1.0.1 ~/.terraform.d/plugins/registry.terraform.io/outscale/outscale/1.0.1/darwin_arm64
```

3. Execute `terraform init`

4. Execute `terraform plan`

### For OpenTofu
#### On Linux

1. Download and install [OpenTofu](https://opentofu.org/docs/intro/install/deb/)

2. Move the plugin to the repository ~/.terraform.d/plugins/registry.opentofu.org/outscale/outscale/1.0.1/linux_amd64/.
```shell
mkdir -p ~/.terraform.d/plugins/registry.opentofu.org/outscale/outscale/1.0.1/linux_amd64
mv terraform-provider-outscale_v1.0.1 ~/.terraform.d/plugins/registry.opentofu.org/outscale/outscale/1.0.1/linux_amd64
```
3. Execute `tofu init`

4. Execute `tofu plan`

#### On macOS
1. Download and install [OpenTofu](https://opentofu.org/docs/intro/install/homebrew/)

2. Move the plugin to the repository ~/.terraform.d/plugins/registry.opentofu.org/outscale/outscale/1.0.1/darwin_arm64
```shell
mkdir -p ~/.terraform.d/plugins/registry.opentofu.org/outscale/outscale/1.0.1/darwin_arm64
mv terraform-provider-outscale_v1.0.1 ~/.terraform.d/plugins/registry.opentofu.org/outscale/outscale/1.0.1/darwin_arm64
```

3. Execute `tofu init`

4. Execute `tofu plan`

## Issues and contributions
Check [CONTRIBUTING.md](./CONTRIBUTING.md) for more details.

## Building the documentation

Requirements:
- make
- python3
- python-venv

```shell
make doc
```
