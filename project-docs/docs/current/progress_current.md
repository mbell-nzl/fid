# Progress - Current

> **Purpose.** In-flight phase status only. Loaded by default in every agent
> session. Historical phase progress lives under `../phases/phase-N/progress.md`.
>
> **Size target:** under 100 lines.

## Project State

**Status:** Phase 10 (Permissions matrix + policy + CF-003) is **CLOSED** on branch
`phase-10`. **v0.1d slice 1 of 4 complete.** QA final review: **PASS for
deterministic QA** (0 blocking / 0 major / 0 improvement). Human dogfood **4/4 PASS**
(HUMAN-09a..d). **VAL-001..009 all PASS** at wrap. Phases 0-9 archived under
`../phases/phase-N/`.

Phase 10 delivered:

- **Task 1 — CF-003 corrupt-tail reconciliation** (`716329c`): `KnownRoles`/
  `collectRoles` report a malformed/un-upgradeable trailing journal line as
  corruption with journal file/line context (DEC-D, VAL-001). **CF-003 closed.**
- **Task 2 — matrix + `Authorize` + `fid policy show`** (`6f4ce6b`): in-code
  single-source matrix, `policy.Authorize` + `PermissionDeniedError`, read-only
  `fid policy show` (json/toon/markdown/yaml) (VAL-002/003).
- **Task 3 — retrofit enforcement onto every write command** (`5293fd2`):
  `Authorize` on create/update/transition/complete, ac add/verify, concern
  raise/address, rebuild; domain-layer action derivation; `--by human` escape valve;
  DEC-E fixture moves; 4 direct-prerequisite matrix rows (VAL-004..008).
- **Commit-field reconciliation** (`eeaf9fa`): cleared the Review 3 blocker
  (ESC-007 / DEC-025).

Planning decisions DEC-A..E (human-confirmed pre-execution) all honoured; in-phase
DEC-025 resolved ESC-007. Spec sync: `03_TRD_SCHEMA.md` §13 touched (4 matrix rows);
all other specs verified-no-change or N/A. Read-ordering rule promoted to
`decisions_principles.md`. See `../phases/phase-10/decisions.md`.

Coverage held: `cmd/fid` 77.8%, `internal/issues` 81.9%, `internal/journal` 80.3%,
`internal/policy` 81.2%, `internal/workflow` 95.5% (COV-001/002 satisfied).

Open carry-forward register: CF-008 (resolution type/reason, v0.2), CF-011
(consistent `activity_log`, v0.1e), CF-013 (`complete.go` orchestration size,
v0.1e). CF-003 closed this phase.

## Next Action

- **Plan Phase 11 — Claims layer** (roadmap §v0.1d slice 2 of 4: 10 permissions →
  **11 claims** → 12 requirements+links → 13 labels+closure). `fid issue claim`
  (atomic acquire, 30m TTL, `--renew`/`--force`/`--note`), `fid issue unclaim`,
  `fid claim list` (`--stale`/`--by`/all formats), claim–transition integration,
  and human claims — all gated by the Phase 10 matrix. `claims` table/triggers
  already exist (v0.1a); no schema migration expected. (PRD-CLAIM-001–008a,
  TRD-ARCH-011)
- **Phase 11 appends its `issue.claim.*` matrix rows** to the in-code matrix and
  `03_TRD_SCHEMA.md` §13 (DEC-B: matrix grows per phase).
- **Carry the standing QA-plan checks:** show-vs-enforce cross-surface parity for
  new claim actions, and validate-before-append on every claim rejection
  (conflict / force-denied / missing / unknown role). `fid claim list` ordering
  defaults to journal append order (now a standing principle).
- **Commit-field hygiene (Phase 10 retro):** pair completion-report commit-hash fill
  with an immediate reconciliation commit so commit-field drift does not resurface as
  a phase-end blocker (DEC-025 recurrence note).
- Standing commitment before release-readiness: decide where held-out scenarios live
  and how implementation agents are prevented from loading them.
