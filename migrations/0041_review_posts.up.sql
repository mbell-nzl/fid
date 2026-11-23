CREATE TABLE review_posts (
 id TEXT PRIMARY KEY,
 task_id TEXT NOT NULL REFERENCES issues(id),
 handoff_id TEXT NOT NULL REFERENCES handoffs(id),
 pr TEXT NOT NULL,
 url TEXT NOT NULL UNIQUE,
 render_version INTEGER NOT NULL CHECK(render_version=1),
 digest TEXT NOT NULL,
 posted_at TEXT NOT NULL
);
