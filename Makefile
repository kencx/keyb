binary = keyb
version = $(shell git describe --tags)
ldflags = -ldflags "-s -w -X main.version=${version}"

.PHONY: help test clean snapshot build

default: help

help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## test: run tests
test:
	go test ./...

## clean: remove binaries, dist
clean:
	if [ -f ${binary} ]; then rm keyb; fi
	go clean
	rm -rf dist

## snapshot: generate unversioned snapshot release
snapshot:
	goreleaser release --snapshot

## build: build binary
build:
	go build ${ldflags} -o ${binary}

## install: install binary at ~/.local/bin
install:
	cp ${binary} ~/.local/bin/
