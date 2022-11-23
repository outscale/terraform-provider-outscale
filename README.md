# 3DS OUTSCALE Terraform Provider

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)
<img  alt="Terraform"  src="https://camo.githubusercontent.com/1a4ed08978379480a9b1ca95d7f4cc8eb80b45ad47c056a7cfb5c597e9315ae5/68747470733a2f2f7777772e6461746f636d732d6173736574732e636f6d2f323838352f313632393934313234322d6c6f676f2d7465727261666f726d2d6d61696e2e737667"  width="200px">

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 1.0.x

- [Go](https://golang.org/doc/install) 1.13 (to build the provider plugin)

## Using the Provider

Add the following lines in the Terraform configuration to permit to get the provider from the Terrafom registry:

```sh
terraform {
  required_providers {
    outscale = {
      source = "outscale-dev/outscale"
      version = "0.7.0"
    }
  }
}

provider "outscale" {
  # Configuration options
}
```

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
  required_version = ">= 0.13"
  required_providers {
    outscale = {
      source = "outscale-dev/outscale"
      version = "0.7.0"
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
Clone repository to: `$GOPATH/src/github.com/terraform-providers/terraform-provider-outscale`
```sh
mkdir -p $GOPATH/src/github.com/terraform-providers
cd  $GOPATH/src/github.com/terraform-providers
git clone --branch v0.7.0 https://github.com/outscale-dev/terraform-provider-outscale
```
Enter the provider directory and build the provider
```sh
cd  $GOPATH/src/github.com/terraform-providers/terraform-provider-outscale
go build -o terraform-provider-outscale_v0.7.0
```
## Using the provider
### On Linux

1. Download and install [Terraform](https://www.terraform.io/downloads.html)

2. Move the plugin to the repository ~/.terraform.d/plugins/registry.terraform.io/outscale-dev/outscale/0.7.0/linux_amd64/.
```shell
mkdir -p ~/.terraform.d/plugins/regisutry.terraform.io/outscale-dev/outscale/0.7.0/linux_amd64
mv terraform-provider-outscale_v0.7.0 ~/.terraform.d/plugins/registry.terraform.io/outscale-dev/outscale/0.7.0/linux_amd64
```
3. Execute `terraform init

4. Execute `terraform plan`

### On macOS
1. Download and install [Terraform](https://www.terraform.io/downloads.html)

2. Move the plugin to the repository ~/.terraform.d/plugins/registry.terraform.io/outscale-dev/outscale/0.7.0/darwin_arm64
```shell
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/outscale-dev/outscale/0.7.0/darwin_arm64
mv terraform-provider-outscale_v0.7.0 ~/.terraform.d/plugins/registry.terraform.io/outscale-dev/outscale/0.7.0/darwin_arm64
```  

3. Execute `terraform init`

4. Execute `terraform plan`

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
