CREATE TABLE review_claims (
    id TEXT PRIMARY KEY,
    task_id TEXT NOT NULL REFERENCES issues(id),
    handoff_id TEXT NOT NULL REFERENCES handoffs(id),
    stream TEXT NOT NULL CHECK (stream IN ('primary', 'separate')),
    worktree_id TEXT NOT NULL,
    worktree_root TEXT NOT NULL,
    actor TEXT NOT NULL,
    model TEXT NOT NULL,
    since TEXT NOT NULL,
    last_heartbeat_at TEXT NOT NULL,
    expires_at TEXT NOT NULL,
    ended_at TEXT,
    binding_json TEXT NOT NULL,
    CHECK (expires_at > since)
);
CREATE UNIQUE INDEX one_open_review_per_task ON review_claims(task_id) WHERE ended_at IS NULL;
CREATE UNIQUE INDEX one_open_review_per_worktree ON review_claims(worktree_id) WHERE ended_at IS NULL;
