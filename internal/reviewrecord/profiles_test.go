package reviewrecord

import (
	"slices"
	"testing"

	"github.com/mbell-nzl/fid/internal/handoffrecord"
	"github.com/mbell-nzl/fid/internal/policy"
)

func TestRequiredProfilesForFrozenLaneAndPaths(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, lane            string
		paths, patterns, want []string
	}{
		{"routine docs", "routine", []string{"README.md"}, []string{"**/*_test.go"}, []string{}},
		{"routine tests", "routine", []string{"internal/example_test.go"}, []string{"**/*_test.go"}, []string{"test-adequacy"}},
		{"standard code", "standard", []string{"internal/example.go"}, []string{"**/*_test.go"}, []string{"code"}},
		{"standard tests", "standard", []string{"internal/example_test.go"}, []string{"**/*_test.go"}, []string{"code", "test-adequacy"}},
		{"assured tests", "assured", []string{"internal/example_test.go"}, []string{"**/*_test.go"}, []string{"code", "test-adequacy"}},
		{"assured empty", "assured", []string{"internal/example_test.go"}, []string{}, []string{"code"}},
		{"unsupported path", "routine", []string{"../test.go"}, []string{"**/*_test.go"}, []string{"test-adequacy"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := RequiredProfiles(handoffrecord.Record{Lane: tc.lane, ChangedPaths: tc.paths}, policy.Effective{TestAdequacySet: true, TestAdequacy: tc.patterns})
			if err != nil || !slices.Equal(got, tc.want) {
				t.Fatalf("profiles=%v/%v want=%v", got, err, tc.want)
			}
		})
	}
}
