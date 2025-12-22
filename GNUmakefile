TEST?=./...
GOFMT_FILES?=$(shell find . -name '*.go')
PKG_NAME=outscale
TEST?=./...
VERSION=$(shell git describe --exact-match 2> /dev/null || \
                 git describe --match=$(git rev-parse --short=8 HEAD) --always --dirty --abbrev=8)
TF_ACC_PARALLEL=10
TF_ACC_NETS_PARALLEL=2

PYTEST_PARALLEL=10
PYTEST_NETS_PARALLEL=4

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
	go test $(TEST) -count 1 -timeout=30s -parallel $(TF_ACC_PARALLEL)

.PHONY: testacc
testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -count 1 -v -parallel $(TF_ACC_NETS_PARALLEL) $(TESTARGS) -timeout 240m -cover

.PHONY: test-net
test-net: fmtcheck
	TF_ACC=1 go test $(TEST) -run=TestAccNet -count 1 -v -parallel $(TF_ACC_NETS_PARALLEL) $(TESTARGS) -timeout 240m -cover

.PHONY: test-vm
test-vm: fmtcheck
	TF_ACC=1 go test $(TEST) -run=TestAccVM -count 1 -v -parallel $(TF_ACC_PARALLEL) $(TESTARGS) -timeout 240m -cover

.PHONY: test-others
test-others: fmtcheck test-gen-cert
	TF_ACC=1 go test $(TEST) -run=TestAccOthers -count 1 -v -parallel $(TF_ACC_PARALLEL) $(TESTARGS) -timeout 240m -cover

.PHONY: fmt
fmt:
	gofmt -s -w ./utils/
	gofmt -s -w ./main.go
	gofmt -s -w ./$(PKG_NAME)

.PHONY: lint
lint:
	@GOGC=30 golangci-lint run ./$(PKG_NAME)  --deadline=30m

.PHONY: tools
tools:
	GO111MODULE=off go get -u github.com/client9/misspell/cmd/misspell
	GO111MODULE=off go get -u github.com/golangci/golangci-lint/cmd/golangci-lint

.PHONY: test-compile
test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
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
test-integration: test-gen-cert
	PYTEST_PARALLEL=$(PYTEST_PARALLEL) PYTEST_NETS_PARALLEL=$(PYTEST_NETS_PARALLEL) sh -c "'$(CURDIR)/scripts/integration.sh'"

.PHONY: test-gen-cert
test-gen-cert:
	@sh -c "'$(CURDIR)/scripts/generate-certificate.sh'"

.PHONY: doc
doc:
	@sh -c "'$(CURDIR)/scripts/generate-doc.sh'"
