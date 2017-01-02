NAME     := slackposter
VERSION  := 0.1.0
REVISION := $(shell git rev-parse --short HEAD)
SRCS     := slackposter.go
LDFLAGS  := -ldflags="-X \"main.Version=$(VERSION)\" -X \"main.Revision=$(REVISION)\""

# Setup
setup:
	go get github.com/Masterminds/glide
	go get github.com/golang/lint/golint
	go get golang.org/x/tools/cmd/goimports

# Install dependencies
deps: setup
	glide install

# Update dependencies
update: setup
	glide update

## Lint
lint: setup
	go vet $$(glide novendor)
	for pkg in $$(glide novendor -x); do \
		golint -set_exit_status $$pkg || exit $$?;\
	done

## Format source codes
fmt: setup
	goimports -w $$(glide nv -x)

# Build binaries ex. make bin/myproj
# bin/%: cmd/%/main.go deps
# 	go build -ldflags "$(LDFLAGS)" -o $@ $<

bin/$(NAME): $(SRCS) fmt
	go build $(LDFLAGS) -o bin/$(NAME)

# Show help
help:
	@make2help $(MAKEFILE_LIST)

clean:
	rm -rf bin/*

.PHONY: setup deps update test lint help
