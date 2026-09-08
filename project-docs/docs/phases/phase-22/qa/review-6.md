# Phase 22 - Mid-Phase Review 6 (tasks 16-18)

Reviewer: qa-reviewer
Date: 2026-09-08
Tasks covered: phase-22/task-16, phase-22/task-17, phase-22/task-18
Review type: INITIAL
Prior review: none

## Summary

- Review outcome: REPAIR
- Must-fix findings: 1
- Improvement findings: 0
- Blocker candidates: 0
- Existing escalation references: none
- Covered validation assertions: VAL-053, VAL-066, VAL-067, VAL-068, VAL-069, VAL-070, VAL-071, VAL-072, VAL-073, VAL-074, VAL-075

## Scope and independent evidence

- Verified product commits with `git rev-parse --verify`: task 16 `644f30e332a9e5099f0cb31c77fc6b9699fb5db8`, task 17 `ccc15f664c74494cf16219ca350fbec96c51808d`, and task 18 `24e3cbdb22dd9e43b55636da313dcef1ef633b6f`; each resolves and has the recorded subject and `Refs: phase-22/task-N` footer. `git diff --check` is clean for each product commit.
- Independently ran `project-docs/tools/report-lint.sh` for each report/spec pair. All three completed `PASS report-lint completed with no failures` (terminated `exit_code: 0`). The completion reports reconcile every required occurrence: task 16 12/12, task 17 17/17 including three distinct reference samples and two distinct full benchmark commands, and task 18 7/7. Each final receipt is `PASS` with a command token and observed termination.
- Reproduced unchanged required command: `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok github.com/mbell-nzl/fid/cmd/fid/integration 101.237s` (terminated `exit_code: 0`).
- Reproduced risk-directed conformance: `go test -tags integration ./cmd/fid/integration/... -run '^TestP22(BudgetCrossingFrameOrder|CompleteOneIssueRecoveryJourney|PortableRecoveryConformance|HumanFixtures|AmendedFailureWireCorpus|ConformanceCorpusHasEveryAcceptedOperationAndFailure|FinalActivated)' -count=1` - PASS - `ok github.com/mbell-nzl/fid/cmd/fid/integration 29.209s` (terminated `exit_code: 0`).
- Reproduced risk-directed maintenance/capability/schema checks: `go test ./internal/maintenance -run '^(TestMaintenanceReportsBrokenRecoveryLinkWithoutMutation|TestMaintenanceRecoveryIntegrityMatrixDoesNotMutate)$' -count=1 && go test -tags integration ./cmd/fid/integration -run '^(TestP22QualifiedCapabilityMatrixMatchesFakeDiscoveryAndExport|TestRecoverySchemaFreshBinaryInitAndRebuild)$' -count=1` - PASS - `ok github.com/mbell-nzl/fid/internal/maintenance 0.193s`; `ok github.com/mbell-nzl/fid/cmd/fid/integration 1.173s` (terminated `exit_code: 0`).

## Task 16 - Retire escalation state and duplicate recovery truth

### Must-fix

- None.

### Improvement

- None.

Observation: `migrations/0016_p22_remove_escalations.up.sql` removes only the accepted escalation state and restores the activity trigger without the retired type. `internal/db/escalation_removal_test.go:9` checks latest-schema absence together with retained tables and seed removal, while `internal/journal/escalation_removal_test.go:9` supplies both retired-Effect rejection and retained-Effect controls with rollback proof. This supports VAL-066 and the Event-derived boundary in VAL-067. The report's 12 exact required command occurrences and retained-row provenance are mechanically complete for VAL-068.

## Task 17 - Maintenance, capability qualification, and rebuild budget

### Must-fix

- None.

### Improvement

- None.

Observation: `internal/maintenance/check_test.go:188` exercises clean controls, a valid later-superseded complete-forward activation, six corruption classes, and before/after journal and SQLite mutation checks. `cmd/fid/integration/p22_capability_matrix_test.go:12` compares the exact fake, live discovery, export, and resume capability views, including the unavailable Claims row. The task report preserves three separately recorded sub-30-second schema-16 reference samples and both full benchmark occurrences; this review did not rerun benchmarks concurrently with the required integration gate.

## Task 18 - Final recovery conformance and Human validation fixture

### Must-fix

- [cmd/fid/integration/p22_conformance_test.go:110] The `closed_values` inventory proves only that its 12 rows equal a parsed documentation snapshot. In contrast to operations, failures, and frames at `:119-124`, it contains no test references and never enters `p22CheckBehavioralReference`; the checker therefore does not establish qualifying and non-qualifying behavioral evidence for every quoted closed-value row and member. The successful corpus test, its 32 operation references, 22 error references, 10 frame references, and independent 23-action acceptance/rejection test do not fill that gap. Owning task: task-18. Repair boundary: extend the Task 18 wire inventory/conformance fixtures and the smallest matching tests so each of the 12 closed-value rows maps to explicit behavioral assertions of every accepted member and rejection of an unaccepted member, without changing any accepted value, wire contract, or gate. Required evidence: focused conformance output covering the strengthened inventory and named row assertions, then all seven byte-identical Task 18 commands with observed termination; preserve the existing failed/passing attempts. Re-review scope: inspect each closed-value row's positive and negative assertion, re-check the independent 23-action matrix and amended failure/frame tests, and rerun the affected focused corpus plus the tagged integration gate. Affected assertion: VAL-075.

### Improvement

- None.

Observation: the committed DEC-205 repair reads the canonical next aggregate sequence before writing the embedded Stop resolution (`internal/workflow/stop.go:265`); DEC-206 derives controls only from the current execution's committed Attempt/Receipt Events (`internal/evidence/budget_frames.go:15`), and DEC-207 is exercised by actual process failure responses in `cmd/fid/integration/p22_wire_failure_test.go:10`. The reproduced journey/frame/failure suite passed. The retained Human fixture manifest is present at `.cache/phase-22-human/checkpoints-2335511884/manifest.json`, identifies source commit `24e3cbd`, contains all four checkpoints and the retained binary, and its recorded SHA-256 matches the binary. VAL-076 through VAL-078 remain blank Human judgments. VAL-080 remains scheduled for fresh final QA after this review.

## Repair Plan

| Finding | Owning task | Affected assertion | Exact repair boundary | Required evidence | Fresh re-review scope |
|---|---|---|---|---|---|
| R6-01 | task-18 | VAL-075 | `cmd/fid/integration/p22_conformance_test.go`, `cmd/fid/integration/testdata/p22_conformance/wire-inventory.json`, and only the matching in-package or integration behavioral tests needed to cover every closed-value row/member and an invalid-member control | Focused strengthened corpus plus all seven exact Task 18 commands, each with token and observed termination | Audit all 12 row mappings and their accepted/rejected behavior; rerun focused corpus and `go test -tags integration -race ./cmd/fid/integration/...` |

## Prior Finding Resolution

- Not applicable: this is the initial Review 6 artifact.

## Cross-cutting observations

- The fake/fid portable conformance test names and asserts all ten semantic-summary fields for each of eight states; the fake imports archive-entry bytes into an in-memory model and does not open a fid workspace. Full named fid-view and archive-byte parity remain separately asserted. This is sufficient evidence for the reviewed portion of VAL-073, but it cannot substitute for the uncovered closed-value conformance required by VAL-075.
- No dependency, schema, public-interface, quality-gate, or accepted threshold change was found. The one finding is decision-free and wholly inside Task 18's accepted conformance-test scope.

## Blocker Candidates

- None.

## Existing Escalation References

- None.

## Final Lines

Review written: project-docs/docs/phases/phase-22/qa/review-6.md
Must-fix: 1  Improvement: 0  Blocker candidates: 0
VERDICT: REPAIR
