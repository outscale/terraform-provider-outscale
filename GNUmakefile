GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
PKG_NAME=outscale
TEST?=./...
VERSION=$(shell git describe --exact-match 2> /dev/null || \
                 git describe --match=$(git rev-parse --short=8 HEAD) --always --dirty --abbrev=8)
WEBSITE_REPO=github.com/hashicorp/terraform-website

.PHONY: default
default: build

.PHONY: build
build: fmtcheck
	go build -ldflags "-X github.com/terraform-providers/terraform-provider-outscale/version.version=${VERSION}"

.PHONY: fmtcheck
fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

.PHONY: test
test: fmtcheck
	go test $(TEST) -count 1 -timeout=30s -parallel=4

.PHONY: testacc
testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -count 1 -v -parallel 4 $(TESTARGS) -timeout 240m -cover

.PHONY: fmt
fmt:
	gofmt -s -w ./main.go
	gofmt -s -w ./$(PKG_NAME)

.PHONY: websitefmtcheck
websitefmtcheck:
	@sh -c "'$(CURDIR)/scripts/websitefmtcheck.sh'"

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

.PHONY: test-examples
test-examples:
	@sh -c "'$(CURDIR)/scripts/test-examples.sh'"

.PHONY: test-integration
test-integration:
	@sh -c "'$(CURDIR)/scripts/integration.sh'"

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
