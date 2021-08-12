SHELL=/bin/bash

# Go parameters
GOCMD=GO111MODULE=on go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test

compile:
	mkdir -p build
	@$(GOBUILD) -o build/deploy-tool cmd/main.go

tool:
	@echo test case $(m)
	./build/deploy-tool -config=build/config.json -m=$(m)

clean: