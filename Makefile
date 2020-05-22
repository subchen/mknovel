ROOT    := $(shell pwd)
NAME    := mknovel
VERSION := 2.0.0

LDFLAGS := -s -w \
           -X 'main.BuildVersion=$(VERSION)' \
		   -X 'main.BuildGitBranch=$(shell git describe --all)' \
           -X 'main.BuildGitRev=$(shell git rev-list HEAD --count)' \
           -X 'main.BuildGitCommit=$(shell git describe --abbrev=0 --always)' \
           -X 'main.BuildDate=$(shell date -u -R)'

default:
	@ echo "no default target for Makefile"

pre-install:
	go env -w GOPROXY=https://goproxy.cn,direct
	go get github.com/go-bindata/go-bindata/... && cp ~/go/bin/go-bindata /usr/local/bin/

clean:
	@ rm -rf $(NAME) ./releases ./build

fmt:
	@ go fmt ./...

lint: fmt
	@ go vet ./...

generate:
	cd generator/epub && go-bindata -pkg=epub -nometadata -nomemcopy -ignore=.DS_Store -o=assets.go template/...

test: clean fmt
	@ go test -v ./... $(args)

run: clean fmt
	@ go build -o $(NAME)
	@ ./$(NAME) $(args)

build: \
    build-linux \
    build-darwin \
    build-windows

build-linux: clean fmt generate
	@ CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags "$(LDFLAGS)" -o releases/$(NAME)-$(VERSION)-linux-amd64

build-darwin: clean fmt generate
	@ GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o releases/$(NAME)-$(VERSION)-darwin-amd64

build-windows: clean fmt generate
	@ GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o releases/$(NAME)-$(VERSION)-windows-amd64.exe

md5sum: build
	@ for f in $(shell ls ./releases); do \
		cd $(ROOT)/releases; md5sum "$$f" >> $$f.md5; \
	done

release: md5sum

