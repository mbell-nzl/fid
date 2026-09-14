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
