.PHONY: fmt-check lint test coverage build check-restore check-two-worktree check-policy-from-main check-leak-test check-approval-binding check-done-refusals

FMT_DIRS = cmd internal .github/scripts

# gofmt -l exits 0 even when it lists files, so any listed file fails the check.
fmt-check:
	@files="$$(gofmt -l $(FMT_DIRS))" || exit 1; \
	if [ -n "$$files" ]; then printf 'gofmt would reformat:\n%s\n' "$$files"; exit 1; fi; \
	echo "gofmt: $(FMT_DIRS) formatted"

# golangci-lint caches each result with the path of the worktree that first
# linted it. Once that worktree is removed, lint elsewhere can't read the
# file's nolint comment and reports the issue. A gitignored cache in each
# worktree goes with it.
lint: export GOLANGCI_LINT_CACHE = $(CURDIR)/.golangci-cache
lint:
	golangci-lint run ./...
	golangci-lint run .github/scripts/race_packages.go

# Decision G (task 0.7): the race detector runs only in packages whose non-test
# code starts goroutines, found from the source the same way CI finds them.
test:
	python3 .github/scripts/ci_tests.py test

# V2 (RESET.md), check:leak-test: human output names no kernel words or IDs.
# make test runs it too; this runs only the leak checks.
check-leak-test:
	go test -count=1 -run '^TestV2Human' ./cmd/fid

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

# V6 (RESET.md), check:approval-binding: a real-terminal approval survives a
# clean rebase, becomes stale on an edit, and is renewed by a new decision.
check-approval-binding:
	go test -count=1 -run '^TestV6ApprovalSurvivesCleanRebaseAndStalesOnEdit$$' ./cmd/fid

# V5 (RESET.md), check:done-refusals: public journeys refuse missing receipts,
# open majors and stale approvals; only merged work with green CI completes.
# make test runs it too; this runs only it.
check-done-refusals:
	go test -count=1 -run '^TestV5DoneRefusesWithoutReceiptsAndAcceptsMergedGreenWork$$' ./cmd/fid
