package export

func capabilitySnapshot() []capability {
	return []capability{
		{true, "attention.v1", map[string]any{"state": true, "human_outcome_resolution": false}, 1, 1, "fid-journal"},
		{true, "evidence.freshness.v1", map[string]any{"branch_labels": false, "target": "repository-tree"}, 1, 1, "fid-git-observation"},
		{true, "evidence.gate-definitions.v1", map[string]any{"executor_policy": "standard-library", "network": "allowed", "sandbox": "none", "write": "allowed"}, 1, 1, "active-kit-activation"},
		{true, "evidence.receipts.v1", map[string]any{"attestation": false, "producer": "fid", "terminal_receipts": true}, 1, 1, "fid-executed"},
		{true, "export.portable.v1", map[string]any{"create": true, "import": true, "validate": true, "import_mode": "empty-target"}, 1, 1, "fid-journal"},
		{false, "human.outcomes.v1", map[string]any{"record": false}, 1, 1, "phase-23-unavailable"},
		{true, "kit.activation.v1", map[string]any{"abort_before_commit": true, "automatic_respawn": false, "finalize": true, "recover": true}, 1, 1, "active-kit-activation"},
		{false, "reviews.coverage.v1", map[string]any{"record": false}, 1, 1, "phase-23-unavailable"},
		{false, "telemetry.trajectory.v1", map[string]any{"ingest": false}, 1, 1, "phase-25-unavailable"},
		{false, "work.claims.v1", map[string]any{"max_active_per_issue": 1, "parallelism": false, "recover": true}, 1, 1, "phase-24-unavailable"},
		{true, "workflow.state.v1", map[string]any{"max_envelope_revisions": 1, "max_primary_issues": 1, "related_issues": false}, 1, 1, "fid-journal"},
		{true, "workflow.transitions.v1", map[string]any{"execution_bound": true, "terminal_evidence_guard": true}, 1, 1, "fid-journal"},
	}
}

// legacyCapabilitySnapshot is the accepted Phase 21 manifest-v1 profile at
// 35a75ba (projection 14). Future capability adoption does not rewrite it.
func legacyCapabilitySnapshot() []capability {
	return []capability{
		{false, "attention.v1", map[string]any{"state": false}, 1, 1, "phase-22-unavailable"},
		{true, "evidence.freshness.v1", map[string]any{"branch_labels": false, "target": "repository-tree"}, 1, 1, "fid-git-observation"},
		{true, "evidence.gate-definitions.v1", map[string]any{"executor_policy": "standard-library", "network": "allowed", "sandbox": "none", "write": "allowed"}, 1, 1, "active-kit-activation"},
		{true, "evidence.receipts.v1", map[string]any{"attestation": false, "producer": "fid", "terminal_receipts": true}, 1, 1, "fid-executed"},
		{false, "export.portable.v1", map[string]any{"create": true, "import": false, "validate": true}, 1, 1, "phase-22-unavailable"},
		{false, "human.outcomes.v1", map[string]any{"record": false}, 1, 1, "phase-23-unavailable"},
		{false, "kit.activation.v1", map[string]any{"abort_before_commit": true, "automatic_respawn": false, "finalize": true, "recover": false}, 1, 1, "phase-22-unavailable"},
		{false, "reviews.coverage.v1", map[string]any{"record": false}, 1, 1, "phase-23-unavailable"},
		{false, "telemetry.trajectory.v1", map[string]any{"ingest": false}, 1, 1, "phase-25-unavailable"},
		{false, "work.claims.v1", map[string]any{"max_active_per_issue": 1, "parallelism": false, "recover": false}, 1, 1, "phase-22-unavailable"},
		{true, "workflow.state.v1", map[string]any{"max_envelope_revisions": 1, "max_primary_issues": 1, "related_issues": false}, 1, 1, "fid-journal"},
		{true, "workflow.transitions.v1", map[string]any{"execution_bound": true, "terminal_evidence_guard": true}, 1, 1, "fid-journal"},
	}
}
