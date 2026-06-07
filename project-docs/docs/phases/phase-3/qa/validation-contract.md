# Phase 3 Validation Contract - Output Formats + Foundation Hardening

> Phase: 3 - Output formats (TOON / markdown / YAML) + foundation hardening
> Author: Tech Lead (draft)
> Last refined: <pending QA Lead refinement>
>
> Lives at `project-docs/docs/phases/phase-3/qa/validation-contract.md`.

## Purpose

This file is the finite set of validation assertions that must be true when
Phase 3 is complete. The validation contract specifies **what must be proven**;
the QA plan (authored next by the QA Lead) is the strategy for **how to prove
it**.

Phase 3 completes the `--format` output engine across the one shipped read
command (`fid issue get`) — JSON shipped in Phase 1; this phase adds TOON,
markdown, and YAML — and hardens the persistence foundation: rebuild performance
at single-project scale and conflict/corruption error fidelity. The binding
rule is the roadmap's: each format renders **coherent output through one code
path**, and the performance/error checks are run on a command an agent can
observe.

If a required validation assertion cannot be satisfied, stop and escalate
instead of weakening the assertion or silently adapting the spec.

## Numbering

- IDs use `VAL-NNN`, scoped to this phase, sequential from `VAL-001`.
- Retired or superseded assertions stay in the table with status and rationale;
  do not reuse IDs.

## Coverage-Evidence Carry-In (from Phase 1 DEC-007, reaffirmed Phase 2)

Where a `VAL-NNN` behavior lives in a COV-001 domain package
(`internal/issues`, `internal/journal`, `internal/policy`; thresholds in
`08_QUALITY_PLAN.md` §4.3), its evidence method **names package-level tests in
that package**. The hardening track lives in `internal/journal` (COV-001) and
`cmd/fid` (COV-002), so those VALs name package-level + CLI tests. The format
engine lives in `internal/format`, which is **not** a COV-001 package — its
tests do not gate COV-001 but must carry deterministic evidence (parity tests,
golden fixtures, module-graph checks). Hold COV-001 ≥ 80% and COV-002
(`cmd/fid`) ≥ 70%.

## Assertions

| ID | Source requirement | Behavior or invariant | Evidence method | Validator | Owning task/window | Pass criterion | Fail criterion | Status | Out-of-scope notes |
|---|---|---|---|---|---|---|---|---|---|
| VAL-001 | PRD-OUT-002, PRD-OUT-004; TRD §5 ("Parity tests against JSON output"); 05_DESIGN_PLAN §5.1; 06_DESIGN_PLAN_CLI §2 (`fid issue get`) | `fid issue get --format toon` renders the issue in TOON carrying the same field names and values as the JSON output (structural parity), for both an empty-arrays issue and a populated issue | Package-level `internal/format` parity test: marshal one value to JSON and to TOON; assert every JSON key/scalar is represented in the TOON output with equal value (decode round-trip or structural cross-check) + golden TOON fixture; `cmd/fid` integration smoke | QA Lead | task-1 | TOON output of the canonical empty-arrays issue **and** a synthetic populated issue (labels + nested `acceptance_criteria`) carries the same keys/values as JSON; nested arrays are preserved, not flattened or dropped | TOON drops, renames, flattens, or reorders fields away from the JSON shape, or nested arrays are lost | Draft | TOON is a minimal subset (TRD §11 open question sanctions "start minimal"); full TOON-spec coverage is out of scope. JSON remains canonical. |
| VAL-002 | TRD-ARCH-012 | The TOON marshaller is implemented in-tree with no third-party TOON dependency | `go list -m all` / `go.mod` inspection shows no TOON dependency; `go build -o /tmp/fid ./cmd/fid` succeeds | QA Lead | task-1 | No third-party TOON package appears in the module graph and the binary builds with TOON working | A third-party TOON dependency is added to `go.mod` | Draft | Single negative supply-chain invariant; not a both-states conditional. |
| VAL-003 | PRD-OUT-003, PRD-OUT-004; 06_DESIGN_PLAN_CLI §21 (markdown example); 05_DESIGN_PLAN §8 parity matrix | `fid issue get --format markdown` renders the issue per the §21 layout (title heading, kind/status/priority line, epic/owner/sequence/spec lines, Acceptance Criteria with `- [x]`/`- [ ]`, Labels, Links, footer), for both a populated issue and an issue with no ACs/labels/links | Golden markdown fixture + `internal/format` renderer test + `cmd/fid` integration; both states | QA Lead | task-2 | Populated issue renders all §21 sections with correct checkbox/status mapping; an issue with no ACs/labels/links renders gracefully (empty or omitted sections, no panic, no stray headers) | Output diverges from the §21 structure, mis-maps verification state to checkboxes, or panics/garbles on the empty state | Draft | TUI board (v0.3) and markdown *export* (v0.4) are out of scope; this is CLI `--format markdown` only. |
| VAL-004 | PRD-OUT-003a (Could, v0.1), PRD-OUT-004; 05_DESIGN_PLAN §5.1; DEC-012 | `fid issue get --format yaml` renders the issue as YAML whose decoded structure equals the JSON output (parity), for both empty-arrays and populated issues | Package-level `internal/format` parity test: decode the YAML output (via `gopkg.in/yaml.v3`) and deep-equal it to the JSON-decoded value; `cmd/fid` integration smoke | QA Lead | task-3 | Decoded YAML equals the decoded JSON for both states | Decoded YAML diverges from JSON, or YAML fails to parse | Draft | YAML approach resolved (DEC-012): adopt `gopkg.in/yaml.v3` (MIT) for emit + parity decode. Dependency-budget rationale recorded in DEC-012 + the task-3 completion report. |
| VAL-005 | PRD-OUT-004; 05_DESIGN_PLAN §5.1 | Every read command accepts `--format` and dispatches to the requested renderer; an unrecognised format value is rejected with an actionable error rather than silently falling back to JSON | `cmd/fid` table test over `json`/`toon`/`markdown`/`yaml` and an unknown value (e.g. `xml`); both qualifying and non-qualifying states | QA Lead | tasks 1-3 | Each supported format (`json`/`toon`/`markdown`/`yaml`) produces its own representation; `--format xml` exits non-zero with an `error:` naming the unsupported format and the known formats | A supported format errors, or an unknown format silently produces JSON / exits 0 | Draft | Supported set is the full four formats (YAML adopted, DEC-012). |
| VAL-006 | TRD-PERF-002 (<5s at 10k events); roadmap Phase 3 ("sub-second rebuild at single-project scale"); PRD-OPS-004 | `fid rebuild` over a seeded journal completes within the TRD-PERF-002 budget at 10k events and sub-second at single-project scale (~1k events) | `go test -bench` benchmark in `internal/journal` seeding a journal of valid events; record elapsed for ~1k and 10k | QA Lead | task-4 | Benchmark output shows the 10k-event rebuild < 5s and the ~1k-event rebuild < 1s, with the run recorded as evidence | Either threshold is exceeded, or no benchmark evidence is produced | Draft | Absolute wall-clock is environment-sensitive; the deterministic correctness anchor is VAL-007. The 100k/TRD-PERF-003 tier is out of scope this phase. |
| VAL-007 | TRD-ARCH-001; PRD-OPS-003/004; (extends Phase-2 VAL-016 to scale) | `fid rebuild` over the single-project-scale journal is deterministic and non-destructive: event count matches the journal, delete-DB + rebuild yields an identical projection, and two successive rebuilds are identical | Package-level `internal/journal` test over the scale fixture: assert replayed event count, structural/row equality pre/post delete-DB, and rebuild-vs-rebuild equality | QA Lead | task-4 | Rebuild reconstructs all events; projection state is identical before/after delete-DB and across two rebuilds | Event count diverges, projection differs after rebuild, or two rebuilds disagree | Draft | Phase 2 VAL-016 proved determinism on a small all-tables journal; this proves it at scale (the deterministic companion to VAL-006). |
| VAL-008 | PRD-OPS-003/004; TRD-ARCH-002; roadmap v0.1a Journal replayer ("malformed trailing line = corruption, reported not skipped") | A malformed trailing journal line is reported as a structured corruption error (file:line, parse cause) — never silently skipped — on **both** the explicit `fid rebuild` path and the automatic read catch-up path (`fid issue get`); a well-formed journal replays cleanly | `internal/journal` test (replay returns `CorruptJournalError`) + `cmd/fid` tests asserting `fid rebuild` and `fid issue get` both render the structured corruption message; both states (corrupt vs clean) | QA Lead | task-5 | Malformed trailing line → both commands stop and emit `error: malformed journal entry at <file>:<line>` with the parse cause and "manual repair" guidance; a clean journal replays and reads succeed | A malformed line is silently skipped, the read path serves stale/partial data, or a clean journal is falsely rejected | Draft | CF-003 (policy bootstrap corrupt-tail tolerance) is a separate code path (`internal/policy` role discovery), unchanged here; this VAL covers the replay/rebuild/read-catch-up path only. |
| VAL-009 | PRD-OPS-012; TRD-ARCH-002/003 | A failed effect precondition during replay stops before serving reads and reports a structured logical conflict naming the event, entity (table + pk), expected, and actual — on **both** `fid rebuild` and the read catch-up path; a journal whose preconditions all hold replays cleanly | `internal/journal` test (replay returns `PreconditionConflictError` with fields) + `cmd/fid` tests asserting both `fid rebuild` and a read command render the structured conflict; both states | QA Lead | task-5 | Failed precondition → both paths stop and emit a structured conflict naming event id/type, entity, expected, actual; a satisfiable journal replays and reads succeed | The conflict is rendered as a generic/opaque error, the read path serves data despite the conflict, or a valid journal is falsely rejected | Draft | The replayer already returns `PreconditionConflictError` (tested in `internal/journal`); the new work is structured CLI rendering on the rebuild **and** read-catch-up paths. No automatic reconciliation in v0.1. |
| VAL-010 | TRD-ARCH-005; (standing invariant, re-asserted from Phase-2 VAL-017) | The adapter boundary holds: no `internal/...` package imports `cmd/fid`; new rendering logic lives in `internal/format` and `cmd/fid` stays a thin adapter | `go list` import-graph check (no `cmd/fid` in `internal/...` imports) + a source scan that rendering/parsing logic lives in `internal/format`, not `cmd/fid` | QA Lead | phase-end | Import graph shows no `internal/...`→`cmd/fid` edge; format rendering lives in `internal/format` | Any `internal/...` imports `cmd/fid`, or render/parse business logic leaks into `cmd/fid` | Draft | Standing architectural invariant re-asserted each phase. |

## Decision Gates — RESOLVED (human-approved 2026-06-07)

All three planning decisions are resolved; execution may proceed.

1. **YAML approach (gated VAL-004 / task-3): RESOLVED → option (b), DEC-012.**
   Adopt `gopkg.in/yaml.v3` (MIT) for the YAML renderer and YAML parity decode.
   The in-tree-emitter option was rejected because the TOON in-tree mandate
   (TRD-ARCH-012, driven by TOON's instability) does not transfer to mature YAML,
   and VAL-004 parity needs a YAML decoder in tests regardless. Dependency-budget
   rationale (need/alternatives/license/reuse) is in DEC-012.
2. **Markdown rendering architecture (task-2): RESOLVED → DEC-013.** Typed
   renderer seam in `internal/format` (a markdown renderer taking the typed
   view); JSON/TOON/YAML stay on the generic JSON→structure pipeline. `cmd/fid`
   passes the typed value and only selects/writes.
3. **Generic render pipeline (tasks 1, 3): RESOLVED → DEC-013.** TOON and YAML
   render a generic structure obtained by `value → json.Marshal → json.Unmarshal
   into any`, so the formats inherit the JSON shape and the engine generalises to
   future read commands without per-type code.

## Coverage Rules Check (Tech Lead draft)

- Every in-scope roadmap/PRD/TRD requirement is referenced by at least one
  `VAL-NNN`: PRD-OUT-001 (JSON) shipped Phase 1; PRD-OUT-002 → VAL-001/002,
  PRD-OUT-003 → VAL-003, PRD-OUT-003a → VAL-004, PRD-OUT-004 → VAL-005;
  TRD-ARCH-012 → VAL-002; TRD-PERF-002 → VAL-006; PRD-OPS-003/004 →
  VAL-006/007/008; PRD-OPS-012 → VAL-009; TRD-ARCH-002/003 → VAL-008/009;
  TRD-ARCH-005 → VAL-010: YES.
- Every `VAL-NNN` is owned by a task or `phase-end`: YES (tasks 1-5 + phase-end).
- Conditional/binary assertions name both states: YES (VAL-001, VAL-003,
  VAL-004, VAL-005, VAL-008, VAL-009).
- Parity/"matches reference" assertions carry deterministic conformance proof:
  VAL-001 (JSON↔TOON structural parity + golden), VAL-003 (golden markdown vs
  §21), VAL-004 (YAML decode deep-equal to JSON).
- COV-001 carry-in applied: VAL-007/008/009 name `internal/journal`
  package-level tests; `internal/format` VALs carry deterministic non-COV-001
  evidence.
- QA Lead still to: turn each evidence method into deterministic QA-plan steps;
  confirm the perf benchmark threshold is recorded as evidence (VAL-006);
  confirm corruption/conflict tests assert the structured message on both the
  rebuild and read-catch-up paths.

## Outcomes

Phase-end final review records every `VAL-NNN` as `PASS`, `FAIL`, `BLOCKED`, or
`CARRY-FORWARD`. `PASS` without cited evidence is invalid.

## Revision Log

| Date | Assertion | Change | Reason | Decision reference |
|---|---|---|---|---|
| 2026-06-07 | VAL-001..010 | created (Tech Lead draft) | Phase 3 planning | N/A |
| 2026-06-07 | VAL-004, VAL-005 | refined: YAML approach resolved to `gopkg.in/yaml.v3`; supported set is the full four formats | Decision gates resolved (human-approved) | DEC-012, DEC-013 |
