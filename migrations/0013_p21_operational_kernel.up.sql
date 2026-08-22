CREATE TABLE system_capabilities (
  capability_id TEXT PRIMARY KEY,
  schema_min INTEGER NOT NULL CHECK (schema_min > 0),
  schema_max INTEGER NOT NULL CHECK (schema_max >= schema_min),
  available INTEGER NOT NULL CHECK (available IN (0, 1)),
  limits_json TEXT NOT NULL CHECK (json_valid(limits_json)),
  trust_source TEXT NOT NULL CHECK (trust_source IN (
    'active-kit-activation','fid-executed','fid-git-observation','fid-journal',
    'phase-22-unavailable','phase-23-unavailable','phase-25-unavailable'
  ))
);

CREATE TABLE idempotency_results (
  operation_scope TEXT NOT NULL,
  idempotency_key TEXT NOT NULL CHECK (length(idempotency_key) BETWEEN 1 AND 128),
  request_digest TEXT NOT NULL,
  result_entity_id TEXT NOT NULL,
  result_event_id TEXT NOT NULL,
  aggregate_sequence INTEGER NOT NULL CHECK (aggregate_sequence > 0),
  result_json TEXT NOT NULL CHECK (json_valid(result_json)),
  created_at DATETIME NOT NULL,
  PRIMARY KEY (operation_scope, idempotency_key)
);

CREATE TABLE kit_activations (
  id TEXT PRIMARY KEY,
  prior_activation_id TEXT REFERENCES kit_activations(id),
  kit_identity TEXT NOT NULL,
  backend_identity TEXT NOT NULL,
  phase INTEGER NOT NULL CHECK (phase > 0),
  status TEXT NOT NULL CHECK (status IN ('prepared','staged','committed','finalized','aborted')),
  point_of_no_return INTEGER NOT NULL DEFAULT 0 CHECK (point_of_no_return IN (0, 1)),
  direction TEXT NOT NULL CHECK (direction = 'activate'),
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL
);
CREATE UNIQUE INDEX kit_activations_one_nonterminal
  ON kit_activations((1)) WHERE status IN ('prepared','staged','committed');
CREATE UNIQUE INDEX kit_activations_one_active
  ON kit_activations((1)) WHERE status = 'finalized';

CREATE TABLE activation_components (
  activation_id TEXT NOT NULL REFERENCES kit_activations(id),
  component_type TEXT NOT NULL CHECK (component_type IN ('authority-profile','gate-definitions')),
  component_id TEXT NOT NULL,
  component_version INTEGER NOT NULL CHECK (component_version > 0),
  content_digest TEXT NOT NULL,
  expected_prior_digest TEXT,
  status TEXT NOT NULL CHECK (status IN ('staged','committed','finalized','aborted')),
  evidence_json TEXT NOT NULL CHECK (json_valid(evidence_json)),
  error_code TEXT,
  PRIMARY KEY (activation_id, component_type, component_id, component_version)
);

CREATE TABLE activation_actor_profiles (
  activation_id TEXT NOT NULL REFERENCES kit_activations(id),
  actor_id TEXT NOT NULL,
  actor_type TEXT NOT NULL CHECK (actor_type IN ('agent','human','service')),
  producer TEXT NOT NULL,
  provenance_ref TEXT NOT NULL,
  session_ref TEXT,
  capabilities_json TEXT NOT NULL CHECK (json_valid(capabilities_json)),
  PRIMARY KEY (activation_id, actor_id)
);

CREATE TABLE activation_gate_definitions (
  activation_id TEXT NOT NULL REFERENCES kit_activations(id),
  definition_id TEXT NOT NULL,
  definition_version INTEGER NOT NULL CHECK (definition_version > 0),
  definition_digest TEXT NOT NULL,
  definition_json TEXT NOT NULL CHECK (json_valid(definition_json)),
  PRIMARY KEY (activation_id, definition_id, definition_version),
  UNIQUE (activation_id, definition_digest)
);

CREATE TABLE gate_registrations (
  id TEXT NOT NULL,
  version INTEGER NOT NULL CHECK (version > 0),
  activation_id TEXT NOT NULL,
  definition_id TEXT NOT NULL,
  definition_version INTEGER NOT NULL,
  definition_digest TEXT NOT NULL,
  binding_digest TEXT NOT NULL,
  project_id TEXT NOT NULL,
  worktree_id TEXT NOT NULL,
  scope_json TEXT NOT NULL CHECK (json_valid(scope_json)),
  policy_json TEXT NOT NULL CHECK (json_valid(policy_json)),
  authorization_provenance TEXT NOT NULL,
  created_event_id TEXT NOT NULL,
  created_at DATETIME NOT NULL,
  PRIMARY KEY (id, version),
  UNIQUE (definition_digest, binding_digest),
  FOREIGN KEY (activation_id, definition_id, definition_version)
    REFERENCES activation_gate_definitions(activation_id, definition_id, definition_version)
);

CREATE TABLE work_envelopes (
  id TEXT PRIMARY KEY,
  current_revision INTEGER NOT NULL CHECK (current_revision = 1),
  status TEXT NOT NULL CHECK (status IN ('draft','active','completed')),
  created_at DATETIME NOT NULL
);

CREATE TABLE envelope_revisions (
  envelope_id TEXT NOT NULL REFERENCES work_envelopes(id),
  revision INTEGER NOT NULL CHECK (revision = 1),
  prior_revision INTEGER CHECK (prior_revision IS NULL),
  outcome TEXT NOT NULL,
  non_goals_json TEXT NOT NULL CHECK (json_valid(non_goals_json)),
  risk_result_json TEXT NOT NULL CHECK (json_valid(risk_result_json)),
  repository_scope_json TEXT NOT NULL CHECK (json_valid(repository_scope_json)),
  status TEXT NOT NULL CHECK (status IN ('draft','active','completed')),
  created_actor_id TEXT NOT NULL,
  approved_actor_id TEXT,
  created_at DATETIME NOT NULL,
  activated_at DATETIME,
  completed_at DATETIME,
  PRIMARY KEY (envelope_id, revision)
);

CREATE TABLE envelope_sources (
  envelope_id TEXT NOT NULL,
  revision INTEGER NOT NULL,
  ordinal INTEGER NOT NULL CHECK (ordinal >= 0),
  source_json TEXT NOT NULL CHECK (json_valid(source_json)),
  PRIMARY KEY (envelope_id, revision, ordinal),
  FOREIGN KEY (envelope_id, revision) REFERENCES envelope_revisions(envelope_id, revision)
);

CREATE TABLE envelope_issue_selection (
  envelope_id TEXT NOT NULL,
  revision INTEGER NOT NULL,
  issue_id TEXT NOT NULL REFERENCES issues(id),
  primary_issue INTEGER NOT NULL CHECK (primary_issue = 1),
  PRIMARY KEY (envelope_id, revision, issue_id),
  UNIQUE (envelope_id, revision, primary_issue),
  FOREIGN KEY (envelope_id, revision) REFERENCES envelope_revisions(envelope_id, revision)
);

CREATE TABLE envelope_authority_grants (
  envelope_id TEXT NOT NULL,
  revision INTEGER NOT NULL,
  ordinal INTEGER NOT NULL CHECK (ordinal >= 0),
  actor_selector_json TEXT NOT NULL CHECK (json_valid(actor_selector_json)),
  actions_json TEXT NOT NULL CHECK (json_valid(actions_json)),
  scope_json TEXT NOT NULL CHECK (json_valid(scope_json)),
  PRIMARY KEY (envelope_id, revision, ordinal),
  FOREIGN KEY (envelope_id, revision) REFERENCES envelope_revisions(envelope_id, revision)
);

CREATE TABLE envelope_obligations (
  id TEXT PRIMARY KEY,
  envelope_id TEXT NOT NULL,
  revision INTEGER NOT NULL,
  obligation_type TEXT NOT NULL CHECK (obligation_type = 'gate-evidence'),
  registration_id TEXT NOT NULL,
  registration_version INTEGER NOT NULL,
  obligation_json TEXT NOT NULL CHECK (json_valid(obligation_json)),
  UNIQUE (envelope_id, revision, id),
  FOREIGN KEY (envelope_id, revision) REFERENCES envelope_revisions(envelope_id, revision),
  FOREIGN KEY (registration_id, registration_version) REFERENCES gate_registrations(id, version)
);

CREATE TABLE envelope_budgets (
  envelope_id TEXT NOT NULL,
  revision INTEGER NOT NULL,
  name TEXT NOT NULL,
  budget_json TEXT NOT NULL CHECK (json_valid(budget_json)),
  PRIMARY KEY (envelope_id, revision, name),
  FOREIGN KEY (envelope_id, revision) REFERENCES envelope_revisions(envelope_id, revision)
);

CREATE TABLE work_claims (
  id TEXT PRIMARY KEY,
  issue_id TEXT NOT NULL REFERENCES issues(id),
  envelope_id TEXT NOT NULL,
  envelope_revision INTEGER NOT NULL,
  actor_id TEXT NOT NULL,
  worktree_id TEXT NOT NULL,
  scope_json TEXT NOT NULL CHECK (json_valid(scope_json)),
  mode TEXT NOT NULL CHECK (mode = 'mutation'),
  version INTEGER NOT NULL CHECK (version > 0),
  status TEXT NOT NULL CHECK (status IN ('active','released','expired')),
  acquired_at DATETIME NOT NULL,
  expires_at DATETIME NOT NULL,
  released_at DATETIME,
  release_reason TEXT CHECK (release_reason IS NULL OR release_reason IN (
    'completed','failed','operator-released'
  )),
  FOREIGN KEY (envelope_id, envelope_revision) REFERENCES envelope_revisions(envelope_id, revision)
);
CREATE UNIQUE INDEX work_claims_one_active_mutation
  ON work_claims(issue_id) WHERE status = 'active';

CREATE TABLE runs (
  id TEXT PRIMARY KEY,
  issue_id TEXT NOT NULL REFERENCES issues(id),
  envelope_id TEXT NOT NULL,
  envelope_revision INTEGER NOT NULL,
  claim_id TEXT NOT NULL REFERENCES work_claims(id),
  claim_version INTEGER NOT NULL,
  actor_id TEXT NOT NULL,
  start_repository_json TEXT NOT NULL CHECK (json_valid(start_repository_json)),
  terminal_repository_json TEXT CHECK (terminal_repository_json IS NULL OR json_valid(terminal_repository_json)),
  status TEXT NOT NULL CHECK (status IN ('active','succeeded','failed','abandoned')),
  started_at DATETIME NOT NULL,
  ended_at DATETIME,
  FOREIGN KEY (envelope_id, envelope_revision) REFERENCES envelope_revisions(envelope_id, revision)
);
CREATE UNIQUE INDEX runs_one_active_per_claim
  ON runs(claim_id) WHERE status = 'active';

CREATE TABLE run_budget_bindings (
  run_id TEXT NOT NULL REFERENCES runs(id),
  name TEXT NOT NULL,
  budget_json TEXT NOT NULL CHECK (json_valid(budget_json)),
  PRIMARY KEY (run_id, name)
);

CREATE TABLE execution_attempts (
  id TEXT PRIMARY KEY,
  registration_id TEXT NOT NULL,
  registration_version INTEGER NOT NULL,
  obligation_id TEXT NOT NULL REFERENCES envelope_obligations(id),
  run_id TEXT NOT NULL REFERENCES runs(id),
  actor_id TEXT NOT NULL,
  start_repository_json TEXT NOT NULL CHECK (json_valid(start_repository_json)),
  resolved_policy_digest TEXT NOT NULL,
  producer TEXT NOT NULL CHECK (producer = 'fid'),
  started_at DATETIME NOT NULL,
  FOREIGN KEY (registration_id, registration_version) REFERENCES gate_registrations(id, version)
);

CREATE TABLE referenced_artifacts (
  id TEXT PRIMARY KEY,
  digest TEXT NOT NULL UNIQUE,
  media_type TEXT NOT NULL,
  size_bytes INTEGER NOT NULL CHECK (size_bytes >= 0),
  retention_class TEXT NOT NULL CHECK (retention_class = 'evidence'),
  availability TEXT NOT NULL CHECK (availability IN ('available','missing')),
  location_policy TEXT NOT NULL CHECK (location_policy = 'content-addressed'),
  created_at DATETIME NOT NULL,
  expires_at DATETIME
);

CREATE TABLE execution_receipts (
  id TEXT PRIMARY KEY,
  attempt_id TEXT NOT NULL UNIQUE REFERENCES execution_attempts(id),
  trust_source TEXT NOT NULL CHECK (trust_source = 'fid-executed'),
  terminal_status TEXT NOT NULL CHECK (terminal_status IN ('succeeded','failed','timed-out','signaled','spawn-failed')),
  exit_code INTEGER,
  signal TEXT,
  timed_out INTEGER NOT NULL CHECK (timed_out IN (0, 1)),
  started_at DATETIME NOT NULL,
  ended_at DATETIME NOT NULL,
  terminal_repository_json TEXT NOT NULL CHECK (json_valid(terminal_repository_json)),
  toolchain_json TEXT NOT NULL CHECK (json_valid(toolchain_json)),
  output_excerpt TEXT NOT NULL,
  output_digest TEXT NOT NULL,
  output_bytes INTEGER NOT NULL CHECK (output_bytes >= 0),
  output_truncated INTEGER NOT NULL CHECK (output_truncated IN (0, 1)),
  redaction_applied INTEGER NOT NULL CHECK (redaction_applied IN (0, 1)),
  CHECK ((terminal_status = 'spawn-failed' AND exit_code IS NULL AND signal IS NULL AND timed_out = 0)
      OR (terminal_status = 'timed-out' AND timed_out = 1)
      OR (terminal_status = 'signaled' AND exit_code IS NULL AND signal IS NOT NULL AND timed_out = 0)
      OR (terminal_status IN ('succeeded','failed') AND exit_code IS NOT NULL AND signal IS NULL AND timed_out = 0))
);

CREATE TABLE receipt_artifacts (
  receipt_id TEXT NOT NULL REFERENCES execution_receipts(id),
  artifact_id TEXT NOT NULL REFERENCES referenced_artifacts(id),
  declared_path TEXT NOT NULL,
  PRIMARY KEY (receipt_id, artifact_id)
);

CREATE TABLE evidence_records (
  id TEXT PRIMARY KEY,
  obligation_id TEXT NOT NULL REFERENCES envelope_obligations(id),
  receipt_id TEXT NOT NULL REFERENCES execution_receipts(id),
  evidence_type TEXT NOT NULL CHECK (evidence_type = 'receipt'),
  producer TEXT NOT NULL CHECK (producer = 'fid'),
  trust_source TEXT NOT NULL CHECK (trust_source = 'fid-executed'),
  target_repository_json TEXT NOT NULL CHECK (json_valid(target_repository_json)),
  validity TEXT NOT NULL CHECK (validity IN ('valid','invalid')),
  created_at DATETIME NOT NULL,
  UNIQUE (obligation_id, receipt_id)
);

CREATE TABLE evidence_artifact_links (
  evidence_id TEXT NOT NULL REFERENCES evidence_records(id),
  artifact_id TEXT NOT NULL REFERENCES referenced_artifacts(id),
  required_material INTEGER NOT NULL CHECK (required_material IN (0, 1)),
  PRIMARY KEY (evidence_id, artifact_id)
);
