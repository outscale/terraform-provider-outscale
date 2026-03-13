PKG?=./...
GOFMT_FILES?=$(shell find . -name '*.go')
PKG_NAME=outscale
VERSION=$(shell git describe --exact-match 2> /dev/null || \
                 git describe --match=$(git rev-parse --short=8 HEAD) --always --dirty --abbrev=8)

TF_ACC_OAPI_PARALLEL=10
TF_ACC_OAPI_NETS_PARALLEL=3
TF_ACC_OKS_PARALLEL=1

PYTEST_OAPI_PARALLEL=10
PYTEST_OAPI_NETS_PARALLEL=6
PYTEST_OKS_PARALLEL=2

# Integration test provider version
TF_PROVIDER_VERSION=1.0.0-test
PROVIDER_BINARY=terraform-provider-outscale_v$(TF_PROVIDER_VERSION)

UNAME_S := $(shell uname -s)
UNAME_M := $(shell uname -m)
ifeq ($(UNAME_S),Linux)
	PLUGIN_ARCH=linux_amd64
else ifeq ($(UNAME_S),Darwin)
	ifeq ($(UNAME_M),arm64)
		PLUGIN_ARCH=darwin_arm64
	else
		PLUGIN_ARCH=darwin_amd64
	endif
else
	$(error OS $(UNAME_S) is not supported)
endif

PLUGIN_DIR=tests/terraform.d/plugins/registry.terraform.io/outscale/outscale/$(TF_PROVIDER_VERSION)/$(PLUGIN_ARCH)

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
	go vet $(PKG)

.PHONY: test
test: fmtcheck vet
	go test $(PKG) -count 1 -timeout=30s -parallel $(TF_ACC_OAPI_PARALLEL)

# Run all acceptance tests
.PHONY: testacc
testacc: fmtcheck
	TF_ACC=1 go test $(PKG) -count 1 -v -parallel $(TF_ACC_OAPI_NETS_PARALLEL) $(TESTARGS) -timeout 240m -cover

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

.PHONY: testacc-oapi-frieza
testacc-oapi-frieza:
	@"$(CURDIR)/scripts/frieza-wrap.sh" testacc-oapi-snapshot testacc-oapi

.PHONY: testacc-oapi-net-frieza
testacc-oapi-net-frieza:
	@"$(CURDIR)/scripts/frieza-wrap.sh" testacc-oapi-net-snapshot testacc-oapi-net

.PHONY: testacc-oapi-vm-frieza
testacc-oapi-vm-frieza:
	@"$(CURDIR)/scripts/frieza-wrap.sh" testacc-oapi-vm-snapshot testacc-oapi-vm

.PHONY: testacc-oapi-others-frieza
testacc-oapi-others-frieza:
	@"$(CURDIR)/scripts/frieza-wrap.sh" testacc-oapi-others-snapshot testacc-oapi-others

.PHONY: testacc-oks-frieza
testacc-oks-frieza:
	@"$(CURDIR)/scripts/frieza-wrap.sh" testacc-oks-snapshot testacc-oks

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
	@if [ "$(PKG)" = "./..." ]; then \
		echo "ERROR: Set PKG to a specific package. For example,"; \
		echo "  make test-compile PKG=./internal/services/oapi"; \
		exit 1; \
	fi
	go test -c $(PKG) $(TESTARGS)

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

.PHONY: build-test-provider
build-test-provider:
	@go build -o $(PROVIDER_BINARY)

.PHONY: install-test-provider
install-test-provider: build-test-provider
	@mkdir -p $(PLUGIN_DIR)
	@cp $(PROVIDER_BINARY) $(PLUGIN_DIR)/

# Function to run pytest with parallelization and retry on failure
# Usage: $(call run-pytest,<test_file>,<parallel_workers>,[extra_env_vars])
define run-pytest
	@echo "Running $(1) with $(2) workers"
	@cd tests && $(3) python3 -m pytest -n $(2) -v $(1) || \
		python3 -m pytest --lf -n $(2) -v $(1)
endef

.PHONY: test-integration
test-integration: test-integration-oapi-nets test-integration-oapi test-integration-oks

.PHONY: test-integration-oapi-nets
test-integration-oapi-nets: install-test-provider test-gen-cert
	$(call run-pytest,test_provider_oapi.py,$(PYTEST_OAPI_NETS_PARALLEL),RUN_NETS_ONLY=1)

.PHONY: test-integration-oapi
test-integration-oapi: install-test-provider test-gen-cert
	$(call run-pytest,test_provider_oapi.py,$(PYTEST_OAPI_PARALLEL),SKIP_NETS=1)

.PHONY: test-integration-oks
test-integration-oks: install-test-provider test-gen-cert
	$(call run-pytest,test_provider_oks.py,$(PYTEST_OKS_PARALLEL))

# Frieza-wrapped integration tests
.PHONY: test-integration-oapi-nets-frieza
test-integration-oapi-nets-frieza:
	@"$(CURDIR)/scripts/frieza-wrap.sh" integration-oapi-nets-snapshot test-integration-oapi-nets

.PHONY: test-integration-oapi-frieza
test-integration-oapi-frieza:
	@"$(CURDIR)/scripts/frieza-wrap.sh" integration-oapi-snapshot test-integration-oapi

.PHONY: test-integration-oks-frieza
test-integration-oks-frieza:
	@"$(CURDIR)/scripts/frieza-wrap.sh" integration-oks-snapshot test-integration-oks

.PHONY: test-integration-single-frieza
test-integration-single-frieza:
	@if [ -z "$(TEST)" ]; then \
		echo "Usage: make test-integration-single-frieza TEST=TF-10"; \
		exit 1; \
	fi
	@TEST_FILTER="$(TEST)" "$(CURDIR)/scripts/frieza-wrap.sh" integration-single-snapshot test-integration-single

.PHONY: test-integration-single
test-integration-single: install-test-provider test-gen-cert
	@if [ -z "$(TEST)" ]; then \
		echo "Usage: make test-integration-single TEST=TF-10"; \
		exit 1; \
	fi
	@cd tests && TF_LOG=DEBUG python3 -m pytest -s -v test_provider_oapi.py -k "$(TEST)_"

.PHONY: testacc-single-frieza
testacc-single-frieza:
	@if [ -z "$(TEST)" ]; then \
		echo "Usage: make testacc-single-frieza TEST=TestAccVM_basic"; \
		exit 1; \
	fi
	@TEST_FILTER="$(TEST)" "$(CURDIR)/scripts/frieza-wrap.sh" testacc-single-snapshot testacc-single

.PHONY: testacc-single
testacc-single: fmtcheck
	@if [ -z "$(TEST)" ]; then \
		echo "Usage: make testacc-single TEST=TestAccVM_basic"; \
		exit 1; \
	fi
	@TF_LOG=DEBUG TF_ACC=1 go test $(PKG) -run="$(TEST)" -count 1 -v -parallel 1 -timeout 240m -cover

.PHONY: clean-test-provider
clean-test-provider:
	@rm -f $(PROVIDER_BINARY)
	@rm -rf tests/terraform.d

.PHONY: test-gen-cert
test-gen-cert:
	@sh -c "'$(CURDIR)/scripts/generate-certificate.sh'"

.PHONY: doc
doc:
	@sh -c "'$(CURDIR)/scripts/generate-doc.sh'"
