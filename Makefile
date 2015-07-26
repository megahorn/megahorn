.PHONEY: deps build test release all
.DEFAULT_GOAL := all

NAME=webminal
ARCH=$(shell uname -m)
VERSION=$(shell cat VERSION)
BUILD=$(shell git rev-parse --short HEAD)

deps:
	go get ./...

build: deps
	go get github.com/mitchellh/gox
	gox -build-toolchain -os="linux darwin" -arch="amd64"
	mkdir -p build/Linux && GOOS=linux go build -ldflags "-X main.Version $(VERSION) -X main.Build $(BUILD)" -o build/Linux/$(NAME)
	mkdir -p build/Darwin && GOOS=darwin go build -ldflags "-X main.Version $(VERSION) -X main.Build $(BUILD)" -o build/Darwin/$(NAME)
	rm -rf release && mkdir release
	tar -zcf release/$(NAME)_linux_$(ARCH).tgz -C build/Linux $(NAME)
	tar -zcf release/$(NAME)_darwin_$(ARCH).tgz -C build/Darwin $(NAME)

release:
ifeq ($(shell git diff --shortstat 2> /dev/null | tail -n1),)
	git tag -f -s v$(VERSION)
	git push --tags --force
	git push --all
	github-changes -o webminal -r webminal --use-commit-body --no-merges
	git add CHANGELOG.md
	git commit -m "Update CHANGELOG.md"
else
	@echo "Please cleanup working directory." && exit 1
endif

test: deps
	go test ./...

all: test
