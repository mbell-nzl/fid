# Phase Brief - Phase 17 (Approved; QA Planning)

> **Purpose.** Default execution brief for the current phase. Full task and
> validation detail lives in the Phase 17 archive.
>
> **Size target:** under 150 lines.

## Phase Goal

Close v0.1 early dogfood with a trustworthy maintenance audit, locked SQLite
VACUUM, explicit version guards for schema-only governance entities, and
deterministic proof of the full implemented surface. Success means structural
violations and operational warnings are actionable, future commands fail with
honest version guidance, and the v0.1 journey/rebuild/format/telemetry criteria
are evidenced rather than inferred.

## Status

**Approved; awaiting QA plan.** The Human approved Heavy Mode and DC-A/DC-B/DC-C
on 2026-07-11; DEC-049..052 are recorded and VAL-001..025 are Ready. No
Developer task may start until QA Lead creates the Heavy Mode QA plan with fresh
command and fixture feasibility validation.

## Operating Mode

**Heavy (approved, DEC-049).** Phase 17 adds public CLI/config contracts, mutates the
derived persistence file under lock, and closes a version. Required workflow:
validation contract, nine task specs, fresh QA-plan feasibility validation,
mid-phase reviews after Tasks 2, 5, and 7, final review, human dogfood, and
full-scope spec sync. This is v0.1 early-dogfood closure, not the charter's v0.3
public release; public-release threat model/scans/pre-release checklist and
held-out scenarios remain deferred under `08_QUALITY_PLAN.md` sections 5.8/7/9.

## Task Index

Full specs: `../phases/phase-17/tasks/`.

| Task | Title | Status | Dependencies |
|---|---|---|---|
| 1 | CF-019 markdown formatter split | Pending | QA plan Ready |
| 2 | Maintenance grace config contract | Pending | DEC-050 |
| 3 | Maintenance structural errors | Pending | DEC-049 |
| 4 | Maintenance warnings and info | Pending | Tasks 2-3 |
| 5 | `fid maintenance check` CLI + formats | Pending | Tasks 1, 3-4 |
| 6 | `fid maintenance vacuum` | Pending | Task 5 parent command |
| 7 | Governance stub-table guards | Pending | DEC-052 |
| 8 | v0.1 read-format conformance | Pending | Task 5 |
| 9 | Full v0.1 closure journey | Pending | Tasks 2-8 |

## Validation Contract

Full contract: `../phases/phase-17/qa/validation-contract.md`.

| VAL range | Scope | Owner | Status |
|---|---|---|---|
| VAL-001..002 | CF-019 + config compatibility | Review 1 | Ready |
| VAL-003..014 | maintenance errors/warnings/check CLI | Review 2 | Ready |
| VAL-015..016 | vacuum + governance guards | Review 3 | Ready |
| VAL-017..023 | read parity + closure journey | QA Lead | Ready |
| VAL-024 | v0.1 PRD traceability + full/static gates | QA Lead, phase-end | Ready |
| VAL-025 | real-work v0.1 dogfood | Human, phase-end | Ready |

## QA Artifacts

| Artifact | Path | Status |
|---|---|---|
| Validation contract | `../phases/phase-17/qa/validation-contract.md` | Ready |
| QA plan | `../phases/phase-17/qa/qa-plan.md` | Missing; next required artifact |
| Mid-phase reviews | `../phases/phase-17/qa/review-1.md` etc. | After Tasks 2, 5, 7 |
| Final review | `../phases/phase-17/qa/final-review.md` | Pending |
| Human checklist | `../phases/phase-17/qa/human-checklist.md` | Pending |

## Approved Decisions

- **DEC-050 (DC-A):** `[maintenance].started_without_claim_grace_seconds`, default 300,
  zero allowed, negative/malformed rejected; age from latest append-ordered
  `status_changed`, fallback `issues.created_at`. Stale verification means a
  superseded completion reference; telemetry gap requires completion linkage.
- **DEC-051 (DC-B):** maintenance status is `pass|warn|fail`; warning promotion changes
  exit/status only; deterministic finding order; all four read formats.
- **DEC-052 (DC-C):** guard `decision create/get/list`, `escalation
  create/resolve/get/list`, and `triage record/start/verdict` with
  `commands for <entity> ship in v0.2`; leave `requirement sync` shipped.

## Carry-Forward Disposition

- CF-019: schedule in Task 1 because Task 5 reopens markdown rendering.
- CF-008: re-defer to v0.2; unrelated product/data-model enhancement.
- CF-017: re-defer to v0.2c; Phase 17 does not modify `internal/dev`.
- CF-018: re-defer to next SQL package reopen; Task 8 exercises but does not
  modify `internal/sqlquery`. A required SQL-package fix stops for re-planning
  and schedules CF-018 first.

## Parity Matrix

| Axis | Maintenance | Other v0.1 reads | Evidence |
|---|---|---|---|
| Formats | json/toon/markdown/yaml | json/toon/markdown/yaml | VAL-012/017 |
| Ordering | check/entity/message | established command order | VAL-011/017 |
| Mutation | check none; vacuum DB-only | reads none | VAL-014/015/023 |
| Failure | structured result + non-zero | documented command error | VAL-012/013/016 |

## Interfaces / Subsystems In Scope

- `internal/maintenance/`, `cmd/fid/maintenance.go` - audit and VACUUM.
- `internal/config/` - additive grace-window config.
- `internal/format/` - CF-019 split and maintenance markdown.
- `cmd/fid/governance_stubs.go` - v0.2 entity guards.
- `cmd/fid/e2e_v01*_test.go` - conformance and closure evidence.

## Known Risks

- Corrupt-state fixtures must bypass triggers only in tests; production pragmas
  and trigger defenses remain enabled.
- Warning-only non-zero promotion must retain structured stdout for CI.
- Task 9 is audit-only: any product defect stops for triage instead of expanding
  the final task across subsystem boundaries.

## Open Questions

None. QA-plan feasibility findings follow the normal hard-stop/escalation path.

## Success / Exit Criteria

- [ ] Tasks 1-9 complete with exact required gates.
- [ ] VAL-001..024 have deterministic terminal outcomes; VAL-025 dogfood is
  completed and accepted or explicitly dispositioned.
- [ ] QA final review has no unresolved blocking finding.
- [ ] Every v0.1 PRD row maps to passing evidence; static binary is <30 MB.
- [ ] Tech Lead wrap records spec sync, triage audit, and v0.1 closure decision.

## Next Planning Checkpoint

After Phase 17 wrap, slice roadmap v0.2a workflow-engine migration into a new
approved phase; do not start v0.2 work inside this closure phase.
