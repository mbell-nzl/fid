# TARGET LOGICAL SCHEMA AND INVARIANTS - fid

This document replaces the speculative pre-ADR-0011 DDL. It is the accepted
logical schema for roadmap planning and conformance. Physical SQLite DDL and
canonical Event payloads are added here slice by slice before product code may
depend on them. The final convergence slice consolidates those migrations into
the first supported pre-dogfood baseline.

## 1. Shared Rules

- Entity IDs use stable opaque strings. fid may generate readable prefixes,
  but consumers never derive semantics from them.
- All timestamps are UTC RFC 3339 values in JSON and normalized UTC values in
  SQLite.
- Mutable projections carry an integer revision for optimistic concurrency.
- Immutable records have no update operation; correction appends a superseding
  or disposition record.
- JSON payload columns are permitted only for versioned bounded structures
  whose complete schema is defined by the public contract.
- Foreign keys and required uniqueness constraints are enabled on every
  connection. Triggers are corruption backstops, not the primary policy layer.
- Secrets, unbounded logs, raw environment snapshots, project prose, and source
  file contents are forbidden from normal projection rows.

## 2. Common Value Types

### 2.1 Actor reference

Required fields: opaque Actor ID, Actor type, producer, provenance reference,
and optional session/provider reference. Capabilities and authority are not a
global role record; applicable grants belong to an Envelope Revision or active
registration/activation.

### 2.2 Repository identity

Required fields: configured repository/root identity, worktree identity,
base/head/tree object IDs where applicable, dirty boolean, deterministic dirty
digest when dirty, and capture source. Missing Git objects serialize as null,
never a branch-name substitute.

### 2.3 Source reference

Required fields: source type, project-relative path or stable external ID,
optional section, repository identity/digest, and display label. It references
project truth without copying its prose.

### 2.4 Budget

Required fields: name, unit, maximum, consumed value or null, warning threshold,
scope, and exhaustion action. Initial units cover attempts, failures, time,
input/output tokens, cost, tool failures, diff size, and interventions.

### 2.5 Obligation

Required fields: ID, type, subject, source reference, applicability, required
lane/profile where supplied by kit policy, accepted evidence/trust sources,
producer/validator independence, Artifact policy, freshness target, guarded
action, and derived status. Status is one of `satisfied`, `unsatisfied`,
`stale`, `blocked`, or `inapplicable`.

## 3. Projection Table Catalog

The catalog names the target tables and minimum columns. Exact SQL types,
indexes, checks, and migration numbers are accepted with the owning slice.

### 3.1 Projection control

| Table | Minimum identity and content | Rules |
|---|---|---|
| `projection_meta` | singleton schema version, last Event ID/hash, applied count | Updated only after the projection transaction commits. |
| `idempotency_results` | operation scope, key, request digest, result entity/Event, created time | Same key/digest returns original result; different digest conflicts. |
| `system_capabilities` | capability ID/version, schema range, availability, limits, trust sources | Derived from running binary and active migration state. |

### 3.2 Issue tracker

| Table | Minimum identity and content | Rules |
|---|---|---|
| `issues` | ID, kind, title, description/reference, parent, priority, sequence, workflow/state, revision, created/updated/closed time | Parent hierarchy and workflow state must be valid. |
| `issue_dependencies` | Issue, prerequisite Issue, created Event | No self-link or directed cycle. |
| `epic_scheduling` | Epic ID, mode | One row per Epic; mode is sequential or prioritized. |
| `acceptance_criteria` | ID, Issue, ordinal, text, created Event | Unique ordinal per Issue. No verification role columns. |
| `labels` | name, colour, description | Stable project label identity. |
| `issue_labels` | Issue, label | Composite primary key. |
| `links` | from type/ID, relation, to type/ID, created Event | Endpoints use a closed supported type registry; dedicated lifecycle fields are not duplicated. |
| `requirement_refs` | stable requirement ID, source reference/digest, status | Project files remain canonical. |
| `requirement_issue_links` | requirement ID, Issue ID, relation | Traceability only. |

### 3.3 Issue workflow registry

| Table | Minimum identity and content | Rules |
|---|---|---|
| `workflows` | ID, slug, name, version, active flag, created Event | Immutable version identity; lifecycle revision semantics must be decided before mutable authoring is added. |
| `workflow_states` | workflow/version, state, category, ordinal | Unique state and ordinal per workflow version. |
| `workflow_transitions` | workflow/version, from state, to state | Both endpoints must exist in the same version. |
| `workflow_transition_requires` | transition, generic precondition code, versioned parameters | No kit role names or risk-policy evaluation. |
| `workflow_assignments` | scope type/value, workflow/version, created Event | Kind assignment precedes project assignment; one active assignment per scope. |

### 3.4 Work Envelopes

| Table | Minimum identity and content | Rules |
|---|---|---|
| `work_envelopes` | Envelope ID, current revision, status, created time | Stable identity only; material contract is revisioned. |
| `envelope_revisions` | Envelope/revision, prior revision, outcome, non-goals, risk result, repository scope, status, created/approved Actor/time | Immutable; active revision unique per Envelope. |
| `envelope_sources` | Envelope/revision, source reference | References only. |
| `envelope_issue_selection` | Envelope/revision, Issue, primary flag, disposition for removed unfinished work | Exactly one primary; selected Issues already exist. |
| `envelope_authority_grants` | Envelope/revision, Actor/capability selector, action or Decision class, scope | Absence is not allow. No fid-global roles. |
| `envelope_protected_decisions` | Envelope/revision, Decision reference, class, status, required authority, affected scope | Unsatisfied required Decision blocks. |
| `envelope_obligations` | Envelope/revision plus Obligation value | Immutable contract; satisfaction is derived from Evidence/Outcomes. |
| `envelope_budgets` | Envelope/revision plus Budget value | Consumption survives recovery and revision according to policy. |
| `envelope_stop_conditions` | Envelope/revision, typed condition and required resolver authority | Immutable contract. |

### 3.5 Claims, Runs, Stops, and Attention

| Table | Minimum identity and content | Rules |
|---|---|---|
| `claims` | ID, subject, Issue, Envelope/revision, Actor, optional Obligation, worktree/scope, mode, version, acquired/expires/released time and reason | One active mutating Claim per Issue; stale version cannot renew or authorize. |
| `runs` | ID, Issue, Envelope/revision, Claim, Actor, optional Obligation, repository start/terminal identity, status, predecessor, started/ended time | Immutable identity; interrupted continuation uses a new linked Run. |
| `run_recoveries` | prior/new Run and Claim, termination provenance, repository observation, reconciliation/adoption outcome, authority, reason | One atomic recovery fact; no destructive Git action. |
| `run_telemetry` | ID, Run, metric, observed value/null, unit, source, observed time, terminal flag | Append-only observed samples. |
| `stops` | ID, subject, code, versioned details, created Actor/time, resolution condition, resolved Event | Distinct semantics remain distinct facts. |
| `attention_items` | ID, semantic deduplication key, subject, class, details, status, raised time, resolution Outcome/fact | Equivalent open requests deduplicate; history remains. |

### 3.6 Registered gates and evidence

| Table | Minimum identity and content | Rules |
|---|---|---|
| `gate_registrations` | registration ID/version/status, definition ID/version/digest, binding digest, project/overlay/repository/scope binding, authorization provenance, structured execution policy, created Event | Immutable; content/binding change creates a new version. |
| `execution_attempts` | execution ID, registration, subject, Run, start repository identity, resolved policy digest, producer, started time | Appended before spawn. |
| `execution_receipts` | Receipt ID, execution ID, trust source, terminal status, exit/signal, timing, terminal repository identity, toolchain, sanitized excerpt/hash/size flags, Artifact ref | Immutable; one per execution. |
| `evidence_records` | Evidence ID, Obligation, evidence type/ref, producer, trust source, target identity, validity, created time | Append-only; no caller-authored pass flag. |
| `referenced_artifacts` | opaque/content ID, digest, media/type, retention class, availability, location policy, created/expiry time | Loss changes availability, not Evidence identity. |
| `evidence_artifact_links` | Evidence ID, Artifact ID, required-material flag | Composite primary key. |

### 3.7 Review and judgment

| Table | Minimum identity and content | Rules |
|---|---|---|
| `review_records` | Review ID, Run, Envelope/revision, exact repository identity, profile, coverage, verdict, evidence refs, created time | Immutable; Findings may be empty. |
| `review_findings` | Finding ID, Review, subject/scope, severity, summary, blocking flag | Created only with a Review Record. |
| `finding_dispositions` | Finding, Actor, addressed/resolved/inapplicable action, qualifying later Review/Outcome/reference, time | Builder may address; qualifying Reviewer/Human authority resolves. |
| `concerns` | ID, subject, Actor, scope, severity, summary, status, references, created time | Generic subject; not completion-linked. |
| `deviations` | ID, subject, expected contract/control, observed departure, authority, disposition, references, created time | Append-only decision-relevant fact. |
| `decision_refs` | ID, class, source reference/digest, authority, outcome, scope, created time | No decision body. |
| `human_outcomes` | ID, Actor/provenance, class, subject, result, scope, references, created time | Only active-activation provenance can exercise protected authority. |

### 3.8 Kit Activation

| Table | Minimum identity and content | Rules |
|---|---|---|
| `kit_activations` | ID, old/new kit and backend identities, phase, point-of-no-return state, authoritative direction, created/updated time | One nonterminal activation blocks normal workflow. |
| `activation_components` | activation, component, old/new identity, expected revision, state, evidence/error | Every required component must agree before active. |

## 4. Critical Invariants

- Every selected Issue exists; one selection is primary.
- Progress is queried from selected Issue state and is never persisted as an
  Envelope or membership status.
- Completion atomically re-evaluates active revision, every selected Issue, and
  every applicable Obligation.
- New Envelope revision activation requires no active mutation against its
  predecessor.
- Actor authorization uses declared identity/capabilities and approved grants;
  missing or forged authority rejects.
- One active mutation Claim per Issue; parallel Claims require distinct
  worktrees and non-overlapping declared scopes.
- Run recovery never resets consumed budgets, Issues, Envelope state,
  Obligations, Concerns, or evidence policy.
- Target changes stale repository-bound Receipts and Review Records according
  to Obligation policy.
- A Builder cannot resolve its own Finding.
- A Human Outcome cannot satisfy deterministic execution.
- An Execution Attempt without a terminal Receipt cannot pass.
- Receipt invocation and policy identity exactly match the resolved immutable
  registration.
- Evidence absence, record validity, and Artifact availability remain distinct.
- Activation disagreement blocks normal work.

## 5. Canonical Event Families

Each accepted family needs a versioned JSON payload fixture and live/rebuild
parity test before implementation:

| Aggregate | Event families |
|---|---|
| Issue | created, updated, transitioned, dependency/label/link/AC changed, workflow assigned |
| Envelope | created, approved, activated, revised, blocked, completed, cancelled |
| Claim/Run | claimed, renewed, released, started, telemetry recorded, terminalized, recovered |
| Stop/Attention | raised, deduplicated, resolved |
| Gate/Evidence | registered, superseded, execution started, Receipt recorded, attestation ingested, Evidence linked, Artifact availability changed |
| Review/Judgment | Review recorded, Finding addressed/resolved, Concern/Deviation recorded, Decision Reference recorded, Human Outcome recorded |
| Activation | prepared, component recorded, committed, aborted, finalized, recovery direction selected/completed |

Every Event carries aggregate sequence, idempotency, Actor/producer,
correlation/causation, time, and repository identity when relevant. Sequence
gaps or conflicts are maintenance failures.

## 6. Removed Pre-Dogfood Schema

The supported baseline shall not contain the old `roles`, `gate_runs`,
role-specific `task_completions`, `triage_events`, `escalations`, QA Plan,
Validation Contract/assertion, template, retrospective, saved-view, cycle, or
Markdown-export tables. Generic replacements above retain only behavior earned
by the ADR-0011 durability test.

The current development migrations remain executable only until their owning
replacement slices land. No compatibility migration or old-Event upgrader is
required. Before first dogfood, a clean init migration and full Event fixture
set replace the development chain.

## 7. Physical DDL Gate

For each D3 vertical slice, the task pack must amend this document with:

- exact `CREATE TABLE`/index/trigger DDL for newly implemented tables;
- exact down/reset behavior for the development baseline;
- JSON Schema or equivalent fixture contract for every new Event payload;
- live/rebuild/export/import parity assertions; and
- deletion mapping for displaced pre-dogfood tables and tests.

Until that amendment is accepted, implementation depending on that table or
Event family is blocked rather than inferred from this logical catalog.

## 8. Accepted Slice DDL

Entries in this section are accepted only with the Human-approved owning slice
contract and task pack.

### 8.1 Phase 20 - clean standalone Issue journey

Migrations `0011_p20_remove_issue_fields` and `0012_p20_issue_revision`
realign the `issues` projection with section 3.2. Full slice contract, including
JSON results and the removal inventory, is in
`../docs/phases/phase-20/contracts.md`.

`0011` up (logical shape; a table-rebuild implementation in the `0008` style is
permitted because the dependent CHECK constraint is removed with its column):

```sql
ALTER TABLE issues DROP COLUMN in_inbox;
ALTER TABLE issues DROP COLUMN est_tokens;
ALTER TABLE issues DROP COLUMN due_date;
```

Its down migration restores `in_inbox BOOLEAN DEFAULT 0`, `est_tokens INTEGER`,
and `due_date DATE` with the original
`kind != 'epic' OR est_tokens IS NULL` check.

`0012` up/down:

```sql
-- up
ALTER TABLE issues ADD COLUMN revision INTEGER NOT NULL DEFAULT 1;
-- down
ALTER TABLE issues DROP COLUMN revision;
```

Rules:

- The caller never supplies the resulting `revision`. Issue revision covers the
  complete aggregate. `issue.created` stores revision 1; `issue.updated`,
  `issue.transitioned`, `issue.ac_added`, `issue.label.added`,
  `issue.label.removed`, `link.added`, `link.removed`, `dev.started`,
  `dev.finished`, and `dev.blocked` advance each affected Issue exactly once.
  Link Events advance one or two Issue endpoints as applicable. The legacy dev
  families retain this invariant only until their Phase 21 deletion.
- Every update-like revision Effect has `op:"update"`, `table:"issues"`,
  `pk:{id}`, `values:{revision: prior+1}`, `prior:{revision: prior}`, and
  `preconditions:{revision: prior}`. Empty optional Effect maps are omitted,
  never `null`. The existing absolute-assignment projector therefore rebuilds
  deterministically without a new Effect operation.
- The `issues_enriched` view is recreated without the removed columns.
- No new table or Event family is introduced. The eleven existing families above
  lose the retired effect keys and carry the exact revision Effects defined in
  the Phase 20 contract section 7. Golden fixtures cover every family, including
  Issue/Issue and Issue/requirement link cases.
- Development journals containing the removed keys are not migrated. A
  projection applier meeting an unknown column key fails loudly. Repository
  fixtures are regenerated with the migration, and live/rebuild parity is proven
  against the regenerated fixtures.
- Deletion mapping for displaced surfaces is in the phase contract, section 5.

*Last updated: 2026-08-09 (Phase 20 accepted slice DDL amendment)*
