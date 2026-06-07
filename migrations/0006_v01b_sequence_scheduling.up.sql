CREATE TRIGGER issues_sequence_scheduling_guard_insert
BEFORE INSERT ON issues
FOR EACH ROW
WHEN NEW.kind = 'task'
  AND (
    (NEW.parent_id IS NULL AND NEW.sequence IS NOT NULL)
    OR (
      EXISTS (
        SELECT 1
        FROM issues
        WHERE id = NEW.parent_id AND kind = 'epic'
      )
      AND (
        NOT EXISTS (
          SELECT 1
          FROM epic_metadata
          WHERE issue_id = NEW.parent_id
        )
        OR EXISTS (
          SELECT 1
          FROM epic_metadata
          WHERE issue_id = NEW.parent_id
            AND scheduling = 'sequential'
            AND NEW.sequence IS NULL
        )
        OR EXISTS (
          SELECT 1
          FROM epic_metadata
          WHERE issue_id = NEW.parent_id
            AND scheduling = 'prioritised'
            AND NEW.sequence IS NOT NULL
        )
      )
    )
  )
BEGIN
  SELECT RAISE(
    ABORT,
    CASE
      WHEN NEW.parent_id IS NULL THEN 'issues: unattached task must not have sequence'
      WHEN NOT EXISTS (
        SELECT 1
        FROM epic_metadata
        WHERE issue_id = NEW.parent_id
      ) THEN 'issues: task parent epic must have epic_metadata'
      WHEN EXISTS (
        SELECT 1
        FROM epic_metadata
        WHERE issue_id = NEW.parent_id AND scheduling = 'sequential'
      ) THEN 'issues: task under sequential epic requires sequence'
      ELSE 'issues: task under prioritised epic must not have sequence'
    END
  );
END;

CREATE TRIGGER issues_sequence_scheduling_guard_update
BEFORE UPDATE OF parent_id, sequence ON issues
FOR EACH ROW
WHEN NEW.kind = 'task'
  AND (
    (NEW.parent_id IS NULL AND NEW.sequence IS NOT NULL)
    OR (
      EXISTS (
        SELECT 1
        FROM issues
        WHERE id = NEW.parent_id AND kind = 'epic'
      )
      AND (
        NOT EXISTS (
          SELECT 1
          FROM epic_metadata
          WHERE issue_id = NEW.parent_id
        )
        OR EXISTS (
          SELECT 1
          FROM epic_metadata
          WHERE issue_id = NEW.parent_id
            AND scheduling = 'sequential'
            AND NEW.sequence IS NULL
        )
        OR EXISTS (
          SELECT 1
          FROM epic_metadata
          WHERE issue_id = NEW.parent_id
            AND scheduling = 'prioritised'
            AND NEW.sequence IS NOT NULL
        )
      )
    )
  )
BEGIN
  SELECT RAISE(
    ABORT,
    CASE
      WHEN NEW.parent_id IS NULL THEN 'issues: unattached task must not have sequence'
      WHEN NOT EXISTS (
        SELECT 1
        FROM epic_metadata
        WHERE issue_id = NEW.parent_id
      ) THEN 'issues: task parent epic must have epic_metadata'
      WHEN EXISTS (
        SELECT 1
        FROM epic_metadata
        WHERE issue_id = NEW.parent_id AND scheduling = 'sequential'
      ) THEN 'issues: task under sequential epic requires sequence'
      ELSE 'issues: task under prioritised epic must not have sequence'
    END
  );
END;
