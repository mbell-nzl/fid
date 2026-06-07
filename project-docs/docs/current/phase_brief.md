# Phase Brief - Phase 3: Output Formats + Foundation Hardening

> **Purpose.** Default execution brief for the current phase. Loaded by every
> agent session. Keep under 150 lines.
>
> **State:** NOT YET PLANNED (stub written at Phase 2 wrap, 2026-06-07). No
> validation contract, QA plan, or task specs exist yet. Phase 2 is CLOSED (see
> `../phases/phase-2/`). The next Tech Lead planning session owns the contract
> and ordered task specs; Developer execution begins only after the human
> approves them.

## Candidate Phase Goal (from roadmap slicing note)

Complete the v0.1a output surface and harden the persistence foundation
(Roadmap "Phase 3 — Output formats + foundation hardening"):

- **Output formats:** finish the `--format` engine across all read commands —
  TOON (opt-in), markdown, and YAML — alongside the JSON shipped in Phase 1.
  (`05_DESIGN_PLAN.md`, `06_DESIGN_PLAN_CLI.md`; PRD-OUT-001..004; TRD-ARCH-012.)
- **Foundation hardening:** verify sub-second `fid rebuild` at single-project
  scale and confirm conflict/corruption error fidelity (malformed trailing line
  = reported corruption; failed effect precondition = structured logical
  conflict). (PRD-OPS-003/004/012; TRD-PERF-002; TRD-ARCH-002/003.)

**Headline end-to-end increment (candidate):** a read command (e.g.
`fid issue get`) returns coherent output in all four formats from one code path,
and a deterministic perf/error-fidelity check passes on rebuild.

> Re-slice freely if dependency analysis at phase-start suggests a different cut.
> The binding rule is *every phase ships an end-to-end verifiable increment* and
> exits on a command an agent can actually run and observe.

## Likely Scope (out — confirm at planning)

- v0.1b lifecycle commands (`issue create` kind hierarchy, transitions, claims,
  ACs, list filters) — later v0.1 slices.
- v0.2+ tables (workflows, role_permissions, comments) and beyond.

## Planning Checklist (Phase Start)

1. Read the roadmap Phase 3 candidate + relevant `05`/`06` design and `02` TRD
   perf/error sections.
2. Review the Carry-Forward Register in `decisions_active.md` (CF-003, CF-004,
   CF-005) — none are obviously exercised by output/hardening, but confirm.
3. Create `../phases/phase-3/` directories and draft
   `../phases/phase-3/qa/validation-contract.md` (VAL-NNN, one observable each;
   binary/conditional assertions name both states; parity assertions need
   deterministic conformance evidence — e.g. format round-trip / golden output).
4. Apply the standing COV-001 carry-in (Phase 1 DEC-007, reaffirmed Phase 2):
   when a VAL's behavior lives in a COV-001 domain package
   (`internal/issues`, `internal/journal`, `internal/policy`), name
   **package-level** tests; format rendering likely lives in `internal/format` /
   `internal/export` (not COV-001 packages) but still needs deterministic
   evidence. Hold COV-001 ≥ 80% and COV-002 (`cmd/fid`) ≥ 70%.
5. Break into ordered one-session Developer tasks; each cites covered VAL-NNN,
   lists exact gate commands (including any QA-plan command the Developer must
   run), and lists documentation updates.
6. Refresh this brief with the real contract index + ordered task list, then
   present to the human.

## Carry-In Notes

- DEC-010 (Phase 2): bump `ProjectionSchemaVersion` only for migrations that
  change projected structure, not for trigger-only migrations. Relevant if any
  view/index change ships.
- DEC-011 (Phase 2): `task_completions` kind guard is task-or-subtask; `02_TRD`
  §4.1 reconciled to match §14.1. No further action.
- Carry-Forward Register (still open): CF-003 (role-validation corrupt-tail
  tolerance — note this if Phase 3 hardening touches journal corruption error
  paths), CF-004 (sequence-vs-scheduling trigger), CF-005 (`activity_log`
  projection-only guard / `activity` op).

## Next Checkpoint

Tech Lead drafts the Phase 3 validation contract + ordered tasks and presents
them to the human. On approval, QA Lead authors
`../phases/phase-3/qa/qa-plan.md` and Developer execution begins on a new
`phase-3/...` branch.
