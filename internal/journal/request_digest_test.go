package journal

import (
	"reflect"
	"testing"
)

func TestComputeRequestDigestNamesPublicFieldsAndActor(t *testing.T) {
	t.Parallel()
	request := map[string]any{"claim_version": 1, "claim_id": "CLM-1"}
	digest, fields, err := ComputeRequestDigest(map[string]any{
		"operation": "run.start", "project_id": "PRJ-1", "actor_id": "actor-one", "request": request,
	}, "actor-one", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(digest) != len("sha256:")+64 {
		t.Fatalf("request digest = %q", digest)
	}
	want := map[string]bool{"claim_version": true, "claim_id": true, "--actor": true}
	if len(fields) != len(want) {
		t.Fatalf("field names = %v, want %v", fields, want)
	}
	for name, fieldDigest := range fields {
		if !want[name] || len(fieldDigest) != len("sha256:")+64 {
			t.Fatalf("field %q = %q", name, fieldDigest)
		}
	}
	_, changed, err := ComputeRequestDigest(map[string]any{
		"operation": "run.start", "project_id": "PRJ-1", "actor_id": "actor-one",
		"request": map[string]any{"claim_version": 2, "claim_id": "CLM-1"},
	}, "actor-one", nil)
	if err != nil {
		t.Fatal(err)
	}
	if got := DifferingRequestFields(fields, changed); !reflect.DeepEqual(got, []string{"claim_version"}) {
		t.Fatalf("differing fields = %v, want [claim_version]", got)
	}
}

func TestComputeRequestDigestOptionalPresence(t *testing.T) {
	t.Parallel()
	components := func(request map[string]any) map[string]any {
		return map[string]any{"operation": "review.record", "project_id": "PRJ-1", "actor_id": "reviewer", "request": request}
	}
	_, omitted, err := ComputeRequestDigest(components(map[string]any{"profile": "default"}), "reviewer", nil)
	if err != nil {
		t.Fatal(err)
	}
	_, supplied, err := ComputeRequestDigest(components(map[string]any{"profile": "default", "inapplicable_findings": []any{map[string]any{"finding_id": "FND-1", "reason": "n/a"}}}), "reviewer", nil)
	if err != nil {
		t.Fatal(err)
	}
	if got := DifferingRequestFields(omitted, supplied); !reflect.DeepEqual(got, []string{"inapplicable_findings"}) {
		t.Fatalf("optional difference = %v", got)
	}
}

func TestComputeRequestDigestIncludesWrittenActorWithOpaqueResolvedActor(t *testing.T) {
	t.Parallel()
	type resolvedActor struct{ id string }
	components := func() map[string]any {
		return map[string]any{"operation": "issue.update", "project_id": "PRJ-1", "actor": resolvedActor{id: "opaque"}, "request": map[string]any{"title": "same"}}
	}
	first, _, err := ComputeRequestDigest(components(), "agent:developer", nil)
	if err != nil {
		t.Fatal(err)
	}
	second, _, err := ComputeRequestDigest(components(), "agent:other", nil)
	if err != nil {
		t.Fatal(err)
	}
	if first == second {
		t.Fatalf("Actor change retained request digest %q", first)
	}
}
