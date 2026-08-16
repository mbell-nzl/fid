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
>
> Under `Recovery directive schema: v2`, every value in the typed
> control-conformance directive is written on one physical line. The Phase
> Runner copies those values byte-for-byte into the run ledger; prose elsewhere
> in the entry does not grant or widen repair authority.

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
- **Verdict:** RESUME | REPAIR | REPAIR / control-conformance | BLOCK
- **Reasoning:** <2-4 lines>
- **Local decisions made:** <bullets or "none">
- **Resolution evidence:** <why cleared; bounded repair evidence; or exact
  protected/exhausted condition>
- **If REPAIR:** owning task <task-M>; category <recovery category>; next action
  <same-task continuation plus exact verification/re-review>; no ESC.
- **Originating task:** <task-M for REPAIR / control-conformance | N/A>
- **Control trigger signature:** <exact stable signature | N/A>
- **Control attempt count:** <exact next mediated attempt | N/A>
- **Control contract:** <binding accepted source and required behavior | N/A>
- **Control mismatch:** <exact reproduced mismatch and termination | N/A>
- **Allowed files:** <exact ordered project-docs/tools/ paths | N/A>
- **Forbidden changes:** <exact protected paths and semantics | N/A>
- **Exact verification:** <exact command bytes and working directory | N/A>
- **Originating report:** <exact repo-relative completion-report path | N/A>
- **Control repair base:** <full runner-supplied pre-triage phase HEAD | N/A>
- **If BLOCK:** category <protected blocker category>;
  ESC-<phase>-<NNN> <raised or referenced> in `tl_escalations_current.md`;
  resolution condition <named Human or external-authority action>.
- **Human attention required:** <no for RESUME/REPAIR | yes for BLOCK>
- **References:** <file:line, doc section, DEC-NNN>

For `REPAIR / control-conformance`, the ten fields from `Originating task`
through `Control repair base` are the authoritative typed directive. The
ordinary `Trigger signature` and `Attempt count` must equal `Control trigger
signature` and `Control attempt count`. TL-Lite copies the runner-supplied full
`Control repair base`; it never selects or verifies a different base.

## Control Directive Amendment For TRG-<phase>-<NNN>

Use this append-only form only when an existing `DEC-<NNN>` resolves the
governing BLOCK into `REPAIR / control-conformance`. Do not edit the original
BLOCK entry and do not allocate a new TRG. TL-Lite or the full Tech Lead
appends the amendment after the decision exists and before RUN-RESUME.

- **Governing TRG:** TRG-<phase>-<NNN>
- **Governing ESC:** ESC-<phase>-<NNN>
- **Resolution authority:** DEC-<NNN>
- **Verdict:** REPAIR / control-conformance
- **Originating task:** <task-M>
- **Control trigger signature:** <exact stable signature retained from STOP>
- **Control attempt count:** <exact next mediated attempt number>
- **Control contract:** <binding accepted source and uniquely required behavior>
- **Control mismatch:** <exact expected/observed reproducer and termination>
- **Allowed files:** <exact ordered project-docs/tools/ paths>
- **Forbidden changes:** <exact protected paths and semantics>
- **Exact verification:** <exact command bytes and working directory>
- **Originating report:** <exact repo-relative completion-report path>
- **Control repair base:** <full phase HEAD supplied and verified by Phase Runner>
- **References:** <existing DEC, task/report, control contract, and evidence>

Every boundary-valid DEC citation in an entry or amendment must exist in the
phase decisions file. One real DEC does not validate another nonexistent
citation, and embedded text such as `NOTDEC-123X` is not a citation.
