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

# BUILD_TAGS is a list of build tags to be passed to the go build command.
BUILD_TAGS 	  := "swagger_docs_enabled"

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

.PHONY: lint
lint: hack/go-lint.sh
	@chmod +x hack/go-lint.sh
	@echo "Running linter."
	@hack/go-lint.sh
	@echo "Linter done."

.PHONY: test
test: hack/go-unittest.sh
	@chmod +x hack/go-unittest.sh
	@echo "Running unit tests."
	@hack/go-unittest.sh
	@echo "Unit tests done."

# The pattern rule bin/% builds binaries in the 'bin' directory.
# It is used by the `make build` target.
#
# Example: 'make bin/cmd-a' builds 'cmd/cmd-a/main.go' and places the binary in 'bin/cmd-a'.
bin/%: $(shell find . -type f -name '*.go') # ensure to rebuild if any go file changed.
	@echo "Building '$@'."
	@mkdir -p $(dir $@) # create `bin` directory if not exist.
	@go build \
		-race=$(CGO_ENABLED) \
		-tags=$(BUILD_TAGS) \
		-ldflags "\
			-X main.buildName=$(@F) \
			-X main.buildTime=$(CURRENT_TIME) \
			-X main.buildVersion=$(BUILD_VERSION)" \
		-o $@ ./cmd/$(@F)


# Install all development tools, these tools are used by pre-commit hook.
tools: hack/install-tools.sh
	@echo "Installing tools"
	@hack/install-tools.sh
	@echo "Tools installed"


# Enable pre-commit hook.
setup-pre-commit:
	@echo "Setting up pre-commit hook"
	@cp -f hack/pre-commit.sh .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-commit


images: $(addprefix docker-image/, $(CMD_SET))
	@echo "All images built."

# Build docker image.
QUITE ?= false
docker-image/%:
	@echo "Building docker image: $(@F), please wait..."
	@docker build \
		-f cmd/admin-restful/Dockerfile \
		-q=$(QUITE) \
		-t josestg/swe-be-mono-$(@F):$(BUILD_VERSION) \
		--build-arg BUILD_VERSION=$(BUILD_VERSION) \
		--build-arg BUILD_DATE=$(CURRENT_TIME) \
		--build-arg IMAGE_NAME=josestg/swe-be-mono-$(@F) \
		--build-arg GITHUB_TOKEN=$(GITHUB_TOKEN) \
		--build-arg CACHEBUST=$(shell date +%s) \
		. ;
	@echo "Docker image built: josestg/swe-be-mono-$(@F):$(BUILD_VERSION)"


# Swagger related targets.
swagger:
	@make swagger-gen/admin-restful \
		SW_TAGS="Admin,System" \
		SW_INSTANCE_NAME="admin"
	@make swagger-gen/enduser-restful \
		SW_TAGS="Enduser,System" \
		SW_INSTANCE_NAME="enduser"

SW_TAGS ?= "System"
SW_INSTANCE_NAME ?= "swagger"
swagger-gen/%:
	@if [[ $(BUILD_TAGS) == *"swagger_docs_enabled"* ]]; then \
  		printf "[tags=\"$(SW_TAGS)\" instance=\"$(SW_INSTANCE_NAME)\"] Generating swagger docs for $(@F) ...\n" && \
		swag init \
				--generalInfo cmd/$(@F)/main*_with_swagger.go \
				--output cmd/$(@F)/swagger-docs \
				--parseDependency \
				--overridesFile .swaggo \
			  	--instanceName $(SW_INSTANCE_NAME) \
				--tags $(SW_TAGS) \
				--outputTypes go,json \
				-q; \
	else \
		echo "swagger_docs_enabled tags not set; skipping swagger docs generation."; \
	fi


# ---- Docker ----
infra-up:
	docker compose -f docker-compose.yaml up -d

infra-down:
	docker compose -f docker-compose.yaml down


# ---- Database Migration ----
MIGRATION_DRIVER ?= postgre
ifeq ($(MIGRATION_DRIVER),postgre)
	db_migration_dsn := postgres://${DB_POSTGRE_USER}:${DB_POSTGRE_PASSWORD}@${DB_POSTGRE_HOST}:${DB_POSTGRE_PORT}/${DB_POSTGRE_DATABASE}?${DB_POSTGRE_CONN_QUERY}
	db_migration_dir := resources/migrations/postgre
endif
define exec_dbmate
	dbmate -d "${db_migration_dir}" -u "${db_migration_dsn}" $(1) $(2)
endef

.PHONY: db-new
db-new: # create database.
	@$(call exec_dbmate,new,$(name))

.PHONY: db-status
db-status: # show database migration status.
	@$(call exec_dbmate,status)

.PHONY: db-migrate
db-migrate: # run database migration.
	@$(call exec_dbmate,migrate)

.PHONY: db-rollback
db-rollback: # rollback database migration.
	@$(call exec_dbmate,rollback)
