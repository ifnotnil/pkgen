SHELL := /usr/bin/env bash

.NOTPARALLEL:
.SECONDEXPANSION:
## NOTINTERMEDIATE requires make >=4.4
.NOTINTERMEDIATE:

include $(CURDIR)/scripts/go.mk
include $(CURDIR)/tools/tools.mk

.DEFAULT_GOAL=default
.PHONY: default
default: checks test

# man git-clean
.PHONY: git-reset
git-reset:
	git reset --hard
	git clean -fd

.PHONY: checks
checks: vet staticcheck gofumpt goimports golangci-lint

.PHONY: ci-gen-n-format
ci-gen-n-format: mockery golangci-lint-fmt
	@./scripts/git-check-dirty

.PHONY: ci-mod
ci-mod: mod
	@./scripts/git-check-dirty

.PHONY: ci-sh
ci-sh: shfmt shellcheck
	@./scripts/git-check-dirty
