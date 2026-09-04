package requirements

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/mbell-nzl/fid/internal/journal"
)

type standaloneActor struct{ id, typeName string }

func standaloneActorForOptions(value string) (standaloneActor, error) {
	parts := strings.SplitN(value, ":", 2)
	if len(parts) != 2 || (parts[0] != "agent" && parts[0] != "human" && parts[0] != "service") || !validActorID(parts[1]) {
		return standaloneActor{}, fmt.Errorf("requirement actor must match <agent|human|service>:<id>")
	}
	return standaloneActor{id: parts[1], typeName: parts[0]}, nil
}

func validActorID(value string) bool {
	if len(value) < 1 || len(value) > 128 {
		return false
	}
	for index, character := range value {
		alphaNumeric := character >= 'A' && character <= 'Z' || character >= 'a' && character <= 'z' || character >= '0' && character <= '9'
		if (!alphaNumeric && character != '.' && character != '_' && character != '-') || (index == 0 && !alphaNumeric) {
			return false
		}
	}
	return true
}

func (actor standaloneActor) eventActor() journal.EventActor {
	return journal.EventActor{ActorID: actor.id, ActorType: actor.typeName, Producer: "fid-standalone", ProvenanceRef: "project-local-attribution"}
}

func requirementRequestDigest(projectID string, actor standaloneActor, request any) (string, error) {
	payload, err := journal.MarshalCanonicalJSON(map[string]any{"actor": actor, "operation": "requirement.sync", "project_id": projectID, "request": request})
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(payload)
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}
