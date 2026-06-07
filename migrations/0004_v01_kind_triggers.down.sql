DROP TRIGGER IF EXISTS issues_kind_immutable;

DROP TRIGGER IF EXISTS decisions_epic_kind_guard_insert;

DROP TRIGGER IF EXISTS triage_events_epic_kind_guard_insert;

DROP TRIGGER IF EXISTS epic_metadata_issue_kind_guard_insert;

DROP TRIGGER IF EXISTS task_completions_issue_kind_guard_insert;

DROP TRIGGER IF EXISTS acceptance_criteria_issue_kind_guard_insert;

DROP TRIGGER IF EXISTS issues_parent_kind_guard_update;

DROP TRIGGER IF EXISTS issues_parent_kind_guard_insert;
