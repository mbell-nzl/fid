# Phase 10 Validation Contract - Permissions Matrix + Policy

> Phase: 10 - Permissions matrix + policy + CF-003 (roadmap §v0.1d, slice 1 of 4)
> Author: Tech Lead (draft)
> Last refined: 2026-06-21 by Tech Lead (awaiting QA Lead evidence refinement)
>
> Location: `project-docs/docs/phases/phase-10/qa/validation-contract.md`.

## Purpose

This file is the finite set of validation assertions that must be true when
Phase 10 is complete. The contract specifies **what must be proven**; the QA
plan specifies **how**.

If a required assertion cannot be satisfied, stop and escalate rather than
weakening it.

## Phase scoping notes

Phase 10 is the first slice of v0.1d (roadmap §v0.1d phase-slicing note,
2026-06-21: permissions → claims → requirements+links → labels). It introduces
the **permissions matrix as a cross-cutting gate on every write command** — the
first version whose writes are gated by role/action and not `RequireKnownRole`
alone. It ships no new entity; it builds the in-code matrix, the enforcement
seam, `fid policy show`, and retrofits enforcement onto every existing write
surface. It also reconciles the role-discovery corrupt-tail tolerance (**CF-003**,
explicitly due "before v0.1d role/permission work").

The `claims`, `links`, `labels`, `issue_labels`, and `requirements` tables and
their triggers already exist (v0.1a schema-widening); Phase 10 adds **no schema
migration**. The permissions matrix is an in-code constant (TRD §13,
TRD-ARCH-010); the `role_permissions` table is v0.2.

Planning decisions confirmed by the human (2026-06-21): DEC-A..E.

- **DEC-A — enforcement seam in `internal/policy`, action key derived in the
  domain layer.** Add `policy.Authorize(workspaceRoot, role, action string)
  error` that subsumes `RequireKnownRole` (preserving `MissingRoleError` and
  `UnknownRoleError` behaviour) and adds an in-code matrix lookup, returning a
  new `PermissionDeniedError{Role, Action}` for a denied combination. Action
  keys are derived in the domain layer (`internal/issues`, `cmd/fid/rebuild.go`)
  from the already-known kind/target-state/role, keeping the CLI adapter thin
  and the business rule in core (AGENTS.md). No new architecture boundary.
- **DEC-B — the in-code matrix is the single source of truth, shared by
  `Authorize` and `fid policy show`, and lists only actions for implemented
  write commands.** Phase 10 seeds the v0.1a–v0.1c write actions
  (`issue.create.{epic,task,subtask}`, `issue.update`, `issue.complete`,
  `issue.ac.add`, `issue.ac.verify.{developer,qa_reviewer,qa_lead,tl,human}`,
  `issue.transition.{in_progress,dev_complete,failed_gate,qa_review,qa_pass,
  qa_block,done,cancelled}`, `concern.raise`, `concern.address`,
  `maintenance.rebuild`). Phases 11–13 append their action rows alongside their
  commands. TRD-ARCH-010 completeness is tested against implemented actions.
- **DEC-C — `fid policy show` is a read command:** no `--by`, no permission
  check, no journal write, no `activity_log`; renders the matrix in all four
  formats (json default, toon, markdown, yaml).
- **DEC-D — CF-003 reconciliation aligns role discovery with the journal
  corruption contract.** `KnownRoles`/`collectRoles` (`internal/policy`) report a
  malformed or un-upgradeable journal line as corruption (non-zero exit, clear
  error naming the journal file) rather than silently tolerating it after seeded
  roles are found, matching the replayer (TRD-ARCH-002) and the journal-canonical
  ADR-0001. The existing `TestKnownRolesTRG002CorruptTailTolerance` "tolerated"
  subtest is rewritten (decision-backed, red/green) to assert reported
  corruption — a correction of acknowledged Phase 1 debt, not a test weakening.
- **DEC-E — the retrofit moves existing test fixtures to permitted role/action
  combinations.** Enforcing the matrix changes the behaviour of existing
  commands (e.g. `developer` may not create an epic; `qa_reviewer` may not
  transition to `dev_complete`). Existing tests/fixtures that used now-denied
  combinations are updated to permitted roles (or `--by tl`/`--by human`),
  recorded as a deliberate, spec-mandated behaviour change, not test weakening.
  Validate-before-append (CF-007) holds for every permission rejection.

`fid init` is bootstrap (it seeds the role set) and is exempt — there is no role
to check before roles exist; it takes no `--by` and no matrix action.

## Numbering

- IDs use `VAL-NNN`, scoped to Phase 10, sequential from `VAL-001`.
- Carry-forward assertions are re-issued with a new `VAL-NNN` and a source note.
- Retired/superseded assertions stay in the table; IDs are never reused.

## Assertions

| ID | Source requirement | Behavior or invariant | Evidence method | Validator | Owning task/window | Pass criterion | Fail criterion | Status | Out-of-scope notes |
|---|---|---|---|---|---|---|---|---|---|
| VAL-001 | CF-003; TRD-ARCH-002; ADR-0001 | Role discovery (`KnownRoles`/`RequireKnownRole`/`Authorize`) reports a malformed or un-upgradeable journal line as corruption (non-zero exit, clear error naming the journal file) instead of silently tolerating it after seeded roles are found; a well-formed journal with valid role events still resolves the role set successfully. | Unit/integration: (a) journal with a valid role event + a malformed trailing line → `KnownRoles` returns a corruption error naming the file; (b) clean journal → roles resolved; plus a CLI write on a corrupt-tail workspace → non-zero exit + corruption error. | QA Lead | task-1 | Corrupt-tail journal is reported as corruption naming the file (non-zero exit); a clean journal resolves roles with no error. | A corrupt tail is silently tolerated/skipped, or a clean journal fails to resolve roles. | Draft | Replayer/rebuild corruption behaviour itself is unchanged; this only reconciles the lighter role-discovery scan with that contract. |
| VAL-002 | PRD-ROLE-005; TRD-ARCH-010; design §11 `fid policy show` | `fid policy show` renders the in-code permissions matrix as one row per action with a boolean per role, identically across all four formats (`json` default, `toon`, `markdown`, `yaml`), and the output includes every implemented v0.1 write action. | CLI: `fid policy show` in each format; assert the action×role booleans render in each and that the implemented write-action set is fully present. | QA Lead | task-2 | Every format renders the matrix (action + per-role booleans) with the same data, and every implemented write action appears. | A format omits/diverges, drops the per-role booleans, or an implemented write action is missing. | Draft | Only implemented (v0.1a–c) actions appear at Phase 10; claim/link/label/requirement.sync rows arrive with their commands in Phases 11–13. |
| VAL-003 | TRD-ARCH-010 | Every implemented write command's permission action has exactly one matrix entry, and `policy.Authorize` returns no error for a permitted (role, action) and a `PermissionDeniedError` naming the role and action for a denied (role, action). | Unit: table-driven test asserting each implemented action key resolves to a matrix row; `Authorize` permits a known allowed combination and returns `PermissionDeniedError` for a known denied combination. | QA Lead | task-2 | All implemented action keys have a matrix entry; `Authorize` permits allowed and denies denied with `PermissionDeniedError`. | An implemented action has no matrix entry, `Authorize` permits a denied combination, or denies a permitted one. | Draft | Wildcard/condition permission rows are out of scope (v0.1 matches actions exactly per TRD §13). |
| VAL-004 | PRD-ROLE-003; TRD-ARCH-010; CF-007 | A write command invoked with a known role that the matrix denies for the derived action is rejected **before any journal append** (non-zero exit, clear error naming the role, the action, and that it is not permitted; zero new journal lines/rows); the same command invoked with a permitted role succeeds and appends exactly one event. | CLI over a representative denied write (e.g. `issue create --kind epic --by developer`) asserting zero new journal lines + non-zero exit + role/action-named error; plus the permitted case (`--by tl`) appending exactly one event/row. | QA Lead | task-3 | The denied write is rejected with a role/action-named error and zero append; the permitted write appends one event and succeeds. | A denied write appends a line/row or succeeds, the error omits role or action, or a permitted write is rejected. | Draft | Validate-before-append holds for the permission gate exactly as it does for the existing domain validations (CF-007). |
| VAL-005 | PRD-ROLE-001; PRD-ROLE-002; CF-007 | Role-identity validation is preserved through the new `Authorize` seam: a write with a missing `--by` is rejected with `--by is required`, and a write with an unknown role is rejected with an error listing the known roles; both reject before any journal append. With `--by <permitted role>` the same command proceeds. | CLI: one write command with no `--by` (missing-role error, zero append) and with an unknown role (unknown-role error listing known roles, zero append); plus the permitted-role success case. | QA Lead | task-3 | Missing `--by` and unknown-role are each rejected before append with their distinct errors; the permitted-role case proceeds. | Either reject path appends/regresses, loses its error, or the permitted-role case is rejected. | Draft | Missing and unknown role are two reject paths of the same identity gate (Phase 9 VAL-002 precedent for clustering reject paths). |
| VAL-006 | PRD-ROLE-004 | `--by human` is permitted on every write command (the product-owner escape valve), including a write a non-human role is denied; and for role-attributed AC verification `--by human` records `human` as the verifying role via action `issue.ac.verify.human` — it does not mint a `developer`/`qa_lead` verification. | CLI: a write denied for a non-human role (e.g. `issue.create.epic` by `developer`) succeeds with `--by human`; plus `issue ac verify --by human` whose persisted/read-back verification records `_role = human` (not another verifier). | QA Lead | task-3 | The human escape valve succeeds where a non-human role is denied, and a human AC verification is attributed to `human`, not another verifier role. | `--by human` is denied on any write, or a human AC verification is recorded/read back as `developer`/`qa_lead`/etc. | Draft | Human still cannot impersonate another verifier role; the escape valve grants invocation, not identity laundering. |
| VAL-007 | PRD-ROLE-003; TRD-ARCH-010 | Permission checks key on the **derived** action: `issue.ac.verify.<role>` (so `issue ac verify --by qa_reviewer` is permitted but `--by workflow_runner` is denied) and `issue.transition.<target-state>` (so `issue transition --to dev_complete --by developer` is permitted but `--by qa_reviewer` is denied). | CLI: AC verify permitted for `qa_reviewer` and denied for `workflow_runner`; transition to `dev_complete` permitted for `developer` and denied for `qa_reviewer` — each with the matching outcome and (for denials) a role/action-named error and zero append. | QA Lead | task-3 | Role-keyed AC-verify and state-keyed transition actions resolve to the correct matrix row, permitting the named qualifying role and denying the named non-qualifying role. | A role-keyed or state-keyed action resolves to the wrong row (e.g. any verify role passes, or any role may transition to `dev_complete`). | Draft | Only the v0.1 default-graph target states and the five AC-verify verifier roles are in scope. |
| VAL-008 | TRD-ARCH-010; phase brief (cross-surface parity) | The matrix `fid policy show` displays and the matrix `Authorize` enforces are the same: for a representative denied cell, `fid policy show` displays it as not-permitted **and** the corresponding write is rejected; for a permitted cell, `fid policy show` displays it as permitted **and** the corresponding write succeeds. | CLI: pick one permitted and one denied (role, action) cell; assert the `fid policy show` boolean matches the actual write outcome for each. | QA Lead | task-3 | Displayed permission and enforced behaviour agree for both the permitted and the denied cell. | `fid policy show` displays a cell as permitted that is enforced as denied (or vice versa). | Draft | Single-source matrix (DEC-B) makes this structural; the check guards against display/enforcement drift, a standing cross-surface parity check. |
| VAL-009 | phase brief; roadmap §v0.1d | Human confirms on a real `.fid/` workspace that `fid policy show` reads as a clear role×action matrix and that a denied write produces an actionable error naming the role, the action, and that it is not permitted, so the operator knows how to proceed. | Human checklist against a real `.fid/` workspace; record outcomes per observable. | Human | phase-end | Human confirms the matrix is readable and a denied-write error is actionable (names role, action, not-permitted, with a usable next step). | Human finds the matrix unreadable/ambiguous, or a denial error that does not name the role/action or gives no actionable signal. | Draft | QA Lead fans this into per-observable human-checklist rows (policy-show readable / denial error actionable). |

## Coverage Rules

Confirmed before task execution begins:

- Phase-scope requirements each map to ≥1 `VAL-NNN`: PRD-ROLE-001→VAL-005;
  PRD-ROLE-002→VAL-005; PRD-ROLE-003→VAL-004/007; PRD-ROLE-004→VAL-006;
  PRD-ROLE-005→VAL-002; TRD-ARCH-010→VAL-002/003/004/007/008; CF-003→VAL-001;
  CF-007 validate-before-append→VAL-004/005.
- Every `VAL-NNN` is referenced by a task spec or marked `phase-end` (VAL-009).
- Binary/conditional assertions name both states — VAL-001, 002, 003, 004, 005,
  006, 007, 008.
- DEC-A..E are human-confirmed (2026-06-21); the contract and tasks cite them.
- The in-code matrix lists only implemented write actions (DEC-B); VAL-002/003
  assert completeness against the implemented set, not against future
  claim/link/label/requirement.sync actions (Phases 11–13).
- Coverage (DEC-007) is enforced at the gate, not as a `VAL-NNN`: COV-001
  `internal/issues|journal|policy|workflow` ≥ 80%; COV-002 `cmd/fid` ≥ 70%. Each
  task runs the coverage gate and the completion report records the percentages.
- QA Lead maps each `VAL-NNN` to deterministic `(qa-lead)`, mid-phase-review, or
  `(human)` verification and fans VAL-009 into per-observable human rows. The QA
  plan carries the two Phase 9 standing checks (cross-surface read-back parity —
  here show-vs-enforce, VAL-008; and validate-before-append on every reject path,
  VAL-004/005).

## Outcomes

Phase-end final review records every `VAL-NNN` outcome: `PASS`, `FAIL`,
`BLOCKED`, or `CARRY-FORWARD`. `PASS` without evidence is invalid.

## Revision Log

| Date | Assertion | Change | Reason | Decision reference |
|---|---|---|---|---|
| 2026-06-21 | VAL-001..009 | created | Phase 10 planning draft | phase brief; roadmap §v0.1d; PRD-ROLE-001..005; TRD-ARCH-010; CF-003; CF-007 |
| 2026-06-21 | n/a | DEC-A..E confirmed | Human accepted all five Tech Lead planning decisions pre-execution | this contract §Phase scoping notes |
