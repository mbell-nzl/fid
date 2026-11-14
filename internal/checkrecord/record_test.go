package checkrecord

import (
	"slices"
	"testing"
)

func TestQualificationBindsTheTreeAndRecordsEveryMissingCondition(t *testing.T) {
	t.Parallel()
	clean := Tree{Head: "commit-a", Tree: "tree-a"}
	dirty := clean
	dirty.Dirty = true
	digest := "sha256:dirty"
	dirty.DirtyDigest = &digest
	for _, row := range []struct {
		name          string
		start, end    Tree
		current, kept bool
		reasons       []string
	}{
		{"same tree", clean, clean, true, true, []string{}},
		{"another commit with the same tree", clean, Tree{Head: "commit-b", Tree: "tree-a"}, true, true, []string{}},
		{"another tree", clean, Tree{Head: "commit-b", Tree: "tree-b"}, true, true, []string{"tree-changed"}},
		{"all missing conditions", dirty, clean, false, false, []string{"tree-dirty", "tree-changed", "ownership-lost", "artifact-not-kept"}},
	} {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			got := Qualification(row.start, row.end, "passed", row.current, row.kept)
			if !slices.Equal(got, row.reasons) {
				t.Fatalf("reasons=%v, want %v", got, row.reasons)
			}
		})
	}
}
