# Phase 21 - Mid-Phase Review 7b (tasks 19-21)

Reviewer: qa-reviewer
Date: 2026-09-02
Tasks covered: phase-21/task-19, phase-21/task-20, phase-21/task-21
Review type: RE-REVIEW
Prior review: project-docs/docs/phases/phase-21/qa/review-7.md

## Summary

- Review outcome: PASS
- Must-fix findings: 0
- Improvement findings: 0
- Blocker candidates: 0
- Existing escalation references: ESC-21-014, closed by DEC-173; provenance only
- Covered validation assertions: VAL-002, VAL-069, VAL-070, VAL-071, VAL-074, VAL-075, VAL-076, VAL-077, VAL-078, VAL-079, VAL-080, VAL-081, VAL-082, VAL-083

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-002 | task-20 | ON TRACK | Unchanged support from review 7: task-20 report retains the migration 0014 evidence and required gate receipts; no task-20 report or implementation file changed in `5cc8fc0..HEAD`. |
| VAL-069 | task-19 | ON TRACK | Unchanged `TestCreateIsByteIdenticalWithCompleteDeterministicInventory` at `internal/export/export_test.go:20` remains in the focused package gate. |
| VAL-070 | task-19 | PASS | QA-68 is now satisfied by `TestValidateRejectsEveryArchiveInvalidCheckWithoutMutation` at `internal/export/export_test.go:97`: eleven isolated fresh workspace/archive rows cover `artifact-reference-invalid`, `capability-snapshot-mismatch`, `digest-mismatch`, `duplicate-entry`, `entry-extra`, `entry-missing`, `journal-invalid`, `manifest-invalid`, `path-unsafe`, `size-mismatch`, and `terminal-hash-mismatch`; each compares pre/post workspace tree digests and asserts its expected `InvalidArchiveError` check. The exact affected race test passed. |
| VAL-071 | task-19 | ON TRACK | Import remains absent: both independently run source-absence controls produced no match with expected `exit_code: 1`; the unchanged CLI test retains unknown-command and unavailable-capability assertions. |
| VAL-074 | task-20 | ON TRACK | Unchanged support from review 7; task-20 report/evidence is outside the repair diff. |
| VAL-075 | task-20 | ON TRACK | Unchanged support from review 7; task-20 report/evidence is outside the repair diff. |
| VAL-076 | task-21 | ON TRACK | Unchanged support from review 7; task-21 report/evidence is outside the repair diff. |
| VAL-077 | task-21 | ON TRACK | Unchanged support from review 7; task-21 report/evidence is outside the repair diff. |
| VAL-078 | task-21 | ON TRACK | Unchanged support from review 7; task-21 report/evidence is outside the repair diff. |
| VAL-079 | task-21 | ON TRACK | Unchanged support from review 7; task-21 report/evidence is outside the repair diff. |
| VAL-080 | task-21 | ON TRACK | Unchanged support from review 7; task-21 report/evidence is outside the repair diff. |
| VAL-081 | task-21 | ON TRACK | Unchanged support from review 7; task-21 report/evidence is outside the repair diff. |
| VAL-082 | task-21 | ON TRACK | Unchanged support from review 7; task-21 report/evidence is outside the repair diff. |
| VAL-083 | task-21 | ON TRACK | Unchanged support from review 7; task-21 report/evidence is outside the repair diff. |

QA-81 through QA-85 remain assigned to task 22 and final QA, outside this window.

## Gate And Evidence Checks

| Check | Command or evidence | Result | Notes |
|---|---|---|---|
| MF-19-01 corpus / QA-68 | `go test -race ./internal/export -run '^TestValidateRejectsEveryArchiveInvalidCheckWithoutMutation$' -count=1` | PASS | `12.106s` (terminated `exit_code: 0`). The source establishes a fresh workspace and archive per row, performs one selected corruption per pristine archive, asserts an unchanged before/after workspace digest, and asserts the expected check code. |
| Required source absence | `rg -n --glob '*.go' --glob '!**/*_test.go' 'Use:\\s*"import|newExportImport|export\\.import' cmd/fid internal` | PASS | `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`). |
| Task-19 gate 1 | `go test -race ./internal/export ./internal/evidence ./internal/journal ./internal/cli ./cmd/fid/...` | PASS | `(cached)` (terminated `exit_code: 0`). |
| Task-19 gate 2 | `rg -n --glob '*.go' --glob '!**/*_test.go' 'Use:\\s*"import|newExportImport|export\\.import' cmd/fid internal` | PASS | `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`). |
| Task-19 gate 3 | `gofmt -l cmd internal` | PASS | `exit 0 (no output)` (terminated `exit_code: 0`). |
| Task-19 gate 4 | `golangci-lint run ./...` | PASS | `0 issues.` (terminated `exit_code: 0`). |
| Task-19 gate 5 | `go test -race -coverprofile=coverage.out ./...` | PASS | `27.5% of statements` (terminated `exit_code: 0`). |
| Task-19 gate 6 | `go tool cover -func=coverage.out` | PASS | `61.5%` (terminated `exit_code: 0`). |
| Task-19 report control | `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-19-done.md project-docs/docs/phases/phase-21/tasks/task-19-export-create-validate.md` | PASS | `PASS report-lint completed with no failures` (terminated `exit_code: 0`). The report lists all six task-spec commands as PASS with observed termination, preserves the historical lint failure/pass sequence, and preserves the QA-rework red/green corpus evidence. |
| Commit and worktree provenance | `git rev-parse --verify`, `git show`, `git diff`, and `git status --short` | PASS | `3e2736549ef100e3d1cdbd3b11d1a51488e10ac6` resolves with subject `test(export): cover archive validation corpus` and `Refs: phase-21/task-19`; report commit `47109dcdf7ff72364467b8f613c2ff83a530900b` resolves with its matching task reference. `git diff --name-only 5cc8fc0..HEAD` contains only `internal/export/export_test.go` and `task-19-done.md`; `git diff --check` passed. The only pre-existing dirty files are parent-owned `run-ledger.md` and `triage-log.md`. |

## Task 19 - Deterministic Export Create And Validate

### Must-fix

- None. MF-19-01 is resolved.

### Improvement

- None.

## Task 20 - V1 Removal And Migration 0014

### Must-fix

- None. Its report and supporting evidence are unchanged from review 7.

### Improvement

- None.

## Task 21 - Final-Schema Maintenance, Init, Rebuild, And Export Parity

### Must-fix

- None. Its report and supporting evidence are unchanged from review 7.

### Improvement

- None.

## Prior Finding Resolution

| Prior finding | Owning task | Result | Fresh evidence |
|---|---|---|---|
| MF-19-01: VAL-070 lacked the QA-68 full archive-invalid/no-mutation corpus. | task-19 | RESOLVED | Commit `3e2736549ef100e3d1cdbd3b11d1a51488e10ac6` adds the closed eleven-code corpus without changing production semantics. Each subtest has a separate initialized workspace/archive, changes only the selected archive representation, checks its expected code, and compares the complete project tree digest before/after `Validate`. The exact corpus race test, the byte-identical task-19 six-command gate, the required source-absence control, and the two-argument report control all passed with observed terminations. The prior aggregate-sequence collision test remains at `internal/export/export_test.go:70`. |

## Cross-cutting observations

- The repair range does not modify `internal/export/export.go`, `cmd/fid/export.go`, `cmd/fid/export_test.go`, or `internal/cli/codes.go`; archive format, validation semantics, public `archive-invalid` handling, import absence, and unavailable capability behavior therefore remain the reviewed behavior.
- The historical failed lint attempt and passing rerun, plus the QA-rework red then green corpus attempts, remain recorded in `task-19-done.md`; this review does not relabel or discard them.
- Task-20/task-21 reports are absent from the repair range, so review-7 support remains intact rather than being silently revalidated from Developer attestation.

## Repair Plan

- None. No Must-fix finding remains.

## Blocker Candidates

- None.

## Existing Escalation References

- ESC-21-014, closed by DEC-173; provenance only and unrelated to this bounded task-19 test-evidence repair.

## Final Lines

Review written: project-docs/docs/phases/phase-21/qa/review-7b.md
Must-fix: 0  Improvement: 0  Blocker candidates: 0
VERDICT: PASS
