# Copyright 2023 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Ensure Make is run with bash shell as some syntax below is bash-specific
# Setting SHELL to bash allows bash commands to be executed by recipes.
# This is a requirement for 'setup-envtest.sh' in the test target.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.

MAKEFLAGS=--warn-undefined-variables

SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

.DEFAULT_GOAL:=help

#
# Directories.
#
# Full directory of where the Makefile resides
ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
BIN_DIR := bin
TEST_DIR := test
TOOLS_DIR := hack/tools
TOOLS_BIN_DIR := $(TOOLS_DIR)/$(BIN_DIR)
export PATH := $(abspath $(TOOLS_BIN_DIR)):$(abspath node_modules/.bin):$(PATH)

#
# Tooling Binaries.
#
GOLANGCI_LINT := $(abspath $(TOOLS_BIN_DIR)/golangci-lint)
#
# Container related variables. Releases should modify and double check these vars.
#
REGISTRY ?= ghcr.io/hivelocity
PROD_REGISTRY := ghcr.io/hivelocity
IMAGE_NAME ?= hivelocity-cloud-controller-manager-staging
CONTROLLER_IMG ?= $(REGISTRY)/$(IMAGE_NAME)
TAG ?= dev
ARCH ?= amd64
# Modify these according to your needs
PLATFORMS  = linux/amd64,linux/arm64
# This option is for running docker manifest command
export DOCKER_CLI_EXPERIMENTAL := enabled


.PHONY: all
all: help

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Binaries / Software

golangci-lint: $(GOLANGCI_LINT) ## Build a local copy of golangci-lint
$(GOLANGCI_LINT): # Download golanci-lint using hack script into tools folder.
	hack/ensure-golangci-lint.sh -b $(TOOLS_DIR)/$(BIN_DIR)

##@ Generate / Manifests

.PHONY: generate
generate: ## Run all generate-manifests, generate-go-deepcopyand generate-go-conversions targets
	$(MAKE) generate-mocks

.PHONY: ensure
ensure: ensure-boilerplate

.PHONY: ensure-boilerplate
ensure-boilerplate: ## Ensures that a boilerplate exists in each file by adding missing boilerplates
	./hack/ensure-boilerplate.sh

##@ Lint and Verify

.PHONY: modules
modules: ## Runs go mod to ensure modules are up to date.
	go mod tidy
	cd $(TOOLS_DIR); go mod tidy

.PHONY: lint
lint: $(GOLANGCI_LINT) lint-yaml lint-helm-charts
	$(GOLANGCI_LINT) run -v

.PHONY: lint-suggestions
lint-suggestions: $(GOLANGCI_LINT) ## Show optional suggestions for the Golang codebase
	## TODO: only show warnings about the new code.
	## https://golangci-lint.run/usage/faq/#how-to-integrate-golangci-lint-into-large-project-with-thousands-of-issues
	## TODO: does not work for me:
	## https://github.com/golangci/golangci/issues/55
	$(GOLANGCI_LINT) run -v -c .golangci-suggestions.yaml --new-from-rev=`git merge-base main HEAD`
## --new=true
## --new-from-rev=HEAD~2

.PHONY: lint-fix
lint-fix: $(GOLANGCI_LINT) ## Lint the Go codebase and run auto-fixers if supported by the linter.
	GOLANGCI_LINT_EXTRA_ARGS=--fix $(MAKE) lint

ALL_VERIFY_CHECKS = boilerplate shellcheck modules gen

.PHONY: verify
verify: lint checkmake $(addprefix verify-,$(ALL_VERIFY_CHECKS)) ## Run all verify-* targets
	@echo "All verify checks passed, congrats!"

.PHONY: lint-helm-charts
lint-helm-charts:
	helm template charts/ccm-hivelocity/| go run github.com/yannh/kubeconform/cmd/kubeconform@latest

.PHONY: checkmake
checkmake:
	go run github.com/mrtazz/checkmake/cmd/checkmake@0.2.1 Makefile

.PHONY: verify-modules
verify-modules: modules  ## Verify go modules are up to date
	@if !(git diff --quiet HEAD -- go.sum go.mod $(TOOLS_DIR)/go.mod $(TOOLS_DIR)/go.sum $(TEST_DIR)/go.mod $(TEST_DIR)/go.sum); then \
		git diff | cat; \
		echo "go module files are out of date"; exit 1; \
	fi
	@if (find . -name 'go.mod' | xargs -n1 grep -q -i 'k8s.io/client-go.*+incompatible'); then \
		find . -name "go.mod" -exec grep -i 'k8s.io/client-go.*+incompatible' {} \; -print; \
		echo "go module contains an incompatible client-go version"; exit 1; \
	fi

.PHONY: verify-gen
verify-gen: generate  ## Verfiy go generated files are up to date
	@if !(git diff --quiet HEAD); then \
		git diff | cat; \
		echo "generated files are out of date, run make generate"; exit 1; \
	fi

.PHONY: verify-boilerplate
verify-boilerplate: ## Verify boilerplate text exists in each file
	./hack/verify-boilerplate.sh

.PHONY: verify-shellcheck
verify-shellcheck: ## Verify shell files
	./hack/verify-shellcheck.sh

.PHONY: lint-yaml
lint-yaml: ## Lint YAML files
	pip install --quiet yamllint
	yamllint -c .github/linters/yaml-lint.yaml --strict .

##@ Clean

.PHONY: clean
clean: ## Remove all generated files
	$(MAKE) clean-bin

.PHONY: clean-bin
clean-bin: ## Remove all generated helper binaries
	rm -rf $(BIN_DIR)
	rm -rf $(TOOLS_BIN_DIR)

##@ Build

build: generate fmt vet ## Build manager binary.
	go build -o bin/manager main.go

run: generate fmt vet ## Run a controller from your host.
	go run ./main.go

## --------------------------------------
## Docker
## --------------------------------------


# Just one of the many files created
QEMUDETECT = /proc/sys/fs/binfmt_misc/qemu-m68k

docker-multiarch: qemu docker-multiarch-builder
	docker buildx build --builder docker-multiarch --pull --push \
		--platform ${PLATFORMS} \
		-t $(CONTROLLER_IMG):$(TAG) .

.PHONY: qemu docker-multiarch-builder

qemu:	${QEMUDETECT}
${QEMUDETECT}:
	docker run --rm --privileged multiarch/qemu-user-static --reset -p yes

docker-multiarch-builder: qemu
	if ! docker buildx ls | grep -w docker-multiarch > /dev/null; then \
		docker buildx create --name docker-multiarch && \
		docker buildx inspect --builder docker-multiarch --bootstrap; \
	fi

##@ Development

fmt: ## Run go fmt against code.
	go fmt ./...

vet: ## Run go vet against code.
	go vet ./...

.PHONY: generate-mocks
generate-mocks:
	cd client; go run github.com/vektra/mockery/v2@v2.21.1

.PHONY: test
test:
	go test ./...
