LDFLAGS := "-X main.version=${VERSION}"
# Required for globs to work correctly
SHELL=/bin/bash

.PHONY: build
build:
	go build -o secure-tiller -ldflags $(LDFLAGS) ./main.go

HAS_DEP := $(shell command -v glide;)

.PHONY: bootstrap
bootstrap:
ifndef HAS_GLIDE
	go get -u github.com/Masterminds/glide
endif
	glide install -v
