DROP INDEX claims_one_active_mutation;
ALTER TABLE claims RENAME TO work_claims;
CREATE UNIQUE INDEX work_claims_one_active_mutation
  ON work_claims(issue_id) WHERE status = 'active';

CREATE TABLE roles (
  name        TEXT PRIMARY KEY,
  description TEXT,
  is_human    BOOLEAN DEFAULT 0,
  active      BOOLEAN DEFAULT 1,
  created_at  DATETIME
);

CREATE TABLE claims (
  ref_type    TEXT NOT NULL,
  ref_id      TEXT NOT NULL,
  claimed_by  TEXT NOT NULL,
  claimed_at  DATETIME NOT NULL,
  ttl_seconds INTEGER NOT NULL DEFAULT 1800,
  note        TEXT,
  PRIMARY KEY (ref_type, ref_id)
);

CREATE TABLE gate_runs (
  id               TEXT PRIMARY KEY,
  issue_id         TEXT NOT NULL REFERENCES issues(id),
  attempt_number   INTEGER NOT NULL,
  agent_run_id     TEXT REFERENCES agent_runs(id),
  recorded_by_role TEXT NOT NULL,
  command          TEXT NOT NULL,
  exit_code        INTEGER NOT NULL,
  output_excerpt   TEXT,
  ran_at           DATETIME NOT NULL,
  CHECK (attempt_number > 0)
);
