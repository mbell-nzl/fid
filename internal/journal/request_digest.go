package journal

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
)

// ComputeRequestDigest hashes the caller's request and each named public input.
// Components retain the operation, project and resolved Actor for the family.
// Aliases map encoded request members to flag or positional names where needed.
func ComputeRequestDigest(components map[string]any, actor any, aliases map[string]string) (string, map[string]string, error) {
	request, ok := components["request"]
	if !ok {
		return "", nil, fmt.Errorf("request digest: request component is required")
	}
	digestComponents := make(map[string]any, len(components)+1)
	for name, value := range components {
		digestComponents[name] = value
	}
	digestComponents["actor_input"] = actor
	payload, err := MarshalCanonicalJSON(digestComponents)
	if err != nil {
		return "", nil, fmt.Errorf("request digest: canonicalize components: %w", err)
	}
	requestJSON, err := MarshalCanonicalJSON(request)
	if err != nil {
		return "", nil, fmt.Errorf("request digest: canonicalize request: %w", err)
	}
	var members map[string]json.RawMessage
	if err := json.Unmarshal(requestJSON, &members); err != nil || members == nil {
		return "", nil, fmt.Errorf("request digest: request must be an object")
	}
	fields := make(map[string]string, len(members)+1)
	for name, value := range members {
		publicName := name
		if alias, found := aliases[name]; found {
			publicName = alias
		}
		if _, exists := fields[publicName]; exists || publicName == "--actor" {
			return "", nil, fmt.Errorf("request digest: duplicate public field %q", publicName)
		}
		fields[publicName] = sha256Digest(value)
	}
	actorJSON, err := MarshalCanonicalJSON(actor)
	if err != nil {
		return "", nil, fmt.Errorf("request digest: canonicalize Actor: %w", err)
	}
	fields["--actor"] = sha256Digest(actorJSON)
	return sha256Digest(payload), fields, nil
}

func sha256Digest(payload []byte) string {
	sum := sha256.Sum256(payload)
	return "sha256:" + hex.EncodeToString(sum[:])
}

// DifferingRequestFields reports each field whose value or presence differs.
func DifferingRequestFields(stored, supplied map[string]string) []string {
	differing := make([]string, 0)
	for name, digest := range stored {
		if other, found := supplied[name]; !found || digest != other {
			differing = append(differing, name)
		}
	}
	for name := range supplied {
		if _, found := stored[name]; !found {
			differing = append(differing, name)
		}
	}
	sort.Strings(differing)
	return differing
}
