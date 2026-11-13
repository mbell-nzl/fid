package journal

import (
	"bytes"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/mbell-nzl/fid/internal/policy"
)

var policyCommitPattern = regexp.MustCompile(`^([0-9a-f]{40}|[0-9a-f]{64})$`)

func validatePolicySnapshot(tx *sql.Tx, event Event) error {
	if event.EventType != "policy.snapshotted" {
		return nil
	}
	if len(event.Effects) != 1 || event.CommandName != "policy.snapshot" {
		return fmt.Errorf("policy snapshot requires exactly one policy Effect and its internal command")
	}
	if event.Aggregate.ID != "policy" || event.ActorV2.ActorID != "fid" || event.ActorV2.ActorType != "service" || event.Idempotency != nil || event.Repository != nil {
		return fmt.Errorf("policy snapshot must be written by fid on the policy aggregate")
	}
	effect := event.Effects[0]
	if err := validateCurrentEffect(event, effect); err != nil {
		return err
	}
	if err := validatePolicyPayload(effect.Values["effective_json"].(string), effect.Values["digest"]); err != nil {
		return err
	}
	if !currentValuesEqual(effect.PK["sequence"], event.Aggregate.Sequence) {
		return fmt.Errorf("policy snapshot sequence differs from aggregate")
	}
	ref, _ := effect.Values["source_ref"].(string)
	commit := currentEventValue(effect.Values["source_commit"])
	if ref == "built-in defaults" {
		if commit != nil {
			return fmt.Errorf("built-in policy snapshot must not have a source commit")
		}
	} else if text, ok := commit.(string); !ok || !policyCommitPattern.MatchString(text) || !strings.HasPrefix(ref, "origin/") || len(ref) <= len("origin/") {
		return fmt.Errorf("policy snapshot requires origin/<branch> and its source commit")
	}
	var sequence int
	var digest string
	err := tx.QueryRow(`SELECT sequence, digest FROM policy_snapshots ORDER BY sequence DESC LIMIT 1`).Scan(&sequence, &digest)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("read prior policy snapshot: %w", err)
	}
	if event.Aggregate.Sequence != sequence+1 || digest == effect.Values["digest"] {
		return fmt.Errorf("policy snapshot must advance the sequence and change the effective digest")
	}
	return nil
}

// Validate historical payloads without consulting today's built-in rules:
// replay must keep the policy fid actually used, across later upgrades.
func validatePolicyPayload(raw string, digest any) error {
	var recorded policy.Effective
	decoder := json.NewDecoder(strings.NewReader(raw))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&recorded); err != nil {
		return fmt.Errorf("decode policy snapshot: %w", err)
	}
	if err := recorded.ValidateRecorded(); err != nil {
		return err
	}
	encoded, err := recorded.Canonical()
	if err != nil || !bytes.Equal(encoded, []byte(raw)) {
		return fmt.Errorf("policy snapshot effective policy must use canonical encoding")
	}
	if digest != fmt.Sprintf("sha256:%x", sha256.Sum256(encoded)) {
		return fmt.Errorf("policy snapshot effective digest mismatch")
	}
	return nil
}
