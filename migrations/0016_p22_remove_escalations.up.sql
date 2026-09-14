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
