.PHONY: fmt-check lint test coverage build

fmt-check:
	gofmt -l cmd internal .github/scripts

lint:
	golangci-lint run ./...
	golangci-lint run .github/scripts/race_packages.go

# Decision G (task 0.7): the race detector runs only in packages whose non-test
# code starts goroutines, found from the source the same way CI finds them.
test:
	python3 .github/scripts/ci_tests.py test

coverage:
	go test -race -gcflags=modernc.org/...=-d=checkptr=0 -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

build:
	go build -o /tmp/fid ./cmd/fid
