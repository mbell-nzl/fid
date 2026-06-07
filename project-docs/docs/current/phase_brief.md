# Phase Brief - Phase 3: Output Formats + Foundation Hardening

> **Purpose.** Default execution brief for the current phase. Loaded by every
> agent session. Keep under 150 lines.
>
> **State:** PLANNED — validation contract, ordered tasks, and the three planning
> decisions are human-approved (2026-06-07; DEC-012, DEC-013). Next: QA Lead
> authors `../phases/phase-3/qa/qa-plan.md`, then Developer execution begins on a
> new `phase-3/...` branch. Phase 2 is CLOSED (`../phases/phase-2/`).

## Phase Goal

Complete the `--format` output engine across the one shipped read command
(`fid issue get`) — TOON, markdown, YAML alongside Phase-1 JSON — and harden the
persistence foundation (rebuild performance + conflict/corruption error
fidelity). Roadmap "Phase 3 — Output formats + foundation hardening";
PRD-OUT-002/003/003a/004, TRD-ARCH-012; PRD-OPS-003/004/012, TRD-PERF-002,
TRD-ARCH-002/003.

**Headline increment:** `fid issue get` returns coherent output in all four
formats from one code path, and deterministic perf + error-fidelity checks pass
on rebuild and read catch-up.

## Validation Contract Index

Full contract: `../phases/phase-3/qa/validation-contract.md` (VAL-001..010).

- VAL-001 — `--format toon` is structurally at parity with JSON (empty + nested). [task-1]
- VAL-002 — TOON marshaller is in-tree; no third-party TOON dependency. [task-1]
- VAL-003 — `--format markdown` renders per `06_DESIGN_PLAN_CLI` §21 (both states). [task-2]
- VAL-004 — `--format yaml` decodes equal to JSON (both states); uses `gopkg.in/yaml.v3` (DEC-012). [task-3]
- VAL-005 — every read command accepts `--format`; unknown format rejected, no silent fallback. [tasks 1-3]
- VAL-006 — `fid rebuild` within TRD-PERF-002 (10k<5s) and sub-second at ~1k (benchmark). [task-4]
- VAL-007 — rebuild at scale is deterministic + non-destructive. [task-4]
- VAL-008 — malformed trailing line reported as corruption on rebuild **and** read catch-up. [task-5]
- VAL-009 — failed precondition → structured logical conflict on both paths; stop before reads. [task-5]
- VAL-010 — adapter boundary holds; render logic in `internal/format`. [phase-end]

## Ordered Tasks

1. **task-1 — Generic render pipeline + in-tree TOON marshaller.** JSON→`any`
   structure seam; in-tree TOON; wire `fid issue get --format toon`. (VAL-001/002/005)
2. **task-2 — Markdown renderer for `fid issue get`.** Typed renderer seam in
   `internal/format` per §21. (VAL-003/005)
3. **task-3 — YAML output.** YAML over the generic pipeline via `gopkg.in/yaml.v3`
   (DEC-012). (VAL-004/005)
4. **task-4 — Rebuild performance verification at scale.** Seeded-journal
   benchmark (~1k + 10k) + deterministic scale-rebuild test. (VAL-006/007)
5. **task-5 — Conflict/corruption error fidelity.** Structured corruption +
   precondition-conflict rendering on rebuild and read catch-up. (VAL-008/009)

Tasks 1→2→3 are sequential (each plugs into the format seam). Tasks 4 and 5 are
independent of the format track and of each other.

## Planning Decisions — RESOLVED (human-approved 2026-06-07)

See `../phases/phase-3/decisions.md`.

1. **YAML approach (DEC-012):** adopt `gopkg.in/yaml.v3` (MIT) for the YAML
   renderer + parity decode. The in-tree-emitter option was rejected (the TOON
   in-tree mandate is driven by TOON's instability and does not transfer to
   mature YAML; VAL-004 parity needs a YAML decoder in tests regardless).
   Dependency-budget rationale in DEC-012.
2. **Markdown rendering architecture (DEC-013):** typed renderer seam in
   `internal/format`; JSON/TOON/YAML stay on the generic pipeline.
3. **Generic render pipeline (DEC-013):** TOON/YAML render a generic structure
   from value→JSON→`any`, so the engine generalises to future read commands.

## Carry-In Notes

- **COV carry-in (DEC-007, reaffirmed Phase 2):** hardening VALs (007/008/009)
  live in `internal/journal` (COV-001 ≥ 80%) — name package-level tests.
  `internal/format` is **not** COV-001 but needs deterministic evidence (parity
  tests, golden fixtures, module-graph check). Hold COV-002 `cmd/fid` ≥ 70%.
- **Carry-Forward Register (still open):** CF-003 (policy bootstrap corrupt-tail
  tolerance) — task-5 hardens the *replay/rebuild/read* corruption path only; the
  `internal/policy` bootstrap path is a separate code path, unchanged. Re-defer
  CF-003 (target v0.1d). CF-004 (sequence-vs-scheduling trigger) and CF-005
  (activity_log projection guard / `activity` op) are not exercised this phase;
  re-defer.
- **No `integration` build tag** is in use; integration-style CLI tests run under
  `go test -race ./...`. The phase-level gate adds `go test -bench . -benchmem`
  (required by task-4) and `go build`.

## Current Code State (planning facts)

- `internal/format/format.go` ships JSON; TOON/markdown/YAML return
  `UnsupportedFormatError`. `fid issue get` is the only read command with
  `--format`.
- Replayer returns `CorruptJournalError` (rendered by `fid rebuild`) and
  `PreconditionConflictError` (NOT yet rendered by `fid rebuild`; read path
  renders neither). No benchmarks exist yet.

## Next Checkpoint

Contract, tasks, and decisions are approved. QA Lead authors
`../phases/phase-3/qa/qa-plan.md` (mapping every VAL-001..010 to deterministic
`(qa-lead)`/mid-phase-review/`(human)` steps; confirm the VAL-006 benchmark
threshold and the VAL-008/009 both-paths rendering). Then Developer execution
begins on a new `phase-3/...` branch, starting at task-1.
