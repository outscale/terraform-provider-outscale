TEST?=./...
GOFMT_FILES?=$(shell find . -name '*.go')
PKG_NAME=outscale
TEST?=./...
VERSION=$(shell git describe --exact-match 2> /dev/null || \
                 git describe --match=$(git rev-parse --short=8 HEAD) --always --dirty --abbrev=8)
WEBSITE_REPO=github.com/hashicorp/terraform-website
TF_ACC_PARALLEL=10
TF_ACC_NETS_PARALLEL=2

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
	@sh -c "'$(CURDIR)/scripts/integration.sh'"

.PHONY: test-gen-cert
test-gen-cert:
	@sh -c "'$(CURDIR)/scripts/generate-certificate.sh'"

.PHONY: website
website:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

.PHONY: website-local
website-local:
	@sh -c "'$(CURDIR)/scripts/test-doc.sh'"

.PHONY: website-lint
website-lint:
	@misspell -error -source=text website/

.PHONY: website-test
website-test:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider-test PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

.PHONY: doc
doc:
	@sh -c "'$(CURDIR)/scripts/generate-doc.sh'"
