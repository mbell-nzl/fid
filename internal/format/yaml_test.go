package format

import (
	"bytes"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestWriteYAMLMatchesJSONStructureForEmptyIssue(t *testing.T) {
	issue := emptyIssueFixture()

	var jsonOutput bytes.Buffer
	if err := Write(&jsonOutput, JSON, issue); err != nil {
		t.Fatalf("Write(JSON) error = %v", err)
	}

	var yamlOutput bytes.Buffer
	if err := Write(&yamlOutput, YAML, issue); err != nil {
		t.Fatalf("Write(YAML) error = %v", err)
	}

	want := decodeJSONStructure(t, jsonOutput.String())
	got := decodeYAMLStructure(t, yamlOutput.String())

	if !structuresEqual(want, got) {
		t.Fatalf("YAML structure = %#v, want %#v", got, want)
	}
}

func TestWriteYAMLMatchesJSONStructureForPopulatedIssue(t *testing.T) {
	issue := populatedIssueFixture()

	var jsonOutput bytes.Buffer
	if err := Write(&jsonOutput, JSON, issue); err != nil {
		t.Fatalf("Write(JSON) error = %v", err)
	}

	var yamlOutput bytes.Buffer
	if err := Write(&yamlOutput, YAML, issue); err != nil {
		t.Fatalf("Write(YAML) error = %v", err)
	}

	want := decodeJSONStructure(t, jsonOutput.String())
	got := decodeYAMLStructure(t, yamlOutput.String())

	if !structuresEqual(want, got) {
		t.Fatalf("YAML structure = %#v, want %#v", got, want)
	}
}

func decodeYAMLStructure(t *testing.T, value string) any {
	t.Helper()

	var result any
	if err := yaml.Unmarshal([]byte(value), &result); err != nil {
		t.Fatalf("Unmarshal YAML structure error = %v", err)
	}
	return result
}
