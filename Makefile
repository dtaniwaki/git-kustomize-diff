NAME        := git-kustomize-diff
PROJECTROOT := $(shell pwd)
VERSION     := $(if $(VERSION),$(VERSION),$(shell cat ${PROJECTROOT}/VERSION)-dev)
REVISION    := $(shell git rev-parse --short HEAD)
OUTDIR      ?= $(PROJECTROOT)/dist

LDFLAGS := -ldflags="-s -w -X \"github.com/dtaniwaki/git-kustomize-diff/cmd.Version=$(VERSION)\" -X \"github.com/dtaniwaki/git-kustomize-diff/cmd.Revision=$(REVISION)\""

.PHONY: build
build:
	go build $(LDFLAGS) -o $(OUTDIR)/$(NAME)

.PHONY: install
install:
	go install $(LDFLAGS)

.PHONY: build-linux-amd64
build-linux-amd64:
	make build \
		GOOS=linux \
		GOARCH=amd64 \
		NAME=git-kustomize-diff-linux-amd64

.PHONY: build-linux
build-linux: build-linux-amd64

.PHONY: build-darwin
build-darwin:
	make build \
		GOOS=darwin \
		NAME=git-kustomize-diff-darwin-amd64

.PHONY: build-windows
build-windows:
	make build \
		GOARCH=amd64 \
		GOOS=windows \
		NAME=git-kustomize-diff-windows-amd64.exe

.PHONY: build-all
build-all: build-linux build-darwin build-windows

.PHONY: lint
lint:
	golangci-lint run

test:
	@go test -v -race -short -tags no_e2e ./...

.PHONY: coverage
coverage:
	@go test -tags no_e2e -covermode=count -coverprofile=profile.cov -coverpkg ./pkg/...,./cmd/... $(shell go list ./... | grep -v /vendor/)
	@go tool cover -func=profile.cov

.PHONY: clean
clean:
	rm -rf $(OUTDIR)/*
