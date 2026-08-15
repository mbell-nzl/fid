# Phase <N> - TL-Lite Triage Log

> **Purpose.** Append-only log of mid-phase triage events handled by TL-Lite or
> the full Tech Lead. Phase Runner reads the `VERDICT:` line to decide whether
> to advance, dispatch a same-task repair, or stop. Tech Lead audits this log at
> wrap-up.
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
> Only TL-Lite or the full Tech Lead may allocate an `ESC`. Developer and QA
> outputs are recovery or blocker candidates; they neither mint nor reserve an
> escalation ID. `REPAIR` creates no ESC and requests no Human attention.
>
> **Write the entry before reporting the verdict.** This entry, not the chat
> reply, is what the Phase Runner confirms and the Tech Lead audits at wrap-up.
> Append new entries at the end in wall-clock order. Never insert, edit, delete,
> or renumber an earlier entry; append a labelled correction instead.

<!--
Phase wrap-up checklist:

- [ ] Walked every TRG-<phase>-<NNN> in this file.
- [ ] Confirmed each RESUME had resolution evidence and every REPAIR completed
      verification plus any required independent re-review.
- [ ] Flagged any TRG that should have escalated.
- [ ] Confirmed every entry carries a real, unique, non-placeholder ID, and
      that every BLOCK names the ESC it raised or referenced.
- [ ] Confirmed QA and Developer candidates did not mint ESC IDs.
-->

---

## TRG-<phase>-<NNN> - <short title>

- **Trigger:** <verbatim trigger line from Phase Runner prompt>
- **Trigger signature:** `<task/review>:<category>:<stable source identity>`
- **Source:** <task report / review / session / existing ESC>
- **Candidate:** <RECOVERY_REQUIRED / QA REPAIR / blocker candidate / other>
- **Attempt count:** <total mediated continuations for this signature>
- **Checkpoint:** <commit/tree/report/review state>
- **Verdict:** RESUME | REPAIR | BLOCK
- **Reasoning:** <2-4 lines>
- **Local decisions made:** <bullets or "none">
- **Resolution evidence:** <why cleared; bounded repair evidence; or exact
  protected/exhausted condition>
- **If REPAIR:** owning task <task-M>; category <recovery category>; next action
  <same-task continuation plus exact verification/re-review>; no ESC.
- **If BLOCK:** category <protected blocker category>;
  ESC-<phase>-<NNN> <raised or referenced> in `tl_escalations_current.md`;
  resolution condition <named Human or external-authority action>.
- **Human attention required:** <no for RESUME/REPAIR | yes for BLOCK>
- **References:** <file:line, doc section, DEC-NNN>
