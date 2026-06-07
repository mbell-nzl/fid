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
