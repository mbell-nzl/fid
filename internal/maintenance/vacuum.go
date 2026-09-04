package maintenance

// VacuumResult reports the derived SQLite projection size before and after compaction.
type VacuumResult struct {
	BeforeBytes int64
	AfterBytes  int64
}

// Vacuum compacts the derived SQLite projection without changing journal history.
func Vacuum(workspaceRoot, migrationsDir, byRole string) (VacuumResult, error) {
	return VacuumResult{}, &CapabilityUnavailableError{CapabilityID: "telemetry.trajectory.v1", OwnerPhase: 25}
}
