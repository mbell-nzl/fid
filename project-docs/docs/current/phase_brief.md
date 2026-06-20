# Phase Brief - Phase 9 (Concerns + v0.1c closure)

> **Purpose.** Default execution brief for the current phase. Loaded by every
> agent session. Keep under 150 lines.
>
> **State:** Phase 8 (Task completions) is **CLOSED** — VAL-001..009 all PASS,
> human dogfood 5/5, zero BLOCKs/escalations. Phase 9 (Concerns + derived flag +
> v0.1c closure, roadmap §v0.1c) is **NOT YET PLANNED**. Tech Lead phase start
> owns the validation contract + task list before any task work begins.

## Phase 9 goal (roadmap §v0.1c — Concerns + v0.1c closure)

Ship the **concern write surface** linked to completions, surface the derived
concern flags on issue reads, wire real concerns into the existing
`concerns_addressed` transition rule, and run the end-to-end v0.1c smoke that
closes the slice.

**Exit (from roadmap):** an agent can raise a concern against a completion, see
open concerns light up the `issues_enriched` derived flag on issue reads, address
it and watch the flag clear, and close out the v0.1c lifecycle end-to-end over
ac add/verify/list → complete → concern raise/address/list.

## Scope (roadmap items)

- `fid concern raise` — structured concern linked to issue + completion.
  `--on <issue-id> --by developer --in-completion <id> --summary --body-md`
  (PRD-DEV-005). Depends on Phase 8 completions for a real `--in-completion`
  endpoint.
- `fid concern address` — transition concern to addressed `--in-completion <id>`
  (PRD-DEV-006).
- `fid concern list` — filter by issue and status (PRD-DEV-006/007).
- Derived concern flag — `issues_enriched` exposes `concerns_raised` and
  `open_concerns_count` on issue reads (PRD-DEV-007).
- Real `concerns_addressed` transition-rule wiring, replacing Phase 5 fixtures.
- End-to-end v0.1c smoke across the whole slice.

Explicit roadmap deferrals: concern `escalated`/`wontfix` lifecycle beyond
raise/address remains v0.1e/later; only the PRD-DEV-006 statuses are modelled.

## Planning inputs the Tech Lead must honor

- Validate-before-append (CF-007): all rejections (e.g. concern on a bad target,
  unknown role, missing completion) occur before any journal append — non-zero
  exit, clear error, zero new lines/rows.
- `RequireKnownRole` permissions only (matrix is v0.1d), matching Phase 7/8.
- DEC-D / CF-011 precedent: decide the `activity_log` stance for concern writes
  (the create/AC/completion precedent is no `activity_log`; consistent timeline
  population is the v0.1e CF-011 job). State it explicitly in the contract.
- ADR-0001 rebuild parity: include a delete-DB + `fid rebuild` parity assertion
  for concerns + derived flags.
- Human dogfood VAL for the concern lifecycle feel (per-observable rows).
- Coverage (DEC-007): COV-001 `internal/issues|journal|policy|workflow` ≥ 80%;
  COV-002 `cmd/fid` ≥ 70%. Phase 8 ended `cmd/fid` 78.7%, `internal/issues`
  81.9%, `internal/journal` 80.3%, `internal/policy` 85.0%, `internal/workflow`
  95.5%.
- Task-spec authoring nudges (carried): scope any file-size acceptance to the
  task's allowed write set, not the whole package; and require finalized,
  self-contained report evidence (literal transcripts/metadata) in the
  code-behavior commit, not a later docs commit.

## Carry-Forward Register (open)

| ID | Item | Disposition |
|---|---|---|
| CF-003 | Role-validation corrupt-tail tolerance | Before v0.1d role/permission work. |
| CF-008 | Resolution type / cancellation reason (DEC-015) | v0.2 (product). |
| CF-011 | Consistent `activity_log` timeline population (DEC-020) | v0.1e (with `fid activity list`). |
| CF-013 | `internal/issues/complete.go` `Complete` orchestration size | v0.1e, alongside `fid dev finish`/`block`. |

CF-012 closed at Phase 8 wrap (Task 4 / VAL-009): `issue.go` split into
`create.go` 387 / `get.go` 120 / residual `issue.go` 192.

## Next Checkpoint

Tech Lead: run Phase 9 phase-start — read the roadmap §v0.1c Phase 9 entry +
relevant PRD-DEV-005..007 / PRD-WF-009 spec sections, draft
`project-docs/docs/phases/phase-9/qa/validation-contract.md`, break into ordered
one-session tasks, and present the contract + task list to the human for
approval before Phase Runner starts Task 1.
</content>
