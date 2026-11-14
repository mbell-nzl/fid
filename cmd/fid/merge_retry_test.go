package main

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/mergework"
)

func TestMergeObservationRetryPreservesOriginalFacts(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"identical", "task", "PR", "number", "repository", "branch", "head", "tree", "merged time"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			f, options, handoff, pr, tree := mergeFixture(t)
			first, err := mergework.RecordMerge(options, handoff, pr, tree)
			if err != nil {
				t.Fatal(err)
			}
			// Keep the natural key unchanged and retry at a later observation time.
			options.Clock = func() time.Time { return statusNow.Add(2 * time.Minute) }
			switch name {
			case "task":
				_, created := runFIDCommand(t, f.workspace, "issue", "create", "--kind", "task", "--title", "Other task", "--actor", "human:mark")
				options.TaskID = parseCreatedIssueID(t, created)
			case "PR":
				pr.URL = "https://github.com/mbell-nzl/fid/pull/64"
				pr.Number = 64
			case "number":
				pr.Number++
			case "repository":
				pr.BaseRepository.Name = "other"
				pr.URL = "https://github.com/mbell-nzl/other/pull/63"
			case "branch":
				pr.BaseRef = "other"
			case "head":
				pr.HeadSHA = strings.Repeat("f", 40)
			case "tree":
				tree = strings.Repeat("b", 40)
			case "merged time":
				changed := pr.MergedAt.Add(time.Second)
				pr.MergedAt = &changed
			}
			before := snapshotWorkspacePersistence(t, f.workspace)
			again, err := mergework.RecordMerge(options, handoff, pr, tree)
			if name == "identical" {
				if err != nil || !reflect.DeepEqual(again, first) {
					t.Errorf("identical retry = %+v/%v, want original %+v", again, err, first)
				}
			} else if err == nil || !strings.Contains(err.Error(), "retry differs from original") {
				t.Errorf("conflicting retry = %+v/%v, want conflict refusal", again, err)
			}
			if after := snapshotWorkspacePersistence(t, f.workspace); !before.equal(after) {
				t.Fatal("retry changed journal or projection")
			}
		})
	}
}
