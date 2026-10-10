# QA-31 / VAL-048 glossary comparison

Working directory: /home/mark/projects/fid. Evidence is the current file content quoted below and exact closed-field membership checks. No glossary or contract was edited.

## Quoted glossary entries

> **Project-owned prose**
> Definition: Requirements, architecture, validation strategy, decisions, and
> other durable Human-authored intent whose canonical content remains in
> project files and Git. fid references this content through source metadata
> such as labels, paths and digests, without reproducing referenced file content
> in Human views.
> Use when: distinguishing referenced project-file content from fid's recorded
> operational facts and reasoning.
> Avoid: “fid document” or “backend prose,” or treating recorded reasoning as
> referenced file content.
>
> **Recorded reasoning**
> Definition: Actor-authored text stored in a fid judgment record: Concern
> `summary` and `reasoning`, Concern disposition `reason`, Finding `summary`
> and every FindingDisposition `reason`, Deviation `expected`, `observed` and
> `reasoning`, and Human Outcome `reasoning`. It belongs to the fid record and
> may be shown in Human views; it does not establish deterministic proof.
> Use when: explaining the recorded basis for a judgment, distinct from the
> content of a referenced project file (DEC-272; Phase 24 contract section 10).
> Avoid: calling this text project-owned prose, or treating it as Evidence.

## Quoted accepted interface section 10 A7

> - **A7 recorded reasoning.** Human views show recorded reasoning. The closed
>   list:
>   - Concern `summary` and `reasoning`;
>   - Concern disposition `reason`;
>   - Finding `summary` and every FindingDisposition `reason`;
>   - Deviation `expected`, `observed` and `reasoning`;
>   - Human Outcome `reasoning`.
>
>   Lists show one line: whitespace collapsed, cut to 80 Unicode scalars with
>   `…`. Detail views show full text. Control characters are always escaped
>   (PRD-OPS-007).
>   - Phase 23 section 7's "never reproduces project prose" now means never
>     reproducing the **content of a referenced project file**.
>   - The `PRIVATE_` leak checks are replaced by a guard that seeds referenced
>     files with sentinel content and proves no Human view prints it.
>   - GLOSSARY distinguishes project-owned prose from recorded reasoning.

## Item-by-item comparison

| Item | Result | Basis |
|---|---|---|
| Concern summary/reasoning | PASS | Each named entity/field string occurs in both quoted definitions. |
| Concern disposition reason | PASS | Each named entity/field string occurs in both quoted definitions. |
| Finding summary/every disposition reason | PASS | Each named entity/field string occurs in both quoted definitions. |
| Deviation expected/observed/reasoning | PASS | Each named entity/field string occurs in both quoted definitions. |
| Human Outcome reasoning | PASS | Each named entity/field string occurs in both quoted definitions. |
| Referenced-file prose | PASS | Glossary excludes reproducing referenced file content; section 10 excludes the content of a referenced project file. Metadata and recorded reasoning remain distinct. |
| Proof meaning | PASS | Recorded reasoning belongs to a fid record and does not establish deterministic proof; the accepted section identifies these fields as recorded reasoning. |

Observed helper outcome: PASS, all five closed field groups and referenced-file meanings match. See the final review execution receipt for actual exit_code: 0.
