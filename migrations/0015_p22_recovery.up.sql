ALTER TABLE kit_activations ADD COLUMN version INTEGER NOT NULL DEFAULT 1 CHECK (version > 0);
CREATE TABLE run_runtime (
  run_id TEXT PRIMARY KEY REFERENCES runs(id),
  version INTEGER NOT NULL CHECK (version > 0),
  executor_id TEXT NOT NULL UNIQUE,
  obligation_id TEXT REFERENCES envelope_obligations(id)
);
CREATE TABLE run_observations (
  id TEXT PRIMARY KEY,
  run_id TEXT NOT NULL REFERENCES runs(id),
  kind TEXT NOT NULL CHECK (kind IN ('liveness','checkpoint')),
  actor_id TEXT NOT NULL,
  observed_at TEXT NOT NULL,
  observation_json TEXT NOT NULL CHECK (json_valid(observation_json)),
  created_event_id TEXT NOT NULL UNIQUE
);
CREATE INDEX run_observations_by_run ON run_observations(run_id, observed_at, id);
CREATE TABLE run_recoveries (
  id TEXT PRIMARY KEY,
  prior_run_id TEXT NOT NULL UNIQUE REFERENCES runs(id),
  successor_run_id TEXT NOT NULL UNIQUE REFERENCES runs(id),
  prior_claim_id TEXT NOT NULL REFERENCES claims(id),
  successor_claim_id TEXT NOT NULL UNIQUE REFERENCES claims(id),
  observation_id TEXT NOT NULL REFERENCES run_observations(id),
  reconciliation_json TEXT NOT NULL CHECK (json_valid(reconciliation_json)),
  adopted_paths_json TEXT NOT NULL CHECK (json_valid(adopted_paths_json)),
  authority_json TEXT NOT NULL CHECK (json_valid(authority_json)),
  reason TEXT NOT NULL,
  created_event_id TEXT NOT NULL UNIQUE,
  CHECK (prior_run_id <> successor_run_id),
  CHECK (prior_claim_id <> successor_claim_id)
);
CREATE TABLE budget_samples (
  id TEXT PRIMARY KEY,
  envelope_id TEXT NOT NULL,
  envelope_revision INTEGER NOT NULL,
  budget_name TEXT NOT NULL,
  run_id TEXT NOT NULL REFERENCES runs(id),
  actor_id TEXT NOT NULL,
  source TEXT NOT NULL CHECK (source IN ('fid','activated-harness')),
  source_sample_id TEXT NOT NULL,
  sample_kind TEXT NOT NULL CHECK (sample_kind IN ('delta','baseline')),
  previous_sample_id TEXT REFERENCES budget_samples(id),
  value INTEGER CHECK (value IS NULL OR value >= 0),
  measurement_ref_json TEXT CHECK (measurement_ref_json IS NULL OR json_valid(measurement_ref_json)),
  observed_at TEXT NOT NULL,
  created_event_id TEXT NOT NULL,
  UNIQUE (envelope_id, envelope_revision, budget_name, source, source_sample_id),
  FOREIGN KEY (envelope_id, envelope_revision, budget_name)
    REFERENCES envelope_budgets(envelope_id, revision, name)
);
CREATE INDEX budget_samples_by_budget
  ON budget_samples(envelope_id, envelope_revision, budget_name, observed_at, id);
CREATE TABLE stops (
  id TEXT PRIMARY KEY,
  subject_type TEXT NOT NULL CHECK (subject_type IN ('run','envelope','activation')),
  subject_id TEXT NOT NULL,
  code TEXT NOT NULL CHECK (code IN (
    'executor-liveness-unknown','recovery-state-conflict','budget-exhausted',
    'budget-data-unavailable','activation-incomplete','capability-unavailable',
    'operator-stop'
  )),
  blocked_actions_json TEXT NOT NULL CHECK (json_valid(blocked_actions_json)),
  details_json TEXT NOT NULL CHECK (json_valid(details_json)),
  resolution_condition_json TEXT NOT NULL CHECK (json_valid(resolution_condition_json)),
  resolver_capability TEXT NOT NULL,
  actor_id TEXT NOT NULL,
  raised_at TEXT NOT NULL,
  created_event_id TEXT NOT NULL
);
CREATE INDEX stops_by_subject ON stops(subject_type, subject_id, raised_at, id);
CREATE TABLE stop_resolutions (
  stop_id TEXT PRIMARY KEY REFERENCES stops(id),
  actor_id TEXT NOT NULL,
  fact_event_id TEXT NOT NULL,
  resolved_at TEXT NOT NULL,
  created_event_id TEXT NOT NULL
);
CREATE TABLE attention_items (
  id TEXT PRIMARY KEY,
  semantic_key TEXT NOT NULL,
  subject_type TEXT NOT NULL CHECK (subject_type IN ('run','envelope','activation')),
  subject_id TEXT NOT NULL,
  code TEXT NOT NULL CHECK (code IN (
    'executor-liveness-unknown','recovery-state-conflict','budget-exhausted',
    'budget-data-unavailable','activation-incomplete','capability-unavailable',
    'operator-stop'
  )),
  status TEXT NOT NULL CHECK (status IN ('open','resolved')),
  version INTEGER NOT NULL CHECK (version > 0),
  raised_at TEXT NOT NULL,
  resolution_event_id TEXT,
  CHECK ((status = 'open' AND resolution_event_id IS NULL)
    OR (status = 'resolved' AND resolution_event_id IS NOT NULL))
);
CREATE UNIQUE INDEX attention_one_open_semantic_key
  ON attention_items(semantic_key) WHERE status = 'open';
CREATE TABLE attention_stop_links (
  attention_id TEXT NOT NULL REFERENCES attention_items(id),
  stop_id TEXT NOT NULL UNIQUE REFERENCES stops(id),
  PRIMARY KEY (attention_id, stop_id)
);
CREATE TABLE activation_resources (
  activation_id TEXT NOT NULL REFERENCES kit_activations(id),
  component TEXT NOT NULL CHECK (component IN ('kit','overlay','provider','project','backend')),
  version INTEGER NOT NULL CHECK (version > 0),
  old_identity TEXT,
  new_identity TEXT NOT NULL,
  observed_identity TEXT,
  state TEXT NOT NULL CHECK (state IN ('pending','old','new','error')),
  evidence_ref TEXT,
  error_code TEXT,
  PRIMARY KEY (activation_id, component)
);
CREATE TABLE activation_recoveries (
  activation_id TEXT PRIMARY KEY REFERENCES kit_activations(id),
  direction TEXT NOT NULL CHECK (direction IN ('restore-prior','complete-forward')),
  actor_attribution TEXT NOT NULL,
  reason TEXT NOT NULL,
  selected_event_id TEXT NOT NULL UNIQUE,
  completed_event_id TEXT UNIQUE
);
