package main

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestApprovalDisplayPreservesKnownEvidenceWithoutFingerprint(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name             string
		review, approval bool
	}{
		{name: "reviews only", review: true},
		{name: "approvals only", approval: true},
		{name: "neither record"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			f := approvalDisplayFixture(t)
			if tc.review {
				lifecycleResult(t, recordReview(t, f, reviewInput(), "review", "human:courier"), "review.record")
			}
			if tc.approval {
				approveSuccess(t, approveForTest(f, "FID-9f13aa", "FID-9f13aa\n", true))
			}
			head := strings.TrimSpace(identityGit(t, f.workspace, "rev-parse", "HEAD"))
			objects := strings.TrimSpace(identityGit(t, f.workspace, "rev-parse", "--path-format=absolute", "--git-path", "objects"))
			if err := os.Remove(filepath.Join(objects, head[:2], head[2:])); err != nil {
				t.Fatal(err)
			}
			// With no approvals, missing is known. With no rounds, the disclosure
			// is known. Only evidence requiring the lost fingerprint becomes null.
			var decision, flag any
			if !tc.approval {
				decision = map[string]any{"state": "missing"}
			}
			if !tc.review {
				flag = true
			}
			before := startEventCount(t, f.workspace)
			for _, command := range []string{"show", "status"} {
				args, operation := []string{command}, "project.status"
				if command == "show" {
					args, operation = append(args, "FID-9f13aa"), "issue.show"
				}
				result, event := lifecycleResult(t, f.run(append([]string{"--json"}, args...)...), operation)
				if event != nil {
					t.Errorf("%s appended an Event: %v", command, event)
				}
				assertUnavailableApprovalJSON(t, command, result, decision, flag, tc.review || tc.approval)
				got := f.run(args...)
				if got.err != nil {
					t.Fatal(got)
				}
				assertUnavailableApprovalText(t, command, got.stdout, !tc.approval, !tc.review, tc.review || tc.approval)
			}
			if startEventCount(t, f.workspace) != before {
				t.Fatal("display appended an Event")
			}
		})
	}
}

func assertUnavailableApprovalJSON(t *testing.T, command string, result map[string]any, decision, flag any, warning bool) {
	t.Helper()
	fields := result
	if command == "status" {
		section := result["in_progress"].(map[string]any)
		tasks := section["tasks"].([]any)
		for _, raw := range section["epics"].([]any) {
			tasks = append(tasks, raw.(map[string]any)["tasks"].([]any)...)
		}
		if len(tasks) != 1 || tasks[0].(map[string]any)["id"] != "FID-9f13aa" {
			t.Fatalf("status tasks=%v", tasks)
		}
		fields = tasks[0].(map[string]any)
		items := []map[string]any{}
		for _, raw := range result["needs_you"].([]any) {
			item := raw.(map[string]any)
			if item["kind"] == "approval" {
				items = append(items, item)
			}
		}
		if len(items) != 0 {
			t.Errorf("unfinished reviews requested approval: %v", items)
		}
	}
	for key, want := range map[string]any{"approval": decision, "no_other_model_review": flag} {
		if got, exists := fields[key]; !exists || !reflect.DeepEqual(got, want) {
			t.Errorf("%s %s=%v (present=%t), want %v", command, key, got, exists, want)
		}
	}
	foundWarning := false
	for _, raw := range result["warnings"].([]any) {
		if raw.(map[string]any)["kind"] == "approval-unavailable" {
			foundWarning = true
		}
	}
	if foundWarning != warning {
		t.Errorf("%s approval-unavailable=%t, want %t", command, foundWarning, warning)
	}
}

func assertUnavailableApprovalText(t *testing.T, command, text string, missing, disclosure, warning bool) {
	t.Helper()
	if got := strings.Contains(strings.Join(strings.Fields(text), " "), "approval unavailable: fingerprint missing-object"); got != warning {
		t.Errorf("%s approval-unavailable=%t, want %t:\n%s", command, got, warning, text)
	}
	parts := map[string]bool{}
	if command == "show" {
		parts["Approval"] = missing
		parts["Approval    ○ missing"] = missing
		parts["no other-model review"] = disclosure
	} else {
		parts["FID-9f13aa approval missing"] = false
		parts["no other-model review"] = false
		parts["fid approve FID-9f13aa --actor human:<name>"] = false
		parts["--reason <text>"] = false
	}
	for part, want := range parts {
		if strings.Contains(text, part) != want {
			t.Errorf("%s contains %q=%t, want %t:\n%s", command, part, !want, want, text)
		}
	}
}
