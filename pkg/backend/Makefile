GO ?= GO111MODULE=on CGO_ENABLED=0 go

run:
	$(GO) run main.go

build:
	$(GO) build -o=bootstrap main.go

test:
	$(GO) test ./...