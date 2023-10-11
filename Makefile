#!/usr/bin/make -f

# Choosing the shell
# - [docs](https://www.gnu.org/software/make/manual/html_node/Choosing-the-Shell.html)
SHELL =/bin/bash

# Flag for enabling CGO, in development environment CGO should be enabled to enable -race flags.
CGO_ENABLED ?= 1

# BIN_DIR is the directory where built binaries will be placed.
BIN_DIR ?= bin

# CMD_DIR is the directory where main.go files are located.
CMD_DIR ?= cmd

# CMD_SET is a list of targets to build.
# The value of CMD_SET is constructed by finding one level of subdirectories under CMD_DIR and then adding the prefix
# BIN_DIR to each directory name.
#
# For example, if BIN_DIR is 'bin' and CMD_DIR is 'cmd', then there are two subdirectories within CMD_DIR:
# 'cmd/cmd-a' and 'cmd/cmd-b'. As a result, CMD_SET becomes 'bin/cmd-a bin/cmd-b'.
CMD_SET ?= $(addprefix $(BIN_DIR)/, $(shell find $(CMD_DIR) -mindepth 1 -maxdepth 1 -type d -exec basename {} \;))

# CURRENT_TIME is the current time in RFC3339 format.
CURRENT_TIME := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

# BUILD_VERSION is the git tag of the current commit.
BUILD_VERSION := $(shell git describe --tags --always --match "[0-9][0-9][0-9][0-9].*.*")

# The default target is to prepare the development environment.
all: tools setup-pre-commit

.PHONY: build
build: $(CMD_SET) # build all binaries in CMD_SET.
	@echo "Build done. Binaries:"
	@for bin in $(CMD_SET); do echo "  - $$bin"; done

.PHONY: clean
clean:
	@echo "Removing directories: '$(BIN_DIR)'."
	@rm -rf $(BIN_DIR)
	@echo "Clean done."

# The pattern rule bin/% builds binaries in the 'bin' directory.
# It is used by the `make build` target.
#
# Example: 'make bin/cmd-a' builds 'cmd/cmd-a/main.go' and places the binary in 'bin/cmd-a'.
bin/%: $(shell find . -type f -name '*.go') # ensure to rebuild if any go file changed.
	@echo "Building '$@'."
	@mkdir -p $(dir $@) # create `bin` directory if not exist.
	@go build \
		-race=$(CGO_ENABLED) \
		-ldflags "\
			-X main.buildName=$(@F) \
			-X main.buildTime=$(CURRENT_TIME) \
			-X main.buildVersion=$(BUILD_VERSION)" \
		-o $@ ./cmd/$(@F)


# Install all development tools, these tools are used by pre-commit hook.
tools: hack/install_tools.sh
	@echo "Installing tools"
	@hack/install_tools.sh
	@echo "Tools installed"


# Enable pre-commit hook.
setup-pre-commit:
	@echo "Setting up pre-commit hook"
	@cp -f hack/pre-commit.sh .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-commit


images: $(addprefix docker-image/, $(CMD_SET))
	@echo "All images built."

# Build docker image.
docker-image/%:
	@echo "Building docker image: $(@F), please wait..."
	@docker build \
		-f cmd/admin-restful/Dockerfile \
		-q \
		-t josestg/swe-be-mono-$(@F):$(BUILD_VERSION) \
		--build-arg BUILD_VERSION=$(BUILD_VERSION) \
		--build-arg BUILD_DATE=$(CURRENT_TIME) \
		--build-arg IMAGE_NAME=josestg/swe-be-mono-$(@F) \
		--build-arg GITHUB_TOKEN=$(GITHUB_TOKEN) \
		--build-arg CACHEBUST=$(shell date +%s) \
		. ;
	@echo "Docker image built: josestg/swe-be-mono-$(@F):$(BUILD_VERSION)"
