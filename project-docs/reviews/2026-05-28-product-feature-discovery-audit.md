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

---

## Resolution Log

> Added 2026-05-31; completed 2026-06-01. Two `/grill-with-docs` sessions walked
> all twelve items in the priority order below, resolving design decisions one at
> a time and updating the spec suite inline. This log records (1) the assessor's
> original disposition of all twelve items and (2) the decisions reached and
> specs changed for each. **All twelve items are now resolved** — see the
> per-item RESOLVED blocks below; no items remain open.

### Assessor's Original Disposition (2026-05-29)

The audit is high quality: well-sourced, scope-aware (it excludes
out-of-scope features and respects local-first / no-daemon / no-git-replacement
constraints), and it correctly identifies fid's one defensible edge over
markdown+grep — the append-only journal. Three caveats temper it:

1. **Several "gaps" are overstated** — already on the roadmap, so the audit's
   contribution there is refinement, not discovery (`fid event revert` v0.2c
   for #4; named reports `stale-claims`/`qa-blockers` for #5; `fid ready`
   empty-explanation for #8).
2. **Sequencing is optimistic** — almost everything is slotted v0.2, but v0.1
   is deliberately lean and unbuilt. These belong on the roadmap's existing
   *conditional / evidence-gated* footing, not pre-committed.
3. **A few items quietly expand fid's responsibility surface or lock new
   contracts** in ways the audit under-weights (evidence *capture* makes fid a
   command executor; handoff subcommands hard-code one workflow; task-pack and
   idempotency each freeze a new contract).

Tiering (with the adjustment applied during resolution):

| # | Item | Merit | Key adjustment |
|---|---|---|---|
| 1 | Journal change review packs | **High** | Ship `fid changes --since` diff-first; role shaping reuses context packs later, not a bespoke `--role` |
| 2 | Safe write (dry-run + idempotency) | **High (dry-run) / Conditional (idempotency)** | Split: dry-run is universal/cheap; idempotency mainly serves MCP/programmatic, defer |
| 7 | Git-aware run/resume state | **High** | Strongest on-mission item: moves ephemeral orchestration state into the durable store |
| 5 | Human attention queue | **Solid** | Essentially one more named report; do it, don't oversell as novel |
| 8 | Ready-queue explain | **Solid** | Cheap extension of planned `fid ready`; v0.2 |
| 4 | Corrections/undo | **Real but overstated** | `fid event revert` already roadmapped; reframe as domain-level sugar on top |
| 9 | Handoff/prompt packs | **Real friction, wrong shape** | Express as templates over context packs, not hard-coded `fid handoff <role>` |
| 3 | Task-pack apply | **Defer/gate** | Locks a new schema; brushes the deferred-import boundary; gate on proven dogfood pain |
| 6 | Evidence registry | **Split** | `evidence check` (verify a ref resolves) fits; `evidence capture --cmd` is scope creep (fid as executor) |
| 11 | Context budget metering | **Low-until-needed** | No value until packs exist and prove heavy; estimate-only |
| 10 | Formal intake | **Low** | Audit admits it's sugar over `in_inbox`; defer |
| 12 | Pattern mining | **Low, correctly placed** | v0.4+; agreed |

**Top recommendation, unchanged:** the journal-as-"what-changed-since" surface
(#1) is the single highest-leverage move; ship the durable diff before the
role-awareness.

### Resolution Priority Order

1. #1 Change review packs → 2. #2 Safe write (dry-run) → 3. #7 Git-aware
run/resume → 4. #5 Attention queue → 5. #8 Ready explain → 6. #4 Corrections →
7. #9 Handoff packs → 8. #3 Task-pack apply → 9. #6 Evidence registry →
10. #11 Context budget → 11. #10 Intake → 12. #12 Pattern mining.

### Resolved Items

#### Item #1 — Change review packs → `fid changes` (RESOLVED, v0.2c)

Decisions:
- **Distinct digest command** `fid changes --since <anchor>`, aggregated by
  entity (net field changes + newly-added child facts) — not a filtered
  timeline. Derived from `activity_log`; `fid activity list` stays the raw
  stream. (New glossary term **Change digest**.)
- **Anchors:** event ID + timestamp in v0.2, via an **extensible resolver**;
  review-id (v0.3) and commit (with run state) slot in additively without a
  contract change. `--since last` symbolic anchor deferred as sugar.
- **Output contract:** type-keyed extensible structure — `field_changes`
  (net `from→to`, intermediate transitions collapsed) + `added`/`removed`
  maps keyed by fact type (new fact types = new keys, no break). Net `from`
  derived from earliest post-anchor `activity_log` entry; current value = `to`.
- **Entity inclusion:** an entity appears iff it, or an append-only child it
  directly owns (ACs/verifications, completions, concerns, its own links),
  changed since the anchor. No implicit parent rollup.
- **Role shaping deferred** (role-agnostic MVP); when added, reuse context
  packs rather than a bespoke `--role`.
- **`fid diff <entity> --from --to` deferred**; single-entity case covered by
  `fid changes --since … --entity <id>`.

Specs changed: GLOSSARY (Change digest); PRD-OUT-019 + traceability;
06_DESIGN_PLAN_CLI §8 renamed "Activity and Change Commands" + full `fid
changes` spec; roadmap v0.2c row; TRD `internal/changes/` package.

#### Item #2 — Safe write protocol → `--dry-run` (RESOLVED, v0.2); idempotency deferred

Decisions:
- **Split.** Ship `--dry-run`; **defer idempotency keys** to a v0.5 conditional
  paired with the MCP evaluation (PRD-OPS-016, *Could*). No pre-work needed —
  additive later via schema versioning + a write-path hook.
- **Universal `--dry-run`** on every write command, built once at the write
  boundary (build event + validate + check preconditions, stop before append).
- **Output:** exit code **mirrors** the would-be outcome (0 apply / non-zero
  reject) so dry-runs compose with `&&`; body carries `dry_run`, `would_apply`,
  the `validation` result (same `rule`/`detail`/`entities` vocabulary as a real
  rejection), and an `event_preview` of effects. IDs minted only at append
  render as `<minted-on-append>`.
- **Placement:** assert the prepare/validate-before-append **seam in v0.1**
  (TRD §4.3); ship the flag in **v0.2** (feature follows its autonomous-agent
  consumers; keeps v0.1 lean).

Specs changed: PRD-OPS-015 (dry-run) + PRD-OPS-016 (idempotency, deferred) +
traceability; TRD §4.3 dry-run seam; GLOSSARY (Dry-run); 06_DESIGN_PLAN_CLI new
`## Conventions` preamble with text+JSON examples; roadmap v0.2c `--dry-run`
row + v0.5 conditional idempotency row.

#### Item #7 — Git-aware run/resume state (RESOLVED, v0.2c)

Decisions:
- **Hybrid git capture** (revised from pure recorder during grilling): fid
  **auto-captures** branch/HEAD/dirty read-only via a new `internal/gitctx/`,
  with **flag override**, **graceful degradation** outside a repo (null context
  + warning, never fail), and a self-describing **`git_source`**
  (`captured`|`provided`|`none`). fid never performs git operations. Rationale:
  auto-capture removes reliance on an agent remembering to pass git state —
  which is the whole point of moving run state out of the ephemeral chat
  ledger; and `dirty` is only meaningful when fid reads it live.
- **State lives on `agent_runs`** (extend, don't add a table or rename): new
  columns `status` (open|closed|abandoned), `branch`, `base_commit`,
  `head_commit`, `dirty`, `worktree_path`, `git_source`, `resumed_from_run_id`,
  `abandoned_reason_md_path`, `abandoned_at`. Additive v0.2 migration. Glossary
  term **Agent run** broadened to the unified execution record.
- **Lifecycle + commands:** `fid run list [--open|--abandoned|--by|--epic]`,
  `fid run show <run-id>`. `fid dev abandon <run-id> --reason-md --by` releases
  the claim + marks the run abandoned (row preserved; task left as-is, surfaced
  by `maintenance check`). `fid dev resume <task-id> --by [--force]` re-acquires
  the claim, opens a **new** run linked via `resumed_from_run_id` (clean
  per-session telemetry), carries forward branch/base, abandons any prior open
  run, returns resume **hints** only.
- **Capture points:** branch/base at `dev start`; head/dirty refreshed at
  `dev gate`/`dev finish`; dirty is "as of last recorded command."

Specs changed: PRD-OPS-017/018/019 + traceability; 03_TRD_SCHEMA §7
`agent_runs` columns; 06_DESIGN_PLAN_CLI `dev start`/`dev finish`/`dev resume`/
`dev abandon`/`run list`/`run show`; GLOSSARY (Agent run); roadmap v0.2c rows;
TRD `internal/gitctx/` module + dependency-policy note.

#### Item #5 — Human attention queue → `fid report attention` (RESOLVED, v0.2c)

Decisions:
- **A named report, not a top-level command.** `fid report attention` joins the
  `PRD-OUT-006` family (one more opinionated SQL wrapper), adding zero new
  top-level contract. A `fid attention` alias is a trivial additive promotion
  later if dogfooding shows the human types it constantly.
- **Human-default; audience-aware; role filter deferred.** v0.2 output is
  human-facing with no `--role` flag. Each reason group carries an internal
  `audience` (e.g. `needs_decision → [tl, human]`), so a later `--role`/`--for`
  filter is purely additive — never an enumerated `human|tl` freeze. Mirrors the
  #1 role-shaping deferral.
- **Open, version-gated reason set keyed by `reason`.** Reasons light up as their
  source entities ship — the same type-keyed extensibility as `fid changes`.
  v0.2: `needs_decision` (open escalation, no resolving decision),
  `needs_approval` (epic wrapped, awaiting human closure), `blocked`,
  `human_claimed`. v0.3 adds `blocking_finding_open`; v0.4 adds
  `validation_failed` / `validation_carried_forward` — additively, no contract
  change.
- **Three noise-control sharpenings.** (1) `blocked` is narrowed to *dead-end*
  blockers — an issue `blocked_by` a **cancelled** issue that will never clear
  without a human (normal blocked-by-in-progress is excluded as self-clearing;
  broader stall detection deferred). (2) `human_intervention_recorded` is
  **dropped** — past-tense, already covered by `fid stats interventions`; the
  queue is open-work-now, not an audit log. (3) `stale_escalation` is **folded
  into** `needs_decision` as an `age_hours` + `stale` attribute (48h default,
  hardcoded in v0.2, promotable to a `[reports]` config knob) rather than a
  double-counting separate group.
- **Defaults:** group order most-actionable first (`needs_decision` →
  `needs_approval` → `blocked` → `human_claimed`); within a group stale/oldest
  lead. **Empty is a healthy zero-exit state** (`note: nothing needs human
  attention`), never an error — a queue that cries wolf trains the human to
  ignore it.

Specs changed: PRD-OUT-020 + v0.2 traceability; GLOSSARY (Attention report);
06_DESIGN_PLAN_CLI §20 `fid report attention` (JSON + markdown + empty state) and
the `fid report list` example; roadmap v0.2c row. No schema/TRD change — it is a
hand-coded-SQL named report in the existing `internal/report/` package.

#### Item #8 — Ready-queue explain → `fid ready --explain` (RESOLVED, v0.2c)

Decisions:
- **A flag on `fid ready`, not a subcommand; default output frozen.** `--explain`
  (with `--limit`, `--include-blocked`) opts into a diagnostic shape; bare
  `fid ready` keeps its single-task contract verbatim, because the Workflow
  Runner depends on it (`PRD-ISSUE-008–012`).
- **Deterministic `order_factors`, not a `score`.** fid's ready queue is
  deterministic lexicographic ordering (sequence / priority / cross-epic rule),
  not a weighted engine. Each candidate exposes the actual sort keys the rule
  applied; explain **surfaces the existing rule and adds no ranking semantics**.
  A numeric score was rejected as fabricating tunable-weight semantics that are
  an explicit non-goal.
- **Typed skip-reason taxonomy** (extensible): `blocked_by`, `claimed`,
  `sequence_wait`, `not_actionable`, plus a `stale_claim` *flag* on candidates
  that are still ready but carry a past-TTL claim.
- **Noise control:** `claimed` / `sequence_wait` skips show under plain
  `--explain`; the potentially numerous `blocked_by` skips require
  `--include-blocked`. `--limit` (default 5) caps candidates and the blocked
  portion of `skipped`.
- **Empty queue is the prize.** `fid ready --explain --include-blocked` on an
  empty queue returns `selected: null` with the full skip list at **zero exit** —
  the canonical "why is nothing ready?" answer, specced as a first-class case.

Specs changed: PRD-ISSUE-024 + v0.2 traceability; 06_DESIGN_PLAN_CLI §4
(`fid ready --explain`/`--limit`/`--include-blocked`, diagnostic JSON + markdown,
empty-queue-explain); roadmap v0.2c row. No GLOSSARY change (blocked/claim/
sequence/stale-claim already defined; factors/skip reasons are output fields). No
schema/TRD change — it extends the existing ready-queue selection.

#### Item #4 — Corrections/undo → `fid event revert` taxonomy (RESOLVED, v0.2c)

Three of the audit's five proposed commands already ship: `fid event show` =
`fid journal inspect` (v0.1); `fid link remove` (v0.1); `fid event revert`
already roadmapped (v0.2c). The work was to lock the correction *mechanism* and
the *reversibility taxonomy*.

Decisions:
- **One generic correction primitive — no bespoke domain-correction verbs.**
  Corrections are compensating events via `fid event revert`, not git-revert and
  not new verbs. `fid ac retract` and `fid issue restore-status` are **not added**
  — calls below make them redundant. If dogfooding proves the generic path too
  clunky on a hot path, a sugar verb is a purely additive follow-up.
- **Two compensation modes, keyed by event type.** *Inverse* restores a
  mutable-state effect's prior value (transitions, field/label/link changes),
  guarded by an effect precondition that current state still equals what the
  event produced. *Retraction* withdraws an append-only audit fact (AC
  verification, concern, comment, triage verdict), excluding it from
  current-state derivation while **preserving its row**.
- **Transitions ARE revert-eligible (inverse-restore)** — this obviates
  `restore-status`. Reverting a transition returns the issue to a
  *previously-valid* prior status under the precondition guard, not an arbitrary
  workflow jump. **Revised the spec**: `06_DESIGN_PLAN_CLI` §16/§17 previously
  used `issue.transitioned` as the *not-eligible* example; flipped — transitions
  are now the eligible case and `issue.created` is the not-eligible exemplar.
- **Append-only facts revert via retraction** — this obviates `ac retract`. A
  *wrong-status* verification is still better fixed by re-verifying (latest-wins);
  retraction is reserved for "this fact should never have counted."
- **Creations, completions, claims, telemetry, and composites are NOT eligible**
  — they route to the forward domain path (cancel via transition; new attempt;
  claim/unclaim; supersede). Deleting a creation/completion would violate the
  append-only audit principle.
- **Eligibility is type-keyed and extensible**; `fid journal inspect` reports
  `revert.eligible` + `mode` (+ precondition) or the forward `correction`.

Specs changed: PRD-OPS-013 (inspect reports mode/precondition or guidance) and
PRD-OPS-014 (two-mode taxonomy, transition-eligible, generic-only); GLOSSARY
(Compensating event — inverse/retraction modes); 06_DESIGN_PLAN_CLI §16 revert
verdict + §17 `fid event revert` (inverse + retraction + not-eligible examples);
roadmap v0.2c `fid event revert` row. No schema change — retraction reuses the
append-only journal + projection; no new module.

#### Item #9 — Handoff/prompt packs → `fid handoff` (RESOLVED, v0.3d)

Resolves the configurability tension the open-item note flagged. The friction is
real (Phase Runner §2.4 #7 and §3.4 force hand-authored child prompts), but the
audit's hard-coded `fid handoff <role>` subcommands would freeze this project's
role set and §4 context-loading policy into fid.

Decisions:
- **A single generic `fid handoff` composite; role is a *parameter*, not a
  subcommand.** `fid handoff --role <role> --issue <id>` composes three
  *configurable* inputs into one read-only call: the role's context pack (v0.2c),
  a derived manifest, and a handoff template (v0.3). Adding a custom role → add
  its template → it works with no new command. Rejected pure compose-it-yourself
  (re-imposes the ceremony the feature removes) and hard-coded subcommands
  (freeze the workflow).
- **Manifest is derived, not declared.** `stop_conditions` and `expected_writes`
  are read off the configured workflow graph + role permissions, so they retrack
  when a project edits its workflow — this is what keeps Option A from smuggling
  hard-coded workflow back in. `required_reads` reuse the context pack;
  `blockers` reuse the #8 `fid ready --explain` skip data + open
  concerns/escalations; `commands` come from the template (fid doesn't own gate
  commands). Extensible/version-gated: `evidence_obligations` lights up v0.4.
- **Dual output, read-only.** JSON manifest by default (for the `workflow_runner`
  to consume) with the rendered prompt embedded as `prompt_md`; `--format
  markdown` renders the prose prompt. No journal event, no claim — generating a
  handoff is a view; claiming stays an explicit `fid dev start`.
- **Sequencing:** v0.3d, composed over v0.2c context packs and the v0.3c template
  engine. It is the on-ramp for fid's own `workflow_runner` to eventually replace
  the Phase Runner's hand-authored prompts at dogfood convergence.
- **Already-shipped overlap:** the audit's `--since` (QA reviewer window) and
  `--epic` (QA Lead final) are flags on the one command, not separate verbs.

Note: `09_AI_WORKFLOW.md` (§2.4 #7, §3.4) is the cited gap, but it is the current
markdown build-pipeline workflow and is intentionally **left untouched** — fid
product features are not threaded into it; fid-aware workflow/agent docs come
later for projects adopting fid.

Specs changed: PRD-OUT-021 + v0.3 traceability; GLOSSARY (Handoff pack);
06_DESIGN_PLAN_CLI §22 `fid handoff` (manifest JSON + rendered prompt); roadmap
v0.3d row; TRD `internal/handoff/` module (tree + boundary table). No schema
change — handoff composes existing read models + templates.

#### Item #3 — Task-pack apply (RESOLVED: deferred/gated, conditional v0.3d)

Disposition matches the assessor's Defer/gate tier — fid is unbuilt, so the
"dozens of fragile CLI calls" planning pain is asserted, not observed.

Decisions:
- **Defer to an evidence-gated conditional, not v0.2.** Recorded as a conditional
  roadmap item evaluated after v0.2 dogfooding; earliest build v0.3+. The input
  schema is NOT frozen now.
- **Gate:** build only if v0.1–v0.2 dogfooding shows TL planning ingestion is a
  measured friction or error source (multi-call planning sequences with real
  retry/partial-failure rate, or planning mistakes one reviewed transaction would
  prevent).
- **Import boundary (the load-bearing guard):** task-pack apply is
  structured-authoring **bulk-create** — structured input only (YAML/JSON/TOON,
  never free-prose markdown parsing), **create-only**, one journal event
  (`task_pack.applied`) with many effects, returning a source-ref→fid-id map. No
  round-trip, no sync, no markdown-canonical mode, no update-existing. The
  deferred markdown import (`fid import`/`init --import`, v0.5 conditional) stays
  separate.
- **Intent recorded, schema unfrozen:** `fid task-pack validate --from` +
  `apply --from --dry-run --by tl`; pack schema versioned and locked deliberately
  at build time.
- **Footprint kept honest:** PRD conditional line + glossary term + roadmap
  conditional only — NO CLI command spec and NO TRD module added (adding them
  would re-freeze the deferred schema).

Specs changed: PRD-OPS-020 (Could, conditional; not in committed traceability);
GLOSSARY (Task pack, disambiguated from [[Handoff pack]] and markdown import);
roadmap v0.3d conditional row. No CLI/TRD/schema change by design.

#### Item #6 — Evidence registry → `fid evidence check` (RESOLVED, v0.4a)

Adopts the assessor's split, sharpened to drop the registry table.

Decisions:
- **Accept read-only `fid evidence check`; reject `fid evidence capture --cmd`.**
  fid is a verifier, never an executor: making it run `go test` would break
  local-first / no-side-effects (PRD-NFR-004/005/006) and the single-binary
  model, and duplicate the gate mechanism. Capture is already met by the agent
  recording evidence (`fid dev gate`, evidence ref/excerpt on `ac verify`).
- **No evidence registry table (sharpening).** Evidence stays embedded
  (ac_verifications, qa_findings, val_assertion_outcomes); check is a derived
  read. This also dissolves `fid evidence show <evidence-id>` (no evidence-id) —
  a scoped check is the per-evidence view.
- **Type-keyed, extensible resolver.** `commit` (local git) and file/source/
  screenshot (filesystem) are checkable → `resolved`/`missing`; `command`,
  `manual`, `url` are `unverifiable`. New types added without a contract change.
- **Two boundary calls.** (1) *Resolution, not re-proof* — verifies the ref
  exists, never re-runs a test or re-proves a claim. (2) *Existence, not content
  drift* — no hash comparison in v1 (would need capture-time hashes stored;
  deferred, additive).
- **Standalone, not folded into `fid maintenance check`** (QA evidence integrity
  vs structural DB integrity — different audiences); `--warn-as-error` for CI /
  wrap-up gating. v0.4 (covers AC v0.1 + finding v0.3 + outcome v0.4 evidence in
  one pass).

Specs changed: PRD-QA-009 + v0.4 traceability; GLOSSARY (Evidence check);
06_DESIGN_PLAN_CLI §23 `fid evidence check`; roadmap v0.4a row; TRD
`internal/evidence/` module (tree + table, execution/hashing forbidden). No
schema change.

#### Item #11 — Context budget metering (RESOLVED: deferred/gated, estimate-only, conditional v0.5)

Tier Low-until-needed. Disposition mirrors #3.

Decisions:
- **Defer to a gated conditional.** Context packs ship v0.2c, but default context
  is deliberately small (workflow principle #4), so "packs heavy enough to
  overload a session" is unproven. Gate: build only if dogfooding shows packs
  overload Developer/QA/TL-Lite sessions.
- **Estimate-only; reject `--budget-tokens` auto-trim.** Metering reports size; it
  never omits pack sections to fit a budget. Auto-trim would make
  `--context-pack <role>` lossy and non-deterministic, undermining packs as a
  stable contract. Over-budget is the operator's call, not fid's.
- **Heuristic proxy, not a tokenizer.** A model-agnostic chars-per-token
  approximation reported as an estimate; fid bundles no per-model tokenizers
  (multi-provider; single static binary).
- **Intent recorded, shape unfrozen:** `fid context estimate --pack <role>
  --issue <id>`, tokens by section.
- **Footprint:** PRD conditional line + roadmap v0.5 conditional only. No
  GLOSSARY/CLI/TRD.

Specs changed: PRD-OUT-022 (Could, conditional; not in committed traceability);
roadmap v0.5 conditional row.

#### Item #10 — Formal intake (RESOLVED: verbs deferred, promote primitive closed v0.2c)

Tier Low. The audit calls intake "sugar over `in_inbox`" — but grilling that claim
exposed that the *promote* half of the primitive was missing.

Decisions:
- **Defer the `fid intake` verbs** (`add`/`promote`/`close`) — premature sugar.
- **Close the real gap so the deferral is honest.** Park (`create --inbox`) and
  list (`list --inbox`) existed, but there was **no way to assign an inbox task to
  an epic or clear the inbox flag** — `fid issue update` (PRD-ISSUE-022) omits
  parent/epic and inbox. Without that, the v0.1 inbox is check-in-only and "it's
  just sugar" is false.
- **Minimal primitive (v0.2):** extend `fid issue update` with `--epic`
  (assign/change/clear a *task's* epic parent — null↔epic only, not subtask
  reparenting; `--sequence` required for a sequential target; kind-hierarchy
  rules still apply) and `--inbox/--no-inbox`. Now park → list → promote all work
  on base commands, and the `fid intake` verbs are correctly deferred ergonomics.
- v0.1's frozen `issue update` field set is left untouched (the addition lands
  v0.2).

Specs changed: PRD-ISSUE-025 + v0.2 traceability; 06_DESIGN_PLAN_CLI
`fid issue update` (`--epic`/`--inbox` flags + promote example); roadmap v0.2c
row. No GLOSSARY change (Inbox term stands); no TRD/schema change (lives in
`internal/issues/`).

#### Item #12 — Pattern mining → recurring-pattern reports (RESOLVED, v0.4b)

Tier Low, placement agreed (v0.4+). Accepted as deterministic pattern surfacing.

Decisions:
- **Deterministic grouping only (v0.4b).** `fid report recurring-findings`,
  `recurring-concerns`, `recurring-patterns` group findings/concerns/gate
  runs/triage/VAL-outcome failures by structured dimensions (trigger class,
  label, area, gate command, finding severity, file-path prefix, repeated
  `VAL-NNN` failures) and rank recurring groups as candidates with counts +
  exemplars. No ML, embeddings, or semantic theme extraction.
- **Surface, don't author (the load-bearing boundary).** The audit's
  `suggest-principles` implied fid generating principle/decision prose — which
  would make fid an AI service, breaking local-first / no-network / no-executor.
  Reframed: fid ranks structured patterns; a human/TL turns a pattern into a
  DEC-NNN or principle. Same analyzer-not-author line as evidence-check and
  context-estimate.
- **Report family, not a new `fid insights` namespace.** Deterministic reads
  belong with `fid report` (where `attention` landed). Read-only; no schema;
  extends the v0.4b computed retro signals.

Specs changed: PRD-OUT-023 + v0.4 traceability; GLOSSARY (Pattern report);
06_DESIGN_PLAN_CLI §20 (recurring reports + `report list` entries); roadmap v0.4b
row. No TRD module — `internal/report/`, like `attention`.

### Open Items

None — all twelve audit items have been walked and resolved. Resolved blocks
above: #1, #2, #7 (2026-05-31 session); #5, #8, #4, #9, #3, #6, #11, #10, #12
(2026-06-01 session, in the resolution priority order).

