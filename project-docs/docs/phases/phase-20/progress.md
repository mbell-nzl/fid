# Phase 20 Progress - Clean Standalone Issue Journey

## Status

**Tasks 1-17 complete; Review 6 is next.** Task 18 has not started. Branch:
`phase-20`. The ESC-20-015 recovery is closed in the run ledger; the workflow
control hardening is undergoing its final independent infrastructure review.

## Planning artifacts

| Artifact | Path | State |
|---|---|---|
| Slice contract and removal inventory | `contracts.md` | Accepted |
| Validation contract | `qa/validation-contract.md` | Accepted and active (VAL-001..VAL-046) |
| Task specs | `tasks/task-1.md` .. `tasks/task-18.md` | Accepted |
| Phase decisions | `decisions.md` | Active through DEC-146 |
| QA plan | `qa/qa-plan.md` | Accepted and in execution |
| Triage log | `triage-log.md` | Current through TRG-20-015 |
| Escalations | `tl-escalations.md` | ESC-20-015 resolved; no open escalation |
| Run ledger | `run-ledger.md` | `ADVANCEABLE`; Task 17 complete |

## Accepted planning-time spec edits

- `03_TRD_SCHEMA.md` - new §8 accepted slice DDL for migrations
  `0011_p20_remove_issue_fields` and `0012_p20_issue_revision`.
- `04_TRD_CONFIG.md` - required config version plus explicit temporary
  `pricing.toml` expiry and atomic-init contract alignment.
- `06_DESIGN_PLAN_CLI.md` - §2 Response Envelope rules, §5 Phase 20 reason
  codes, new §8 accepted slice leaf contracts.
- `GLOSSARY.md` - new "Machine Interface" section resolving the
  Envelope/Response Envelope collision.

## Task status

| Task | Title | Status |
|---|---|---|
| 1 | Response contract foundation | Complete |
| 2 | Project selection and config namespace | Complete |
| 3 | Atomic project initialization | Complete |
| 4 | Rebuild, authority, and lock failures | Complete |
| 5 | Remove orphaned Issue fields | Complete |
| 6 | Replayable Issue aggregate revision | Complete |
| 7 | Issue hierarchy invariant | Complete |
| 8 | Issue read contract | Complete |
| 9 | Issue mutation and concurrency contract | Complete |
| 10 | Acceptance-criterion and label contracts | Complete |
| 11 | Link contract and dependency invariants | Complete |
| 12 | Workflow command contract | Complete |
| 13 | Maintenance inspection contract | Complete |
| 14 | Legacy JSON renderer conversion | Complete |
| 15 | Remaining format caller/test convergence | Complete |
| 16 | Delete universal formatters and YAML | Complete |
| 17 | Delete promised-command stubs | Complete |
| 18 | Journey acceptance and documentation | Pending |

Recovery slices `task-3a`..`task-3l`, `task-9a`, `task-16a`, and `task-16b`
are complete. Reviews 1-5 and their linked re-reviews are closed. Task 17 is
complete at product commit `95ee13d` and report finalization `0c1094c`.

## Next action

Start a fresh Phase Runner session on `phase-20`, append a full schema-v2 normal
`RUN-RESUME`, and obtain Review 6 covering Tasks 16-17. If Review 6 passes,
dispatch Task 18. Do not reopen Task 17 or rerun its completed product work.
