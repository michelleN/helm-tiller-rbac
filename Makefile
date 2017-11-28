# Required for globs to work correctly
SHELL=/bin/bash

HELM_HOME ?= $(shell helm home)
HELM_PLUGIN_DIR ?= $(HELM_HOME)/plugins/helm-secure-tiller
PLUGIN_REPO ?= helm-secure-tiller
PLUGIN_BINARY ?= secure-tiller
HAS_DEP := $(shell command -v glide;)
VERSION := $(shell sed -n -e 's/version:[ "]*\([^"]*\).*/\1/p' plugin.yaml)
DIST := $(CURDIR)/_dist
LDFLAGS := "-X main.version=${VERSION}"

.PHONY: install
install: bootstrap build
	helm plugin install $(GOPATH)/src/github.com/michelleN/helm-secure-tiller

.PHONY: build
build:
	go build -o secure-tiller -ldflags $(LDFLAGS) ./main.go

.PHONY: dist
dist:
	mkdir -p $(DIST)
	GOOS=linux GOARCH=amd64 go build -o $(PLUGIN_BINARY) -ldflags $(LDFLAGS) ./main.go
	tar -zcvf $(DIST)/$(PLUGIN_REPO)-tiller-linux-$(VERSION).tgz $(PLUGIN_BINARY) README.md LICENSE.txt plugin.yaml
	GOOS=darwin GOARCH=amd64 go build -o $(PLUGIN_BINARY) -ldflags $(LDFLAGS) ./main.go
	tar -zcvf $(DIST)/$(PLUGIN_REPO)-macos-$(VERSION).tgz $(PLUGIN_BINARY) README.md LICENSE.txt plugin.yaml
	GOOS=windows GOARCH=amd64 go build -o $(PLUGIN_BINARY).exe -ldflags $(LDFLAGS) ./main.go
	tar -zcvf $(DIST)/$(PLUGIN_REPO)-windows-$(VERSION).tgz $(PLUGIN_BINARY).exe README.md LICENSE.txt plugin.yaml

.PHONY: bootstrap
bootstrap:
ifndef HAS_GLIDE
	go get -u github.com/Masterminds/glide
endif
	glide install -v
