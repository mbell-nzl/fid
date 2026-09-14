# Phase 23 Decisions

> Purpose: phase-local decision history, spec sync verdicts, workflow notes,
> and wrap-up decisions. DEC numbering is global and never resets.
>
> Status: planning draft. Nothing below is accepted yet.

## Pending Acceptance

### DEC-222 (PROPOSED): Accept the Phase 23 slice contract and task pack

- **Date:** proposed 2026-09-14
- **Context:** The approved roadmap calls Phase 23 "Review and Human judgment".
  Its exact contracts and task pack were explicitly left to this planning step
  by the Phase 22 wrap. Roadmap section 5 requires the JSON, event, and schema
  contracts to be accepted before any product code in the slice changes.
- **Decision requested:** Accept, together and as one unit,
  `contracts.md`, `public-interface.md`, `schema.md`,
  `execution-scope.md`, `v1-convergence-inventory.md`, `contract-coverage.md`,
  `qa/validation-contract.md`, and the ordered 18-task pack; declare the phase
  Heavy; and accept the four canonical spec amendments listed in
  `contracts.md` section 15.
- **Consequences if accepted:** Implementation may begin at task 1 after the
  Runner's QA readiness check. `reviews.coverage.v1` and `human.outcomes.v1`
  become advertisable at task 17. The named V1 judgment surfaces are removed at
  task 15 under replacement-bound removal. `work.claims.v1` stays unavailable
  until Phase 24 and `telemetry.trajectory.v1` until Phase 25.
- **Consequences if not accepted:** No Phase 23 product code, migration, or
  test may start; the planning artifacts remain drafts.
- **Refs:** `10_IMPLEMENTATION_ROADMAP.md` sections 4, 5, 5.1, 7; PRD sections
  7 and 8; TRD sections 4, 5.4, and `03_TRD_SCHEMA.md` sections 3.7 and 7;
  Phase 22 DEC-187 kit dispositions C-03, C-04, C-10; DEC-220.

### Proposed carry-forward dispositions

| ID | Proposed Phase 23 disposition |
|---|---|
| CF-008 | Re-defer to Phase 24 with the Envelope disposition work for removed unfinished Issues. |
| CF-018 | Re-defer; no planned production edit to `internal/sqlquery/sqlquery.go` in this phase. |
| CF-023 | Re-defer to dedicated workflow maintenance. |
| CF-026a | Re-defer; unscheduled product candidate. |
| CF-026b | Re-defer; separate Human product decision. |
| CF-029 | Re-defer to dedicated workflow maintenance. |
| CF-030 | Re-defer to dedicated workflow maintenance. |
| CF-033 | Re-defer to dedicated workflow maintenance. |
| CF-036 | Re-defer to the dedicated workflow retrospective. |
| CF-037 | Re-defer to Phase 25; Phase 23 must not claim any coordinator obligation. |
| CF-038 | Re-defer; the trigger fires only on a causal edit to `cmd/fid/live_output_test.go`. |
| CF-039 | Re-defer the executable cross-check; Phase 23 adopts the manual mitigation in `contract-coverage.md` and records whether it worked in the wrap-up workflow notes. |
| CF-040 | Partially scheduled: the guidance is applied to the Human views this phase introduces (task 14, VAL-073). The polish and verbosity review of pre-existing views stays with Phase 25. |
| CF-041 | New: hosted CI runs only the ordinary suite. The tagged integration suite and the rebuild benchmark are local-only gates, so a green PR is not phase-level evidence. Owner: workflow and CI maintainer, in dedicated maintenance. |

## Decisions

None recorded yet.

## Spec Sync Verdict

Recorded at wrap-up. `verified-no-change` is a positive assertion after review,
not the default.

| Spec | Verdict | Evidence / Reference |
|---|---|---|
| `00_PROJECT_CHARTER.md` | pending | |
| `01_PRD.md` | pending | |
| `02_TRD.md` | pending | expected `touched` at task 17 |
| `03_TRD_SCHEMA.md` | pending | expected `touched` at task 1 (section 8.4) |
| `04_TRD_CONFIG.md` | pending | |
| `05_DESIGN_PLAN.md` | pending | |
| `06_DESIGN_PLAN_CLI.md` | pending | expected `touched` at tasks 8, 9, 12, 13 |
| `07_DESIGN_PLAN_FLOWS.md` | pending | |
| `08_QUALITY_PLAN.md` | pending | |
| `09_AI_WORKFLOW.md` | pending | |
| `10_IMPLEMENTATION_ROADMAP.md` | pending | |
| `GLOSSARY.md` | pending | |

## Triage-Log Audit

- [ ] Walked every TRG-23-NNN in `triage-log.md`.
- [ ] Confirmed each RESUME verdict was correct in hindsight.
- [ ] Flagged any TRG that should have escalated.
- [ ] Confirmed autonomous TL-Lite resume remains in effect (DEC-002).

## Workflow Notes

Recorded at wrap-up.
