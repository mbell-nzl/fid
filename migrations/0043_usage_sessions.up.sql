CREATE TABLE usage_sessions (
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
CREATE INDEX usage_sessions_task ON usage_sessions(task_id);
