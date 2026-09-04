# Phase 22 DDL And Event Contract

> PROPOSED with DEC-180. Incorporated by `03_TRD_SCHEMA.md` section 8.3.
> No implementation may infer additional tables, fields, or Event families.

## 1. Development Baseline

Task 1 owns the following exact textual changes to migration
`0013_p21_operational_kernel.up.sql`, plus new migration `0015_p22_recovery`.
This deliberate pre-adoption baseline replacement requires fresh init/rebuild;
it is not an upgrade of old fid data. All other 0013 SQL stays byte-identical.

| Location | Replace | With |
|---|---|---|
| `system_capabilities.trust_source` CHECK member list | `'phase-22-unavailable','phase-23-unavailable','phase-25-unavailable'` | `'phase-22-unavailable','phase-23-unavailable','phase-24-unavailable','phase-25-unavailable'` |
| `kit_activations.status` CHECK | `('prepared','staged','committed','finalized','aborted')` | `('prepared','staged','committed','finalized','aborted','superseded')` |
| `work_claims.status` CHECK | `('active','released','expired')` | `('active','released','expired','superseded')` |
| `work_claims.release_reason` CHECK member list | `'completed','failed','operator-released'` | `'completed','failed','operator-released','interrupted','recovery-superseded','abandoned'` |
| `runs.status` CHECK | `('active','succeeded','failed','abandoned')` | `('active','succeeded','failed','abandoned','interrupted')` |

The historical section 8.2 of the canonical schema remains the Phase 21 record;
section 8.3 and this replacement table define the Phase 22 baseline. Migration
0013 down still drops its tables; it has no status-conversion obligation.

## 2. Migration 0015 Up

```sql
ALTER TABLE kit_activations ADD COLUMN version INTEGER NOT NULL DEFAULT 1 CHECK (version > 0);
CREATE TABLE run_runtime (
  run_id TEXT PRIMARY KEY REFERENCES runs(id),
  version INTEGER NOT NULL CHECK (version > 0),
  executor_id TEXT NOT NULL UNIQUE,
  obligation_id TEXT REFERENCES envelope_obligations(id)
);
CREATE TABLE run_observations (
  id TEXT PRIMARY KEY,
  run_id TEXT NOT NULL REFERENCES runs(id),
  kind TEXT NOT NULL CHECK (kind IN ('liveness','checkpoint')),
  actor_id TEXT NOT NULL,
  observed_at TEXT NOT NULL,
  observation_json TEXT NOT NULL CHECK (json_valid(observation_json)),
  created_event_id TEXT NOT NULL UNIQUE
);
CREATE INDEX run_observations_by_run ON run_observations(run_id, observed_at, id);
CREATE TABLE run_recoveries (
  id TEXT PRIMARY KEY,
  prior_run_id TEXT NOT NULL UNIQUE REFERENCES runs(id),
  successor_run_id TEXT NOT NULL UNIQUE REFERENCES runs(id),
  prior_claim_id TEXT NOT NULL REFERENCES claims(id),
  successor_claim_id TEXT NOT NULL UNIQUE REFERENCES claims(id),
  observation_id TEXT NOT NULL REFERENCES run_observations(id),
  reconciliation_json TEXT NOT NULL CHECK (json_valid(reconciliation_json)),
  adopted_paths_json TEXT NOT NULL CHECK (json_valid(adopted_paths_json)),
  authority_json TEXT NOT NULL CHECK (json_valid(authority_json)),
  reason TEXT NOT NULL,
  created_event_id TEXT NOT NULL UNIQUE,
  CHECK (prior_run_id <> successor_run_id),
  CHECK (prior_claim_id <> successor_claim_id)
);
CREATE TABLE budget_samples (
  id TEXT PRIMARY KEY,
  envelope_id TEXT NOT NULL,
  envelope_revision INTEGER NOT NULL,
  budget_name TEXT NOT NULL,
  run_id TEXT NOT NULL REFERENCES runs(id),
  actor_id TEXT NOT NULL,
  source TEXT NOT NULL CHECK (source IN ('fid','activated-harness')),
  source_sample_id TEXT NOT NULL,
  sample_kind TEXT NOT NULL CHECK (sample_kind IN ('delta','baseline')),
  previous_sample_id TEXT REFERENCES budget_samples(id),
  value INTEGER CHECK (value IS NULL OR value >= 0),
  measurement_ref_json TEXT CHECK (measurement_ref_json IS NULL OR json_valid(measurement_ref_json)),
  observed_at TEXT NOT NULL,
  created_event_id TEXT NOT NULL,
  UNIQUE (envelope_id, envelope_revision, budget_name, source, source_sample_id),
  FOREIGN KEY (envelope_id, envelope_revision, budget_name)
    REFERENCES envelope_budgets(envelope_id, revision, name)
);
CREATE INDEX budget_samples_by_budget
  ON budget_samples(envelope_id, envelope_revision, budget_name, observed_at, id);
CREATE TABLE stops (
  id TEXT PRIMARY KEY,
  subject_type TEXT NOT NULL CHECK (subject_type IN ('run','envelope','activation')),
  subject_id TEXT NOT NULL,
  code TEXT NOT NULL CHECK (code IN (
    'executor-liveness-unknown','recovery-state-conflict','budget-exhausted',
    'budget-data-unavailable','activation-incomplete','capability-unavailable',
    'operator-stop'
  )),
  blocked_actions_json TEXT NOT NULL CHECK (json_valid(blocked_actions_json)),
  details_json TEXT NOT NULL CHECK (json_valid(details_json)),
  resolution_condition_json TEXT NOT NULL CHECK (json_valid(resolution_condition_json)),
  resolver_capability TEXT NOT NULL,
  actor_id TEXT NOT NULL,
  raised_at TEXT NOT NULL,
  created_event_id TEXT NOT NULL
);
CREATE INDEX stops_by_subject ON stops(subject_type, subject_id, raised_at, id);
CREATE TABLE stop_resolutions (
  stop_id TEXT PRIMARY KEY REFERENCES stops(id),
  actor_id TEXT NOT NULL,
  fact_event_id TEXT NOT NULL,
  resolved_at TEXT NOT NULL,
  created_event_id TEXT NOT NULL
);
CREATE TABLE attention_items (
  id TEXT PRIMARY KEY,
  semantic_key TEXT NOT NULL,
  subject_type TEXT NOT NULL CHECK (subject_type IN ('run','envelope','activation')),
  subject_id TEXT NOT NULL,
  code TEXT NOT NULL CHECK (code IN (
    'executor-liveness-unknown','recovery-state-conflict','budget-exhausted',
    'budget-data-unavailable','activation-incomplete','capability-unavailable',
    'operator-stop'
  )),
  status TEXT NOT NULL CHECK (status IN ('open','resolved')),
  version INTEGER NOT NULL CHECK (version > 0),
  raised_at TEXT NOT NULL,
  resolution_event_id TEXT,
  CHECK ((status = 'open' AND resolution_event_id IS NULL)
    OR (status = 'resolved' AND resolution_event_id IS NOT NULL))
);
CREATE UNIQUE INDEX attention_one_open_semantic_key
  ON attention_items(semantic_key) WHERE status = 'open';
CREATE TABLE attention_stop_links (
  attention_id TEXT NOT NULL REFERENCES attention_items(id),
  stop_id TEXT NOT NULL UNIQUE REFERENCES stops(id),
  PRIMARY KEY (attention_id, stop_id)
);
CREATE TABLE activation_resources (
  activation_id TEXT NOT NULL REFERENCES kit_activations(id),
  component TEXT NOT NULL CHECK (component IN ('kit','overlay','provider','project','backend')),
  version INTEGER NOT NULL CHECK (version > 0),
  old_identity TEXT,
  new_identity TEXT NOT NULL,
  observed_identity TEXT,
  state TEXT NOT NULL CHECK (state IN ('pending','old','new','error')),
  evidence_ref TEXT,
  error_code TEXT,
  PRIMARY KEY (activation_id, component)
);
CREATE TABLE activation_recoveries (
  activation_id TEXT PRIMARY KEY REFERENCES kit_activations(id),
  direction TEXT NOT NULL CHECK (direction IN ('restore-prior','complete-forward')),
  actor_attribution TEXT NOT NULL,
  reason TEXT NOT NULL,
  selected_event_id TEXT NOT NULL UNIQUE,
  completed_event_id TEXT UNIQUE
);
```

Foreign keys are enabled. Polymorphic Stop subjects and Event references are
validated by the common accepted Event applier before applying any Effects,
live or replay; they cannot be checked by a SQLite FK to the journal. The
domain additionally validates same-Issue/revision ownership, capability, Run
chain acyclicity, and condition semantics. Direct SQL is not an authority path.
No budget total, Next Action, or second recovery-status aggregate is stored.

## 3. Migration 0015 Down

```sql
DROP TABLE activation_recoveries;
DROP TABLE activation_resources;
DROP TABLE attention_stop_links;
DROP TABLE attention_items;
DROP TABLE stop_resolutions;
DROP TABLE stops;
DROP TABLE budget_samples;
DROP TABLE run_recoveries;
DROP TABLE run_observations;
DROP TABLE run_runtime;
ALTER TABLE kit_activations DROP COLUMN version;
```

Full down/up is a disposable development reset. It destroys new projection
facts; reconstruction uses the same accepted journal and matching baseline.
It never translates Phase 22 history into Phase 21 history.

## 4. Migration 0016 - Escalation Removal

Task 16 applies this after replacement Stop/Attention conformance. Existing
escalation-bearing development data is not a migration input. Retained triage
and Concern rows must have null removed references; non-null data requires a
fresh disposable project, not silent conversion to a Stop.

Up is exactly:

```sql
DROP TRIGGER activity_log_ref_integrity_guard_insert;
ALTER TABLE triage_events DROP COLUMN esc_id;
ALTER TABLE concerns DROP COLUMN escalated_to_esc_id;
DELETE FROM workflow_transition_requires WHERE workflow_id = 'escalation-default';
DELETE FROM workflow_transitions WHERE workflow_id = 'escalation-default';
DELETE FROM workflow_states WHERE workflow_id = 'escalation-default';
DELETE FROM workflows WHERE id = 'escalation-default';
DROP TABLE escalations;
CREATE TRIGGER activity_log_ref_integrity_guard_insert
BEFORE INSERT ON activity_log
FOR EACH ROW
WHEN NEW.ref_type NOT IN ('issue', 'requirement', 'decision')
  OR (NEW.ref_type = 'issue' AND NOT EXISTS (
    SELECT 1 FROM issues WHERE id = NEW.ref_id
  ))
  OR (NEW.ref_type = 'requirement' AND NOT EXISTS (
    SELECT 1 FROM requirements WHERE id = NEW.ref_id
  ))
  OR (NEW.ref_type = 'decision' AND NOT EXISTS (
    SELECT 1 FROM decisions WHERE id = NEW.ref_id
  ))
BEGIN
  SELECT RAISE(ABORT, 'activity_log: ref_type/ref_id must resolve to an existing row');
END;
```

Down is exactly:

```sql
DROP TRIGGER activity_log_ref_integrity_guard_insert;
CREATE TABLE escalations (
  id TEXT PRIMARY KEY,
  status TEXT NOT NULL,
  trigger_class TEXT,
  raised_by_role TEXT,
  raised_at DATETIME,
  resolved_at DATETIME,
  resolution_dec TEXT REFERENCES decisions(id)
);
ALTER TABLE triage_events ADD COLUMN esc_id TEXT REFERENCES escalations(id);
ALTER TABLE concerns ADD COLUMN escalated_to_esc_id TEXT REFERENCES escalations(id);
CREATE TRIGGER activity_log_ref_integrity_guard_insert
BEFORE INSERT ON activity_log
FOR EACH ROW
WHEN NEW.ref_type NOT IN ('issue', 'requirement', 'decision', 'escalation')
  OR (NEW.ref_type = 'issue' AND NOT EXISTS (
    SELECT 1 FROM issues WHERE id = NEW.ref_id
  ))
  OR (NEW.ref_type = 'requirement' AND NOT EXISTS (
    SELECT 1 FROM requirements WHERE id = NEW.ref_id
  ))
  OR (NEW.ref_type = 'decision' AND NOT EXISTS (
    SELECT 1 FROM decisions WHERE id = NEW.ref_id
  ))
  OR (NEW.ref_type = 'escalation' AND NOT EXISTS (
    SELECT 1 FROM escalations WHERE id = NEW.ref_id
  ))
BEGIN
  SELECT RAISE(ABORT, 'activity_log: ref_type/ref_id must resolve to an existing row');
END;
```

Down restores empty table/column structure for the remainder of the down chain,
not deleted rows or the retired reference workflow. Full reset proceeds to zero
then up to 16. Partial down is not a supported product operation. Historical
0002/0005/0010 fixtures remain scoped to their exact migration, and are removed
at Phase 25 consolidation. Latest-schema tests must assert absence instead.

## 5. Closed Event Fixtures

All top-level fields and Effect grammar are Phase 21 `contracts.md` section 14.
Every new payload has required keys, explicit nulls, canonical ordering,
expected-version preconditions, captured prior update values, trusted Actor,
correlation to Run/activation where relevant, and operation/project idempotency.
An Event row below is the equivalent fixture contract: only its listed domain
Effects, plus the existing idempotency Effect, are allowed. Column keys are
exactly the DDL columns; JSON values use `public-interface.md` section 2.
All inserted entity IDs/times are bound consistently across Effects/results.

| Event type | Aggregate | Ordered domain Effects | Owner |
|---|---|---|---|
| `run.started` (replacement) | run | existing Run and immutable budget bindings; insert run_runtime with version 1 | 3 |
| `claim.expired` | claim | update active claims to expired, version +1, preserve expires_at | 3 |
| `run.observed` | run | insert liveness run_observations; increment run_runtime.version | 3 |
| `run.checkpointed` | run | insert checkpoint run_observations; increment run_runtime.version | 4 |
| `run.recovered` | predecessor run | terminalize active predecessor if needed; increment predecessor runtime version; dispose prior Claim if needed; insert successor Claim, Run, budget bindings, runtime; insert run_recoveries; optional fid intervention budget effects | 5, 7 |
| `stop.raised` | stop | insert stops; insert or reuse open attention_items; insert attention_stop_links | 2 |
| `stop.resolved` | stop | insert stop_resolutions only | 2 |
| `attention.resolved` | attention | update open attention_items to resolved/version +1/resolution Event | 2 |
| `run.budget-observed` | run | insert budget_samples; increment runtime version; append Stop/Attention effects if boundary crossed | 7 |
| `gate.execution-started` (extension) | attempt | existing Attempt; fid budget sample; optional Stop/Attention boundary effects | 7 |
| `gate.receipt-recorded` (extension) | receipt | existing Receipt/Artifact effects; fid budget samples; optional Stop/Attention boundary effects | 7, 8 |
| `run.execution-blocked` | run | unknown/exhausted budget Stop/Attention effects, no Attempt; increment runtime version | 7 |
| `run.closed` | run | terminal Run and repository; increment runtime version | 8 |
| `run.abandoned` | run | terminal Run, repository and runtime version; release Claim/version +1 | 8 |
| `activation.prepared` (extension) | activation | existing prepare plus prior activation identity and five pending activation_resources | 10 |
| `activation.resource-observed` | activation | expected-version update activation_resources | 10 |
| `activation.recovery-selected` | activation | insert activation_recoveries; raise activation-incomplete Stop/Attention | 10 |
| `activation.recovery-completed` | activation | record completion Event; finalize new/supersede prior or abort new/retain prior; disposition qualifying activation Stop | 10 |

Existing activation commit/finalize now validate resource agreement before
their existing Effects; ordinary finalize supersedes the predecessor in the
same Event. Readiness no longer equates two staged internal components with
cross-resource agreement. Budget and recovery Event extensions expire their
Phase 21 fixture exclusions in the task that changes each producer.

Activation prepare inserts version 1. Every subsequent activation Event also
updates kit_activations.version with the exact expected-version precondition;
resource updates check the resource version separately. Finalization that
supersedes a predecessor checks/increments its projection version in the same
Event. The predecessor's historical finalized Event is never edited.

`run.finished` and its caller-reported outcome payload have no final Phase 22
producer or accepted current fixture. Pure inspection, resume, export,
validation, import publication, and rebuild append no domain Event.

## 6. Parity And Corruption Proof

Each owning task tests its Events through live apply and replay using separate
temporary databases. Test malformed/unknown fields, bad cross-references,
duplicate immutable records, stale preconditions, sequence gaps, and rejected
effects with no partial projection. Task 15 extends every family to import and
fresh re-export; task 18 checks semantic JSON projections and original Event
identity/order against the same contract in the test-only fake.

Task 1 must prove fresh init, full disposable down/up, and migration-loaded
constraints on both the ordinary SQLite suite and the tagged built-binary
suite. Task 16 repeats that proof at schema 16. Schema values in rebuild,
export manifests, fixtures, and discovery derive from the active migration;
they cannot remain hard-coded to 14 or be prematurely set to 16.
