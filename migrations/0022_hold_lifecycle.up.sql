-- Keep an interrupted started task reserved even after its hold ends. The
-- same partial indexes cover both live holds and explicit reservations.
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
  end_reason TEXT CHECK (end_reason IN ('replaced', 'task-completed', 'task-cancelled', 'interrupted', 'takeover')),
  last_heartbeat_at TEXT,
  kept INTEGER NOT NULL DEFAULT 0 CHECK (kept IN (0, 1)),
  park_reason TEXT,
  CHECK ((ended_at IS NULL) = (end_reason IS NULL)),
  CHECK (expires_at > since),
  CHECK (last_heartbeat_at IS NULL OR (last_heartbeat_at >= since AND last_heartbeat_at < expires_at)),
  CHECK (kept = 0 OR (ended_at IS NOT NULL AND end_reason = 'interrupted')),
  CHECK ((park_reason IS NOT NULL) = (end_reason IS 'interrupted'))
);
INSERT INTO holds_new (id, task_id, worktree_id, worktree_root, actor, since, expires_at, replaces_hold_id, ended_at, end_reason, last_heartbeat_at)
SELECT id, task_id, worktree_id, worktree_root, actor, since, expires_at, replaces_hold_id, ended_at, end_reason, since FROM holds;
DROP TABLE holds;
ALTER TABLE holds_new RENAME TO holds;
CREATE UNIQUE INDEX holds_one_open_per_task ON holds(task_id) WHERE ended_at IS NULL OR kept = 1;
CREATE UNIQUE INDEX holds_one_open_per_worktree ON holds(worktree_id) WHERE ended_at IS NULL OR kept = 1;

-- Older hold.taken Events predate last_heartbeat_at. Replaying them and
-- upgrading their existing projection must produce the same initial time.
CREATE TRIGGER holds_initial_heartbeat AFTER INSERT ON holds
WHEN NEW.last_heartbeat_at IS NULL
BEGIN UPDATE holds SET last_heartbeat_at = NEW.since WHERE id = NEW.id; END;

CREATE TABLE person_decisions (
  id TEXT PRIMARY KEY,
  kind TEXT NOT NULL CHECK (kind IN ('takeover', 'cancel')),
  actor TEXT NOT NULL CHECK (actor GLOB 'human:?*'),
  subject_type TEXT NOT NULL CHECK (subject_type IN ('task', 'epic')),
  subject_id TEXT NOT NULL REFERENCES issues(id),
  reason TEXT NOT NULL CHECK (length(trim(reason)) > 0),
  binding_id TEXT NOT NULL REFERENCES issues(id),
  decided_at TEXT NOT NULL,
  CHECK (binding_id = subject_id)
);
