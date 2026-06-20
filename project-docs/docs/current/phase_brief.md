# Phase Brief - Phase 10 (Claims, Permissions, Links — v0.1d)

> **Purpose.** Default execution brief for the current phase. Loaded by every
> agent session. Keep under 150 lines.
>
> **State:** Phase 9 (Concerns + derived flag + v0.1c closure) is **CLOSED**
> (VAL-001..011 all PASS; human dogfood 5/5; ESC-006 resolved by DEC-023/DEC-024 +
> Task 4). **v0.1c is complete.** Phase 10 (roadmap §v0.1d) is **NOT YET
> PLANNED** — Tech Lead phase start owns the slice, validation contract, and task
> list before any execution.

## Phase 10 goal (roadmap §v0.1d — Claims, Permissions, Links)

Ship the **claims layer**, the **permissions matrix**, and the **links/labels
surface** so the full issue lifecycle (create → claim → transition → complete →
unclaim) works with workflow validation and permission enforcement. v0.1d is the
first phase whose writes are gated by the permission matrix rather than
`RequireKnownRole` only.

**Roadmap §v0.1d items (TL slices these into coherent phases at phase start):**

- `fid issue claim` — atomic acquire with TTL (default 30m); `--by <role>
  --ttl`, `--renew` for the holder, `--force` for permitted force-acquire; reject
  on live claim unless renewing/stale/force. (PRD-CLAIM-001–004, 008a)
- `fid issue unclaim` — release caller's claim, `--by <role>`. (PRD-CLAIM-004)
- `fid claim list` — live claims; `--stale`, `--by <role>`. (PRD-CLAIM-007)
- Claim-transition integration — started-category transitions require a live
  claim by the caller; transitions out release it. (PRD-CLAIM-005–006)
- Human claims — `--by human` signals agents should not advance. (PRD-CLAIM-008)
- Permissions matrix — in-code constant for v0.1; enforce `--by <role>` on every
  write; reject unknown/denied role/action with actionable error; `fid policy
  show` displays the read-only matrix. (PRD-ROLE-001–005, TRD-ARCH-010)
- `fid requirement sync` (minimal) — parse PRD/TRD requirement IDs and upsert
  pointer rows so issue→requirement links have valid endpoints. (PRD-GOV-006,
  PRD-LINK-002)
- `fid link add/remove/list` — typed links with auto-inverse projection; v0.1
  scope issue↔issue and issue→requirement; relation/entity-type validation;
  reject links that duplicate dedicated schema. (PRD-LINK-001–006a, 008)
- `fid label add/list` + `fid issue label add/remove` — additive label CRUD and
  per-issue assignment without replacing the full set (`fid issue update --labels`
  stays replace-all). (PRD-ISSUE-015, 023)

**No pre-baked roadmap slice for v0.1d** — unlike v0.1a/b/c there is no
phase-slicing note. The Tech Lead must produce a vertical, agent-runnable slice at
phase start (candidate cut: claims → permissions/policy → requirement-sync + links
+ labels), each phase exiting on real CLI commands.

## Phase start checklist (Tech Lead — do before execution)

1. Read roadmap §v0.1d and the PRD-CLAIM / PRD-ROLE / PRD-LINK / PRD-ISSUE refs
   plus `03_TRD_SCHEMA.md` for `claims`, `links`, `labels`, `issue_labels`,
   `requirements`, and the claim/link/scope triggers.
2. **Schedule CF-003** — role-validation corrupt-tail tolerance is explicitly
   targeted "before v0.1d role/permission work"; decide do-now vs early-phase task.
3. Slice v0.1d into coherent phases; draft the validation contract first.
4. Carry the two Phase 9 standing QA checks into the QA plan (below).

## Carried-forward planning inputs (from prior wraps)

- **CF-003 (DUE)** — role-validation corrupt-tail tolerance, before v0.1d
  role/permission work. Schedule at Phase 10 planning.
- **Standing QA checks (Phase 9 deep-retro rec):** every write-surface phase QA
  plan must include (a) a cross-surface read-back parity check (a persisted field
  must be exposed on every read surface that should carry it) and (b) a read
  ordering determinism check (journal append order / `rowid DESC`, not
  `created_at`). These caught ESC-006 late in Phase 9.
- **Validate-before-append (CF-007)** on every reject path — claims, links, and
  label writes must reject before any journal append (non-zero exit, clear error,
  zero new journal lines/rows).
- **Permission enforcement is new this version** — v0.1a–v0.1c used
  `RequireKnownRole` only. v0.1d introduces the matrix; expect public-interface
  and testing-strategy implications (escalation triggers) — plan accordingly.
- **Coverage (DEC-007):** COV-001 `internal/issues|journal|policy|workflow` ≥ 80%;
  COV-002 `cmd/fid` ≥ 70% — enforced at the gate each task. Phase 9 ended
  `cmd/fid` 77.8%, `internal/issues` 81.8%.
- **Read ordering principle candidate** — DEC-019/DEC-021/DEC-024 are three
  append-order fixes; promote toward a standing principle in
  `decisions_principles.md` during Phase 10 (deep-retro rec).

## Carry-Forward Register (open)

| ID | Item | Disposition |
|---|---|---|
| CF-003 | Role-validation corrupt-tail tolerance | **DUE** — schedule at Phase 10 / v0.1d role/permission planning. |
| CF-008 | Resolution type / cancellation reason (DEC-015) | v0.2 (product). |
| CF-011 | Consistent `activity_log` timeline population (DEC-020) | v0.1e (with `fid activity list`). |
| CF-013 | `internal/issues/complete.go` `Complete` orchestration size | v0.1e, alongside `fid dev finish`/`block`. |

## Next Checkpoint

Tech Lead opens Phase 10 planning: read roadmap §v0.1d + refs, decide the slice,
schedule CF-003, and draft `project-docs/docs/phases/phase-10/qa/
validation-contract.md` before any task execution. Present the slice + contract +
task list to the human for approval.
