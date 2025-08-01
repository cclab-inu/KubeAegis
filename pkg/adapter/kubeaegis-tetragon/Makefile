# Image URL to use all building/pushing image targets
IMG ?= 5gsec/kubeaegis-tetragon
# Image Tag to use all building/pushing image targets
TAG ?= v0.1

CONTAINER_TOOL ?= docker
BINARY ?= bin/kubeaegis-tetragon

build:
	@go build -ldflags="-w" -o ${BINARY}  main.go

run: build
	@./${BINARY}

.PHONY: docker-build
docker-build:
	$(CONTAINER_TOOL) build -t ${IMG}:${TAG} -t ${IMG}:latest --build-arg VERSION=${TAG} .

.PHONY: docker-push
docker-push:
	$(CONTAINER_TOOL) push ${IMG}:${TAG}
	$(CONTAINER_TOOL) push ${IMG}:latest

PLATFORMS ?= linux/arm64,linux/amd64,linux/s390x,linux/ppc64le
.PHONY: docker-buildx
docker-buildx:
	# copy existing Dockerfile and insert --platform=${BUILDPLATFORM} into Dockerfile.cross, and preserve the original Dockerfile
	sed -e '1 s/\(^FROM\)/FROM --platform=\$$\{BUILDPLATFORM\}/; t' -e ' 1,// s//FROM --platform=\$$\{BUILDPLATFORM\}/' Dockerfile > Dockerfile.cross
	- $(CONTAINER_TOOL) buildx create --name project-v3-builder
	$(CONTAINER_TOOL) buildx use project-v3-builder
	- $(CONTAINER_TOOL) buildx build --push --platform=$(PLATFORMS) --build-arg VERSION=${TAG} --tag ${IMG}:${TAG} -f Dockerfile.cross . || { $(CONTAINER_TOOL) buildx rm project-v3-builder; rm Dockerfile.cross; exit 1; }
	- $(CONTAINER_TOOL) buildx rm project-v3-builder
	rm Dockerfile.cross