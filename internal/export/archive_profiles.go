package export

// The portable archive is the surviving capability from the former profile.
// Lean work records travel as Events, without inventing new capability IDs.
func capabilitySnapshot() []capability {
	return []capability{{true, "export.portable.v1", map[string]any{"create": true, "import": true, "validate": true, "import_mode": "empty-target"}, 1, 1, "fid-journal"}}
}
