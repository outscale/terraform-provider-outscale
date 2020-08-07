3DS OUTSCALE Terraform Provider
===============================

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img alt="Terraform" src="https://www.terraform.io/assets/images/logo-hashicorp-3f10732f.svg" width="600px">

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.12.x
-	[Go](https://golang.org/doc/install) 1.13 (to build the provider plugin)


Installing The Provider on Linux
--------------------------------

Download the binary and install it in ~/.terraform.d/plugins/linux_amd64/.

```sh
$ wget https://github.com/outscale-dev/terraform-provider-outscale/releases/download/release-0.1.0/terraform-provider-outscale_linux_amd64_v0.1.0.zip
$ unzip terraform-provider-outscale_linux_amd64_v0.1.0.zip
$ mv terraform-provider-outscale_v0.1.0 ~/.terraform.d/plugins/linux_amd64/.
```

Installing The Provider on MacOs
--------------------------------

Download the binary and install it in ~/.terraform/plugins/darwin_amd64/.

```sh
$ wget https://github.com/outscale-dev/terraform-provider-outscale/releases/download/release-0.1.0/terraform-provider-outscale_darwin_amd64_v0.1.0.zip
$ unzip terraform-provider-outscale_darwin_amd64_v0.1.0.zip
$ mv terraform-provider-outscale_v0.1.0 ~/.terraform.d/plugins/darwin_amd64/.
```

Installing The Provider on Windows
----------------------------------

Download the binary and install it in %APPDATA%\terraform.d\plugins

```sh
$ download the plugin from https://github.com/outscale-dev/terraform-provider-outscale/releases/download/release-0.1.0/terraform-provider-outscale_windows_386_v0.1.0.zip
$ unzip terraform-provider-outscale_windows_386_v0.1.0.zip
$ move terraform-provider-outscale_v0.1.0 %APPDATA%\terraform.d\plugins
```


Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/terraform-providers/terraform-provider-outscale`

```sh
$ mkdir -p $GOPATH/src/github.com/terraform-providers; cd $GOPATH/src/github.com/terraform-providers
$ git clone --branch release-0.1.0 https://github.com/outscale-dev/terraform-provider-outscale
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/terraform-providers/terraform-provider-outscale
$ go build -o terraform-provider-outscale_v0.1.0
```

Using the provider
----------------------
1. Download and install [Terraform](https://www.terraform.io/downloads.html)
2. Move the plugin to the repository ~/.terraform.d/plugins/linux_amd64/.

```shell
  $ mv terraform-provider-outscale_v0.1.0 ~/.terraform.d/plugins/linux_amd64/.
```

3. Execute `terraform plan`
4. `terraform init`

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.8+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-outscale
...
```

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

*Note:* The following environment variables must be set prior to run Acceptance Tests

```sh
$ export OUTSCALE_IMAGEID="ami-xxxxxxxx"    # i.e. "ami-4a7bf2b3"
$ export OUTSCALE_ACCESSKEYID="<ACCESSKEY>" # i.e. "XXXXXXXXXXXXXXXXXXXX"
$ export OUTSCALE_SECRETKEYID="<SECRETKEY>" # i.e. "YYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYY"
$ export OUTSCALE_REGION="<REGION>"         # i.e. "eu-west-2"
$ export OUTSCALE_ACCOUNT="<ACCOUNTPID>"    # i.e. "XXXXXXXXXXXX"
```

```sh
$ make testacc
```
