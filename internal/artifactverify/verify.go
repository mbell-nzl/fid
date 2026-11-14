// Package artifactverify observes retained Artifact availability independently of checks.
package artifactverify

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/artifactstore"
	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/policy"
)

type Request struct {
	ArtifactID string `json:"artifact_id"`
}

type Options struct {
	WorkspaceRoot  string
	MigrationsDir  string
	ActorID        string
	IdempotencyKey string
	Request        Request
}

type Result struct {
	Artifact        Artifact `json:"artifact"`
	Changed         bool     `json:"changed"`
	IdempotencyKey  string   `json:"idempotency_key"`
	Replayed        bool     `json:"replayed"`
	EventID         string
	EventOccurredAt string
}

// Artifact is content-addressed material retained by fid.
type Artifact struct {
	ID             string  `json:"id"`
	Digest         string  `json:"digest"`
	MediaType      string  `json:"media_type"`
	SizeBytes      int64   `json:"size_bytes"`
	RetentionClass string  `json:"retention_class"`
	Availability   string  `json:"availability"`
	LocationPolicy string  `json:"location_policy"`
	CreatedAt      string  `json:"created_at"`
	ExpiresAt      *string `json:"expires_at"`
	DeclaredPath   string  `json:"-"`
}

// DeniedActorError refuses actor types outside the work set.
type DeniedActorError struct{ Actor string }

func (err *DeniedActorError) Error() string {
	return fmt.Sprintf("actor %s can't verify Artifacts", err.Actor)
}

var actorPattern = regexp.MustCompile(policy.ActorPattern)

// Verify uses the lean core's built-in authority. Check results remain immutable;
// only the retained Artifact's current availability is observed.
func Verify(options Options) (Result, error) {
	if !actorPattern.MatchString(options.ActorID) {
		return Result{}, fmt.Errorf("verify Artifact: invalid actor %q", options.ActorID)
	}
	if !policy.Allows(options.ActorID, policy.WorkSet) {
		return Result{}, &DeniedActorError{Actor: options.ActorID}
	}
	kind, id, _ := strings.Cut(options.ActorID, ":")
	return VerifyWithHooks(options, Hooks{Authorize: func(*sql.DB) (journal.EventActor, error) {
		return journal.EventActor{ActorID: id, ActorType: kind, Producer: "fid-standalone", ProvenanceRef: "project-local-attribution"}, nil
	}})
}

// Hooks isolates the legacy authority and readiness effects until the kernel's
// removal in RESET 1.5c2. Lean callers use Verify.
type Hooks struct {
	Authorize func(*sql.DB) (journal.EventActor, error)
	Restored  func(*sql.DB, string, journal.EventActor) ([]journal.Effect, error)
	Prepare   func(journal.Event)
}

func eventOccurredAt(root, id string) (string, error) {
	found, err := journal.FindEvent(filepath.Join(config.LedgerDir(root), config.JournalDirName), id)
	if err != nil {
		return "", fmt.Errorf("verify Artifact: read Event: %w", err)
	}
	return found.Event.OccurredAt.UTC().Format(time.RFC3339Nano), nil
}

func VerifyWithHooks(options Options, hooks Hooks) (Result, error) {
	if strings.TrimSpace(options.WorkspaceRoot) == "" || strings.TrimSpace(options.MigrationsDir) == "" || strings.TrimSpace(options.ActorID) == "" || strings.TrimSpace(options.Request.ArtifactID) == "" {
		return Result{}, fmt.Errorf("verify Artifact: workspace root, migrations directory, actor, and Artifact are required")
	}
	if hooks.Authorize == nil {
		return Result{}, fmt.Errorf("verify Artifact: authority is required")
	}
	var result Result
	err := journal.WithLock(options.WorkspaceRoot, journal.LockOwner{Command: "fid artifact verify"}, func(lock *journal.WriteLock) error {
		if _, err := journal.EnsureProjectionFreshWithLock(lock, options.MigrationsDir); err != nil {
			return err
		}
		conn, err := fiddb.Open(filepath.Join(config.LedgerDir(options.WorkspaceRoot), config.DatabaseFileName))
		if err != nil {
			return err
		}
		defer func() { _ = conn.Close() }()
		project, err := config.ReadProjectConfig(filepath.Join(config.LedgerDir(options.WorkspaceRoot), config.ConfigFileName))
		if err != nil {
			return fmt.Errorf("verify Artifact: read project config: %w", err)
		}
		digest, fields, err := journal.ComputeRequestDigest(map[string]any{"actor_id": options.ActorID, "artifact_id": options.Request.ArtifactID, "project_id": project.ProjectID, "request": options.Request}, options.ActorID, nil)
		if err != nil {
			return err
		}
		idempotency := journal.IdempotencyRequest{OperationScope: "artifact.verify:" + project.ProjectID, IdempotencyKey: options.IdempotencyKey, RequestDigest: digest, RequestFieldDigests: fields}
		if outcome, found, err := journal.ReplayIdempotent(conn, idempotency); err != nil {
			return err
		} else if found {
			if err := json.Unmarshal([]byte(outcome.ResultJSON), &result); err != nil {
				return fmt.Errorf("verify Artifact: decode replay: %w", err)
			}
			result.EventID = outcome.EventID
			result.IdempotencyKey = options.IdempotencyKey
			result.Replayed = true
			return nil
		}
		actor, err := hooks.Authorize(conn)
		if err != nil {
			return err
		}
		artifact, found, err := loadArtifactByID(conn, options.Request.ArtifactID)
		if err != nil {
			return err
		}
		if !found {
			return fmt.Errorf("verify Artifact: Artifact %q was not found", options.Request.ArtifactID)
		}

		storeRoot, err := artifactstore.Root(config.LedgerDir(options.WorkspaceRoot), project.ArtifactRoot)
		if err != nil {
			return err
		}
		available := artifactstore.Verify(filepath.Join(storeRoot, strings.TrimPrefix(artifact.Digest, "sha256:")), artifact.Digest) == nil
		observed := "missing"
		if available {
			observed = "available"
		}
		priorAvailability := artifact.Availability
		changed := priorAvailability != observed
		artifact.Availability = observed
		sequence, err := journal.NextAggregateSequence(filepath.Join(config.LedgerDir(options.WorkspaceRoot), config.JournalDirName), journal.Aggregate{ID: artifact.ID, Type: "artifact"})
		if err != nil {
			return err
		}
		effects := []journal.Effect{}
		if changed {
			effects = append(effects, journal.Effect{Op: "update", Table: "referenced_artifacts", PK: map[string]any{"id": artifact.ID}, Values: map[string]any{"availability": observed}, Preconditions: map[string]any{"availability": priorAvailability}})
		}
		if changed && observed == "available" && hooks.Restored != nil {
			extra, err := hooks.Restored(conn, artifact.ID, actor)
			if err != nil {
				return err
			}
			effects = append(effects, extra...)
		}
		event := journal.Event{EventSchema: journal.EventSchemaV2, EventType: "artifact.availability-observed", ActorV2: actor, Aggregate: journal.Aggregate{ID: artifact.ID, Sequence: sequence, Type: "artifact"}, CommandName: "artifact.verify", Effects: effects}
		result = Result{Artifact: artifact, Changed: changed, IdempotencyKey: options.IdempotencyKey, Replayed: false}
		resultJSON, err := journal.MarshalCanonicalJSON(result)
		if err != nil {
			return err
		}
		outcome, err := journal.CommitIdempotent(lock, conn, idempotency, event, func(prepared journal.Event) (journal.IdempotencyMutation, error) {
			if hooks.Prepare != nil {
				hooks.Prepare(prepared)
			}
			return journal.IdempotencyMutation{EntityID: artifact.ID, ResultJSON: string(resultJSON)}, nil
		})
		if err != nil {
			return err
		}
		result.EventID = outcome.EventID
		return nil
	})
	if err != nil {
		return Result{}, err
	}
	if result.EventID != "" {
		occurredAt, err := eventOccurredAt(options.WorkspaceRoot, result.EventID)
		if err != nil {
			return Result{}, err
		}
		result.EventOccurredAt = occurredAt
	}
	return result, nil
}

func loadArtifactByID(conn *sql.DB, id string) (Artifact, bool, error) {
	var artifact Artifact
	err := conn.QueryRow(`SELECT id, digest, media_type, size_bytes, retention_class, availability, location_policy, created_at, expires_at FROM referenced_artifacts WHERE id = ?`, id).Scan(
		&artifact.ID, &artifact.Digest, &artifact.MediaType, &artifact.SizeBytes, &artifact.RetentionClass,
		&artifact.Availability, &artifact.LocationPolicy, &artifact.CreatedAt, &artifact.ExpiresAt,
	)
	if err == sql.ErrNoRows {
		return Artifact{}, false, nil
	}
	if err != nil {
		return Artifact{}, false, fmt.Errorf("load Artifact: %w", err)
	}
	return artifact, true, nil
}
