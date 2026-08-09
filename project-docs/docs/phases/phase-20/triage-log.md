# Phase 20 - TL-Lite Triage Log

> **Purpose.** Append-only log of mid-phase triage events handled by TL-Lite.
> Phase Runner reads TL-Lite's `VERDICT:` line to decide whether to resume or
> stop. Tech Lead audits this log at phase wrap-up.
>
> **Numbering.** Identifiers are phase-scoped: `TRG-20-NNN` and `ESC-20-NNN`.
>
> **This file is the allocation source for both.** Read this file, take the
> highest `NNN` already recorded in that series, plus one; start at `001` when
> empty. Never recall a number from memory, never renumber a published entry,
> and never emit a bare placeholder.
>
> **Write the entry before reporting the verdict.** Newest entry at top.

| TRG-20-NNN | Date | Trigger | Source | Verdict | ESC-20-NNN (if BLOCK) | Summary |
|---|---|---|---|---|---|---|
| _(none yet)_ | | | | | | |

---

<!--
Phase wrap-up checklist:

- [ ] Walked every TRG-20-NNN above.
- [ ] Confirmed each RESUME verdict was correct in hindsight.
- [ ] Flagged any TRG that should have escalated.
- [ ] Confirmed every entry carries a real, unique, non-placeholder ID, and
      that every BLOCK names the ESC it raised.
- [ ] Confirmed autonomous TL-Lite resume in effect (DEC-002; no warm-up tally).
-->
