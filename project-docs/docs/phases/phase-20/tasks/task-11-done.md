# Phase 20 - Task 11: Link contract and dependency invariants - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 1 (`review-4:qa-rework:MF-01-typed-link-evidence-matrix`)
- Next action: Independent QA re-review of MF-01.
- Human attention required: no

## Summary of what was implemented

- Migrated `link add`, `link remove`, and `link list` to typed endpoint requests, canonical JSON envelopes, sorted Issue revisions, inverse projection, bounded cycle detection, and stdout-only Human success.
- Review 4 MF-01 added the complete typed-link evidence matrix and corrected three revealed defects: Issue/requirement endpoints sharing an opaque ID now honor their supplied types; a self dependency emits `dependency-cycle` with its path; and add/remove retain the appended Event timestamp at nanosecond precision.

## Files created or modified

- `internal/links/add.go`
- `internal/links/remove.go`
- `internal/links/cycle.go`
- `internal/links/cycle_test.go`
- `internal/links/list_test.go`
- `cmd/fid/link_test.go`
- `cmd/fid/link_limit_test.go`
- `project-docs/docs/phases/phase-20/tasks/task-11-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-20/tasks/task-11-done.md` — refreshed Review 4 MF-01 evidence, repair history, and current exact gate record.

## Validation Contract Coverage

- VAL-008: PASS — `cmd/fid/link_test.go:13` retains Human add stdout-only/empty-stderr behavior; `TestLinkMutationsJSONMatchContract` requires JSON mutation stdout-only success.
- VAL-010: PASS — `cmd/fid/link_test.go:573` requires add/remove sorted two-Issue revisions (2 then 3), and `:593` keeps the journal count unchanged for every rejection.
- VAL-011: PASS — `cmd/fid/link_test.go:703` compares each returned Event ID/timestamp to the journal tail and requires the five-effect Event shape; this caught the original second-precision timestamp defect.
- VAL-017: PASS — `cmd/fid/link_test.go:573` asserts non-null matching mutation Events; `:446` asserts `link.list` has `event:null`.
- VAL-023: PASS — `cmd/fid/link_test.go:573` asserts closed add/remove Link/revision schemas; `:446` asserts typed list count/order/item schemas and inverse directions.
- VAL-024: PASS — `internal/links/cycle_test.go:49` covers self; `:76` covers two-node, three-node, corrupt pre-existing cyclic storage, exact paths, zero append, and diamond success; `cmd/fid/link_test.go:679` decodes the JSON `dependency-cycle` path.
- VAL-025: PASS — `cmd/fid/link_test.go:593` decodes exact `invalid-argument` self-link and `link-not-found` subject/details and proves zero append.
- VAL-041: PASS — `cmd/fid/link_test.go:703` requires Actor `developer`, the closed five-effect shape, and absence of Evidence/Outcome/Claim protected effects; the retained Human tests cover stdout-only confirmations.
- VAL-044: PASS — `cmd/fid/link_test.go:593` covers omitted add/remove `--from-type`, omitted add/remove `--to-type`, and omitted list `--type`, each as exact `invalid-argument` with no append.
- VAL-045: PASS — `cmd/fid/link_test.go:446` verifies list order; `internal/links/list_test.go:111` proves 10,000 success and 10,001 no-truncation domain behavior; `cmd/fid/link_limit_test.go:12` proves both states and typed failure through the CLI with unchanged journal/projection.
- VAL-046: PASS — `cmd/fid/link_test.go:533` proves supplied-type resolution for same-ID Issue/requirement endpoints; `:593` asserts distinct typed missing Issue and requirement codes/details with no append.

## Tests run

- Original TDD RED: `go test ./internal/links -run '^TestAddRejectsThreeNodeDependencyCycle$' -count=1` - `FAIL` (terminated `exit_code: 1`). Initial compile-red named missing `DependencyCycleError`; behavior-red was `error = <nil> <nil>, want *DependencyCycleError` (terminated `exit_code: 1`).
- Original TDD GREEN: `go test ./internal/links -run '^TestAddRejectsThreeNodeDependencyCycle$' -count=1` - `ok  	github.com/mbell-nzl/fid/internal/links	0.069s` (terminated `exit_code: 0`).
- Review 4 RED: `go test ./cmd/fid -run '^TestLinkAddUsesSuppliedTypeForAmbiguousSameIDEndpoints$' -count=1` first found a test-fixture parser mismatch, `requirement id "FID-e97095" not found in "project-docs/specs/link_requirements.md"` (terminated `exit_code: 1`); after seeding the typed requirement directly, its byte-identical rerun found `*links.SelfLinkError ... want JSON success on stdout only` (terminated `exit_code: 1`). Repair: self-link matching includes endpoint type.
- Review 4 GREEN: `go test ./cmd/fid -run '^TestLinkAddUsesSuppliedTypeForAmbiguousSameIDEndpoints$' -count=1` - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.057s` (terminated `exit_code: 0`).
- Review 4 RED: `go test ./internal/links -run '^TestAddRejectsSelfDependencyCycle$' -count=1` - `error = *links.SelfLinkError ... want *DependencyCycleError` (terminated `exit_code: 1`). Repair: same-Issue `blocks`/`blocked_by` returns a two-ID `DependencyCycleError` before append.
- Review 4 GREEN: `go test ./internal/links -run '^TestAddRejectsSelfDependencyCycle$' -count=1` - `ok  	github.com/mbell-nzl/fid/internal/links	0.049s` (terminated `exit_code: 0`).
- Review 4 RED: `go test ./cmd/fid -run '^(TestLinkMutationsJSONMatchContract|TestLinkJSONRequestAndTypedFailuresAreAtomic|TestLinkAddUsesSuppliedTypeForAmbiguousSameIDEndpoints)$' -count=1` - `link.add envelope event = map[string]interface {}... journal tail = journal.Event{...}` (terminated `exit_code: 1`). Repair: add/remove response Event timestamps use `time.RFC3339Nano`.
- Review 4 GREEN: `go test ./cmd/fid -run '^(TestLinkMutationsJSONMatchContract|TestLinkJSONRequestAndTypedFailuresAreAtomic|TestLinkAddUsesSuppliedTypeForAmbiguousSameIDEndpoints)$' -count=1` - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.092s` (terminated `exit_code: 0`).
- Review 4 GREEN: `go test ./internal/links -run '^TestDependencyCycleTraversalMatrix$' -count=1` - `ok  	github.com/mbell-nzl/fid/internal/links	0.262s` (terminated `exit_code: 0`).
- Review 4 GREEN: `go test ./cmd/fid -run '^(TestLinkJSONDependencyCycleHasExactPathAndNoAppend|TestLinkListJSONHonorsAtAndOverLimitWithoutTruncation)$' -count=1` - `ok  	github.com/mbell-nzl/fid/cmd/fid	1.996s` (terminated `exit_code: 0`).
- Review 4 GREEN: `go test ./internal/links ./cmd/fid -count=1` - `ok  	github.com/mbell-nzl/fid/cmd/fid	8.369s` (terminated `exit_code: 0`).
- Real entrypoint: the first fresh `/tmp/fid` probe parsed its Event ID as an Issue ID and returned `"code":"issue-not-found"` (terminated `exit_code: 1`); after correcting only the probe extractor, `/tmp/fid --project <fresh-workspace> --json link add ...` returned `"operation":"link.add"` with two sorted revision-2 endpoints and Event `occurred_at:"2026-08-15T06:05:27.445473919Z"` (terminated `exit_code: 0`).

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 80.5% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)					77.7%` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	4.568s` (terminated `exit_code: 0`)
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)

## Deviations from spec

- None. The Review 4 repair stayed within Task 11 link implementation/tests and conformed to the accepted typed-endpoint, cycle, Event, and collection contracts.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `cf096da6391c020c24e3392996548158271bc080`

## Commit subject

- `fix(links): complete typed link evidence matrix`
