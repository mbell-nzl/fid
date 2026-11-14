package main

import (
	"path/filepath"
	"regexp"
	"testing"
)

func TestLaneMutationGolden(t *testing.T) {
	t.Parallel()
	for _, row := range []struct {
		name string
		lane *string
		args []string
	}{
		{"create", nil, []string{"issue", "create", "--kind", "task", "--title", "Retry a rejected backup push", "--lane", "standard", "--actor", "human:mark", "--idempotency-key", "k1"}},
		{"raise", stringPointer("standard"), []string{"issue", "update", "FID-9f13aa", "--lane", "assured", "--actor", "agent:builder", "--idempotency-key", "k2"}},
		{"person-lower-todo", stringPointer("standard"), []string{"issue", "update", "FID-9f13aa", "--lane", "routine", "--actor", "human:mark", "--idempotency-key", "k2"}},
		{"agent-set-missing", nil, []string{"issue", "update", "FID-9f13aa", "--lane", "standard", "--actor", "agent:builder", "--idempotency-key", "k2"}},
	} {
		t.Run(row.name, func(t *testing.T) {
			f := plannedStartFixture(t, "todo", row.lane)
			for _, jsonMode := range []bool{true, false} {
				args := row.args
				ext := ".txt"
				if jsonMode {
					args = append([]string{"--json"}, args...)
					ext = ".json"
				}
				result := f.run(args...)
				if result.err != nil || result.stderr != "" {
					t.Fatalf("lane write = %#v", result)
				}
				output := regexp.MustCompile(`"occurred_at":"[^"]+"`).ReplaceAllString(result.stdout, `"occurred_at":"<TIME>"`)
				if row.name == "create" {
					output = regexp.MustCompile(`FID-[0-9a-f]{6}`).ReplaceAllString(output, "FID-<new>")
				}
				lifecycleGolden(t, f, filepath.Join("..", "start", "lane-"+row.name+ext), output)
			}
		})
	}
}

func TestShowTaskLaneGolden(t *testing.T) {
	t.Parallel()
	f := plannedStartFixture(t, "todo", stringPointer("standard"))
	f.golden(filepath.Join("start", "show-standard.txt"), "show", "FID-9f13aa")
	f.golden(filepath.Join("start", "show-standard.json"), "--json", "show", "FID-9f13aa")
}
