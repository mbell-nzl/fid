# TL Escalations - <Current | Phase N>

> **Purpose.** Unresolved escalations that exceed local-decision authority and
> need full Tech Lead, Human, or named external-authority action. Ordinary
> incomplete work, red tests, failed gates, report defects, and bounded QA
> rework are recovery candidates, not escalations.
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
>
> Only TL-Lite or the full Tech Lead allocates an ESC. Developer and QA roles
> may emit a blocker candidate with evidence, but must not mint or reserve an
> identifier. Human attention is requested only for a protected category below
> or exhausted bounded recovery.

### ESC-<phase>-<NNN>: <short title>

- **Raised:** <YYYY-MM-DD>
- **Allocated by:** <TL-Lite or full Tech Lead / session reference>
- **Candidate source:** <Developer/QA report, review, Phase Runner, or other>
- **Protected category:** <closed category below>
- **Trigger:** <one escalation trigger>
- **Trigger signature:** `<task/review>:<category>:<stable source identity>`
- **Recovery attempts:** <count and outcome; 0 when recovery is unsafe or
  unauthorized>
- **Checkpoint:** <commit/tree/report/review state>
- **Evidence:** <exact protected or exhausted condition>
- **Question / proposed direction:** <what must be decided>
- **Resolution condition:** <named Human or external-authority action>
- **Links:** <task spec, review, code reference, related DEC>
- **Human attention required:** yes
- **Status:** Open

---

(Add new escalations above this divider. Remove this block when none are open.)

## Protected Blocker Categories

TL-Lite or the full Tech Lead may allocate an ESC only for:

1. `product-decision`, `architecture-decision`, `data-model-decision`, or
   `public-interface-decision` requiring protected authority.
2. `dependency-approval` absent from the accepted dependency budget.
3. `security-or-data-loss` where continued work may be unsafe or destructive.
4. `external-infrastructure` that cannot be recovered locally.
5. `specification-contradiction` that leaves no authorized implementation.
6. `recovery-budget-exhausted` after two consecutive Phase Runner-mediated
   continuation attempts for one trigger signature make no material progress.

A cross-task or high-severity candidate is not automatically an ESC. Classify
the authority and repair boundary first: bounded, decision-free work is
`REPAIR`; only a category above is `BLOCK`.
