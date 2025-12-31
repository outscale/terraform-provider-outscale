# Contributing Guidelines

## Opening an Issue

Feel free to open a [GitHub issue](https://github.com/outscale/terraform-provider-outscale/issues) to report a problem or suggest an improvement.

Please include the following information when opening an issue:

- **Terraform version**: Specify the version you are using.
- **Steps to reproduce the issue**: Provide detailed steps.
- **Command output**: Include the output of your command with `TF_LOG=TRACE` enabled (e.g., `TF_LOG=TRACE terraform apply`).
- **Large outputs**: Attach large outputs (e.g., logs or traces) as a file instead of pasting them into the issue.
- **Sensitive information**: Ensure no sensitive information (e.g., credentials) is included in your report.

---

## Developing the Provider

### Prerequisites

To work on the provider, you need:

1. **Go**: Install [Go](https://golang.org) (version 1.23+ required).
2. **GOPATH setup**: Configure your [GOPATH](https://golang.org/doc/code.html#GOPATH) and add `$GOPATH/bin` to your `$PATH`.

### Build the Provider

To compile the provider, use the following command:

```sh
make build
```

This will build the provider and put the binary in the `$GOPATH/bin` directory.

Example:

```sh
$ make build
$ $GOPATH/bin/terraform-provider-outscale
```

### Test the Provider

#### Unit Tests

Run the unit tests with:

```sh
make test
```

#### Acceptance Tests

To run the full suite of acceptance tests, use:

```sh
make testacc
```

You can also run tests for specific services or categories:

```sh
make testacc-oapi         # All OAPI service tests
make testacc-oapi-net     # OAPI network tests only
make testacc-oapi-vm      # OAPI VM tests only
make testacc-oapi-others  # OAPI other tests only
make testacc-oks          # All OKS service tests
make testacc-provider     # Provider-level tests only
```

**Notes:**

- **Cost**: Acceptance tests create real resources and may incur costs.
- **Environment variables**: Set the following variables before running the tests:

```sh
export OUTSCALE_IMAGEID="ami-xxxxxxxx"    # Example: "ami-e58ac287"
export OUTSCALE_ACCESSKEYID="<ACCESSKEY>" # Example: "XXXXXXXXXXXXXXXXXXXX"
export OUTSCALE_SECRETKEYID="<SECRETKEY>" # Example: "YYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYY"
export OUTSCALE_REGION="<REGION>"         # Example: "eu-west-2"
export OUTSCALE_ACCOUNT="<ACCOUNTPID>"    # Example: "XXXXXXXXXXXX"
```

---

### Local Testing with Ricochet-2

You can test Terraform locally using [Ricochet-2](https://github.com/outscale/osc-ricochet-2). To do this:

1. Download the [latest Ricochet-2 release](https://github.com/outscale/osc-ricochet-2/tags).
2. Extract and start it:

```sh
tar -xvf osc-ricochet-2_v0.2.0_x86_64-unknown-linux-musl.tar.gz
./ricochet-2 ./ricochet.json
```

3. In another terminal, run:

```sh
make test-locally
```

Alternatively, run specific tests manually:

```sh
scripts/local-test.sh TestAccVM_withFlexibleGpuLink_basic
```

**Limitations:**

- Ricochet-2 is still experimental and does not support all Outscale API calls. Some tests may fail.
- Ricochet-2 currently works only on Linux.