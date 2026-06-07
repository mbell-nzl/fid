CREATE TABLE epic_metadata (
  issue_id      TEXT PRIMARY KEY REFERENCES issues(id),
  scheduling    TEXT NOT NULL,
  brief_md_path TEXT,
  target_start  DATE,
  target_end    DATE
);

CREATE TABLE requirements (
  id               TEXT PRIMARY KEY,
  type             TEXT NOT NULL,
  summary          TEXT NOT NULL,
  source_doc       TEXT NOT NULL,
  source_section   TEXT,
  last_synced_hash TEXT,
  status           TEXT
);

CREATE TABLE decisions (
  id           TEXT PRIMARY KEY,
  epic_id      TEXT REFERENCES issues(id),
  author_role  TEXT NOT NULL,
  summary      TEXT NOT NULL,
  body_md_path TEXT NOT NULL,
  area         TEXT,
  status       TEXT NOT NULL,
  created_at   DATETIME
);

CREATE TABLE escalations (
  id            TEXT PRIMARY KEY,
  status        TEXT NOT NULL,
  trigger_class TEXT,
  raised_by_role TEXT,
  raised_at     DATETIME,
  resolved_at   DATETIME,
  resolution_dec TEXT REFERENCES decisions(id)
);

CREATE TABLE triage_events (
  id            TEXT PRIMARY KEY,
  epic_id       TEXT NOT NULL REFERENCES issues(id),
  issue_id      TEXT NOT NULL REFERENCES issues(id),
  triager_role  TEXT NOT NULL,
  source        TEXT NOT NULL,
  source_ref_id TEXT,
  verdict       TEXT NOT NULL,
  esc_id        TEXT REFERENCES escalations(id),
  reasoning_md  TEXT,
  created_at    DATETIME
);

CREATE TABLE task_completions (
  id                      TEXT PRIMARY KEY,
  issue_id                TEXT NOT NULL REFERENCES issues(id),
  attempt_number          INTEGER NOT NULL,
  developer_role          TEXT,
  model                   TEXT,
  status                  TEXT,
  notes_md_path           TEXT,
  actual_tokens           INTEGER,
  commit_hashes           TEXT,
  files_touched           TEXT,
  expected_files_touched  TEXT,
  gates_passed            INTEGER DEFAULT 0,
  gates_failed            INTEGER DEFAULT 0,
  gate_results            TEXT,
  human_intervention      BOOLEAN DEFAULT 0,
  human_intervention_note TEXT,
  completed_at            DATETIME
);

CREATE TABLE acceptance_criteria (
  id       TEXT PRIMARY KEY,
  issue_id TEXT NOT NULL REFERENCES issues(id),
  ordinal  INTEGER NOT NULL,
  text     TEXT NOT NULL
);

CREATE TABLE ac_verifications (
  id             TEXT PRIMARY KEY,
  ac_id          TEXT NOT NULL REFERENCES acceptance_criteria(id),
  verifier_role  TEXT NOT NULL,
  verifier_actor TEXT,
  status         TEXT NOT NULL,
  evidence_type  TEXT,
  evidence_ref   TEXT,
  notes          TEXT,
  completion_id  TEXT REFERENCES task_completions(id),
  verified_at    DATETIME NOT NULL
);

CREATE TABLE agent_runs (
  id                   TEXT PRIMARY KEY,
  issue_id             TEXT REFERENCES issues(id),
  completion_id        TEXT REFERENCES task_completions(id),
  role                 TEXT NOT NULL,
  model                TEXT NOT NULL,
  tokens_in            INTEGER,
  tokens_out           INTEGER,
  price_in_per_mtok    REAL,
  price_out_per_mtok   REAL,
  cost_usd             REAL,
  token_breakdown_json TEXT,
  tool_call_count      INTEGER,
  provider_run_id      TEXT,
  provider_cost_usd    REAL,
  started_at           DATETIME,
  ended_at             DATETIME
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

CREATE TABLE concerns (
  id                         TEXT PRIMARY KEY,
  issue_id                   TEXT NOT NULL REFERENCES issues(id),
  raised_in_completion_id    TEXT NOT NULL REFERENCES task_completions(id),
  raised_by_role             TEXT NOT NULL,
  summary                    TEXT NOT NULL,
  body_md_path               TEXT,
  status                     TEXT NOT NULL,
  addressed_in_completion_id TEXT REFERENCES task_completions(id),
  escalated_to_esc_id        TEXT REFERENCES escalations(id),
  triaged_in_event_id        TEXT REFERENCES triage_events(id),
  created_at                 DATETIME,
  updated_at                 DATETIME
);

CREATE TABLE labels (
  name        TEXT PRIMARY KEY,
  colour      TEXT,
  description TEXT
);

CREATE TABLE issue_labels (
  issue_id TEXT NOT NULL REFERENCES issues(id),
  label    TEXT NOT NULL REFERENCES labels(name),
  PRIMARY KEY (issue_id, label)
);

CREATE TABLE links (
  from_type  TEXT NOT NULL,
  from_id    TEXT NOT NULL,
  relation   TEXT NOT NULL,
  to_type    TEXT NOT NULL,
  to_id      TEXT NOT NULL,
  created_at DATETIME,
  PRIMARY KEY (from_type, from_id, relation, to_type, to_id)
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

CREATE TABLE activity_log (
  id               TEXT PRIMARY KEY,
  journal_event_id TEXT NOT NULL,
  ref_type         TEXT NOT NULL,
  ref_id           TEXT NOT NULL,
  event            TEXT NOT NULL,
  details          TEXT,
  actor            TEXT,
  occurred_at      DATETIME
);
