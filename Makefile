
# Image URL to use all building/pushing image targets
IMG ?= cloudctl:latest
# ENVTEST_K8S_VERSION refers to the version of kubebuilder assets to be downloaded by envtest binary.
ENVTEST_K8S_VERSION = 1.25.3

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Get the WorkPath
WORKPATH = $(shell pwd)

# Setting SHELL to bash allows bash commands to be executed by recipes.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

.PHONY: all
all: build

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./src/*

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./src/*

##@ Build
.PHONY: build
build: fmt vet ## Build manager binary.
	go build -o cloudctl main.go

.PHONY: run
run: fmt vet ## Run a controller from your host.
	export CLOUDCTL_CONFIG_PATH=./src/config/crd_configs/openstack
	go run ./main.go

# If you wish built the manager image targeting other platforms you can use the --platform flag.
# (i.e. docker build --platform linux/arm64 ). However, you must enable docker buildKit for it.
# More info: https://docs.docker.com/develop/develop-images/build_enhancements/
.PHONY: docker-build
docker-build:## Build docker image with the manager.
	docker build -t ${IMG} .
	docker save -o ${IMG}.tar ${IMG}
	ctr -n=k8s.io image import ${IMG}.tar
	rm ./${IMG}.tar

.PHONY: docker-push
docker-push: ## Push docker image with the manager.
	docker push ${IMG}


##@ Deployment
.PHONY: install
install: ## Install CRDs into the K8s cluster
	kubectl apply -f ./config/crd_yamls/openstack/

.PHONY: uninstall
uninstall:  ## Uninstall CRDs from the K8s cluster.
	kubectl delete -f ./config/crd_yamls/openstack/

.PHONY: deploy
deploy:  ## Deploy controller to the K8s cluster specified in ~/.kube/config.
	kubectl apply -f ./deployment/rbac.yaml
	kubectl apply -f ./deployment/cloudctl.yaml

.PHONY: undeploy
undeploy: ## Undeploy controller from the K8s cluster specified in ~/.kube/config. Call with ignore-not-found=true to ignore resource not found errors during deletion.
	kubectl delete -f ./deployment/cloudctl.yaml
	kubectl delete -f ./deployment/rbac.yaml
