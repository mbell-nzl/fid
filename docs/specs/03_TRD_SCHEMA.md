# TARGET LOGICAL SCHEMA AND INVARIANTS - fid

> **Kernel chain superseded on 2026-11-10 by ADR-0014.** The golden path is
> built on a lean core: Issues, the journal and check results, plus holds,
> policy snapshots, review rounds and a person's decisions. The kernel chain
> this spec describes (Kit Activation, Work Envelope, Obligation, Claim, Run,
> Attempt, Evidence Record, Stop, Attention Item and their budgets and
> recovery) is removed after task 1.5a. Where they conflict, ADR-0014 and
> `RESET.md` win. Later tasks update this spec as they build.

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
| `idempotency_results` | operation scope, key, request digest, request field digests JSON, result entity/Event, created time | Same key/digest returns original result; different digest conflicts. Field digests are immutable after insertion. |
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
| `links` | from type/ID, relation, to type/ID, created Event | Both endpoints are Issues. The closed relation registry excludes `implements` and `implemented_by`; dedicated lifecycle fields are not duplicated. |
| `requirement_refs` | stable requirement ID, source reference/digest, status | Project files remain canonical. |
| `requirement_issue_links` | requirement ID, Issue ID, relation | Traceability only; relation is exactly `implements`, `verifies`, or `relates-to`. |

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
| `execution_receipts` | Receipt ID, execution ID, trust source, terminal status, exit/signal, timing, terminal repository identity, toolchain, sanitized excerpt/hash/size flags, canonical rejected-Artifact metadata | Immutable; one per execution. |
| `evidence_records` | Evidence ID, Obligation, evidence type/ref, producer, trust source, target identity, validity, created time | Append-only; no caller-authored pass flag. |
| `referenced_artifacts` | opaque/content ID, digest, media/type, retention class, availability, location policy, created/expiry time | Loss changes availability, not Evidence identity. |
| `evidence_artifact_links` | Evidence ID, Artifact ID, required-material flag | Composite primary key. |

### 3.7 Review and judgment

| Table | Minimum identity and content | Rules |
|---|---|---|
| `review_records` | Review ID, Run, Envelope/revision, exact repository identity, profile, coverage, verdict, evidence refs, created time | Immutable; Findings may be empty. |
| `review_findings` | Finding ID, Review, subject/scope, severity, summary, blocking flag | Created only with a Review Record. |
| `finding_dispositions` | Finding, Actor, addressed/resolved/inapplicable action, Claim or qualifying later Review, reason, time and Event | Addressed is nonterminal and requires a mutating Claim; only a later qualifying independent Review records resolved or inapplicable, with a mandatory reason for inapplicable. At most one terminal disposition; no Human Outcome substitutes for the Review (DEC-242). |
| `concerns` | ID, subject, raising Actor, scope, severity, summary, status, references, created time; disposition action, Issue/superseded Concern/superseded Finding reference, reason, disposing Actor, time and Event | Generic subject; not completion-linked. One irreversible `carried-forward`, `rejected`, or `superseded` disposition; the reference columns are mutually constrained by action. |
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
| Review/Judgment | Review recorded, Finding addressed/resolved, Concern recorded/disposed, Deviation recorded, Decision Reference recorded, Human Outcome recorded |
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
  `issue.transitioned`, `issue.ac-added`, `issue.label.added`,
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

### 8.2 Phase 21 - one-Issue operational walking skeleton

DEC-153 accepts this slice subject to Human approval of the complete planning
pack. The normative semantic/Event contract is
`../docs/phases/phase-21/contracts.md`; the exact wire contract is
`../docs/phases/phase-21/public-interface.md`; and the live removal authority
is `../docs/phases/phase-21/v1-convergence-inventory.md`. Migrations are deliberately split:
`0013_p21_operational_kernel` adds the generic replacement under the temporary
name `work_claims`; `0014_p21_retire_v1_authority_execution` runs only after
replacement conformance, drops the V1 collision, and renames it to `claims`.

`0013` up is exactly:

```sql
CREATE TABLE system_capabilities (
  capability_id TEXT PRIMARY KEY,
  schema_min INTEGER NOT NULL CHECK (schema_min > 0),
  schema_max INTEGER NOT NULL CHECK (schema_max >= schema_min),
  available INTEGER NOT NULL CHECK (available IN (0, 1)),
  limits_json TEXT NOT NULL CHECK (json_valid(limits_json)),
  trust_source TEXT NOT NULL CHECK (trust_source IN (
    'active-kit-activation','fid-executed','fid-git-observation','fid-journal',
    'phase-22-unavailable','phase-23-unavailable','phase-25-unavailable'
  ))
);

CREATE TABLE idempotency_results (
  operation_scope TEXT NOT NULL,
  idempotency_key TEXT NOT NULL CHECK (length(idempotency_key) BETWEEN 1 AND 128),
  request_digest TEXT NOT NULL,
  result_entity_id TEXT NOT NULL,
  result_event_id TEXT NOT NULL,
  aggregate_sequence INTEGER NOT NULL CHECK (aggregate_sequence > 0),
  result_json TEXT NOT NULL CHECK (json_valid(result_json)),
  created_at DATETIME NOT NULL,
  PRIMARY KEY (operation_scope, idempotency_key)
);

CREATE TABLE kit_activations (
  id TEXT PRIMARY KEY,
  prior_activation_id TEXT REFERENCES kit_activations(id),
  kit_identity TEXT NOT NULL,
  backend_identity TEXT NOT NULL,
  phase INTEGER NOT NULL CHECK (phase > 0),
  status TEXT NOT NULL CHECK (status IN ('prepared','staged','committed','finalized','aborted')),
  point_of_no_return INTEGER NOT NULL DEFAULT 0 CHECK (point_of_no_return IN (0, 1)),
  direction TEXT NOT NULL CHECK (direction = 'activate'),
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL
);
CREATE UNIQUE INDEX kit_activations_one_nonterminal
  ON kit_activations((1)) WHERE status IN ('prepared','staged','committed');
CREATE UNIQUE INDEX kit_activations_one_active
  ON kit_activations((1)) WHERE status = 'finalized';

CREATE TABLE activation_components (
  activation_id TEXT NOT NULL REFERENCES kit_activations(id),
  component_type TEXT NOT NULL CHECK (component_type IN ('authority-profile','gate-definitions')),
  component_id TEXT NOT NULL,
  component_version INTEGER NOT NULL CHECK (component_version > 0),
  content_digest TEXT NOT NULL,
  expected_prior_digest TEXT,
  status TEXT NOT NULL CHECK (status IN ('staged','committed','finalized','aborted')),
  evidence_json TEXT NOT NULL CHECK (json_valid(evidence_json)),
  error_code TEXT,
  PRIMARY KEY (activation_id, component_type, component_id, component_version)
);

CREATE TABLE activation_actor_profiles (
  activation_id TEXT NOT NULL REFERENCES kit_activations(id),
  actor_id TEXT NOT NULL,
  actor_type TEXT NOT NULL CHECK (actor_type IN ('agent','human','service')),
  producer TEXT NOT NULL,
  provenance_ref TEXT NOT NULL,
  session_ref TEXT,
  capabilities_json TEXT NOT NULL CHECK (json_valid(capabilities_json)),
  PRIMARY KEY (activation_id, actor_id)
);

CREATE TABLE activation_gate_definitions (
  activation_id TEXT NOT NULL REFERENCES kit_activations(id),
  definition_id TEXT NOT NULL,
  definition_version INTEGER NOT NULL CHECK (definition_version > 0),
  definition_digest TEXT NOT NULL,
  definition_json TEXT NOT NULL CHECK (json_valid(definition_json)),
  PRIMARY KEY (activation_id, definition_id, definition_version),
  UNIQUE (activation_id, definition_digest)
);

CREATE TABLE gate_registrations (
  id TEXT NOT NULL,
  version INTEGER NOT NULL CHECK (version > 0),
  activation_id TEXT NOT NULL,
  definition_id TEXT NOT NULL,
  definition_version INTEGER NOT NULL,
  definition_digest TEXT NOT NULL,
  binding_digest TEXT NOT NULL,
  project_id TEXT NOT NULL,
  worktree_id TEXT NOT NULL,
  scope_json TEXT NOT NULL CHECK (json_valid(scope_json)),
  policy_json TEXT NOT NULL CHECK (json_valid(policy_json)),
  authorization_provenance TEXT NOT NULL,
  created_event_id TEXT NOT NULL,
  created_at DATETIME NOT NULL,
  PRIMARY KEY (id, version),
  UNIQUE (definition_digest, binding_digest),
  FOREIGN KEY (activation_id, definition_id, definition_version)
    REFERENCES activation_gate_definitions(activation_id, definition_id, definition_version)
);

CREATE TABLE work_envelopes (
  id TEXT PRIMARY KEY,
  current_revision INTEGER NOT NULL CHECK (current_revision = 1),
  status TEXT NOT NULL CHECK (status IN ('draft','active','completed')),
  created_at DATETIME NOT NULL
);

CREATE TABLE envelope_revisions (
  envelope_id TEXT NOT NULL REFERENCES work_envelopes(id),
  revision INTEGER NOT NULL CHECK (revision = 1),
  prior_revision INTEGER CHECK (prior_revision IS NULL),
  outcome TEXT NOT NULL,
  non_goals_json TEXT NOT NULL CHECK (json_valid(non_goals_json)),
  risk_result_json TEXT NOT NULL CHECK (json_valid(risk_result_json)),
  repository_scope_json TEXT NOT NULL CHECK (json_valid(repository_scope_json)),
  status TEXT NOT NULL CHECK (status IN ('draft','active','completed')),
  created_actor_id TEXT NOT NULL,
  approved_actor_id TEXT,
  created_at DATETIME NOT NULL,
  activated_at DATETIME,
  completed_at DATETIME,
  PRIMARY KEY (envelope_id, revision)
);

CREATE TABLE envelope_sources (
  envelope_id TEXT NOT NULL,
  revision INTEGER NOT NULL,
  ordinal INTEGER NOT NULL CHECK (ordinal >= 0),
  source_json TEXT NOT NULL CHECK (json_valid(source_json)),
  PRIMARY KEY (envelope_id, revision, ordinal),
  FOREIGN KEY (envelope_id, revision) REFERENCES envelope_revisions(envelope_id, revision)
);

CREATE TABLE envelope_issue_selection (
  envelope_id TEXT NOT NULL,
  revision INTEGER NOT NULL,
  issue_id TEXT NOT NULL REFERENCES issues(id),
  primary_issue INTEGER NOT NULL CHECK (primary_issue = 1),
  PRIMARY KEY (envelope_id, revision, issue_id),
  UNIQUE (envelope_id, revision, primary_issue),
  FOREIGN KEY (envelope_id, revision) REFERENCES envelope_revisions(envelope_id, revision)
);

CREATE TABLE envelope_authority_grants (
  envelope_id TEXT NOT NULL,
  revision INTEGER NOT NULL,
  ordinal INTEGER NOT NULL CHECK (ordinal >= 0),
  actor_selector_json TEXT NOT NULL CHECK (json_valid(actor_selector_json)),
  actions_json TEXT NOT NULL CHECK (json_valid(actions_json)),
  scope_json TEXT NOT NULL CHECK (json_valid(scope_json)),
  PRIMARY KEY (envelope_id, revision, ordinal),
  FOREIGN KEY (envelope_id, revision) REFERENCES envelope_revisions(envelope_id, revision)
);

CREATE TABLE envelope_obligations (
  id TEXT PRIMARY KEY,
  envelope_id TEXT NOT NULL,
  revision INTEGER NOT NULL,
  obligation_type TEXT NOT NULL CHECK (obligation_type = 'gate-evidence'),
  registration_id TEXT NOT NULL,
  registration_version INTEGER NOT NULL,
  obligation_json TEXT NOT NULL CHECK (json_valid(obligation_json)),
  UNIQUE (envelope_id, revision, id),
  FOREIGN KEY (envelope_id, revision) REFERENCES envelope_revisions(envelope_id, revision),
  FOREIGN KEY (registration_id, registration_version) REFERENCES gate_registrations(id, version)
);

CREATE TABLE envelope_budgets (
  envelope_id TEXT NOT NULL,
  revision INTEGER NOT NULL,
  name TEXT NOT NULL,
  budget_json TEXT NOT NULL CHECK (json_valid(budget_json)),
  PRIMARY KEY (envelope_id, revision, name),
  FOREIGN KEY (envelope_id, revision) REFERENCES envelope_revisions(envelope_id, revision)
);

CREATE TABLE work_claims (
  id TEXT PRIMARY KEY,
  issue_id TEXT NOT NULL REFERENCES issues(id),
  envelope_id TEXT NOT NULL,
  envelope_revision INTEGER NOT NULL,
  actor_id TEXT NOT NULL,
  worktree_id TEXT NOT NULL,
  scope_json TEXT NOT NULL CHECK (json_valid(scope_json)),
  mode TEXT NOT NULL CHECK (mode = 'mutation'),
  version INTEGER NOT NULL CHECK (version > 0),
  status TEXT NOT NULL CHECK (status IN ('active','released','expired')),
  acquired_at DATETIME NOT NULL,
  expires_at DATETIME NOT NULL,
  released_at DATETIME,
  release_reason TEXT CHECK (release_reason IS NULL OR release_reason IN (
    'completed','failed','operator-released'
  )),
  FOREIGN KEY (envelope_id, envelope_revision) REFERENCES envelope_revisions(envelope_id, revision)
);
CREATE UNIQUE INDEX work_claims_one_active_mutation
  ON work_claims(issue_id) WHERE status = 'active';

CREATE TABLE runs (
  id TEXT PRIMARY KEY,
  issue_id TEXT NOT NULL REFERENCES issues(id),
  envelope_id TEXT NOT NULL,
  envelope_revision INTEGER NOT NULL,
  claim_id TEXT NOT NULL REFERENCES work_claims(id),
  claim_version INTEGER NOT NULL,
  actor_id TEXT NOT NULL,
  start_repository_json TEXT NOT NULL CHECK (json_valid(start_repository_json)),
  terminal_repository_json TEXT CHECK (terminal_repository_json IS NULL OR json_valid(terminal_repository_json)),
  status TEXT NOT NULL CHECK (status IN ('active','succeeded','failed','abandoned')),
  started_at DATETIME NOT NULL,
  ended_at DATETIME,
  FOREIGN KEY (envelope_id, envelope_revision) REFERENCES envelope_revisions(envelope_id, revision)
);
CREATE UNIQUE INDEX runs_one_active_per_claim
  ON runs(claim_id) WHERE status = 'active';

CREATE TABLE run_budget_bindings (
  run_id TEXT NOT NULL REFERENCES runs(id),
  name TEXT NOT NULL,
  budget_json TEXT NOT NULL CHECK (json_valid(budget_json)),
  PRIMARY KEY (run_id, name)
);

CREATE TABLE execution_attempts (
  id TEXT PRIMARY KEY,
  registration_id TEXT NOT NULL,
  registration_version INTEGER NOT NULL,
  obligation_id TEXT NOT NULL REFERENCES envelope_obligations(id),
  run_id TEXT NOT NULL REFERENCES runs(id),
  actor_id TEXT NOT NULL,
  start_repository_json TEXT NOT NULL CHECK (json_valid(start_repository_json)),
  resolved_policy_digest TEXT NOT NULL,
  producer TEXT NOT NULL CHECK (producer = 'fid'),
  started_at DATETIME NOT NULL,
  FOREIGN KEY (registration_id, registration_version) REFERENCES gate_registrations(id, version)
);

CREATE TABLE referenced_artifacts (
  id TEXT PRIMARY KEY,
  digest TEXT NOT NULL UNIQUE,
  media_type TEXT NOT NULL,
  size_bytes INTEGER NOT NULL CHECK (size_bytes >= 0),
  retention_class TEXT NOT NULL CHECK (retention_class = 'evidence'),
  availability TEXT NOT NULL CHECK (availability IN ('available','missing')),
  location_policy TEXT NOT NULL CHECK (location_policy = 'content-addressed'),
  created_at DATETIME NOT NULL,
  expires_at DATETIME
);

CREATE TABLE execution_receipts (
  id TEXT PRIMARY KEY,
  attempt_id TEXT NOT NULL UNIQUE REFERENCES execution_attempts(id),
  trust_source TEXT NOT NULL CHECK (trust_source = 'fid-executed'),
  terminal_status TEXT NOT NULL CHECK (terminal_status IN ('succeeded','failed','timed-out','signaled','spawn-failed')),
  exit_code INTEGER,
  signal TEXT,
  timed_out INTEGER NOT NULL CHECK (timed_out IN (0, 1)),
  started_at DATETIME NOT NULL,
  ended_at DATETIME NOT NULL,
  terminal_repository_json TEXT NOT NULL CHECK (json_valid(terminal_repository_json)),
  toolchain_json TEXT NOT NULL CHECK (json_valid(toolchain_json)),
  output_excerpt TEXT NOT NULL,
  output_digest TEXT NOT NULL,
  output_bytes INTEGER NOT NULL CHECK (output_bytes >= 0),
  output_truncated INTEGER NOT NULL CHECK (output_truncated IN (0, 1)),
  redaction_applied INTEGER NOT NULL CHECK (redaction_applied IN (0, 1)),
  rejected_artifacts_json TEXT NOT NULL
    CHECK (json_valid(rejected_artifacts_json) AND json_type(rejected_artifacts_json) = 'array'),
  CHECK ((terminal_status = 'spawn-failed' AND exit_code IS NULL AND signal IS NULL AND timed_out = 0)
      OR (terminal_status = 'timed-out' AND timed_out = 1)
      OR (terminal_status = 'signaled' AND exit_code IS NULL AND signal IS NOT NULL AND timed_out = 0)
      OR (terminal_status IN ('succeeded','failed') AND exit_code IS NOT NULL AND signal IS NULL AND timed_out = 0))
);

CREATE TABLE receipt_artifacts (
  receipt_id TEXT NOT NULL REFERENCES execution_receipts(id),
  artifact_id TEXT NOT NULL REFERENCES referenced_artifacts(id),
  declared_path TEXT NOT NULL,
  PRIMARY KEY (receipt_id, artifact_id)
);

CREATE TABLE evidence_records (
  id TEXT PRIMARY KEY,
  obligation_id TEXT NOT NULL REFERENCES envelope_obligations(id),
  receipt_id TEXT NOT NULL REFERENCES execution_receipts(id),
  evidence_type TEXT NOT NULL CHECK (evidence_type = 'receipt'),
  producer TEXT NOT NULL CHECK (producer = 'fid'),
  trust_source TEXT NOT NULL CHECK (trust_source = 'fid-executed'),
  target_repository_json TEXT NOT NULL CHECK (json_valid(target_repository_json)),
  validity TEXT NOT NULL CHECK (validity IN ('valid','invalid')),
  created_at DATETIME NOT NULL,
  UNIQUE (obligation_id, receipt_id)
);

CREATE TABLE evidence_artifact_links (
  evidence_id TEXT NOT NULL REFERENCES evidence_records(id),
  artifact_id TEXT NOT NULL REFERENCES referenced_artifacts(id),
  required_material INTEGER NOT NULL CHECK (required_material IN (0, 1)),
  PRIMARY KEY (evidence_id, artifact_id)
);
```

`0013` down is exactly the following reverse-order destructive development
rollback. It is valid only after `0014` has first been rolled down:

```sql
DROP TABLE evidence_artifact_links;
DROP TABLE evidence_records;
DROP TABLE receipt_artifacts;
DROP TABLE execution_receipts;
DROP TABLE referenced_artifacts;
DROP TABLE execution_attempts;
DROP TABLE run_budget_bindings;
DROP INDEX runs_one_active_per_claim;
DROP TABLE runs;
DROP INDEX work_claims_one_active_mutation;
DROP TABLE work_claims;
DROP TABLE envelope_budgets;
DROP TABLE envelope_obligations;
DROP TABLE envelope_authority_grants;
DROP TABLE envelope_issue_selection;
DROP TABLE envelope_sources;
DROP TABLE envelope_revisions;
DROP TABLE work_envelopes;
DROP TABLE gate_registrations;
DROP TABLE activation_gate_definitions;
DROP TABLE activation_actor_profiles;
DROP TABLE activation_components;
DROP INDEX kit_activations_one_active;
DROP INDEX kit_activations_one_nonterminal;
DROP TABLE kit_activations;
DROP TABLE idempotency_results;
DROP TABLE system_capabilities;
```

`0014` up is exactly:

```sql
DROP TABLE roles;
DROP TABLE claims;
DROP TABLE gate_runs;

DROP INDEX work_claims_one_active_mutation;
ALTER TABLE work_claims RENAME TO claims;
CREATE UNIQUE INDEX claims_one_active_mutation
  ON claims(issue_id) WHERE status = 'active';

DELETE FROM workflow_transition_requires
WHERE rule IN ('all_acs_verified', 'has_completion_report', 'concerns_addressed')
  AND workflow_id IN (
    SELECT id
    FROM workflows
    WHERE applies_to IN ('epic', 'task', 'subtask')
  );
```

`0014` down is an empty-state development rollback, not data reconstruction:

```sql
DROP INDEX claims_one_active_mutation;
ALTER TABLE claims RENAME TO work_claims;
CREATE UNIQUE INDEX work_claims_one_active_mutation
  ON work_claims(issue_id) WHERE status = 'active';

CREATE TABLE roles (
  name        TEXT PRIMARY KEY,
  description TEXT,
  is_human    BOOLEAN DEFAULT 0,
  active      BOOLEAN DEFAULT 1,
  created_at  DATETIME
);

CREATE TABLE claims (
  ref_type    TEXT NOT NULL,
  ref_id      TEXT NOT NULL,
  claimed_by  TEXT NOT NULL,
  claimed_at  DATETIME NOT NULL,
  ttl_seconds INTEGER NOT NULL DEFAULT 1800,
  note        TEXT,
  PRIMARY KEY (ref_type, ref_id)
);

CREATE TABLE gate_runs (
  id               TEXT PRIMARY KEY,
  issue_id         TEXT NOT NULL REFERENCES issues(id),
  attempt_number   INTEGER NOT NULL,
  agent_run_id     TEXT REFERENCES agent_runs(id),
  recorded_by_role TEXT NOT NULL,
  command          TEXT NOT NULL,
  exit_code        INTEGER NOT NULL,
  output_excerpt   TEXT,
  ran_at           DATETIME NOT NULL,
  CHECK (attempt_number > 0)
);
```

Event schema v2 and its canonical fixtures are in the phase contract section
14. The final `project.initialized` fixture has the `fid-bootstrap` service
Actor, project aggregate sequence 1, null repository/idempotency, and zero
domain Effects. It never inserts a role or policy row. Each other listed
family gets one success fixture plus its independent conflict/absence fixture.
Live apply and rebuild must produce byte-identical canonical reads. Phase 21
export create/validate covers v2; import parity is explicitly excluded.
Migration `0014` does not drop `task_completions`, `agent_runs`,
`ac_verifications`, or `concerns`; the convergence inventory freezes them
without writers or protected-semantic readers and names Phase 23/25 ownership.

### 8.3 Phase 22 - accepted failure-and-recovery amendment

Status: delivered and accepted under DEC-188/220. The exact SQL and closed
Event fixture contract are incorporated here from
`../docs/phases/phase-22/schema.md`. Its sections 1-3 define the precise 0013
development-baseline replacements and 0015 up/down SQL; section 4 defines
replacement-bound 0016 escalation removal/reset; sections 5-6 define Event
Effects, references, rejection, and live/rebuild/import proof obligations.

The associated JSON values and mutation semantics are
`../docs/phases/phase-22/public-interface.md` and `contracts.md`. Their explicit
replacements are the delivered Phase 22 baseline, and shipped SQL must match
this incorporated source. Section 8.2 retains historical Phase 21 SQL; these
pre-adoption replacements do not authorize an old-data migration. DEC-195
adds the accepted cancelled Receipt member and composite CHECK clause to 0013;
DEC-216 specifies complete typed current-family idempotency snapshots, including
pending Attempts, terminal Receipts and committed denials.

DEC-184's QA refinement makes run_recoveries.observation_id non-null and gives
attention_items.code the same seven-code CHECK as stops.code. Disposition Event
references are distinct from their underlying qualifying fact references.
DEC-187 keeps recovery-export binding in the kit coordinator's activation
record; no such column is implied in fid activation_recoveries.

### 8.4 Phase 23 - review and judgment baseline

Status: accepted under DEC-222 and delivered incrementally by the Phase 23 task
pack. The exact development-baseline replacements to migrations 0013 and 0015,
the full migration 0017 up/down DDL, the evidence-artifact-link-preserving
rebuild order, and the nine Event family Effect rules are incorporated by
reference from `../docs/phases/phase-23/schema.md` sections 1-3 and 6-7.
Those sections are the physical DDL authority for this slice; shipped migration
SQL and the incorporated DDL must agree exactly. The baseline remains
pre-adoption: clean initialization and rebuild are required, with no old-data
converter implied.

Migration 0018 (V1 judgment removal) is incorporated from sections 4-5 of the
same file. Its down migration restores the schema-17 state exactly (DEC-264);
no migration from 0001 to 0017 is edited. Journals written by an earlier
Phase 23 development binary are not supported: where a stored result no longer
matches the current closed shape, rebuild stops with `journal-replay-stopped`
and leaves the source journal unchanged (DEC-270, ADR-0011).

### 8.5 Phase 24 - idempotency field digest storage

Migration 0019 adds one column to `idempotency_results`, after the schema-18
baseline (DEC-275):

```sql
ALTER TABLE idempotency_results ADD COLUMN request_field_digests_json TEXT NOT NULL CHECK (json_valid(request_field_digests_json));
```

The column holds a canonical JSON object mapping public input names to SHA-256
digests. Every idempotency insert Effect carries it. The receipt-recorded
update guards the stored value and never changes it (DEC-278). Migration 0019
down drops only this column and restores schema 18 byte for byte. The active
projection and export manifest schema are 19; projection meta format remains 5.
Earlier journals that lack the field are unsupported and fail closed (DEC-270).

*Last updated: 2026-10-04 (Phase 24 task 1 storage baseline; DEC-275/278)*
