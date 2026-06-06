.PHONY: fmt-check lint test coverage build

fmt-check:
	gofmt -l cmd internal

lint:
	golangci-lint run ./...

test:
	go test -race ./...

coverage:
	go test -race -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

build:
	go build -o /tmp/fid ./cmd/fid
