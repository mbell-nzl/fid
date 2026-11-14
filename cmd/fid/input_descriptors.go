package main

import (
	"reflect"

	"github.com/mbell-nzl/fid/internal/artifactverify"
	"github.com/mbell-nzl/fid/internal/cli"
)

var inputDescriptorRegistry = map[string]cli.InputDescriptor{
	"artifact verify": {Operation: "artifact.verify", Purpose: "Observe retained Artifact availability.", Type: reflect.TypeFor[artifactverify.Request](), Fields: []cli.InputField{{Name: "artifact_id", Type: "string", Requirement: "required", Meaning: "Artifact ID.", Format: "ID", Placeholder: "<artifact-id>"}}},
}
