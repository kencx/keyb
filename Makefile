binary = keyb
version = $(shell git describe --tags)
ldflags = -ldflags "-s -w"

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

## build: build binary
build:
	go build ${ldflags} -o ${binary}
