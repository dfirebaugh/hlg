
GO ?= go

GOLANGCI_LINT_PACKAGE ?= github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.6
GOFUMPT_PACKAGE ?= mvdan.cc/gofumpt@latest
GODOC_PACKAGE ?= golang.org/x/tools/cmd/godoc@latest

build-docs:
	bash ./scripts/build_docs.sh

serve-docs:
	bash ./scripts/serve_docs.sh

deploy-docs:
	bash scripts/deploy_docs.sh

clean:
	rm -rf .dist
	rm -rf docs/user_docs/.book

.PHONY: go-docs
go-docs: ## serve up the go-docs
	@echo "go docs serve up at http://localhost:6060/pkg/github.com/dfirebaugh/hlg/"
	$(GO) run $(GODOC_PACKAGE) -http=:6060

.PHONY: deps-tools
deps-tools: ## install tool dependencies
	$(GO) install $(GOLANGCI_LINT_PACKAGE)
	$(GO) install $(GOFUMPT_PACKAGE)

.PHONY: format
format: ## checks formatting
	gofumpt -l -w .

.PHONY: lint
lint: ## lint 
	$(GO) run $(GOLANGCI_LINT_PACKAGE) run

