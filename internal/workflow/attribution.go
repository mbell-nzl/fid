package workflow

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"sync/atomic"

	"github.com/mbell-nzl/fid/internal/journal"
)

type standaloneActor struct {
	id       string
	typeName string
}

var legacyIdempotencySequence uint64

func standaloneActorForOptions(value, legacyRole string) (standaloneActor, error) {
	if strings.TrimSpace(value) == "" && strings.TrimSpace(legacyRole) != "" {
		value = "agent:" + strings.TrimSpace(legacyRole)
	}
	parts := strings.SplitN(value, ":", 2)
	if len(parts) != 2 || !validActorType(parts[0]) || !validActorID(parts[1]) {
		return standaloneActor{}, fmt.Errorf("workflow actor must match <agent|human|service>:<id>")
	}
	return standaloneActor{id: parts[1], typeName: parts[0]}, nil
}

func workflowIdempotencyKey(key, legacyRole string) string {
	if strings.TrimSpace(key) != "" || strings.TrimSpace(legacyRole) == "" {
		return key
	}
	return fmt.Sprintf("legacy-in-process-%d", atomic.AddUint64(&legacyIdempotencySequence, 1))
}

func validActorType(value string) bool {
	return value == "agent" || value == "human" || value == "service"
}

func validActorID(value string) bool {
	if len(value) < 1 || len(value) > 128 {
		return false
	}
	for index, character := range value {
		alphaNumeric := character >= 'A' && character <= 'Z' || character >= 'a' && character <= 'z' || character >= '0' && character <= '9'
		if !alphaNumeric && character != '.' && character != '_' && character != '-' || index == 0 && !alphaNumeric {
			return false
		}
	}
	return true
}

func (actor standaloneActor) eventActor() journal.EventActor {
	return journal.EventActor{ActorID: actor.id, ActorType: actor.typeName, Producer: "fid-standalone", ProvenanceRef: "project-local-attribution"}
}

func workflowRequestDigest(operation, projectID string, actor standaloneActor, request any) (string, error) {
	payload, err := journal.MarshalCanonicalJSON(map[string]any{"actor": actor, "operation": operation, "project_id": projectID, "request": request})
	if err != nil {
		return "", fmt.Errorf("canonicalize %s request: %w", operation, err)
	}
	sum := sha256.Sum256(payload)
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}
