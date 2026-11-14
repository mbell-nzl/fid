package handoffrecord

import "testing"

func TestGitHubPullRequestURLs(t *testing.T) {
	t.Parallel()
	for _, row := range []struct {
		pr    string
		valid bool
	}{
		{"https://github.com/mbell-nzl/fid/pull/63", true},
		{"https://github.com/owner/repo.name-1/pull/1", true},
		{"http://github.com/owner/repo/pull/1", false},
		{"https://example.com/owner/repo/pull/1", false},
		{"https://github.com/owner/repo/pull/0", false},
		{"https://github.com/owner/repo/pull/1/files", false},
		{"https://github.com/owner/repo/pull/1?x=1", false},
		{"https://github.com/owner/repo/pull/1#discussion", false},
		{"https://github.com/owner/repo/issues/1", false},
		{"https://github.com/owner/repo/pull/-1", false},
		{"https://github.com/owner/repo/pull/01", false},
		{"https://github.com/owner/repo/pull/1\n", false},
		{"https://github.com@evil.example/owner/repo/pull/1", false},
	} {
		t.Run(row.pr, func(t *testing.T) {
			t.Parallel()
			if ValidPR(row.pr) != row.valid {
				t.Fatalf("ValidPR(%q)=%v", row.pr, !row.valid)
			}
		})
	}
}
