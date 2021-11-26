SHELL=/bin/bash

# Go parameters
GOCMD=GO111MODULE=on go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
ENV=$(DEPLOYTOOL)

compile:
	mkdir -p build
	@$(GOBUILD) -o build/$(ENV)/deploy-tool cmd/main.go

compile-linux-amd64:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 @$(GOBUILD) -o build/$(ENV)/deploy-tool-linux-amd64 cmd/main.go

tool:
	@echo test case $(m)
	./build/$(ENV)/deploy-tool -config=build/config.json -m=$(m)

clean: