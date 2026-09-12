package workflow

import "testing"

func TestRecoveryReadContextWithoutImportPreservesRecordedLiveness(t *testing.T) {
	options, request := checkpointFixture(t)
	before := inspectRepositorySnapshot(t, options.WorkspaceRoot)
	context, err := ReadRecoveryContext(options.WorkspaceRoot, request.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if context.Import != nil || context.Liveness == nil || context.Liveness.State != "terminated" || context.LivenessHistorical {
		t.Fatalf("context invented history or liveness: %#v", context)
	}
	if after := inspectRepositorySnapshot(t, options.WorkspaceRoot); after != before {
		t.Error("context read changed repository")
	}
}
