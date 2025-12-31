TEST?=./...
GOFMT_FILES?=$(shell find . -name '*.go')
PKG_NAME=outscale
VERSION=$(shell git describe --exact-match 2> /dev/null || \
                 git describe --match=$(git rev-parse --short=8 HEAD) --always --dirty --abbrev=8)

TF_ACC_OAPI_PARALLEL=10
TF_ACC_OAPI_NETS_PARALLEL=4
TF_ACC_OKS_PARALLEL=1

PYTEST_OAPI_PARALLEL=10
PYTEST_OAPI_NETS_PARALLEL=6
PYTEST_OKS_PARALLEL=2

# Service paths
OAPI_PKG=./internal/services/oapi
OKS_PKG=./internal/services/oks
PROVIDER_PKG=./provider

.PHONY: default
default: build

.PHONY: build
build: fmtcheck
	go build -ldflags "-X github.com/outscale/terraform-provider-outscale/version.version=${VERSION}"

.PHONY: fmtcheck
fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

.PHONY: vet
vet: fmtcheck
	go vet $(TEST)

.PHONY: test
test: fmtcheck vet
	go test $(TEST) -count 1 -timeout=30s -parallel $(TF_ACC_OAPI_PARALLEL)

# Run all acceptance tests
.PHONY: testacc
testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -count 1 -v -parallel $(TF_ACC_OAPI_NETS_PARALLEL) $(TESTARGS) -timeout 240m -cover

.PHONY: testacc-oapi
testacc-oapi: fmtcheck
	TF_ACC=1 go test $(OAPI_PKG) -count 1 -v -parallel $(TF_ACC_OAPI_PARALLEL) $(TESTARGS) -timeout 240m -cover

.PHONY: testacc-oapi-net
testacc-oapi-net: fmtcheck
	TF_ACC=1 go test $(OAPI_PKG) -run=TestAccNet -count 1 -v -parallel $(TF_ACC_OAPI_NETS_PARALLEL) $(TESTARGS) -timeout 240m -cover

.PHONY: testacc-oapi-vm
testacc-oapi-vm: fmtcheck
	TF_ACC=1 go test $(OAPI_PKG) -run=TestAccVM -count 1 -v -parallel $(TF_ACC_OAPI_PARALLEL) $(TESTARGS) -timeout 240m -cover

.PHONY: testacc-oapi-others
testacc-oapi-others: fmtcheck test-gen-cert
	TF_ACC=1 go test $(OAPI_PKG) -run=TestAccOthers -count 1 -v -parallel $(TF_ACC_OAPI_PARALLEL) $(TESTARGS) -timeout 240m -cover

.PHONY: testacc-oks
testacc-oks: fmtcheck
	TF_ACC=1 go test $(OKS_PKG) -count 1 -v -parallel $(TF_ACC_OKS_PARALLEL) $(TESTARGS) -timeout 240m -cover

# Provider tests
.PHONY: testacc-provider
testacc-provider: fmtcheck
	TF_ACC=1 go test $(PROVIDER_PKG) -count 1 -v -parallel $(TF_ACC_OAPI_PARALLEL) $(TESTARGS) -timeout 240m -cover

.PHONY: test-net
test-net: testacc-oapi-net

.PHONY: test-vm
test-vm: testacc-oapi-vm

.PHONY: test-others
test-others: testacc-oapi-others

.PHONY: fmt
fmt:
	gofmt -s -w ./internal/
	gofmt -s -w ./provider/
	gofmt -s -w ./main.go

.PHONY: lint
lint:
	@GOGC=30 golangci-lint run ./internal/...  --deadline=30m
	@GOGC=30 golangci-lint run ./provider/...  --deadline=30m

.PHONY: tools
tools:
	GO111MODULE=off go get -u github.com/client9/misspell/cmd/misspell
	GO111MODULE=off go get -u github.com/golangci/golangci-lint/cmd/golangci-lint

.PHONY: test-compile
test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./internal/services/oapi"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

.PHONY: test-locally
test-locally:
	## "$(CURDIR)/scripts/local-test.sh" TestAccOthers_Volume_io1Type  ##add request_id parameter in ricochet
	"$(CURDIR)/scripts/local-test.sh" TestAccVM_withFlexibleGpuLink_basic

.PHONY: terraform-examples
terraform-examples:
	@sh -c "'$(CURDIR)/scripts/terraform-examples.sh'"

.PHONY: tofu-examples
tofu-examples:
	@sh -c "'$(CURDIR)/scripts/tofu-examples.sh'"

.PHONY: test-integration
test-integration: test-integration-oapi-nets test-integration-oapi test-integration-oks

.PHONY: test-integration-oapi-nets
test-integration-oapi-nets: test-gen-cert
	PYTEST_PARALLEL=$(PYTEST_OAPI_NETS_PARALLEL) RUN_NETS_ONLY=1 sh -c "'$(CURDIR)/scripts/test-integration.sh' test_provider_oapi.py"

.PHONY: test-integration-oapi
test-integration-oapi: test-gen-cert
	PYTEST_PARALLEL=$(PYTEST_OAPI_PARALLEL) SKIP_NETS=1 sh -c "'$(CURDIR)/scripts/test-integration.sh' test_provider_oapi.py"

.PHONY: test-integration-oks
test-integration-oks: test-gen-cert
	PYTEST_PARALLEL=$(PYTEST_OKS_PARALLEL) sh -c "'$(CURDIR)/scripts/test-integration.sh' test_provider_oks.py"

.PHONY: test-gen-cert
test-gen-cert:
	@sh -c "'$(CURDIR)/scripts/generate-certificate.sh'"

.PHONY: doc
doc:
	@sh -c "'$(CURDIR)/scripts/generate-doc.sh'"
