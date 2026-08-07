# TL Escalations - Current

> **Purpose.** Open items that exceed local-decision authority and need full
> Tech Lead or human input.
>
> **Scope.** Open items only. Resolved phase items move to the phase archive.
>
> **Size target:** under 80 lines.

## Open Escalations

**None.** No phase is active.

## Recent Archive Notes

- **Phase 19** — **eight escalations, ESC-19-001..008, all resolved, zero
  product defects.** Every one landed in verification apparatus or evidence
  records: four (003/004/005/007) in one observable's shell ceremony for
  VAL-059, one (006) in its fallout — a false gate attestation disproven by a
  lint regression the same task introduced — one (001) an evidence-token
  defect, one (002) a test-only import cycle, and one (008) the final review's
  five findings, none in product code. Resolved by DEC-094..DEC-106.
  All eight carried real, unique, phase-scoped IDs (DEC-084); the ID-collision
  defect that affected Phases 17-18 did not recur. Three of eight triage records
  required a Tech Lead backfill (CF-023). Archived to
  `../phases/phase-19/tl-escalations.md`.
- **Phase 18** — ESC-001 (final-review BLOCK, `workflow.Load` deadlock) resolved
  by DEC-072; five other triages adjudicated as DEC-067..071 without an
  `ESC-NNN`. In five of six cases TL-Lite reported an escalation it never minted
  and wrote no triage log — fixed by DEC-085, tracked as **CF-023**. Archived to
  `../phases/phase-18/tl-escalations.md`.
- **Phase 17** — five escalations, all resolved; two retrospective items became
  **CF-023**. Archived to `../phases/phase-17/tl-escalations.md`.

## Standing Escalation Triggers

TL-Lite escalates instead of deciding locally when an issue:

1. Changes a public interface, CLI contract, output shape, journal event schema,
   data model, persistence contract, dependency, tool, quality gate, or workflow.
2. Contradicts `decisions_active.md`, `decisions_principles.md`, `AGENTS.md`, or
   a recorded phase decision.
3. Affects more than one task/subsystem, exposes product-intent ambiguity, or
   changes testing strategy / coverage policy.
4. Resumes context-exhausted or partial work spanning schema/data model,
   dependencies, public interfaces, canonical docs, or multiple Developer
   sessions.

QA findings escalate here if they have architecture implications, ambiguous
acceptance criteria, unclear security/privacy severity, contradict a standing
decision, span tasks/subsystems beyond the review window, plausibly indicate a
material spec deviation, or expose a quality-gate anomaly not anticipated by the
QA plan.
