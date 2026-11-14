-- ADR-0014 clean cut: journal/archive admission already refuses removed history.
-- Keep Artifact/retry storage, lean records, fixed states and frozen V1 telemetry.
DROP VIEW issues_enriched;
DROP TRIGGER activity_log_ref_integrity_guard_insert;

DROP TABLE activation_actor_profiles;
DROP TABLE activation_components;
DROP TABLE activation_recoveries;
DROP TABLE activation_resources;
DROP TABLE attention_stop_links;
DROP TABLE budget_samples;
DROP TABLE concerns;
DROP TABLE decision_refs;
DROP TABLE deviations;
DROP TABLE envelope_authority_grants;
DROP TABLE envelope_issue_selection;
DROP TABLE envelope_protected_decisions;
DROP TABLE envelope_sources;
DROP TABLE evidence_artifact_links;
DROP TABLE finding_dispositions;
DROP TABLE receipt_artifacts;
DROP TABLE requirement_issue_links;
DROP TABLE requirement_refs;
DROP TABLE run_budget_bindings;
DROP TABLE run_recoveries;
DROP TABLE run_runtime;
DROP TABLE stop_resolutions;
DROP TABLE system_capabilities;
DROP TABLE workflow_assignments;
DROP TABLE attention_items;
DROP TABLE envelope_budgets;
DROP TABLE evidence_records;
DROP TABLE review_findings;
DROP TABLE run_observations;
DROP TABLE stops;
DROP TABLE execution_receipts;
DROP TABLE human_outcomes;
DROP TABLE review_records;
DROP TABLE execution_attempts;
DROP TABLE envelope_obligations;
DROP TABLE runs;
DROP TABLE claims;
DROP TABLE envelope_revisions;
DROP TABLE gate_registrations;
DROP TABLE activation_gate_definitions;
DROP TABLE work_envelopes;
DROP TABLE kit_activations;

CREATE VIEW issues_enriched AS
SELECT i.*,
  CASE (
    SELECT te.verdict FROM triage_events te
    WHERE te.issue_id = i.id ORDER BY te.created_at DESC LIMIT 1
  )
    WHEN 'RESUME' THEN 'resumed'
    WHEN 'BLOCK' THEN 'blocked'
    ELSE NULL
  END AS triage_status
FROM issues i;

CREATE TRIGGER activity_log_ref_integrity_guard_insert
BEFORE INSERT ON activity_log
FOR EACH ROW
WHEN NEW.ref_type != 'issue' OR NOT EXISTS (
  SELECT 1 FROM issues WHERE id = NEW.ref_id
)
BEGIN
  SELECT RAISE(ABORT, 'activity_log: ref_type/ref_id must resolve to an existing row');
END;
