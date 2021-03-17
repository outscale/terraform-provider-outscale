# Opening an issue

Feel free to open a [Github issue](https://github.com/outscale-dev/terraform-provider-outscale/issues) and explain your problem.

Please provide at least those informations:
- terraform version
- how to reproduce the issue
- output of your command with `TF_LOG=TRACE` set (e.g. `TF_LOG=TRACE terraform apply`)
- please store large output (like traces) as as an attached file
- make sure your don't leak any sensible informations (credentials, ...)

# Developing the Provider

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
