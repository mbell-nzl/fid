package links

import "testing"

func TestLinkErrorsRenderActionableMessages(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
		want string
	}{
		{
			name: "unknown relation",
			err:  (&UnknownRelationError{Relation: "parent_of"}),
			want: `link add: unknown relation "parent_of"`,
		},
		{
			name: "self link",
			err:  (&SelfLinkError{FromID: "FID-1", Relation: "relates_to", ToID: "FID-1"}),
			want: `link add: self link is not allowed: FID-1 relates_to FID-1`,
		},
		{
			name: "missing endpoint",
			err: (&EndpointNotFoundError{
				EndpointType: "requirement",
				EndpointID:   "PRD-LINK-002",
				Relation:     "implements",
				OtherType:    "issue",
				OtherID:      "FID-1",
			}),
			want: `link add: requirement "PRD-LINK-002" does not resolve for relation "implements" between requirement "PRD-LINK-002" and issue "FID-1"`,
		},
		{
			name: "incompatible relation",
			err: (&IncompatibleRelationError{
				Relation: "blocks",
				FromType: "issue",
				FromID:   "FID-1",
				ToType:   "requirement",
				ToID:     "PRD-LINK-002",
			}),
			want: `link add: relation "blocks" is not valid between issue "FID-1" and requirement "PRD-LINK-002"`,
		},
		{
			name: "duplicate link",
			err:  (&DuplicateLinkError{FromID: "FID-1", Relation: "blocks", ToID: "FID-2"}),
			want: `link add: link already exists: FID-1 blocks FID-2`,
		},
		{
			name: "missing link",
			err: (&MissingLinkError{
				Operation: "link remove",
				FromID:    "FID-1",
				Relation:  "blocks",
				ToID:      "FID-2",
			}),
			want: `link remove: link does not exist: FID-1 blocks FID-2`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Fatalf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRegistryHelpersCoverEndpointAndCompatibilityRules(t *testing.T) {
	t.Parallel()

	if _, err := endpointType("requirement", "type", "link add"); err == nil {
		t.Fatal("endpointType(requirement) accepted removed endpoint")
	}
	if got, err := endpointType("", "type", "link add"); err != nil || got != "" {
		t.Fatalf("endpointType(empty) = %q, %v; want empty, nil", got, err)
	}
	if _, err := endpointType("claim", "type", "link add"); err == nil {
		t.Fatal("endpointType(claim) error = nil, want invalid endpoint type")
	}

	if relationAllowedForTypes("relates_to", "issue", "requirement") {
		t.Fatal("relationAllowedForTypes accepted requirement endpoint")
	}
	for _, relation := range []string{"implements", "implemented_by"} {
		if relationAllowedForTypes(relation, "issue", "issue") {
			t.Fatalf("relationAllowedForTypes accepted removed %s relation", relation)
		}
	}
	if !relationAllowedForTypes("discovered", "issue", "issue") {
		t.Fatal("relationAllowedForTypes(discovered, issue, issue) = false, want true")
	}
	if relationAllowedForTypes("blocks", "issue", "requirement") {
		t.Fatal("relationAllowedForTypes(blocks, issue, requirement) = true, want false")
	}
	if relationAllowedForTypes("mystery", "issue", "issue") {
		t.Fatal("relationAllowedForTypes(mystery, issue, issue) = true, want false")
	}

	if !isV01EndpointType("issue") || isV01EndpointType("requirement") {
		t.Fatal("endpoint registry must admit only Issues")
	}
	if isV01EndpointType("decision") {
		t.Fatal("isV01EndpointType(decision) = true, want false")
	}
}
