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
    'carried-forward','rejected','superseded'
  )),
  disposition_issue_id TEXT REFERENCES issues(id),
  supersedes_concern_id TEXT REFERENCES judgment_concerns(id),
  supersedes_finding_id TEXT REFERENCES review_findings(id),
  disposition_reason TEXT,
  disposed_actor_id TEXT,
  disposed_at TEXT,
  disposed_event_id TEXT,
  created_at TEXT NOT NULL,
  created_event_id TEXT NOT NULL UNIQUE,
  CHECK ((status = 'open'
          AND disposition_action IS NULL AND disposition_reason IS NULL
          AND disposed_actor_id IS NULL AND disposed_at IS NULL
          AND disposed_event_id IS NULL)
      OR (status = 'disposed'
          AND disposition_action IS NOT NULL AND disposition_reason IS NOT NULL
          AND disposed_actor_id IS NOT NULL AND disposed_at IS NOT NULL
          AND disposed_event_id IS NOT NULL)),
  CHECK ((disposition_action IS NULL
          AND disposition_issue_id IS NULL AND supersedes_concern_id IS NULL
          AND supersedes_finding_id IS NULL)
      OR (disposition_action = 'carried-forward'
          AND disposition_issue_id IS NOT NULL AND supersedes_concern_id IS NULL
          AND supersedes_finding_id IS NULL)
      OR (disposition_action = 'rejected'
          AND disposition_issue_id IS NULL AND supersedes_concern_id IS NULL
          AND supersedes_finding_id IS NULL)
      OR (disposition_action = 'superseded'
          AND disposition_issue_id IS NULL
          AND ((supersedes_concern_id IS NOT NULL AND supersedes_finding_id IS NULL)
            OR (supersedes_concern_id IS NULL AND supersedes_finding_id IS NOT NULL)))),
  CHECK (supersedes_concern_id IS NULL OR supersedes_concern_id <> id)
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
CREATE TEMP TABLE evidence_artifact_links_carry AS
  SELECT evidence_id, artifact_id, required_material FROM evidence_artifact_links;
DELETE FROM evidence_artifact_links;
DROP TABLE evidence_records;
ALTER TABLE evidence_records_new RENAME TO evidence_records;
INSERT INTO evidence_artifact_links (evidence_id, artifact_id, required_material)
  SELECT evidence_id, artifact_id, required_material FROM evidence_artifact_links_carry;
DROP TABLE evidence_artifact_links_carry;
CREATE UNIQUE INDEX evidence_records_one_receipt
  ON evidence_records(obligation_id, receipt_id) WHERE receipt_id IS NOT NULL;
CREATE UNIQUE INDEX evidence_records_one_review
  ON evidence_records(obligation_id, review_id) WHERE review_id IS NOT NULL;
CREATE UNIQUE INDEX evidence_records_one_outcome
  ON evidence_records(obligation_id, human_outcome_id) WHERE human_outcome_id IS NOT NULL;
