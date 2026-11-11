-- A hold is a task handed out by fid next --claim to one actor in one
-- worktree (RESET 1.4a, ADR-0014). A hold is open until it ends; whether an
-- open hold has lapsed is read from the clock against expires_at, so it
-- can't be part of an index. The two partial indexes back up the write
-- lock: a task, and a worktree, has at most one open hold.
CREATE TABLE holds (
  id TEXT PRIMARY KEY,
  task_id TEXT NOT NULL REFERENCES issues(id),
  worktree_id TEXT NOT NULL,
  worktree_root TEXT NOT NULL,
  actor TEXT NOT NULL,
  since TEXT NOT NULL,
  expires_at TEXT NOT NULL,
  replaces_hold_id TEXT UNIQUE REFERENCES holds(id),
  ended_at TEXT,
  end_reason TEXT CHECK (end_reason IN ('replaced', 'task-completed', 'task-cancelled')),
  CHECK ((ended_at IS NULL) = (end_reason IS NULL)),
  CHECK (expires_at > since)
);
CREATE UNIQUE INDEX holds_one_open_per_task ON holds(task_id) WHERE ended_at IS NULL;
CREATE UNIQUE INDEX holds_one_open_per_worktree ON holds(worktree_id) WHERE ended_at IS NULL;
