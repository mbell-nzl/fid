CREATE TABLE projection_meta (
  key                       TEXT PRIMARY KEY CHECK (key = 'default'),
  journal_file              TEXT,
  journal_offset            INTEGER NOT NULL DEFAULT 0,
  last_event_id             TEXT,
  last_event_hash           TEXT,
  event_schema              INTEGER NOT NULL,
  projection_schema_version INTEGER NOT NULL,
  replay_status             TEXT NOT NULL,
  updated_at                DATETIME NOT NULL,
  CHECK (journal_offset >= 0),
  CHECK (replay_status IN ('clean', 'catching_up', 'rebuilding', 'error'))
);

CREATE TABLE roles (
  name        TEXT PRIMARY KEY,
  description TEXT,
  is_human    BOOLEAN DEFAULT 0,
  active      BOOLEAN DEFAULT 1,
  created_at  DATETIME
);

CREATE TABLE issues (
  id           TEXT PRIMARY KEY,
  kind         TEXT NOT NULL,
  parent_id    TEXT REFERENCES issues(id),
  title        TEXT NOT NULL,
  workflow_id  TEXT NOT NULL,
  status       TEXT NOT NULL,
  priority     INTEGER,
  sequence     INTEGER,
  spec_md_path TEXT,
  est_tokens   INTEGER,
  due_date     DATE,
  in_inbox     BOOLEAN DEFAULT 0,
  owner_role   TEXT,
  owner_actor  TEXT,
  created_at   DATETIME,
  updated_at   DATETIME,
  CHECK (kind IN ('epic', 'task', 'subtask')),
  CHECK ((kind = 'epic' AND workflow_id = 'epic-default')
      OR (kind = 'task' AND workflow_id = 'task-default')
      OR (kind = 'subtask' AND workflow_id = 'subtask-default')),
  CHECK ((workflow_id = 'task-default' AND status IN (
          'todo', 'in_progress', 'dev_complete', 'qa_review',
          'qa_block', 'failed_gate', 'qa_pass', 'done', 'cancelled'))
      OR (workflow_id = 'epic-default' AND status IN (
          'todo', 'in_progress', 'wrapped', 'closed', 'cancelled'))
      OR (workflow_id = 'subtask-default' AND status IN (
          'todo', 'in_progress', 'done', 'cancelled'))),
  CHECK (priority IS NULL OR (priority >= 0 AND priority <= 5)),
  CHECK (kind != 'epic' OR est_tokens IS NULL),
  CHECK (sequence IS NULL OR sequence > 0),
  UNIQUE (parent_id, sequence)
);
