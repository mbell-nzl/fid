# Prompt: Workflow Retrospective + One-Shot Improvement Pass

> Hand this to a fresh Opus session at the repository root. It is a one-shot
> job: conduct the deferred deep retrospective, then make reasonable workflow
> improvements. It is **not** the agentic-workflow overhaul.

---

## Your Role

You are acting as **workflow owner** for this repository. You have authority
over the agent workflow surface — role prompts, skills, templates, repo
tooling, and the workflow spec. You have **no** authority over product code.

## Mission

Two parts, in order:

**Part A — Conduct the deep retrospective** that Phases 17 and 18 both
recommended and that has been deferred twice.

**Part B — Make the workflow changes** the retrospective justifies.

## The Bar

This is explicitly a **one-shot, reasonable-improvement pass**, not a
rigorous eval-driven iteration. The owner is separately building a larger
overhaul of the agentic engineering workflow that will supersede much of this
surface.

That sets your optimisation target:

- **"At least no worse than current" is the acceptance bar**, not "optimal".
- **Prefer cheap, high-leverage, low-sunk-cost changes.** Favour changing
  *affordances* (templates, worked examples, tool error messages, drill
  coverage) over adding *obligations* (new required steps, new artifacts, new
  checklists).
- **Do not build frameworks.** Anything elaborate is wasted work — the
  overhaul will replace it.
- If the highest-leverage fix for an item is "wait for the overhaul", say so
  and leave it. A short, honest "not fixed, here is why" is a better outcome
  than a speculative mechanism.

## Context To Load

In this order. Do not load the full spec suite.

**Workflow surface (what you may change):**

- `AGENTS.md` — project rules, commands, hard stops
- `project-docs/specs/09_AI_WORKFLOW.md` — the canonical workflow
- `project-docs/agents/*.md` — canonical role prompts (6 roles)
- `project-docs/skills/*/SKILL.md` — canonical skills
- `project-docs/templates/` — especially `completion-report-template.md`,
  `triage-log-template.md`, `tl-escalations-template.md`,
  `run-ledger-template.md`, `phase-runner-prompt.md`
- `project-docs/tools/README.md`, `report-lint.sh`, `workflow-lint.sh`,
  `sync-agents.sh`, `sync-skills.sh`
- `project-docs/tools/drills/README.md` and `run-drills.sh`
- `project-docs/docs/playbook.md`

**Evidence base (read-only history):**

- `project-docs/docs/phases/phase-17/decisions.md` — §"Interim Workflow
  Notes" (WN-1), DEC-058, the triage audit, and the deep-retrospective
  recommendation at the end
- `project-docs/docs/phases/phase-17/triage-log.md`,
  `phase-17/tl-escalations.md`, `phase-17/run-ledger.md`
- `project-docs/docs/phases/phase-18/decisions.md` — §"Triage-Log Audit",
  §"Workflow Notes", §"Deep Retrospective Recommendation", and DEC-072/073/074
- `project-docs/docs/phases/phase-18/triage-log.md`,
  `phase-18/tl-escalations.md`, `phase-18/run-ledger.md`
- `project-docs/docs/phases/phase-18/qa/final-review.md`
- `project-docs/docs/current/decisions_active.md` — the CF-023 row and
  DEC-083

## The Evidence Base, Summarised

You should verify these against the sources above rather than trusting this
summary, but it will save you discovery time.

### Theme 1 — Workflow-tooling reliability (Phase 17, recurred worse in Phase 18)

1. **Report-lint / completion-report commit-metadata friction.** Five of
   eight Phase 17 triages (TRG-002/003/004/005/007) were the *same* defect:
   commit metadata not backtick-wrapped, or the report left uncommitted so
   the lint evidence was invalid. Three escalated (ESC-004/005/011). Phase 17
   flagged a suspected contributing factor: the recurrence began after the
   Developer prompt and the DEC-053 Codex routing change, so it is plausibly
   a prompt/model-fit and template-affordance problem, not a judgment problem.
2. **Escalation-register integrity.** The template claims `ESC-NNN` IDs never
   reset, but ESC-001, ESC-004, and ESC-005 all collide across phases.
   Separately, escalations were raised without ever being written to the
   register and had to be backfilled at wrap.
3. **TRG numbering discipline.** Phase 17 produced a duplicate (Task 6
   recorded as both TRG-005 and TRG-017) that disagreed with the run ledger.
   There is no single authoritative numbering source.
4. **Phase 18 reproduced the whole class.** Six triages (TRG-001..006), all
   `BLOCK`, all correct on substance. But **five of six** reported
   "escalation required" or "escalated as ESC-NNN" while minting **no**
   identifier and writing **no** triage-log entry. Only TRG-006 produced a
   real `ESC-NNN` — and it collided with Phase 17's ESC-001. Run-ledger
   entries were also out of order / duplicated.

### Theme 2 — Evidence over-claiming (new in Phase 18; the expensive one)

Three separate reports asserted `PASS` where the evidence did not exist:
VAL-020, VAL-006, and four Task 5 gate commands. The Task 5 case is the
sharpest: a `workflow.Load` connection deadlock left four required commands
**non-terminating**, and they were recorded `PASS` anyway.

Every one was caught **only by QA Lead's independent rerun** — never by the
review cadence. Phase 18's own note: *"the cadence cannot detect a gate
recorded `PASS` without terminating."* Cost: a complete final-review rerun
(the most expensive artifact in the phase) plus a Developer remediation cycle.

The product signal was healthy — human checklist 0% fail, no `(human)`
failure mapping to a `(qa-lead)` PASS. The **process** signal was not.

Note that the existing drill set does not cover this class: `drill-03` is a
*silently skipped* required command and `drill-04` a *narrower substituted*
one. "Required command ran but never terminated, recorded `PASS`" is a
distinct defect class with no drill.

### Secondary observations worth carrying in

- **Task sizing.** Phase 18 Tasks 1-4 were right-sized; **Task 5 was
  overloaded** — consumer swap + maintenance parity + closure evidence +
  eight VAL assertions + a ~30s benchmark in one task. It absorbed four of
  the six triages.
- **QA cadence was correct; input trust was the gap.** Mid-phase reviews
  caught both evidence gaps before they compounded.
- **Token cost.** The seven-command gate ran end-to-end at least four times
  (~30s+ for the 100k-event rebuild row alone), one full final review was
  rerun from scratch, and DEC-067..DEC-072 total roughly 500 lines of
  decision prose that every subsequent session reloads.

### The pre-existing candidate scope (CF-023)

CF-023 in `project-docs/docs/current/decisions_active.md` lists eight
candidate items: (a) worked backtick-wrapped commit-metadata example in the
completion-report template; (b) a Developer-prompt rule that the report is its
own committed `docs(task)` step; (c) report-lint messaging distinguishing
dirty vs parse failure; (d) `ESC-NNN`/`TRG-NNN` global uniqueness; (e)
TL-Lite/Runner must mint a concrete `ESC-NNN` **and** a triage-log entry on
every `BLOCK`; (f) run-ledger chronological append with no duplicates; (g) a
wrap-time reopen audit for trigger-based carry-forwards; (h) evidence-honesty
hardening so a non-terminating required command cannot be recorded `PASS`.

**You are not bound to implement all eight.** Treat the list as candidate
scope, prioritise by leverage against the evidence, and state explicitly what
you left undone and why.

---

## Part A — The Retrospective

Write `project-docs/docs/reviews/2026-07-26-workflow-retrospective.md`.

Answer these directly, grounded in cited incidents (TRG/ESC/VAL/DEC IDs), not
in generalities:

1. **What is the actual failure mechanism behind Theme 1?** Is this a prompt
   defect, a template-affordance defect, a tooling defect, a model-fit
   artefact of the DEC-053 routing change, or a missing single source of
   truth for identifiers? Distinguish the ID-*allocation* problem from the
   ID-*recording* problem — they may have different fixes.
2. **What is the actual failure mechanism behind Theme 2?** Specifically: why
   did an agent record `PASS` for a command it never saw terminate, and what
   is the cheapest structural change that makes that harder? Consider what
   evidence a report *can* carry that a claim cannot fake.
3. **Which control actually caught these?** Be precise about what worked
   (QA Lead's independent rerun; the validation contract; the "supply
   evidence or fix the product, never weaken the assertion" rule) and what
   structurally cannot catch them.
4. **What should change, ranked by leverage per unit of effort?** For each,
   name the exact file to change and the change in one sentence.
5. **What should NOT change**, and what should wait for the overhaul.
6. **Was the deferral in DEC-083 correct in hindsight?** Record a plain
   yes/no with reasoning.

Keep it under about 200 lines. Cite; do not narrate.

## Part B — The Improvement Pass

### Surfaces you may change

- `project-docs/agents/*.md` — canonical role prompts (then run
  `sync-agents.sh`)
- `project-docs/skills/*/SKILL.md` — canonical skills (then run
  `sync-skills.sh`)
- `project-docs/templates/*.md`
- `project-docs/tools/*.sh` and `project-docs/tools/drills/`
- `project-docs/specs/09_AI_WORKFLOW.md`
- `AGENTS.md` — workflow rules only
- `project-docs/docs/playbook.md`

### Surfaces you must NOT change

- Any product code: `cmd/`, `internal/`, `migrations/`, `go.mod`.
- Specs `00`-`08` and `10`, unless a workflow rule genuinely lives there —
  and then say why.
- **Historical phase artefacts.** Do not renumber or rewrite past TRG/ESC
  IDs, triage logs, run ledgers, or completion reports. History is evidence.
- `.claude/agents/`, `.codex/agents/`, `.agents/skills/`, `.claude/skills/` —
  these are generated. Edit the canonical source and run the sync scripts.

### Constraints

- **No net increase in required per-task obligations.** If you add a required
  step, it must replace or subsume an existing one, or be directly justified
  by a named incident. Affordances over obligations.
- **Do not grow default-loaded context.** The size budgets in
  `decisions_active.md` are real: `phase_brief.md` 150 lines,
  `decisions_active.md` 150, `progress_current.md` 100,
  `tl_escalations_current.md` 80. Role prompts should not grow materially; if
  one must, justify it.
- **Mode-invariant rules never turn off** (`09_AI_WORKFLOW.md` §4): exact
  command gates, evidence rules, commit discipline, glossary discipline,
  every hard stop. You may strengthen these; you may not weaken them.
- **If you add a drill, follow the harness's own protocol** in
  `project-docs/tools/drills/README.md`: seed the defect, confirm the drill
  **fails** first, fix the workflow surface, confirm it **passes**. Report
  both observations. One defect per drill.

### Required verification before you finish

Run these and report each result verbatim:

```sh
project-docs/tools/sync-agents.sh
project-docs/tools/sync-skills.sh
project-docs/tools/sync-agents.sh --check
project-docs/tools/sync-skills.sh --check
project-docs/tools/workflow-lint.sh
project-docs/tools/report-lint.sh --self-test
project-docs/tools/drills/run-drills.sh
```

**The "no worse than current" gate:** `run-drills.sh` mechanical drills must
finish at **0 misses**. Capture the baseline **before** you change anything so
you can prove you did not regress it. If any command does not terminate, say
so — do not record it as passing. (That rule is the subject of this
retrospective; do not be its next data point.)

Agent drills (05-08) require spawning subagents. If you cannot run them,
record that plainly as unverified rather than implying they passed.

### Recording decisions

Global `DEC-NNN` numbering is continuous and currently ends at **DEC-083**.
Record your decisions as **DEC-084** onward inside the retrospective document,
then:

- Add a short pointer block to `project-docs/docs/current/decisions_active.md`
  (a few lines, not the full text — respect the size budget).
- Update the **CF-023** row there to reflect what you actually landed versus
  what remains for the overhaul. Close it only if genuinely satisfied;
  otherwise narrow its scope honestly.

### Interaction with the pending Phase 19 plan

Phase 19 is planned and awaiting approval; its artefacts are in
`project-docs/docs/phases/phase-19/`. Its brief currently carries two
stand-in mitigations from DEC-083: the evidence-honesty rule, and phase-scoped
`ESC-19-NNN` / `TRG-19-NNN` identifiers.

If your changes supersede either, you may update
`project-docs/docs/current/phase_brief.md` and the Phase 19 task specs to
match — but **flag every such edit explicitly** in your summary. Do not
silently alter an approved-pending plan.

---

## Definition Of Done

- [ ] Retrospective written, citing real incident IDs, under ~200 lines.
- [ ] Changes made, each traceable to a named incident or theme.
- [ ] Canonical sources edited, never generated copies; both sync scripts run
      and `--check` clean.
- [ ] `workflow-lint.sh` and `report-lint.sh --self-test` green.
- [ ] `run-drills.sh` mechanical drills at 0 misses, with the pre-change
      baseline recorded for comparison.
- [ ] Any new drill demonstrated failing-then-passing.
- [ ] Decisions recorded as DEC-084+; `decisions_active.md` pointer added and
      CF-023 updated.
- [ ] Conventional Commit(s) for the work.

## Final Output

Close with a concise summary containing:

1. **A change table:** file → change → the incident it addresses.
2. **What you deliberately did not fix**, and whether it is "wait for the
   overhaul" or "not worth it".
3. **Verification results**, verbatim, including the drill baseline and final
   counts.
4. **Residual risk** in one short paragraph: what can still go wrong in the
   next phase, and what would catch it.
5. **Any edit you made to the pending Phase 19 plan.**

Be direct about uncertainty. If a change is a guess that the overhaul should
re-examine, label it as such rather than presenting it as settled.
