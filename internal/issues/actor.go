package issues

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

var legacyInProcessIdempotencySequence uint64

func parseStandaloneActor(value string) (standaloneActor, error) {
	parts := strings.SplitN(value, ":", 2)
	if len(parts) != 2 || !validStandaloneActorType(parts[0]) || !validStandaloneActorID(parts[1]) {
		return standaloneActor{}, fmt.Errorf("Issue actor must match <agent|human|service>:<id>")
	}
	return standaloneActor{id: parts[1], typeName: parts[0]}, nil
}

// standaloneActorForOptions preserves retained in-process test callers while
// CLI transport remains restricted to the typed Actor field.
func standaloneActorForOptions(actor, legacyRole string) (standaloneActor, error) {
	if strings.TrimSpace(actor) != "" {
		return parseStandaloneActor(actor)
	}
	if strings.TrimSpace(legacyRole) != "" {
		return parseStandaloneActor("agent:" + legacyRole)
	}
	return parseStandaloneActor(actor)
}

func idempotencyKeyForOptions(key, legacyRole string) string {
	if strings.TrimSpace(key) != "" || strings.TrimSpace(legacyRole) == "" {
		return key
	}
	return fmt.Sprintf("legacy-in-process-%d", atomic.AddUint64(&legacyInProcessIdempotencySequence, 1))
}

func validStandaloneActorType(value string) bool {
	return value == "agent" || value == "human" || value == "service"
}

func validStandaloneActorID(value string) bool {
	if len(value) < 1 || len(value) > 128 {
		return false
	}
	for index, character := range value {
		if !(character >= 'A' && character <= 'Z' || character >= 'a' && character <= 'z' || character >= '0' && character <= '9' || character == '.' || character == '_' || character == '-') {
			return false
		}
		if index == 0 && !(character >= 'A' && character <= 'Z' || character >= 'a' && character <= 'z' || character >= '0' && character <= '9') {
			return false
		}
	}
	return true
}

func (actor standaloneActor) eventActor() journal.EventActor {
	return journal.EventActor{ActorID: actor.id, ActorType: actor.typeName, Producer: "fid-standalone", ProvenanceRef: "project-local-attribution"}
}

func issueRequestDigest(operation, projectID string, actor standaloneActor, request any) (string, error) {
	payload, err := journal.MarshalCanonicalJSON(map[string]any{"actor": actor, "operation": operation, "project_id": projectID, "request": request})
	if err != nil {
		return "", fmt.Errorf("canonicalize %s request: %w", operation, err)
	}
	sum := sha256.Sum256(payload)
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}
