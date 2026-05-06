GOFLAGS := -tags=stdjson,gjson
GOCACHE := /tmp/go-build-cache
GOPATH := /tmp/go

.PHONY: build test run tidy

build:
	GOCACHE="$(GOCACHE)" GOPATH="$(GOPATH)" GOFLAGS="$(GOFLAGS)" go build ./...

test:
	GOCACHE="$(GOCACHE)" GOPATH="$(GOPATH)" GOFLAGS="$(GOFLAGS)" go test ./...

run:
	GOCACHE="$(GOCACHE)" GOPATH="$(GOPATH)" GOFLAGS="$(GOFLAGS)" go run ./cmd/server

tidy:
	GOCACHE="$(GOCACHE)" GOPATH="$(GOPATH)" GOFLAGS="$(GOFLAGS)" go mod tidy
