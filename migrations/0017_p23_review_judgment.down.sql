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
CREATE TEMP TABLE evidence_artifact_links_carry AS
  SELECT evidence_id, artifact_id, required_material FROM evidence_artifact_links;
DELETE FROM evidence_artifact_links;
DROP TABLE evidence_records;
ALTER TABLE evidence_records_old RENAME TO evidence_records;
INSERT INTO evidence_artifact_links (evidence_id, artifact_id, required_material)
  SELECT evidence_id, artifact_id, required_material FROM evidence_artifact_links_carry
  WHERE evidence_id IN (SELECT id FROM evidence_records);
DROP TABLE evidence_artifact_links_carry;
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
