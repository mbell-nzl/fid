-- The journal remains canonical; rebuild with the previous fid to roll back.
DROP TABLE person_decisions;
CREATE TABLE holds_old (
  id TEXT PRIMARY KEY,
  task_id TEXT NOT NULL REFERENCES issues(id),
  worktree_id TEXT NOT NULL,
  worktree_root TEXT NOT NULL,
  actor TEXT NOT NULL,
  since TEXT NOT NULL,
  expires_at TEXT NOT NULL,
  replaces_hold_id TEXT UNIQUE REFERENCES holds_old(id),
  ended_at TEXT,
  end_reason TEXT CHECK (end_reason IN ('replaced', 'task-completed', 'task-cancelled')),
  CHECK ((ended_at IS NULL) = (end_reason IS NULL)), CHECK (expires_at > since)
);
INSERT INTO holds_old SELECT id, task_id, worktree_id, worktree_root, actor, since, expires_at, replaces_hold_id, ended_at,
 CASE WHEN end_reason IN ('interrupted', 'takeover') THEN 'replaced' ELSE end_reason END FROM holds;
DROP TABLE holds;
ALTER TABLE holds_old RENAME TO holds;
CREATE UNIQUE INDEX holds_one_open_per_task ON holds(task_id) WHERE ended_at IS NULL;
CREATE UNIQUE INDEX holds_one_open_per_worktree ON holds(worktree_id) WHERE ended_at IS NULL;
