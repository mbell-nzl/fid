# Phase 11 Decisions

> Purpose: phase-local decision history, spec sync verdicts, workflow notes, and
> wrap-up decisions. DEC numbering is global and never resets.

## Decisions

### DEC-028: Claim writes must refresh stale projections without nested locks

- **Date:** 2026-06-21
- **Context:** ESC-008 / TRG-001 found that `internal/claims.Acquire`,
  `Renew`, and `Release` call `journal.EnsureProjectionFresh` from inside their
  claim write-lock callback. When the journal is ahead of SQLite, that helper
  can catch up or rebuild by trying to acquire `.fid/write.lock` again, which
  produced nested-lock failures in Task 4 fixture work.
- **Decision:** The current lock ordering is not acceptable for Phase 11 final
  QA. Add do-now Task 5 before final review. Claim write paths must validate
  against a fresh projection while avoiding nested acquisition of the same
  checkout-scoped write lock. A lock-aware journal helper or equivalent narrow
  refactor is acceptable; public claim CLI flags/output and claim semantics must
  not change.
- **Consequences:** Validation contract gains VAL-015 and QA plan gains
  corresponding evidence. QA review may resume after Task 5 is complete and
  reviewed. No schema migration or spec change is planned because this enforces
  the existing TRD write-path/concurrency contract.
- **Refs:** ESC-008; phase-11/task-5; PRD-NFR-008; `02_TRD.md` §4.3 and §4.9;
  ADR-0001; VAL-015.

### DEC-027: `fid claim list --by` is a read filter, not an actor flag

- **Date:** 2026-06-21
- **Context:** `06_DESIGN_PLAN_CLI.md` defines `fid claim list [--stale] [--by
  <role>] [--format ...]` and describes `--by` as a filter by claiming role.
  Phase 10 permission enforcement applies to write commands via acting roles;
  `claim list` is read-only.
- **Decision:** Treat `fid claim list --by` as a filter over stored
  `claimed_by`, not as an acting role. The command does not call
  `policy.Authorize`, does not require the role to be known, and appends no
  journal/activity rows. Unknown filter strings return zero rows unless an
  existing read-filter pattern forces a contrary decision, in which case
  Developer must escalate before implementation.
- **Consequences:** This keeps read surfaces inspectable by agents before they
  act and avoids overloading `--by` with two meanings on the same command.
- **Refs:** PRD-CLAIM-007; `06_DESIGN_PLAN_CLI.md` `fid claim list`;
  `decisions_principles.md` read-ordering rule; phase-11/task-3.

### DEC-026: Claim schema and matrix handling at kickoff

- **Date:** 2026-06-21
- **Context:** Phase 10 intentionally kept the in-code matrix limited to
  implemented write actions. Phase 11 now implements claim writes. The claims
  table/triggers already exist from v0.1a, and `03_TRD_SCHEMA.md` section 13
  already contains the `issue.claim.*` default-profile rows.
- **Decision:** Phase 11 adds no schema migration by plan. Task 1 appends the
  four claim actions to the in-code matrix by copying the existing section 13
  role booleans. `03_TRD_SCHEMA.md` is edited only if implementation proves the
  current spec wrong; otherwise it receives a verified-no-change verdict at wrap.
- **Consequences:** Matrix growth remains per-implemented-action, while the
  schema/spec remains the durable source for default cell values.
- **Refs:** TRD-ARCH-010; TRD-ARCH-011; `03_TRD_SCHEMA.md` sections 11.1 and 13;
  phase-11/task-1.

## Carry-Forward Review At Kickoff

| ID | Disposition | Rationale |
|---|---|---|
| CF-008 | Re-defer to v0.2 | Resolution reason is a product/data-model enhancement unrelated to claims. |
| CF-011 | Re-defer broad reconciliation to v0.1e; schedule claim renewal projection now | PRD-CLAIM-008a explicitly requires renewal to project to `activity_log`, so Phase 11 handles that per-command case without widening the broader timeline policy. |
| CF-013 | Re-defer to v0.1e | `complete.go` orchestration cleanup belongs with `fid dev finish`/`block`, not primitive claims. |

## Spec Sync Verdict

Record a verdict for every spec at phase wrap-up. `verified-no-change` is a
positive assertion after review, not the default.

| Spec | Verdict | Evidence / Reference |
|---|---|---|
| `00_PROJECT_CHARTER.md` | pending wrap review | Phase in progress. |
| `01_PRD.md` | pending wrap review | Phase in progress. |
| `02_TRD.md` | pending wrap review | Phase in progress. |
| `03_TRD_SCHEMA.md` | pending wrap review | DEC-026; verify claim rows/schema. |
| `04_TRD_CONFIG.md` | pending wrap review | Phase in progress. |
| `05_DESIGN_PLAN.md` | pending wrap review | Phase in progress. |
| `06_DESIGN_PLAN_CLI.md` | pending wrap review | Claim command/list output verification. |
| `07_DESIGN_PLAN_FLOWS.md` | pending wrap review | Phase in progress. |
| `08_QUALITY_PLAN.md` | pending wrap review | Phase in progress. |
| `09_AI_WORKFLOW.md` | pending wrap review | Phase in progress. |
| `10_IMPLEMENTATION_ROADMAP.md` | pending wrap review | Phase in progress. |

## Triage-Log Audit

- [ ] Walked every TRG-NNN in `triage-log.md`.
- [ ] Confirmed each RESUME verdict was correct in hindsight.
- [ ] Flagged any TRG that should have escalated.
- [ ] Confirmed autonomous TL-Lite resume remains in effect (DEC-002; no warm-up tally).

## Workflow Notes

Answer at wrap.

- **What worked well?** Pending.
- **What was inefficient?** Pending.
- **Were tasks the right size?** Pending.
- **Did the QA review cadence feel right?** Pending.
- **Were there token-efficiency opportunities missed?** Pending.

### Deep Retrospective Recommendation

Pending wrap assessment.
