# Phase 20 Progress - Clean Standalone Issue Journey

## Status

**Implementation pack Human-approved; QA planning is next.** No implementation
has started. Branch: `main`.

## Planning artifacts

| Artifact | Path | State |
|---|---|---|
| Slice contract and removal inventory | `contracts.md` | Accepted |
| Validation contract | `qa/validation-contract.md` | Human-approved; Draft pending QA-plan acceptance (VAL-001..VAL-046) |
| Task specs | `tasks/task-1.md` .. `tasks/task-18.md` | Accepted |
| Phase decisions | `decisions.md` | DEC-107..DEC-129 |
| QA plan | `qa/qa-plan.md` | Missing - QA Lead |
| Triage log | `triage-log.md` | Seeded |
| Escalations | `tl-escalations.md` | Seeded |
| Run ledger | `run-ledger.md` | Seeded |

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
| 1 | Response contract foundation | Pending |
| 2 | Project selection and config namespace | Pending |
| 3 | Atomic project initialization | Pending |
| 4 | Rebuild, authority, and lock failures | Pending |
| 5 | Remove orphaned Issue fields | Pending |
| 6 | Replayable Issue aggregate revision | Pending |
| 7 | Issue hierarchy invariant | Pending |
| 8 | Issue read contract | Pending |
| 9 | Issue mutation and concurrency contract | Pending |
| 10 | Acceptance-criterion and label contracts | Pending |
| 11 | Link contract and dependency invariants | Pending |
| 12 | Workflow command contract | Pending |
| 13 | Maintenance inspection contract | Pending |
| 14 | Legacy JSON renderer conversion | Pending |
| 15 | Remaining format caller/test convergence | Pending |
| 16 | Delete universal formatters and YAML | Pending |
| 17 | Delete promised-command stubs | Pending |
| 18 | Journey acceptance and documentation | Pending |

## Next action

QA Lead writes and validates `qa/qa-plan.md`; Phase Runner then starts Task 1.
