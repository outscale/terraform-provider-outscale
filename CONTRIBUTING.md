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

If you want, you can tests terraform locally using [ricochet-2](https://github.com/outscale/osc-ricochet-2/)
to do so, you need to get ricochet-2 [release](https://github.com/outscale/osc-ricochet-2/tags)
extract and start it:
```
tar -xvf osc-ricochet-2_v0.2.0_x86_64-unknown-linux-musl.tar.gz
./ricochet-2 ./ricochet.json
```
and in another terminal, either call `make test-locally`, or call the script manually
```
# if you want TestAccVM_withFlexibleGpuLink_basic
scripts/local-test.sh TestAccVM_withFlexibleGpuLink_basic
```

Note that ricochet-2 been fearly new, doesn't support all Outscale Calls, and some tests will fails.
Also ricochet-2 work only on Linux for now.
