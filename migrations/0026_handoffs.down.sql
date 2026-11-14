PRAGMA defer_foreign_keys = ON;
DROP TABLE handoffs;
CREATE TABLE holds_new (
  id TEXT PRIMARY KEY,
  task_id TEXT NOT NULL REFERENCES issues(id),
  worktree_id TEXT NOT NULL,
  worktree_root TEXT NOT NULL,
  actor TEXT NOT NULL,
  since TEXT NOT NULL,
  expires_at TEXT NOT NULL,
  replaces_hold_id TEXT UNIQUE REFERENCES holds_new(id),
  ended_at TEXT,
  end_reason TEXT CHECK (end_reason IN ('replaced', 'task-completed', 'task-cancelled', 'interrupted', 'takeover', 'waiting')),
  last_heartbeat_at TEXT,
  kept INTEGER NOT NULL DEFAULT 0 CHECK (kept IN (0, 1)),
  park_reason TEXT,
  CHECK ((ended_at IS NULL) = (end_reason IS NULL)),
  CHECK (expires_at > since),
  CHECK (last_heartbeat_at IS NULL OR (last_heartbeat_at >= since AND last_heartbeat_at < expires_at)),
  CHECK (kept = 0 OR (ended_at IS NOT NULL AND end_reason IN ('interrupted', 'waiting'))),
  CHECK ((park_reason IS NOT NULL) = (end_reason IS 'interrupted'))
);
INSERT INTO holds_new (id, task_id, worktree_id, worktree_root, actor, since, expires_at, replaces_hold_id, ended_at, end_reason, last_heartbeat_at, kept, park_reason)
SELECT id, task_id, worktree_id, worktree_root, actor, since, expires_at, replaces_hold_id, ended_at, CASE WHEN end_reason='review' THEN 'interrupted' ELSE end_reason END, last_heartbeat_at, kept, CASE WHEN end_reason='review' THEN 'In review' ELSE park_reason END FROM holds;
DROP TABLE holds;
ALTER TABLE holds_new RENAME TO holds;
CREATE UNIQUE INDEX holds_one_open_per_task ON holds(task_id) WHERE ended_at IS NULL OR kept = 1;
CREATE UNIQUE INDEX holds_one_open_per_worktree ON holds(worktree_id) WHERE ended_at IS NULL OR kept = 1;

-- Older hold.taken Events predate last_heartbeat_at. Replaying them and
-- upgrading their existing projection must produce the same initial time.
CREATE TRIGGER holds_initial_heartbeat AFTER INSERT ON holds
WHEN NEW.last_heartbeat_at IS NULL
BEGIN UPDATE holds SET last_heartbeat_at = NEW.since WHERE id = NEW.id; END;


-- Replacing a referenced table leaves SQLite's deferred DROP counters even
-- when the replacement satisfies every child reference. Verify the complete
-- final graph before clearing those counters; any violation aborts migration.
CREATE TEMP TABLE fid_migration_26_fk_check (violation INTEGER CHECK (violation = 0));
INSERT INTO fid_migration_26_fk_check SELECT 1 FROM pragma_foreign_key_check;
DROP TABLE fid_migration_26_fk_check;
PRAGMA defer_foreign_keys = OFF;
