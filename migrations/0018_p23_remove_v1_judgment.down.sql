DROP TRIGGER activity_log_ref_integrity_guard_insert;
DROP TRIGGER links_relation_compatibility_guard_insert;
DROP TRIGGER links_to_ref_integrity_guard_insert;
DROP TRIGGER links_from_ref_integrity_guard_insert;
DROP TRIGGER links_endpoint_scope_guard_insert;
DROP VIEW issues_enriched;
DROP INDEX concerns_by_subject_status;
ALTER TABLE concerns RENAME TO judgment_concerns;
CREATE INDEX judgment_concerns_by_subject
  ON judgment_concerns(subject_type, subject_id, status, created_at, id);

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

CREATE TABLE concerns (
  id                         TEXT PRIMARY KEY,
  issue_id                   TEXT NOT NULL REFERENCES issues(id),
  raised_in_completion_id    TEXT NOT NULL REFERENCES task_completions(id),
  raised_by_role             TEXT NOT NULL,
  summary                    TEXT NOT NULL,
  body_md_path               TEXT,
  status                     TEXT NOT NULL,
  addressed_in_completion_id TEXT REFERENCES task_completions(id),
  triaged_in_event_id        TEXT REFERENCES triage_events(id),
  created_at                 DATETIME,
  updated_at                 DATETIME
);

CREATE TRIGGER task_completions_issue_kind_guard_insert
BEFORE INSERT ON task_completions
FOR EACH ROW
WHEN NOT EXISTS (
  SELECT 1
  FROM issues
  WHERE id = NEW.issue_id AND kind IN ('task', 'subtask')
)
BEGIN
  SELECT RAISE(ABORT, 'task_completions: issue_id must reference a task or subtask');
END;

CREATE TRIGGER decisions_epic_kind_guard_insert
BEFORE INSERT ON decisions
FOR EACH ROW
WHEN NEW.epic_id IS NOT NULL
  AND NOT EXISTS (
    SELECT 1
    FROM issues
    WHERE id = NEW.epic_id AND kind = 'epic'
  )
BEGIN
  SELECT RAISE(ABORT, 'decisions: epic_id must reference an epic when set');
END;

CREATE INDEX idx_concerns_issue_status ON concerns(issue_id, status);

PRAGMA foreign_keys=OFF;
PRAGMA legacy_alter_table=ON;
PRAGMA defer_foreign_keys=ON;

CREATE TABLE agent_runs_new (
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

INSERT INTO agent_runs_new (
  id, issue_id, role, model, tokens_in, tokens_out, price_in_per_mtok,
  price_out_per_mtok, cost_usd, token_breakdown_json, tool_call_count,
  provider_run_id, provider_cost_usd, started_at, ended_at
)
SELECT
  id, issue_id, role, model, tokens_in, tokens_out, price_in_per_mtok,
  price_out_per_mtok, cost_usd, token_breakdown_json, tool_call_count,
  provider_run_id, provider_cost_usd, started_at, ended_at
FROM agent_runs;

DROP TABLE agent_runs;
ALTER TABLE agent_runs_new RENAME TO agent_runs;

PRAGMA foreign_keys=ON;
PRAGMA legacy_alter_table=OFF;
PRAGMA defer_foreign_keys=OFF;

CREATE VIEW issues_enriched AS
WITH open_concerns AS (
  SELECT i.id AS issue_id, (
    SELECT COUNT(*) FROM concerns c WHERE c.issue_id = i.id AND c.status = 'open'
  ) AS n
  FROM issues i
), latest_triage AS (
  SELECT i.id AS issue_id, (
    SELECT te.verdict FROM triage_events te WHERE te.issue_id = i.id
    ORDER BY te.created_at DESC LIMIT 1
  ) AS verdict
  FROM issues i
)
SELECT i.*, (COALESCE(oc.n, 0) > 0) AS concerns_raised,
  COALESCE(oc.n, 0) AS open_concerns_count,
  CASE lt.verdict WHEN 'RESUME' THEN 'resumed' WHEN 'BLOCK' THEN 'blocked' ELSE NULL END AS triage_status
FROM issues i
LEFT JOIN open_concerns oc ON oc.issue_id = i.id
LEFT JOIN latest_triage lt ON lt.issue_id = i.id;

CREATE TRIGGER links_endpoint_scope_guard_insert
BEFORE INSERT ON links
FOR EACH ROW
WHEN NEW.from_type NOT IN ('issue', 'requirement')
  OR NEW.to_type NOT IN ('issue', 'requirement')
BEGIN
  SELECT RAISE(ABORT, 'links: endpoint types must be issue or requirement in v0.1');
END;

CREATE TRIGGER links_from_ref_integrity_guard_insert
BEFORE INSERT ON links
FOR EACH ROW
WHEN (NEW.from_type = 'issue' AND NOT EXISTS (
    SELECT 1
    FROM issues
    WHERE id = NEW.from_id
  ))
  OR (NEW.from_type = 'requirement' AND NOT EXISTS (
    SELECT 1
    FROM requirements
    WHERE id = NEW.from_id
  ))
BEGIN
  SELECT RAISE(ABORT, 'links: from endpoint must resolve to an existing row');
END;

CREATE TRIGGER links_to_ref_integrity_guard_insert
BEFORE INSERT ON links
FOR EACH ROW
WHEN (NEW.to_type = 'issue' AND NOT EXISTS (
    SELECT 1
    FROM issues
    WHERE id = NEW.to_id
  ))
  OR (NEW.to_type = 'requirement' AND NOT EXISTS (
    SELECT 1
    FROM requirements
    WHERE id = NEW.to_id
  ))
BEGIN
  SELECT RAISE(ABORT, 'links: to endpoint must resolve to an existing row');
END;

CREATE TRIGGER links_relation_compatibility_guard_insert
BEFORE INSERT ON links
FOR EACH ROW
WHEN NEW.from_type IN ('issue', 'requirement')
  AND NEW.to_type IN ('issue', 'requirement')
  AND NOT (
  NEW.relation = 'relates_to'
  OR (
    NEW.relation = 'implements'
    AND NEW.from_type = 'issue'
    AND NEW.to_type = 'requirement'
  )
  OR (
    NEW.relation = 'implemented_by'
    AND NEW.from_type = 'requirement'
    AND NEW.to_type = 'issue'
  )
  OR (
    NEW.relation IN (
      'blocks',
      'blocked_by',
      'duplicates',
      'duplicated_by',
      'supersedes',
      'superseded_by',
      'caused_by',
      'causes',
      'discovered_from',
      'discovered'
    )
    AND NEW.from_type = 'issue'
    AND NEW.to_type = 'issue'
  )
)
BEGIN
  SELECT RAISE(ABORT, 'links: relation/endpoint combination is not allowed in v0.1');
END;

CREATE TRIGGER activity_log_ref_integrity_guard_insert
BEFORE INSERT ON activity_log
FOR EACH ROW
WHEN NEW.ref_type NOT IN ('issue', 'requirement', 'decision')
  OR (NEW.ref_type = 'issue' AND NOT EXISTS (
    SELECT 1 FROM issues WHERE id = NEW.ref_id
  ))
  OR (NEW.ref_type = 'requirement' AND NOT EXISTS (
    SELECT 1 FROM requirements WHERE id = NEW.ref_id
  ))
  OR (NEW.ref_type = 'decision' AND NOT EXISTS (
    SELECT 1 FROM decisions WHERE id = NEW.ref_id
  ))
BEGIN
  SELECT RAISE(ABORT, 'activity_log: ref_type/ref_id must resolve to an existing row');
END;
