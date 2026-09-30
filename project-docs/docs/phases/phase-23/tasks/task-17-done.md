# Phase 23 - Task 17: Final capability matrix and availability audit - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 0
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Qualified `reviews.coverage.v1` and `human.outcomes.v1`, and set `attention.v1` `human_outcome_resolution:true` in discovery and the current export profile, with the exact section-15 limits and trust sources. The Phase 21 legacy archive profile remains its historical snapshot. `work.claims.v1` and `telemetry.trajectory.v1` remain false with Phase 24 and 25 ownership.
- The shared fake matrix now compares all twelve rows with built-fid discovery and the export manifest. The same journey checks the `resume` next action's required capability. Retained CLI and Attention tests now assert the qualified rows, while removed V1 leaves still reject without a journal or projection change.
- Task 15 removed `internal/issues/capability_unavailable.go`; `TestV1JudgmentProjectionIsAbsentAfterReplacement` remains as a replacement-schema absence assertion. A current Go-source search for `phase-23-unavailable`, `V1JudgmentUnavailable`, and `judgmentUnavailable` in the live judgment packages found no match. The only live `writeCapabilityUnavailable` call is maintenance vacuum for `telemetry.trajectory.v1` (Phase 25).

## Files created or modified

- `internal/workflow/capabilities.go`; `internal/workflow/stop_conditions_test.go`.
- `internal/export/archive_profiles.go`.
- `cmd/fid/capabilities_test.go`; `cmd/fid/integration/p22_capability_matrix_test.go`; `cmd/fid/integration/p23_attention_resolution_test.go`.
- `project-docs/specs/02_TRD.md`; this report.

## Documentation updated

- `project-docs/specs/02_TRD.md` section 4 records both qualifications, the Attention limit, and the two deferred owners.
- This report cites the prerequisite closure and final command evidence. The Runner-owned phase ledger was left untouched.

## Validation Contract Coverage

- VAL-047: PASS - Phase 23 task 11's report, `task-11-done.md:34`, records built-binary `human_outcome_resolution:false` before task 17 and the working resolution path in `TestP23AttentionResolutionHumanOutcomeReadinessRejectsWrongClassAndWrongResultAndLeavesUnrelatedStopsOpen`. The retained real-process `TestP23AttentionResolutionThroughReviewRecordedAndFindingResolvedRealProcess` now checks `true` after exercising Review and Finding resolution (`cmd/fid/integration/p23_attention_resolution_test.go:118`); the direct task-17 CLI JSON output also contained `"human_outcome_resolution":true`.
- VAL-066: PASS - `TestFinalCapabilityMatrixMatchesAcceptedTable` (`cmd/fid/integration/p22_capability_matrix_test.go:12`) compares the twelve-row fake fixture at lines 43-54 exactly with built-fid discovery and the export manifest, then checks the resume `run.observe` required capability. `TestCapabilitiesJSONGoldenAndProjectionRows` checks the exact CLI JSON and twelve materialized rows. Both false rows and their owner phases are fixed in the shared fixture at lines 52-53 and observed in the direct CLI output.
- VAL-079: PASS - the exact ordinary gate reported `cmd/fid 236.172s` and `internal/workflow 90.207s` (each under 600s); the exact tagged gate reported `cmd/fid/integration 234.507s` (under 900s). None reached 80% of its budget.

Prerequisite closure: task 4 proved immutable Review coverage and independence; task 5 proved Finding dispositions; task 6 proved Review Evidence and review-cycle charging; task 7 proved freshness and Outcome supersession; task 8 proved Concerns and Deviations; task 9 proved Decision References; task 10 proved activated Human Outcome provenance and linkage; task 11 proved Stop and Attention resolution; task 12 proved Requirement References; task 13 proved CLI contract and denial behavior; task 14 proved Human projections and resume; task 15 proved V1 removal and replacement; task 16 proved portable schema-18 recovery. Each `task-4-done.md` through `task-16-done.md` reports `Task outcome: COMPLETE` and names its executing tests. In particular, task 15's VAL-061/063 evidence includes `TestLatestSchemaHasNoV1JudgmentTables`, `TestV1JudgmentProjectionIsAbsentAfterReplacement`, removed-leaf real-process probes, and maintenance checks; task 16's VAL-064/065 evidence covers export, validate, import, and exact schema acceptance.

Per-row provenance and execution: the Phase 21 base rows were introduced at `6ff90152`; Phase 22 qualification rows at `ccc15f66` (`Refs: phase-22/task-17`). Every row below is asserted by `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^(TestFinalCapabilityMatrixMatchesAcceptedTable|TestP23AttentionResolutionThroughReviewRecordedAndFindingResolvedRealProcess)$' ./cmd/fid/integration/` and by the mandatory full tagged gate. `TestFinalCapabilityMatrixMatchesAcceptedTable` checks exact availability, limits, schema bounds, and trust source in discovery and manifest for each row; the direct CLI probe below confirms the binary's visible state.

| Row | Accepted provenance and task-17 disposition |
|---|---|
| `attention.v1` | Phase 22 qualified base (`ccc15f66`); task 11 resolution path; task 17 flips only `human_outcome_resolution`. |
| `evidence.freshness.v1` | Phase 21 base (`6ff90152`); unchanged. |
| `evidence.gate-definitions.v1` | Phase 21 base (`6ff90152`); unchanged. |
| `evidence.receipts.v1` | Phase 21 base (`6ff90152`); unchanged. |
| `export.portable.v1` | Phase 22 qualification (`ccc15f66`) and task 16 schema-18 parity; unchanged. |
| `human.outcomes.v1` | Phase 21 unavailable row (`6ff90152`); tasks 7 and 10 prove supersession, provenance, and recording; task 17 qualifies it. |
| `kit.activation.v1` | Phase 22 qualification (`ccc15f66`); unchanged. |
| `reviews.coverage.v1` | Phase 21 unavailable row (`6ff90152`); tasks 4-6 prove independent Review, Findings, and Evidence; task 17 qualifies it. |
| `telemetry.trajectory.v1` | Phase 21 deferred row (`6ff90152`); remains false, owner Phase 25. |
| `work.claims.v1` | Phase 22 deferred row (`ccc15f66`); remains false, owner Phase 24. |
| `workflow.state.v1` | Phase 21 base (`6ff90152`); unchanged. |
| `workflow.transitions.v1` | Phase 21 base (`6ff90152`); unchanged. |

Retained assertion provenance and convergence: `cmd/fid/capabilities_test.go` exact golden and qualification assertions and the integration fake test came from `ccc15f6` (phase-22/task-17); the removed-leaf/no-state-change probe came from `401f609` (phase-21/task-3); the Attention false assertion came from `c5bf0f8` (phase-23/task-11); the unavailable-capability Stop test came from `f973bca` (phase-22/task-2). This task's accepted matrix changes the first four expected states. The Stop test keeps its unavailable-capability subject and strength, now using still-false `telemetry.trajectory.v1`. All run under the final ordinary or tagged race commands below.

## Tests run

- First RED: `go test -tags integration -run '^TestFinalCapabilityMatrixMatchesAcceptedTable$' -count=1 ./cmd/fid/integration/` - FAIL - `FAIL` (terminated `exit_code: 1`); discovery still reported `human_outcome_resolution:false` and both Phase 23 rows false against the final fixture.
- First GREEN: the byte-identical focused command - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	2.034s` (terminated `exit_code: 0`).
- Retained CLI focused: `go test -run '^TestCapabilities(JSONGoldenAndProjectionRows|AdvertiseOnlyFinalQualifiedRows|RejectsOverridesAndKeepsDeferredRowsUnavailable|RemovedLeavesRejectWithoutSideEffectsAndDeferredRowsStayUnavailable)$' -count=1 ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.192s` (terminated `exit_code: 0`).
- Retained Attention focused: `go test -tags integration -run '^TestP23AttentionResolutionThroughReviewRecordedAndFindingResolvedRealProcess$' -count=1 ./cmd/fid/integration/` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	3.411s` (terminated `exit_code: 0`).
- First touched-package race: `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./internal/workflow ./internal/export ./cmd/fid` - FAIL - `FAIL` (terminated `exit_code: 1`); `TestStopResolutionChecksCurrentActivationAgreementAndCapability/unavailable_capability` used the newly available Human Outcome row.
- Focused Stop repair: `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestStopResolutionChecksCurrentActivationAgreementAndCapability/unavailable_capability$' ./internal/workflow/` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	3.181s` (terminated `exit_code: 0`).
- Touched-package race rerun, byte-identical to the failed command - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	231.658s` (terminated `exit_code: 0`); `internal/workflow 85.367s` and `internal/export 117.755s` also passed.
- Focused tagged race: `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^(TestFinalCapabilityMatrixMatchesAcceptedTable|TestP23AttentionResolutionThroughReviewRecordedAndFindingResolvedRealProcess)$' ./cmd/fid/integration/` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	4.974s` (terminated `exit_code: 0`).
- `go build -o /tmp/fid-task17 ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`).
- `/tmp/fid-task17 --project /tmp/fid-task17-cli-hHF8rK capabilities` - PASS - `workflow.transitions.v1	true` (terminated `exit_code: 0`); the same output listed `attention.v1`, `human.outcomes.v1`, and `reviews.coverage.v1` as true, and `work.claims.v1` and `telemetry.trajectory.v1` as false. The preceding JSON invocation showed the exact new limits and trust sources.
- Source absence observation: `rg -n 'phase-23-unavailable|V1JudgmentUnavailable|judgmentUnavailable' internal/workflow internal/review internal/issues cmd/fid --glob '*.go' --glob '!**/*_test.go'` found no match (`exit_code: 1`, no output). This is an inventory observation, not a declared non-zero PASS gate. `rg -n 'writeCapabilityUnavailable\(' cmd/fid --glob '!**/*_test.go'` found only the helper and the Phase-25 maintenance call (terminated `exit_code: 0`).

## Exact quality-gate results

All five commands ran in order from `/home/mark/projects/fid` on the final product tree. The last printed line of each command is quoted below.

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	90.207s	coverage: 62.1% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:											(statements)			60.5%` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	234.507s` (terminated `exit_code: 0`)

## Deviations from spec

- No scope or contract deviation. The first touched-package race attempt failed with `FAIL` (`exit_code: 1`) because a retained Phase 22 test treated `human.outcomes.v1` as unavailable. Its original subject is rejection of an unavailable capability; the accepted task-17 matrix invalidated only that example row. Replacing it with the still-unavailable `telemetry.trajectory.v1` kept the assertion and its no-journal-change check. The byte-identical touched-package race command later passed with `ok  	github.com/mbell-nzl/fid/cmd/fid	231.658s` (`exit_code: 0`).

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `994a8734d15580a23a2a79d3c28f5f33d9dbf1e5`

## Commit subject

- `feat(workflow): qualify final Phase 23 capability matrix`
