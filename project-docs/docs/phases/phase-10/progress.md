# Phase 10 Progress (archived)

> Phase: 10 - Permissions matrix + policy + CF-003 (roadmap §v0.1d, slice 1 of 4)
> Status: CLOSED on branch `phase-10`, 2026-06-21. **v0.1d slice 1 of 4 complete.**

## Outcome

Phase 10 introduced the **permissions matrix as a cross-cutting gate on every
write command** — the first version gated by role/action, not `RequireKnownRole`
alone. It built the in-code matrix (TRD §13), the `policy.Authorize` seam, the
read-only `fid policy show`, retrofitted enforcement onto every existing v0.1a–c
write surface, and reconciled the role-discovery corrupt-tail tolerance (CF-003).
No schema migration.

QA final review was **PASS for deterministic QA** (0 blocking / 0 major / 0
improvement); VAL-009 was CARRY-FORWARD to the human checklist and is now **PASS**
(HUMAN-09a..d 4/4 PASS, overall PASS). **VAL-001..VAL-009 all PASS** at wrap.

## Delivered

- **Task 1 — CF-003 corrupt-tail reconciliation** (`716329c`). `KnownRoles`/
  `collectRoles` (`internal/policy`) report a malformed or un-upgradeable trailing
  journal line as corruption with journal file/line context, matching the replayer
  (TRD-ARCH-002) and ADR-0001; TRG-002 test rewritten decision-backed (DEC-D).
  (VAL-001)
- **Task 2 — matrix + `Authorize` + `fid policy show`** (`6f4ce6b`). In-code
  single-source matrix (`internal/policy/matrix.go`), `policy.Authorize` +
  `PermissionDeniedError` (`authorize.go`), and the read-only `fid policy show`
  rendering json/toon/markdown/yaml (`cmd/fid/policy.go`). (VAL-002/003)
- **Task 3 — retrofit enforcement onto every write command** (`5293fd2`).
  `Authorize` on issue create/update/transition/complete, ac add/verify, concern
  raise/address, and rebuild; domain-layer action derivation
  (`internal/issues/policy_actions.go`); `--by human` escape valve; DEC-E fixture
  moves to permitted roles; four direct-prerequisite matrix rows added. (VAL-004..008)
- **Commit-field reconciliation** (`eeaf9fa`). `docs(task): reconcile phase-10 task
  report commit fields` cleared the Review 3 blocker (ESC-007 / DEC-025).

## Decisions

- DEC-A..E (planning decisions, human-confirmed 2026-06-21) all honoured.
- DEC-025 (ESC-007 resolution: reconcile uncommitted task-report commit fields via
  a §10-sanctioned `docs(task):` commit) issued mid-wrap. See `decisions.md`.
- Spec sync: `03_TRD_SCHEMA.md` §13 touched (four direct-prerequisite matrix rows);
  all other specs verified-no-change or N/A. See `decisions.md` Spec Sync Verdict.

## Escalations

- **ESC-007 — RESOLVED.** Review 3 BLOCK (task-report commit-field drift); direction
  DEC-025; reconciliation commit `eeaf9fa`; Review 3 addendum verified clearance.
  Archived in `tl-escalations.md`.

## Carry-forward

- **CF-003 — CLOSED** (delivered by Task 1, VAL-001 PASS).
- CF-008 (v0.2), CF-011 (v0.1e), CF-013 (v0.1e) re-deferred. No new CF opened.

## Coverage at close

`cmd/fid` 77.8%, `internal/issues` 81.9%, `internal/journal` 80.3%,
`internal/policy` 81.2%, `internal/workflow` 95.5% (COV-001/COV-002 satisfied).
Total statement coverage 76.3%.
