//go:build integration

package integration

import (
	"encoding/json"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func p22PortableThreeWay(t *testing.T, a *p22FIDRecoveryAdapter, activations []string) (map[string][]byte, map[string]any) {
	t.Helper()
	dir := t.TempDir()
	archive := filepath.Join(dir, "live.zip")
	p21Run(t, a.binary, a.workspace, "--json", "export", "create", archive)
	p21Run(t, a.binary, dir, "--json", "export", "validate", archive)
	live := p22PortableNamedProjection(t, a, a.workspace, activations)
	entries := p22ArchiveEntries(t, archive)
	p21Run(t, a.binary, a.workspace, "--json", "rebuild")
	rebuilt := p22PortableNamedProjection(t, a, a.workspace, activations)
	rebuiltArchive := filepath.Join(dir, "rebuilt.zip")
	p21Run(t, a.binary, a.workspace, "--json", "export", "create", rebuiltArchive)
	target := filepath.Join(t.TempDir(), "destination")
	closureRun(t, filepath.Dir(target), "git", "clone", "-q", a.workspace, target)
	result := p21Run(t, a.binary, target, "--json", "export", "import", archive, "--repository", target, "--idempotency-key", "portable-roundtrip")
	if p21Path(t, result, "result", "import", "project_id") != a.projectID || result["event"] != nil {
		t.Fatalf("import identity/Event: %#v", result)
	}
	imported := p22PortableNamedProjection(t, a, target, activations)
	importedArchive := filepath.Join(dir, "imported.zip")
	p21Run(t, a.binary, target, "--json", "export", "create", importedArchive)
	for name, projection := range map[string]map[string]any{"rebuilt": rebuilt, "imported": imported} {
		if !reflect.DeepEqual(live, projection) {
			t.Fatalf("%s full named views differ: live=%#v got=%#v", name, live, projection)
		}
	}
	for name, path := range map[string]string{"rebuilt": rebuiltArchive, "imported": importedArchive} {
		if !reflect.DeepEqual(entries, p22ArchiveEntries(t, path)) {
			t.Fatalf("%s canonical entry bytes differ", name)
		}
	}
	// Independently rebuild the destination as well; installed delivery metadata
	// is not canonical state and must not change either views or re-export.
	p21Run(t, a.binary, target, "--json", "rebuild")
	if got := p22PortableNamedProjection(t, a, target, activations); !reflect.DeepEqual(live, got) {
		t.Fatalf("imported rebuild changed views: %#v", got)
	}
	return entries, live
}

// Every promised durable dimension is kept as a named JSON read. All three
// projections receive the SAME observed repository input; destination-local
// observation and mutation are separately proven by the built recovery journey.
func p22PortableNamedProjection(t *testing.T, a *p22FIDRecoveryAdapter, workspace string, activations []string) map[string]any {
	t.Helper()
	views := p22NamedRecoveryProjection(t, a, workspace, a.workspace)
	for name, args := range map[string][]string{
		"run.list":     {"run", "list", "--issue", a.issueID, "--limit", "100"},
		"receipt.list": {"receipt", "list", "--limit", "100"},
		"issue.get":    {"issue", "get", a.issueID},
	} {
		response := p21Run(t, a.binary, workspace, append([]string{"--project", workspace, "--json"}, args...)...)
		views[name] = response["result"]
	}
	states := []any{}
	for _, id := range activations {
		response := p21Run(t, a.binary, workspace, "--project", workspace, "--json", "activation", "status", id)
		states = append(states, p21Path(t, response, "result", "activation"))
	}
	views["activation.status"] = states
	return views
}

func p22PortableNamedSummary(views map[string]any, entries map[string][]byte) map[string]any {
	// Named view keys contain dots; traverse their response body explicitly.
	runs := views["run.list"].(map[string]any)["items"]
	statuses := []string{}
	for _, item := range runs.([]any) {
		statuses = append(statuses, item.(map[string]any)["status"].(string))
	}
	sort.Strings(statuses)
	summary := map[string]any{"run_statuses": statuses}
	resume := views["resume"].(map[string]any)
	for _, item := range resume["budgets"].([]any) {
		budget := item.(map[string]any)
		if budget["name"] == "tokens" {
			summary["tokens_state"], summary["tokens_consumed"], summary["tokens_remaining"] = budget["state"], budget["consumed"], budget["remaining"]
		}
	}
	attention := views["attention.list"].(map[string]any)["items"].([]any)
	summary["attention_count"] = len(attention)
	summary["stop_count"] = len(views["stop.list"].(map[string]any)["items"].([]any))
	links := 0
	for _, item := range attention {
		links += len(item.(map[string]any)["stop_ids"].([]any))
	}
	summary["attention_link_count"] = links
	activationStatuses := []string{}
	for _, item := range views["activation.status"].([]any) {
		activationStatuses = append(activationStatuses, item.(map[string]any)["status"].(string))
	}
	sort.Strings(activationStatuses)
	summary["activation_statuses"] = activationStatuses
	availability := []string{}
	seen := map[string]bool{}
	for _, item := range views["receipt.list"].(map[string]any)["items"].([]any) {
		for _, value := range item.(map[string]any)["artifacts"].([]any) {
			artifact := value.(map[string]any)
			id := artifact["id"].(string)
			if !seen[id] {
				availability = append(availability, artifact["availability"].(string))
				seen[id] = true
			}
		}
	}
	sort.Strings(availability)
	summary["artifact_availability"] = availability
	count := 0
	for path := range entries {
		if strings.HasPrefix(path, "artifacts/") {
			count++
		}
	}
	summary["artifact_body_count"] = count
	encoded, _ := json.Marshal(summary)
	_ = json.Unmarshal(encoded, &summary)
	return summary
}
