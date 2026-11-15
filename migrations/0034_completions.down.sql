-- A completed journal needs fid with task.finished support. Downgrade only
-- after restoring a pre-completion backup; a projection is always rebuildable.
DROP TABLE completions;
CREATE TEMP TABLE questions_saved AS SELECT * FROM questions;
DROP TABLE questions;
CREATE TABLE questions (
 id TEXT PRIMARY KEY,
 task_id TEXT NOT NULL REFERENCES issues(id),
 hold_id TEXT NOT NULL REFERENCES holds(id),
 text TEXT NOT NULL CHECK (length(trim(text)) > 0),
 asker TEXT NOT NULL,
 asked_at TEXT NOT NULL,
 resolved_at TEXT,
 resolution TEXT CHECK (resolution IN ('answered', 'cancelled')),
 decision_id TEXT UNIQUE REFERENCES person_decisions(id),
 CHECK ((resolved_at IS NULL) = (resolution IS NULL)),
 CHECK ((resolution IS 'answered') = (decision_id IS NOT NULL))
);
CREATE UNIQUE INDEX questions_one_open_per_task ON questions(task_id) WHERE resolved_at IS NULL;
INSERT INTO questions SELECT * FROM questions_saved;
DROP TABLE questions_saved;
