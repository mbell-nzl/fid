# Phase 11 - Task 4: Claim And Transition Integration

## Objective

Integrate claims with issue workflow transitions: entering a started-category
state requires a live claim held by the caller, and leaving a started-category
state automatically releases that caller's claim in the same journal event.

## Requirements

- PRD refs: PRD-CLAIM-005, PRD-CLAIM-006, PRD-CLAIM-008
- TRD refs: `03_TRD_SCHEMA.md` section 11.1; TRD-ARCH-004; TRD-ARCH-005
- Design refs: roadmap v0.1d Phase 11; `06_DESIGN_PLAN_CLI.md` `fid issue transition`
- Quality refs: `08_QUALITY_PLAN.md` COV-001/COV-002/COV-003; CF-007

## Validation Contract Coverage

- VAL-010: human-held claims block agents because agents do not hold the live claim.
- VAL-011: started-category transitions require a live caller claim.
- VAL-012: transitions out of started auto-release the caller's claim in the same event.

## Scope

### Files To Create Or Modify

- `internal/issues/transition.go` and tests - enforce live caller claim before
  transitions into started-category states; append claim-release effect when
  transitioning out of started-category states.
- `internal/claims/` - expose narrow helpers for live-claim lookup and release
  effect construction if needed; avoid duplicating claim semantics in
  `internal/issues`.
- `cmd/fid/issue_transition_test.go` - CLI coverage for no-claim, wrong-holder,
  stale-holder, matching-holder, and auto-release paths.
- Journal/project tests only if a new combined status+claim event shape needs
  focused projection coverage.

### Files To Read For Context

- `internal/workflow/workflow.go` status categories.
- `internal/issues/transition_rules.go` and `internal/issues/transition.go`.
- `internal/claims/` from Tasks 1-2.
- Existing transition tests and Phase 10 permission retrofit tests.

### Out Of Scope

- Claim acquire/release command semantics already delivered by Tasks 1-2.
- `fid dev start`, `fid dev finish`, `fid dev block`, and agent runs (v0.1e+).
- Ready queue claimed/stale behavior (v0.2).
- New workflow graph states or edges.

## Acceptance Criteria

- [ ] Transition into a started-category state with no claim rejects before
  status append and leaves status unchanged (VAL-011).
- [ ] Transition into started with a live claim held by another role rejects
  before append and names the claim holder/expiry or otherwise gives an
  actionable claim error (VAL-011).
- [ ] Transition into started with a stale claim held by the caller rejects
  before append until the caller renews or reacquires the claim (VAL-011).
- [ ] Transition into started with a matching live caller claim succeeds and
  still enforces existing workflow rules and permission action
  `issue.transition.<target-state>` (VAL-011).
- [ ] A human-held claim prevents a developer/agent transition into started
  because the caller is not the holder (VAL-010).
- [ ] Transition out of a started-category state appends one journal event that
  updates issue status and releases the caller's claim atomically; read-back
  shows the new status and no current claim (VAL-012).
- [ ] Reject paths append zero journal lines/rows and do not delete claims
  (VAL-011).

## TDD Plan

- First failing test: transition `todo -> in_progress` by `developer` without a
  live developer claim now rejects before append.
- Additional tests: wrong-holder, stale-holder, matching-holder success, human
  holder blocks developer, transition out auto-release, existing invalid-edge
  and blocked-transition tests still pass, projection of combined status+claim
  event is atomic.
- Regression risks: preserve existing workflow rule errors where the claim gate
  is not the relevant failure; avoid a second journal event for auto-release;
  do not put output formatting or flag parsing into domain code.

## Documentation Updates

- `project-docs/docs/phases/phase-11/tasks/task-4-done.md` - completion report
  with command results, coverage, and commit hash/subject.
- `project-docs/docs/phases/phase-11/decisions.md` - record any ordering
  decision if claim gating and workflow-rule gating interact in a way not
  already specified.
- `06_DESIGN_PLAN_CLI.md` or `05_DESIGN_PLAN.md` - edit only if implementation
  proves the transition/claim integration contract wrong; otherwise record
  verified-no-change in the completion report.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None.
- Cross-surface parity impact: transition claim errors should be consistent with
  claim command conflict errors where they describe holder/expiry.
- Audit or release-readiness impact: this is the phase's cross-subsystem task;
  QA Reviewer should review it as its own window.

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

The full gate above is required before task close because this task changes an
existing public write path and cross-cuts `internal/issues`, `internal/claims`,
workflow state categories, journal effects, and CLI behavior. A focused test
command may be used during the red/green loop, but it does not substitute for
the full gate at close. Every required command must be recorded in the
completion report as `PASS`, or as `BLOCKED` with reason in a blocker report.

Conditional commands:

```text
go test -bench . -benchmem ./...
```

- `go test -bench . -benchmem ./...` - required if transition benchmark
  fixtures or journal scale helpers are touched; otherwise record N/A with
  reason.

## Notes For Developer

- If the order between claim-gate errors and existing workflow-rule errors is
  not already clear from code and specs, stop and escalate before locking in a
  public error precedence.
- Auto-release must be part of the transition event, not a follow-up event.
- Keep the release-effect construction reusable, but do not introduce a broad
  abstraction unless it removes real duplication.
- Pair completion-report commit-hash filling with the task commit or immediate
  reconciliation commit.
