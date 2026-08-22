# Phase 20 - Mid-Phase Review 5b (Task 13, Review 5 MF-01 only)

Reviewer: qa-reviewer
Date: 2026-08-16
Tasks covered: phase-20/task-13
Review type: FRESH INDEPENDENT RE-REVIEW
Prior review: project-docs/docs/phases/phase-20/qa/review-5.md

## Scope Reviewed

- This fresh independent re-review is limited to Review 5 MF-01 and its
  bounded Task 13 repair. It does not reopen Task 13 implementation outside
  that repair.
- Review 5 MF-02 remains pending for Task 15. This review does not disposition,
  weaken, or otherwise re-review MF-02.
- Reconciled the Task 13 repair commits:
  - `78f2ed7ebaa9078579e42f4e71493f4d0cad114a` —
    `test(maintenance): prove exact grace boundary`
    (`Refs: phase-20/task-13`)
  - `c529a786e257ec104c19d1a5a0a05e2f6e255300` —
    `docs(task): finalize completion report for phase-20 task-13`
    (`Refs: phase-20/task-13`)
- Both hashes resolve with `git rev-parse --verify <hash>^{commit}`. The repair
  diff is limited to the maintenance fixture and completion report; the final
  report commit is clean. The only dirty files remain the pre-existing Phase
  Runner/user-authored documentation files, none of which this review modifies.

## Summary

- Review outcome: PASS
- Must-fix findings: 0
- Improvement findings: 0
- Blocker candidates: 0
- Existing escalation references: none
- Covered validation assertions: VAL-027

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Fresh evidence |
|---|---|---|---|
| VAL-027 | task-13 | ON TRACK | The fixed command-level fixture overrides `abandoned_attempt_grace_seconds` to `0`, seeds `FID-299` at `fixtureNow-299s` and `FID-300` at `fixtureNow-300s`, then decodes the canonical success result to require 299 absent and 300 present (`cmd/fid/maintenance_test.go:113-139`). It decodes the same fixture's promoted failure and compares `errors`, `warnings`, and `info` exactly (`:141-154`). The production command passes the fixed `legacyStartedWithoutClaimGraceSeconds = 300`, rather than Attempt config, to maintenance checking (`cmd/fid/maintenance.go:18,102-105`), matching the accepted result/channel contract (`contracts.md:479-495`). |

## Gate And Evidence Checks

| Check | Command or evidence | Result | Notes |
|---|---|---|---|
| Task-spec/QA-plan required-command coverage | Task 13 exact six-command gate | PASS | `task-13-done.md:69-74` records each byte-identical required command from `task-13.md` as PASS with a command-produced token and observed `exit_code: 0`; QA-plan Task 13 adds no different required gate. |
| Report outcome, prior evidence, and deviation record | `task-13-done.md:4-86` | PASS | The report is COMPLETE after one bounded recovery, retains the original red/green repair evidence, records the MF-01 focused command and real-entrypoint attempts, preserves observed termination for every gate, and explains that the prior 298/300 claim was corrected without changing the accepted check, registry, request/result contract, or gate definition. |
| Report commit identity | `78f2ed7…`, `c529a78…`, and `task-13-done.md:83-88` | PASS | The report records repair commit `78f2ed7…` and its exact subject; independent verification confirms that repair commit and the committed final-report follow-up `c529a78…` with their required Task 13 footer. |
| Mechanical report verification | `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-20/tasks/task-13-done.md project-docs/docs/phases/phase-20/tasks/task-13.md` | PASS | `PASS report-lint completed with no failures` (terminated `exit_code: 0`). It verifies the six exact required commands, tokens, observed termination, repair hash/subject, and committed-clean report. |
| Reproduced focused maintenance matrix | `go test ./cmd/fid -run '^(TestMaintenanceCheckJSONMatchesContract|TestMaintenanceCheckWarningPromotionPreservesFindingsAndInfo|TestMaintenanceCheckUsesFixed300SecondGraceAndRejectsFormat|TestMaintenanceCheckHumanSuccessAndVacuumJSONRejection|TestMaintenanceCheckUsesExplicitProjectRoot|TestMaintenanceCheckRejectsFindingLimitWithoutTruncation)$' -count=1` | PASS | From `/home/mark/projects/fid`: `ok  github.com/mbell-nzl/fid/cmd/fid  0.290s` (terminated `exit_code: 0`). The matrix includes contract, same-fixture warning-promotion, repaired 299/300, legacy/root, and no-truncation limit cases. |
| Required independent review-window reproduction | `go test -tags integration -race ./cmd/fid/integration/...` | PASS | From `/home/mark/projects/fid`: `ok  github.com/mbell-nzl/fid/cmd/fid/integration  4.725s` (terminated `exit_code: 0`). |

## Task 13 - Maintenance inspection contract

### Must-fix

- None.

### Improvement

- None.

## Prior Finding Resolution

| Prior finding | Owning task | Result | Fresh evidence |
|---|---|---|---|
| MF-01 — the reported command fixture was labelled 299 seconds but seeded 298 seconds, so it did not prove VAL-027's exact command-level boundary. | task-13 | RESOLVED | Repair commit `78f2ed7…` changes the fixture to exact nanosecond-preserving 299/300 seeds (`cmd/fid/maintenance_test.go:128-130`) after the Attempt-config override (`:113-121`). The decoded success result rejects a `FID-299` warning and requires `FID-300` (`:132-139`); the decoded promoted failure is `projection-divergence` and has equal `errors`, `warnings`, and `info` for that same fixture (`:141-154`). The independently rerun focused matrix and required integration command both passed. |

## Cross-cutting observations

- Review 5 MF-02 is intentionally untouched and remains the pending Task 15
  command-tree guard repair.
- No evidence warrants reopening resolved ESC-20-011: MF-01 was repaired
  entirely inside Task 13 without a product, architecture, data-model,
  public-interface, dependency, security, validation-contract, or quality-gate
  change.

## Blocker Candidates

- None.

## Existing Escalation References

- None. ESC-20-011 remains resolved and is not reopened.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-20/qa/review-5b.md
Must-fix: 0  Improvement: 0  Blocker candidates: 0
VERDICT: CONTINUE
```
