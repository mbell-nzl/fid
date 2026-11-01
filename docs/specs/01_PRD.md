# PRODUCT REQUIREMENTS - fid

This document defines the aligned product after ADR-0011. Requirements use
stable IDs and testable “shall” language. `P0` is required for the first
supported fid/V2 baseline; `P1` is retained product scope that may follow the
backend critical path. Evidence-gated ideas live only in ADR-0011.

## 1. Foundation

| ID | Priority | Requirement |
|---|---|---|
| PRD-FOUND-001 | P0 | fid shall operate as a local project-scoped CLI without a daemon, listener, hosted service, or mandatory network access. |
| PRD-FOUND-002 | P0 | Every durable mutation shall append an immutable ordered journal Event before or atomically with projection application. |
| PRD-FOUND-003 | P0 | SQLite shall be a derived projection that can be deleted and rebuilt from the journal plus versioned system reference data. |
| PRD-FOUND-004 | P0 | Public machine operations shall accept and return versioned JSON. JSON shall be the only canonical machine representation. |
| PRD-FOUND-005 | P0 | Human CLI output shall be concise terminal text unless a named Human projection explicitly defines Markdown or HTML rendering. |
| PRD-FOUND-006 | P0 | fid shall not expose universal TOON, YAML, or Markdown format options. |
| PRD-FOUND-007 | P0 | Capability discovery shall report independently versioned semantic capabilities, supported schema ranges, limits, availability, and trust sources. |
| PRD-FOUND-008 | P0 | A consumer requiring an unavailable or incompatible capability shall receive a typed rejection; fid shall not silently use a weaker fallback. |

## 2. Issue Tracker

| ID | Priority | Requirement |
|---|---|---|
| PRD-ISSUE-001 | P0 | fid shall create, read, list, update, and transition Issues of kind epic, task, or subtask. |
| PRD-ISSUE-002 | P0 | Epics shall be top-level; tasks may be top-level or belong to one Epic; subtasks shall belong to one Task. |
| PRD-ISSUE-003 | P0 | Issue parentage and dependencies shall be the only durable work graph. |
| PRD-ISSUE-004 | P0 | fid shall reject parent cycles, dependency cycles, invalid hierarchy, and references to missing Issues. |
| PRD-ISSUE-005 | P0 | Sequential Epics shall expose only the lowest-sequence unfinished unblocked child as ready; prioritized Epics shall rank ready children deterministically. |
| PRD-ISSUE-006 | P0 | fid shall support stable labels and typed links without duplicating dedicated lifecycle fields as links. |
| PRD-ISSUE-007 | P0 | An Issue may contain ordered acceptance criteria describing observable required outcomes. |
| PRD-ISSUE-008 | P0 | Acceptance state shall be derived from qualifying current Evidence, Review, or Human Outcome records where an active Envelope requires them; a caller assertion alone shall not prove acceptance. |
| PRD-ISSUE-009 | P0 | Actionable-work queries shall return typed reasons for both actionable and non-actionable Issues. |
| PRD-ISSUE-010 | P0 | Inventory authoring outside an active kit workflow shall remain available as attributed local Issue-tracker operations and shall not imply protected Human authority or trusted evidence. |
| PRD-ISSUE-011 | P1 | fid may synchronize stable requirement identity, source reference, digest, and affected Issue links while project files remain canonical. |
| PRD-ISSUE-012 | P1 | Requirement drift shall report changed source digests and affected open work without parsing requirement prose into a second requirements model. |

## 3. Issue Workflows

| ID | Priority | Requirement |
|---|---|---|
| PRD-WF-001 | P0 | Issue lifecycle shall be defined by named workflows containing states, status categories, legal transitions, and generic transition preconditions. |
| PRD-WF-002 | P0 | fid shall support workflow definition, assignment by project or Issue kind, and deterministic assignment precedence. |
| PRD-WF-003 | P0 | Unknown states, illegal edges, and unsupported preconditions shall reject before any durable effect. |
| PRD-WF-004 | P0 | Execution-bound transitions shall atomically enforce the active Envelope revision, actor authority, Claim lease, expected state revision, and applicable Obligations. |
| PRD-WF-005 | P0 | fid shall not encode kit risk lanes, provider roles, or assurance choreography as Issue workflow definitions. |

## 4. Actor And Authority

| ID | Priority | Requirement |
|---|---|---|
| PRD-AUTH-001 | P0 | Every durable operation shall record an Actor identity and producer provenance. |
| PRD-AUTH-002 | P0 | Backend-port authorization shall evaluate generic Actor capabilities and approved Work Envelope grants, not a fid-global role taxonomy. |
| PRD-AUTH-003 | P0 | Absence, expiry, forgery, or mismatch of required authority shall reject without partial effects. |
| PRD-AUTH-004 | P0 | fid shall enforce the recorded approved authority but shall not recalculate kit role, routing, risk, or Obligation-selection policy. |
| PRD-AUTH-005 | P0 | A protected Decision shall block affected execution until a qualifying Decision Reference and Human Outcome satisfy its recorded authority requirement. |
| PRD-AUTH-006 | P0 | A caller-supplied role label or `human` string shall not by itself establish protected Human authority. |
| PRD-AUTH-007 | P0 | Standalone Issue inventory operations may use project-local attributed authority but shall not create trusted execution, protected Decisions, or Human Outcomes. |
| PRD-AUTH-008 | P0 | Kit-mode Human Outcomes shall require provenance authorized by the active Kit Activation. |

## 5. Work Envelopes And Obligations

| ID | Priority | Requirement |
|---|---|---|
| PRD-ENV-001 | P0 | Each active kit-managed outcome shall have one versioned Work Envelope selecting one primary fid Issue and an explicitly bounded related Issue set. |
| PRD-ENV-002 | P0 | A Work Envelope shall not copy Issue title, acceptance state, dependencies, or lifecycle into persistent Work Items. |
| PRD-ENV-003 | P0 | An approved Envelope revision shall bind outcome and non-goals, source references, repository scope, risk result, authority, protected Decisions, selected Issues, Obligations, budgets, stop conditions, and approval provenance. |
| PRD-ENV-004 | P0 | Material selection or contract changes shall create a new immutable revision; prior Claims, Runs, Reviews, and Receipts shall retain their original revision binding. |
| PRD-ENV-005 | P0 | A new revision shall not activate while a mutating Claim or Run remains active against its predecessor. |
| PRD-ENV-006 | P0 | Removing unfinished selected work shall record an authorized successor, deferral, inapplicability, or approved outcome-change disposition. |
| PRD-ENV-007 | P0 | Partial progress shall be derived from selected Issue states and shall have no independent durable lifecycle. |
| PRD-ENV-008 | P0 | Completion shall require every active-revision Issue to be in a satisfying terminal state and every applicable Obligation to be current and satisfied. |
| PRD-ENV-009 | P0 | Cancellation shall not satisfy completion or dependencies. |
| PRD-ENV-010 | P0 | Obligations shall directly encode applicability, accepted evidence/trust, producer independence, artifact policy, freshness, and guarded action without a QA Plan or Validation Contract hierarchy. |

**Phase 21 qualification (DEC-156).** For the one selected Issue, the supported
gate-evidence Obligation guards the terminal-success Issue transition. Gate
execution, qualifying Receipt linking, and satisfied Evidence evaluation occur
before that transition. Envelope completion is a later mutation that
re-evaluates the same Obligation and also requires the terminal-success Issue.
Both actions use one evaluator and return the same all-reasons semantics from
PRD-EVID-006. This qualification does not change the broader target ownership
of Review/Human Obligations in later phases.

## 6. Claims, Runs, Stops, And Recovery

| ID | Priority | Requirement |
|---|---|---|
| PRD-RUN-001 | P0 | A Claim shall bind subject, Issue, Envelope revision, Actor, capabilities, optional Obligation, lease version, and expiry. |
| PRD-RUN-002 | P0 | One Issue shall have at most one active mutating Claim. |
| PRD-RUN-003 | P0 | Parallel mutation shall require distinct selected Issues, separate worktrees, and declared provably disjoint mutation scopes; unknown or overlapping scope shall return `mutation-conflict`. |
| PRD-RUN-004 | P0 | Claim acquire, renewal, release, expiry, and authorized supersession shall use optimistic versions and idempotent semantics. |
| PRD-RUN-005 | P0 | A Run shall bind its Issue, Envelope revision, Claim, Actor, repository identity, budgets, optional Obligation, and predecessor/successor recovery links. |
| PRD-RUN-006 | P0 | Budget consumption and exhaustion shall be durable; exhaustion shall atomically create the applicable typed Stop and Attention Item. |
| PRD-RUN-007 | P0 | Continuation after interruption shall create a linked new Run and shall never reuse the prior Run identity. |
| PRD-RUN-008 | P0 | Recovery shall require confirmed prior-executor termination or authorized force recovery, eligible Claim disposition, and non-destructive repository reconciliation. |
| PRD-RUN-009 | P0 | Unknown executor liveness shall return `executor-liveness-unknown`; unsafe or ambiguous repository state shall return `recovery-state-conflict`. |
| PRD-RUN-010 | P0 | Recovery shall not automatically reset, clean, stash, checkout, commit, or discard repository state. |
| PRD-RUN-011 | P0 | Recovered work may be adopted under recorded authority but shall not become trusted evidence by adoption. |
| PRD-RUN-012 | P0 | Stops shall be typed, scoped, immutable facts with explicit resolution conditions. |
| PRD-RUN-013 | P0 | Equivalent Attention requests shall deduplicate by stable semantic identity without discarding distinct Stop facts. |
| PRD-RUN-014 | P0 | A fresh session shall recover the next valid action, outstanding Obligations, evidence state, unresolved Findings, Stops, Attention, and remaining budgets without prior chat. |

## 7. Registered Execution And Evidence

| ID | Priority | Requirement |
|---|---|---|
| PRD-GATE-001 | P0 | The project and kit shall author and validate versioned Gate Definitions; fid shall create immutable authorized registrations from the activated values. |
| PRD-GATE-002 | P0 | Changing definition content or binding shall create a new registration and shall not mutate an existing one. |
| PRD-GATE-003 | P0 | Trusted local execution shall accept only a Gate Definition Registration reference plus subject, Run, and repository identity; callers shall not substitute invocation or policy fields. |
| PRD-GATE-004 | P0 | fid shall append an Execution Attempt before process spawn and at most one immutable terminal Receipt after observing termination. |
| PRD-GATE-005 | P0 | A missing Receipt, crash, timeout, cancellation, signal, or harness failure shall never imply pass. |
| PRD-GATE-006 | P0 | A Receipt shall bind registration, producer, trust source, subject/Run, repository identity, invocation, timing, terminal outcome, toolchain, and sanitized output identity. |
| PRD-GATE-007 | P0 | Gate output shall stream live to the invoking Human or harness while persisted excerpts remain bounded, sanitized, hashed, and retention-aware. |
| PRD-GATE-008 | P0 | A separately trusted producer may submit an attested Receipt only through an accepted versioned attestation contract. |
| PRD-GATE-009 | P0 | Caller-reported commands or exit codes shall not be persisted as deterministic evidence. |
| PRD-EVID-001 | P0 | Evidence Records shall link Obligations to qualifying Receipts, Reviews, Human Outcomes, or external artifact identities. |
| PRD-EVID-002 | P0 | Obligation absence, Evidence Record validity, and Referenced Artifact availability shall be serialized and evaluated as distinct dimensions. |
| PRD-EVID-003 | P0 | Initial deterministic freshness shall conservatively stale repository-bound proof after a target-tree change unless the Obligation explicitly targets an unchanged external artifact. |
| PRD-EVID-004 | P0 | Evidence Records shall retain identity, provenance, reference, and digest after a Referenced Artifact expires or becomes unavailable. |
| PRD-EVID-005 | P0 | An Obligation shall specify whether metadata/digest is sufficient or full material must remain available. |
| PRD-EVID-006 | P0 | Obligation evaluation and guarded transitions shall use the same deterministic rules and return all unsatisfied reasons. |

## 8. Review, Concerns, And Human Judgment

| ID | Priority | Requirement |
|---|---|---|
| PRD-REV-001 | P0 | A completed Review Run shall create an immutable Review Record bound to exact Envelope revision and repository identity. |
| PRD-REV-002 | P0 | A Review Record shall declare profile, coverage, verdict, evidence, and zero or more Findings. |
| PRD-REV-003 | P0 | `changes-required` shall leave the review Obligation unsatisfied and expose Builder rework on the original Issue/Envelope. |
| PRD-REV-004 | P0 | A Builder may record a Finding as addressed, which remains nonterminal; only a later qualifying independent Review Record may resolve it or mark it inapplicable with a reason. |
| PRD-REV-005 | P0 | Target changes shall stale target-bound Review Records and Receipts according to the active Obligation policy. |
| PRD-REV-006 | P0 | Rework shall consume existing budgets; exhaustion shall return `review-rework-exhausted` plus Attention. |
| PRD-REV-007 | P0 | Concerns and Deviations shall record only the identity, subject, scope, severity, status, authority, disposition, references, and concise reasoning required to govern work. |
| PRD-REV-008 | P0 | Blocking work shall not be deferred merely to permit completion. |
| PRD-REV-009 | P0 | Project Decision prose shall remain in project files; fid shall store only an operational Decision Reference and outcome. |
| PRD-HUM-001 | P0 | Human Outcomes shall satisfy only Obligations that explicitly require or permit Human judgment and shall never be presented as deterministic proof. |
| PRD-HUM-002 | P0 | Attention shall be raised only for a real decision, policy Stop, circuit breaker, unresolved risk, or readiness for Human validation; routine status shall remain queryable without interruption. |
| PRD-HUM-003 | P0 | Resolving Attention shall reference the qualifying Outcome or operational fact rather than erase the original request. |

## 9. Projections, Recovery, And Telemetry

| ID | Priority | Requirement |
|---|---|---|
| PRD-PROJ-001 | P0 | fid shall provide versioned JSON projections for actionable work, context, changed-since state, resume, completion/evidence trace, Attention, and Human validation. |
| PRD-PROJ-002 | P0 | Projections shall be reproducible from durable state and shall not be editable sources. |
| PRD-PROJ-003 | P0 | Context projections shall be bounded to the requested subject, Actor needs, and optional changed-since cursor without copying workflow protocol prose. |
| PRD-PROJ-004 | P0 | The Human Validation Brief shall present delivered outcomes, material Decisions, evidence roll-up, Deviations, unresolved risk, economics, and remaining Human-only checks without a full transcript. |
| PRD-PROJ-005 | P0 | Purpose-built Human Markdown or HTML shall be derived from the same projection schema as JSON and shall not define a second contract. |
| PRD-PROJ-006 | P0 | Human views shall render Stop, Attention and repository state as labelled text naming the qualifying resolution and its command, never as raw JSON or opaque internal identifiers (DEC-272). |
| PRD-PROJ-007 | P0 | Human views shall show recorded reasoning (summarized in lists and in full in detail views) and shall never reproduce the content of a referenced project file (DEC-272, DEC-274). |
| PRD-REC-001 | P0 | Portable export shall be versioned, backend-neutral, integrity-protected, and sufficient to reconstruct retained operational state and current projections. |
| PRD-REC-002 | P0 | Export validation shall support a non-mutating validation-only mode. |
| PRD-REC-003 | P0 | Import shall preserve stable identities and history unless an explicit versioned migration maps them. |
| PRD-REC-004 | P0 | Rebuild, export/import, and live projection shall produce semantically equivalent current state. |
| PRD-TELEM-001 | P0 | Run telemetry shall record observed token, cost, duration, retry, intervention, failure, and tool data incrementally and shall preserve unknown values as null. |
| PRD-TELEM-002 | P0 | Telemetry shall bind Actor, Run, Issue, Envelope revision, provider/model identity, observation source, and timestamp. |
| PRD-TELEM-003 | P0 | Basic trajectory summaries shall support comparison of workflow cost, quality, retries, and Human attention without agent-entered estimates. |

## 10. Operations And Security

| ID | Priority | Requirement |
|---|---|---|
| PRD-OPS-001 | P0 | Writes shall be serialized across processes and one logical mutation shall commit one journal transaction or no effects. |
| PRD-OPS-002 | P0 | Create and transition operations shall accept idempotency keys and expected revisions where retry or concurrency is possible. |
| PRD-OPS-003 | P0 | Maintenance shall detect journal/projection divergence, invalid references, broken sequences, abandoned execution Attempts, stale Claims, binding conflicts, missing artifacts, and activation disagreement. |
| PRD-OPS-004 | P0 | Kit Activation shall prepare, validate, commit, abort, finalize, and recover kit/backend changes under one durable activation identity. |
| PRD-OPS-005 | P0 | Normal kit workflow shall stop while activation is incomplete or components disagree. |
| PRD-OPS-006 | P0 | Before an irreversible backend commit recovery may restore the prior active version; after it recovery shall complete forward and shall not claim rollback. |
| PRD-OPS-007 | P0 | Sensitive gate output shall be sanitized before persistence; secrets and control-sequence injection shall not enter Human projections. |
| PRD-OPS-008 | P0 | Errors shall be typed and machine-readable in JSON, include no partial success ambiguity, and use non-zero process status for failed operations. |
| PRD-OPS-009 | P0 | The production port shall never read fid SQLite tables, journal files, or `.fid/` internals directly. |
| PRD-OPS-010 | P0 | The in-memory backend fake shall be test-only and shall never be selected as a production fallback. |
| PRD-OPS-011 | P0 | fid shall observe repository identity itself. A mutation input shall carry only expectations that a check uses (expected HEAD and uncommitted-change state, and an optional expected project-file digest). fid shall record its own observation (DEC-272). |
| PRD-OPS-012 | P0 | An idempotency request digest shall cover the caller's input as written and never fid-observed values. An idempotency conflict shall name each differing input (DEC-272, DEC-275). |
| PRD-OPS-013 | P0 | A Human shall be able to author any mutation input from fid's own output: help shall document every input field, a template command shall print an editable input, and a malformed input shall name the offending field (DEC-272, DEC-276). |
| PRD-OPS-014 | P0 | Human failure output shall state the JSON failure message, the failure reason, any blocking Stop and Attention, and the next legitimate action where one exists (DEC-272). |

## 11. First Supported Baseline

The first supported baseline exists only when:

- the shared V2 fake/fid conformance suite passes for every required capability;
- interruption, review rework, stale evidence, activation failure, rebuild, and
  portable export scenarios pass;
- all universal TOON/YAML/Markdown and superseded V1 role/composite surfaces
  have been removed;
- the migration and journal baseline has been consolidated; and
- a disposable project rehearsal completes before live dogfood.

*Last updated: 2026-10-04 (Phase 24 planning: PRD-OPS-011..014 and PRD-PROJ-006/007 under DEC-272..276)*
