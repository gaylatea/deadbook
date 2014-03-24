all: build

build:
	@gox -osarch "darwin/amd64 linux/amd64" -output "./bin/deadbook_{{.OS}}.{{.Arch}}"

setup:
	@go get github.com/mitchellh/gox
	@gox -build-toolchain -osarch "darwin/amd64 linux/amd64"

.PHONY: all build setup