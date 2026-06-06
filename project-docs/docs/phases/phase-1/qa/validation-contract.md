# Phase 1 Validation Contract - End-to-End Persistence Loop (Tracer Bullet)

> Phase: 1 - End-to-End Persistence Loop (Tracer Bullet)
> Author: Tech Lead (draft)
> Last refined: Not yet refined by QA Lead

## Purpose

This file is the finite set of validation assertions that must be true when
Phase 1 is complete. The validation contract is the specification of **what
must be proven**; the QA plan is the strategy for **how to prove it**.

Phase 1 is the v0.1a *tracer bullet*: the thinnest end-to-end path through the
journal → projection → rebuild machinery, exercised by real commands. The
hardest-to-change contracts (journal event v1 format, ID shapes, projection
checkpoint, write lock) are locked here per Roadmap Guiding Principle #3.

If a required validation assertion cannot be satisfied, stop and escalate
instead of weakening the assertion or silently adapting the spec.

## Numbering

- IDs use `VAL-NNN`, scoped to Phase 1 and sequential from `VAL-001`.
- Carry-forward assertions are re-issued with a new `VAL-NNN` and a source note.
- Retired or superseded assertions stay in the table; do not reuse IDs.

## Assertions

| ID | Source requirement | Behavior or invariant | Evidence method | Validator | Owning task/window | Pass criterion | Fail criterion | Status | Out-of-scope notes |
|---|---|---|---|---|---|---|---|---|---|
| VAL-001 | PRD-INIT-001; TRD CONFIG §1, §4; TRD §4.3 | `fid init` creates the `.fid/` workspace: `journal/` dir, `config.toml`, `pricing.toml`, `fid.sqlite`, and `.gitignore` that ignores `fid.sqlite`. | Integration test + file-state proof after `fid init` in a temp dir. | QA Lead | task 5 | All listed paths exist after init; `.gitignore` contains `fid.sqlite`; `config.toml`/`pricing.toml` parse with default content. | Any path is missing, `.gitignore` omits `fid.sqlite`, or a config file is absent/unparseable. | Draft | Export `[export]` config section (v0.4) is not written in v0.1. |
| VAL-002 | PRD-INIT-002; TRD CONFIG §2 | `fid init --prefix X` records prefix `X` in `config.toml`; with no flag the prefix defaults to the uppercased project directory name. | Integration test asserting `config.toml` `[project].prefix` for both flagged and default invocations. | QA Lead | task 5 | Flagged init stores the given prefix; default init stores the uppercased directory name. | Prefix is not stored, ignores the flag, or default differs from the uppercased directory name. | Draft | Prefix immutability after init beyond rejecting re-init is not separately tested here. |
| VAL-003 | PRD-INIT-003; TRD §4.3 bootstrap; TRD CONFIG §4 | `fid init` appends exactly one `project.initialized` event and projects it, seeding exactly the seven default roles (developer, qa_reviewer, qa_lead, tl_lite, tl, workflow_runner, human). | Integration test: assert one journal line of type `project.initialized`; query `roles` for the seven names. | QA Lead | task 5 | Exactly one `project.initialized` event exists and the `roles` table contains exactly the seven named roles. | More/fewer than one bootstrap event, or the roles set differs from the seven. | Draft | Role permission matrix rows (v0.1d/v0.2) are not seeded in Phase 1. |
| VAL-004 | PRD-INIT-004 | Re-running `fid init` in a directory that already has `.fid/` is rejected with a clear error and makes no changes; a fresh directory initializes successfully. | Integration test for both states; compare workspace state before/after the rejected re-init. | QA Lead | task 5 | Fresh dir init succeeds; re-init exits non-zero with a clear `.fid/ already exists` error and leaves the workspace unchanged. | Re-init succeeds, partially overwrites, or fails without a clear error. | Draft | Force re-init flags are not in scope. |
| VAL-005 | PRD-OPS-004; TRD-ARCH-001; TRD §4.3 rebuild | After `fid init`, deleting `fid.sqlite` and running `fid rebuild` reconstructs the seeded roles and `projection_meta` from the journal without consulting the old DB; the journal is unchanged. | Integration test: init, delete DB, rebuild, diff `roles` query and `projection_meta`; assert journal bytes unchanged. | QA Lead | task 5 | Roles and projection checkpoint after rebuild equal the post-init state; journal files are byte-identical. | Rebuild produces different roles/checkpoint, errors, or mutates the journal. | Draft | Multi-file journal performance targets are Phase 3. |
| VAL-006 | PRD-OPS-001; PRD-OUT-001; PRD-OUT-004 | A minimal issue created through the write path is projected to SQLite and read back as JSON via `fid issue get <id> --format json` (and the JSON default). | Integration test: create one issue, read back via `fid issue get`, assert JSON fields match the created values. | QA Lead | task 6 | Created issue is returned as well-formed JSON whose fields match what was written. | Issue is not persisted, not readable, or JSON fields disagree with the write. | Draft | Full v0.1b lifecycle (kinds hierarchy, workflows, transitions) and TOON/markdown/YAML formats are out of scope. |
| VAL-007 | TRD-ARCH-004; TRD §4.3 event v1 contract; PRD-OPS-001 | One logical CLI write produces exactly one journal JSONL line carrying the v1 envelope: `event_id`, `event_schema` (=1), `event_type`, `occurred_at`, `actor`, `command`, ordered `effects`. | Integration test: count journal lines after one write; assert the parsed envelope has all required v1 fields. | QA Lead | tasks 3, 6 | Each single write adds exactly one well-formed v1 event line with all required fields. | A write adds zero, multiple, or partially-formed event lines, or omits a required envelope field. | Draft | Composite multi-effect commands (`fid dev finish`) are later phases; effects array is exercised but kept minimal. |
| VAL-008 | TRD-ARCH-001; TRD-ARCH-002 | Replay is pure data: after creating an issue, deleting `fid.sqlite`, and rebuilding, issue query results are identical, with no command logic re-run and no DB reads during effect application. | Integration test: create issue, delete DB, rebuild, assert identical `fid issue get` output; source review confirms replay applies resolved effects only. | QA Lead | task 4 | Post-rebuild issue read equals pre-delete read; replay applies stored effects without re-deriving values or reading prior DB state. | Read differs after rebuild, or replay re-runs command logic / reads DB state to resolve effect values. | Draft | Schema upgraders beyond v1 are exercised structurally only (see VAL-009). |
| VAL-009 | TRD-ARCH-003; TRD §4.3 event schema versioning | Replay routes every event through the version-upgrader chain before applying effects; a golden-file v1 journal replays to the expected DB state and the chain is a verified no-op at schema v1. | Golden-file replay test + unit test that the chain dispatches by `event_schema`. | QA Lead | task 4 | Golden v1 journal replays to expected state; chain is invoked and is identity at v1. | Effects are applied without version routing, or the golden replay output mismatches. | Draft | No N→N+1 upgrader exists yet; only the mechanism and v1 path are proven. |
| VAL-010 | TRD §4.3 malformed journal; PRD-OPS-003 | A malformed trailing JSONL line is reported as corruption (non-zero exit naming file and line) and never silently skipped or truncated; a well-formed journal replays cleanly. | Integration test for both a corrupt-tail journal and a clean journal. | QA Lead | task 4 | Corrupt tail produces a non-zero error naming the file/line; clean journal replays without error. | Corrupt tail is skipped/truncated/ignored, or a clean journal is wrongly flagged. | Draft | Mid-file corruption recovery tooling is future work. |
| VAL-011 | PRD-OPS-012; TRD §4.3 preconditions | When an effect precondition does not match during replay, replay stops with a structured logical-conflict error naming event_id, event_type, table, pk, expected, actual, and journal file/line; matching preconditions apply cleanly. | Integration test with a crafted conflicting precondition and a matching-precondition control. | QA Lead | task 4 | Mismatch stops replay with all named fields; matching precondition applies the effect. | Mismatch is applied anyway, or the error omits any required field. | Draft | v0.2 `fid merge --journal` reconciliation is out of scope. |
| VAL-012 | TRD §4.3 rebuild/checkpoint; PRD-OPS-003 | `projection_meta` governs read freshness: a checkpoint matching the journal serves reads from SQLite without rebuild; a missing checkpoint, schema-version mismatch, position mismatch, or last-event-hash mismatch forces rebuild/catch-up before serving. | Integration tests covering clean-match, missing-meta, and hash/position-mismatch cases. | QA Lead | task 4 | Matching checkpoint serves without rebuild; each mismatch class triggers rebuild/catch-up before answering. | Stale data is served on mismatch, or a clean match triggers an unnecessary full rebuild error. | Draft | Performance of catch-up vs full rebuild is Phase 3. |
| VAL-013 | PRD-NFR-008; TRD-PERF-006; TRD §4.3 concurrency | Same-checkout writes are serialized by `.fid/write.lock` (a lock dir with PID/hostname/started_at/command/role metadata) without data loss; a held lock causes a waiting caller to exit non-zero with an error naming the holder and path. | Integration test: two concurrent writes both persist; separate test asserts a held lock yields a holder-named error. | QA Lead | task 3 | Concurrent writes both land as distinct events with no loss; a held lock produces a clear holder-named error. | A concurrent write is lost/corrupted, or a held lock is silently broken. | Draft | v0.1 never silently breaks a live lock; cross-worktree text merge is Phase 3+. |
| VAL-014 | TRD-ARCH-007; TRD §4.4 | Every SQLite connection opens with `PRAGMA foreign_keys = ON` and `PRAGMA recursive_triggers = ON`. | Unit/integration test querying both pragma values on a freshly opened connection. | QA Lead | task 1 | Both pragmas report ON on every opened connection. | Either pragma is OFF or unset on any connection path. | Draft | Trigger DDL itself ships in Phase 2; only the pragmas are asserted here. |
| VAL-015 | TRD §4.2; TRD SCHEMA §1; TRD §4.3 event IDs | Generated IDs match their shapes: issue IDs `<PREFIX>-<hash>` (6 hex default), event IDs `EVT-<UTC YYYYMMDD>-<sequence>` allocated per UTC journal file under the lock; insert collisions retry with a fresh hash. | Unit tests for ID shape/sequence allocation + a forced-collision retry test. | QA Lead | tasks 2, 3 | IDs match the documented shapes; per-day event sequence increments correctly; a simulated collision retries and succeeds. | An ID violates its shape, sequence numbering is wrong, or a collision is not retried. | Draft | Progressive 6→8 scaling thresholds are exercised by unit test only; large-scale tuning is later. |
| VAL-016 | PRD-ROLE-001; PRD-ROLE-002 | Every Phase 1 write command (`fid issue create`, `fid rebuild`) requires `--by <role>` and rejects a role absent from the `roles` table with a clear error; a seeded role is accepted. | Integration test: missing `--by`, unknown role, and valid role for each write command. | QA Lead | task 6 | Missing `--by` and unknown roles are rejected with clear errors; a seeded role is accepted. | A write proceeds without `--by` or with an unknown role. | Draft | The full role×action permission matrix (PRD-ROLE-003) is deferred to v0.1d; Phase 1 checks role existence only. `fid init` is the bootstrap exception and takes no `--by`. |
| VAL-017 | TRD-ARCH-005; TRD §2 architecture | `internal/` packages never import `cmd/fid/`, and `cmd/fid/` holds no business logic (only flag parsing, delegation, formatting, exit codes). | Import-graph check (`go list`/grep) + source review of `cmd/fid/`. | QA Lead | phase-end | No `internal/` package imports `cmd/fid/`; `cmd/fid/` contains no domain logic. | Any `internal/` package imports `cmd/fid/`, or domain logic lives in `cmd/fid/`. | Draft | Per-package internal layering beyond this rule is reviewed but not separately asserted. |

## Field Guidance

- Source requirement cites the durable PRD/TRD/design source justifying the
  assertion.
- Evidence method names the smallest reproducible proof; QA Lead expands these
  into QA-plan steps.
- Validator is one role; all deterministic assertions here use `QA Lead`.

## Coverage Rules

Before task execution begins:

- Every Phase 1 roadmap item is referenced by at least one `VAL-NNN`, or the
  phase brief records why no assertion applies.
- Every `VAL-NNN` is referenced by at least one task spec or marked
  `phase-end`.
- Every Developer task spec lists the `VAL-NNN` IDs it advances or states why
  none apply.
- Every `VAL-NNN` has a planned evidence method, one validator role, and an
  owning task/window.
- QA Lead maps every `VAL-NNN` to deterministic verification in the QA plan.

Gaps are QA-plan defects and block implementation until resolved.

## Outcomes

Phase-end final review records every `VAL-NNN` outcome as `PASS`, `FAIL`,
`BLOCKED`, or `CARRY-FORWARD`. `PASS` without evidence is invalid.

## Revision Log

| Date | Assertion | Change | Reason | Decision reference |
|---|---|---|---|---|
| 2026-06-06 | VAL-001..VAL-017 | Created | Phase 1 planning kickoff (v0.1a tracer-bullet slice) | N/A |
