-- Restore schema 26 definitions only; removed data requires a pre-upgrade backup.
DROP VIEW issues_enriched;
DROP TRIGGER activity_log_ref_integrity_guard_insert;

CREATE TABLE workflow_assignments (
  scope       TEXT NOT NULL,
  kind        TEXT,
  workflow_id TEXT NOT NULL REFERENCES workflows(id),
  assigned_at DATETIME,
  CHECK (scope IN ('kind', 'project')),
  CHECK ((scope = 'kind' AND kind IS NOT NULL)
      OR (scope = 'project' AND kind IS NULL))
);

CREATE TABLE system_capabilities (
  capability_id TEXT PRIMARY KEY,
  schema_min INTEGER NOT NULL CHECK (schema_min > 0),
  schema_max INTEGER NOT NULL CHECK (schema_max >= schema_min),
  available INTEGER NOT NULL CHECK (available IN (0, 1)),
  limits_json TEXT NOT NULL CHECK (json_valid(limits_json)),
  trust_source TEXT NOT NULL CHECK (trust_source IN (
    'active-kit-activation','fid-executed','fid-git-observation','fid-journal',
    'phase-22-unavailable','phase-23-unavailable','phase-24-unavailable','phase-25-unavailable'
  ))
);

CREATE TABLE kit_activations (
  id TEXT PRIMARY KEY,
  prior_activation_id TEXT REFERENCES kit_activations(id),
  kit_identity TEXT NOT NULL,
  backend_identity TEXT NOT NULL,
  phase INTEGER NOT NULL CHECK (phase > 0),
  status TEXT NOT NULL CHECK (status IN ('prepared','staged','committed','finalized','aborted','superseded')),
  point_of_no_return INTEGER NOT NULL DEFAULT 0 CHECK (point_of_no_return IN (0, 1)),
  direction TEXT NOT NULL CHECK (direction = 'activate'),
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL
, version INTEGER NOT NULL DEFAULT 1 CHECK (version > 0));

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
  obligation_type TEXT NOT NULL CHECK (obligation_type IN ('gate-evidence','review','human-outcome')),
  registration_id TEXT,
  registration_version INTEGER,
  obligation_json TEXT NOT NULL CHECK (json_valid(obligation_json)), review_profile TEXT, outcome_class TEXT,
  CHECK ((obligation_type = 'gate-evidence'
          AND registration_id IS NOT NULL AND registration_version IS NOT NULL)
      OR (obligation_type IN ('review','human-outcome')
          AND registration_id IS NULL AND registration_version IS NULL)),
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

CREATE TABLE "claims" (
  id TEXT PRIMARY KEY,
  issue_id TEXT NOT NULL REFERENCES issues(id),
  envelope_id TEXT NOT NULL,
  envelope_revision INTEGER NOT NULL,
  actor_id TEXT NOT NULL,
  worktree_id TEXT NOT NULL,
  scope_json TEXT NOT NULL CHECK (json_valid(scope_json)),
  mode TEXT NOT NULL CHECK (mode IN ('mutation','review')),
  version INTEGER NOT NULL CHECK (version > 0),
  status TEXT NOT NULL CHECK (status IN ('active','released','expired','superseded')),
  acquired_at DATETIME NOT NULL,
  expires_at DATETIME NOT NULL,
  released_at DATETIME,
  release_reason TEXT CHECK (release_reason IS NULL OR release_reason IN (
    'completed','failed','operator-released','interrupted','recovery-superseded','abandoned'
  )),
  FOREIGN KEY (envelope_id, envelope_revision) REFERENCES envelope_revisions(envelope_id, revision)
);

CREATE TABLE runs (
  id TEXT PRIMARY KEY,
  issue_id TEXT NOT NULL REFERENCES issues(id),
  envelope_id TEXT NOT NULL,
  envelope_revision INTEGER NOT NULL,
  claim_id TEXT NOT NULL REFERENCES "claims"(id),
  claim_version INTEGER NOT NULL,
  actor_id TEXT NOT NULL,
  start_repository_json TEXT NOT NULL CHECK (json_valid(start_repository_json)),
  terminal_repository_json TEXT CHECK (terminal_repository_json IS NULL OR json_valid(terminal_repository_json)),
  status TEXT NOT NULL CHECK (status IN ('active','succeeded','failed','abandoned','interrupted')),
  started_at DATETIME NOT NULL,
  ended_at DATETIME,
  FOREIGN KEY (envelope_id, envelope_revision) REFERENCES envelope_revisions(envelope_id, revision)
);

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

CREATE TABLE execution_receipts (
  id TEXT PRIMARY KEY,
  attempt_id TEXT NOT NULL UNIQUE REFERENCES execution_attempts(id),
  trust_source TEXT NOT NULL CHECK (trust_source = 'fid-executed'),
  terminal_status TEXT NOT NULL CHECK (terminal_status IN ('succeeded','failed','timed-out','signaled','spawn-failed','cancelled')),
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
  rejected_artifacts_json TEXT NOT NULL
    CHECK (json_valid(rejected_artifacts_json) AND json_type(rejected_artifacts_json) = 'array'),
  CHECK ((terminal_status = 'spawn-failed' AND exit_code IS NULL AND signal IS NULL AND timed_out = 0)
      OR (terminal_status = 'timed-out' AND timed_out = 1)
      OR (terminal_status = 'signaled' AND exit_code IS NULL AND signal IS NOT NULL AND timed_out = 0)
      OR (terminal_status IN ('succeeded','failed') AND exit_code IS NOT NULL AND signal IS NULL AND timed_out = 0)
      OR (terminal_status = 'cancelled' AND timed_out = 0 AND NOT (exit_code IS NOT NULL AND signal IS NOT NULL)))
);

CREATE TABLE receipt_artifacts (
  receipt_id TEXT NOT NULL REFERENCES execution_receipts(id),
  artifact_id TEXT NOT NULL REFERENCES referenced_artifacts(id),
  declared_path TEXT NOT NULL,
  PRIMARY KEY (receipt_id, artifact_id)
);

CREATE TABLE evidence_artifact_links (
  evidence_id TEXT NOT NULL REFERENCES evidence_records(id),
  artifact_id TEXT NOT NULL REFERENCES referenced_artifacts(id),
  required_material INTEGER NOT NULL CHECK (required_material IN (0, 1)),
  PRIMARY KEY (evidence_id, artifact_id)
);

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

CREATE TABLE stops (
  id TEXT PRIMARY KEY,
  subject_type TEXT NOT NULL CHECK (subject_type IN ('run','envelope','activation')),
  subject_id TEXT NOT NULL,
  code TEXT NOT NULL CHECK (code IN (
    'executor-liveness-unknown','recovery-state-conflict','budget-exhausted',
    'budget-data-unavailable','activation-incomplete','capability-unavailable',
    'operator-stop','review-rework-exhausted','protected-decision-required',
    'human-outcome-required'
  )),
  blocked_actions_json TEXT NOT NULL CHECK (json_valid(blocked_actions_json)),
  details_json TEXT NOT NULL CHECK (json_valid(details_json)),
  resolution_condition_json TEXT NOT NULL CHECK (json_valid(resolution_condition_json)),
  resolver_capability TEXT NOT NULL,
  actor_id TEXT NOT NULL,
  raised_at TEXT NOT NULL,
  created_event_id TEXT NOT NULL
);

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
    'operator-stop','review-rework-exhausted','protected-decision-required',
    'human-outcome-required'
  )),
  status TEXT NOT NULL CHECK (status IN ('open','resolved')),
  version INTEGER NOT NULL CHECK (version > 0),
  raised_at TEXT NOT NULL,
  resolution_event_id TEXT,
  CHECK ((status = 'open' AND resolution_event_id IS NULL)
    OR (status = 'resolved' AND resolution_event_id IS NOT NULL))
);

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

CREATE TABLE review_records (
  id TEXT PRIMARY KEY,
  run_id TEXT NOT NULL UNIQUE REFERENCES runs(id),
  issue_id TEXT NOT NULL REFERENCES issues(id),
  envelope_id TEXT NOT NULL,
  envelope_revision INTEGER NOT NULL,
  actor_id TEXT NOT NULL,
  profile TEXT NOT NULL,
  repository_json TEXT NOT NULL CHECK (json_valid(repository_json)),
  coverage_json TEXT NOT NULL CHECK (json_valid(coverage_json)),
  verdict TEXT NOT NULL CHECK (verdict IN ('pass','changes-required','blocked')),
  evidence_refs_json TEXT NOT NULL CHECK (json_valid(evidence_refs_json)),
  receipt_ids_json TEXT NOT NULL CHECK (json_valid(receipt_ids_json)),
  independence TEXT NOT NULL CHECK (independence = 'distinct-actor'),
  recorded_at TEXT NOT NULL,
  created_event_id TEXT NOT NULL UNIQUE,
  FOREIGN KEY (envelope_id, envelope_revision) REFERENCES envelope_revisions(envelope_id, revision)
);

CREATE TABLE review_findings (
  id TEXT PRIMARY KEY,
  review_id TEXT NOT NULL REFERENCES review_records(id),
  subject_type TEXT NOT NULL CHECK (subject_type IN (
    'issue','envelope','run','review','finding','requirement'
  )),
  subject_id TEXT NOT NULL,
  scope_json TEXT NOT NULL CHECK (json_valid(scope_json)),
  severity TEXT NOT NULL CHECK (severity IN ('blocker','major','minor','nit')),
  summary TEXT NOT NULL,
  blocking INTEGER NOT NULL CHECK (blocking IN (0, 1)),
  created_at TEXT NOT NULL,
  created_event_id TEXT NOT NULL,
  CHECK ((severity IN ('blocker','major') AND blocking = 1)
      OR (severity IN ('minor','nit') AND blocking = 0))
);

CREATE TABLE finding_dispositions (
  id TEXT PRIMARY KEY,
  finding_id TEXT NOT NULL REFERENCES review_findings(id),
  action TEXT NOT NULL CHECK (action IN ('addressed','resolved','inapplicable')),
  actor_id TEXT NOT NULL,
  claim_id TEXT REFERENCES claims(id),
  qualifying_review_id TEXT REFERENCES review_records(id),
  reason TEXT NOT NULL,
  disposed_at TEXT NOT NULL,
  created_event_id TEXT NOT NULL,
  CHECK ((action = 'addressed'
          AND claim_id IS NOT NULL AND qualifying_review_id IS NULL)
      OR (action IN ('resolved','inapplicable')
          AND claim_id IS NULL AND qualifying_review_id IS NOT NULL)),
  UNIQUE (finding_id, claim_id)
);

CREATE TABLE "concerns" (
  id TEXT PRIMARY KEY,
  subject_type TEXT NOT NULL CHECK (subject_type IN (
    'issue','envelope','run','review','finding','requirement'
  )),
  subject_id TEXT NOT NULL,
  actor_id TEXT NOT NULL,
  scope_json TEXT NOT NULL CHECK (json_valid(scope_json)),
  severity TEXT NOT NULL CHECK (severity IN ('blocker','major','minor','nit')),
  summary TEXT NOT NULL,
  reasoning TEXT NOT NULL,
  references_json TEXT NOT NULL CHECK (json_valid(references_json)),
  status TEXT NOT NULL CHECK (status IN ('open','disposed')),
  disposition_action TEXT CHECK (disposition_action IS NULL OR disposition_action IN (
    'carried-forward','rejected','superseded'
  )),
  disposition_issue_id TEXT REFERENCES issues(id),
  supersedes_concern_id TEXT REFERENCES "concerns"(id),
  supersedes_finding_id TEXT REFERENCES review_findings(id),
  disposition_reason TEXT,
  disposed_actor_id TEXT,
  disposed_at TEXT,
  disposed_event_id TEXT,
  created_at TEXT NOT NULL,
  created_event_id TEXT NOT NULL UNIQUE,
  CHECK ((status = 'open'
          AND disposition_action IS NULL AND disposition_reason IS NULL
          AND disposed_actor_id IS NULL AND disposed_at IS NULL
          AND disposed_event_id IS NULL)
      OR (status = 'disposed'
          AND disposition_action IS NOT NULL AND disposition_reason IS NOT NULL
          AND disposed_actor_id IS NOT NULL AND disposed_at IS NOT NULL
          AND disposed_event_id IS NOT NULL)),
  CHECK ((disposition_action IS NULL
          AND disposition_issue_id IS NULL AND supersedes_concern_id IS NULL
          AND supersedes_finding_id IS NULL)
      OR (disposition_action = 'carried-forward'
          AND disposition_issue_id IS NOT NULL AND supersedes_concern_id IS NULL
          AND supersedes_finding_id IS NULL)
      OR (disposition_action = 'rejected'
          AND disposition_issue_id IS NULL AND supersedes_concern_id IS NULL
          AND supersedes_finding_id IS NULL)
      OR (disposition_action = 'superseded'
          AND disposition_issue_id IS NULL
          AND ((supersedes_concern_id IS NOT NULL AND supersedes_finding_id IS NULL)
            OR (supersedes_concern_id IS NULL AND supersedes_finding_id IS NOT NULL)))),
  CHECK (supersedes_concern_id IS NULL OR supersedes_concern_id <> id)
);

CREATE TABLE deviations (
  id TEXT PRIMARY KEY,
  subject_type TEXT NOT NULL CHECK (subject_type IN (
    'issue','envelope','run','review','finding','requirement'
  )),
  subject_id TEXT NOT NULL,
  actor_id TEXT NOT NULL,
  expected TEXT NOT NULL,
  observed TEXT NOT NULL,
  authority_actor_id TEXT NOT NULL,
  authority_reference TEXT,
  disposition TEXT NOT NULL CHECK (disposition IN ('accepted','rejected','pending-decision')),
  reasoning TEXT NOT NULL,
  references_json TEXT NOT NULL CHECK (json_valid(references_json)),
  supersedes_deviation_id TEXT REFERENCES deviations(id),
  created_at TEXT NOT NULL,
  created_event_id TEXT NOT NULL UNIQUE,
  CHECK (supersedes_deviation_id IS NULL OR supersedes_deviation_id <> id)
);

CREATE TABLE decision_refs (
  id TEXT PRIMARY KEY,
  class TEXT NOT NULL,
  source_json TEXT NOT NULL CHECK (json_valid(source_json)),
  authority_actor_id TEXT NOT NULL,
  authority_capabilities_json TEXT NOT NULL CHECK (json_valid(authority_capabilities_json)),
  outcome TEXT NOT NULL CHECK (outcome IN ('approved','rejected','deferred')),
  scope_issue_id TEXT NOT NULL REFERENCES issues(id),
  created_at TEXT NOT NULL,
  created_event_id TEXT NOT NULL UNIQUE
);

CREATE TABLE human_outcomes (
  id TEXT PRIMARY KEY,
  actor_id TEXT NOT NULL,
  activation_id TEXT NOT NULL REFERENCES kit_activations(id),
  activation_version INTEGER NOT NULL CHECK (activation_version > 0),
  profile_ref TEXT NOT NULL,
  class TEXT NOT NULL,
  subject_type TEXT NOT NULL CHECK (subject_type IN (
    'issue','envelope','run','review','finding','requirement'
  )),
  subject_id TEXT NOT NULL,
  result TEXT NOT NULL CHECK (result IN (
    'approved','rejected','changes-required','acknowledged'
  )),
  scope_issue_id TEXT NOT NULL REFERENCES issues(id),
  reasoning TEXT NOT NULL,
  references_json TEXT NOT NULL CHECK (json_valid(references_json)),
  created_at TEXT NOT NULL,
  created_event_id TEXT NOT NULL UNIQUE
);

CREATE TABLE envelope_protected_decisions (
  envelope_id TEXT NOT NULL,
  revision INTEGER NOT NULL,
  ordinal INTEGER NOT NULL CHECK (ordinal >= 0),
  decision_class TEXT NOT NULL,
  required_capabilities_json TEXT NOT NULL CHECK (json_valid(required_capabilities_json)),
  affected_actions_json TEXT NOT NULL CHECK (json_valid(affected_actions_json)),
  PRIMARY KEY (envelope_id, revision, ordinal),
  UNIQUE (envelope_id, revision, decision_class),
  FOREIGN KEY (envelope_id, revision) REFERENCES envelope_revisions(envelope_id, revision)
);

CREATE TABLE requirement_refs (
  id TEXT PRIMARY KEY,
  requirement_id TEXT NOT NULL,
  source_json TEXT NOT NULL CHECK (json_valid(source_json)),
  status TEXT NOT NULL CHECK (status IN ('active','superseded','withdrawn')),
  superseded_by_id TEXT REFERENCES requirement_refs(id),
  created_at TEXT NOT NULL,
  created_event_id TEXT NOT NULL,
  CHECK ((status = 'superseded' AND superseded_by_id IS NOT NULL)
      OR (status IN ('active','withdrawn') AND superseded_by_id IS NULL)),
  CHECK (superseded_by_id IS NULL OR superseded_by_id <> id)
);

CREATE TABLE requirement_issue_links (
  requirement_id TEXT NOT NULL,
  issue_id TEXT NOT NULL REFERENCES issues(id),
  relation TEXT NOT NULL CHECK (relation IN ('implements','verifies','relates-to')),
  created_at TEXT NOT NULL,
  created_event_id TEXT NOT NULL,
  PRIMARY KEY (requirement_id, issue_id, relation)
);

CREATE TABLE "evidence_records" (
  id TEXT PRIMARY KEY,
  obligation_id TEXT NOT NULL REFERENCES envelope_obligations(id),
  evidence_type TEXT NOT NULL CHECK (evidence_type IN ('receipt','review','human-outcome')),
  receipt_id TEXT REFERENCES execution_receipts(id),
  review_id TEXT REFERENCES review_records(id),
  human_outcome_id TEXT REFERENCES human_outcomes(id),
  producer TEXT NOT NULL CHECK (producer IN ('fid','activated-actor','human')),
  trust_source TEXT NOT NULL CHECK (trust_source IN (
    'fid-executed','activated-actor','activated-human'
  )),
  target_repository_json TEXT NOT NULL CHECK (json_valid(target_repository_json)),
  validity TEXT NOT NULL CHECK (validity IN ('valid','invalid')),
  created_at DATETIME NOT NULL,
  CHECK ((evidence_type = 'receipt' AND producer = 'fid'
          AND trust_source = 'fid-executed'
          AND receipt_id IS NOT NULL AND review_id IS NULL AND human_outcome_id IS NULL)
      OR (evidence_type = 'review' AND producer = 'activated-actor'
          AND trust_source = 'activated-actor'
          AND receipt_id IS NULL AND review_id IS NOT NULL AND human_outcome_id IS NULL)
      OR (evidence_type = 'human-outcome' AND producer = 'human'
          AND trust_source = 'activated-human'
          AND receipt_id IS NULL AND review_id IS NULL AND human_outcome_id IS NOT NULL))
);

CREATE UNIQUE INDEX idx_workflow_assignments_kind
  ON workflow_assignments(kind) WHERE scope = 'kind';

CREATE UNIQUE INDEX idx_workflow_assignments_project
  ON workflow_assignments(scope) WHERE scope = 'project';

CREATE UNIQUE INDEX kit_activations_one_nonterminal
  ON kit_activations((1)) WHERE status IN ('prepared','staged','committed');

CREATE UNIQUE INDEX kit_activations_one_active
  ON kit_activations((1)) WHERE status = 'finalized';

CREATE UNIQUE INDEX runs_one_active_per_claim
  ON runs(claim_id) WHERE status = 'active';

CREATE UNIQUE INDEX claims_one_active_mutation
  ON claims(issue_id) WHERE status = 'active';

CREATE INDEX run_observations_by_run ON run_observations(run_id, observed_at, id);

CREATE INDEX budget_samples_by_budget
  ON budget_samples(envelope_id, envelope_revision, budget_name, observed_at, id);

CREATE INDEX stops_by_subject ON stops(subject_type, subject_id, raised_at, id);

CREATE UNIQUE INDEX attention_one_open_semantic_key
  ON attention_items(semantic_key) WHERE status = 'open';

CREATE INDEX review_records_by_issue ON review_records(issue_id, recorded_at, id);

CREATE INDEX review_records_by_envelope ON review_records(envelope_id, envelope_revision, recorded_at, id);

CREATE INDEX review_findings_by_review ON review_findings(review_id, id);

CREATE UNIQUE INDEX finding_dispositions_one_terminal
  ON finding_dispositions(finding_id) WHERE action IN ('resolved','inapplicable');

CREATE INDEX finding_dispositions_by_finding
  ON finding_dispositions(finding_id, disposed_at, id);

CREATE INDEX deviations_by_subject
  ON deviations(subject_type, subject_id, created_at, id);

CREATE INDEX decision_refs_by_class
  ON decision_refs(class, scope_issue_id, created_at, id);

CREATE INDEX human_outcomes_by_class
  ON human_outcomes(class, scope_issue_id, created_at, id);

CREATE UNIQUE INDEX requirement_refs_one_active
  ON requirement_refs(requirement_id) WHERE status = 'active';

CREATE INDEX requirement_refs_by_requirement
  ON requirement_refs(requirement_id, created_at, id);

CREATE INDEX requirement_issue_links_by_issue
  ON requirement_issue_links(issue_id, requirement_id);

CREATE TRIGGER requirement_issue_links_ref_guard_insert
BEFORE INSERT ON requirement_issue_links
FOR EACH ROW
WHEN NOT EXISTS (
  SELECT 1 FROM requirement_refs WHERE requirement_id = NEW.requirement_id
)
BEGIN
  SELECT RAISE(ABORT, 'requirement_issue_links: requirement_id must resolve to a recorded reference');
END;

CREATE UNIQUE INDEX envelope_obligations_one_review_profile
  ON envelope_obligations(envelope_id, revision, review_profile)
  WHERE review_profile IS NOT NULL;

CREATE UNIQUE INDEX envelope_obligations_one_outcome_class
  ON envelope_obligations(envelope_id, revision, outcome_class)
  WHERE outcome_class IS NOT NULL;

CREATE UNIQUE INDEX envelope_obligations_one_gate_evidence
  ON envelope_obligations(envelope_id, revision)
  WHERE obligation_type = 'gate-evidence';

CREATE UNIQUE INDEX evidence_records_one_receipt
  ON evidence_records(obligation_id, receipt_id) WHERE receipt_id IS NOT NULL;

CREATE UNIQUE INDEX evidence_records_one_review
  ON evidence_records(obligation_id, review_id) WHERE review_id IS NOT NULL;

CREATE UNIQUE INDEX evidence_records_one_outcome
  ON evidence_records(obligation_id, human_outcome_id) WHERE human_outcome_id IS NOT NULL;

CREATE INDEX concerns_by_subject_status
  ON concerns(subject_type, subject_id, status, created_at, id);

CREATE VIEW issues_enriched AS
WITH open_concerns AS (
  SELECT
    i.id AS issue_id,
    (
      SELECT COUNT(*)
      FROM concerns c
      WHERE c.subject_type = 'issue' AND c.subject_id = i.id AND c.status = 'open'
    ) AS n
  FROM issues i
),
latest_triage AS (
  SELECT
    i.id AS issue_id,
    (
      SELECT te.verdict
      FROM triage_events te
      WHERE te.issue_id = i.id
      ORDER BY te.created_at DESC
      LIMIT 1
    ) AS verdict
  FROM issues i
)
SELECT
  i.*,
  (COALESCE(oc.n, 0) > 0) AS concerns_raised,
  COALESCE(oc.n, 0) AS open_concerns_count,
  CASE lt.verdict
    WHEN 'RESUME' THEN 'resumed'
    WHEN 'BLOCK' THEN 'blocked'
    ELSE NULL
  END AS triage_status
FROM issues i
LEFT JOIN open_concerns oc ON oc.issue_id = i.id
LEFT JOIN latest_triage lt ON lt.issue_id = i.id;

CREATE TRIGGER activity_log_ref_integrity_guard_insert
BEFORE INSERT ON activity_log
FOR EACH ROW
WHEN NEW.ref_type NOT IN ('issue', 'requirement', 'decision')
  OR (NEW.ref_type = 'issue' AND NOT EXISTS (
    SELECT 1 FROM issues WHERE id = NEW.ref_id
  ))
  OR (NEW.ref_type = 'requirement' AND NOT EXISTS (
    SELECT 1 FROM requirement_refs WHERE requirement_id = NEW.ref_id
  ))
  OR (NEW.ref_type = 'decision' AND NOT EXISTS (
    SELECT 1 FROM decision_refs WHERE id = NEW.ref_id
  ))
BEGIN
  SELECT RAISE(ABORT, 'activity_log: ref_type/ref_id must resolve to an existing row');
END;
