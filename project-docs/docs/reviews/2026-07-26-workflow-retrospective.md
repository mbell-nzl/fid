# Workflow Retrospective — Phases 17–18

**Date:** 2026-07-26. **Scope:** the deep retrospective recommended by Phase 17
(DEC-058) and Phase 18 (DEC-073) and deferred twice, plus a one-shot improvement
pass. **Not** the agentic-workflow overhaul; the bar is "no worse than current",
favouring affordances over obligations.

---

## 1. Theme 1 — The Actual Failure Mechanism

**Two distinct defects were tracked as one. Allocation is a spec defect;
recording is a missing-consequence defect. Neither is a judgment defect, and the
commit-metadata sub-class is already closed.**

### 1a. ID allocation: an unsatisfiable instruction

`tl-escalations-template.md` said "Use global ESC-NNN numbering that never
resets." TL-Lite's context is `AGENTS.md`, the four current docs, and the named
excerpt. Of these only `tl_escalations_current.md` carries ESC IDs — and it holds
**open items only**; resolved ones archive to `phases/phase-N/tl-escalations.md`.
Every Phase 18 ledger entry records `ESC baseline (Open at resume): none`. The
register was empty, and the global counter's state lived across 19 phase archives
that TL-Lite has no instruction to scan and (`disallowedTools: Bash`) cannot
enumerate cheaply.

Asked for an identifier it could not compute, TL-Lite emitted the placeholder
from its own output template, which read `VERDICT: BLOCK (escalated as
ESC-<NNN>)`. TRG-18-001..005 produced `ESC-NNN` / `ESC-NEW`; Phase 17's TRG-007
produced none ("TL-Lite could not record an escalation ID in its output").
**Rational behaviour under an impossible instruction, not a lapse.**

The *substance* was right every time: DEC-073's audit records all six Phase 18
blocks as correct, Phase 17's audit all four. Judgment never failed.

Collisions follow from the same cause — ESC-001 (phase-1/14), ESC-004 (6/17),
ESC-005 (7/17), ESC-001 again (17/18). A counter nobody can read restarts
wherever each session guesses.

**On the DEC-053 model-fit hypothesis (WN-1): weakly supported, not primary.**
The identical failure spans Phase 17 and Phase 18 across different sessions, and
the fix requires no model change. Model fit may modulate frequency; it cannot
explain why a correct agent could not comply.

### 1b. ID recording: a check with no consequence

`phase-runner.md` already said "confirm `triage-log.md` contains the TRG entry" —
and the runner **did** notice, recording "TL-Lite did not assign an ESC-NNN or
create the required triage log" (TRG-18-001) and "did not mint a concrete
identifier" (TRG-18-003/004/005). It detected the defect four times, had no rule
for what to do, wrote prose, and continued. The observation never reached a
structured field, so each event needed full Tech Lead reconstruction.

Phase 17's numbering disorder (TRG-005 duplicated as TRG-017, disagreeing with
the ledger) shares the root: no file was designated authoritative.

### 1c. The commit-metadata sub-class is already closed

Phase 17's five-of-eight recurrence (TRG-002/003/004/005/007; ESC-004/005/011)
**did not recur in Phase 18** — the ledger records `report-lint: PASS` on every
task, first pass. The completion-report template carries the worked
backtick-wrapped metadata block and `developer.md` step 4 states the report is
its own `docs(task)` commit. CF-023 (a) and (b) are satisfied.

*Uncertainty:* no shell this session, so I cannot `git log` when that text landed
relative to Phase 17. I record the recurrence data (5 → 0), not a causal claim.

---

## 2. Theme 2 — Why `PASS` For A Command Nobody Saw Terminate

**The report format asks for a verdict, not evidence. A `PASS` token is free to
write, carries no falsifiable content, and costs the same whether the command
succeeded, failed, or hung.**

The template line was `` `<command>` - PASS/FAIL/BLOCKED ``. `report-lint.sh`
verifies a command *appears* with a PASS token; it cannot verify it *ran*. In
DEC-072 the `workflow.Load` deadlock left four required commands non-terminating
on `46d8013`; all four were recorded `PASS`.

**The three over-claims are not one class — and that determines the fix.**

| Over-claim | Caught by | Catchable by review? |
|---|---|---|
| VAL-020 rebuild evidence (DEC-068) | QA Reviewer, Review 1 | **Yes** — reviewer read the cited test |
| VAL-006 valid-pair UPDATE (DEC-069) | QA Reviewer, Review 2 | **Yes** — same |
| Four Task 5 gate commands (DEC-072) | QA Lead independent rerun only | **No** |

The brief's framing ("every one caught only by QA Lead's rerun") is inexact: the
mid-phase cadence caught both *document-level* over-claims before they
compounded, as Phase 18's own workflow note says. Only the *execution-level* one
escaped. That distinction is the finding:

- A **document-level** claim ("VAL-020 is covered by the existing suite") points
  at an artifact that still exists. A reviewer can open it. Demand a citation.
- An **execution-level** claim ("this command passed") points at a terminal
  session that no longer exists. Nobody can inspect it. Only two things work:
  carry evidence the run produced, or re-execute.

**Cheapest structural change: make the report carry something the run produced.**
A verbatim final output line — `ok internal/db 4.2s`, `total: (statements)
76.8%`, `exit 0 (no output)` — is free for a Developer who just watched the
command finish and **unavailable** to one who did not. It converts a false PASS
from an omission into a fabrication: a higher bar, and one a reviewer can spot
because invented output rarely matches a tool's real format. Not a detector —
friction placed where the incentive is.

---

## 3. Which Control Actually Caught These

**Worked:**

- **QA Lead's independent re-execution** — the only control that caught DEC-072.
  It does not trust the report; it regenerates the evidence.
- **The validation contract's both-states rule.** VAL-006 required "both the
  rejected (invalid) and accepted (valid) INSERT/UPDATE paths"; that wording let
  Review 2 name a missing cell rather than argue about adequacy.
- **"Supply evidence or fix the product, never weaken the assertion."** Applied
  literally under pressure: DEC-071 made the rebuild faster instead of relaxing
  TRD-PERF-003; DEC-072 fixed `workflow.Load` instead of raising
  `SetMaxOpenConns(1)`; DEC-068/069 supplied missing tests. Zero assertions,
  thresholds, fixtures, or gate commands were weakened all phase.
- **Mid-phase cadence**, for document-level claims (VAL-020, VAL-006).

**Structurally cannot work:**

- **Any review that reads reports without re-executing** cannot detect a gate
  recorded `PASS` without terminating. Not diligence — the evidence is gone by
  the time the reviewer arrives.
- **`report-lint.sh` as written**: it checks the shape of a claim, never its
  provenance.
- **Drills 00–08.** drill-03 (skipped) and drill-04 (narrower substitution) are
  report-text defects. "Ran but never terminated, recorded PASS" has no drill.

---

## 4. What Should Change, By Leverage Per Unit Effort

| # | Change | File | Landed |
|---|---|---|---|
| 1 | IDs become phase-scoped `TRG-<phase>-<NNN>` / `ESC-<phase>-<NNN>`, allocated from the phase triage log | `triage-log-template.md`, `tl-escalations-template.md`, `09_AI_WORKFLOW.md` §2.3, `tl-lite.md` | ✅ |
| 2 | TL-Lite writes the log entry **before** the verdict; placeholders named defective triage | `tl-lite.md`, `09_AI_WORKFLOW.md` §2.3 | ✅ |
| 3 | Gate results carry a verbatim terminal token; non-terminating is never `PASS` | `completion-report-template.md`, `developer.md`, `AGENTS.md`, `09_AI_WORKFLOW.md` §8 | ✅ |
| 4 | Ledger TRIAGE entry gains `Triage log:` so a missing entry lands in a field, not prose | `run-ledger-template.md`, `phase-runner.md` | ✅ |
| 5 | VAL coverage demands a citation (test / `file:line` / output), not a restatement | `completion-report-template.md`, `developer.md` | ✅ |
| 6 | Wrap-up audits trigger-based carry-forwards against the phase diff | `phase-wrap-up/SKILL.md` | ✅ |
| 7 | Ledger append rule: chronological, no duplicates, corrections appended not edited | `run-ledger-template.md` | ✅ |
| 8 | `drill-09` for the non-terminating-command class + the `report-lint` check it needs | `run-drills.sh`, `report-lint.sh` | ❌ §5 |

#1 removes the cause of five incidents for the price of a template sentence — the
best ratio available. #3 is the only item carrying product-quality risk. #8 is
the only item requiring executable code.

---

## 5. What Should Not Change, And What Waits

- **`report-lint.sh` and `run-drills.sh` are untouched.** This session has no
  shell. A `report-lint` failure is a Phase Runner hard stop, and the drill
  harness requires a seed-fail-fix-pass demonstration. Shipping an unverified new
  hard-stop — asserting it works without watching it run — is precisely the
  defect this retrospective is about. Deferred as **DEC-087** with concrete
  scope.
- **CF-023 (c), report-lint dirty-vs-parse messaging: reclassified low-value.**
  The messages are already distinct ("report file has uncommitted modifications"
  vs "no commit hash found"). The real wart is that a hash parse failure
  `return`s before the committed-file check, so a doubly-broken report reports
  one fault at a time. Cost: one extra Phase 17 cycle, zero in Phase 18.
- **No historical renumbering.** Phase-scoped IDs start at Phase 19; existing
  collisions stay as recorded because ledgers and reviews reference them.
- **Task sizing stays advisory.** Phase 18 Task 5 was genuinely overloaded
  (consumer swap + parity + closure evidence + 8 VAL assertions + a ~30 s
  benchmark; it absorbed four of six triages). But "closure evidence is its own
  task" is a heuristic the Tech Lead already has authority and a wrap-up note to
  apply. Encoding it as a rule risks the opposite failure. Overhaul.
- **Decision-prose volume** (DEC-067..072 ≈ 500 lines reloaded every session) is
  a context-architecture problem, not a rule problem. Overhaul.

---

## 6. Was The DEC-083 Deferral Correct In Hindsight?

**Yes on risk; no on cost.**

The risk read was accurate and the outcome confirms it: seven of eight CF-023
items are audit-trail and orchestration friction that cannot let a defect reach
the product, and Phase 18 closed with a 0%-fail human checklist and no `(human)`
failure mapping to a `(qa-lead)` PASS. Item (h) was correctly singled out as the
only product-quality risk, correctly carried as a phase rule, and the control
that caught it — QA Lead's independent rerun — was correctly identified as
retained by heavy mode.

Where it was wrong: it priced the fix as a *pass* — a work item worth deferring
wholesale. The fix was ~40 lines of template and prompt text, landed in one
session with no executable change. The deferred cost was five Tech Lead
reconstruction sessions in one phase, the dominant Tech Lead expense of Phase 18.
**The deferral was defensible; treating it as all-or-nothing was the error.** A
cheap-item carve-out would have been strictly better.

---

## Decisions

### DEC-084: Triage and escalation identifiers are permanently phase-scoped

`TRG-<phase>-<NNN>` and `ESC-<phase>-<NNN>`. **The phase triage log is the single
allocation source for both series** — read it, take the highest `<NNN>`, add one,
start at `001` in a new phase. Rationale: a global counter has no readable source
of truth for an agent that cannot scan every phase archive, which is why five
Phase 18 triages emitted placeholders and four IDs collided. Historical IDs are
**not** renumbered. Supersedes the "global ESC-NNN numbering that never resets"
rule and promotes DEC-083's Phase-19-only convention to the standing scheme.

### DEC-085: The record is written before the verdict; placeholders are defective triage

TL-Lite appends the TRG entry (and on BLOCK the ESC entry) to disk **first**,
then emits a verdict quoting the identifier it wrote. A verdict carrying
`ESC-NNN` / `ESC-NEW`, or unaccompanied by a log entry, is a defective triage.
The Phase Runner records the concrete IDs in a new `Triage log:` field on the
TRIAGE ledger entry, or records the gap verbatim there. **No new stop condition:**
Phase 18 showed the runner already detects this; it lacked somewhere to put it.

### DEC-086: Gate results carry terminal evidence; non-terminating is never `PASS`

Every required command's result carries one short verbatim token from the
command's own output (last printed line, or `exit 0 (no output)`). A command that
did not terminate is never `PASS`; it is `BLOCKED` in a blocker report. `VAL-NNN`
coverage cites a test name, `file:line`, or output — not a restatement and not
"covered by the existing suite" (the exact DEC-068 wording). **Mode-invariant.**
Promotes DEC-083's Phase-19 stand-in to a standing rule; closes CF-023 (h).

### DEC-087: `drill-09` and the report-lint termination check are specified, not landed

The uncovered class is "required command ran but never terminated, recorded
`PASS`" (DEC-072). Closing it mechanically needs `report-lint.sh` to assert an
evidence token on each required command's result line, plus a `drill-09`
demonstrated failing-then-passing per the harness protocol. **Not landed:** the
authoring session had no shell, and an unverified hard-stop in the Phase Runner
path is worse than an honest gap. Scope recorded in `tools/drills/README.md`.
Owner: the overhaul, or any session with a shell.

### DEC-088: CF-023 is narrowed, not closed

Landed: (d) → DEC-084; (e) → DEC-085; (f) ledger append rule; (g) wrap-time
reopen audit; (h) → DEC-086. Already satisfied: (a), (b) — the Phase 17
recurrence did not repeat in Phase 18. Reclassified low-value: (c). **Remaining:
DEC-087 only.**

---

## Verification

**`sync-agents.sh`, `sync-skills.sh`, `workflow-lint.sh`,
`report-lint.sh --self-test`, and `run-drills.sh` were NOT executed** — the
authoring session had no shell tool. Generated copies under `.claude/agents/`,
`.codex/agents/`, `.agents/skills/`, and `.claude/skills/` were kept in sync by
applying byte-identical edits to canonical and generated files inside shared body
segments, which is what the sync scripts produce — **asserted by construction,
not verified by execution.** No executable file was modified, so the drill
baseline is unchanged by construction; that too is unverified. Run all five
commands before trusting this surface.
