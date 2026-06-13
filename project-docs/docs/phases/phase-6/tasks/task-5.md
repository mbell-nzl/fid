# Phase 6 - Task 5: end-to-end v0.1b smoke + held-out scenario recommendation

## Objective

Prove v0.1b closure with an end-to-end smoke that exercises the full v0.1b
command surface — create → list → get → update → transition → ready — against a
real `.fid/` workspace, with consistent read-back at each step and `fid rebuild`
parity. Also produce a short recommendation on whether this smoke should become
the first held-out scenario harness input (feeding the standing
held-out-scenario scheduled commitment).

## Requirements

- PRD refs: PRD-ISSUE-001–012, 015, 021; PRD-WF-001–010, 013 (the v0.1b surface
  under test)
- TRD refs: TRD §4.3 (journal → projection → rebuild parity)
- Design refs: `07_DESIGN_PLAN_FLOWS.md` §1 (developer task loop / ready flow)
- Quality refs: `08_QUALITY_PLAN.md` §3 (held-out scenarios), §5.2 (phase QA)

## Validation Contract Coverage

- VAL-011: end-to-end v0.1b smoke over create → list → get → update →
  transition → ready produces consistent read-back at each step, and
  delete-DB + `fid rebuild` reproduces identical projected state.

## Scope

### Files To Create Or Modify

- `cmd/fid/e2e_v01b_test.go` (integration-tagged, or an equivalent integration
  location consistent with existing patterns) — a scripted sequence over a temp
  `.fid/` workspace asserting state at each step plus post-rebuild parity.
- `project-docs/docs/phases/phase-6/held-out-recommendation.md` — short TL/QA
  recommendation: should this smoke seed the first held-out scenario, where
  would it live, and how are implementation agents prevented from loading it.

### Files To Read For Context

- `cmd/fid/init_test.go`, `cmd/fid/rebuild_test.go`,
  `cmd/fid/issue_transition_test.go` (existing CLI/integration test harness
  patterns and temp-workspace helpers)
- `project-docs/specs/08_QUALITY_PLAN.md` §3 (held-out scenario policy)
- `project-docs/docs/current/decisions_active.md` (Scheduled Commitments:
  held-out scenario location)

### Out Of Scope

- v0.1c+ commands (ac/complete/concern), claims, links — not part of the v0.1b
  surface.
- Actually creating/moving a held-out scenario repository (recommendation only;
  the location decision is a separate scheduled commitment).

## Acceptance Criteria

- [ ] The smoke runs `fid init` → `fid issue create` (epic + tasks/subtask) →
  `fid issue list` → `fid issue get` → `fid issue update` → `fid issue
  transition` → `fid ready`, asserting expected read-back at each step
  (VAL-011).
- [ ] After the sequence, deleting the SQLite DB and running `fid rebuild`
  reproduces byte-identical projected state for the touched issues (VAL-011).
- [ ] `fid ready` within the smoke returns the expected scheduling answer for at
  least one sequential and one prioritised epic (VAL-011, supports VAL-002/004).
- [ ] `held-out-recommendation.md` records a clear recommend / defer decision
  with rationale and the isolation mechanism if recommended.

## TDD Plan

- First failing test: the end-to-end sequence test (it fails until all v0.1b
  commands behave consistently together).
- Additional tests: post-rebuild parity assertion; an empty/all-done `fid ready`
  case to exercise the empty-state contract end-to-end.
- Regression risks: test flakiness from non-deterministic ordering or timestamps
  — pin inputs; reuse deterministic fixtures/helpers from existing tests.

## Documentation Updates

- `project-docs/docs/phases/phase-6/held-out-recommendation.md` — new
  recommendation artifact (see Scope).
- Note in the completion report whether AGENTS.md "End-to-end tests: N/A" should
  change once a harness is accepted (flag for TL at wrap; do not edit AGENTS.md
  here unless the held-out decision is accepted).

## Scenario / Parity / Audit Notes

- Held-out scenario impact: this task explicitly evaluates seeding the first
  held-out scenario — keep any hidden scenario content out; this smoke uses only
  in-repo fixtures.
- Cross-surface parity impact: exercises create/list/get/update/transition/ready
  together; rebuild parity is the cross-surface invariant.
- Audit or release-readiness impact: feeds the standing held-out-scenario
  scheduled commitment and the v0.1b exit criteria.

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

This task adds an integration-tagged test, so `go test -tags integration -race
./...` is a **required** command for task close (not optional).

Conditional commands:

- `go build -o /tmp/fid ./cmd/fid`: build the binary for the scripted smoke /
  manual confirmation.

## Notes For Developer

- This is the v0.1b closure task: prefer asserting observable command output and
  read-back over internal state. The smoke is the cross-command safety net.
- Keep the smoke deterministic and self-contained in a temp workspace; do not
  depend on developer machine state or network.
- The held-out recommendation is a judgment artifact for the Tech Lead and QA
  Lead — keep it short (one page) and concrete about the isolation mechanism.
