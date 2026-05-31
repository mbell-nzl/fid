# PRODUCT REQUIREMENTS DOCUMENT - fid

> This document defines product behavior in testable language. It should say
> what the system must do and why, not how it is implemented.

## 1. Requirement Style

All requirements use EARS-style language:

- Ubiquitous: "The system shall ..."
- Event-driven: "When <event>, the system shall ..."
- State-driven: "While <state>, the system shall ..."
- Optional feature: "Where <feature>, the system shall ..."
- Unwanted behavior: "If <error condition>, the system shall ..."

Stable requirement IDs use the form `PRD-<AREA>-<NNN>`:

- `INIT` - project initialisation and setup
- `ISSUE` - issue lifecycle (create, read, list, scheduling, ready queue)
- `AC` - acceptance criteria and verification
- `DEV` - developer workflow (completions, concerns, gates)
- `WF` - workflow transitions and state enforcement
- `CLAIM` - exclusive-access leases
- `ROLE` - role-based permissions
- `LINK` - cross-entity traceability links
- `QA` - QA reviews, findings, and plans
- `VAL` - validation contracts and assertion outcomes
- `GOV` - governance (decisions, escalations, triage, requirements,
  carry-forward, retrospectives)
- `TELEM` - agent telemetry and cost tracking
- `OUT` - output formats, export, reports, and templates
- `OPS` - operations and maintenance

Priority uses MoSCoW: Must / Should / Could / Won't. Target version noted
where a requirement is deferred past v0.1.

## 2. Personas / User Jobs

| Persona | Job To Be Done | Success Signal |
|---|---|---|
| Developer agent | Implement a task following a spec, run quality gates, verify ACs, raise concerns | Task transitions through `dev_complete` with all ACs verified and concerns recorded |
| QA Reviewer agent | Review completed tasks mid-epic, surface code quality issues | Blocking findings raised before bad patterns compound across tasks |
| QA Lead agent | Design verification plans, run final epic reviews, generate human checklists | Epic quality verified with evidence before wrap-up |
| TL-Lite agent | Triage developer concerns and soft blockers with a single verdict | RESUME/BLOCK verdict delivered without invoking full TL |
| TL agent | Plan epics, create task specs, wrap up completed epics, record decisions | Epics scoped, reviewed, and closed with full audit trail |
| Workflow Runner agent | Drive the task loop: pick next task, spawn agents, dispatch on results | Next ready task identified, agents spawned, concerns dispatched, QA triggered on schedule |
| Human product owner | Approve epic closure, review human checklist, make product decisions, override role boundaries | Project state inspectable at any time, quality verified, decisions recorded durably |

## 3. Functional Requirements

### 3.1 Initialisation (INIT)

| ID | Requirement | Priority | Notes |
|---|---|---|---|
| PRD-INIT-001 | The system shall initialise a project-scoped workspace in the current directory. | Must | v0.1. Creates `.fid/` directory. |
| PRD-INIT-002 | When initialising, the system shall accept a configurable issue ID prefix. | Must | v0.1. Defaults to project directory name. |
| PRD-INIT-003 | When initialising, the system shall seed the project with the default role profile. | Must | v0.1. Seven roles: developer, qa_reviewer, qa_lead, tl_lite, tl, workflow_runner, human. |
| PRD-INIT-004 | If the current directory already contains a fid workspace, the system shall reject re-initialisation with a clear error. | Must | v0.1 |

### 3.2 Issue Lifecycle (ISSUE)

| ID | Requirement | Priority | Notes |
|---|---|---|---|
| PRD-ISSUE-001 | The system shall support three issue kinds: epic, task, and subtask. | Must | v0.1 |
| PRD-ISSUE-002 | When creating an issue, the system shall generate a unique ID composed of the project prefix and a hash suffix. | Must | v0.1 |
| PRD-ISSUE-003 | The system shall enforce that epics have no parent, tasks may optionally belong to one epic, and subtasks belong to exactly one task. | Must | v0.1 |
| PRD-ISSUE-004 | If an issue creation would violate the kind hierarchy (e.g. subtask without a task parent, task with a subtask parent), the system shall reject it with an error identifying the constraint. | Must | v0.1 |
| PRD-ISSUE-005 | The system shall not allow an issue's kind to be changed after creation. | Must | v0.1 |
| PRD-ISSUE-006 | When queried by ID, the system shall return the issue with its acceptance criteria, latest verification state per role, open concern count, and latest triage status. | Must | v0.1 |
| PRD-ISSUE-007 | The system shall list issues filterable by kind, status category, epic, label, priority, and inbox state. | Must | v0.1 |
| PRD-ISSUE-008 | The system shall support a ready queue that returns the next actionable task. | Must | v0.1 |
| PRD-ISSUE-009 | While an epic uses sequential scheduling, the ready queue shall return only the lowest-sequence unfinished task, and only when that task is unblocked. | Must | v0.1 |
| PRD-ISSUE-010 | While an epic uses sequential scheduling, if the lowest-sequence unfinished task is blocked, the ready queue shall return no task for that epic and report the blocker. | Must | v0.1 |
| PRD-ISSUE-011 | While an epic uses prioritised scheduling, the ready queue shall return the highest-priority unblocked task within the epic. | Must | v0.1 |
| PRD-ISSUE-012 | When the ready queue is invoked without an epic filter, the system shall return the top task across all prioritised epics and unattached tasks, with sequential epics contributing only their current task when unblocked. | Must | v0.1 |
| PRD-ISSUE-013 | When creating a task under a sequential epic, the system shall require a positive integer sequence value. | Must | v0.1 |
| PRD-ISSUE-014 | If a task sequence value would duplicate an existing sibling's sequence under the same epic, the system shall reject the creation. | Must | v0.1 |
| PRD-ISSUE-015 | The system shall support optional labels, priority (0–5, where 0 is critical), estimates (in tokens), and due dates on issues. | Should | v0.1 |
| PRD-ISSUE-016 | The system shall support an inbox boolean on issues to mark unsorted incoming items awaiting assignment. | Should | v0.1 |
| PRD-ISSUE-017 | The system shall support optional cycle assignment on issues. | Could | v0.5 |
| PRD-ISSUE-018 | When reading an issue, the system shall support on-demand section parsing from the referenced spec markdown file via a `--section` flag. | Must | v0.1 |
| PRD-ISSUE-019 | When creating an issue, the system shall assign a default owner role: `tl` for epics, `developer` for tasks. The owner role shall be overridable at creation time. | Should | v0.1 |
| PRD-ISSUE-020 | The system shall support an optional owner actor field on issues to distinguish specific actors within a role. | Should | v0.1 |
| PRD-ISSUE-021 | The system shall allow owner role and owner actor to be changed via an explicit update command. Owner changes shall be journaled and projected to the activity log. | Should | v0.1 |
| PRD-ISSUE-022 | The system shall allow an issue's title, priority, labels, spec markdown path, token estimate, and due date to be changed via an explicit update command. Field updates shall be journaled and projected to the activity log. | Should | v0.1 |
| PRD-ISSUE-023 | The system shall allow labels to be added to or removed from an existing issue without replacing the full label set. Label changes shall be journaled and projected to the activity log. | Should | v0.1 |

### 3.3 Acceptance Criteria (AC)

| ID | Requirement | Priority | Notes |
|---|---|---|---|
| PRD-AC-001 | The system shall support ordered acceptance criteria on tasks and subtasks. | Must | v0.1 |
| PRD-AC-002 | If acceptance criteria are added to an epic, the system shall reject the operation. | Must | v0.1 |
| PRD-AC-003 | The system shall support multi-role AC verification, allowing independent outcomes from different verifier roles on the same AC. | Must | v0.1 |
| PRD-AC-004 | When an AC is verified, the system shall record the verifier role, status (pass/fail/blocked/not_applicable), evidence type, evidence reference, and timestamp as an append-only record. | Must | v0.1 |
| PRD-AC-005 | The system shall preserve all verification history; a failed verification followed by a passing one from the same role results in two records. | Must | v0.1 |
| PRD-AC-006 | When listing ACs for an issue, the system shall include the latest verification status per verifier role. | Must | v0.1 |

### 3.4 Developer Workflow (DEV)

| ID | Requirement | Priority | Notes |
|---|---|---|---|
| PRD-DEV-001 | When a developer completes a task attempt, the system shall create an append-only completion record with attempt number, tokens used, commit hashes, files touched, and gate results. | Must | v0.1 |
| PRD-DEV-002 | The system shall support multiple completion attempts per task, each as a separate numbered record. | Must | v0.1 |
| PRD-DEV-003 | When a previous attempt exists on a task, the system shall mark it as superseded when a new attempt is recorded. | Must | v0.1 |
| PRD-DEV-004 | The system shall record gate results on completion records as structured data: command, exit code, timestamp, and output excerpt. | Must | v0.1 |
| PRD-DEV-005 | When a developer raises a concern during a task attempt, the system shall create a structured concern record linked to both the issue and the completion attempt. | Must | v0.1 |
| PRD-DEV-006 | The system shall track concern lifecycle through statuses: open, addressed, escalated, wontfix. | Must | v0.1 |
| PRD-DEV-007 | When an issue has one or more open concerns, the system shall surface this as a derived concern flag on issue reads. | Must | v0.1 |
| PRD-DEV-008 | The system shall support composite developer commands that combine claim acquisition, status transition, agent run recording, and completion creation into a single atomic operation. | Should | v0.1 |
| PRD-DEV-009 | When a developer records a human intervention during a task attempt, the system shall flag it on the completion record and may store a short intervention note. | Should | v0.1 |
| PRD-DEV-010 | The system shall allow gate runs to be recorded before a completion exists. Gate runs shall be append-only and shall be rolled into the completion record for the same issue and attempt at finish time. | Must | v0.1 |

### 3.5 Workflow Transitions (WF)

| ID | Requirement | Priority | Notes |
|---|---|---|---|
| PRD-WF-001 | The system shall enforce that issue status transitions follow the assigned workflow's allowed transition edges. | Must | v0.1 |
| PRD-WF-002 | If a requested status transition is not an allowed edge in the workflow, the system shall reject it with an error naming the current state and the requested state. | Must | v0.1 |
| PRD-WF-003 | When a transition has validation rules, the system shall evaluate all rules before allowing the transition. | Must | v0.1 |
| PRD-WF-004 | If a transition rule fails, the system shall report which specific conditions are unmet, including entity IDs where applicable (e.g. "ACs not yet verified by developer: AC-a1b2c3, AC-7c3e9f"). | Must | v0.1 |
| PRD-WF-005 | The system shall ship default workflows for epic, task, and subtask kinds. | Must | v0.1 |
| PRD-WF-005a | The system shall ship a default workflow for the escalation kind. | Must | v0.2 |
| PRD-WF-006 | The system shall map every workflow state to a status category (unstarted, started, completed, cancelled) for portable queries. | Must | v0.1 |
| PRD-WF-007 | The system shall support the `all_acs_verified` transition rule, parameterised by a set of required verifier roles. | Must | v0.1 |
| PRD-WF-008 | The system shall support the `has_completion_report` transition rule. | Must | v0.1 |
| PRD-WF-009 | The system shall support the `concerns_addressed` transition rule, satisfied when no open concerns exist on the issue. | Must | v0.1 |
| PRD-WF-010 | The system shall support the `all_child_tasks_done_or_cancelled` transition rule for epics. | Must | v0.1 |
| PRD-WF-011 | Where the configurable workflow engine is enabled, the system shall allow users to define custom workflows with states, transitions, and per-transition rules. | Should | v0.2 |
| PRD-WF-012 | Where the configurable workflow engine is enabled, the system shall allow assigning a custom workflow per kind or per project. | Should | v0.2 |
| PRD-WF-013 | The system shall allow any issue to transition to the cancelled state from any non-cancelled state. | Must | v0.1 |

### 3.6 Exclusive Access (CLAIM)

| ID | Requirement | Priority | Notes |
|---|---|---|---|
| PRD-CLAIM-001 | The system shall support exclusive-access claims on issues with a configurable TTL (default 30 minutes). | Must | v0.1 |
| PRD-CLAIM-002 | When a live claim exists on an entity, the system shall reject new claim attempts with a message identifying the current holder and expiry time. | Must | v0.1 |
| PRD-CLAIM-003 | When a claim has expired past its TTL, the system shall allow a new claimant to acquire it. | Must | v0.1 |
| PRD-CLAIM-004 | When a caller with a force flag acquires an entity with a live claim, the system shall supersede the prior claim and record the takeover. | Must | v0.1 |
| PRD-CLAIM-005 | When transitioning an issue into a started-category state, the system shall require a live claim held by the caller. | Must | v0.1 |
| PRD-CLAIM-006 | When transitioning an issue out of a started-category state, the system shall automatically release the claim. | Should | v0.1 |
| PRD-CLAIM-007 | The system shall support listing all live claims and filtering by stale (past-TTL) status. | Must | v0.1 |
| PRD-CLAIM-008 | The system shall support human claims (`--by human`) to signal agents should not advance the entity. | Should | v0.1 |
| PRD-CLAIM-008a | The system shall allow the current claim holder to renew a live claim and optionally set a new TTL. Renewal shall be journaled and projected to the activity log. | Should | v0.1 |
| PRD-CLAIM-009 | Where QA review claims are enabled, the system shall extend claim support to QA review entities. | Should | v0.3 |

### 3.7 Role-Based Permissions (ROLE)

| ID | Requirement | Priority | Notes |
|---|---|---|---|
| PRD-ROLE-001 | The system shall require a `--by <role>` flag on every write command. | Must | v0.1 |
| PRD-ROLE-002 | If the specified role is not in the project's role set, the system shall reject the command with a clear error. | Must | v0.1 |
| PRD-ROLE-003 | If the specified role lacks permission for the requested action, the system shall reject the command with an error naming the role, the action, and that it is not permitted. | Must | v0.1 |
| PRD-ROLE-004 | The system shall allow the `human` role to perform any action via `--by human`. | Must | v0.1 |
| PRD-ROLE-005 | The system shall provide a command to display the active permissions matrix. | Should | v0.1 |
| PRD-ROLE-006 | Where configurable permissions are enabled, the system shall allow per-project role creation and permission customisation. | Should | v0.2 |

### 3.8 Cross-Entity Links (LINK)

| ID | Requirement | Priority | Notes |
|---|---|---|---|
| PRD-LINK-001 | The system shall support typed directional links between issues. | Must | v0.1 |
| PRD-LINK-002 | The system shall support links from issues to requirements. | Must | v0.1 |
| PRD-LINK-003 | When a link is created, the system shall automatically create the inverse link. | Must | v0.1 |
| PRD-LINK-004 | The system shall support the relations: blocks/blocked_by, relates_to, duplicates/duplicated_by, supersedes/superseded_by, implements/implemented_by, caused_by/causes, discovered_from/discovered. | Must | v0.1 |
| PRD-LINK-005 | If a link relation is not valid for the given entity type pair, the system shall reject it with an error. | Must | v0.1 |
| PRD-LINK-006 | The system shall list all links from and to a given entity. | Must | v0.1 |
| PRD-LINK-006a | The system shall allow an existing link to be removed. Removing a link shall remove its auto-created inverse link and project both removals to the activity log. | Should | v0.1 |
| PRD-LINK-007 | The system shall extend link endpoints to decisions and escalations when their CLI surface ships. | Should | v0.2 |
| PRD-LINK-008 | The system shall not allow links to duplicate relationships that have dedicated schema (parent_id, escalated_to_esc_id, resolution_dec, etc.). | Must | v0.1 |

### 3.9 Quality Assurance (QA)

| ID | Requirement | Priority | Notes |
|---|---|---|---|
| PRD-QA-001 | The system shall support QA reviews scoped as: epic_window (covering N tasks), single_task, or epic_full. | Must | v0.3 |
| PRD-QA-002 | When a review scope is epic_window, the system shall record which specific tasks are covered. | Must | v0.3 |
| PRD-QA-003 | The system shall support QA findings with severity (blocking or improvement), status tracking, and an optional task subject. | Must | v0.3 |
| PRD-QA-004 | The system shall support literal evidence excerpts (CLI output, file excerpts) on findings as proof. | Should | v0.3 |
| PRD-QA-005 | When a blocking finding is raised, the system shall support creating a formal escalation from that finding. | Should | v0.3 |
| PRD-QA-006 | The system shall support QA plans with ordered verification steps, each specifying an evidence method and validator role. | Should | v0.3 |
| PRD-QA-007 | Where a QA plan cites a validation contract, the system shall enforce coverage rules at plan finalisation. | Should | v0.4 |
| PRD-QA-008 | The system shall support composite QA commands that combine claim acquisition, review creation, finding addition, and review completion into role-convenient operations. | Should | v0.3 |

### 3.10 Validation Contracts (VAL)

| ID | Requirement | Priority | Notes |
|---|---|---|---|
| PRD-VAL-001 | The system shall support validation contracts per epic, each containing ordered assertions with `VAL-NNN` display IDs. | Should | v0.4 |
| PRD-VAL-002 | The system shall support required validator specifications per assertion, each declaring a validator role and evidence method. | Should | v0.4 |
| PRD-VAL-003 | When an assertion outcome is recorded, the system shall create an append-only record with outcome (pass/fail/carry_forward), evidence, and recorder role. | Should | v0.4 |
| PRD-VAL-004 | If an outcome is pass, the system shall require a non-null evidence method and evidence reference. | Should | v0.4 |
| PRD-VAL-005 | If an outcome is carry_forward, the system shall require a non-null rationale. | Should | v0.4 |
| PRD-VAL-006 | The system shall support linking acceptance criteria to validation assertions for traceability. | Should | v0.4 |
| PRD-VAL-007 | The system shall support linking QA plan steps to validation assertions for coverage tracking. | Should | v0.4 |
| PRD-VAL-008 | Where both a contract and a plan are active for the same epic, the system shall enforce that every assertion has at least one required validator and at least one covering plan step. | Should | v0.4 |

### 3.11 Governance (GOV)

| ID | Requirement | Priority | Notes |
|---|---|---|---|
| PRD-GOV-001 | The system shall support decision records with summary, prose body, area classification, and active/superseded status. | Must | v0.2. Stub tables in v0.1. |
| PRD-GOV-002 | The system shall support escalation records with open/resolved status. | Must | v0.2. Stub tables in v0.1. |
| PRD-GOV-003 | When an escalation is resolved, the system shall require a link to a resolving decision. | Must | v0.2 |
| PRD-GOV-004 | The system shall support triage event records with RESUME/BLOCK verdict, source classification (dev_concern, qa_block, esc_open), source reference, and reasoning. | Must | v0.2. Stub tables in v0.1. |
| PRD-GOV-005 | When a triage verdict is BLOCK, the system shall support creating an escalation from the triaged concern or finding. | Must | v0.2 |
| PRD-GOV-006 | The system shall support requirement pointer records with source document, source section, summary label, and a drift-detection hash. | Must | v0.2. Stub tables in v0.1. |
| PRD-GOV-007 | When a requirement's source section content has changed relative to the stored hash, the system shall report the drift. | Should | v0.2 |
| PRD-GOV-008 | The system shall support comments on issues, QA reviews, decisions, escalations, and requirements. | Should | v0.2 |
| PRD-GOV-009 | The system shall support carry-forward items with statuses: open, scheduled, closed, principled. | Should | v0.4 |
| PRD-GOV-010 | When a carry-forward item is promoted to a principle, the system shall flag it for rendering in the active decisions view. | Should | v0.4 |
| PRD-GOV-011 | The system shall support retrospective records (lightweight and deep types) with a recommendation field. | Should | v0.4 |
| PRD-GOV-012 | The system shall support composite TL-Lite triage commands that combine source inspection, verdict recording, and optional escalation creation into a single operation. | Should | v0.2 |
| PRD-GOV-013 | The system shall support a composite TL wrap-up command that generates a retrospective, closes resolved escalations, and transitions the epic. | Should | v0.2 |

### 3.12 Agent Telemetry (TELEM)

| ID | Requirement | Priority | Notes |
|---|---|---|---|
| PRD-TELEM-001 | When an agent run completes, the system shall record the role, model, input tokens, output tokens, start/end timestamps, and snapshot pricing. | Must | v0.1 |
| PRD-TELEM-002 | When an agent run is recorded, the system shall compute USD cost using snapshot pricing at run end. | Must | v0.1 |
| PRD-TELEM-003 | If the model is not found in the pricing configuration, the system shall leave cost null and emit a warning. | Must | v0.1 |
| PRD-TELEM-003a | When provider details are available, the system shall optionally record token breakdown, tool-call count, provider run ID, and provider-billed cost without replacing the normalized token and list-price summary fields. | Should | v0.1 |
| PRD-TELEM-004 | The system shall aggregate token and cost statistics by role, by model, and by epic. | Must | v0.1 |
| PRD-TELEM-005 | The system shall aggregate task statistics including estimate-vs-actual tokens, rework rate (attempts per task), and gate pass/fail counts. | Should | v0.2 |
| PRD-TELEM-006 | The system shall aggregate human intervention frequency per epic. | Should | v0.2 |

### 3.13 Output, Export, Reports, and Templates (OUT)

| ID | Requirement | Priority | Notes |
|---|---|---|---|
| PRD-OUT-001 | The system shall output JSON by default for all read commands. | Must | v0.1 |
| PRD-OUT-002 | The system shall support TOON format via `--format toon` for token-efficient agent reads. | Must | v0.1 |
| PRD-OUT-003 | The system shall support markdown format via `--format markdown` for human inspection. | Must | v0.1 |
| PRD-OUT-003a | The system shall support YAML format via `--format yaml` for config and template authoring. | Could | v0.1 |
| PRD-OUT-004 | Every read command shall accept `--format` to override the default output format. | Must | v0.1 |
| PRD-OUT-005 | Where context packs are available, the system shall return role-tuned bundles of issue data via `--context-pack <role>`. | Should | v0.2 |
| PRD-OUT-006 | The system shall provide named reports for common cross-cutting queries: epic-status, open-escalations, qa-blockers, ac-matrix, epic-wrapup-pack, task-hygiene, stale-claims. | Should | v0.2 |
| PRD-OUT-007 | The system shall list all available named reports. | Should | v0.2 |
| PRD-OUT-008 | The system shall support a read-only SQL escape hatch that rejects any non-SELECT statement. | Should | v0.1 |
| PRD-OUT-009 | The system shall support markdown export with YAML frontmatter to a configurable directory tree. | Should | v0.4 |
| PRD-OUT-010 | When exporting, the system shall include a "generated by fid" banner so humans do not mistake exports for source-of-truth files. | Should | v0.4 |
| PRD-OUT-011 | The system shall support contract-coverage, assertion-trace, and contract-outcomes reports. | Should | v0.4 |
| PRD-OUT-012 | The system shall support reusable templates with monotonic versioning. | Should | v0.3 |
| PRD-OUT-013 | When an entity is created from a template, the system shall record the template kind, name, and version on the entity. | Should | v0.3 |
| PRD-OUT-014 | Where template-driven validation is enabled, the system shall validate artifact structure (required sections) against the template at write time and reject malformed artifacts. | Should | v0.3 |
| PRD-OUT-015 | The system shall provide a TUI kanban board grouped by status category and filterable by epic, labels, and cycle. | Should | v0.3 |
| PRD-OUT-016 | The system shall support saved views (named queries) that can be run by name. | Should | v0.4 |
| PRD-OUT-017 | The system shall provide an entity activity-history command that lists projected activity log entries by entity, event type, and date range. | Should | v0.1 |
| PRD-OUT-018 | The system shall provide a lightweight `fid search <term>` command for literal, case-insensitive discovery across structured projected fields. | Should | v0.2. Search complements named reports and `fid sql`; it does not replace shell `rg` for full markdown body search. |
| PRD-OUT-019 | The system shall provide a change digest command (`fid changes --since <anchor>`) that returns, grouped by entity, the net field changes and newly-added child facts for entities changed since the anchor, scoped by epic or entity. Anchors are journal event IDs or timestamps in v0.2, resolved by an extensible resolver. | Should | v0.2. Aggregated from `activity_log`; distinct from the chronological `fid activity list` (PRD-OUT-017). Net changes collapse intermediate transitions. Role-tuned shaping deferred to a fast-follow. |

### 3.14 Operations and Maintenance (OPS)

| ID | Requirement | Priority | Notes |
|---|---|---|---|
| PRD-OPS-001 | The system shall durably persist every write as a single atomic event; no partial writes shall be observable. | Must | v0.1 |
| PRD-OPS-002 | If a composite write partially fails (e.g. one step of a multi-table operation), the system shall reject the entire operation with no state change. | Must | v0.1 |
| PRD-OPS-003 | When the query store is missing or stale, the system shall automatically rebuild it from the durable event log before answering any read. | Must | v0.1 |
| PRD-OPS-004 | The system shall support an explicit rebuild command that reconstructs the query store from scratch. | Must | v0.1. `fid rebuild`. |
| PRD-OPS-005 | The system shall provide a consistency audit that checks for kind-hierarchy violations, orphaned polymorphic references, FK integrity, and workflow state validity. | Must | v0.1 |
| PRD-OPS-006 | When the consistency audit finds errors, the system shall exit non-zero and list each violation. | Must | v0.1 |
| PRD-OPS-007 | The system shall support `--warn-as-error` on the consistency audit for CI and pre-merge use. | Should | v0.1 |
| PRD-OPS-008 | The consistency audit shall report warnings for stale claims, started-category issues without a live claim (with a configurable grace window), and telemetry gaps. | Should | v0.1 |
| PRD-OPS-009 | The system shall distribute as a single binary with no runtime dependencies. | Must | v0.1 |
| PRD-OPS-010 | The system shall support agent integration through the CLI contract, structured output, and repo skills. MCP integration is deferred until post-dogfood evaluation and requires a recorded product/security decision before entering a roadmap. | Must | v0.1 CLI/skills. MCP uncommitted. |
| PRD-OPS-011 | The system shall support replaying a foreign event log into the current project with conflict detection on entity, claim, and status. | Could | v0.2. `fid merge --journal`. Builds on v0.1 journal preconditions. |
| PRD-OPS-012 | When journal replay encounters a failed event precondition, the system shall stop before serving reads and report a structured logical conflict naming the event, entity, expected value, and actual value. | Must | v0.1. No automatic intent reconciliation in v0.1. |
| PRD-OPS-013 | The system shall provide a journal inspection command that explains a journal event by ID, including actor, command metadata, effects, and whether the event is eligible for compensating revert. | Should | v0.1. Read-only. |
| PRD-OPS-014 | The system shall support scoped compensating event reverts for explicitly reversible event types. Reverts shall append a new event and preserve the original event. | Should | v0.2. `fid event revert`. |
| PRD-OPS-015 | Every write command shall accept `--dry-run`, returning the validation result and a preview of the would-be journal event effects without appending an event or modifying state. The command exit code shall mirror the simulated outcome (zero if the write would apply, non-zero if it would be rejected). | Should | v0.2. Hooks the prepare/validate-before-append seam (TRD §4.3). Rejection vocabulary matches the real command. |
| PRD-OPS-016 | The system may accept a caller-supplied idempotency key on writes such that a repeated key returns the original event and result instead of creating a duplicate. | Could | Deferred. Conditional on a concrete programmatic/MCP retry need, recorded with a product decision. Additive to the write-path seam; pairs with the MCP evaluation. |
| PRD-OPS-017 | The system shall record git context on agent runs — branch, base commit, head commit, dirty flag, and worktree path — auto-captured read-only from git by default, overridable by flags, and degrading to a recorded null context with a warning outside a git worktree. Each run shall record whether its git context was captured, provided, or absent. | Should | v0.2. Read-only capture via `internal/gitctx/`; fid never performs git operations. Dirty is point-in-time as of the last recorded command. |
| PRD-OPS-018 | The system shall track an agent-run lifecycle (open, closed, abandoned) and provide commands to list runs filtered by lifecycle, role, and epic, and to show one run's full detail. | Should | v0.2. `fid run list --open\|--abandoned`, `fid run show <run-id>`. |
| PRD-OPS-019 | The system shall support resuming and abandoning runs. Resume shall re-acquire the claim and return the recorded git resume context, opening a new run linked to the prior; abandon shall release the claim and mark the run abandoned while preserving the run row. | Should | v0.2. `fid dev resume <task-id>`, `fid dev abandon <run-id> --reason-md --by`. Resume returns hints only; no git operations. |

## 4. Non-Functional Product Requirements

These are user-visible qualities, not implementation choices.

| ID | Requirement | Priority | Notes |
|---|---|---|---|
| PRD-NFR-001 | The system shall respond to any single read command within 500ms at single-project scale (up to 10,000 issues, 100,000 journal events). | Must | |
| PRD-NFR-002 | The system shall rebuild the query store from 10,000 journal events in under 5 seconds. | Must | |
| PRD-NFR-003 | When rejecting a command, the system shall provide an error message that names the specific constraint violated, the entities involved, and (where applicable) the corrective action. | Must | |
| PRD-NFR-004 | The system shall never store project data outside the project's `.fid/` directory and referenced markdown files. | Must | Local-first guarantee. |
| PRD-NFR-005 | The system shall never transmit project data to any external service. | Must | No remote sync in v1. |
| PRD-NFR-006 | The system shall operate without a background daemon, network listener, or external service dependency. | Must | |
| PRD-NFR-007 | The system shall merge cleanly across git branches using standard text-based merge for its durable event log. | Must | No binary-blob merge conflicts. |
| PRD-NFR-008 | When two agents in the same checkout issue concurrent writes, the system shall serialise them with a checkout-scoped write lock before validation, journal append, and projection, without data loss. | Must | v0.1 uses `.fid/write.lock` plus SQLite projection transactions. |
| PRD-NFR-009 | The durable event log shall be human-readable (one JSON event per line) so that merge conflicts can be resolved by inspecting the text. | Should | |
| PRD-NFR-010 | The system shall ship with default pricing data for current Anthropic, OpenAI, and Google models, user-overridable via configuration. | Should | |

## 5. Error and Empty States

Define how the product behaves when data is missing, invalid, incomplete, or
ambiguous.

| ID | Condition | Required Behavior |
|---|---|---|
| PRD-ERR-001 | Issue ID does not exist | Return "issue <id> not found" with non-zero exit code. |
| PRD-ERR-002 | Ready queue has no actionable tasks (all blocked, all completed, or no tasks exist) | Return an empty result set with zero exit code and a human-readable note explaining why no tasks are ready. |
| PRD-ERR-003 | Claim acquisition on an entity with a live (non-expired) claim | Return "already claimed by <role> until <expiry>" with non-zero exit code. |
| PRD-ERR-004 | Transition request violating the workflow state machine | Return "transition from <current> to <requested> is not allowed; permitted transitions: <list>" with non-zero exit code. |
| PRD-ERR-005 | Transition request failing a validation rule | Return the specific unmet conditions (e.g. "3 of 5 ACs not verified by developer: <ids>") with non-zero exit code. |
| PRD-ERR-006 | Write command with `--by <unknown-role>` | Return "role <X> not in project profile; known roles: <list>" with non-zero exit code. |
| PRD-ERR-007 | Write command with `--by <role>` that lacks permission for the action | Return "role <role> is not permitted to <action>" with non-zero exit code. |
| PRD-ERR-008 | Hash ID collision on insert | Retry with a new hash transparently; if collision persists after retry, return an error. |
| PRD-ERR-009 | Malformed trailing line in the durable event log | Stop processing, report the corrupt line number and content, and prompt for manual repair. Do not silently skip or truncate. |
| PRD-ERR-010 | Link creation with an invalid relation/entity-type combination | Return "<relation> is not valid between <from_type> and <to_type>" with non-zero exit code. |
| PRD-ERR-011 | AC addition or completion creation on an epic-kind issue | Return "acceptance criteria / completions are not valid on epics" with non-zero exit code. |
| PRD-ERR-012 | Command targeting a v0.2+ entity whose CLI surface has not shipped (e.g. `fid decision create` in v0.1) | Return "commands for <entity> ship in <version>" with non-zero exit code. |
| PRD-ERR-013 | Export to a directory that does not exist | Create the directory tree rather than failing. |

## 6. Traceability

Each roadmap version should reference PRD IDs. Each implemented PRD ID should
have tests or verification evidence. Traceability to tasks populated during
phase planning.

| PRD ID(s) | Roadmap Version | Task(s) | Verification |
|---|---|---|---|
| PRD-INIT-001 – 004 | v0.1 | | |
| PRD-ISSUE-001 – 021 | v0.1 | | |
| PRD-AC-001 – 006 | v0.1 | | |
| PRD-DEV-001 – 009 | v0.1 | | |
| PRD-WF-001 – 010, 013 | v0.1 | | |
| PRD-CLAIM-001 – 008 | v0.1 | | |
| PRD-ROLE-001 – 005 | v0.1 | | |
| PRD-LINK-001 – 006, 008 | v0.1 | | |
| PRD-TELEM-001 – 004 | v0.1 | | |
| PRD-OUT-001 – 003a, 004, 008 | v0.1 | | |
| PRD-OPS-001 – 009 | v0.1 | | |
| PRD-WF-005a, PRD-WF-011 – 012, PRD-ROLE-006, PRD-LINK-007 | v0.2 | | |
| PRD-GOV-001 – 008, 012 – 013 | v0.2 | | |
| PRD-TELEM-005 – 006 | v0.2 | | |
| PRD-OUT-005 – 007, 018 – 019 | v0.2 | | |
| PRD-OPS-011, 014 – 015 | v0.2 | | |
| PRD-OPS-017 – 019 | v0.2 | | |
| PRD-QA-001 – 006, 008, PRD-CLAIM-009 | v0.3 | | |
| PRD-OUT-012 – 015 | v0.3 | | |
| PRD-QA-007, PRD-VAL-001 – 008 | v0.4 | | |
| PRD-GOV-009 – 011 | v0.4 | | |
| PRD-OUT-009 – 011, 016 | v0.4 | | |
| PRD-ISSUE-017 | v0.5 | | |

## 7. Open Product Questions

These are product and scope questions. Technical open questions are tracked
separately in `02_TRD.md` §11.

- [ ] Should the `workflow_runner` role become a deterministic harness rather
  than an LLM role? Deferred from spec §15.1; revisit when routine
  workflow_runner cost or prompt drift causes a real failure.
- [ ] Should bidirectional markdown sync (`fid import`, `fid sync`) ship in
  v0.5, or is the agent-with-fid-skill alternative sufficient? Deferred
  pending evidence from v0.4 export usage.
- [ ] What is the licensing model for open-source release?

*Last updated: 2026-05-31*
