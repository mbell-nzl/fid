package project

import "testing"

func TestExpectedRepositoryChecksOnlyHeadAndDirty(t *testing.T) {
	t.Parallel()
	observed := RecordedRepository{HeadOID: "0123456789abcdef0123456789abcdef01234567", Dirty: true, TreeOID: "tree", WorktreeID: "worktree"}
	for _, test := range []struct {
		name     string
		expected ExpectedRepository
		matches  bool
	}{
		{"match", ExpectedRepository{HeadOID: observed.HeadOID, Dirty: true}, true},
		{"HEAD differs", ExpectedRepository{HeadOID: "another", Dirty: true}, false},
		{"dirty differs", ExpectedRepository{HeadOID: observed.HeadOID, Dirty: false}, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			if got := test.expected.MatchesRecordedRepository(observed); got != test.matches {
				t.Fatalf("MatchesRecordedRepository() = %t, want %t", got, test.matches)
			}
		})
	}
}
