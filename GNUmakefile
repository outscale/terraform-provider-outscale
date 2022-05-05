TEST?=./...
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
PKG_NAME=outscale
WEBSITE_REPO=github.com/hashicorp/terraform-website

default: build

VERSION=$(shell git describe --exact-match 2> /dev/null || \
                 git describe --match=$(git rev-parse --short=8 HEAD) --always --dirty --abbrev=8)

build: fmtcheck
	go build -ldflags "-X github.com/terraform-providers/terraform-provider-outscale/version.version=${VERSION}"


test: fmtcheck
	go test $(TEST) -count 1 -timeout=30s -parallel=4

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -count 1 -v -parallel 4 $(TESTARGS) -timeout 240m -cover

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -s -w ./main.go
	gofmt -s -w ./$(PKG_NAME)

# Currently required by tf-deploy compile
fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

websitefmtcheck:
	@sh -c "'$(CURDIR)/scripts/websitefmtcheck.sh'"

lint:
	@echo "==> Checking source code against linters..."
	@GOGC=30 golangci-lint run ./$(PKG_NAME)  --deadline=30m

tools:
	GO111MODULE=off go get -u github.com/client9/misspell/cmd/misspell
	GO111MODULE=off go get -u github.com/golangci/golangci-lint/cmd/golangci-lint

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)
test-integration:
	@sh -c "'$(CURDIR)/scripts/integration.sh'"
website:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

website-local:
	@sh -c "'$(CURDIR)/scripts/test-doc.sh'"

website-lint:
	@echo "==> Checking website against linters..."
	@misspell -error -source=text website/

website-test:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider-test PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

# In order to test examples, you will need to set
# - TF_VAR_access_key_id
# - TF_VAR_secret_key_id
# - TF_VAR_region
examples-test:
	find ./examples -mindepth 1 -maxdepth 1 -type d | \
	while read example; do \
		cd "$$example"; \
		terraform init; \
		terraform apply -auto-approve; \
		terraform destroy -auto-approve; \
		cd ../..; \
	done\

doc:
	@sh -c "'$(CURDIR)/scripts/generate-doc.sh'"

.PHONY: build test testacc fmt fmtcheck lint tools test-compile website website-lint website-test examples-test website-local test-integration doc
