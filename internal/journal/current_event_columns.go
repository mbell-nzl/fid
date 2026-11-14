package journal

// Column domains copied from the accepted migration baseline. No per-Event
// schema introspection is performed. The contract test checks these against DDL.
var currentEventTables = map[string]currentTable{
	"handoffs":             currentTableColumns("id", "id:s task_id:s hold_id:s policy_snapshot:i record_json:s"),
	"check_requests":       currentTableColumns("id", "id:s task_id:s hold_id:s actor:s names_json:s response_json:s?"),
	"check_starts":         currentTableColumns("id", "id:s request_id:s task_id:s name:s start_json:s"),
	"check_processes":      currentTableColumns("check_id", "check_id:s process_json:s"),
	"check_results":        currentTableColumns("check_id", "check_id:s result_json:s"),
	"check_files":          currentTableColumns("check_id path", "check_id:s path:s artifact_id:s? reason:s?"),
	"questions":            currentTableColumns("id", "id:s task_id:s hold_id:s text:s asker:s asked_at:s resolved_at:s? resolution:s? decision_id:s?"),
	"person_decisions":     currentTableColumns("id", "id:s kind:s actor:s subject_type:s subject_id:s reason:s binding_id:s decided_at:s"),
	"policy_snapshots":     currentTableColumns("sequence", "sequence:i source_ref:s source_commit:s? digest:s effective_json:s"),
	"activity_log":         currentTableColumns("id", "id:s journal_event_id:s ref_type:s ref_id:s event:s details:s? actor:s? occurred_at:s?"),
	"holds":                currentTableColumns("id", "id:s task_id:s worktree_id:s worktree_root:s actor:s since:s expires_at:s replaces_hold_id:s? ended_at:s? end_reason:s? last_heartbeat_at:s? kept:b park_reason:s?"),
	"idempotency_results":  currentTableColumns("operation_scope idempotency_key", "operation_scope:s idempotency_key:s request_digest:s result_entity_id:s result_event_id:s aggregate_sequence:i result_json:s created_at:s request_field_digests_json:s"),
	"referenced_artifacts": currentTableColumns("id", "id:s digest:s media_type:s size_bytes:i retention_class:s availability:s location_policy:s created_at:s expires_at:s?"),
}
