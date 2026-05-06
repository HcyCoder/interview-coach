GOFLAGS := -tags=stdjson,gjson
GOCACHE := /tmp/go-build-cache
GOPATH := /tmp/go
GOPROXY := off
GOSUMDB := off

.PHONY: build test run tidy

build:
	GOCACHE="$(GOCACHE)" GOPATH="$(GOPATH)" GOPROXY="$(GOPROXY)" GOSUMDB="$(GOSUMDB)" GOFLAGS="$(GOFLAGS)" go build ./...

test:
	GOCACHE="$(GOCACHE)" GOPATH="$(GOPATH)" GOPROXY="$(GOPROXY)" GOSUMDB="$(GOSUMDB)" GOFLAGS="$(GOFLAGS)" go test ./...

run:
	GOCACHE="$(GOCACHE)" GOPATH="$(GOPATH)" GOPROXY="$(GOPROXY)" GOSUMDB="$(GOSUMDB)" GOFLAGS="$(GOFLAGS)" go run ./cmd/server

tidy:
	GOCACHE="$(GOCACHE)" GOPATH="$(GOPATH)" GOPROXY="$(GOPROXY)" GOSUMDB="$(GOSUMDB)" GOFLAGS="$(GOFLAGS)" go mod tidy
