# fid Spec Suite - Independent Second-Pass Audit

**Date:** 2026-05-28
**Scope:** All files under `project-docs/specs/`, compared with historical
input in `scaffold/initial-spec.md`.
**Prior audit reviewed:** `project-docs/docs/reviews/2026-05-27-spec-audit.md`
**Lens:** Watertightness plus in-scope feature gaps for an agent-first,
journal-canonical CLI issue tracker.

## Audit Notes

The first audit resolution pass materially improved the suite: schema DDL,
config, CLI flag/output detail, end-to-end flows, whet mapping, and roadmap
content are now in formal spec files. This pass does not re-raise findings
marked `DONE` in the prior audit. It also does not re-litigate intentionally
deferred items #18 (named reports), #19 (traceability table fill-in), #20
(quality-gate command validation), or #22 (MCP design).

The findings below focus on places where the current formal suite would still
let two implementation agents build incompatible behavior, or where a
predictable mid-loop user need is not expressible.

---

## Watertightness Findings

### P0 - Same-checkout concurrency is not actually specified

**Citation:** `project-docs/specs/01_PRD.md` section 4, PRD-NFR-008 (line 275);
`project-docs/specs/02_TRD.md` section 4.3, write path (lines 198-205) and
concurrency (lines 228-232).

**Problem:** PRD-NFR-008 requires two agents in one checkout to serialize
writes without data loss, but the TRD says the journal append is the commit
point and separately claims SQLite default locking makes single-worktree
writes safe. SQLite locking cannot serialize the pre-DB journal append unless
the spec defines a process lock or a transaction held before append. Two
processes can both validate against the same DB state and append conflicting
journal events unless the lock boundary is explicit.

**Proposed fix:** Add a v0.1 requirement and architecture clause: acquire an
exclusive `.fid/lock` or `BEGIN IMMEDIATE` projection lock before validation;
hold it through validation, journal append+fsync, and DB projection; release
only after projection succeeds or fails safely. Add an integration test with
two OS processes or goroutines racing the same write target.

### P0 - Prior #5 resolved the hidden initial-spec dependency incompletely: default workflows still live only in historical input

**Citation:** `project-docs/specs/01_PRD.md` section 3.5, PRD-WF-005 through
PRD-WF-010 (lines 120-126); `project-docs/specs/10_IMPLEMENTATION_ROADMAP.md`
section v0.1b, "Hard-coded workflows" (line 77);
`project-docs/specs/03_TRD_SCHEMA.md` section 2, `issues.workflow_id` and `status`
comments (lines 90-94); historical source `scaffold/initial-spec.md` section 7.2-7.4
(lines 1530-1587).

**Problem:** The prior audit #5 moved much of the hidden initial-spec content
into formal docs, but the exact v0.1 state machines remain outside the source
of truth. PRD requires default epic/task/subtask workflows and status-category
mappings, the roadmap points to "section 7.2-7.4", and schema comments say status is
CHECK-constrained per workflow. Those sections only exist in
`scaffold/initial-spec.md`, which is explicitly historical input here. An
implementation agent reading only formal specs cannot know whether task states
include `qa_review`, `qa_block`, `failed_gate`, `qa_pass`, etc., or which rules
attach to each transition.

**Proposed fix:** Move the exact task, epic, and subtask workflow state lists,
category mappings, allowed transitions, and transition requirements from
`scaffold/initial-spec.md` section 7.2-7.4 into `02_TRD.md` or a new formal workflow
spec. Update the roadmap references away from naked `section 7.2-7.4` references.

### P1 - Journal event contract is too loose for the durable source of truth

**Citation:** `project-docs/specs/02_TRD.md` section 4.3, write path (lines 198-202);
`project-docs/specs/02_TRD.md` section 5, journal contract row (line 304);
`project-docs/specs/03_TRD_SCHEMA.md` section 9.4, `activity_log.journal_event_id`
(line 498).

**Problem:** The journal is the source of truth and a compatibility contract,
but the formal spec only names an envelope and generic effects. It does not
define `event_id`, actor shape, command metadata shape, effect schema,
update/delete patch semantics, preconditions, checksums, or how
`activity_log.journal_event_id` is derived. Without this, replay, merge,
activity projection, golden fixtures, and upgrader tests can diverge.

**Proposed fix:** Add a JSON schema and at least one full v1 event example to
`02_TRD.md` or `03_TRD_SCHEMA.md`: `event_id`, `event_schema`, `event_type`,
timestamp, actor role/actor, sanitized command metadata, ordered effects,
primary-key representation, update preconditions/prior values, and
activity-log effect conventions.

### P1 - Git merge conflict behavior is overclaimed

**Citation:** `project-docs/specs/02_TRD.md` section 4.3, concurrency (lines
230-234); `project-docs/specs/01_PRD.md` section 4, PRD-NFR-007 (line 274);
`project-docs/specs/01_PRD.md` section 3.14, PRD-OPS-011 (line 260).

**Problem:** The TRD says logically conflicting writes surface as JSONL
conflicts on the same lines. In normal git merges, two branches appending
different valid lines often merge cleanly while still expressing conflicting
status, claim, or field updates. `fid merge --journal` conflict detection is
deferred to v0.2, leaving v0.1 replay behavior ambiguous after a clean textual
merge that is logically invalid.

**Proposed fix:** Add event preconditions and replay conflict handling in v0.1,
or explicitly state that v0.1 applies merged events in deterministic order and
`fid maintenance check` reports logical conflicts. If preconditions are chosen,
define them in the journal event contract.

### P1 - Projection freshness has no metadata home

**Citation:** `project-docs/specs/01_PRD.md` section 3.14, PRD-OPS-003 and
PRD-OPS-004 (lines 252-253); `project-docs/specs/02_TRD.md` section 4.3, rebuild
(lines 248-251); `project-docs/specs/03_TRD_SCHEMA.md` section 9, supporting tables
(lines 453-540).

**Problem:** The DB must detect whether it is missing or stale before reads,
but the schema has no projection metadata table recording last journal file,
byte offset, event ID/hash, projection schema version, or replay status. The
spec says "checked against the journal tail" but not how.

**Proposed fix:** Add `projection_meta` or equivalent to `03_TRD_SCHEMA.md`.
Define the exact freshness comparison and recovery behavior: full rebuild when
metadata is absent or inconsistent; catch-up when the last applied event/offset
matches the journal prefix; stop on mismatch.

### P1 - Permissions matrix is not complete enough to enforce PRD-ROLE

**Citation:** `project-docs/specs/01_PRD.md` section 3.7, PRD-ROLE-001 through
PRD-ROLE-005 (lines 149-153); `project-docs/specs/02_TRD.md` section 10,
TRD-ARCH-010 (line 394); `project-docs/specs/03_TRD_SCHEMA.md` section 13 (lines
729-755); `project-docs/specs/10_IMPLEMENTATION_ROADMAP.md` section v0.1d,
"Permissions matrix" (line 102).

**Problem:** TRD-ARCH-010 says every CLI write action must have a matrix entry,
but the default matrix omits many v0.1 writes: issue update, issue complete,
claim/unclaim, AC add, link add, label add, run record, developer composite
commands, and maintenance writes. The roadmap and PRD require a command to show
the matrix (`fid policy show`), but the CLI design does not specify it.

**Proposed fix:** Expand the matrix to every write action, including
composites, or explicitly define wildcard/action-family semantics. Add
`fid policy show [--format json|toon|markdown|yaml]` to
`05_DESIGN_PLAN.md` and `06_DESIGN_PLAN_CLI.md`.

### P1 - `dev gate` has nowhere to persist before completion exists

**Citation:** `project-docs/specs/06_DESIGN_PLAN_CLI.md` section 5, `fid dev start`,
`fid dev gate`, and `fid dev finish` (lines 311-378);
`project-docs/specs/03_TRD_SCHEMA.md` section 2.3, `task_completions` (lines
209-226); `project-docs/specs/GLOSSARY.md` "Gate" (lines 101-108).

**Problem:** `fid dev gate` records gate output during an attempt before
`fid dev finish` creates the `task_completions` row where `gate_results` lives.
There is no draft attempt row, gate table, or active attempt concept. This
makes the command impossible to implement without inventing storage behavior.

**Proposed fix:** Choose one model:

- Create a draft `task_completions` row at `dev start` and finalize it at
  `dev finish`.
- Add a `gate_runs` table keyed by issue/run/attempt and roll it into the
  completion output.
- Remove standalone `dev gate` until finish-time gate capture is specified.

### P2 - Mandatory `--by` conflicts with write command specs

**Citation:** `project-docs/specs/01_PRD.md` section 3.7, PRD-ROLE-001 (line 149);
`project-docs/specs/03_TRD_SCHEMA.md` section 13 final note (line 755);
`project-docs/specs/06_DESIGN_PLAN_CLI.md` section 2, `fid issue unclaim` (lines
194-205); `project-docs/specs/05_DESIGN_PLAN.md` section 5.4 claim error example
(lines 207-214).

**Problem:** `fid issue unclaim` is a state-changing command but lacks
`--by`; even the force-acquire corrective hint omits actor identity. That
breaks auditability, permissions, and the global "every write command"
requirement.

**Proposed fix:** Change the command to
`fid issue unclaim <id> --by <role> [--force]`; update examples and errors.
If force-acquire is really a claim operation rather than an unclaim operation,
rename it to `fid issue claim <id> --by <role> --force`.

### P2 - Human intervention requirement has no CLI path

**Citation:** `project-docs/specs/01_PRD.md` section 3.4, PRD-DEV-009 (line 110);
`project-docs/specs/03_TRD_SCHEMA.md` section 2.3, `task_completions.human_intervention`
(line 224); `project-docs/specs/06_DESIGN_PLAN_CLI.md` section 2 `fid issue complete`
(lines 155-175) and section 5 `fid dev finish` (lines 354-371).

**Problem:** PRD-DEV-009 and the schema support a `human_intervention` flag,
but neither `fid issue complete` nor `fid dev finish` exposes a flag or parsing
rule to set it. Implementers may ignore the field or invent incompatible
behavior.

**Proposed fix:** Add `--human-intervention` and optional
`--intervention-note` to completion commands, or specify a template section that
`fid dev finish` parses into this field.

### P2 - Stats versioning contradicts PRD and roadmap

**Citation:** `project-docs/specs/01_PRD.md` section 3.12, PRD-TELEM-004 through
PRD-TELEM-006 (lines 221-222); `project-docs/specs/10_IMPLEMENTATION_ROADMAP.md`
section v0.2c, stats expansion (line 177); `project-docs/specs/06_DESIGN_PLAN_CLI.md`
section 11, `fid stats tasks` and `fid stats attempts` (lines 628-657).

**Problem:** `stats tasks` and `stats attempts` are fully specified in the CLI
doc without a v0.2 marker, but PRD-TELEM-005/006 and the roadmap defer those
aggregations to v0.2. A v0.1 implementation agent could reasonably implement
them early or stub them, depending on which doc it privileges.

**Proposed fix:** Mark `fid stats tasks` and `fid stats attempts` as v0.2 in
`06_DESIGN_PLAN_CLI.md`, or move those metrics into v0.1 requirements and
roadmap.

### P2 - MCP violates the absolute "no network listener" requirement

**Citation:** `project-docs/specs/01_PRD.md` section 3.14, PRD-OPS-010 (line 259);
`project-docs/specs/01_PRD.md` section 4, PRD-NFR-006 (line 273);
`project-docs/specs/02_TRD.md` section 6, TRD-SEC-001 and TRD-SEC-005 (lines
317-321); `project-docs/specs/10_IMPLEMENTATION_ROADMAP.md` section v0.3d, MCP
server (line 231).

**Problem:** PRD-NFR-006 says the system operates without a network listener,
but TRD and roadmap require `fid mcp start` on localhost. This is distinct
from prior deferred issue #22; the defect is the current NFR/security
contradiction, not the detailed MCP design.

**Proposed fix:** Change PRD-NFR-006 to "no listener except explicit
localhost-only `fid mcp start`" and cross-reference TRD-SEC-005, or specify MCP
over stdio rather than a socket.

### P3 - Minor flag/name drift will trip implementation

**Citation:** `project-docs/specs/06_DESIGN_PLAN_CLI.md` section 2,
`fid issue create` flag table (line 66); `project-docs/specs/10_IMPLEMENTATION_ROADMAP.md`
section v0.1b, `fid issue create` (line 70); `project-docs/specs/06_DESIGN_PLAN_CLI.md`
section 8, label commands (lines 516-542).

**Problem:** CLI design uses `--inbox`; roadmap says `--in-inbox`. Label
commands can create and list labels, while issue label assignment is only
implicit through issue create/update replacement. There is no explicit
add/remove lifecycle even though labels are a daily-use issue-triage surface.

**Proposed fix:** Standardize on `--inbox`. Define label assignment semantics:
either `fid issue update --labels` replaces all labels, or add
`fid issue label add/remove` commands with journaled activity-log projection.

---

## Feature-Gap Findings

### P2 - In-scope gap: no entity history command despite activity-log model

**Citation:** `project-docs/specs/GLOSSARY.md` "Activity log" (lines 137-142);
`project-docs/specs/03_TRD_SCHEMA.md` section 9.4, `activity_log` (lines 496-504);
`project-docs/specs/10_IMPLEMENTATION_ROADMAP.md` section v0.1e, "Activity log"
(line 121).

**Problem:** The schema and roadmap imply `fid activity list`, but PRD and CLI
design do not specify it. Agents and humans will repeatedly need to answer
"what happened to this task?" without reading raw journal events.

**Proposed fix:** Add a PRD-OUT or PRD-OPS requirement and CLI:
`fid activity list [--on issue/FID-x] [--since <date>] [--event <name>] [--format ...]`.
Define JSON output with event, actor, occurred_at, ref, and details.

### P2 - In-scope gap: no claim renewal

**Citation:** `project-docs/specs/01_PRD.md` section 3.6, claim requirements (lines
132-140); `project-docs/specs/03_TRD_SCHEMA.md` section 11, claim semantics (lines
584-626); `project-docs/specs/06_DESIGN_PLAN_CLI.md` section 2, claim/unclaim
commands (lines 177-205).

**Problem:** Claims default to 30 minutes, but long-running agent tasks need a
way to extend a live claim without releasing/reacquiring or force-stealing
their own work. Stale-claim recovery is specified; proactive renewal is not.

**Proposed fix:** Add `fid issue claim <id> --by <role> --renew [--ttl <duration>]`
or `fid claim renew <type>/<id> --by <role> [--ttl <duration>]`. Record renewal
in the journal and activity log.

### P2 - In-scope gap: no compensating correction or undo story

**Citation:** `project-docs/specs/02_TRD.md` section 4.3, malformed journal handling
(lines 253-256); `project-docs/specs/02_TRD.md` section 8, rollback process
(lines 356-357).

**Problem:** Recovery is manual journal repair or git revert. For an
append-only audit system, a bad-but-valid event needs a first-class correction
path that preserves history. Otherwise routine mistakes such as wrong label,
wrong link, wrong status, or bad AC verification become manual surgery.

**Proposed fix:** Add `fid event revert <event-id> --by human|tl --reason <text>`
that writes a compensating event where possible, plus `fid journal inspect`
for locating event IDs. Scope v0.1 to reversible event types if necessary.

### P3 - In-scope gap: link lifecycle lacks unlink

**Citation:** `project-docs/specs/01_PRD.md` section 3.8, PRD-LINK-001 through
PRD-LINK-008 (lines 162-167); `project-docs/specs/06_DESIGN_PLAN_CLI.md` section 7
(lines 464-512).

**Problem:** Links can be added and listed, but not removed or superseded.
Bad dependency and traceability links are predictable during planning and
task triage. Without unlink, users must either leave wrong links or perform
manual correction.

**Proposed fix:** Add `fid link remove <from-id> <relation> <to-id> --by <role>`
and define whether it removes both inverse rows. Project removal to
`activity_log`.

### P3 - In-scope gap: telemetry is too coarse for billing analysis

**Citation:** `project-docs/specs/01_PRD.md` section 3.12, telemetry requirements
(lines 217-222); `project-docs/specs/03_TRD_SCHEMA.md` section 7, `agent_runs`
(lines 405-423); `project-docs/specs/06_DESIGN_PLAN_CLI.md` section 10,
`fid run record` (lines 578-604).

**Problem:** Only input/output tokens and computed list-price cost are modeled.
Modern agent runs often need cached tokens, reasoning tokens, tool-call counts,
provider request IDs, and provider-billed cost to reconcile actual spend and
debug cost anomalies.

**Proposed fix:** Add optional `token_breakdown_json`, `tool_call_count`,
`provider_run_id`, and `provider_cost_usd` fields. Keep current columns as the
normalized summary used by v0.1 stats.

### P3 - Later-version gap: search/query remains conditional and underspecified

**Citation:** `project-docs/specs/10_IMPLEMENTATION_ROADMAP.md` section v0.3d,
conditional `fid query` DSL (line 232); historical source
`scaffold/initial-spec.md` section 9.5, example queries (lines 1950-1952).

**Problem:** SQL exists, named reports are deferred, and `fid query` is
conditional. A human or agent mid-loop will predictably need simple search
across titles, summaries, labels, and markdown paths without writing SQL.

**Proposed fix:** Keep SQL as the escape hatch, but add a later-version
`fid search <term> [--kind --status --label --epic]`, or make `fid query` a
tracked PRD item once report gaps are known from dogfooding.

---

## Overall Assessment

The suite is much stronger after the first resolution pass, but it is not yet
implementation-ready for Phase 1. The top three fixes before phase planning
are: specify the v0.1 workflow graphs in the formal specs; nail down journal
locking, event schema, and replay/precondition semantics; and complete the
v0.1 permission/CLI traceability gaps (`--by`, `policy show`, gate
persistence, activity history). Once those are resolved, the remaining issues
are mostly lifecycle polish and version-boundary cleanup rather than blockers.

