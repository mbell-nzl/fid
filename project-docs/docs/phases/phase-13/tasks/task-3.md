# Phase 13 - Task 3: v0.1d Closure Smoke (End-to-End Lifecycle)

## Objective

Prove the v0.1 headline exit criterion: the full gated issue lifecycle runs
end-to-end through real commands and is rebuild-stable. The scripted flow is
create -> claim -> transition into a started-category state (which requires a
live caller claim) -> complete -> unclaim, with permission enforcement on every
write, plus a typed link and a label assignment on the issue. This task adds an
integration test only; it ships no new product behaviour (claim, transition,
complete, link, and label all shipped in earlier phases / Tasks 1-2).

## Requirements

- PRD refs: PRD-CLAIM-005 (started-category requires live claim), PRD-WF-001..010,
  PRD-ISSUE-023, PRD-NFR-008
- TRD refs: ADR-0001 (journal canonical, DB derived), TRD-ARCH-011
- Design refs: roadmap v0.1d Phase 13 exit; v0.1 exit criteria
- Quality refs: `08_QUALITY_PLAN.md` COV-001/COV-002
- Decisions: DEC-028 (held-lock projection freshness; no nested lock)

## Validation Contract Coverage

Assertions live in `project-docs/docs/phases/phase-13/qa/validation-contract.md`.

- VAL-009: the full create -> claim -> transition -> complete -> unclaim
  lifecycle with permission enforcement, a link, and a label runs with one event
  per write, no nested-lock failure, a denied-role rejection somewhere in the
  flow, and `fid rebuild` parity of issues/claims/links/issue_labels.

## Scope

### Files To Create Or Modify

- `cmd/fid/e2e_v01d_closure_test.go` (new, `//go:build integration`) - the
  scripted end-to-end lifecycle smoke and rebuild-parity assertions.
- Reuse existing e2e helpers; add a helper only if an existing one does not
  cover a step (keep helpers in the existing test-helper location).

### Files To Read For Context

- `cmd/fid/e2e_v01d_link_test.go` for the integration-test shape: `t.TempDir()`,
  `runInitCommand`, `createCommandIssue`, `runFIDCommand`, `mustJournalPath`,
  `splitJournalLines`, `assertNoNestedLockFailure`, and rebuild-parity helpers.
- `cmd/fid/e2e_v01b_test.go` and `e2e_v01c_complete_test.go` for the
  create/transition/complete command sequences and claim usage.
- `cmd/fid/claim.go`, `cmd/fid/issue.go` (transition/complete) for exact flag
  spellings used in the scripted flow.

### Out Of Scope

- Any new product command or behaviour change. If the smoke reveals a real
  defect, stop and escalate; do not patch behaviour inside a test-only task.
- New domain packages or migrations.

## Acceptance Criteria

- [ ] An integration test drives, against a temp workspace, the sequence:
  `fid issue create` -> `fid issue claim` -> `fid issue transition` into a
  started-category state -> `fid issue complete` -> `fid issue unclaim`, plus
  `fid label add` + `fid issue label add` and one `fid link add`, all with
  `--by` roles permitted by the matrix (VAL-009).
- [ ] The transition into the started-category state is rejected when no live
  claim is held and accepted when the caller holds the claim (PRD-CLAIM-005),
  asserted in the test (VAL-009).
- [ ] At least one write in the flow is attempted with a denied role and rejected
  before append (non-zero exit, zero new journal lines), proving enforcement is
  live end-to-end (VAL-009).
- [ ] Each successful write appends exactly one journal event and produces no
  nested-lock failure (VAL-009).
- [ ] `fid rebuild` after the flow reproduces identical issue status, claim
  state, `links`, and `issue_labels` rows (VAL-009).

## TDD Plan

- First failing test: the full scripted sequence asserting per-write journal-line
  increments and final issue status.
- Additional assertions: started-category transition rejected without a claim;
  denied-role write rejected before append; rebuild parity of
  issues/claims/links/issue_labels via a fresh DB comparison.
- Regression risks: keep the flow deterministic (fixed roles, no time-dependent
  TTL flake); assert journal-line counts rather than wall-clock timing; do not
  assert on incidental ordering that the projection does not guarantee.

## Documentation Updates

- `project-docs/docs/phases/phase-13/tasks/task-3-done.md` - completion report
  with command results (including `go test -tags integration`), coverage, the
  scripted sequence, and commit hash.
- No spec edits expected; record verified-no-change for the v0.1d exit criteria
  in the completion report unless the smoke surfaces a spec gap.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None.
- Cross-surface parity impact: this task is itself the cross-command parity proof
  for the v0.1d lifecycle (claims + transitions + links + labels under one gate).
- Audit or release-readiness impact: this is the v0.1 headline exit criterion;
  flag in the completion report that v0.1d is closed pending final review.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./...
```

The full task gate above is required before task close. `go test -tags
integration -race ./...` is **required** for this task (not conditional) because
the deliverable is the integration smoke. Every required command must be recorded
in the completion report as `PASS`, or as `BLOCKED` with reason in a blocker
report.

Conditional commands:

```text
go build -o /tmp/fid ./cmd/fid
```

- `go build -o /tmp/fid ./cmd/fid` - run if any non-test wiring changed;
  otherwise record N/A (this task is expected to be test-only).

## Notes For Developer

- This task is test-only. If the smoke fails because of a real product defect,
  stop and escalate via TL-Lite rather than changing behaviour inside the task.
- Depends on Tasks 1-2 (labels) and the already-shipped claim/transition/complete
  commands.
- Keep the test under the `integration` build tag so the default unit gate stays
  fast.
- Completion report commit fields must be filled from `git rev-parse HEAD`.
