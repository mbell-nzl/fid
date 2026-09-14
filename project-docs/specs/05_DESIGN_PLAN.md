# INTERACTION DESIGN - fid

## 1. Product Experience

fid is a quiet operational CLI. Its primary interactions are:

- a Human maintaining and inspecting the Issue inventory;
- a deterministic kit/harness invoking versioned JSON operations;
- an operator investigating Attention, evidence, recovery, or integrity; and
- a Human reviewing a purpose-built validation brief.

The CLI should make the normal action obvious, return stable typed reasons when
an action is unavailable, and avoid exposing storage or workflow-procedure
detail.

## 2. Surface Rules

1. **Nouns follow the durable model.** Commands use Issue, Envelope, Claim,
   Run, gate, Evidence, Review, Finding, Stop, Attention, Outcome, activation,
   export, and projection.
2. **No role choreography.** There are no Developer, QA, TL, triage, handoff,
   or completion-report command families.
3. **One machine format.** Explicit machine mode emits versioned JSON only.
4. **Human output answers one question.** Default terminal text is concise and
   never an alternate machine schema.
5. **Purpose-built rendering only.** Markdown/HTML may render a named Human
   projection such as a validation brief; it is not a universal `--format`.
6. **Typed unavailability.** “Not ready,” missing authority, stale evidence,
   mutation conflict, and activation conflict are distinct reasons.
7. **No invisible fallback.** A missing capability or stronger evidence
   requirement stops; the CLI does not call an older composite.
8. **Stable dimensions of meaning.** Existence, validity/freshness, and
   Artifact availability are shown separately.

## 3. Human And Machine Modes

### Human mode

- Default for direct operator commands.
- Uses short headings, aligned tables only where comparison helps, and stable
  IDs suitable for a follow-up command.
- Mutations print one confirmation line naming the entity and resulting state.
- Long logs stream directly during gate execution and are not repeated in the
  final summary.

### Machine mode

- Selected explicitly by `--json` on operator commands; backend-port command
  surfaces may be JSON-only.
- Success JSON goes to stdout and ends with one newline.
- Typed error JSON goes to stderr with non-zero process status.
- No progress text, ANSI control sequences, or mixed Human prose appears on
  the JSON channel.
- Every operation envelope names its contract version and operation.

## 4. Information Hierarchy

### Issue view

Show identity, kind, current workflow state/category, parent/dependencies,
priority/sequence, acceptance criteria, labels/links, current Envelope
selection, active Claim/Run, open Concern/Finding/Stop/Attention, and typed
readiness reasons. History is fetched separately.

### Envelope view

Show active revision, outcome, selected Issues and derived progress, authority,
Obligations and their three evidence dimensions, budgets, Stops, Reviews,
Attention, repository identity, and allowed Next Actions. Do not reproduce
project source prose.

### Receipt/evidence view

Show exact registration/definition identity, subject/Run/target identity,
producer/trust, terminal result, freshness, Artifact availability, bounded
excerpts, and stable references to retained material.

### Human Validation Brief

Show delivered outcome, selected Issues, material project Decision references,
evidence and review roll-up, Deviations/unresolved risk, trajectory economics,
and only the remaining Human checks or approvals. Omit routine Run narration.

## 5. Live Gate Execution

- Print the resolved gate identity before execution in Human mode.
- Stream stdout/stderr with clear framing while preserving original order as
  far as the runtime permits.
- On termination, show terminal status, duration, Receipt ID, repository
  identity, truncation/redaction state, and Artifact reference if retained.
- A crash or lost terminal result is shown as abandoned/unknown, never pass.
- Phase 21 JSON machine mode is `fid.live.v1` NDJSON. One globally sequenced
  stream carries sanitized `stdout`, `stderr`, typed `control`, and exactly one
  final `terminal` frame; the terminal frame contains the complete canonical
  Response Envelope and no extra JSON follows it.
- Payload is UTF-8 with invalid bytes replaced and control bytes JSON-escaped.
  Redaction precedes framing, excerpting, hashing, retention, errors, and Human
  rendering. A combined output ceiling emits one `output-truncated` control
  frame while fid continues draining and observes process termination.
  Byte counts are exact until that ceiling fires; afterwards the reported
  counts are lower bounds that depend on pipe read boundaries and are not
  reproducible across runs, so no fixture may pin an exact post-truncation
  count (DEC-189).
- Human mode prefixes sanitized output by stream, renders no raw control byte,
  and derives its final summary from the same terminal result.
- The exact stdout/stderr/control/terminal payload union, terminal outcomes,
  exit classification, and Human obligations are fixed by
  `../docs/phases/phase-21/public-interface.md` section 8; this design text does
  not introduce an alternate rendering or frame schema.

## 6. Errors And Stops

Every failure maps to a stable code and category. Important categories:

- invalid input or unknown reference;
- unsupported capability/schema;
- authority or protected-Decision failure;
- state/revision/idempotency conflict;
- mutation isolation or lease failure;
- unsatisfied/stale/unavailable evidence;
- gate terminal failure or abandoned Attempt;
- review rework or budget exhaustion;
- activation conflict; and
- maintenance/recovery conflict.

Human messages state what failed, the subject, and the next legitimate action.
They never advise editing `.fid/` internals or weakening an Obligation.

## 7. Accessibility And Safety

- Color is optional and never the only state signal.
- `NO_COLOR` is respected.
- IDs and reason codes remain copyable plain text.
- Untrusted output is escaped or sanitized before terminal/Markdown/HTML
  rendering.
- Destructive or irreversible operations require explicit flags and present
  the affected scope in Human mode; kit-mode authority is additionally checked
  through the active contract.

## 8. Deferred Experience Surfaces

TUI/dashboard and search are the leading evidence-gated candidates in
ADR-0011. They are not represented by placeholder commands or speculative
dependencies. The first supported CLI must prove whether concise semantic
projections are sufficient.

*Last updated: 2026-08-22 (Phase 21 DEC-157 live-contract amendment)*
