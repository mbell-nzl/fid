CREATE TRIGGER issues_parent_kind_guard_insert
BEFORE INSERT ON issues
FOR EACH ROW
WHEN
  (NEW.kind = 'epic' AND NEW.parent_id IS NOT NULL)
  OR (
    NEW.kind = 'task'
    AND NEW.parent_id IS NOT NULL
    AND NOT EXISTS (
      SELECT 1
      FROM issues
      WHERE id = NEW.parent_id AND kind = 'epic'
    )
  )
  OR (
    NEW.kind = 'subtask'
    AND NOT EXISTS (
      SELECT 1
      FROM issues
      WHERE id = NEW.parent_id AND kind = 'task'
    )
  )
BEGIN
  SELECT RAISE(
    ABORT,
    CASE
      WHEN NEW.kind = 'epic' THEN 'issues: epic parent_id must be NULL'
      WHEN NEW.kind = 'task' THEN 'issues: task parent_id must reference an epic or be NULL'
      ELSE 'issues: subtask parent_id must reference a task'
    END
  );
END;

CREATE TRIGGER issues_parent_kind_guard_update
BEFORE UPDATE OF parent_id ON issues
FOR EACH ROW
WHEN
  (NEW.kind = 'epic' AND NEW.parent_id IS NOT NULL)
  OR (
    NEW.kind = 'task'
    AND NEW.parent_id IS NOT NULL
    AND NOT EXISTS (
      SELECT 1
      FROM issues
      WHERE id = NEW.parent_id AND kind = 'epic'
    )
  )
  OR (
    NEW.kind = 'subtask'
    AND NOT EXISTS (
      SELECT 1
      FROM issues
      WHERE id = NEW.parent_id AND kind = 'task'
    )
  )
BEGIN
  SELECT RAISE(
    ABORT,
    CASE
      WHEN NEW.kind = 'epic' THEN 'issues: epic parent_id must be NULL'
      WHEN NEW.kind = 'task' THEN 'issues: task parent_id must reference an epic or be NULL'
      ELSE 'issues: subtask parent_id must reference a task'
    END
  );
END;

CREATE TRIGGER acceptance_criteria_issue_kind_guard_insert
BEFORE INSERT ON acceptance_criteria
FOR EACH ROW
WHEN NOT EXISTS (
  SELECT 1
  FROM issues
  WHERE id = NEW.issue_id AND kind IN ('task', 'subtask')
)
BEGIN
  SELECT RAISE(ABORT, 'acceptance_criteria: issue_id must reference a task or subtask');
END;

CREATE TRIGGER task_completions_issue_kind_guard_insert
BEFORE INSERT ON task_completions
FOR EACH ROW
WHEN NOT EXISTS (
  SELECT 1
  FROM issues
  WHERE id = NEW.issue_id AND kind IN ('task', 'subtask')
)
BEGIN
  SELECT RAISE(ABORT, 'task_completions: issue_id must reference a task or subtask');
END;

CREATE TRIGGER epic_metadata_issue_kind_guard_insert
BEFORE INSERT ON epic_metadata
FOR EACH ROW
WHEN NOT EXISTS (
  SELECT 1
  FROM issues
  WHERE id = NEW.issue_id AND kind = 'epic'
)
BEGIN
  SELECT RAISE(ABORT, 'epic_metadata: issue_id must reference an epic');
END;

CREATE TRIGGER triage_events_epic_kind_guard_insert
BEFORE INSERT ON triage_events
FOR EACH ROW
WHEN NOT EXISTS (
  SELECT 1
  FROM issues
  WHERE id = NEW.epic_id AND kind = 'epic'
)
BEGIN
  SELECT RAISE(ABORT, 'triage_events: epic_id must reference an epic');
END;

CREATE TRIGGER decisions_epic_kind_guard_insert
BEFORE INSERT ON decisions
FOR EACH ROW
WHEN NEW.epic_id IS NOT NULL
  AND NOT EXISTS (
    SELECT 1
    FROM issues
    WHERE id = NEW.epic_id AND kind = 'epic'
  )
BEGIN
  SELECT RAISE(ABORT, 'decisions: epic_id must reference an epic when set');
END;

CREATE TRIGGER issues_kind_immutable
BEFORE UPDATE OF kind ON issues
FOR EACH ROW
WHEN NEW.kind != OLD.kind
BEGIN
  SELECT RAISE(ABORT, 'issues: kind is immutable');
END;
