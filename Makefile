# The development version of clang is distributed as the 'clang' binary,
# while stable/released versions have a version number attached.
# Pin the default clang to a stable version.
CLANG ?= clang-14
STRIP ?= llvm-strip-14
OBJCOPY ?= llvm-objcopy-14
TARGET_ARCH ?= arm64	# x86 or arm64
CFLAGS := -O2 -g -Wall -Werror -D__TARGET_ARCH_$(TARGET_ARCH) $(CFLAGS)

# Obtain an absolute path to the directory of the Makefile.
# Assume the Makefile is in the root of the repository.
REPODIR := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
UIDGID := $(shell stat -c '%u:%g' ${REPODIR})

# Prefer podman if installed, otherwise use docker.
# Note: Setting the var at runtime will always override.
CONTAINER_ENGINE ?= docker
CONTAINER_RUN_ARGS ?= $(--user "${UIDGID}")

IMAGE_GENERATE := ebpf-builder
VERSION_GENERATE := v1-1.22.1
GENERATE_DOCKERFILE := ebpf-builder/Dockerfile

# clang <8 doesn't tag relocs properly (STT_NOTYPE)
# clang 9 is the first version emitting BTF
TARGETS := \

.PHONY: go_builder_image_build
go_builder_image_build:
	docker build -t ${IMAGE_GENERATE}:${VERSION_GENERATE} -f ${GENERATE_DOCKERFILE} .


.PHONY: all clean go_generate container-shell generate

.DEFAULT_TARGET = go_generate

# Build all ELF binaries using a containerized LLVM toolchain.
go_generate:
	+${CONTAINER_ENGINE} run --rm ${CONTAINER_RUN_ARGS} \
		-v "${REPODIR}":/ebpf -w /ebpf --env MAKEFLAGS \
		--env CFLAGS="-fdebug-prefix-map=/ebpf=." \
		--env HOME="/tmp" \
		"${IMAGE_GENERATE}:${VERSION_GENERATE}" \
		make all

# (debug) Drop the user into a shell inside the container as root.
container-shell:
	${CONTAINER_ENGINE} run --rm -ti \
		-v "${REPODIR}":/ebpf -w /ebpf \
		"${IMAGE_GENERATE}:${VERSION_GENERATE}"


all: generate

# $BPF_CLANG is used in go:generate invocations.
generate: export BPF_CLANG := $(CLANG)
generate: export BPF_CFLAGS := $(CFLAGS)
generate:
	go generate ./...

%-el.elf: %.c
	$(CLANG) $(CFLAGS) -target bpfel -g -c $< -o $@
	$(STRIP) -g $@

%-eb.elf : %.c
	$(CLANG) $(CFLAGS) -target bpfeb -c $< -o $@
	$(STRIP) -g $@


## Oodle Image

OODLE_IMAGE_NAME := k8sebpf
OODLE_TAG ?= latest
REGISTRY ?= public.ecr.aws/h3e9h5g1/oodle
OODLE_DOCKERFILE ?= Dockerfile.default
BUILDX_BUILDER := buildx-multi-arch

ifeq ($(TARGET_ARCH), arm64)
	DOCKER_PLATFORM := linux/arm64
else
	DOCKER_PLATFORM := linux/amd64
endif

.PHONY: build_push_buildx
build_push_buildx:
	docker buildx inspect $(BUILDX_BUILDER) || \
	docker buildx create --name=$(BUILDX_BUILDER) && \
	docker buildx build --push --platform=$(DOCKER_PLATFORM) --builder=$(BUILDX_BUILDER) --build-arg OODLE_TAG=$(OODLE_TAG) --build-arg VERSION=$(OODLE_TAG) --tag=$(REGISTRY)/$(OODLE_IMAGE_NAME):$(OODLE_TAG)-$(TARGET_ARCH) -f $(OODLE_DOCKERFILE) .

.PHONY: docker_merge_platforms
docker_merge_platforms:
	docker buildx imagetools create --tag $(REGISTRY)/$(OODLE_IMAGE_NAME):$(OODLE_TAG) $(REGISTRY)/$(OODLE_IMAGE_NAME):$(OODLE_TAG)-arm64 $(REGISTRY)/$(OODLE_IMAGE_NAME):$(OODLE_TAG)-x86

.PHONY: build_push
build_push:
	docker build --build-arg VERSION=$(OODLE_TAG) -t $(REGISTRY)/$(OODLE_IMAGE_NAME):$(OODLE_TAG)  -f $(OODLE_DOCKERFILE) .
	docker push $(REGISTRY)/$(OODLE_IMAGE_NAME):$(OODLE_TAG)

# make go_builder_image_build
# OODLE_TAG=latest
# make go_generate TARGET_ARCH=arm64
# make build_push_buildx TARGET_ARCH=arm64 OODLE_TAG=$OODLE_TAG

# make go_generate TARGET_ARCH=x86
# make build_push_buildx TARGET_ARCH=x86 OODLE_TAG=$OODLE_TAG

# make docker_merge_platforms OODLE_TAG=$OODLE_TAG
