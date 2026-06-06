# Phase Brief - Phase 1: End-to-End Persistence Loop (Tracer Bullet)

> **Purpose.** Default execution brief for the current phase. Loaded by every
> agent session. Keep under 150 lines.
>
> **State:** PLANNED (2026-06-06). Validation contract + 6 ordered task specs
> drafted. Awaiting QA Lead QA plan and human review before Developer execution.

## Phase Goal

Build the thinnest **end-to-end** path through fid's persistence machinery —
journal → projection → rebuild — exercised by real commands. This is the v0.1a
tracer-bullet slice (Roadmap §v0.1a "Phase-slicing note"). It locks the
hardest-to-change contracts early (Guiding Principle #3): journal event v1
format, ID shapes, projection checkpoint, and the write lock.

**Vertical, not horizontal.** Each task adds plumbing, but the phase exits on
commands an agent can actually run: `fid init`, `fid rebuild`,
`fid issue create`, `fid issue get --format json`.

## Scope (in)

- `.fid/` workspace layout, `config.toml`, `pricing.toml`, `.gitignore`.
- `golang-migrate` harness + **minimal** schema: `projection_meta`, `roles`,
  `issues` only. Required pragmas on every connection.
- Journal event v1 model, canonical JSON + hashing, ID generation, upgrader
  chain (mechanism + v1 identity).
- `.fid/write.lock` (lock-dir + owner metadata) and the journal writer
  (single clock read, per-day sequence, append+flush as commit point,
  clock-skew detection).
- Replayer + projection + `projection_meta` checkpoint; clean-read vs catch-up
  vs full rebuild; corruption + precondition conflict handling; `fid rebuild`.
- `fid init` (seeds 7 roles via `project.initialized`; rejects re-init).
- Minimal `fid issue create` / `fid issue get` + JSON `--format` engine seam.
- `--by <role>` with role-existence validation on writes.

## Scope (out — deferred)

- Remaining v0.1 tables, stub tables, triggers, `issues_enriched` view → Phase 2.
- TOON / markdown / YAML output + perf/error-path hardening → Phase 3.
- Full v0.1b issue lifecycle (kinds hierarchy, workflows, transitions, claims,
  ACs, list filters) → v0.1b+.
- Full role×action permission matrix (PRD-ROLE-003) → v0.1d. Phase 1 checks
  role existence only; `fid init` is the bootstrap exception (no `--by`).

## Validation Contract Index

Contract: `../phases/phase-1/qa/validation-contract.md` (VAL-001..VAL-017).

| VAL | Behavior | Owner |
|---|---|---|
| VAL-001 | `fid init` creates `.fid/` workspace + config files | task 5 |
| VAL-002 | `--prefix` sets prefix; default = uppercased dir name | task 5 |
| VAL-003 | One `project.initialized` event seeds 7 roles | task 5 |
| VAL-004 | Re-init rejected, no changes | task 5 |
| VAL-005 | Delete DB + rebuild reconstructs roles + checkpoint | task 5 |
| VAL-006 | Issue created via event, read back as JSON | task 6 |
| VAL-007 | One logical write = one v1 event line | tasks 3, 6 |
| VAL-008 | Replay pure data; rebuild yields identical state | task 4 |
| VAL-009 | Replay routes via upgrader chain (golden v1) | task 4 |
| VAL-010 | Malformed trailing line reported as corruption | task 4 |
| VAL-011 | Precondition mismatch → structured conflict error | task 4 |
| VAL-012 | `projection_meta` drives clean/catch-up/rebuild | task 4 |
| VAL-013 | `.fid/write.lock` serializes writes; holder-named error | task 3 |
| VAL-014 | Pragmas ON on every connection | task 1 |
| VAL-015 | ID shapes + per-day event sequence + collision retry | tasks 2, 3 |
| VAL-016 | Writes require `--by`; unknown roles rejected | task 6 |
| VAL-017 | `internal/` never imports `cmd/fid/`; thin adapters | phase-end |

## Ordered Tasks

1. **SQLite connection, pragmas, minimal migration harness** — `internal/db/`,
   migration `0001` (projection_meta, roles, issues). (VAL-014)
2. **Journal event v1 model, ID gen, canonical JSON, upgrader chain** — pure
   data, no I/O. (VAL-007, VAL-009, VAL-015)
3. **Write lock + journal writer (append+flush)** — `.fid/write.lock`, sequence
   allocation, clock-skew. (VAL-007, VAL-013, VAL-015)
4. **Replayer, projection, checkpoint, catch-up, `fid rebuild`** — corruption +
   precondition errors, freshness logic. (VAL-008..012)
5. **`fid init` bootstrap** — workspace + `project.initialized` + 7 roles;
   reject re-init. (VAL-001..005)
6. **Minimal `fid issue create` / `get` + JSON format engine** — closes the
   tracer bullet; `--by` role check. (VAL-006, 007, 008, 016)

## Exit Criteria

- `fid init` then `fid rebuild` reconstructs the seeded roles from the journal;
  delete-and-rebuild is non-destructive.
- One issue can be created via an event and read back as JSON.
- One logical write = exactly one v1 journal event; replay is pure data.
- Corruption and precondition conflicts stop replay with structured errors.
- Writes require `--by`; unknown roles rejected.
- All VAL-001..VAL-017 PASS with deterministic evidence; QA final review has no
  unresolved blocking findings; human checklist PASS.

## Planning Notes / Carry-In

- Stack locked (Go, Cobra, modernc.org/sqlite, golang-migrate, testify);
  golangci-lint v2 (DEC-003). Sandbox cache redirection in `AGENTS.md`
  (DEC-001/005). Use `gofmt -l cmd internal`, not `gofmt -l .`.
- Cut a new phase branch for Phase 1.
- Carry-Forward Register: no open items at Phase 1 start.
- Roadmap source: `../../specs/10_IMPLEMENTATION_ROADMAP.md` v0.1a +
  slicing note. Re-slice freely if dependency analysis warrants, but every
  phase must ship an end-to-end verifiable increment.

## Next Checkpoint

QA Lead writes `../phases/phase-1/qa/qa-plan.md` mapping every VAL to
deterministic verification; Tech Lead presents contract + tasks to the human;
Developer execution begins only after the human approves.
