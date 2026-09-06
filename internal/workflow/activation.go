package workflow

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/journal"
)

// Activation is the durable lifecycle record for one kit trust selection.
type Activation struct {
	ID               string   `json:"id"`
	Status           string   `json:"status"`
	KitIdentity      string   `json:"kit_identity"`
	BackendIdentity  string   `json:"backend_identity"`
	Phase            int      `json:"phase"`
	Direction        string   `json:"direction"`
	PointOfNoReturn  bool     `json:"point_of_no_return"`
	ComponentDigests []string `json:"component_digests"`
	CreatedAt        string   `json:"created_at"`
	UpdatedAt        string   `json:"updated_at"`
}

type ActivationPrepareRequest struct {
	KitIdentity     string `json:"kit_identity"`
	BackendIdentity string `json:"backend_identity"`
	Phase           int    `json:"phase"`
	Direction       string `json:"direction"`
}

type ActivationPrepareOptions struct {
	WorkspaceRoot  string
	MigrationsDir  string
	BootstrapActor string
	IdempotencyKey string
	Request        ActivationPrepareRequest
}

type ActivationPrepareResult struct {
	Activation     Activation `json:"activation"`
	IdempotencyKey string     `json:"idempotency_key"`
	Replayed       bool       `json:"replayed"`
	EventID        string
	OccurredAt     time.Time
}

type ActivationComponentRequest struct {
	ActivationID        string          `json:"activation_id"`
	ComponentType       string          `json:"component_type"`
	ComponentID         string          `json:"component_id"`
	ComponentVersion    int             `json:"component_version"`
	ExpectedPriorDigest *string         `json:"expected_prior_digest"`
	Content             json.RawMessage `json:"content"`
}

type ActivationComponentOptions struct {
	WorkspaceRoot  string
	MigrationsDir  string
	BootstrapActor string
	IdempotencyKey string
	Request        ActivationComponentRequest
}

type ActivationComponent struct {
	Type    string `json:"type"`
	ID      string `json:"id"`
	Version int    `json:"version"`
	Digest  string `json:"digest"`
}

type ActivationComponentResult struct {
	Activation     Activation          `json:"activation"`
	Component      ActivationComponent `json:"component"`
	IdempotencyKey string              `json:"idempotency_key"`
	Replayed       bool                `json:"replayed"`
	EventID        string
}

type ActivationCommitRequest struct {
	ActivationID     string   `json:"activation_id"`
	ExpectedStatus   string   `json:"expected_status"`
	ComponentDigests []string `json:"component_digests"`
}

type ActivationAbortRequest struct {
	ActivationID   string `json:"activation_id"`
	ExpectedStatus string `json:"expected_status"`
	Reason         string `json:"reason"`
}

type ActivationFinalizeRequest struct {
	ActivationID   string `json:"activation_id"`
	ExpectedStatus string `json:"expected_status"`
}

type ActivationTransitionOptions struct {
	WorkspaceRoot  string
	MigrationsDir  string
	BootstrapActor string
	IdempotencyKey string
}

type ActorProfile struct {
	ActorID       string   `json:"actor_id"`
	ActorType     string   `json:"actor_type"`
	Producer      string   `json:"producer"`
	ProvenanceRef string   `json:"provenance_ref"`
	SessionRef    *string  `json:"session_ref"`
	Capabilities  []string `json:"capabilities"`
}

type ActivationNotFoundError struct{ ID string }

func (err *ActivationNotFoundError) Error() string {
	return fmt.Sprintf("activation %q was not found", err.ID)
}

type ActivationConflictError struct {
	ActiveID      *string
	NonterminalID *string
}

func (err *ActivationConflictError) Error() string {
	return "another activation is already active or nonterminal"
}

type ActivationInvalidError struct {
	Status           string
	ExpectedStatuses []string
}

func (err *ActivationInvalidError) Error() string {
	return fmt.Sprintf("activation status %q is not permitted", err.Status)
}

type ActivationRecoveryUnavailableError struct{ Status string }

func (err *ActivationRecoveryUnavailableError) Error() string {
	return fmt.Sprintf("activation recovery is unavailable after %s", err.Status)
}

type ActivationOperationError struct{ Operation string }

func (err *ActivationOperationError) Error() string {
	return fmt.Sprintf("activation operation %q is invalid", err.Operation)
}

type TrustedActorUnavailableError struct{ ActorID string }

func (err *TrustedActorUnavailableError) Error() string {
	return fmt.Sprintf("trusted actor %q is unavailable", err.ActorID)
}

func PrepareActivation(options ActivationPrepareOptions) (ActivationPrepareResult, error) {
	if strings.TrimSpace(options.WorkspaceRoot) == "" {
		return ActivationPrepareResult{}, fmt.Errorf("prepare activation: workspace root is required")
	}
	if strings.TrimSpace(options.MigrationsDir) == "" {
		return ActivationPrepareResult{}, fmt.Errorf("prepare activation: migrations dir is required")
	}
	if err := validateActivationPrepare(options.Request); err != nil {
		return ActivationPrepareResult{}, err
	}

	var result ActivationPrepareResult
	err := journal.WithLock(options.WorkspaceRoot, journal.LockOwner{Command: "fid activation prepare"}, func(lock *journal.WriteLock) error {
		if _, err := journal.EnsureProjectionFreshWithLock(lock, options.MigrationsDir); err != nil {
			return err
		}
		conn, err := openWorkflowDB(options.WorkspaceRoot)
		if err != nil {
			return err
		}
		defer closeWorkflowDB(conn)
		project, err := config.ReadProjectConfig(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, "config.toml"))
		if err != nil {
			return fmt.Errorf("prepare activation: read project config: %w", err)
		}
		requestDigest, err := activationRequestDigest("activation.prepare", project.ProjectID, options.BootstrapActor, options.Request)
		if err != nil {
			return err
		}
		idempotency := journal.IdempotencyRequest{OperationScope: "activation.prepare:" + project.ProjectID, IdempotencyKey: options.IdempotencyKey, RequestDigest: requestDigest}
		if outcome, found, err := journal.ReplayIdempotent(conn, idempotency); err != nil {
			return err
		} else if found {
			decoded, err := activationPrepareResultFromJSON(outcome.ResultJSON, outcome.EventID, outcome.Replayed)
			if err != nil {
				return err
			}
			result = decoded
			return nil
		}
		if err := assertActivationAvailable(conn); err != nil {
			return err
		}
		activationID, err := journal.GenerateID("ACT", func(candidate string) (bool, error) {
			var found int
			err := conn.QueryRow(`SELECT 1 FROM kit_activations WHERE id = ?`, candidate).Scan(&found)
			return err == nil, nil
		})
		if err != nil {
			return err
		}
		createdAt := time.Now().UTC()
		activation := Activation{ID: activationID, Status: "prepared", KitIdentity: options.Request.KitIdentity, BackendIdentity: options.Request.BackendIdentity, Phase: 21, Direction: "activate", ComponentDigests: []string{}, CreatedAt: createdAt.Format(time.RFC3339Nano), UpdatedAt: createdAt.Format(time.RFC3339Nano)}
		event := journal.Event{
			EventSchema: journal.EventSchemaV2, EventType: "activation.prepared",
			ActorV2:     journal.EventActor{ActorID: "fid-bootstrap", ActorType: "service", Producer: "fid", ProvenanceRef: "activation:bootstrap"},
			Aggregate:   journal.Aggregate{ID: activationID, Sequence: 1, Type: "activation"},
			CommandName: "activation.prepare", Effects: []journal.Effect{{
				Op: "insert", Table: "kit_activations", PK: map[string]any{"id": activationID}, Values: map[string]any{
					"prior_activation_id": nil, "kit_identity": activation.KitIdentity, "backend_identity": activation.BackendIdentity,
					"phase": activation.Phase, "status": activation.Status, "point_of_no_return": false, "direction": activation.Direction,
					"created_at": activation.CreatedAt, "updated_at": activation.UpdatedAt,
				},
			}},
		}
		resultJSON, err := activationResultJSON(activation, options.IdempotencyKey)
		if err != nil {
			return err
		}
		outcome, err := journal.CommitIdempotent(lock, conn, idempotency, event, func(journal.Event) (journal.IdempotencyMutation, error) {
			return journal.IdempotencyMutation{EntityID: activationID, ResultJSON: resultJSON}, nil
		})
		if err != nil {
			return err
		}
		decoded, err := activationPrepareResultFromJSON(outcome.ResultJSON, outcome.EventID, outcome.Replayed)
		if err != nil {
			return err
		}
		result = decoded
		return nil
	})
	if err != nil {
		return ActivationPrepareResult{}, err
	}
	return result, nil
}

func ActivationStatus(workspaceRoot, migrationsDir, activationID string) (Activation, error) {
	if _, err := journal.EnsureProjectionFresh(workspaceRoot, migrationsDir); err != nil {
		return Activation{}, err
	}
	conn, err := openWorkflowDB(workspaceRoot)
	if err != nil {
		return Activation{}, err
	}
	defer closeWorkflowDB(conn)
	return loadActivation(conn, strings.TrimSpace(activationID))
}

// ResolveTrustedActor returns the Actor profile selected by the finalized activation only.
func ResolveTrustedActor(workspaceRoot, migrationsDir, actorID string) (ActorProfile, error) {
	if _, err := journal.EnsureProjectionFresh(workspaceRoot, migrationsDir); err != nil {
		return ActorProfile{}, err
	}
	conn, err := openWorkflowDB(workspaceRoot)
	if err != nil {
		return ActorProfile{}, err
	}
	defer closeWorkflowDB(conn)
	var activationID string
	if err := conn.QueryRow(`SELECT id FROM kit_activations WHERE status = 'finalized' LIMIT 1`).Scan(&activationID); err == sql.ErrNoRows {
		return ActorProfile{}, &TrustedActorUnavailableError{ActorID: actorID}
	} else if err != nil {
		return ActorProfile{}, fmt.Errorf("resolve trusted actor: read finalized activation: %w", err)
	}
	var actor ActorProfile
	var capabilities string
	err = conn.QueryRow(`SELECT actor_id, actor_type, producer, provenance_ref, session_ref, capabilities_json FROM activation_actor_profiles WHERE activation_id = ? AND actor_id = ?`, activationID, actorID).Scan(&actor.ActorID, &actor.ActorType, &actor.Producer, &actor.ProvenanceRef, &actor.SessionRef, &capabilities)
	if err == sql.ErrNoRows {
		return ActorProfile{}, &TrustedActorUnavailableError{ActorID: actorID}
	}
	if err != nil {
		return ActorProfile{}, fmt.Errorf("resolve trusted actor %q: %w", actorID, err)
	}
	if err := json.Unmarshal([]byte(capabilities), &actor.Capabilities); err != nil {
		return ActorProfile{}, fmt.Errorf("resolve trusted actor %q capabilities: %w", actorID, err)
	}
	return actor, nil
}

func StageActivationComponent(options ActivationComponentOptions) (ActivationComponentResult, error) {
	if strings.TrimSpace(options.WorkspaceRoot) == "" || strings.TrimSpace(options.MigrationsDir) == "" {
		return ActivationComponentResult{}, fmt.Errorf("stage activation component: workspace root and migrations directory are required")
	}
	component, profiles, definitions, err := validateActivationComponent(options.Request)
	if err != nil {
		return ActivationComponentResult{}, err
	}
	var result ActivationComponentResult
	err = journal.WithLock(options.WorkspaceRoot, journal.LockOwner{Command: "fid activation component"}, func(lock *journal.WriteLock) error {
		if _, err := journal.EnsureProjectionFreshWithLock(lock, options.MigrationsDir); err != nil {
			return err
		}
		conn, err := openWorkflowDB(options.WorkspaceRoot)
		if err != nil {
			return err
		}
		defer closeWorkflowDB(conn)
		project, err := config.ReadProjectConfig(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, "config.toml"))
		if err != nil {
			return fmt.Errorf("stage activation component: read project config: %w", err)
		}
		requestDigest, err := activationRequestDigest("activation.component", project.ProjectID, options.BootstrapActor, options.Request)
		if err != nil {
			return err
		}
		idempotency := journal.IdempotencyRequest{OperationScope: "activation.component:" + project.ProjectID, IdempotencyKey: options.IdempotencyKey, RequestDigest: requestDigest}
		if outcome, found, err := journal.ReplayIdempotent(conn, idempotency); err != nil {
			return err
		} else if found {
			decoded, err := activationComponentResultFromJSON(outcome.ResultJSON, outcome.EventID, outcome.Replayed)
			if err != nil {
				return err
			}
			result = decoded
			return nil
		}
		activation, err := loadActivation(conn, options.Request.ActivationID)
		if err != nil {
			return err
		}
		if activation.Status != "prepared" && activation.Status != "staged" {
			return &ActivationInvalidError{Status: activation.Status, ExpectedStatuses: []string{"prepared", "staged"}}
		}
		nextSequence, err := activationComponentSequence(conn, activation.ID)
		if err != nil {
			return err
		}
		updatedAt := time.Now().UTC().Format(time.RFC3339Nano)
		effects := activationComponentEffects(activation, options.Request, component, profiles, definitions, updatedAt)
		event := journal.Event{
			EventSchema: journal.EventSchemaV2, EventType: "activation.component-staged",
			ActorV2:     journal.EventActor{ActorID: "fid-bootstrap", ActorType: "service", Producer: "fid", ProvenanceRef: "activation:bootstrap"},
			Aggregate:   journal.Aggregate{ID: activation.ID, Sequence: nextSequence, Type: "activation"},
			CommandName: "activation.component", Effects: effects,
		}
		staged := activation
		staged.Status, staged.UpdatedAt = "staged", updatedAt
		staged.ComponentDigests = append(staged.ComponentDigests, component.Digest)
		resultJSON, err := activationComponentResultJSON(staged, component, options.IdempotencyKey)
		if err != nil {
			return err
		}
		outcome, err := journal.CommitIdempotent(lock, conn, idempotency, event, func(journal.Event) (journal.IdempotencyMutation, error) {
			return journal.IdempotencyMutation{EntityID: activation.ID, ResultJSON: resultJSON}, nil
		})
		if err != nil {
			return err
		}
		decoded, err := activationComponentResultFromJSON(outcome.ResultJSON, outcome.EventID, outcome.Replayed)
		if err != nil {
			return err
		}
		result = decoded
		return nil
	})
	if err != nil {
		return ActivationComponentResult{}, err
	}
	return result, nil
}

func CommitActivation(options ActivationTransitionOptions, request ActivationCommitRequest) (ActivationPrepareResult, error) {
	if request.ExpectedStatus != "staged" || len(request.ComponentDigests) != 2 || request.ComponentDigests[0] >= request.ComponentDigests[1] {
		return ActivationPrepareResult{}, fmt.Errorf("activation commit input is invalid")
	}
	return transitionActivation(options, request.ActivationID, "activation.commit", "committed", func(conn *sql.DB, activation Activation) ([]journal.Effect, error) {
		if activation.Status != request.ExpectedStatus {
			return nil, &ActivationInvalidError{Status: activation.Status, ExpectedStatuses: []string{"staged"}}
		}
		rows, err := conn.Query(`SELECT component_type, component_id, component_version, content_digest FROM activation_components WHERE activation_id = ? AND status = 'staged' ORDER BY content_digest`, activation.ID)
		if err != nil {
			return nil, err
		}
		defer func() { _ = rows.Close() }()
		var digests []string
		var components []activationComponentIdentity
		for rows.Next() {
			var component activationComponentIdentity
			var digest string
			if err := rows.Scan(&component.Type, &component.ID, &component.Version, &digest); err != nil {
				return nil, err
			}
			digests = append(digests, digest)
			components = append(components, component)
		}
		if err := rows.Err(); err != nil {
			return nil, err
		}
		if len(digests) != 2 || digests[0] != request.ComponentDigests[0] || digests[1] != request.ComponentDigests[1] {
			return nil, &ActivationInvalidError{Status: activation.Status, ExpectedStatuses: []string{"staged"}}
		}
		effects := []journal.Effect{
			{Op: "update", Table: "kit_activations", PK: map[string]any{"id": activation.ID}, Values: map[string]any{"status": "committed", "point_of_no_return": true, "updated_at": time.Now().UTC().Format(time.RFC3339Nano)}},
		}
		return append(effects, activationComponentStatusEffects(activation.ID, components, "committed")...), nil
	}, request)
}

func AbortActivation(options ActivationTransitionOptions, request ActivationAbortRequest) (ActivationPrepareResult, error) {
	if (request.ExpectedStatus != "prepared" && request.ExpectedStatus != "staged") || strings.TrimSpace(request.Reason) == "" || len(request.Reason) > 256 {
		return ActivationPrepareResult{}, fmt.Errorf("activation abort input is invalid")
	}
	return transitionActivation(options, request.ActivationID, "activation.abort", "aborted", func(conn *sql.DB, activation Activation) ([]journal.Effect, error) {
		if activation.Status == "committed" || activation.Status == "finalized" {
			return nil, &ActivationRecoveryUnavailableError{Status: activation.Status}
		}
		if activation.Status != request.ExpectedStatus {
			return nil, &ActivationInvalidError{Status: activation.Status, ExpectedStatuses: []string{"prepared", "staged"}}
		}
		effects := []journal.Effect{
			{Op: "update", Table: "kit_activations", PK: map[string]any{"id": activation.ID}, Values: map[string]any{"status": "aborted", "updated_at": time.Now().UTC().Format(time.RFC3339Nano)}},
		}
		components, err := activationComponentIdentities(conn, activation.ID)
		if err != nil {
			return nil, err
		}
		return append(effects, activationComponentStatusEffects(activation.ID, components, "aborted")...), nil
	}, request)
}

func FinalizeActivation(options ActivationTransitionOptions, request ActivationFinalizeRequest) (ActivationPrepareResult, error) {
	if request.ExpectedStatus != "committed" || strings.TrimSpace(request.ActivationID) == "" {
		return ActivationPrepareResult{}, fmt.Errorf("activation finalize input is invalid")
	}
	return transitionActivation(options, request.ActivationID, "activation.finalize", "finalized", func(conn *sql.DB, activation Activation) ([]journal.Effect, error) {
		if activation.Status != "committed" {
			return nil, &ActivationInvalidError{Status: activation.Status, ExpectedStatuses: []string{"committed"}}
		}
		components, err := activationComponentIdentities(conn, activation.ID)
		if err != nil {
			return nil, err
		}
		effects := []journal.Effect{{Op: "update", Table: "kit_activations", PK: map[string]any{"id": activation.ID}, Values: map[string]any{"status": "finalized", "updated_at": time.Now().UTC().Format(time.RFC3339Nano)}}}
		return append(effects, activationComponentStatusEffects(activation.ID, components, "finalized")...), nil
	}, request)
}

type activationComponentIdentity struct {
	Type, ID string
	Version  int
}

func activationComponentIdentities(conn *sql.DB, activationID string) ([]activationComponentIdentity, error) {
	rows, err := conn.Query(`SELECT component_type, component_id, component_version FROM activation_components WHERE activation_id = ?`, activationID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	var components []activationComponentIdentity
	for rows.Next() {
		var component activationComponentIdentity
		if err := rows.Scan(&component.Type, &component.ID, &component.Version); err != nil {
			return nil, err
		}
		components = append(components, component)
	}
	return components, rows.Err()
}

func activationComponentStatusEffects(activationID string, components []activationComponentIdentity, status string) []journal.Effect {
	effects := make([]journal.Effect, 0, len(components))
	for _, component := range components {
		effects = append(effects, journal.Effect{Op: "update", Table: "activation_components", PK: map[string]any{"activation_id": activationID, "component_type": component.Type, "component_id": component.ID, "component_version": component.Version}, Values: map[string]any{"status": status}})
	}
	return effects
}

func transitionActivation(options ActivationTransitionOptions, activationID, operation, status string, build func(*sql.DB, Activation) ([]journal.Effect, error), request any) (ActivationPrepareResult, error) {
	if strings.TrimSpace(options.WorkspaceRoot) == "" || strings.TrimSpace(options.MigrationsDir) == "" {
		return ActivationPrepareResult{}, fmt.Errorf("%s: workspace root and migrations directory are required", operation)
	}
	var result ActivationPrepareResult
	err := journal.WithLock(options.WorkspaceRoot, journal.LockOwner{Command: "fid " + strings.ReplaceAll(operation, ".", " ")}, func(lock *journal.WriteLock) error {
		if _, err := journal.EnsureProjectionFreshWithLock(lock, options.MigrationsDir); err != nil {
			return err
		}
		conn, err := openWorkflowDB(options.WorkspaceRoot)
		if err != nil {
			return err
		}
		defer closeWorkflowDB(conn)
		project, err := config.ReadProjectConfig(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, "config.toml"))
		if err != nil {
			return err
		}
		requestDigest, err := activationRequestDigest(operation, project.ProjectID, options.BootstrapActor, request)
		if err != nil {
			return err
		}
		idempotency := journal.IdempotencyRequest{OperationScope: operation + ":" + project.ProjectID, IdempotencyKey: options.IdempotencyKey, RequestDigest: requestDigest}
		if outcome, found, err := journal.ReplayIdempotent(conn, idempotency); err != nil {
			return err
		} else if found {
			decoded, err := activationPrepareResultFromJSON(outcome.ResultJSON, outcome.EventID, outcome.Replayed)
			if err != nil {
				return err
			}
			result = decoded
			return nil
		}
		activation, err := loadActivation(conn, strings.TrimSpace(activationID))
		if err != nil {
			return err
		}
		effects, err := build(conn, activation)
		if err != nil {
			return err
		}
		nextSequence, err := activationTransitionSequence(conn, activation)
		if err != nil {
			return err
		}
		var eventType string
		switch operation {
		case "activation.commit":
			eventType = "activation.committed"
		case "activation.abort":
			eventType = "activation.aborted"
		case "activation.finalize":
			eventType = "activation.finalized"
		default:
			return &ActivationOperationError{Operation: operation}
		}
		event := journal.Event{EventSchema: journal.EventSchemaV2, EventType: eventType, ActorV2: journal.EventActor{ActorID: "fid-bootstrap", ActorType: "service", Producer: "fid", ProvenanceRef: "activation:bootstrap"}, Aggregate: journal.Aggregate{ID: activation.ID, Sequence: nextSequence, Type: "activation"}, CommandName: operation, Effects: effects}
		updated := activation
		updated.Status = status
		updated.PointOfNoReturn = status == "committed" || activation.PointOfNoReturn
		updated.UpdatedAt = time.Now().UTC().Format(time.RFC3339Nano)
		resultJSON, err := activationResultJSON(updated, options.IdempotencyKey)
		if err != nil {
			return err
		}
		outcome, err := journal.CommitIdempotent(lock, conn, idempotency, event, func(journal.Event) (journal.IdempotencyMutation, error) {
			return journal.IdempotencyMutation{EntityID: activation.ID, ResultJSON: resultJSON}, nil
		})
		if err != nil {
			return err
		}
		decoded, err := activationPrepareResultFromJSON(outcome.ResultJSON, outcome.EventID, outcome.Replayed)
		if err != nil {
			return err
		}
		result = decoded
		return nil
	})
	if err != nil {
		return ActivationPrepareResult{}, err
	}
	return result, nil
}

func activationTransitionSequence(conn *sql.DB, activation Activation) (int, error) {
	sequence, err := activationComponentSequence(conn, activation.ID)
	if err != nil {
		return 0, err
	}
	if activation.Status == "committed" {
		return sequence + 1, nil
	}
	return sequence, nil
}

func validateActivationPrepare(request ActivationPrepareRequest) error {
	if strings.TrimSpace(request.KitIdentity) == "" || strings.TrimSpace(request.BackendIdentity) == "" || request.Phase != 21 || request.Direction != "activate" {
		return fmt.Errorf("activation prepare input is invalid")
	}
	return nil
}

func validateActivationComponent(request ActivationComponentRequest) (ActivationComponent, []ActorProfile, []map[string]any, error) {
	if strings.TrimSpace(request.ActivationID) == "" || strings.TrimSpace(request.ComponentID) == "" || request.ComponentVersion < 1 {
		return ActivationComponent{}, nil, nil, fmt.Errorf("activation component identity is invalid")
	}
	payload, err := journal.MarshalCanonicalJSON(json.RawMessage(request.Content))
	if err != nil {
		return ActivationComponent{}, nil, nil, fmt.Errorf("activation component content: %w", err)
	}
	sum := sha256.Sum256(payload)
	component := ActivationComponent{Type: request.ComponentType, ID: request.ComponentID, Version: request.ComponentVersion, Digest: "sha256:" + hex.EncodeToString(sum[:])}
	switch request.ComponentType {
	case "authority-profile":
		var profiles []ActorProfile
		decoder := json.NewDecoder(strings.NewReader(string(request.Content)))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&profiles); err != nil {
			return ActivationComponent{}, nil, nil, fmt.Errorf("decode authority profile: %w", err)
		}
		if err := validateActorProfiles(profiles); err != nil {
			return ActivationComponent{}, nil, nil, err
		}
		return component, profiles, nil, nil
	case "gate-definitions":
		var definitions []map[string]any
		if err := json.Unmarshal(request.Content, &definitions); err != nil || len(definitions) == 0 {
			return ActivationComponent{}, nil, nil, fmt.Errorf("gate definitions component must contain a nonempty array")
		}
		return component, nil, definitions, nil
	default:
		return ActivationComponent{}, nil, nil, fmt.Errorf("activation component type %q is invalid", request.ComponentType)
	}
}

func validateActorProfiles(profiles []ActorProfile) error {
	if len(profiles) == 0 {
		return fmt.Errorf("authority profile must contain at least one actor")
	}
	previous := ""
	for _, profile := range profiles {
		if strings.TrimSpace(profile.ActorID) == "" || strings.TrimSpace(profile.Producer) == "" || strings.TrimSpace(profile.ProvenanceRef) == "" || len(profile.Capabilities) == 0 {
			return fmt.Errorf("authority profile actor is incomplete")
		}
		if profile.ActorType != "agent" && profile.ActorType != "human" && profile.ActorType != "service" {
			return fmt.Errorf("authority profile actor type %q is invalid", profile.ActorType)
		}
		if previous >= profile.ActorID {
			return fmt.Errorf("authority profile actors must be sorted and unique")
		}
		previous = profile.ActorID
		if !sortedUniqueCapabilities(profile.Capabilities) {
			return fmt.Errorf("authority profile capabilities must be sorted, unique, and accepted")
		}
	}
	return nil
}

func sortedUniqueCapabilities(capabilities []string) bool {
	allowed := map[string]bool{"artifact.verify": true, "attention.resolve": true, "claim.acquire": true, "claim.expire": true, "claim.release": true, "claim.renew": true, "envelope.activate": true, "envelope.complete": true, "envelope.create": true, "evidence.link": true, "gate.register": true, "gate.run": true, "issue.transition": true, "run.checkpoint": true, "run.finish": true, "run.force-recover": true, "run.observe": true, "run.observe-budget": true, "run.recover": true, "run.start": true, "stop.raise": true, "stop.resolve": true}
	previous := ""
	for _, capability := range capabilities {
		if !allowed[capability] || previous >= capability {
			return false
		}
		previous = capability
	}
	return true
}

func activationComponentSequence(conn *sql.DB, activationID string) (int, error) {
	var count int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM activation_components WHERE activation_id = ?`, activationID).Scan(&count); err != nil {
		return 0, fmt.Errorf("count activation components: %w", err)
	}
	return count + 2, nil
}

func activationComponentEffects(activation Activation, request ActivationComponentRequest, component ActivationComponent, profiles []ActorProfile, definitions []map[string]any, updatedAt string) []journal.Effect {
	effects := []journal.Effect{
		{Op: "update", Table: "kit_activations", PK: map[string]any{"id": activation.ID}, Values: map[string]any{"status": "staged", "updated_at": updatedAt}},
		{Op: "insert", Table: "activation_components", PK: map[string]any{"activation_id": activation.ID, "component_type": component.Type, "component_id": component.ID, "component_version": component.Version}, Values: map[string]any{"content_digest": component.Digest, "expected_prior_digest": request.ExpectedPriorDigest, "status": "staged", "evidence_json": string(request.Content), "error_code": nil}},
	}
	for _, profile := range profiles {
		capabilities, _ := journal.MarshalCanonicalJSON(profile.Capabilities)
		effects = append(effects, journal.Effect{Op: "insert", Table: "activation_actor_profiles", PK: map[string]any{"activation_id": activation.ID, "actor_id": profile.ActorID}, Values: map[string]any{"actor_type": profile.ActorType, "producer": profile.Producer, "provenance_ref": profile.ProvenanceRef, "session_ref": profile.SessionRef, "capabilities_json": string(capabilities)}})
	}
	for _, definition := range definitions {
		definitionID, _ := definition["definition_id"].(string)
		definitionVersion, _ := definition["definition_version"].(float64)
		payload, _ := journal.MarshalCanonicalJSON(definition)
		sum := sha256.Sum256(payload)
		effects = append(effects, journal.Effect{Op: "insert", Table: "activation_gate_definitions", PK: map[string]any{"activation_id": activation.ID, "definition_id": definitionID, "definition_version": int(definitionVersion)}, Values: map[string]any{"definition_digest": "sha256:" + hex.EncodeToString(sum[:]), "definition_json": string(payload)}})
	}
	return effects
}

func assertActivationAvailable(conn *sql.DB) error {
	var active sql.NullString
	if err := conn.QueryRow(`SELECT id FROM kit_activations WHERE status = 'finalized' LIMIT 1`).Scan(&active); err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("read active activation: %w", err)
	}
	var nonterminal sql.NullString
	if err := conn.QueryRow(`SELECT id FROM kit_activations WHERE status IN ('prepared', 'staged', 'committed') LIMIT 1`).Scan(&nonterminal); err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("read nonterminal activation: %w", err)
	}
	if active.Valid || nonterminal.Valid {
		return &ActivationConflictError{ActiveID: optionalString(active), NonterminalID: optionalString(nonterminal)}
	}
	return nil
}

func loadActivation(conn *sql.DB, id string) (Activation, error) {
	var activation Activation
	var pointOfNoReturn int
	err := conn.QueryRow(`
		SELECT id, status, kit_identity, backend_identity, phase, direction, point_of_no_return, created_at, updated_at
		FROM kit_activations WHERE id = ?`, id).Scan(
		&activation.ID, &activation.Status, &activation.KitIdentity, &activation.BackendIdentity, &activation.Phase,
		&activation.Direction, &pointOfNoReturn, &activation.CreatedAt, &activation.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return Activation{}, &ActivationNotFoundError{ID: id}
	}
	if err != nil {
		return Activation{}, fmt.Errorf("load activation %q: %w", id, err)
	}
	activation.PointOfNoReturn = pointOfNoReturn != 0
	rows, err := conn.Query(`SELECT content_digest FROM activation_components WHERE activation_id = ? ORDER BY content_digest`, id)
	if err != nil {
		return Activation{}, fmt.Errorf("load activation components: %w", err)
	}
	defer func() { _ = rows.Close() }()
	activation.ComponentDigests = make([]string, 0)
	for rows.Next() {
		var digest string
		if err := rows.Scan(&digest); err != nil {
			return Activation{}, err
		}
		activation.ComponentDigests = append(activation.ComponentDigests, digest)
	}
	return activation, rows.Err()
}

func activationRequestDigest(operation, projectID, bootstrapActor string, request any) (string, error) {
	payload, err := journal.MarshalCanonicalJSON(map[string]any{"bootstrap_actor": bootstrapActor, "operation": operation, "project_id": projectID, "request": request})
	if err != nil {
		return "", fmt.Errorf("canonical activation request: %w", err)
	}
	sum := sha256.Sum256(payload)
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}

func activationResultJSON(activation Activation, key string) (string, error) {
	payload, err := journal.MarshalCanonicalJSON(map[string]any{"activation": activation, "idempotency_key": key, "replayed": false})
	if err != nil {
		return "", err
	}
	return string(payload), nil
}

func activationComponentResultJSON(activation Activation, component ActivationComponent, key string) (string, error) {
	payload, err := journal.MarshalCanonicalJSON(map[string]any{"activation": activation, "component": component, "idempotency_key": key, "replayed": false})
	if err != nil {
		return "", err
	}
	return string(payload), nil
}

func activationPrepareResultFromJSON(raw, eventID string, replayed bool) (ActivationPrepareResult, error) {
	var value struct {
		Activation     Activation `json:"activation"`
		IdempotencyKey string     `json:"idempotency_key"`
	}
	if err := json.Unmarshal([]byte(raw), &value); err != nil {
		return ActivationPrepareResult{}, fmt.Errorf("decode stored activation result: %w", err)
	}
	return ActivationPrepareResult{Activation: value.Activation, IdempotencyKey: value.IdempotencyKey, Replayed: replayed, EventID: eventID}, nil
}

func activationComponentResultFromJSON(raw, eventID string, replayed bool) (ActivationComponentResult, error) {
	var value struct {
		Activation     Activation          `json:"activation"`
		Component      ActivationComponent `json:"component"`
		IdempotencyKey string              `json:"idempotency_key"`
	}
	if err := json.Unmarshal([]byte(raw), &value); err != nil {
		return ActivationComponentResult{}, fmt.Errorf("decode stored activation component result: %w", err)
	}
	return ActivationComponentResult{Activation: value.Activation, Component: value.Component, IdempotencyKey: value.IdempotencyKey, Replayed: replayed, EventID: eventID}, nil
}

func optionalString(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}
	result := value.String
	return &result
}
