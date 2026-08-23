package evidence

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type GateDefinition struct {
	ArtifactBytes     int64              `json:"artifact_bytes"`
	ArtifactPaths     []string           `json:"artifact_paths"`
	Argv              []string           `json:"argv"`
	CredentialNames   []string           `json:"credential_names"`
	CWD               string             `json:"cwd"`
	DefinitionID      string             `json:"definition_id"`
	DefinitionVersion int                `json:"definition_version"`
	Environment       []EnvironmentEntry `json:"environment"`
	ExcerptBytes      int64              `json:"excerpt_bytes"`
	Network           string             `json:"network"`
	OutputBytes       int64              `json:"output_bytes"`
	RedactLiterals    []string           `json:"redact_literals"`
	Sandbox           string             `json:"sandbox"`
	Shell             bool               `json:"shell"`
	SuccessExitCodes  []int              `json:"success_exit_codes"`
	TimeoutMS         int64              `json:"timeout_ms"`
	Write             string             `json:"write"`
}

type EnvironmentEntry struct {
	Name   string  `json:"name"`
	Source string  `json:"source"`
	Value  *string `json:"value"`
}

type ExecutorPolicyUnsupportedError struct{ Fields []string }

func (err *ExecutorPolicyUnsupportedError) Error() string {
	return fmt.Sprintf("executor policy is unsupported for fields: %s", strings.Join(err.Fields, ", "))
}

var environmentNamePattern = regexp.MustCompile(`^[A-Z_][A-Z0-9_]{0,127}$`)

var inheritedNames = map[string]bool{
	"CI": true, "GOCACHE": true, "GOMODCACHE": true, "GOPATH": true, "HOME": true,
	"LANG": true, "LC_ALL": true, "PATH": true, "TMPDIR": true, "XDG_CACHE_HOME": true,
}

func ValidateGateDefinition(workspaceRoot string, definition GateDefinition) error {
	fields := make([]string, 0)
	if strings.TrimSpace(definition.DefinitionID) == "" || definition.DefinitionVersion < 1 {
		fields = append(fields, "definition")
	}
	if len(definition.Argv) == 0 || len(definition.Argv) > 128 || !allBoundedStrings(definition.Argv, 1, 4096) {
		fields = append(fields, "argv")
	}
	if definition.Shell {
		fields = append(fields, "shell")
	}
	if definition.Sandbox != "none" {
		fields = append(fields, "sandbox")
	}
	if definition.Network != "allowed" {
		fields = append(fields, "network")
	}
	if definition.Write != "allowed" {
		fields = append(fields, "write")
	}
	if definition.TimeoutMS < 1 || definition.TimeoutMS > 600000 {
		fields = append(fields, "timeout_ms")
	}
	if definition.OutputBytes < 1 || definition.OutputBytes > 10485760 {
		fields = append(fields, "output_bytes")
	}
	if definition.ExcerptBytes < 0 || definition.ExcerptBytes > 16384 || definition.ExcerptBytes > definition.OutputBytes {
		fields = append(fields, "excerpt_bytes")
	}
	if definition.ArtifactBytes < 1 || definition.ArtifactBytes > 1073741824 {
		fields = append(fields, "artifact_bytes")
	}
	if !sortedUniqueIntegers(definition.SuccessExitCodes, -255, 255) {
		fields = append(fields, "success_exit_codes")
	}
	credentialsValid := sortedUniqueStrings(definition.CredentialNames)
	if !credentialsValid {
		fields = append(fields, "credential_names")
	}
	if !sortedUniqueStrings(definition.RedactLiterals) {
		fields = append(fields, "redact_literals")
	}
	if credentialsValid && !validEnvironment(definition.Environment, definition.CredentialNames) {
		fields = append(fields, "environment")
	}
	if !sortedUniqueStrings(definition.ArtifactPaths) || !containedPath(workspaceRoot, definition.CWD) {
		fields = append(fields, "cwd")
	}
	for _, path := range definition.ArtifactPaths {
		if !containedPath(workspaceRoot, path) {
			fields = append(fields, "artifact_paths")
			break
		}
	}
	if len(fields) > 0 {
		sort.Strings(fields)
		return &ExecutorPolicyUnsupportedError{Fields: fields}
	}
	return nil
}

func allBoundedStrings(values []string, minimum, maximum int) bool {
	for _, value := range values {
		if len(value) < minimum || len(value) > maximum {
			return false
		}
	}
	return true
}

func sortedUniqueStrings(values []string) bool {
	previous := ""
	for _, value := range values {
		if value == "" || previous >= value {
			return false
		}
		previous = value
	}
	return true
}

func sortedUniqueIntegers(values []int, minimum, maximum int) bool {
	if len(values) == 0 {
		return false
	}
	previous := minimum - 1
	for index, value := range values {
		if value < minimum || value > maximum || (index > 0 && previous >= value) {
			return false
		}
		previous = value
	}
	return true
}

func validEnvironment(entries []EnvironmentEntry, credentialNames []string) bool {
	credentials := make([]string, 0)
	previous := ""
	for _, entry := range entries {
		if !environmentNamePattern.MatchString(entry.Name) || previous >= entry.Name {
			return false
		}
		previous = entry.Name
		switch entry.Source {
		case "literal":
			if entry.Value == nil {
				return false
			}
		case "inherit":
			if entry.Value != nil || !inheritedNames[entry.Name] {
				return false
			}
		case "credential":
			if entry.Value != nil {
				return false
			}
			credentials = append(credentials, entry.Name)
		default:
			return false
		}
	}
	return strings.Join(credentials, "\x00") == strings.Join(credentialNames, "\x00")
}

func containedPath(workspaceRoot, value string) bool {
	if strings.TrimSpace(value) == "" || filepath.IsAbs(value) {
		return false
	}
	clean := filepath.Clean(value)
	if clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return false
	}
	root, err := filepath.EvalSymlinks(workspaceRoot)
	if err != nil {
		return false
	}
	path := filepath.Join(root, clean)
	for {
		resolved, err := filepath.EvalSymlinks(path)
		if err == nil {
			relative, relErr := filepath.Rel(root, resolved)
			return relErr == nil && relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator))
		}
		if !os.IsNotExist(err) {
			return false
		}
		parent := filepath.Dir(path)
		if parent == path {
			return false
		}
		path = parent
	}
}
