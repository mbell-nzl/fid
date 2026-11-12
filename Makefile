.PHONY: fmt-check lint test coverage build check-restore check-two-worktree check-policy-from-main

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

# V8 (RESET.md), check:restore: back up, wipe, restore and compare. make test
# runs it too; this runs only it.
check-restore:
	go test -count=1 -run '^TestV8BackupWipeRestoreYieldsIdenticalState$$' ./cmd/fid

# V3 (RESET.md), check:two-worktree: worktrees sharing one ledger never hold
# the same task, and fid next respects dependencies. make test runs it too;
# this runs only it.
check-two-worktree:
	go test -count=1 -run '^TestV3TwoWorktreesNeverHoldTheSameTask$$' ./cmd/fid

# V7 (RESET.md), check:policy-from-main: working-tree, index and local-branch
# edits cannot change policy; only the caller refreshing origin/<default> can.
# make test runs it too; this runs only it.
check-policy-from-main:
	go test -count=1 -run '^TestV7WorkingTreeEditsToPolicyChangeNothingUntilMerged$$' ./cmd/fid

coverage:
	go test -race -gcflags=modernc.org/...=-d=checkptr=0 -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

build:
	go build -o /tmp/fid ./cmd/fid
