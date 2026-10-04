package project

// ExpectedRepository contains the caller's two repository preconditions.
type ExpectedRepository struct {
	HeadOID string `json:"head_oid"`
	Dirty   bool   `json:"dirty"`
}

// MatchesRecordedRepository compares only the caller's preconditions with
// the single repository observation used for the Event and stored record.
func (expected ExpectedRepository) MatchesRecordedRepository(observed RecordedRepository) bool {
	return expected.HeadOID == observed.HeadOID && expected.Dirty == observed.Dirty
}
