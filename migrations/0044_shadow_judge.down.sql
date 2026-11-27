-- Down refuses while judge sessions or usage remain; retain the journal and use the new binary.
CREATE TABLE review_claims_old (
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
INSERT INTO review_claims_old SELECT * FROM review_claims;
DROP TABLE review_claims;
ALTER TABLE review_claims_old RENAME TO review_claims;
CREATE UNIQUE INDEX one_open_review_per_task ON review_claims(task_id) WHERE ended_at IS NULL;
CREATE UNIQUE INDEX one_open_review_per_worktree ON review_claims(worktree_id) WHERE ended_at IS NULL;
CREATE TABLE usage_sessions_old (
 id TEXT PRIMARY KEY,
 task_id TEXT NOT NULL REFERENCES issues(id),
 hold_id TEXT NOT NULL,
 actor TEXT NOT NULL,
 model TEXT NOT NULL,
 role TEXT NOT NULL CHECK(role IN ('build','fix','review')),
 stream TEXT NOT NULL CHECK((role='review' AND stream IN ('primary','separate')) OR (role!='review' AND stream='')),
 format TEXT NOT NULL CHECK(format IN ('codex','claude')),
 input_digest TEXT NOT NULL,
 parser_revision TEXT NOT NULL,
 tool_version TEXT NOT NULL,
 version_source TEXT NOT NULL,
 uncached_input INTEGER NOT NULL CHECK(uncached_input>=0),
 cache_read INTEGER NOT NULL CHECK(cache_read>=0),
 cache_write INTEGER NOT NULL CHECK(cache_write>=0),
 output INTEGER NOT NULL CHECK(output>=0),
 reasoning INTEGER NOT NULL CHECK(reasoning>=0 AND reasoning<=output),
 total INTEGER NOT NULL CHECK(total>=0),
 record_json TEXT NOT NULL,
 recorded_at TEXT NOT NULL,
 UNIQUE(hold_id,input_digest)
);
INSERT INTO usage_sessions_old SELECT * FROM usage_sessions;
DROP TABLE usage_sessions;
ALTER TABLE usage_sessions_old RENAME TO usage_sessions;
CREATE INDEX usage_sessions_task ON usage_sessions(task_id);
DROP TABLE merge_calls;
DROP TABLE merge_holds;
