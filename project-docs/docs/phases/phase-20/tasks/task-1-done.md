# Phase 20 - Task 1: Response contract foundation - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 1
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Added the domain-neutral `internal/cli` `fid-cli.v1` Response Envelope writer for success and failure payloads using `encoding/json`.
- Added the closed Phase 20 reason-code registry, with code-owned category, retryability, allowed subject types, exact detail-object schemas, and closed-details redaction validation.
- Kept the package free of Cobra and domain-package imports; no command migration, project discovery, flag change, or domain-error mapping was introduced.

## Files created or modified

- `internal/cli/envelope.go`
- `internal/cli/envelope_test.go`
- `internal/cli/codes.go`
- `internal/cli/codes_test.go`
- `project-docs/docs/phases/phase-20/tasks/task-1-done.md`

## Documentation updated

- None required.

## Validation Contract Coverage

- VAL-001: Task-1 foundation PASS — `TestWriteSuccessUsesExactContractKeys` at `internal/cli/envelope_test.go:9` decodes the success output and proves exactly the five required keys, literal contract version, `ok:true`, and `event:null`; its assertion at `internal/cli/envelope_test.go:17` proves the trailing newline. Real-command stdout/channel/exit evidence remains owned by Task 3.
- VAL-002: Task-1 foundation PASS — `TestWriteFailureUsesExactContractKeys` at `internal/cli/envelope_test.go:33` proves exactly the four failure-envelope keys and six error keys, with registry-derived `issue-not-found` metadata; `TestWriteFailurePreservesRequiredNullFields` at `internal/cli/envelope_test.go:75` proves required `subject:null` and nullable `details.value:null` are emitted. `TestReasonCodeRegistryIsClosedAndComplete` at `internal/cli/codes_test.go:8` proves all 23 registered codes and their category/retryability/subject/details/redaction metadata, while `TestNewFailureRejectsInvalidSubjectDetailsAndRedaction` at `internal/cli/codes_test.go:69` proves unknown, wrong-subject, missing-detail, and attempted secret-detail cases reject. Real-command stderr/channel/exit evidence and domain mappings remain owned by Tasks 3-13.

## Tests run

Historical Task 1 test evidence is preserved below. Its original direct-test and
unchanged-surface PASS records predate the required termination syntax and are
not used as current gate evidence.

- `go test ./internal/cli -run '^TestWriteSuccessUsesExactContractKeys$'` — EXPECTED RED — `FAIL`
- `go test ./internal/cli` — EXPECTED RED — `FAIL`
> - `go test ./internal/cli` — PASS — `ok  	github.com/mbell-nzl/fid/internal/cli	0.003s`
> - `go test ./internal/cli` — PASS — `ok  	github.com/mbell-nzl/fid/internal/cli	0.002s`
> - `/tmp/fid --help` — PASS (unchanged-surface check) — `Use "fid [command] --help" for more information about a command.` The CLI has no Task 1 envelope integration by design; this task changed no live Cobra surface.
- QA-plan applicability: QA-01 applies through the five exact task gates below. QA-02 and QA-21 apply to this task's deterministic unit foundation through the envelope, registry, null-field, and redaction tests cited above; their command-channel, path-containment, and all-leaf evidence remains assigned to later tasks. The QA plan schedules integration from Task 2 and benchmarks only for Tasks 5, 6, and 18, so neither is applicable to Task 1.

## Exact quality-gate results

> - `gofmt -l cmd internal` — PASS — `exit 0 (no output)`
> - `golangci-lint run ./...` — PASS — `0 issues.`
> - `go test -race -coverprofile=coverage.out ./...` — UNVERIFIED — see Gate-evidence correction (DEC-131).
> - `go tool cover -func=coverage.out` — PASS — `total:										(statements)				74.6%`
> - `go build -o /tmp/fid ./cmd/fid` — PASS — `exit 0 (no output)`

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 80.6% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)					78.7%` (terminated `exit_code: 0`)
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)

## Deviations from spec

- Historical Task 1 implementation had no product or specification deviations.
- DEC-131 gate-evidence correction remains immutable historical context: the
  original full-suite record was UNVERIFIED because its silent token could not
  prove termination. This continuation preserves that record and adds a fresh,
  session-captured PASS result without rewriting it.
- MF-01 report/control-conformance repair: the initial exact report-lint
  attempt failed with `FAIL report-lint found 19 failure(s)` (terminated
  `exit_code: 19`) because the report lacked the canonical typed workflow
  outcome and used legacy evidence punctuation/termination syntax. This
  report-only repair adds the required typed fields, preserves historical
  evidence verbatim as quoted context, and records fresh terminal evidence for
  every exact Task 1 command.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `ab4ac8f288b3f180b698bab8924b92f375d09eb0`

## Commit subject

- `feat(cli): add response contract foundation`

## Gate-evidence correction (DEC-131)

- The full-suite `go test -race -coverprofile=coverage.out ./...` record is UNVERIFIED. Its original recorded token `exit 0 (no output)` could not have come from a completed `go test` run: long runs may print nothing until their packages finish, while a killed or timed-out run is BLOCKED, never PASS with a silent token.
- Preserved original record: `go test -race -coverprofile=coverage.out ./...` — PASS — original recorded token `exit 0 (no output)`.
- `task-3c-done.md`'s current-tree full-gate run is the Phase 20 verified baseline; this historical command was not re-run at Task 1's commit.
- Task 1 remains complete, but its own full-suite gate evidence is UNVERIFIED as of its task commit.
