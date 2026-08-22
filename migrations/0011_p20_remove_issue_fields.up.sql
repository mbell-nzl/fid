PRAGMA foreign_keys=OFF;
PRAGMA legacy_alter_table=ON;
PRAGMA defer_foreign_keys=ON;

DROP VIEW issues_enriched;

CREATE TABLE issues_new (
  id           TEXT PRIMARY KEY,
  kind         TEXT NOT NULL,
  parent_id    TEXT REFERENCES issues(id),
  title        TEXT NOT NULL,
  workflow_id  TEXT NOT NULL REFERENCES workflows(id),
  status       TEXT NOT NULL,
  priority     INTEGER,
  sequence     INTEGER,
  spec_md_path TEXT,
  owner_role   TEXT,
  owner_actor  TEXT,
  created_at   DATETIME,
  updated_at   DATETIME,
  CHECK (kind IN ('epic', 'task', 'subtask')),
  CHECK (priority IS NULL OR (priority >= 0 AND priority <= 5)),
  CHECK (sequence IS NULL OR sequence > 0),
  UNIQUE (parent_id, sequence),
  FOREIGN KEY (workflow_id, status)
    REFERENCES workflow_states(workflow_id, state)
);

INSERT INTO issues_new (
  id, kind, parent_id, title, workflow_id, status, priority, sequence,
  spec_md_path, owner_role, owner_actor, created_at, updated_at
)
SELECT
  id, kind, parent_id, title, workflow_id, status, priority, sequence,
  spec_md_path, owner_role, owner_actor, created_at, updated_at
FROM issues;

DROP TABLE issues;
ALTER TABLE issues_new RENAME TO issues;

CREATE VIEW issues_enriched AS
WITH open_concerns AS (
  SELECT i.id AS issue_id, (
    SELECT COUNT(*) FROM concerns c WHERE c.issue_id = i.id AND c.status = 'open'
  ) AS n
  FROM issues i
), latest_triage AS (
  SELECT i.id AS issue_id, (
    SELECT te.verdict FROM triage_events te WHERE te.issue_id = i.id
    ORDER BY te.created_at DESC LIMIT 1
  ) AS verdict
  FROM issues i
)
SELECT i.*, (COALESCE(oc.n, 0) > 0) AS concerns_raised,
  COALESCE(oc.n, 0) AS open_concerns_count,
  CASE lt.verdict WHEN 'RESUME' THEN 'resumed' WHEN 'BLOCK' THEN 'blocked' ELSE NULL END AS triage_status
FROM issues i
LEFT JOIN open_concerns oc ON oc.issue_id = i.id
LEFT JOIN latest_triage lt ON lt.issue_id = i.id;

CREATE TRIGGER issues_parent_kind_guard_insert
BEFORE INSERT ON issues FOR EACH ROW WHEN
  (NEW.kind = 'epic' AND NEW.parent_id IS NOT NULL) OR
  (NEW.kind = 'task' AND NEW.parent_id IS NOT NULL AND NOT EXISTS (
    SELECT 1 FROM issues WHERE id = NEW.parent_id AND kind = 'epic'
  )) OR
  (NEW.kind = 'subtask' AND NOT EXISTS (
    SELECT 1 FROM issues WHERE id = NEW.parent_id AND kind = 'task'
  ))
BEGIN
  SELECT RAISE(ABORT, CASE
    WHEN NEW.kind = 'epic' THEN 'issues: epic parent_id must be NULL'
    WHEN NEW.kind = 'task' THEN 'issues: task parent_id must reference an epic or be NULL'
    ELSE 'issues: subtask parent_id must reference a task'
  END);
END;

CREATE TRIGGER issues_parent_kind_guard_update
BEFORE UPDATE OF parent_id ON issues FOR EACH ROW WHEN
  (NEW.kind = 'epic' AND NEW.parent_id IS NOT NULL) OR
  (NEW.kind = 'task' AND NEW.parent_id IS NOT NULL AND NOT EXISTS (
    SELECT 1 FROM issues WHERE id = NEW.parent_id AND kind = 'epic'
  )) OR
  (NEW.kind = 'subtask' AND NOT EXISTS (
    SELECT 1 FROM issues WHERE id = NEW.parent_id AND kind = 'task'
  ))
BEGIN
  SELECT RAISE(ABORT, CASE
    WHEN NEW.kind = 'epic' THEN 'issues: epic parent_id must be NULL'
    WHEN NEW.kind = 'task' THEN 'issues: task parent_id must reference an epic or be NULL'
    ELSE 'issues: subtask parent_id must reference a task'
  END);
END;

CREATE TRIGGER issues_kind_immutable
BEFORE UPDATE OF kind ON issues FOR EACH ROW WHEN NEW.kind != OLD.kind
BEGIN
  SELECT RAISE(ABORT, 'issues: kind is immutable');
END;

CREATE TRIGGER issues_sequence_scheduling_guard_insert
BEFORE INSERT ON issues FOR EACH ROW WHEN NEW.kind = 'task' AND (
  (NEW.parent_id IS NULL AND NEW.sequence IS NOT NULL) OR
  (EXISTS (SELECT 1 FROM issues WHERE id = NEW.parent_id AND kind = 'epic') AND (
    NOT EXISTS (SELECT 1 FROM epic_metadata WHERE issue_id = NEW.parent_id) OR
    EXISTS (SELECT 1 FROM epic_metadata WHERE issue_id = NEW.parent_id AND scheduling = 'sequential' AND NEW.sequence IS NULL) OR
    EXISTS (SELECT 1 FROM epic_metadata WHERE issue_id = NEW.parent_id AND scheduling = 'prioritised' AND NEW.sequence IS NOT NULL)
  ))
)
BEGIN
  SELECT RAISE(ABORT, CASE
    WHEN NEW.parent_id IS NULL THEN 'issues: unattached task must not have sequence'
    WHEN NOT EXISTS (SELECT 1 FROM epic_metadata WHERE issue_id = NEW.parent_id) THEN 'issues: task parent epic must have epic_metadata'
    WHEN EXISTS (SELECT 1 FROM epic_metadata WHERE issue_id = NEW.parent_id AND scheduling = 'sequential') THEN 'issues: task under sequential epic requires sequence'
    ELSE 'issues: task under prioritised epic must not have sequence'
  END);
END;

CREATE TRIGGER issues_sequence_scheduling_guard_update
BEFORE UPDATE OF parent_id, sequence ON issues FOR EACH ROW WHEN NEW.kind = 'task' AND (
  (NEW.parent_id IS NULL AND NEW.sequence IS NOT NULL) OR
  (EXISTS (SELECT 1 FROM issues WHERE id = NEW.parent_id AND kind = 'epic') AND (
    NOT EXISTS (SELECT 1 FROM epic_metadata WHERE issue_id = NEW.parent_id) OR
    EXISTS (SELECT 1 FROM epic_metadata WHERE issue_id = NEW.parent_id AND scheduling = 'sequential' AND NEW.sequence IS NULL) OR
    EXISTS (SELECT 1 FROM epic_metadata WHERE issue_id = NEW.parent_id AND scheduling = 'prioritised' AND NEW.sequence IS NOT NULL)
  ))
)
BEGIN
  SELECT RAISE(ABORT, CASE
    WHEN NEW.parent_id IS NULL THEN 'issues: unattached task must not have sequence'
    WHEN NOT EXISTS (SELECT 1 FROM epic_metadata WHERE issue_id = NEW.parent_id) THEN 'issues: task parent epic must have epic_metadata'
    WHEN EXISTS (SELECT 1 FROM epic_metadata WHERE issue_id = NEW.parent_id AND scheduling = 'sequential') THEN 'issues: task under sequential epic requires sequence'
    ELSE 'issues: task under prioritised epic must not have sequence'
  END);
END;

PRAGMA foreign_keys=ON;
PRAGMA legacy_alter_table=OFF;
PRAGMA defer_foreign_keys=OFF;
PRAGMA foreign_key_check;
