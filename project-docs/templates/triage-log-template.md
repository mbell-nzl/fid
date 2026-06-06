# Phase <N> - TL-Lite Triage Log

> **Purpose.** Append-only log of mid-phase triage events handled by TL-Lite.
> Phase Runner reads TL-Lite's `VERDICT:` line to decide whether to resume or
> stop. Tech Lead audits this log at phase wrap-up.
>
> **Numbering.** TRG-NNN numbering is per-phase and resets each phase. Newest
> entry at top.

| TRG-NNN | Date | Trigger | Source | Verdict | ESC-NNN (if BLOCK) | Summary |
|---|---|---|---|---|---|---|
| TRG-001 | YYYY-MM-DD | <one line> | <source> | RESUME / BLOCK | <ESC-NNN or -> | <2 lines max> |

---

## TRG-001 - <short title>

- **Trigger:** <verbatim trigger line from Phase Runner prompt>
- **Source:** <task report / review / ESC-NNN>
- **Verdict:** RESUME | BLOCK
- **Reasoning:** <2-4 lines>
- **Local decisions made:** <bullets or "none">
- **References:** <file:line, doc section, DEC-NNN>
- **If BLOCK:** ESC-<NNN> raised in `tl_escalations_current.md`; recommended
  next action: <action>

<!--
Phase wrap-up checklist:

- [ ] Walked every TRG-NNN above.
- [ ] Confirmed each RESUME verdict was correct in hindsight.
- [ ] Flagged any TRG that should have escalated.
- [ ] Confirmed autonomous TL-Lite resume in effect (DEC-002; no warm-up tally).
-->

