# Product And Feature-Discovery Audit

**Date:** 2026-05-28
**Scope:** fid specification suite, roadmap, and origin scaffold
**Lens:** Product strategy and net-new feature opportunity
**Non-goal:** Correctness, consistency, traceability, or under-specification
review.

## Summary

The current fid suite already covers most obvious issue-tracker primitives:
blockers, labels, inbox, cycles, saved views, named reports, comments, TUI,
MCP, markdown export, and conditional import/sync. Those are not treated as
findings here.

The strongest missing product opportunities are not generic Jira or Linear
features. They are agent-loop accelerators that exploit fid's append-only
journal and local-first structure:

- role-aware "what changed since..." review packs
- safe retryable writes via dry-run, preview, and idempotency keys
- one-transaction task-pack application
- append-only correction/undo commands
- a pull-based human attention queue

The strategic opportunity is to make fid better than markdown artifacts at
review, recovery, and resumption. Markdown can store prose; fid can explain
state changes, prove evidence, and make agent retries safe.

## Existing Coverage To Exclude

These primitives are already planned or explicitly excluded:

| Capability | Current status |
|---|---|
| Issue dependencies/blockers | `blocks` / `blocked_by` links in PRD-LINK-004 and ready-queue behavior |
| Labels, priority, estimates, due dates, inbox | v0.1 issue metadata |
| Cycles/sprints | v0.5 optional cycle support |
| Saved views | v0.4 `views` table and CLI |
| Reports/search | v0.2 named reports, conditional v0.3 query DSL, read-only SQL |
| Comments | v0.2 comments, v0.3 QA-review scope expansion |
| TUI/dashboard | v0.3 TUI board, v0.5 web dashboard |
| MCP integration | v0.3 curated MCP resources, prompts, and tools |
| Bidirectional markdown sync | v0.5+ conditional, pending evidence |
| Watchers/subscribers, SLAs, OKRs, customer requests | Explicitly out of scope |

External analogous-tool scan confirmed that dependencies, saved views,
notifications, and bulk selection are common issue-tracker primitives in
Linear, GitHub Issues, and Jira. They are useful context but not the main
opportunity for fid because the suite already borrows or rejects them
deliberately:

- Linear issue relations: <https://linear.app/docs/issue-relations/>
- GitHub Issues and Projects feature overview:
  <https://github.com/features/issues>
- Jira dependency documentation:
  <https://support.atlassian.com/jira-software-cloud/docs/create-or-remove-dependencies-on-your-timeline/>

## Prioritized Feature Opportunities

### 1. Journal-Powered Change Review Packs

**Pitch:** Show exactly what changed since a checkpoint, in the format the
current role needs.

**User and job-to-be-done:** QA Reviewer, QA Lead, TL, and Human need to
review only the delta since the last task, review, phase checkpoint, event, or
date.

**Gap today:** fid has a canonical journal, projected `activity_log`, basic
`fid activity list`, and role-tuned context packs, but no first-class
"what changed since" or entity diff surface. See:

- `project-docs/specs/00_PROJECT_CHARTER.md` section 4
- `project-docs/specs/10_IMPLEMENTATION_ROADMAP.md` v0.1e
- `project-docs/specs/07_DESIGN_PLAN_FLOWS.md` section 4

**Proposed capability:**

```text
fid changes --since <event|time|review|commit> [--epic <id>] [--role <role>]
fid diff <entity-id> --from <event> --to <event>
fid issue get <id> --context-pack qa_reviewer --since <review-id>
```

Return changed entities, changed fields, new evidence, new concerns/findings,
new links, status transitions, claims, and task completions. Role mode filters
the output: a QA Reviewer sees code/review deltas; a TL sees decisions,
escalations, carry-forward, and validation movement; a Human sees approval
and judgment deltas.

**Why fid specifically:** The append-only journal makes this cheap and
trustworthy. Markdown artifacts can only approximate it through git diff plus
agent interpretation.

**Value and fit:** High value, M effort, suggested v0.2. This is the highest
leverage event-log feature.

### 2. Safe Write Protocol: Dry-Run, Preview, And Idempotency Keys

**Pitch:** Let agents preview writes and retry commands without duplicating
state.

**User and job-to-be-done:** Developer, Workflow Runner, TL, and QA agents
need to issue composite commands safely, recover from timeouts, and avoid
duplicate completions, AC verifications, concerns, findings, or task creation.

**Gap today:** Errors are designed to teach, and composites are atomic, but
there is no write preview and no retry-safe command identity. See:

- `project-docs/specs/05_DESIGN_PLAN.md` section 1
- `project-docs/specs/02_TRD.md` section 4.3
- `project-docs/specs/01_PRD.md` section 3.14

**Proposed capability:**

```text
fid dev finish <id> ... --dry-run
fid issue create ... --dry-run
fid issue ac verify ... --idempotency-key <key>
fid dev finish ... --idempotency-key <key>
```

`--dry-run` returns the validation result and would-be journal event effects
without appending. `--idempotency-key` stores a command key in the journal
event envelope; a repeated key returns the original event/result rather than
creating a duplicate row.

**Why fid specifically:** One logical write already maps to one journal event,
so idempotency can live at the event envelope layer without a daemon or server.

**Value and fit:** High value, M effort, suggested v0.1e or v0.2. This is
especially important for MCP and autonomous agent use.

### 3. Structured Task-Pack Apply

**Pitch:** Create an epic, tasks, ACs, links, order, and metadata from one
reviewed task-pack transaction.

**User and job-to-be-done:** TL agents need to turn a planned phase/task pack
into fid state without dozens of fragile CLI calls.

**Gap today:** The planning flow shows repeated `issue create`, `ac add`, and
`link add` commands. Composites focus on execution, not planning ingestion.
See:

- `project-docs/specs/07_DESIGN_PLAN_FLOWS.md` section 1.1
- `project-docs/specs/05_DESIGN_PLAN.md` section 4
- `project-docs/specs/09_AI_WORKFLOW.md` section 5

**Proposed capability:**

```text
fid task-pack validate --from <task-pack.yaml|task-pack.md>
fid task-pack apply --from <path> --dry-run --by tl
```

The input defines epic metadata, tasks, ACs, linked requirements, dependencies,
labels, token estimates, and sequence/priority. `apply` produces one journal
event with many effects and returns a mapping from source task references to
generated fid IDs.

**Why fid specifically:** The journal can record the phase plan as one
auditable event while SQLite receives normalized tasks, ACs, and links.

**Value and fit:** High value, M effort, suggested v0.2. Consider v0.1 if
dogfood planning ceremony becomes a serious friction point.

### 4. Append-Only Corrections And Undo

**Pitch:** Fix mistaken writes without rewriting history or manually editing
journal files.

**User and job-to-be-done:** Humans and agents need to recover from wrong
status transitions, bad labels, mistaken AC verification, duplicate links, or
accidental completion records.

**Gap today:** The rollback process is git revert on journal files plus
`fid rebuild`, which is coarse. There is no domain-level correction command.
See:

- `project-docs/specs/02_TRD.md` section 8
- `project-docs/specs/01_PRD.md` section 3.3
- `project-docs/specs/01_PRD.md` section 3.4

**Proposed capability:**

```text
fid event show <event-id>
fid correction create --reverses <event-id> --reason-md <path> --dry-run
fid ac retract <acv-id> --reason-md <path> --by human
fid link remove <from> <relation> <to> --by <role>
fid issue restore-status <id> --to <state> --reason-md <path> --by human
```

Correction commands append explicit reversal or supersession events. They do
not delete historical rows.

**Why fid specifically:** Append-only corrections preserve audit value while
avoiding manual journal surgery.

**Value and fit:** High value, M effort, suggested v0.2.

### 5. Human Attention Queue

**Pitch:** Tell the human what needs judgment now, without adding
notifications or watchers.

**User and job-to-be-done:** The Human orchestrator needs a single low-noise
command before approving, interrupting, or resuming a phase.

**Gap today:** Open escalations, human checklist items, human claims, failed
validations, intervention flags, blocked outcomes, and unresolved decisions
exist across separate entities and reports. Watchers/subscribers and
notifications are explicitly out of scope. See:

- `project-docs/specs/00_PROJECT_CHARTER.md` section 5.3
- `project-docs/specs/01_PRD.md` section 2
- `project-docs/specs/09_AI_WORKFLOW.md` section 3.1

**Proposed capability:**

```text
fid attention [--epic <id>] [--role human|tl] [--format json|toon|markdown]
```

Return items grouped by reason:

- needs decision
- needs approval
- blocked
- human-claimed
- human intervention recorded
- validation failed
- validation carried forward
- stale unresolved escalation

**Why fid specifically:** Pull-based attention respects local-first,
no-daemon, no-notification constraints while replacing human grep across
current files.

**Value and fit:** High value, S effort, suggested v0.2.

### 6. Evidence Registry And Evidence Check

**Pitch:** Evidence refs should be verifiable, not just strings.

**User and job-to-be-done:** QA Lead and QA Reviewer need to know whether
cited tests, commands, commits, screenshots, files, and generated artifacts
still exist and still support claims.

**Gap today:** AC verifications, QA findings, and validation outcomes store
evidence refs/excerpts, but no command verifies evidence integrity. See:

- `project-docs/specs/01_PRD.md` section 3.3
- `project-docs/specs/01_PRD.md` section 3.9
- `project-docs/specs/01_PRD.md` section 3.10

**Proposed capability:**

```text
fid evidence check --epic <id>
fid evidence show <evidence-id|ref>
fid evidence capture --cmd "go test ./..." --attach-to AC-...
```

Track evidence type, path/ref, hash where applicable, command output excerpt,
captured-at timestamp, and verification status.

**Why fid specifically:** fid already treats deterministic evidence as central.
Local filesystem and git refs can be checked without a network service.

**Value and fit:** High value, M effort, suggested v0.3 or v0.4.

### 7. Git-Aware Run And Resume State

**Pitch:** Make branch, commit, dirty state, and abandoned-run recovery
first-class.

**User and job-to-be-done:** Workflow Runner and Human need to resume after
crashes, inspect active runs, and verify that a completion's commit/files match
the recorded work.

**Gap today:** `task_completions` stores commits/files and Phase Runner
verifies commits outside fid. `agent_runs` has token metadata but no
branch/base/dirty/resume model. See:

- `project-docs/specs/03_TRD_SCHEMA.md` section 2.3
- `project-docs/specs/03_TRD_SCHEMA.md` section 7
- `project-docs/specs/09_AI_WORKFLOW.md` section 2.4

**Proposed capability:**

```text
fid run list --open|--abandoned
fid run show <run-id>
fid dev resume <task-id>
fid dev abandon <run-id> --reason-md <path> --by <role>
```

Capture branch, base commit, head commit, dirty flag, worktree path, and
resume hints at `dev start` / `dev finish`.

**Why fid specifically:** Agent work is local and git-backed. fid can make
resume and review state queryable without becoming a git replacement.

**Value and fit:** High value, M effort, suggested v0.2.

### 8. Ready Queue Explain And Candidate Ranking

**Pitch:** Return not just the next task, but why it won.

**User and job-to-be-done:** Workflow Runner and Human need confidence that
`fid ready` chose the same work a good operator would choose.

**Gap today:** `fid ready` returns one task or an empty explanation.
`issue list` can filter, but there is no ranked candidate set with skipped
reasons. See:

- `project-docs/specs/01_PRD.md` section 3.2
- `project-docs/specs/06_DESIGN_PLAN_CLI.md` section 4
- `project-docs/specs/05_DESIGN_PLAN.md` section 9

**Proposed capability:**

```text
fid ready --limit 5 --explain
fid ready --epic <id> --explain --include-blocked
```

Return candidates, score/order factors, dependency blockers, claim blockers,
sequence blockers, stale claims, and skipped reasons.

**Why fid specifically:** Ready queue selection is where agent orchestration
meets project policy. Explainability reduces wrong autonomous picks.

**Value and fit:** High value, S effort, suggested v0.1 or v0.2.

### 9. Role Handoff And Prompt Packs

**Pitch:** Generate exact child-agent prompts and context manifests from fid
state.

**User and job-to-be-done:** Phase Runner needs to spawn Developer, QA
Reviewer, QA Lead, and TL-Lite with explicit reads, commands, stop conditions,
output paths, and current blockers.

**Gap today:** Context packs and MCP prompts are planned, but
`09_AI_WORKFLOW.md` still requires Phase Runner to hand-author explicit child
prompts and maintain a ledger in response notes. See:

- `project-docs/specs/09_AI_WORKFLOW.md` section 2.4
- `project-docs/specs/09_AI_WORKFLOW.md` section 3.4
- `project-docs/specs/10_IMPLEMENTATION_ROADMAP.md` v0.3d

**Proposed capability:**

```text
fid handoff developer <task-id>
fid handoff qa-reviewer --since <task-id>
fid handoff qa-lead --epic <id> --final
fid handoff tl-lite --source <concern-id|finding-id>
```

Output prompt text plus a machine-readable manifest of required reads,
expected writes, commands, stop conditions, and evidence obligations.

**Why fid specifically:** fid knows the role model, task state, validation
obligations, current blockers, and context-budget concerns.

**Value and fit:** High value, M effort, suggested v0.2 or v0.3.

### 10. Formal Intake Workflow, Not Just An Inbox Flag

**Pitch:** Capture loose discoveries and promote them deliberately.

**User and job-to-be-done:** Developers, QA, and Human need to park bugs,
tech debt, follow-ups, and discoveries without prematurely assigning them to
an epic.

**Gap today:** `in_inbox` exists and is carefully distinguished from TL-Lite
triage, but there is no intake command set beyond issue create/list/update.
See:

- `project-docs/specs/GLOSSARY.md` section "Adjudication and Intake"
- `project-docs/specs/01_PRD.md` section 3.2

**Proposed capability:**

```text
fid intake add --from-issue <id> --summary <text> --by <role>
fid intake add --from-finding <id> --by <role>
fid intake list
fid intake promote <id> --epic <id> --priority <n> --links ...
fid intake close <id> --duplicate-of <issue-id>
```

Could initially be syntactic sugar over `issues.in_inbox`, labels, links, and
comments; later it can become a stronger workflow if dogfooding shows value.

**Why fid specifically:** Keeps local-first structured capture while avoiding
watchers/customer-request/product-feedback scope creep.

**Value and fit:** Medium-high value, M effort, suggested v0.2.

### 11. Context Budget Metering

**Pitch:** Show how much context a pack or handoff is expected to consume.

**User and job-to-be-done:** Agents and Human orchestrators need to avoid
overloading a Developer, QA Reviewer, or TL-Lite session.

**Gap today:** fid optimizes for token-efficient output through TOON and
context packs, and the workflow keeps default context small. It does not
estimate or report context-pack size before use. See:

- `project-docs/specs/00_PROJECT_CHARTER.md` section 4
- `project-docs/specs/05_DESIGN_PLAN.md` section 1
- `project-docs/specs/09_AI_WORKFLOW.md` section 4

**Proposed capability:**

```text
fid context estimate --pack developer --issue <id>
fid issue get <id> --context-pack developer --budget-tokens 12000
fid context explain --pack qa_reviewer --issue <id>
```

Return estimated tokens by section and optionally omit lower-priority sections
when a budget is passed.

**Why fid specifically:** Agent-native context control is core to fid's
positioning, and all source material is local and structured enough to measure
roughly.

**Value and fit:** Medium-high value, S/M effort, suggested v0.2.

### 12. Pattern Mining From Rework And Findings

**Pitch:** Turn repeated failures into candidate decisions, principles, or
templates.

**User and job-to-be-done:** TL and QA Lead need to detect recurring failure
patterns across attempts, findings, concerns, and retrospectives.

**Gap today:** Stats and computed retro signals cover estimates, attempts,
gates, interventions, finding density, contract pass-rate, and carry-forward
count. They do not mine recurring textual or structured themes into actionable
recommendations. See:

- `project-docs/specs/01_PRD.md` section 3.12
- `project-docs/specs/10_IMPLEMENTATION_ROADMAP.md` v0.4b
- `project-docs/specs/09_AI_WORKFLOW.md` section 9

**Proposed capability:**

```text
fid insights patterns --epic <id>
fid insights recurring-findings --since <date>
fid insights suggest-principles --from findings,concerns,retro
```

Initial implementation can be deterministic: group by trigger class,
label, area, gate command, finding severity, affected file path prefix,
and repeated validation assertion failures.

**Why fid specifically:** The graph connects tasks, attempts, gates,
findings, concerns, decisions, and retrospectives in one queryable place.

**Value and fit:** Medium value, M effort, suggested v0.4.

## Top Five Highest-Leverage Opportunities

1. **Journal-powered change review packs.** Biggest direct advantage over
   markdown and git grep.
2. **Safe write protocol.** Makes autonomous and MCP-driven writes much less
   brittle.
3. **Structured task-pack apply.** Removes planning ceremony and reduces
   multi-command agent mistakes.
4. **Append-only corrections and undo.** Makes the event-sourced model usable
   when humans or agents make mistakes.
5. **Human attention queue.** Gives the Human a clear operational surface
   without violating local-first or no-notification principles.

## Strategic Note

The current feature set is likely sufficient to hit the stated v0.1 dogfood
goal of replacing roughly 70 percent of daily interaction with whet task
artifacts. The core replacement path is covered: issue CRUD, ACs, completions,
concerns, claims, ready queue, links, telemetry, and journal/rebuild.

The other 30 percent is orchestration friction:

- What changed since the last agent touched this?
- Can this command be retried safely?
- How do we recover from a mistaken write?
- What exact prompt and context should the next agent receive?
- What evidence is stale or missing?
- What needs human judgment now?

The single capability that would most move fid from "viable replacement" to
"clearly better than the artifacts it replaces" is:

```text
fid changes --since ... --role ...
```

made role-aware and context-pack aware. It directly exploits the journal as a
product feature. Plain markdown can preserve history through git; fid can turn
history into targeted, low-token, role-specific review input.

