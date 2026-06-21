# Phase Brief - Phase 10 CLOSED → Phase 11 next (v0.1d slice 2/4)

> **Purpose.** Default execution brief for the current phase. Loaded by every
> agent session. Keep under 150 lines.
>
> **State:** Phase 10 (Permissions matrix + policy + CF-003) is **CLOSED** on
> branch `phase-10`, 2026-06-21. **v0.1d slice 1 of 4 complete.** Phase 11
> (Claims layer) is the next target and is **not yet planned** — no validation
> contract or task specs exist yet. The Tech Lead plans Phase 11 next.

## Phase 10 outcome (closed)

Introduced the **permissions matrix as a cross-cutting gate on every write
command** — the first version gated by role/action, not `RequireKnownRole` alone.
Delivered the in-code matrix (TRD §13), `policy.Authorize` + `PermissionDeniedError`,
read-only `fid policy show` (json/toon/markdown/yaml), retrofit enforcement onto
every existing v0.1a–c write surface, and the CF-003 corrupt-tail reconciliation.
No schema migration.

- QA final review: **PASS for deterministic QA** (0 blocking / 0 major / 0
  improvement). Human dogfood **4/4 PASS** (HUMAN-09a..d). **VAL-001..009 all PASS.**
- Coverage at close: `cmd/fid` 77.8%, `internal/issues` 81.9%, `internal/journal`
  80.3%, `internal/policy` 81.2%, `internal/workflow` 95.5% (COV-001/COV-002 met).
- DEC-A..E (planning, human-confirmed) honoured; DEC-025 resolved ESC-007.
- **CF-003 CLOSED** (Task 1). `03_TRD_SCHEMA.md` §13 touched (4 direct-prerequisite
  matrix rows); all other specs verified-no-change or N/A.
- Read-ordering rule promoted to `decisions_principles.md` (Phase 9 deep-retro item
  discharged).

Archive: `project-docs/docs/phases/phase-10/` (decisions, progress, tl-escalations,
triage-log, qa/, tasks/).

## Phase 11 goal (to plan)

**Claims layer**, all gated by the Phase 10 matrix (roadmap §v0.1d slice 2):

- `fid issue claim` — atomic acquire via `INSERT ... ON CONFLICT DO NOTHING`,
  default 30m TTL, `--ttl`, `--renew` for the holder, `--force` gated by
  `issue.claim.force_acquire`, `--note`.
- `fid issue unclaim` — release caller's current claim (`--by` required).
- `fid claim list` — live claims; `--stale` for past-TTL; `--by`; all four formats.
- Claim–transition integration — started-category transitions require a live caller
  claim; transitions out auto-release.
- Human claims — `--by human` as a do-not-advance signal.

Scope refs: PRD-CLAIM-001–008a, TRD-ARCH-011. The `claims` table + triggers already
exist (v0.1a schema-widening); no schema migration expected.

**Exit (roadmap):** an agent can claim/renew/force/release with correct TTL and
atomic conflict rejection, list live and stale claims, drive a started-category
transition only with a live claim, and see the claim auto-released on the way out.

## Planning inputs the Phase 11 Tech Lead must carry

- **New matrix actions.** Phase 11 appends its `issue.claim.*` /
  `issue.claim.force_acquire` action rows to the in-code matrix (DEC-B: matrix lists
  only implemented write actions and grows per phase) and to `03_TRD_SCHEMA.md` §13;
  TRD-ARCH-010 completeness is tested against the implemented set.
- **Standing QA-plan checks** (carry forward): show-vs-enforce cross-surface parity
  for the new claim actions, and validate-before-append on every claim rejection
  (claim conflict, force denied, missing/unknown role) — non-zero exit, role/action
  named, zero new journal lines/rows.
- **Append-order read ordering** is now a standing principle
  (`decisions_principles.md`): `fid claim list` ordering defaults to journal append
  order (`rowid`) unless the spec names another order.
- **CF-011 discipline.** Claim renewal (PRD-CLAIM-008a) projects to `activity_log`;
  follow the existing per-command discipline. CF-011's broader reconciliation stays
  v0.1e.
- **Commit-field hygiene (Phase 10 retro).** Pair each task's completion-report
  commit-hash fill with an immediate reconciliation commit (or a deferred-hash
  placeholder reconciled in the task commit) so commit-field drift does not resurface
  as a phase-end blocker (DEC-025 recurrence note).

## Carry-Forward Register (open)

| ID | Item | Disposition |
|---|---|---|
| CF-008 | Resolution type / cancellation reason (DEC-015) | v0.2 (product). |
| CF-011 | Consistent `activity_log` timeline population (DEC-020) | v0.1e. v0.1d writes that project to `activity_log` (claim renewal in Phase 11, link removal, label changes) follow the existing per-command discipline; CF-011's reconciliation stays v0.1e. |
| CF-013 | `complete.go` `Complete` orchestration size | v0.1e (with `fid dev finish`/`block`). |

CF-003 closed at Phase 10 wrap (Task 1, VAL-001 PASS).

## Next Checkpoint

Tech Lead plans Phase 11 (Claims): read the roadmap §v0.1d Phase 11 entry and
relevant PRD-CLAIM/TRD-ARCH-011 sections, draft the Phase 11 validation contract
(VAL-NNN) and ordered task specs under
`project-docs/docs/phases/phase-11/`, and present the contract + task list to the
human before the dev/QA loop runs.
