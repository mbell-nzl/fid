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
