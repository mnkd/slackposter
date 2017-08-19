PACKAGES = $(shell go list)

default: test

test:
	go test ${PACKAGES}

vet:
	go vet ${PACKAGES}

lint:
	golint ${PACKAGES}

.PHONY: test vet lint
