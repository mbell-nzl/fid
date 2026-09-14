# Phase 23 DDL And Event Contract

> DRAFT, pending Human acceptance with `contracts.md` and `public-interface.md`.
> On acceptance this document is incorporated by `03_TRD_SCHEMA.md` section 8.4.
> No implementation may infer additional tables, columns, indexes, triggers,
> Event families, or effect rows beyond what is written here.

## 1. Development Baseline

Task 1 owns the following exact textual changes to already-shipped migrations,
plus the new migration `0017_p23_review_judgment`. This is a deliberate
pre-adoption baseline replacement requiring fresh init and rebuild; it is not
an upgrade of old fid data. Every other line of those files stays byte-identical
(DEC-178: shipped migration SQL and canonical slice DDL must agree exactly).

| File / location | Replace | With |
|---|---|---|
| `0013` `envelope_obligations.obligation_type` CHECK | `CHECK (obligation_type = 'gate-evidence')` | `CHECK (obligation_type IN ('gate-evidence','review','human-outcome'))` |
| `0013` `envelope_obligations.registration_id` | `registration_id TEXT NOT NULL,` | `registration_id TEXT,` |
| `0013` `envelope_obligations.registration_version` | `registration_version INTEGER NOT NULL,` | `registration_version INTEGER,` |
| `0013` `envelope_obligations` constraint block | (no registration discriminator) | add the CHECK below, immediately before `UNIQUE (envelope_id, revision, id)` |
| `0013` `work_claims.mode` CHECK | `CHECK (mode = 'mutation')` | `CHECK (mode IN ('mutation','review'))` |
| `0015` `stops.code` CHECK member list | the seven Phase 22 codes | the same seven plus `'review-rework-exhausted','protected-decision-required','human-outcome-required'` |
| `0015` `attention_items.code` CHECK member list | the seven Phase 22 codes | the same ten-member list used by `stops.code` |

The added `envelope_obligations` discriminator is exactly:

```sql
  CHECK ((obligation_type = 'gate-evidence'
          AND registration_id IS NOT NULL AND registration_version IS NOT NULL)
      OR (obligation_type IN ('review','human-outcome')
          AND registration_id IS NULL AND registration_version IS NULL)),
```

`work_claims` keeps its existing partial unique index unchanged: Phase 23
retains one active Claim per Issue across both modes, and mode-aware
concurrency is Phase 24 work. Migration 0013 and 0015 down behavior is
unchanged; neither gains a conversion obligation.

## 2. Migration 0017 Up

```sql
CREATE TABLE review_records (
  id TEXT PRIMARY KEY,
  run_id TEXT NOT NULL UNIQUE REFERENCES runs(id),
  issue_id TEXT NOT NULL REFERENCES issues(id),
  envelope_id TEXT NOT NULL,
  envelope_revision INTEGER NOT NULL,
  actor_id TEXT NOT NULL,
  profile TEXT NOT NULL,
  repository_json TEXT NOT NULL CHECK (json_valid(repository_json)),
  coverage_json TEXT NOT NULL CHECK (json_valid(coverage_json)),
  verdict TEXT NOT NULL CHECK (verdict IN ('pass','changes-required','blocked')),
  evidence_refs_json TEXT NOT NULL CHECK (json_valid(evidence_refs_json)),
  receipt_ids_json TEXT NOT NULL CHECK (json_valid(receipt_ids_json)),
  independence TEXT NOT NULL CHECK (independence = 'distinct-actor'),
  recorded_at TEXT NOT NULL,
  created_event_id TEXT NOT NULL UNIQUE,
  FOREIGN KEY (envelope_id, envelope_revision) REFERENCES envelope_revisions(envelope_id, revision)
);
CREATE INDEX review_records_by_issue ON review_records(issue_id, recorded_at, id);
CREATE INDEX review_records_by_envelope ON review_records(envelope_id, envelope_revision, recorded_at, id);

CREATE TABLE review_findings (
  id TEXT PRIMARY KEY,
  review_id TEXT NOT NULL REFERENCES review_records(id),
  subject_type TEXT NOT NULL CHECK (subject_type IN (
    'issue','envelope','run','review','finding','requirement'
  )),
  subject_id TEXT NOT NULL,
  scope_json TEXT NOT NULL CHECK (json_valid(scope_json)),
  severity TEXT NOT NULL CHECK (severity IN ('blocker','major','minor','nit')),
  summary TEXT NOT NULL,
  blocking INTEGER NOT NULL CHECK (blocking IN (0, 1)),
  created_at TEXT NOT NULL,
  created_event_id TEXT NOT NULL,
  CHECK ((severity IN ('blocker','major') AND blocking = 1)
      OR (severity IN ('minor','nit') AND blocking = 0))
);
CREATE INDEX review_findings_by_review ON review_findings(review_id, id);

CREATE TABLE finding_dispositions (
  id TEXT PRIMARY KEY,
  finding_id TEXT NOT NULL REFERENCES review_findings(id),
  action TEXT NOT NULL CHECK (action IN ('addressed','resolved','inapplicable')),
  actor_id TEXT NOT NULL,
  claim_id TEXT REFERENCES claims(id),
  qualifying_review_id TEXT REFERENCES review_records(id),
  reason TEXT NOT NULL,
  disposed_at TEXT NOT NULL,
  created_event_id TEXT NOT NULL,
  CHECK ((action = 'addressed'
          AND claim_id IS NOT NULL AND qualifying_review_id IS NULL)
      OR (action IN ('resolved','inapplicable')
          AND claim_id IS NULL AND qualifying_review_id IS NOT NULL)),
  UNIQUE (finding_id, claim_id)
);
CREATE UNIQUE INDEX finding_dispositions_one_terminal
  ON finding_dispositions(finding_id) WHERE action IN ('resolved','inapplicable');
CREATE INDEX finding_dispositions_by_finding
  ON finding_dispositions(finding_id, disposed_at, id);

CREATE TABLE judgment_concerns (
  id TEXT PRIMARY KEY,
  subject_type TEXT NOT NULL CHECK (subject_type IN (
    'issue','envelope','run','review','finding','requirement'
  )),
  subject_id TEXT NOT NULL,
  actor_id TEXT NOT NULL,
  scope_json TEXT NOT NULL CHECK (json_valid(scope_json)),
  severity TEXT NOT NULL CHECK (severity IN ('blocker','major','minor','nit')),
  summary TEXT NOT NULL,
  reasoning TEXT NOT NULL,
  references_json TEXT NOT NULL CHECK (json_valid(references_json)),
  status TEXT NOT NULL CHECK (status IN ('open','disposed')),
  disposition_action TEXT CHECK (disposition_action IS NULL OR disposition_action IN (
    'accepted','rejected','superseded','converted-to-finding'
  )),
  disposition_finding_id TEXT REFERENCES review_findings(id),
  disposition_reason TEXT,
  disposed_actor_id TEXT,
  disposed_at TEXT,
  disposed_event_id TEXT,
  created_at TEXT NOT NULL,
  created_event_id TEXT NOT NULL UNIQUE,
  CHECK ((status = 'open'
          AND disposition_action IS NULL AND disposition_finding_id IS NULL
          AND disposition_reason IS NULL AND disposed_actor_id IS NULL
          AND disposed_at IS NULL AND disposed_event_id IS NULL)
      OR (status = 'disposed'
          AND disposition_action IS NOT NULL AND disposition_reason IS NOT NULL
          AND disposed_actor_id IS NOT NULL AND disposed_at IS NOT NULL
          AND disposed_event_id IS NOT NULL)),
  CHECK (disposition_action IS NULL
      OR (disposition_action = 'converted-to-finding' AND disposition_finding_id IS NOT NULL)
      OR (disposition_action <> 'converted-to-finding' AND disposition_finding_id IS NULL))
);
CREATE INDEX judgment_concerns_by_subject
  ON judgment_concerns(subject_type, subject_id, status, created_at, id);

CREATE TABLE deviations (
  id TEXT PRIMARY KEY,
  subject_type TEXT NOT NULL CHECK (subject_type IN (
    'issue','envelope','run','review','finding','requirement'
  )),
  subject_id TEXT NOT NULL,
  actor_id TEXT NOT NULL,
  expected TEXT NOT NULL,
  observed TEXT NOT NULL,
  authority_actor_id TEXT NOT NULL,
  authority_reference TEXT,
  disposition TEXT NOT NULL CHECK (disposition IN ('accepted','rejected','pending-decision')),
  reasoning TEXT NOT NULL,
  references_json TEXT NOT NULL CHECK (json_valid(references_json)),
  supersedes_deviation_id TEXT REFERENCES deviations(id),
  created_at TEXT NOT NULL,
  created_event_id TEXT NOT NULL UNIQUE,
  CHECK (supersedes_deviation_id IS NULL OR supersedes_deviation_id <> id)
);
CREATE INDEX deviations_by_subject
  ON deviations(subject_type, subject_id, created_at, id);

CREATE TABLE decision_refs (
  id TEXT PRIMARY KEY,
  class TEXT NOT NULL,
  source_json TEXT NOT NULL CHECK (json_valid(source_json)),
  authority_actor_id TEXT NOT NULL,
  authority_capabilities_json TEXT NOT NULL CHECK (json_valid(authority_capabilities_json)),
  outcome TEXT NOT NULL CHECK (outcome IN ('approved','rejected','deferred')),
  scope_issue_id TEXT NOT NULL REFERENCES issues(id),
  created_at TEXT NOT NULL,
  created_event_id TEXT NOT NULL UNIQUE
);
CREATE INDEX decision_refs_by_class
  ON decision_refs(class, scope_issue_id, created_at, id);

CREATE TABLE human_outcomes (
  id TEXT PRIMARY KEY,
  actor_id TEXT NOT NULL,
  activation_id TEXT NOT NULL REFERENCES kit_activations(id),
  activation_version INTEGER NOT NULL CHECK (activation_version > 0),
  profile_ref TEXT NOT NULL,
  class TEXT NOT NULL,
  subject_type TEXT NOT NULL CHECK (subject_type IN (
    'issue','envelope','run','review','finding','requirement'
  )),
  subject_id TEXT NOT NULL,
  result TEXT NOT NULL CHECK (result IN (
    'approved','rejected','changes-required','acknowledged'
  )),
  scope_issue_id TEXT NOT NULL REFERENCES issues(id),
  reasoning TEXT NOT NULL,
  references_json TEXT NOT NULL CHECK (json_valid(references_json)),
  created_at TEXT NOT NULL,
  created_event_id TEXT NOT NULL UNIQUE
);
CREATE INDEX human_outcomes_by_class
  ON human_outcomes(class, scope_issue_id, created_at, id);

CREATE TABLE envelope_protected_decisions (
  envelope_id TEXT NOT NULL,
  revision INTEGER NOT NULL,
  ordinal INTEGER NOT NULL CHECK (ordinal >= 0),
  decision_class TEXT NOT NULL,
  required_capabilities_json TEXT NOT NULL CHECK (json_valid(required_capabilities_json)),
  affected_actions_json TEXT NOT NULL CHECK (json_valid(affected_actions_json)),
  PRIMARY KEY (envelope_id, revision, ordinal),
  UNIQUE (envelope_id, revision, decision_class),
  FOREIGN KEY (envelope_id, revision) REFERENCES envelope_revisions(envelope_id, revision)
);

CREATE TABLE requirement_refs (
  id TEXT PRIMARY KEY,
  requirement_id TEXT NOT NULL,
  source_json TEXT NOT NULL CHECK (json_valid(source_json)),
  status TEXT NOT NULL CHECK (status IN ('active','superseded','withdrawn')),
  superseded_by_id TEXT REFERENCES requirement_refs(id),
  created_at TEXT NOT NULL,
  created_event_id TEXT NOT NULL,
  CHECK ((status = 'superseded' AND superseded_by_id IS NOT NULL)
      OR (status IN ('active','withdrawn') AND superseded_by_id IS NULL)),
  CHECK (superseded_by_id IS NULL OR superseded_by_id <> id)
);
CREATE UNIQUE INDEX requirement_refs_one_active
  ON requirement_refs(requirement_id) WHERE status = 'active';
CREATE INDEX requirement_refs_by_requirement
  ON requirement_refs(requirement_id, created_at, id);

CREATE TABLE requirement_issue_links (
  requirement_id TEXT NOT NULL,
  issue_id TEXT NOT NULL REFERENCES issues(id),
  relation TEXT NOT NULL CHECK (relation IN ('implements','verifies','relates-to')),
  created_at TEXT NOT NULL,
  created_event_id TEXT NOT NULL,
  PRIMARY KEY (requirement_id, issue_id, relation)
);
CREATE INDEX requirement_issue_links_by_issue
  ON requirement_issue_links(issue_id, requirement_id);

CREATE TRIGGER requirement_issue_links_ref_guard_insert
BEFORE INSERT ON requirement_issue_links
FOR EACH ROW
WHEN NOT EXISTS (
  SELECT 1 FROM requirement_refs WHERE requirement_id = NEW.requirement_id
)
BEGIN
  SELECT RAISE(ABORT, 'requirement_issue_links: requirement_id must resolve to a recorded reference');
END;

ALTER TABLE envelope_obligations ADD COLUMN review_profile TEXT;
ALTER TABLE envelope_obligations ADD COLUMN outcome_class TEXT;
CREATE UNIQUE INDEX envelope_obligations_one_review_profile
  ON envelope_obligations(envelope_id, revision, review_profile)
  WHERE review_profile IS NOT NULL;
CREATE UNIQUE INDEX envelope_obligations_one_outcome_class
  ON envelope_obligations(envelope_id, revision, outcome_class)
  WHERE outcome_class IS NOT NULL;
CREATE UNIQUE INDEX envelope_obligations_one_gate_evidence
  ON envelope_obligations(envelope_id, revision)
  WHERE obligation_type = 'gate-evidence';

CREATE TABLE evidence_records_new (
  id TEXT PRIMARY KEY,
  obligation_id TEXT NOT NULL REFERENCES envelope_obligations(id),
  evidence_type TEXT NOT NULL CHECK (evidence_type IN ('receipt','review','human-outcome')),
  receipt_id TEXT REFERENCES execution_receipts(id),
  review_id TEXT REFERENCES review_records(id),
  human_outcome_id TEXT REFERENCES human_outcomes(id),
  producer TEXT NOT NULL CHECK (producer IN ('fid','activated-actor','human')),
  trust_source TEXT NOT NULL CHECK (trust_source IN (
    'fid-executed','activated-actor','activated-human'
  )),
  target_repository_json TEXT NOT NULL CHECK (json_valid(target_repository_json)),
  validity TEXT NOT NULL CHECK (validity IN ('valid','invalid')),
  created_at DATETIME NOT NULL,
  CHECK ((evidence_type = 'receipt' AND producer = 'fid'
          AND trust_source = 'fid-executed'
          AND receipt_id IS NOT NULL AND review_id IS NULL AND human_outcome_id IS NULL)
      OR (evidence_type = 'review' AND producer = 'activated-actor'
          AND trust_source = 'activated-actor'
          AND receipt_id IS NULL AND review_id IS NOT NULL AND human_outcome_id IS NULL)
      OR (evidence_type = 'human-outcome' AND producer = 'human'
          AND trust_source = 'activated-human'
          AND receipt_id IS NULL AND review_id IS NULL AND human_outcome_id IS NOT NULL))
);
INSERT INTO evidence_records_new (
  id, obligation_id, evidence_type, receipt_id, review_id, human_outcome_id,
  producer, trust_source, target_repository_json, validity, created_at)
SELECT id, obligation_id, evidence_type, receipt_id, NULL, NULL,
  producer, trust_source, target_repository_json, validity, created_at
FROM evidence_records;
DROP TABLE evidence_records;
ALTER TABLE evidence_records_new RENAME TO evidence_records;
CREATE UNIQUE INDEX evidence_records_one_receipt
  ON evidence_records(obligation_id, receipt_id) WHERE receipt_id IS NOT NULL;
CREATE UNIQUE INDEX evidence_records_one_review
  ON evidence_records(obligation_id, review_id) WHERE review_id IS NOT NULL;
CREATE UNIQUE INDEX evidence_records_one_outcome
  ON evidence_records(obligation_id, human_outcome_id) WHERE human_outcome_id IS NOT NULL;
```

`evidence_artifact_links.evidence_id` keeps its declared reference. The rebuild
runs inside the migration transaction with `PRAGMA legacy_alter_table` left at
its default; task 1 proves `PRAGMA foreign_key_check` is empty after the
rebuild, after `0018`, and after a full rebuild from the journal.

## 3. Migration 0017 Down

```sql
DROP INDEX evidence_records_one_outcome;
DROP INDEX evidence_records_one_review;
DROP INDEX evidence_records_one_receipt;
CREATE TABLE evidence_records_old (
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
INSERT INTO evidence_records_old
SELECT id, obligation_id, receipt_id, evidence_type, producer, trust_source,
  target_repository_json, validity, created_at
FROM evidence_records WHERE evidence_type = 'receipt';
DROP TABLE evidence_records;
ALTER TABLE evidence_records_old RENAME TO evidence_records;
DROP INDEX envelope_obligations_one_gate_evidence;
DROP INDEX envelope_obligations_one_outcome_class;
DROP INDEX envelope_obligations_one_review_profile;
ALTER TABLE envelope_obligations DROP COLUMN outcome_class;
ALTER TABLE envelope_obligations DROP COLUMN review_profile;
DROP TRIGGER requirement_issue_links_ref_guard_insert;
DROP TABLE requirement_issue_links;
DROP TABLE requirement_refs;
DROP TABLE envelope_protected_decisions;
DROP TABLE human_outcomes;
DROP TABLE decision_refs;
DROP TABLE deviations;
DROP TABLE judgment_concerns;
DROP TABLE finding_dispositions;
DROP TABLE review_findings;
DROP TABLE review_records;
```

Down is a development reset, not a data-preserving downgrade: rows whose
evidence type is not `receipt` are discarded and that loss is explicit.

## 4. Migration 0018 Up - V1 Judgment Removal

```sql
DROP TRIGGER activity_log_ref_integrity_guard_insert;
DROP TRIGGER links_to_ref_integrity_guard_insert;
DROP TRIGGER links_from_ref_integrity_guard_insert;
DROP TRIGGER links_endpoint_scope_guard_insert;
DROP TRIGGER links_relation_compatibility_guard_insert;
DROP VIEW issues_enriched;
DROP INDEX idx_concerns_issue_status;
DROP INDEX judgment_concerns_by_subject;

DELETE FROM links WHERE from_type <> 'issue' OR to_type <> 'issue';
DELETE FROM activity_log WHERE ref_type <> 'issue';

ALTER TABLE agent_runs DROP COLUMN completion_id;
DROP TABLE ac_verifications;
DROP TABLE concerns;
DROP TABLE task_completions;
DROP TABLE requirements;
DROP TABLE decisions;
ALTER TABLE judgment_concerns RENAME TO concerns;

DELETE FROM workflow_transition_requires
WHERE rule IN ('all_acs_verified','has_completion_report','concerns_addressed');

CREATE INDEX concerns_by_subject_status
  ON concerns(subject_type, subject_id, status, created_at, id);

CREATE VIEW issues_enriched AS
WITH open_concerns AS (
  SELECT
    i.id AS issue_id,
    (
      SELECT COUNT(*)
      FROM concerns c
      WHERE c.subject_type = 'issue' AND c.subject_id = i.id AND c.status = 'open'
    ) AS n
  FROM issues i
),
latest_triage AS (
  SELECT
    i.id AS issue_id,
    (
      SELECT te.verdict
      FROM triage_events te
      WHERE te.issue_id = i.id
      ORDER BY te.created_at DESC
      LIMIT 1
    ) AS verdict
  FROM issues i
)
SELECT
  i.*,
  (COALESCE(oc.n, 0) > 0) AS concerns_raised,
  COALESCE(oc.n, 0) AS open_concerns_count,
  CASE lt.verdict
    WHEN 'RESUME' THEN 'resumed'
    WHEN 'BLOCK' THEN 'blocked'
    ELSE NULL
  END AS triage_status
FROM issues i
LEFT JOIN open_concerns oc ON oc.issue_id = i.id
LEFT JOIN latest_triage lt ON lt.issue_id = i.id;

CREATE TRIGGER links_endpoint_scope_guard_insert
BEFORE INSERT ON links
FOR EACH ROW
WHEN NEW.from_type != 'issue' OR NEW.to_type != 'issue'
BEGIN
  SELECT RAISE(ABORT, 'links: endpoint types must be issue');
END;

CREATE TRIGGER links_from_ref_integrity_guard_insert
BEFORE INSERT ON links
FOR EACH ROW
WHEN NOT EXISTS (SELECT 1 FROM issues WHERE id = NEW.from_id)
BEGIN
  SELECT RAISE(ABORT, 'links: from endpoint must resolve to an existing row');
END;

CREATE TRIGGER links_to_ref_integrity_guard_insert
BEFORE INSERT ON links
FOR EACH ROW
WHEN NOT EXISTS (SELECT 1 FROM issues WHERE id = NEW.to_id)
BEGIN
  SELECT RAISE(ABORT, 'links: to endpoint must resolve to an existing row');
END;

CREATE TRIGGER links_relation_compatibility_guard_insert
BEFORE INSERT ON links
FOR EACH ROW
WHEN NEW.relation NOT IN (
  'relates_to','blocks','blocked_by','duplicates','duplicated_by',
  'supersedes','superseded_by','caused_by','causes',
  'discovered_from','discovered'
)
BEGIN
  SELECT RAISE(ABORT, 'links: relation/endpoint combination is not allowed');
END;

CREATE TRIGGER activity_log_ref_integrity_guard_insert
BEFORE INSERT ON activity_log
FOR EACH ROW
WHEN NEW.ref_type NOT IN ('issue', 'requirement', 'decision')
  OR (NEW.ref_type = 'issue' AND NOT EXISTS (
    SELECT 1 FROM issues WHERE id = NEW.ref_id
  ))
  OR (NEW.ref_type = 'requirement' AND NOT EXISTS (
    SELECT 1 FROM requirement_refs WHERE requirement_id = NEW.ref_id
  ))
  OR (NEW.ref_type = 'decision' AND NOT EXISTS (
    SELECT 1 FROM decision_refs WHERE id = NEW.ref_id
  ))
BEGIN
  SELECT RAISE(ABORT, 'activity_log: ref_type/ref_id must resolve to an existing row');
END;
```

`agent_runs.completion_id` is dropped because its foreign key targets the
removed `task_completions`; the rest of `agent_runs` is retained for Phase 25.
`triage_events` and its index are untouched: the triage-status projection and
its Phase 24 removal are not in this slice. `epic_metadata` and
`acceptance_criteria` are retained; `acceptance_criteria` already has no
verification role column, matching `03_TRD_SCHEMA.md` section 3.2.
Deleting every non-`issue` `activity_log` row is deliberate: the V1 requirement,
decision, and escalation subjects no longer resolve, and the replacement
`requirement` and `decision` activity subjects are written only by the new
Events.

## 5. Migration 0018 Down

```sql
DROP TRIGGER activity_log_ref_integrity_guard_insert;
DROP TRIGGER links_relation_compatibility_guard_insert;
DROP TRIGGER links_to_ref_integrity_guard_insert;
DROP TRIGGER links_from_ref_integrity_guard_insert;
DROP TRIGGER links_endpoint_scope_guard_insert;
DROP VIEW issues_enriched;
DROP INDEX concerns_by_subject_status;
ALTER TABLE concerns RENAME TO judgment_concerns;
CREATE INDEX judgment_concerns_by_subject
  ON judgment_concerns(subject_type, subject_id, status, created_at, id);
```

followed by the byte-identical `0002` definitions of `requirements`,
`decisions`, `task_completions`, `ac_verifications`, and `concerns`, the
byte-identical `0003` `idx_concerns_issue_status` and `issues_enriched`, the
byte-identical `0005`/`0016` link and activity trigger definitions, and an
`agent_runs` rebuild that restores its `completion_id` column. Down restores
empty V1 tables; it does not reconstruct deleted rows, and task 15 records that
explicitly.

## 6. Event Contract

All Phase 23 Events retain `fid.event.v2`. Nine Event types are added and one
is removed.

| Event type | Aggregate | Operation ID | Fence owner |
|---|---|---|---|
| `review.recorded` | `review` | `review.record` | task 4 |
| `finding.addressed` | `finding` | `finding.address` | task 5 |
| `concern.recorded` | `concern` | `concern.record` | task 8 |
| `concern.disposed` | `concern` | `concern.dispose` | task 8 |
| `deviation.recorded` | `deviation` | `deviation.record` | task 8 |
| `decision.referenced` | `decision` | `decision.record-reference` | task 9 |
| `human-outcome.recorded` | `human-outcome` | `human-outcome.record` | task 10 |
| `requirement.referenced` | `requirement` | `requirement.record-reference` | task 12 |
| `requirement.linked` | `requirement` | `requirement.link` | task 12 |

`requirement.synced` is removed by task 12 together with its effect rule,
operation mapping, projection writer, and fixtures.

### 6.1 Progressive Availability Fence

Task 1 adds `unavailableJudgmentEvents` to
`internal/journal/project_effects.go`, rejecting every one of the nine types
above until its domain exists, exactly as DEC-192 established for recovery
Events. The owning task named in the table deletes exactly its own key or keys
as part of implementing that Event and proves live-apply and replay parity for
it. Task 12 is the last owner and also removes the then-empty map and its guard
branch. Task 1's fence assertion is scheduled convergence for every later owner
and is never deleted as a parity assertion.

### 6.2 Current Event Effect Families

`currentEventRules` gains exactly these rows. Stage 1 is required, stage 0 is
optional, matching the existing encoding.

| Event | Required effects | Optional effects |
|---|---|---|
| `review.recorded` | `review_records/insert` | `review_findings/insert`, `finding_dispositions/insert`, `budget_samples/insert`, `stops/insert`, `attention_items/insert`, `attention_stop_links/insert`, `stop_resolutions/insert`, `attention_items/update`, `activity_log/activity`, `idempotency_results/insert` |
| `finding.addressed` | `finding_dispositions/insert` | `activity_log/activity`, `idempotency_results/insert` |
| `concern.recorded` | `judgment_concerns/insert` | `activity_log/activity`, `idempotency_results/insert` |
| `concern.disposed` | `judgment_concerns/update` | `activity_log/activity`, `idempotency_results/insert` |
| `deviation.recorded` | `deviations/insert` | `activity_log/activity`, `idempotency_results/insert` |
| `decision.referenced` | `decision_refs/insert` | `stop_resolutions/insert`, `attention_items/update`, `activity_log/activity`, `idempotency_results/insert` |
| `human-outcome.recorded` | `human_outcomes/insert` | `stop_resolutions/insert`, `attention_items/update`, `activity_log/activity`, `idempotency_results/insert` |
| `requirement.referenced` | `requirement_refs/insert` | `requirement_refs/update`, `activity_log/activity`, `idempotency_results/insert` |
| `requirement.linked` | `requirement_issue_links/insert` | `activity_log/activity`, `idempotency_results/insert` |

Amended existing rules:

| Event | Change | Owner |
|---|---|---|
| `envelope.created` | add optional `envelope_protected_decisions/insert` | task 9 |
| `envelope.activated` | add optional `stops/insert`, `attention_items/insert`, `attention_stop_links/insert` | task 9 |
| `evidence.linked` | unchanged effect rows; the record shape widens only | task 6 |

After task 15 renames `judgment_concerns` to `concerns`, the two concern effect
rules name `concerns`. Task 15 owns that single-word rule change and its
fixtures; no other effect row moves.

### 6.3 Typed Result Snapshots

DEC-216 typed idempotency snapshots extend to every Phase 23 mutation. A
committed denial for `review-rework-exhausted` or `protected-decision-required`
stores its complete typed snapshot - original Budget or protected-Decision
state, Stop and Attention IDs, and EventRef - and replays it without
re-deriving current state.

## 7. Parity And Corruption Proof

Task 1 and task 16 each prove, with fresh workspaces:

- live apply and full journal rebuild produce byte-identical projections for
  every Phase 23 table, including the `evidence_records` rebuild and the
  `judgment_concerns` rename;
- `PRAGMA foreign_key_check` and `PRAGMA integrity_check` are clean after
  `0017`, after `0018`, and after a rebuild;
- export, validate, and import round-trip every Phase 23 table and Event with
  semantically equivalent current state (PRD-REC-004);
- a corrupted or unknown Phase 23 Event payload is rejected by the current-Event
  validator rather than partially applied; and
- every removed table, index, trigger, view column, Event type, and rule name
  in section 4 is proven absent from the latest schema by a direct projection
  query, not by a source-text search alone.
