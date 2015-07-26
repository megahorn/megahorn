.PHONEY: deps build test release all
.DEFAULT_GOAL := all

NAME=webminal
ARCH=$(shell uname -m)
VERSION=$(shell cat VERSION)
BUILD=$(shell git rev-parse --short HEAD)

deps:
	go get ./...

build: deps
	mkdir -p build/Linux && GOOS=linux go build -ldflags "-X main.Version $(VERSION) -X main.Build $(BUILD)" -o build/Linux/$(NAME)
	mkdir -p build/Darwin && GOOS=darwin go build -ldflags "-X main.Version $(VERSION) -X main.Build $(BUILD)" -o build/Darwin/$(NAME)
	rm -rf release && mkdir release
	tar -zcf release/$(NAME)_linux_$(ARCH).tgz -C build/Linux $(NAME)
	tar -zcf release/$(NAME)_darwin_$(ARCH).tgz -C build/Darwin $(NAME)

release: test build
	git tag -s -v v$(VERSION)
	git push --tags

test: deps
	go test ./...

all: test
