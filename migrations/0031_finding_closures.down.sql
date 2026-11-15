PRAGMA defer_foreign_keys=ON;
-- Rebuild the referencing questions too: dropping a populated parent leaves
-- SQLite's deferred FK counter nonzero even if that parent is recreated.
CREATE TABLE questions_saved AS SELECT * FROM questions;
DROP TABLE questions;
DROP TABLE finding_closures;
CREATE TABLE person_decisions_new (
 id TEXT PRIMARY KEY,
 kind TEXT NOT NULL CHECK (kind IN ('takeover', 'cancel', 'answer')),
 actor TEXT NOT NULL CHECK (actor GLOB 'human:?*'),
 subject_type TEXT NOT NULL CHECK (subject_type IN ('task', 'epic')),
 subject_id TEXT NOT NULL REFERENCES issues(id),
 reason TEXT NOT NULL CHECK (length(trim(reason)) > 0),
 binding_id TEXT NOT NULL,
 decided_at TEXT NOT NULL,
 CHECK (kind = 'answer' OR binding_id = subject_id)
);
INSERT INTO person_decisions_new SELECT * FROM person_decisions;
DROP TABLE person_decisions;
ALTER TABLE person_decisions_new RENAME TO person_decisions;
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

