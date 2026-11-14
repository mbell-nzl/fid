CREATE TABLE check_requests (
    id TEXT PRIMARY KEY,
    task_id TEXT NOT NULL REFERENCES issues(id),
    hold_id TEXT NOT NULL REFERENCES holds(id),
    actor TEXT NOT NULL,
    names_json TEXT NOT NULL,
    response_json TEXT
);
CREATE TABLE check_starts (
    id TEXT PRIMARY KEY,
    request_id TEXT NOT NULL REFERENCES check_requests(id),
    task_id TEXT NOT NULL REFERENCES issues(id),
    name TEXT NOT NULL,
    start_json TEXT NOT NULL,
    UNIQUE(request_id, name)
);
CREATE INDEX check_starts_task ON check_starts(task_id);
CREATE TABLE check_processes (
    check_id TEXT PRIMARY KEY REFERENCES check_starts(id),
    process_json TEXT NOT NULL
);
CREATE TABLE check_results (
    check_id TEXT PRIMARY KEY REFERENCES check_starts(id),
    result_json TEXT NOT NULL
);
CREATE TABLE check_files (
    check_id TEXT NOT NULL REFERENCES check_starts(id),
    path TEXT NOT NULL,
    artifact_id TEXT REFERENCES referenced_artifacts(id),
    reason TEXT,
    PRIMARY KEY(check_id, path),
    CHECK ((artifact_id IS NOT NULL AND reason IS NULL) OR (artifact_id IS NULL AND reason IS NOT NULL))
);
