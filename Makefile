#!/usr/bin/make -f

# Choosing the shell
# - [docs](https://www.gnu.org/software/make/manual/html_node/Choosing-the-Shell.html)
SHELL =/bin/bash

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
	@go build -race -o $@ ./cmd/$(@F)