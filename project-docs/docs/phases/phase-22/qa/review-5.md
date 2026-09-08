# Phase 22 - Mid-Phase Review 5 (tasks 13-15)

Reviewer: qa-reviewer
Date: 2026-09-08
Tasks covered: phase-22/task-13, phase-22/task-14, phase-22/task-15
Review type: INITIAL
Prior review: none

## Summary

- Review outcome: PASS
- Must-fix findings: 0
- Improvement findings: 2
- Blocker candidates: 0
- Existing escalation references: none
- Covered validation assertions: VAL-054, VAL-055, VAL-056, VAL-057, VAL-058, VAL-059, VAL-060, VAL-061, VAL-062, VAL-063, VAL-064, VAL-065

## Scope Reviewed

- Task specifications and completion reports: `tasks/task-13.md` through `tasks/task-15.md`, and their corresponding `-done.md` reports.
- QA artifacts: `qa/qa-plan.md` and `qa/validation-contract.md`.
- Product commits: `acb569bb9447a3259991520669167fdc819c1706` (task 13), `8c322b53c749a2a7f0dd0be3e03dbb5034d70e40` (task 14), and `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee` (task 15). All three resolve as commits and carry the recorded task references.
- Completion-report commits: `59d86a315b4c042c79583bd5ca3cac885ee385b0` and `b9ad453efb52c27f2705ea136acb997de42d9a25` resolve. `report-lint.sh` independently passed for each report/spec pair.

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-054 | task-13 | ON TRACK | `archive_v2_test.go` creates the v2 config and ordinal journal control; validation admits only the accepted v2 profile and validates the closed manifest/inventory. |
| VAL-055 | task-13 | ON TRACK | `archive_corruption_test.go` rejects integrity-valid unknown Event/table/column/type, reference, aggregate, and immutable-Effect forgeries; `archive_semantics.go` replays against trusted migration profiles. |
| VAL-056 | task-13 | ON TRACK | `archive_scan.go` measures streamed member and cumulative bytes, paths, regular-file type, count, and line limits; `archive_limits_test.go` provides below/equal/above boundary cases. |
| VAL-057 | task-13 | ON TRACK | `TestValidateHasNoSourceOrDestinationEffects` covers valid and invalid validation with source/destination tree digests. |
| VAL-058 | task-14 | ON TRACK | staged import records ordered immutable journal delivery metadata; `TestImportBindingRejectsAlteredHistoryIdentity` and the focused reproduced suite cover changed imported bytes/identity rejection. |
| VAL-059 | task-14 | ON TRACK | target-inode lock and `RENAME_NOREPLACE` publication are in `import_publication.go`; real abrupt-exit and repeated-contender coverage is in `import_boundary_test.go`. |
| VAL-060 | task-14 | ON TRACK | installed binding validates archive/key/result before returning the replayed result; changed digest/key and altered installed history are negative controls. |
| VAL-061 | task-14 | ON TRACK | target safety, empty-target enforcement, staged same-filesystem publication, and source-file preservation are exercised by the import boundary tests. |
| VAL-062 | task-14 | ON TRACK | `TestImportRequiresFreshDestinationObservationBeforeRecovery` rejects imported old-observation recovery, close, abandon, and historical execution before allowing a fresh destination-bound successor and Receipt. |
| VAL-063 | task-15 | ON TRACK | `p22PortableThreeWay` compares full named JSON reads between live, rebuilt, imported, and imported-rebuilt states; no entity/Event/status/error/null field is normalized away. |
| VAL-064 | task-15 | ON TRACK | the same helper compares all archive-entry bytes for source, rebuilt, and imported re-exports. |
| VAL-065 | task-15 | ON TRACK | built-CLI integration journeys execute export, validate, import, rebuild, resume, fresh destination recovery, and a successor Gate Receipt. |

## Gate And Evidence Checks

| Check | Command or evidence | Result | Notes |
|---|---|---|---|
| Task 13 exact commands | five required commands in `task-13.md` | PASS (report evidence) | All appear as final `PASS` receipts with output token and observed termination in `task-13-done.md`; historical failed attempts remain preserved. |
| Task 14 exact commands | five required commands in `task-14.md` | PASS (report evidence) | All appear as final `PASS` receipts with output token and observed termination in `task-14-done.md`; no required command was narrowed or omitted. |
| Task 15 exact commands | five required commands in `task-15.md` | PASS (report evidence) | All appear as final `PASS` receipts with output token and observed termination in `task-15-done.md`; the earlier full-race fixture failure and report-heading repair remain disclosed. |
| QA-plan task gate | `gofmt`, lint, race/coverage, coverage report, tagged integration | PASS (report reconciliation) | The QA-plan required base is byte-for-byte the same five-command sequence required by each task specification. |
| Reproduced required command | `go test -tags integration -race ./cmd/fid/integration/...` | PASS | `ok  github.com/mbell-nzl/fid/cmd/fid/integration  72.844s` (terminated `exit_code: 0`). Independent Review5 execution from repository root. |
| Reproduced focused high-risk evidence | `go test -race ./internal/export -run '^(TestImportProcessExitAtPublicationBoundaries|TestImportLockExcludesRepeatedContenders|TestImportRequiresFreshDestinationObservationBeforeRecovery|TestArchiveRejectsIntegrityValidSemanticCorruptions|TestValidateRetainsPinnedV1Profile)$' -count=1` | PASS | `ok  github.com/mbell-nzl/fid/internal/export  46.447s` (terminated `exit_code: 0`). Covers process-exit publication boundaries, exclusion, imported-history fencing, semantic corruption, and both v1 qualifying/non-qualifying outcomes. |
| Mechanical evidence | `project-docs/tools/report-lint.sh <report> <task-spec>` for tasks 13, 14, and 15 | PASS | Each independent invocation completed with `PASS report-lint completed with no failures` (terminated `exit_code: 0`). |

## Task 13 - Recovery archive format and validation

### Must-fix

- None.

### Improvement

- [internal/export/archive_profiles.go:3] `capabilitySnapshot` intentionally delegates to the retained profile until the final capability matrix is qualified in task 17. Keep the schema/profile transition visible in task 17's capability-matrix tests so the historical-profile naming cannot conceal a premature capability advertisement. Fix direction: task 17 should assert the final manifest/discovery/fake matrix as its owned VAL-070 work.

Observation: v1 is a separate read-only schema-14 path at `archive_validation.go:125-136`; the focused reproduced test exercised both required states: v1 validation succeeds and v1 import returns `schema-unsupported` without changing the target. V2 does not fall back to this profile. V2 semantic checks replay only against trusted migrations and reject integrity-valid semantic corruption, while the streaming scan applies actual cumulative/member/line bounds before archive use.

## Task 14 - Atomic import into an empty target

### Must-fix

- None.

### Improvement

- None.

Observation: `import.go:72-155` snapshots and validates the archive before staging, then writes the binding and synchronizes staged state before publication. `import_publication.go:27-47` locks the target directory inode and publishes with no-replace rename, avoiding the prior removable lock-path race. The independently reproduced focused suite covers all five real crash boundaries, repeated contenders, source preservation, immutable installed history, and the old-observation/no-authority versus fresh-local-successor states.

## Task 15 - Import CLI and three-way recovery parity

### Must-fix

- None.

### Improvement

- [cmd/fid/integration/p22_portable_conformance_test.go:20] The eight portable fixture expectations intentionally assert the scenario-relevant semantic summaries, while full named-view identity/history parity is asserted separately by fid-to-fid three-way comparisons. Carry a field-inventory assertion into task 18's VAL-073 conformance audit so every fake/fid shared expected-output dimension is explicitly enumerated as the schema evolves. Fix direction: preserve the current independent in-memory fake and extend the task-18 inventory rather than adding an alternate fake expectation surface now.

Observation: `p22PortableThreeWay` compares complete named JSON results and every re-exported entry byte for live, rebuilt, imported, and imported-rebuilt state. `p22PortableFakeImport` consumes only archive-entry bytes into an in-memory Effect model; it neither opens a fid workspace nor invokes production appliers. The five valid-hash negative mutations are rejected by both fake and fid. The independent tagged suite passed all integration cases, including the eight active/interrupted/abandoned/Attention/Budget/activation/Artifact portable states and the CLI import-recovery journey.

## Cross-cutting observations

- The phase parity matrix is represented in executable evidence: archive identity/history is preserved, imported observations remain historical, source/rebuilt/imported views and canonical re-exports agree, and actual destination observation/recovery is separately represented as a new local fact. No Human judgment has been treated as deterministic evidence; the required Human outcomes remain pending.
- `git diff --check` passed for each product commit. The reviewed code keeps archive validation in `internal/export`, import publication in its adapter, and recovery authority decisions in `internal/workflow`; no new runtime dependency was added.

## Repair Plan

- Not required: no Must-fix finding remains.

## Blocker Candidates

- None.

## Existing Escalation References

- None.

## Final Lines

Review written: project-docs/docs/phases/phase-22/qa/review-5.md
Must-fix: 0  Improvement: 2  Blocker candidates: 0
VERDICT: PASS
