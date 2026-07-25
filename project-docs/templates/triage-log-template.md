# Phase <N> - TL-Lite Triage Log

> **Purpose.** Append-only log of mid-phase triage events handled by TL-Lite.
> Phase Runner reads TL-Lite's `VERDICT:` line to decide whether to resume or
> stop. Tech Lead audits this log at phase wrap-up.
>
> **Numbering.** Identifiers are **phase-scoped**: `TRG-<phase>-<NNN>` and
> `ESC-<phase>-<NNN>` — for example `TRG-19-001`, `ESC-19-001`. The phase
> number comes from the phase brief.
>
> **This file is the allocation source for both.** To allocate the next ID,
> read this file and take the highest `<NNN>` already recorded in that series,
> plus one; if the file is new, start at `001`. Never recall a number from
> memory, never renumber a published entry, and never emit a bare `TRG-NNN` or
> `ESC-NNN` placeholder — a placeholder is not an identifier and leaves the
> event with no durable record.
>
> **Write the entry before reporting the verdict.** This row, not the chat
> reply, is what the Phase Runner confirms and the Tech Lead audits at wrap-up.
> Newest entry at top.

| TRG-<phase>-<NNN> | Date | Trigger | Source | Verdict | ESC-<phase>-<NNN> (if BLOCK) | Summary |
|---|---|---|---|---|---|---|
| TRG-19-001 | YYYY-MM-DD | <one line> | <source> | RESUME / BLOCK | <ESC-19-001 or -> | <2 lines max> |

---

## TRG-<phase>-<NNN> - <short title>

- **Trigger:** <verbatim trigger line from Phase Runner prompt>
- **Source:** <task report / review / ESC-<phase>-<NNN>>
- **Verdict:** RESUME | BLOCK
- **Reasoning:** <2-4 lines>
- **Local decisions made:** <bullets or "none">
- **References:** <file:line, doc section, DEC-NNN>
- **If BLOCK:** ESC-<phase>-<NNN> raised in `tl_escalations_current.md`;
  recommended next action: <action>

<!--
Phase wrap-up checklist:

- [ ] Walked every TRG-<phase>-<NNN> above.
- [ ] Confirmed each RESUME verdict was correct in hindsight.
- [ ] Flagged any TRG that should have escalated.
- [ ] Confirmed every entry carries a real, unique, non-placeholder ID, and
      that every BLOCK names the ESC it raised.
- [ ] Confirmed autonomous TL-Lite resume in effect (DEC-002; no warm-up tally).
-->

