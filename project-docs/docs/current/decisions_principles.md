# Decisions - Principles

> **Purpose.** Standing engineering and workflow rules that apply across
> phases. Loaded on demand when an agent needs to cross-check a convention.
>
> **Promotion rule.** Promote a rule here only if it is enforced across
> multiple phases or cited by `AGENTS.md`.

## Source Of Truth

`AGENTS.md` is the authoritative project grounding file. When this file and
`AGENTS.md` disagree, update this file to match `AGENTS.md` unless the human
explicitly decides otherwise.

## Standing Principles

- Agents communicate through repository artifacts, not chat memory.
- Tests are written before behavior changes.
- Deterministic gates run before subjective review.
- Open carry-forward constraints stay in `decisions_active.md`.
- Closed decisions live in per-phase `decisions.md` archives.
- Specs are living guides: update them when implementation proves them wrong.
- Dependency additions require documented need, alternatives, and license check.
- Deterministic QA PASS claims require reproducible evidence, not code-reading
  alone.
- Human checklist items are only for perception, product-owner judgment,
  dogfooding, tone, and flow coherence.
- Held-out scenarios remain outside normal implementation agent context.
- Phase wrap-up records spec sync verdicts, audits TRG entries, and
  appends workflow notes.

## Candidates For Future Promotion

- <Add recurring phase-local decisions here when they appear more than once.>
