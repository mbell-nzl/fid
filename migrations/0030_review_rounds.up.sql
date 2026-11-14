CREATE TABLE review_rounds (
 id TEXT PRIMARY KEY,
 task_id TEXT NOT NULL REFERENCES issues(id),
 handoff_id TEXT NOT NULL REFERENCES handoffs(id),
 record_json TEXT NOT NULL
);
CREATE INDEX review_rounds_by_task ON review_rounds(task_id);
CREATE TABLE round_findings (
 id TEXT PRIMARY KEY,
 round_id TEXT NOT NULL REFERENCES review_rounds(id),
 record_json TEXT NOT NULL
);
