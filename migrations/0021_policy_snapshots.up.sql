-- The journal is canonical. This append-only projection makes the latest
-- effective digest cheap to compare while the ledger's write lock is held.
CREATE TABLE policy_snapshots (
  sequence INTEGER PRIMARY KEY CHECK (sequence > 0),
  source_ref TEXT NOT NULL,
  source_commit TEXT,
  digest TEXT NOT NULL,
  effective_json TEXT NOT NULL CHECK (json_valid(effective_json)),
  CHECK ((source_ref = 'built-in defaults' AND source_commit IS NULL)
    OR (source_ref LIKE 'origin/%' AND source_commit IS NOT NULL))
);
