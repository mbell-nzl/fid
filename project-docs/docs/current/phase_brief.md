# Phase Brief - Phase 10 (Permissions Matrix + Policy — v0.1d slice 1/4)

> **Purpose.** Default execution brief for the current phase. Loaded by every
> agent session. Keep under 150 lines.
>
> **State:** Phase 10 is **PLANNED and READY for the Phase Runner** — DEC-A..E
> human-confirmed 2026-06-21. v0.1c is closed (Phase 9). The roadmap §v0.1d
> phase-slicing note (2026-06-21) cuts v0.1d into four phases: **10 permissions →
> 11 claims → 12 requirements+links → 13 labels+closure.** Validation contract +
> 3 task specs under `project-docs/docs/phases/phase-10/`. The Phase Runner spawns
> the QA Lead first to write the QA plan, then runs the dev/QA loop T1→T2→T3.

## Phase 10 goal

Introduce the **permissions matrix as a cross-cutting gate on every write
command** — the first version gated by role/action, not `RequireKnownRole`
alone. Build the in-code matrix (TRD §13), the `policy.Authorize` seam, and
`fid policy show`; retrofit enforcement onto every existing write surface; and
reconcile the role-discovery corrupt-tail tolerance (**CF-003**). No schema
migration — the v0.1 tables/triggers already exist.

Scope refs: PRD-ROLE-001–005, TRD-ARCH-010, CF-003, CF-007.

## Planning decisions (DEC-A..E — human-confirmed 2026-06-21)

- **DEC-A** — `policy.Authorize(ws, role, action)` subsumes `RequireKnownRole`
  and adds the matrix lookup, returning `PermissionDeniedError{Role, Action}`.
  Action keys are derived in the **domain layer** (`internal/issues`,
  `rebuild.go`) from known kind/target-state/role. No new boundary.
- **DEC-B** — the in-code matrix is the **single source of truth** for both
  `Authorize` and `fid policy show`; it lists only **implemented** write actions
  and grows per phase. Phase 10 seeds the v0.1a–c actions (DEC-B list in the
  contract); Phases 11–13 append theirs.
- **DEC-C** — `fid policy show` is a **read**: no `--by`, no journal write, no
  `activity_log`; renders json/toon/markdown/yaml.
- **DEC-D** — CF-003: role discovery **reports** a malformed/un-upgradeable
  journal line as corruption (aligns with the replayer / ADR-0001) instead of
  tolerating it after seeded roles are found. The `TestKnownRolesTRG002...`
  "tolerated" subtest is rewritten (decision-backed, not weakening).
- **DEC-E** — the retrofit moves existing test fixtures to **permitted** role/
  action combos (enforcement changes existing command behaviour; spec-mandated,
  not weakening). Validate-before-append (CF-007) holds on every denial.

## Ordered tasks

1. **Task 1 — CF-003 corrupt-tail reconciliation** (`internal/policy`). VAL-001.
2. **Task 2 — matrix + `Authorize` + `fid policy show`** (infra + new read
   command). VAL-002, VAL-003.
3. **Task 3 — retrofit enforcement onto every write command** (cross-cutting,
   high blast radius; fixtures move to permitted roles). VAL-004..008.

Human dogfood (VAL-009) is phase-end.

## Validation contract index

`project-docs/docs/phases/phase-10/qa/validation-contract.md` (VAL-001..009):

- VAL-001 (T1) — role discovery reports corrupt tail; clean journal resolves.
- VAL-002 (T2) — `fid policy show` all 4 formats; includes every implemented
  write action.
- VAL-003 (T2) — matrix completeness + `Authorize` permit/deny (unit).
- VAL-004 (T3) — denied write rejected before append (role/action-named error,
  zero lines); permitted succeeds.
- VAL-005 (T3) — missing `--by` + unknown role preserved through `Authorize`.
- VAL-006 (T3) — `--by human` escape valve; human AC verify attributed to
  `human`, no impersonation.
- VAL-007 (T3) — action-key derivation: AC-verify role-keyed, transition
  state-keyed (both qualifying/non-qualifying named).
- VAL-008 (T3) — `fid policy show` display matches enforced behaviour
  (cross-surface parity).
- VAL-009 (phase-end, Human) — matrix readable; denial error actionable.

## Carried QA-plan inputs (QA Lead must include)

- **Show-vs-enforce parity (VAL-008)** — the Phase 9 standing *cross-surface
  read-back parity* check, instantiated for the matrix: what `fid policy show`
  displays must equal what `Authorize` enforces.
- **Validate-before-append on every reject path (CF-007)** — every permission
  denial exits non-zero, names role+action, and appends **zero** journal
  lines/rows (assert explicitly, per VAL-004/005).
- **Coverage (DEC-007):** COV-001 `internal/issues|journal|policy|workflow` ≥
  80%; COV-002 `cmd/fid` ≥ 70%, enforced at each task gate. Phase 9 ended
  `cmd/fid` 77.8%, `internal/issues` 81.8%, `internal/policy` 85.0%.
- **Read-ordering principle promotion** — Phase 9 deep-retro rec: promote the
  append-order-not-`created_at` rule (DEC-019/021/024) toward a standing
  principle in `decisions_principles.md` during this phase wrap.

## Carry-Forward Register (open)

| ID | Item | Disposition |
|---|---|---|
| CF-003 | Role-validation corrupt-tail tolerance | **Scheduled — Phase 10 Task 1** (DEC-D). |
| CF-008 | Resolution type / cancellation reason (DEC-015) | v0.2 (product). |
| CF-011 | Consistent `activity_log` timeline population (DEC-020) | v0.1e. New v0.1d writes that the spec says project to `activity_log` (claim renewal, link removal, label changes) arrive in Phases 11–13 and follow the existing per-command discipline; CF-011's reconciliation stays v0.1e. |
| CF-013 | `complete.go` `Complete` orchestration size | v0.1e (with `fid dev finish`/`block`). |

## Next Checkpoint

Human confirms (or revises) DEC-A..E and the ordered task list. On confirmation,
QA Lead refines the validation contract evidence and writes the Phase 10 QA plan
(carrying the inputs above); then the dev/QA loop runs Tasks 1→2→3.
