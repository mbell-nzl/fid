# Phase 23 - Task 7: Review freshness and Human Outcome supersession - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 0
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Made Evidence freshness follow the accepted Obligation policy. Repository-tree Review Evidence becomes `stale` after a target tree change, while `freshness_target:"none"` Human Outcome evaluations remain `fresh`, including when the Outcome is absent.
- Added Human Outcome evaluation state and a shared resolver that selects the governing Outcome by immutable `created_event_id` order for the same class and Issue scope. Caller-authored timestamps do not influence supersession.
- A later differing Outcome now displaces an earlier linked approved Outcome with reason `outcome-superseded`; a later approved Outcome governs the reverse sequence. Validity, freshness, and Artifact availability remain independently serialized.
- Added compact seeded-projection tests until Task 10 supplies the Human Outcome writer, plus a real built-`fid` journey that proves Review staling, both Outcome directions, and unchanged Review and Outcome readback.

## Files created or modified

- `internal/evidence/evidence.go` and `internal/evidence/evidence_task7_test.go`
- `internal/workflow/evaluation.go`, `internal/workflow/resolve.go`, and `internal/workflow/evaluation_task7_test.go`
- `cmd/fid/integration/p22_fake_test.go` and `cmd/fid/integration/p23_freshness_supersession_test.go`
- `project-docs/docs/phases/phase-23/tasks/task-7-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-23/tasks/task-7-done.md`

## Validation Contract Coverage

- VAL-031: PASS - `TestReviewEvidenceStalesOnTargetTreeChange` in `internal/evidence/evidence_task7_test.go:16` proves repository-tree freshness becomes `stale`. `TestP23FreshnessAndOutcomeSupersessionRealProcessJourney` in `cmd/fid/integration/p23_freshness_supersession_test.go:12` invokes the built `fid` binary, observes `valid/stale/not-applicable` with `evidence-stale`, and compares the complete Review projection before and after evaluation without a changed field.
- VAL-032: PASS - `TestAbsentHumanOutcomeFreshnessRemainsFresh` in `internal/workflow/evaluation_task7_test.go:143` proves an absent `freshness_target:"none"` Outcome remains wire value `fresh`; `TestHumanOutcomeEvaluationKeepsFreshnessAndUsesSupersession` in `internal/evidence/evidence_task7_test.go:29` proves a linked Outcome remains `valid/fresh/not-applicable` across a changed tree. The built-binary journey at `cmd/fid/integration/p23_freshness_supersession_test.go:35` independently observes the linked release Outcome satisfied and fresh against that changed tree.
- VAL-033: PASS - `TestHumanOutcomeSupersessionFollowsRecordedEventOrderBothDirections` in `internal/workflow/evaluation_task7_test.go:12` uses timestamps ordered opposite to Event IDs and selects rejected after approved and approved after rejected. `TestHumanOutcomeEvaluationUsesLaterDifferingResult` at line 95 proves the first direction returns `outcome-superseded` and the reverse is satisfied. The built-binary journey at `cmd/fid/integration/p23_freshness_supersession_test.go:42` exercises both directions through real `fid evidence evaluate` calls and re-reads every seeded Outcome field unchanged.
- VAL-079: PASS - the exact ordinary race gate reported `cmd/fid` at 499.695s and `internal/workflow` at 211.628s, both below 600s. The exact tagged integration gate completed in 648.604s, below 900s.

## Tests run

- Required first RED: `go test ./internal/evidence -run '^TestReviewEvidenceStalesOnTargetTreeChange$' -count=1` - FAIL - `FAIL` after `undefined: freshnessForTarget` (terminated `exit_code: 1`).
- Required first GREEN: the byte-identical command above - PASS - `ok  github.com/mbell-nzl/fid/internal/evidence 0.004s` (terminated `exit_code: 0`).
- Event-order RED: `go test ./internal/workflow -run '^TestHumanOutcomeSupersessionFollowsRecordedEventOrderBothDirections$' -count=1` - FAIL - `FAIL` after `undefined: resolveGoverningHumanOutcome` (terminated `exit_code: 1`).
- Event-order GREEN: the byte-identical command above after adding the shared resolver - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow 0.123s` (terminated `exit_code: 0`).
- Outcome evaluation RED: `go test ./internal/evidence -run '^TestHumanOutcomeEvaluationKeepsFreshnessAndUsesSupersession$' -count=1` - FAIL - `FAIL` because the later rejected Outcome still evaluated as satisfied (terminated `exit_code: 1`).
- Outcome evaluation GREEN: the byte-identical command above after loading governing Outcome state - PASS - `ok  github.com/mbell-nzl/fid/internal/evidence 0.008s` (terminated `exit_code: 0`).
- Human no-target RED: `go test ./internal/workflow -run '^TestAbsentHumanOutcomeFreshnessRemainsFresh$' -count=1` - FAIL - `FAIL` because absent Human Outcome freshness was `stale` (terminated `exit_code: 1`).
- Human no-target GREEN: the byte-identical command above - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow 0.126s` (terminated `exit_code: 0`).
- Focused package suite: `go test ./internal/evidence ./internal/workflow` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow 20.644s` (terminated `exit_code: 0`).
- Focused real-process journey: `go test -tags integration ./cmd/fid/integration -run '^TestP23FreshnessAndOutcomeSupersessionRealProcessJourney$' -count=1 -timeout 5m` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 2.329s` (terminated `exit_code: 0`).

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow 211.628s coverage: 64.0% of statements` (terminated `exit_code: 0`)
- Ordinary package budgets: `cmd/fid` was 499.695s and `internal/workflow` was 211.628s; every reported ordinary package remained below 600s.
- `go tool cover -func=coverage.out` - PASS - `total: (statements) 64.5%` (terminated `exit_code: 0`)
- `go test -tags integration -race -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 648.604s` (terminated `exit_code: 0`)
- Tagged integration budget: `cmd/fid/integration` was 648.604s, below 900s.

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- No retained assertion or fixture was displaced; `p22_fake_test.go` was extended only with opt-in Human Outcome Obligations for the Task 7 built-binary journey.
- Human Outcome rows are explicitly seeded projection fixtures because Task 10 owns the writer, authority, and provenance path. No temporary writer or alternate freshness mechanism was added.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `e4fe99111e035ebd72c547eda3a9dbf828cafa14`

## Commit subject

- `feat(evidence): enforce freshness and outcome supersession`
