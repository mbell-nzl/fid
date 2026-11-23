CREATE TABLE default_conflicts (
    id TEXT PRIMARY KEY,
    task_id TEXT NOT NULL REFERENCES issues(id),
    hold_id TEXT NOT NULL REFERENCES holds(id),
    record_json TEXT NOT NULL
);
