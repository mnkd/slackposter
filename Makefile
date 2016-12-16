NAME     := slackposter
VERSION  := 0.0.1
REVISION := $(shell git rev-parse --short HEAD)
SRCS     := slackposter.go
LDFLAGS  := -ldflags="-X \"main.Version=$(VERSION)\" -X \"main.Revision=$(REVISION)\""

bin/$(NAME): $(SRCS) format
	go build $(LDFLAGS) -o bin/$(NAME)

.PHONY: format
format:
	go fmt $(SRCS)

.PHONY: clean
clean:
	rm -rf bin/*

.PHONY: install
install:
	go install $(LDFLAGS)
