package journal

// Column domains copied from the accepted migration baseline. No per-Event
// schema introspection is performed. The contract test checks these against DDL.
var currentEventTables = map[string]currentTable{
	"activation_actor_profiles":   currentTableColumns("activation_id actor_id", "activation_id:s actor_id:s actor_type:s producer:s provenance_ref:s session_ref:s? capabilities_json:s"),
	"activation_components":       currentTableColumns("activation_id component_type component_id component_version", "activation_id:s component_type:s component_id:s component_version:i content_digest:s expected_prior_digest:s? status:s evidence_json:s error_code:s?"),
	"activation_gate_definitions": currentTableColumns("activation_id definition_id definition_version", "activation_id:s definition_id:s definition_version:i definition_digest:s definition_json:s"),
	"activation_recoveries":       currentTableColumns("activation_id", "activation_id:s direction:s actor_attribution:s reason:s selected_event_id:s completed_event_id:s?"),
	"activation_resources":        currentTableColumns("activation_id component", "activation_id:s component:s version:i old_identity:s? new_identity:s observed_identity:s? state:s evidence_ref:s? error_code:s?"),
	"attention_items":             currentTableColumns("id", "id:s semantic_key:s subject_type:s subject_id:s code:s status:s version:i raised_at:s resolution_event_id:s?"),
	"attention_stop_links":        currentTableColumns("attention_id stop_id", "attention_id:s stop_id:s"),
	"budget_samples":              currentTableColumns("id", "id:s envelope_id:s envelope_revision:i budget_name:s run_id:s actor_id:s source:s source_sample_id:s sample_kind:s previous_sample_id:s? value:i? measurement_ref_json:s? observed_at:s created_event_id:s"),
	"claims":                      currentTableColumns("id", "id:s issue_id:s envelope_id:s envelope_revision:i actor_id:s worktree_id:s scope_json:s mode:s version:i status:s acquired_at:s expires_at:s released_at:s? release_reason:s?"),
	"execution_attempts":          currentTableColumns("id", "id:s registration_id:s registration_version:i obligation_id:s run_id:s actor_id:s start_repository_json:s resolved_policy_digest:s producer:s started_at:s"),
	"execution_receipts":          currentTableColumns("id", "id:s attempt_id:s trust_source:s terminal_status:s exit_code:i? signal:s? timed_out:b started_at:s ended_at:s terminal_repository_json:s toolchain_json:s output_excerpt:s output_digest:s output_bytes:i output_truncated:b redaction_applied:b rejected_artifacts_json:s"),
	"idempotency_results":         currentTableColumns("operation_scope idempotency_key", "operation_scope:s idempotency_key:s request_digest:s result_entity_id:s result_event_id:s aggregate_sequence:i result_json:s created_at:s"),
	"kit_activations":             currentTableColumns("id", "id:s prior_activation_id:s? kit_identity:s backend_identity:s phase:i status:s point_of_no_return:b direction:s created_at:s updated_at:s version:i"),
	"receipt_artifacts":           currentTableColumns("receipt_id artifact_id", "receipt_id:s artifact_id:s declared_path:s"),
	"referenced_artifacts":        currentTableColumns("id", "id:s digest:s media_type:s size_bytes:i retention_class:s availability:s location_policy:s created_at:s expires_at:s?"),
	"run_budget_bindings":         currentTableColumns("run_id name", "run_id:s name:s budget_json:s"),
	"run_observations":            currentTableColumns("id", "id:s run_id:s kind:s actor_id:s observed_at:s observation_json:s created_event_id:s"),
	"run_recoveries":              currentTableColumns("id", "id:s prior_run_id:s successor_run_id:s prior_claim_id:s successor_claim_id:s observation_id:s reconciliation_json:s adopted_paths_json:s authority_json:s reason:s created_event_id:s"),
	"run_runtime":                 currentTableColumns("run_id", "run_id:s version:i executor_id:s obligation_id:s?"),
	"runs":                        currentTableColumns("id", "id:s issue_id:s envelope_id:s envelope_revision:i claim_id:s claim_version:i actor_id:s start_repository_json:s terminal_repository_json:s? status:s started_at:s ended_at:s?"),
	"stop_resolutions":            currentTableColumns("stop_id", "stop_id:s actor_id:s fact_event_id:s resolved_at:s created_event_id:s"),
	"stops":                       currentTableColumns("id", "id:s subject_type:s subject_id:s code:s blocked_actions_json:s details_json:s resolution_condition_json:s resolver_capability:s actor_id:s raised_at:s created_event_id:s"),
}
