ROOT    := $(shell pwd)
NAME    := mknovel
VERSION := 2.0.0

GOPATH  := $(ROOT)/../../../../

LDFLAGS := -s -w \
           -X 'main.BuildVersion=$(VERSION)' \
           -X 'main.BuildGitRev=$(shell git rev-list HEAD --count)' \
           -X 'main.BuildGitCommit=$(shell git describe --abbrev=0 --always)' \
           -X 'main.BuildDate=$(shell date -u -R)'

PACKAGES := $(shell go list ./... | grep -v /vendor/)

default:
	@ echo "no default target for Makefile"

clean:
	@ rm -rf $(NAME) ./releases ./build

glide-vc:
	@ glide-vc --only-code --no-tests --no-legal-files

fmt:
	@ go fmt $(PACKAGES)

lint: fmt
	@ go vet $(PACKAGES)

test: clean fmt
	@ go test -v $(PACKAGES) $(ARGS)

run: clean fmt
	@ go build -o $(NAME)
	@ ./$(NAME) $(args)

build: \
    build-linux \
    build-darwin \
    build-windows

build-linux: clean fmt
	@ GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o releases/$(NAME)-$(VERSION)-linux-amd64

build-darwin: clean fmt
	@ GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o releases/$(NAME)-$(VERSION)-darwin-amd64

build-windows: clean fmt
	@ GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o releases/$(NAME)-$(VERSION)-windows-amd64.exe

md5sum: build
	@ for f in $(shell ls ./releases); do \
		cd $(ROOT)/releases; md5sum "$$f" >> $$f.md5; \
	done

release: md5sum

