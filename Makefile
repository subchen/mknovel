ROOT    := $(shell pwd)
NAME    := mknovel
VERSION := $(shell cat VERSION.txt)

GOPATH  := $(ROOT)/../../../../

LDFLAGS := -s -w \
           -X 'main.BuildVersion=$(VERSION)' \
           -X 'main.BuildGitRev=$(shell git rev-list HEAD --count)' \
           -X 'main.BuildGitCommit=$(shell git describe --abbrev=0 --always)' \
           -X 'main.BuildDate=$(shell date -u -R)'

PACKAGES := $(shell go list ./... | grep -v /vendor/)

clean:
	@ rm -rf $(NAME)

fmt:
	@ go fmt $(PACKAGES)

vet:
	@go vet $(PACKAGES)

test: clean
	@ go test -v $(PACKAGES) $(args)

run: clean fmt
	@ go build -o $(NAME)
	@ ./$(NAME) $(args)

build: \
    build-linux \
    build-darwin \
    build-windows

build-linux: clean
	@ GOOS=linux GOARCH=amd64 \
	go build -ldflags "$(LDFLAGS)" -o bin/$(NAME)-linux-$(VERSION)

build-darwin: clean
	@ GOOS=darwin GOARCH=amd64 \
	go build -ldflags "$(LDFLAGS)" -o bin/$(NAME)-darwin-$(VERSION)

build-windows: clean
	@ GOOS=windows GOARCH=amd64 \
	go build -ldflags "$(LDFLAGS)" -o bin/$(NAME)-windows-$(VERSION).exe

install-deps:
	@ glide install

update-deps:
	@ glide update
