SHELL := /bin/bash
export PATH := $(CURDIR)/_tools/bin:$(PATH)
USER ?= `whoami`

install:
	# linter
	go build -o _tools/bin/golangci-lint github.com/golangci/golangci-lint/cmd/golangci-lint

format:vendor
	go run github.com/golangci/golangci-lint/cmd/golangci-lint run --fix

vendor:
	go mod vendor

test:
 	#nothing here