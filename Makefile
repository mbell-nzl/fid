.PHONY: fmt-check lint test coverage build

FMT_DIRS = cmd internal .github/scripts

# gofmt -l exits 0 even when it lists files, so any listed file fails the check.
fmt-check:
	@files="$$(gofmt -l $(FMT_DIRS))" || exit 1; \
	if [ -n "$$files" ]; then printf 'gofmt would reformat:\n%s\n' "$$files"; exit 1; fi; \
	echo "gofmt: $(FMT_DIRS) formatted"

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
