PROJECT := "docker-machine-driver-kubevirt"
USER := linkacloud
REPO := go.linka.cloud/$(PROJECT)
DOCKER_MACHINE := "docker-machine"
VERSION := latest
IMAGE := $(USER)/$(DOCKER_MACHINE):$(VERSION)
D2VM_DOCKER_IMAGE := d2vm-$(DOCKER_MACHINE)
ALPINE_IMAGE := $(USER)/$(D2VM_DOCKER_IMAGE):alpine
UBUNTU_IMAGE := $(USER)/$(D2VM_DOCKER_IMAGE):ubuntu
D2VM_DOCKER_IMAGE := $(USER)/$(D2VM_DOCKER_IMAGE)

D2VM := $(shell which d2vm || "docker run --rm -i -t -v $(PWD):/build -w /build -v /var/run/docker.sock:/var/run/docker.sock $(D2VM_DOCKER_IMAGE)")

.PHONY: build-docker
build-docker:
	@docker build -t $(IMAGE) .

.PHONY: push-docker
push-docker:
	@docker push $(IMAGE)

.PHONY: docker
docker: build-docker push-docker

.PHONY: docker-alpine
docker-alpine:
	@docker build -t $(IMAGE) -f Dockerfile.alpine .


.PHONY: d2vm-docker-images
d2vm-docker-images:
	@$(D2VM) build -o images/alpine.qcow2 -t $(ALPINE_IMAGE) --push --force -v -f images/alpine.Dockerfile images
	@$(D2VM) build -o images/ubuntu.qcow2 -t $(UBUNTU_IMAGE) --push --force -v -f images/ubuntu.Dockerfile images
