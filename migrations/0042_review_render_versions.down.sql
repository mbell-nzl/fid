-- Refuse the downgrade if version-2 receipts exist; never discard receipts.
CREATE TABLE review_posts_new (
 id TEXT PRIMARY KEY,
 task_id TEXT NOT NULL REFERENCES issues(id),
 handoff_id TEXT NOT NULL REFERENCES handoffs(id),
 pr TEXT NOT NULL,
 url TEXT NOT NULL UNIQUE,
 render_version INTEGER NOT NULL CHECK(render_version=1),
 digest TEXT NOT NULL,
 posted_at TEXT NOT NULL
);
INSERT INTO review_posts_new SELECT * FROM review_posts;
DROP TABLE review_posts;
ALTER TABLE review_posts_new RENAME TO review_posts;
