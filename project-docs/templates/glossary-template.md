# GLOSSARY - <PROJECT_NAME>

> This document defines project language. It is not an architecture spec,
> product requirements document, implementation guide, or decision log.

## Purpose

<Explain why terminology matters for this project and which artifacts should
use this glossary.>

## Usage Rules

- Define project-specific terms, not general programming vocabulary.
- Prefer one term for one concept. Put rejected aliases in `Avoid`.
- Keep definitions short and implementation-free.
- `Definition`, `Use when`, and `Avoid` are required. `Related` and
  `Examples` are optional — include them when they clarify boundaries
  between nearby terms.
- Use `project-docs/specs/02_TRD.md` for architecture and module boundaries.
- Use phase `DEC-NNN` entries for durable tradeoff decisions; do not create a
  parallel ADR system unless explicitly adopted.
- Flag unresolved ambiguity before implementation depends on the term.

## Terms

### <Preferred Term>

Definition: <one or two sentences>
Use when: <where this term applies>
Avoid: <synonyms or ambiguous terms to avoid>
Related: <nearby terms and distinction> *(optional)*
Examples: <short usage examples> *(optional)*

## Flagged Ambiguities

| Term or phrase | Why it is ambiguous | Owner | Resolution status |
|---|---|---|---|
| <term> | <ambiguity> | <Human / Tech Lead / role> | <open/resolved with DEC-NNN> |

*Last updated: <YYYY-MM-DD>*
