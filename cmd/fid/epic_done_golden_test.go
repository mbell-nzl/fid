package main

import (
	"path/filepath"
	"testing"
)

func TestEpicNeedsYouGolden(t *testing.T) {
	t.Parallel()
	for _, format := range []string{"txt", "json"} {
		t.Run(format, func(t *testing.T) {
			t.Parallel()
			f := newStatusFixture(t, 0)
			f.seed([]readyCommandRow{
				statusEpic("FID-e00001", "Finished epic", "todo", nil),
				statusTask("FID-a00001", "Cancelled task", "cancelled", stringPointer("FID-e00001"), nil, nil),
				statusEpic("FID-e00002", "Empty epic", "todo", nil),
				statusEpic("FID-e00003", "Open epic", "in_progress", nil),
				statusTask("FID-a00003", "Task in review", "in_review", stringPointer("FID-e00003"), nil, nil),
				statusEpic("FID-e00004", "Closed epic", "done", nil),
				statusTask("FID-a00004", "Finished task", "done", stringPointer("FID-e00004"), nil, nil),
				statusEpic("FID-e00005", "Cancelled epic", "cancelled", nil),
				statusTask("FID-a00005", "Cancelled task", "cancelled", stringPointer("FID-e00005"), nil, nil),
			}, nil)
			args := []string{"status"}
			if format == "json" {
				args = append(args, "--json")
			}
			f.golden(filepath.Join("epic-done", "status."+format), args...)
		})
	}
}
