# Phase 3 Decisions

> Purpose: phase-local decision history, spec sync verdicts, triage audit,
> workflow notes, and wrap-up decisions. DEC numbering is global and never
> resets (Phase 2 ended at DEC-011).

## Decisions

### DEC-012: Adopt `gopkg.in/yaml.v3` for `--format yaml`

- **Date:** 2026-06-07
- **Context:** Phase 3 completes the `--format` engine, including YAML
  (PRD-OUT-003a, "Could"/v0.1). Go has no standard-library YAML. Three options
  were weighed at planning: (a) an in-tree minimal YAML emitter (dependency-free,
  by analogy to the in-tree TOON mandate); (b) take the de-facto YAML library;
  (c) defer YAML and carry the assertion forward. The TOON analogy was found not
  to transfer: TRD-ARCH-012 mandates TOON in-tree **because TOON is a young,
  unstable format**; YAML is a mature, stable spec. A correct YAML emitter
  (quoting/escaping, null handling, nested sequences) is a real bug surface, and
  VAL-004's parity evidence ("decode the YAML output and deep-equal it to JSON")
  needs a YAML *decoder* in tests regardless — so an in-tree emitter would force
  carrying the dependency as a test-only dep **and** maintaining a fragile
  emitter (the worst combination).
- **Decision:** Adopt **`gopkg.in/yaml.v3`** as a direct runtime dependency for
  the YAML renderer (emit) and YAML parity tests (decode). Human-approved
  (2026-06-07).
- **Dependency budget (AGENTS.md):**
  - *Need:* render `--format yaml` correctly and prove structural parity with
    JSON; stdlib has no YAML; an in-tree emitter is a non-trivial
    correctness/maintenance burden and still needs a decoder to test parity.
  - *Alternatives:* in-tree emitter (more code + bug surface; still needs a
    test-only decoder); defer YAML (rejected — ship the full four-format surface
    this phase).
  - *License:* MIT — passes the default permissive-license policy.
  - *Reuse:* also reusable for the v0.4 markdown-export YAML frontmatter
    (TRD-ARCH-014), though v0.4 may layer a frontmatter helper on top.
- **Consequences:** `go.mod` gains a direct dependency (its first beyond
  migrate/cobra/sqlite). Record it in the task-3 completion report and surface it
  to the dependency-vulnerability / license gates before release readiness. Spec
  impact assessed at wrap: `02_TRD.md` §9 dependency table likely gains a
  `gopkg.in/yaml.v3` row (flag for spec sync).
- **Refs:** VAL-004, VAL-005 (`qa/validation-contract.md`); `tasks/task-3.md`;
  PRD-OUT-003a, PRD-OUT-004; `02_TRD.md` §9.

### DEC-013: Format engine architecture — generic pipeline + typed markdown seam

- **Date:** 2026-06-07
- **Context:** `internal/format.Write(writer, formatName, value any)` serializes
  JSON generically. Phase 3 adds TOON, YAML, and markdown. TOON and YAML can be
  produced generically from any value; markdown (06_DESIGN_PLAN_CLI §21) is a
  bespoke per-command layout that cannot be derived from a generic field-walk.
  Recording the chosen pattern up front avoids an unrecorded-architecture hard
  stop for the Developer (TRD-ARCH-005 / VAL-010 keep `cmd/fid` thin and
  rendering in `internal/`).
- **Decision:**
  1. **Generic render pipeline for JSON/TOON/YAML.** Render from a generic tree
     obtained by `value → json.Marshal → json.Unmarshal into any`
     (`map[string]any` / `[]any` / scalars). TOON and YAML walk that tree. This
     makes the three formats inherit the JSON field names/null handling (they
     cannot drift — serves VAL-001/VAL-004 parity) and gives future read commands
     TOON/YAML with no per-type code.
  2. **Typed renderer seam for markdown.** Markdown is a typed per-command
     renderer in `internal/format` (e.g. a function taking the typed issue view),
     not a generic path and not built in `cmd/fid`. `cmd/fid` passes the typed
     value and only selects/writes.
- **Consequences:** New internal pattern in `internal/format` (generic pipeline +
  typed renderers side by side). No public-interface change; this is presentation
  layering under TRD-ARCH-005. Likely no durable-spec edit (the design plan already
  states "TOON/YAML mirror JSON"); confirm at spec sync (`verified-no-change`
  expected for the seven specs, pending wrap audit).
- **Refs:** VAL-001, VAL-003, VAL-004, VAL-005, VAL-010
  (`qa/validation-contract.md`); `tasks/task-1.md`, `tasks/task-2.md`,
  `tasks/task-3.md`; TRD-ARCH-005, TRD-ARCH-012; 05_DESIGN_PLAN §5.1.

## QA Recommendation Disposition

QA final review recorded **0 blocking / 0 major / 0 improvement** findings
required for PASS. The three non-blocking improvements came from mid-phase
`review-3.md`; each is dispositioned below.

| # | Source | Recommendation | Disposition | Rationale |
|---|---|---|---|---|
| 1 | review-3 Task 1 | Include a short stdout/stderr excerpt for user-facing output VALs, not just test names. | Cancel (no action) | Documentation preference for future reports; the final review already carried deterministic test + parity + golden evidence for VAL-001/003/005. No phase-gate impact. |
| 2 | review-3 Task 2 / TRG-001 | Centralise the enriched issue read model in `internal/issues` before reusing the markdown-only read shape outside `fid issue get`. | Defer → CF-006 | Read-shape/business-placement debt living in `cmd/fid/issue.go`. Acceptable for one read command, but it brushes TRD-ARCH-005 if a second read command reuses it. Bite the cost at v0.1b when `fid issue get` is enriched (`--section`) and `fid issue list` lands. Added to the Carry-Forward Register. |
| 3 | review-3 Task 3 | Carry the DEC-012 approval note into phase-end dependency evidence so `go list -m all` is not misread as an unapproved add. | Do now | Done via spec sync: `02_TRD.md` §9 now lists `gopkg.in/yaml.v3` with `DEC-012` as its decision ref, so release-readiness dependency/license gates have an authoritative approval record. |

## Validation Contract Audit

- All **VAL-001..VAL-010** final outcomes are **PASS** with cited deterministic
  evidence (`qa/final-review.md` Validation Contract Outcomes).
- No FAIL, BLOCKED, or CARRY-FORWARD assertion to disposition.
- Every binary/conditional assertion was proven in both states: VAL-001 (empty +
  populated/nested TOON parity), VAL-003 (populated §21 + empty omitted
  sections), VAL-004 (empty + populated YAML decode parity), VAL-005 (supported
  set dispatch + `xml` rejection), VAL-008 (corrupt tail + clean journal on
  rebuild **and** read catch-up), VAL-009 (failed precondition + satisfiable
  journal on both paths).
- No validation-contract deviation was required; the contract was satisfied as
  written. The VAL-006 wall-clock thresholds (~1k 348.6 ms < 1s; 10k 3431 ms <
  5s) held with margin and are anchored deterministically by VAL-007.

## Spec Sync Verdict

Record a verdict for every spec at phase wrap-up. `verified-no-change` is a
positive assertion after review, not a default skip.

| Spec | Verdict | Evidence / Reference |
|---|---|---|
| `00_PROJECT_CHARTER.md` | verified-no-change | Phase 3 completed the in-scope v0.1 `--format` surface and hardened replay; no identity/scope/constraint change. |
| `01_PRD.md` | verified-no-change | PRD-OUT-002/003/003a/004 and PRD-OPS-003/004/012 MET as written (final review Requirement Compliance); no requirement reworded. |
| `02_TRD.md` | updated | §9 Dependency Policy table gains a `gopkg.in/yaml.v3` (MIT, DEC-012) row for `--format yaml` (do-now sync). TRD-ARCH-012 in-tree TOON (VAL-002), TRD §5 TOON-parity (VAL-001), TRD-ARCH-005 boundary (VAL-010), TRD-PERF-002 (VAL-006), TRD-ARCH-001/002/003 (VAL-007/008/009) all MET as written. §11 TOON open question stays open by design — Phase 3 shipped the sanctioned minimal subset ("start minimal and extend on demand"). |
| `03_TRD_SCHEMA.md` | N/A | Phase 3 added no schema, migration, trigger, or view. |
| `04_TRD_CONFIG.md` | N/A | Phase 3 added no `.fid/` layout or config surface. |
| `05_DESIGN_PLAN.md` | verified-no-change | §5.1 "TOON/YAML mirror JSON" and the §8 parity matrix already describe the realized behavior; DEC-013's generic-pipeline + typed-markdown seam is presentation layering, no design-language change. |
| `06_DESIGN_PLAN_CLI.md` | verified-no-change | `fid issue get --format markdown` renders the §21 layout as documented; TOON/YAML/`xml`-rejection match the §2 `--format` contract. No CLI spec edit needed. |
| `07_DESIGN_PLAN_FLOWS.md` | verified-no-change | No flow change; corruption/conflict rendering follows the documented journal→projection→rebuild / read-catch-up flow. |
| `08_QUALITY_PLAN.md` | verified-no-change | COV-001 (`internal/journal` 80.6%) and COV-002 (`cmd/fid` 74.3%) held; `internal/format` carried non-COV-001 deterministic evidence (parity tests, golden fixtures, module-graph + import-graph checks) per the DEC-007 carry-in. No gate/process change. |
| `09_AI_WORKFLOW.md` | verified-no-change | Workflow ran as specified (validation contract → QA plan → dev/QA loop → one mid-phase review → final review → human checklist); no workflow change. |
| `10_IMPLEMENTATION_ROADMAP.md` | verified-no-change | Phase 3 delivered the v0.1a "Output formats + foundation hardening" slice per the slicing note; CF-003/004/005 re-deferred, no re-slice required. |

## Triage-Log Audit

- [x] Walked `triage-log.md`: one entry, **TRG-001** (Task 2 markdown read-shape
  duplication note), verdict **RESUME**.
- [x] RESUME confirmed correct in hindsight: the note recorded no spec deviation,
  no dependency change, and PASS on required gates and VAL-003/VAL-005. It is
  technical debt, not a phase stopper, and is now tracked as CF-006 for v0.1b
  review — exactly the "revisit when a later read command needs the shared view"
  condition the triage entry named.
- [x] No TRG should have escalated: TRG-001 matched none of the TL-Lite or QA
  escalation triggers (single task, no public-interface/data-model change, no
  standing-decision conflict).
- [x] Autonomous TL-Lite resume remains permanently ENABLED (DEC-002); no warm-up
  tally to maintain. (Phase 3 uses per-phase TRG numbering, consistent with the
  Phase 1 log convention.)

## Workflow Notes

Answer in 1-3 sentences each. Completed at phase wrap-up.

- **What worked well?** Resolving the two architecture decisions (DEC-012 YAML
  dependency, DEC-013 generic-pipeline + typed-markdown seam) at planning time
  meant the Developer never hit an unrecorded-architecture hard stop, and the
  generic `value→JSON→any` pipeline made TOON and YAML inherit JSON parity for
  free (VAL-001/004 could not drift). Independent slicing of the format track
  (tasks 1-3) from the hardening track (tasks 4-5) kept each session focused.
- **What was inefficient?** Very little. The markdown path had to assemble its
  enriched view inside `cmd/fid` because `issues.Get` does not yet expose all the
  fields — the only read-shape debt of the phase (now CF-006), and it was
  correctly scoped rather than expanded mid-phase.
- **Were tasks the right size?** Yes. Five one-session tasks, each red/green TDD,
  no blockers; the format tasks plugged into a shared seam and the hardening
  tasks were each a single deterministic verification slice.
- **Did the QA review cadence feel right?** Yes. One mid-phase review after the
  three-task format track (review-3, CONTINUE, 0 blocking) plus the phase-end
  final review was right for a five-task phase; tasks 4 and 5 were independent
  and well covered by task reports + final QA Lead execution.
- **Were there token-efficiency opportunities missed?** Minor. The final review
  re-ran the full gate and benchmarks (appropriate for VAL-006 evidence) rather
  than only consolidating task evidence; that re-run is justified for the
  performance assertion and was not wasteful.

### Deep Retrospective Recommendation

**No retro warranted** beyond this note. Phase 3 ran clean: all 10 VALs PASS
first time, 0 blocking/major findings, one RESUME triage that was correctly a
deferral, human checklist 2/2 PASS. The durable lessons — resolve format/engine
architecture decisions at planning time, and keep the generic render pipeline so
parity cannot drift — are captured in DEC-012/DEC-013 and need no separate
retrospective. The one open thread (CF-006 enriched read-model centralisation) is
scheduled, not systemic.
