include Makefile.env

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -o test main.go

include hack/make-rules/tools.mk
include hack/make-rules/verify.mk

