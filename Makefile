
# Image URL to use all building/pushing image targets
IMG ?= controller:latest
# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:trivialVersions=true"

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

all: manager

# Run tests
test: generate fmt vet manifests
	go test ./... -coverprofile cover.out --race

# Build manager binary
manager: generate fmt vet
	go build -o bin/manager main.go

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet manifests
	go run ./main.go

# Install CRDs into a cluster
install: manifests
	kustomize build config/crd | kubectl apply -f -

# Uninstall CRDs from a cluster
uninstall: manifests
	kustomize build config/crd | kubectl delete -f -

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: manifests
	cd config/manager && kustomize edit set image controller=${IMG}
	kustomize build config/default | kubectl apply -f -

# Generate manifests e.g. CRD, RBAC etc.
manifests: controller-gen
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

# Generate code
generate: controller-gen
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

# Build the docker image
docker-build: test
	docker build . -t ${IMG}

# Push the docker image
docker-push:
	docker push ${IMG}

# find or download controller-gen
# download controller-gen if necessary
controller-gen:
ifeq (, $(shell which controller-gen))
	@{ \
	set -e ;\
	CONTROLLER_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$CONTROLLER_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.4.1;\
	rm -rf $$CONTROLLER_GEN_TMP_DIR ;\
	}
CONTROLLER_GEN=$(GOBIN)/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif

# find or download controller-gen v0.3.0
# used generate CRD for Kubernetes < v1.16
controller-gen-v3:
ifeq (, $(shell which controller-gen-v3))
	@{ \
	set -e ;\
	CONTROLLER_GEN_TMP_DIR=$$(mktemp -d) ;\
	git -c advice.detachedHead=false clone --single-branch -b v0.3.0 \
	  https://github.com/kubernetes-sigs/controller-tools.git $$CONTROLLER_GEN_TMP_DIR ;\
	cd $$CONTROLLER_GEN_TMP_DIR ;\
	go build -o $(GOBIN)/controller-gen-v3 ./cmd/controller-gen/... ;\
	rm -rf $$CONTROLLER_GEN_TMP_DIR ;\
	}
CONTROLLER_GEN_V3=$(GOBIN)/controller-gen-v3
else
CONTROLLER_GEN_V3=$(shell which controller-gen-v3)
endif

# install tools
install-tools: controller-gen controller-gen-v3

# generate all
generate-all: generate manifests generate-install generate-install-legacy

# generate install manifests
generate-install: controller-gen
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=scheduled-pod-autoscaler-role paths="./..." \
	  output:crd:artifacts:config=manifests/crd \
	  output:rbac:artifacts:config=manifests/rbac
	kustomize build ./manifests/install/ > ./manifests/install/install.yaml

# generate install manifests (Kubernetes < v1.16)
generate-install-legacy: controller-gen-v3
	$(CONTROLLER_GEN_V3) $(CRD_OPTIONS) paths="./..." output:crd:artifacts:config=manifests/crd/legacy
	kustomize build ./manifests/install/legacy/ > ./manifests/install/legacy/install.yaml

# check generated files up to date
# If this fails, try "make generate-all"
check-generated-files-up-to-date: generate-all
	git diff --exit-code
