# Phase 7 - Task 4: `all_acs_verified` real-data gating + rebuild parity (phase closure)

## Objective

Prove the Phase 7 write surface drives the existing `all_acs_verified`
transition rule with **real** `fid issue ac verify` data (replacing the Phase 5
journal fixtures), and that AC adds + verifications survive a delete-DB +
`fid rebuild` with identical projected state. This is the Phase 7 end-to-end
closure task.

## Requirements

- PRD refs: PRD-WF-007 (`all_acs_verified`), PRD-AC-001..006 (exercised end-to-end)
- TRD refs: TRD_SCHEMA §2.3; ADR-0001 (journal canonical, DB derived)
- Design refs: `06_DESIGN_PLAN_CLI.md` §3 (`ac` commands), §2 (`issue transition`)
- Quality refs: `08_QUALITY_PLAN.md` integration/rebuild-parity; DEC-007 coverage

## Validation Contract Coverage

- VAL-008: A real `all_acs_verified`-gated transition rejects when required-role passes are missing (naming AC IDs/roles) and succeeds once real `ac verify` passes exist.
- VAL-009: Delete DB + `fid rebuild` reproduces identical `acceptance_criteria` + `ac_verifications` state.

## Scope

### Files To Create Or Modify

- `cmd/fid/e2e_v01c_ac_test.go` (new) - end-to-end AC test over create → ac add → transition-reject → ac verify → transition-pass → delete-DB → rebuild → parity.
- (If needed) a focused integration helper; do **not** modify the rule functions or write paths from Tasks 1–3 except to fix defects this task surfaces.

### Files To Read For Context

- `internal/issues/transition.go`, `internal/issues/transition_rules.go` (`all_acs_verified` loader + rule).
- `internal/workflow/rules.go` (`AllACsVerified`, required-roles parameterisation).
- `cmd/fid/e2e_v01b_test.go` (Phase 6 smoke pattern + rebuild-parity assertions).
- Phase 5 `task-5` / `qa-plan` notes on the fixture-based gating evidence being replaced.

### Out Of Scope

- `has_completion_report` and `concerns_addressed` gating (Phases 8–9).
- New ranking/selection semantics; only the existing `all_acs_verified` rule.
- Held-out scenario harness relocation (separate scheduled commitment; recommend only).

## Acceptance Criteria

- [ ] End-to-end: create a task whose workflow requires `all_acs_verified`, add ≥1 AC, attempt the gated transition → **rejected** with an error naming unmet AC IDs and missing required role(s), and the issue status is unchanged (VAL-008).
- [ ] After `fid issue ac verify` records the required-role pass(es) via the real command, the same transition **succeeds** (VAL-008).
- [ ] The reject→pass evidence uses real `ac add`/`ac verify` commands, not journal fixtures (VAL-008).
- [ ] After exercising adds + verifications, deleting the SQLite DB and running `fid rebuild` reproduces identical `acceptance_criteria` and `ac_verifications` rows and identical latest-per-role read-back (VAL-009).
- [ ] The test asserts no extra journal events are produced by rejected transitions (validate-before-append holds across the integration).

## TDD Plan

- First failing test: gated transition before any pass verification is rejected with unmet AC/role detail and zero state change.
- Additional tests: transition succeeds after real pass verifications; post-rebuild AC + verification parity; rejected transition adds no journal line.
- Regression risks: must use the real default task workflow's required roles; confirm the rule's `params_json` roles match what the test verifies.

## Documentation Updates

- `project-docs/docs/phases/phase-7/` completion report - paste full end-to-end transcript (reject, verify, pass, rebuild parity).
- Note in the completion report whether the Phase 5 fixture-based gating evidence is now superseded by real commands (input for the TL spec-sync audit at wrap).

## Scenario / Parity / Audit Notes

- Held-out scenario impact: this smoke may seed the held-out scenario library; recommend placement only — do not paste hidden scenario content.
- Cross-surface parity impact: AC state parity across `ac list`, `issue get --section acceptance`, and the transition rule loader before vs after rebuild.
- Audit or release-readiness impact: confirms Phase 5 fixture evidence is replaced by real write commands (spec-sync input for wrap).

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

Coverage targets (DEC-007): `cmd/fid` ≥ 70%, `internal/issues`/`internal/workflow` ≥ 80%.

Conditional commands:

```text
go test -tags integration -race ./...
go build -o /tmp/fid ./cmd/fid
```

List any conditional command required by this task and why:

- `go test -tags integration -race ./...` only if the end-to-end test is placed behind the `integration` build tag (match the existing e2e test's tagging convention).
- `go build -o /tmp/fid ./cmd/fid` if pasting real CLI transcripts.

## Notes For Developer

- Mirror the Phase 6 `e2e_v01b_test.go` structure for workspace setup, command invocation, and the delete-DB + `fid rebuild` parity assertion.
- The gated transition must reject **without** appending a journal event (CF-007 / VAL-008); assert journal line count is unchanged across the rejected attempt.
- If you discover the default task workflow does not actually require `all_acs_verified` on the transition under test, stop and raise it (triage) rather than altering the workflow graph to make the test pass.
</content>
