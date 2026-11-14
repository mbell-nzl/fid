//go:build integration

package integration

// Independent test-model column/default snapshot through the accepted Phase 23
// Envelope authoring additions. The fake remains independent of fid's appliers.
// The fake uses only these literals and in-memory Effect application; it never
// opens a fid database or calls its appliers. Task16 owns schema16 convergence.
var p22PortableColumns = map[string]string{
	"acceptance_criteria":          "id issue_id ordinal text",
	"activation_actor_profiles":    "activation_id actor_id actor_type producer provenance_ref session_ref capabilities_json",
	"activation_components":        "activation_id component_type component_id component_version content_digest expected_prior_digest status evidence_json error_code",
	"activation_gate_definitions":  "activation_id definition_id definition_version definition_digest definition_json",
	"activation_recoveries":        "activation_id direction actor_attribution reason selected_event_id completed_event_id",
	"activation_resources":         "activation_id component version old_identity new_identity observed_identity state evidence_ref error_code",
	"activity_log":                 "id journal_event_id ref_type ref_id event details actor occurred_at",
	"attention_items":              "id semantic_key subject_type subject_id code status version raised_at resolution_event_id",
	"attention_stop_links":         "attention_id stop_id",
	"budget_samples":               "id envelope_id envelope_revision budget_name run_id actor_id source source_sample_id sample_kind previous_sample_id value measurement_ref_json observed_at created_event_id",
	"claims":                       "id issue_id envelope_id envelope_revision actor_id worktree_id scope_json mode version status acquired_at expires_at released_at release_reason",
	"envelope_authority_grants":    "envelope_id revision ordinal actor_selector_json actions_json scope_json",
	"envelope_budgets":             "envelope_id revision name budget_json",
	"envelope_issue_selection":     "envelope_id revision issue_id primary_issue",
	"envelope_obligations":         "id envelope_id revision obligation_type registration_id registration_version review_profile outcome_class obligation_json",
	"envelope_protected_decisions": "envelope_id revision ordinal decision_class required_capabilities_json affected_actions_json",
	"envelope_revisions":           "envelope_id revision prior_revision outcome non_goals_json risk_result_json repository_scope_json status created_actor_id approved_actor_id created_at activated_at completed_at",
	"envelope_sources":             "envelope_id revision ordinal source_json",
	"epic_metadata":                "issue_id scheduling brief_md_path target_start target_end",
	"evidence_artifact_links":      "evidence_id artifact_id required_material",
	"evidence_records":             "id obligation_id receipt_id review_id human_outcome_id evidence_type producer trust_source target_repository_json validity created_at",
	"execution_attempts":           "id registration_id registration_version obligation_id run_id actor_id start_repository_json resolved_policy_digest producer started_at",
	"execution_receipts":           "id attempt_id trust_source terminal_status exit_code signal timed_out started_at ended_at terminal_repository_json toolchain_json output_excerpt output_digest output_bytes output_truncated redaction_applied rejected_artifacts_json",
	"gate_registrations":           "id version activation_id definition_id definition_version definition_digest binding_digest project_id worktree_id scope_json policy_json authorization_provenance created_event_id created_at",
	"holds":                        "id task_id worktree_id worktree_root actor since expires_at replaces_hold_id ended_at end_reason",
	"idempotency_results":          "operation_scope idempotency_key request_digest result_entity_id result_event_id aggregate_sequence result_json created_at request_field_digests_json",
	"issue_labels":                 "issue_id label",
	"issues":                       "id kind parent_id title workflow_id status priority sequence spec_md_path owner_role owner_actor created_at updated_at revision lane",
	"kit_activations":              "id prior_activation_id kit_identity backend_identity phase status point_of_no_return direction created_at updated_at version",
	"labels":                       "name colour description",
	"links":                        "from_type from_id relation to_type to_id created_at",
	"receipt_artifacts":            "receipt_id artifact_id declared_path",
	"referenced_artifacts":         "id digest media_type size_bytes retention_class availability location_policy created_at expires_at",
	"requirements":                 "id type summary source_doc source_section last_synced_hash status",
	"run_budget_bindings":          "run_id name budget_json",
	"run_observations":             "id run_id kind actor_id observed_at observation_json created_event_id",
	"run_recoveries":               "id prior_run_id successor_run_id prior_claim_id successor_claim_id observation_id reconciliation_json adopted_paths_json authority_json reason created_event_id",
	"run_runtime":                  "run_id version executor_id obligation_id",
	"runs":                         "id issue_id envelope_id envelope_revision claim_id claim_version actor_id start_repository_json terminal_repository_json status started_at ended_at",
	"stop_resolutions":             "stop_id actor_id fact_event_id resolved_at created_event_id",
	"stops":                        "id subject_type subject_id code blocked_actions_json details_json resolution_condition_json resolver_capability actor_id raised_at created_event_id",
	"work_envelopes":               "id current_revision status created_at",
	"workflow_assignments":         "scope kind workflow_id assigned_at",
	"workflow_states":              "workflow_id state category ordinal",
	"workflow_transitions":         "workflow_id from_state to_state",
	"workflows":                    "id name applies_to active created_at",
}

var p22PortableDefaults = map[string]map[string]any{
	"issues":          {"revision": 1},
	"kit_activations": {"point_of_no_return": 0, "version": 1},
	"workflows":       {"active": 1},
}
