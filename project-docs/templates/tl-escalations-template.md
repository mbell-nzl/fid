# TL Escalations - <Current | Phase N>

> **Purpose.** Unresolved escalations that exceed local-decision authority and
> need full Tech Lead or human input.
>
> **Scope.** Open items only for the current file. Resolved items move to the
> phase archive at wrap-up.
>
> **Size target:** under 80 lines.

## Open Escalations

> Newest at top. Identifiers are **phase-scoped**: `ESC-<phase>-<NNN>`, for
> example `ESC-19-001`. Allocate from the phase triage log
> (`project-docs/docs/phases/phase-<N>/triage-log.md`), which records every ESC
> minted in the phase and survives an escalation being closed and archived out
> of this file. Take the highest `<NNN>` recorded there, plus one; start at
> `001` in a new phase. Never emit a bare `ESC-NNN` placeholder — an escalation
> without an identifier cannot be resumed, audited, or closed.

### ESC-<phase>-<NNN>: <short title>

- **Raised:** <YYYY-MM-DD>
- **By:** <agent / session reference>
- **Trigger:** <one escalation trigger>
- **Context:** <2-4 lines: where it surfaced and what was observed>
- **Question / proposed direction:** <what must be decided>
- **Links:** <task spec, review, code reference, related DEC>
- **Status:** Open

---

(Add new escalations above this divider. Remove this block when none are open.)

## Standing Escalation Triggers

TL-Lite must escalate if the issue:

1. Changes a public interface, CLI command contract, output format, journal
   event schema, or user-visible workflow.
2. Changes the data model or persistence contract, including SQLite schema,
   migrations, journal replay semantics, or storage invariants.
3. Contradicts `decisions_active.md`, `decisions_principles.md`, or `AGENTS.md`.
4. Affects more than one task or subsystem.
5. Introduces a dependency, tool, quality gate, or workflow.
6. Changes testing strategy or coverage policy.
7. Exposes product-intent, acceptance-criteria, or design ambiguity.
8. Resumes a context-exhausted or partial task where remaining or completed
   work spans schema/data model, dependency/platform configuration, public
   interfaces, canonical docs, or multiple Developer sessions.

QA Reviewer must also escalate if the finding:

1. Is a suspected architecture or layer-boundary violation with implications
   beyond the review window.
2. Is a test gap that implies the acceptance criterion itself is ambiguous.
3. Is a security, data-loss, journal-integrity, or privacy concern where
   severity is genuinely unclear.
4. Contradicts a standing decision in `decisions_principles.md` or `AGENTS.md`.
5. Spans tasks or subsystems beyond the review window.
6. Is a plausible but material spec deviation.
7. Is a coverage, lint, generated-artifact, or command-gate anomaly the phase
   QA plan did not anticipate.
