# Phase 11 - Task 2: `fid issue claim` And `fid issue unclaim`

## Objective

Wire the internal claims domain service into the public issue command surface.
This task gives users the primitive write commands for issue claim acquisition,
renewal, force-acquire, and release, with Phase 10 permission enforcement and
validate-before-append behavior visible at the CLI.

## Requirements

- PRD refs: PRD-CLAIM-001, PRD-CLAIM-002, PRD-CLAIM-003, PRD-CLAIM-004,
  PRD-CLAIM-008, PRD-CLAIM-008a, PRD-ERR-003, PRD-ERR-006
- TRD refs: TRD-ARCH-010, TRD-ARCH-011; `03_TRD_SCHEMA.md` section 11.1
- Design refs: `06_DESIGN_PLAN_CLI.md` `fid issue claim` and `fid issue unclaim`
- Quality refs: `08_QUALITY_PLAN.md` COV-001/COV-002/COV-003; CF-007
- Decisions: Phase 10 DEC-B; Phase 11 DEC-026

## Validation Contract Coverage

- VAL-002: CLI claim creates a live issue claim with default/custom TTL and note.
- VAL-003: live-claim conflict rejects before append and names holder/expiry.
- VAL-004: stale takeover works from the CLI and records superseded holder.
- VAL-005: CLI force-acquire is gated by `issue.claim.force_acquire`.
- VAL-006: CLI renewal works only for holder and projects activity.
- VAL-007: CLI unclaim releases only the caller's live claim.
- VAL-010: `--by human` can hold an issue claim and read back as human.

## Scope

### Files To Create Or Modify

- `cmd/fid/issue.go` - register `issue claim` and `issue unclaim`; parse
  `--by`, `--ttl`, `--force`, `--renew`, and `--note`; render confirmations.
- `cmd/fid/issue_claim_test.go` (new or equivalent) - CLI coverage for claim and
  unclaim behavior.
- `internal/claims/` - add release support if Task 1 did not include it; expose
  result types needed by CLI confirmations.
- Shared CLI error rendering in `cmd/fid/issue.go` or nearby helpers, if claim
  errors need holder/expiry-specific messages.

### Files To Read For Context

- `internal/claims/` from Task 1.
- Existing command patterns in `cmd/fid/issue.go` (`issue transition`,
  `issue complete`, `issue ac verify`) and their tests.
- `cmd/fid/policy_test.go` for show-vs-enforce parity style.
- `project-docs/docs/phases/phase-11/qa/validation-contract.md`.

### Out Of Scope

- `fid claim list` (Task 3).
- Transition claim requirement and auto-release (Task 4).
- JSON output formats for claim writes; the spec defines stderr confirmations.
- Force-release by a non-holder; v0.1 uses force-acquire only.

## Acceptance Criteria

- [ ] `fid issue claim <id> --by developer` succeeds with stderr matching the
  design shape and creates the expected current claim (VAL-002).
- [ ] `--ttl` parses Go-style durations and writes `ttl_seconds`; invalid TTL is
  rejected before append with a clear error (VAL-002).
- [ ] `--note` persists and reads back through the existing issue claims section
  or domain query (VAL-002).
- [ ] A second live claimant without `--force` exits non-zero, names holder and
  expiry, and appends zero journal lines/rows (VAL-003).
- [ ] A stale claim can be taken by a new holder without force, and a permitted
  `--force` can supersede a live claim; denied force names
  `issue.claim.force_acquire` and appends nothing (VAL-004, VAL-005).
- [ ] `--renew` by the holder updates expiry and activity evidence; non-holder
  renewal rejects before append (VAL-006).
- [ ] `fid issue unclaim <id> --by <holder>` releases the claim; wrong-holder,
  repeat release, missing `--by`, and unknown role reject before append (VAL-007).
- [ ] `fid issue claim <id> --by human` succeeds and reads back with
  `claimed_by = human` (VAL-010).

## TDD Plan

- First failing test: `fid issue claim` is registered and a basic claim writes
  one current claim row with the default TTL.
- Additional tests: TTL/note parsing, conflict error, stale takeover,
  force-acquire allowed/denied, renew allowed/denied, unclaim holder/wrong
  holder, missing/unknown role, human claim read-back.
- Regression risks: command flags must not bypass domain validation; denied
  paths must assert zero new journal lines; confirmation text should remain
  stable enough for operators without overfitting tests to whitespace.

## Documentation Updates

- `project-docs/docs/phases/phase-11/tasks/task-2-done.md` - completion report
  with command results, coverage, and commit hash/subject.
- `project-docs/docs/phases/phase-11/decisions.md` - record any CLI output or
  flag behavior deviation from `06_DESIGN_PLAN_CLI.md`.
- `06_DESIGN_PLAN_CLI.md` - edit only if implementation proves the existing
  command contract wrong; otherwise record verified-no-change in the completion
  report.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None.
- Cross-surface parity impact: validate `fid policy show` claim action rows
  against one permitted and one denied CLI outcome.
- Audit or release-readiness impact: public write commands require build
  verification.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./...
go build -o /tmp/fid ./cmd/fid
```

The full gate above is required before task close because this task adds public
CLI write commands. A focused test command may be used during the red/green loop,
but it does not substitute for the full gate at close. Every required command
must be recorded in the completion report as `PASS`, or as `BLOCKED` with reason
in a blocker report.

Conditional commands:

```text
go test -bench . -benchmem ./...
```

- `go test -bench . -benchmem ./...` - required if benchmark fixtures or journal
  scale helpers are touched; otherwise record N/A with reason.

## Notes For Developer

- `--force` and `--renew` are distinct modes. If both are supplied and the
  domain cannot define a spec-backed behavior from existing docs, stop and
  escalate rather than inventing semantics.
- `--by human` is a real role and should flow through `Authorize`; do not special
  case it in the CLI.
- Pair completion-report commit-hash filling with the task commit or immediate
  reconciliation commit, per the Phase 10 deep-retro note.
