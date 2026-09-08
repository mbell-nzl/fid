package export

import (
	"bytes"
	"strings"

	"github.com/mbell-nzl/fid/internal/journal"
)

func validateSemanticManifest(m manifest, entries map[string]scannedEntry, facts archiveSemanticFacts) []Check {
	checks := []Check{}
	for projectID := range facts.ProjectIDs {
		if projectID != m.ProjectID {
			checks = append(checks, check("reference-invalid", nil, "journal Repository project identity differs from manifest"))
		}
	}
	want, err := journal.MarshalCanonicalJSON(facts.Sources)
	got, encodeErr := journal.MarshalCanonicalJSON(m.Sources)
	if err != nil || encodeErr != nil || !bytes.Equal(want, got) {
		checks = append(checks, check("reference-invalid", "manifest.json", "Source inventory differs from journal facts"))
	}
	available := map[string]semanticArtifact{}
	declared := make(map[string]entry, len(m.Entries))
	for _, item := range m.Entries {
		declared[item.Path] = item
	}
	for _, artifact := range facts.Artifacts {
		if artifact.Availability != "available" {
			continue
		}
		name, err := artifactName(artifact.Digest)
		if err != nil {
			checks = append(checks, check("artifact-reference-invalid", nil, "journal Artifact digest is invalid"))
			continue
		}
		path := "artifacts/" + name
		available[path] = artifact
		actual, found := entries[path]
		if !found {
			checks = append(checks, check("entry-missing", path, "available journal Artifact body is absent"))
			continue
		}
		if actual.Digest != artifact.Digest || actual.SizeBytes != artifact.SizeBytes {
			checks = append(checks, check("artifact-reference-invalid", path, "Artifact body differs from immutable journal metadata"))
		}
		if declared[path].MediaType != artifact.MediaType {
			checks = append(checks, check("artifact-reference-invalid", path, "Artifact media type differs from journal metadata"))
		}
	}
	for path := range entries {
		if strings.HasPrefix(path, "artifacts/") {
			if _, found := available[path]; !found {
				checks = append(checks, check("artifact-reference-invalid", path, "Artifact body has no available journal reference"))
			}
		}
	}
	return checks
}
