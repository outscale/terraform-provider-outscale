Terraform Provider
==================

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/source/assets/images/logo-text.svg" width="600px">

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.11.x
-	[Go](https://golang.org/doc/install) 1.12 (to build the provider plugin)

Installing The Provider
-----------------------

Download the binary and install it in ~/.terraform.d/plugins/linux_amd64/.

```sh
$ wget https://github.com/outscale/terraform-provider-outscale/releases/download/release-0.1.0RC3/terraform-provider-outscale_linux_amd64_v0.1.0-rc3.zip
$ unzip terraform-provider-outscale_linux_amd64_v0.1.0-rc3.zip
$ mv terraform-provider-outscale_v0.1.0-rc3 ~/.terraform.d/plugins/linux_amd64/.
```


Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/terraform-providers/terraform-provider-outscale`

```sh
$ mkdir -p $GOPATH/src/github.com/terraform-providers; cd $GOPATH/src/github.com/terraform-providers
$ git clone --branch release-0.1.0RC3 git@github.com:outscale/terraform-provider-outscale
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/terraform-providers/terraform-provider-outscale
$ go build -o terraform-provider-outscale_v0.1.0-rc3
```

Using the provider
----------------------
1. Download and install [Terraform](https://www.terraform.io/downloads.html)
2. Move the plugin to the repository ~/.terraform.d/plugins/linux_amd64/.

```shell
  $ mv terraform-provider-outscale_v0.1.0-rc3 ~/.terraform.d/plugins/linux_amd64/.
```

3. Execute `terraform init`
4. Execute `terraform plan`
5. oAPI beta documentation is available at https://docs-beta.outscale.com

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

```sh
$ make testacc
```
