package evidence

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
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
)

type RegisterRequest struct {
	DefinitionID      string     `json:"definition_id"`
	DefinitionVersion int        `json:"definition_version"`
	WorktreeID        string     `json:"worktree_id"`
	Scope             GateScope  `json:"scope"`
	Repository        Repository `json:"repository"`
}

type GateScope struct {
	IssueID string `json:"issue_id"`
}

type Repository struct {
	BaseOID       *string `json:"base_oid"`
	CaptureSource string  `json:"capture_source"`
	Dirty         bool    `json:"dirty"`
	DirtyDigest   *string `json:"dirty_digest"`
	HeadOID       string  `json:"head_oid"`
	ProjectID     string  `json:"project_id"`
	TreeOID       string  `json:"tree_oid"`
	WorktreeID    string  `json:"worktree_id"`
}

type Registration struct {
	ID                      string         `json:"id"`
	Version                 int            `json:"version"`
	ActivationID            string         `json:"activation_id"`
	DefinitionID            string         `json:"definition_id"`
	DefinitionVersion       int            `json:"definition_version"`
	DefinitionDigest        string         `json:"definition_digest"`
	BindingDigest           string         `json:"binding_digest"`
	ProjectID               string         `json:"project_id"`
	WorktreeID              string         `json:"worktree_id"`
	Scope                   GateScope      `json:"scope"`
	AuthorizationProvenance string         `json:"authorization_provenance"`
	Policy                  GateDefinition `json:"policy"`
	CreatedAt               string         `json:"created_at"`
}

type RegisterOptions struct {
	WorkspaceRoot  string
	MigrationsDir  string
	ActorID        string
	IdempotencyKey string
	Request        RegisterRequest
}

type RegisterResult struct {
	Registration   Registration `json:"registration"`
	IdempotencyKey string       `json:"idempotency_key"`
	Replayed       bool         `json:"replayed"`
	EventID        string
}

type ActiveActivationUnavailableError struct{}

func (*ActiveActivationUnavailableError) Error() string { return "a finalized activation is required" }

type GateDefinitionNotFoundError struct {
	DefinitionID string
	Version      int
}

func (err *GateDefinitionNotFoundError) Error() string {
	return fmt.Sprintf("activated gate definition %q version %d was not found", err.DefinitionID, err.Version)
}

type RegistrationNotFoundError struct {
	ID      string
	Version int
}

func (err *RegistrationNotFoundError) Error() string {
	return fmt.Sprintf("gate registration %q version %d was not found", err.ID, err.Version)
}

type TrustedActorUnavailableError struct{ ActorID string }

func (err *TrustedActorUnavailableError) Error() string {
	return fmt.Sprintf("trusted actor %q is unavailable", err.ActorID)
}

type AuthorityDeniedError struct{ ActorID string }

func (err *AuthorityDeniedError) Error() string {
	return fmt.Sprintf("trusted actor %q lacks gate.register", err.ActorID)
}

func Register(options RegisterOptions) (RegisterResult, error) {
	if strings.TrimSpace(options.WorkspaceRoot) == "" || strings.TrimSpace(options.MigrationsDir) == "" {
		return RegisterResult{}, fmt.Errorf("register gate: workspace root and migrations directory are required")
	}
	if err := validateRegisterRequest(options.Request); err != nil {
		return RegisterResult{}, err
	}

	var result RegisterResult
	err := journal.WithLock(options.WorkspaceRoot, journal.LockOwner{Command: "fid gate register"}, func(lock *journal.WriteLock) error {
		if _, err := journal.EnsureProjectionFreshWithLock(lock, options.MigrationsDir); err != nil {
			return err
		}
		conn, err := openEvidenceDB(options.WorkspaceRoot)
		if err != nil {
			return err
		}
		defer func() { _ = conn.Close() }()

		project, err := config.ReadProjectConfig(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, "config.toml"))
		if err != nil {
			return fmt.Errorf("register gate: read project config: %w", err)
		}
		activationID, actor, err := loadActiveActor(conn, options.ActorID)
		if err != nil {
			return err
		}
		if !hasCapability(actor.Capabilities, "gate.register") {
			return &AuthorityDeniedError{ActorID: options.ActorID}
		}
		if err := validateRepositoryBinding(options.Request, project.ProjectID); err != nil {
			return err
		}
		definition, digest, err := loadActivatedDefinition(conn, activationID, options.Request.DefinitionID, options.Request.DefinitionVersion)
		if err != nil {
			return err
		}
		if err := ValidateGateDefinition(options.WorkspaceRoot, definition); err != nil {
			return err
		}
		bindingDigest, err := gateBindingDigest(options.Request)
		if err != nil {
			return err
		}

		if existing, found, err := loadRegistrationByBinding(conn, digest, bindingDigest); err != nil {
			return err
		} else if found {
			result = RegisterResult{Registration: existing, IdempotencyKey: options.IdempotencyKey, Replayed: true, EventID: existingCreatedEventID(conn, existing)}
			return nil
		}

		registrationID, version, err := nextRegistrationIdentity(conn, options.Request.DefinitionID)
		if err != nil {
			return err
		}
		createdAt := time.Now().UTC().Format(time.RFC3339Nano)
		registration := Registration{ID: registrationID, Version: version, ActivationID: activationID, DefinitionID: definition.DefinitionID, DefinitionVersion: definition.DefinitionVersion, DefinitionDigest: digest, BindingDigest: bindingDigest, ProjectID: project.ProjectID, WorktreeID: options.Request.WorktreeID, Scope: options.Request.Scope, AuthorizationProvenance: actor.ProvenanceRef, Policy: definition, CreatedAt: createdAt}
		policyJSON, err := journal.MarshalCanonicalJSON(definition)
		if err != nil {
			return err
		}
		scopeJSON, err := journal.MarshalCanonicalJSON(options.Request.Scope)
		if err != nil {
			return err
		}
		requestDigest, fields, err := journal.ComputeRequestDigest(map[string]any{"actor_id": options.ActorID, "project_id": project.ProjectID, "request": options.Request}, options.ActorID, nil)
		if err != nil {
			return err
		}
		resultJSON, err := resultJSON(registration, options.IdempotencyKey)
		if err != nil {
			return err
		}
		event := journal.Event{EventSchema: journal.EventSchemaV2, EventType: "gate.registered", ActorV2: journal.EventActor{ActorID: actor.ActorID, ActorType: actor.ActorType, Producer: actor.Producer, ProvenanceRef: actor.ProvenanceRef, SessionRef: actor.SessionRef}, Aggregate: journal.Aggregate{ID: registration.ID, Sequence: registration.Version, Type: "gate-registration"}, CommandName: "gate.register", Repository: repositoryEventValue(options.Request.Repository), Effects: []journal.Effect{{Op: "insert", Table: "gate_registrations", PK: map[string]any{"id": registration.ID, "version": registration.Version}, Values: map[string]any{"activation_id": activationID, "definition_id": definition.DefinitionID, "definition_version": definition.DefinitionVersion, "definition_digest": digest, "binding_digest": bindingDigest, "project_id": project.ProjectID, "worktree_id": options.Request.WorktreeID, "scope_json": string(scopeJSON), "policy_json": string(policyJSON), "authorization_provenance": actor.ProvenanceRef, "created_event_id": "<event-id>", "created_at": createdAt}}}}
		outcome, err := journal.CommitIdempotent(lock, conn, journal.IdempotencyRequest{OperationScope: "gate.register:" + project.ProjectID, IdempotencyKey: options.IdempotencyKey, RequestDigest: requestDigest, RequestFieldDigests: fields}, event, func(prepared journal.Event) (journal.IdempotencyMutation, error) {
			prepared.Effects[0].Values["created_event_id"] = prepared.EventID
			return journal.IdempotencyMutation{EntityID: registration.ID, ResultJSON: resultJSON}, nil
		})
		if err != nil {
			return err
		}
		loaded, err := Resolve(options.WorkspaceRoot, options.MigrationsDir, registration.ID, registration.Version)
		if err != nil {
			return err
		}
		result = RegisterResult{Registration: loaded, IdempotencyKey: options.IdempotencyKey, Replayed: outcome.Replayed, EventID: outcome.EventID}
		return nil
	})
	if err != nil {
		return RegisterResult{}, err
	}
	return result, nil
}

func Resolve(workspaceRoot, migrationsDir, id string, version int) (Registration, error) {
	if _, err := journal.EnsureProjectionFresh(workspaceRoot, migrationsDir); err != nil {
		return Registration{}, err
	}
	conn, err := openEvidenceDB(workspaceRoot)
	if err != nil {
		return Registration{}, err
	}
	defer func() { _ = conn.Close() }()
	return loadRegistration(conn, id, version)
}

func openEvidenceDB(workspaceRoot string) (*sql.DB, error) {
	return fiddb.Open(filepath.Join(workspaceRoot, config.WorkspaceDirName, config.DatabaseFileName))
}

type storedActor struct {
	ActorID, ActorType, Producer, ProvenanceRef string
	SessionRef                                  *string
	Capabilities                                []string
}

func loadActiveActor(conn *sql.DB, actorID string) (string, storedActor, error) {
	var incomplete int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM kit_activations WHERE status IN ('prepared', 'staged', 'committed')`).Scan(&incomplete); err != nil {
		return "", storedActor{}, fmt.Errorf("load active activation: read incomplete activation: %w", err)
	}
	if incomplete != 0 {
		return "", storedActor{}, &ActiveActivationUnavailableError{}
	}
	var activationID string
	if err := conn.QueryRow(`SELECT id FROM kit_activations WHERE status = 'finalized' LIMIT 1`).Scan(&activationID); err == sql.ErrNoRows {
		return "", storedActor{}, &ActiveActivationUnavailableError{}
	} else if err != nil {
		return "", storedActor{}, fmt.Errorf("load active activation: %w", err)
	}
	var actor storedActor
	var capabilities string
	err := conn.QueryRow(`SELECT actor_id, actor_type, producer, provenance_ref, session_ref, capabilities_json FROM activation_actor_profiles WHERE activation_id = ? AND actor_id = ?`, activationID, strings.TrimSpace(actorID)).Scan(&actor.ActorID, &actor.ActorType, &actor.Producer, &actor.ProvenanceRef, &actor.SessionRef, &capabilities)
	if err == sql.ErrNoRows {
		return "", storedActor{}, &TrustedActorUnavailableError{ActorID: actorID}
	}
	if err != nil {
		return "", storedActor{}, fmt.Errorf("load trusted actor: %w", err)
	}
	if err := json.Unmarshal([]byte(capabilities), &actor.Capabilities); err != nil {
		return "", storedActor{}, fmt.Errorf("decode trusted actor capabilities: %w", err)
	}
	return activationID, actor, nil
}

func loadActivatedDefinition(conn *sql.DB, activationID, definitionID string, version int) (GateDefinition, string, error) {
	var raw, digest string
	err := conn.QueryRow(`SELECT definition_json, definition_digest FROM activation_gate_definitions WHERE activation_id = ? AND definition_id = ? AND definition_version = ?`, activationID, definitionID, version).Scan(&raw, &digest)
	if err == sql.ErrNoRows {
		return GateDefinition{}, "", &GateDefinitionNotFoundError{DefinitionID: definitionID, Version: version}
	}
	if err != nil {
		return GateDefinition{}, "", fmt.Errorf("load activated definition: %w", err)
	}
	var definition GateDefinition
	decoder := json.NewDecoder(strings.NewReader(raw))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&definition); err != nil {
		return GateDefinition{}, "", fmt.Errorf("decode activated definition: %w", err)
	}
	return definition, digest, nil
}

func nextRegistrationIdentity(conn *sql.DB, definitionID string) (string, int, error) {
	var id string
	var version int
	err := conn.QueryRow(`SELECT id, version FROM gate_registrations WHERE definition_id = ? ORDER BY version DESC LIMIT 1`, definitionID).Scan(&id, &version)
	if err == sql.ErrNoRows {
		generated, err := journal.GenerateID("REG", func(candidate string) (bool, error) {
			var found int
			err := conn.QueryRow(`SELECT 1 FROM gate_registrations WHERE id = ?`, candidate).Scan(&found)
			return err == nil, nil
		})
		return generated, 1, err
	}
	if err != nil {
		return "", 0, fmt.Errorf("load registration version: %w", err)
	}
	return id, version + 1, nil
}

func loadRegistrationByBinding(conn *sql.DB, definitionDigest, bindingDigest string) (Registration, bool, error) {
	var id string
	var version int
	err := conn.QueryRow(`SELECT id, version FROM gate_registrations WHERE definition_digest = ? AND binding_digest = ?`, definitionDigest, bindingDigest).Scan(&id, &version)
	if err == sql.ErrNoRows {
		return Registration{}, false, nil
	}
	if err != nil {
		return Registration{}, false, fmt.Errorf("load registration binding: %w", err)
	}
	registration, err := loadRegistration(conn, id, version)
	return registration, true, err
}

func loadRegistration(conn *sql.DB, id string, version int) (Registration, error) {
	var registration Registration
	var scopeJSON, policyJSON string
	err := conn.QueryRow(`SELECT id, version, activation_id, definition_id, definition_version, definition_digest, binding_digest, project_id, worktree_id, scope_json, authorization_provenance, policy_json, created_at FROM gate_registrations WHERE id = ? AND version = ?`, id, version).Scan(&registration.ID, &registration.Version, &registration.ActivationID, &registration.DefinitionID, &registration.DefinitionVersion, &registration.DefinitionDigest, &registration.BindingDigest, &registration.ProjectID, &registration.WorktreeID, &scopeJSON, &registration.AuthorizationProvenance, &policyJSON, &registration.CreatedAt)
	if err == sql.ErrNoRows {
		return Registration{}, &RegistrationNotFoundError{ID: id, Version: version}
	}
	if err != nil {
		return Registration{}, fmt.Errorf("load registration: %w", err)
	}
	if err := json.Unmarshal([]byte(scopeJSON), &registration.Scope); err != nil {
		return Registration{}, fmt.Errorf("decode registration scope: %w", err)
	}
	if err := json.Unmarshal([]byte(policyJSON), &registration.Policy); err != nil {
		return Registration{}, fmt.Errorf("decode registration policy: %w", err)
	}
	return registration, nil
}

func existingCreatedEventID(conn *sql.DB, registration Registration) string {
	var eventID string
	if err := conn.QueryRow(`SELECT created_event_id FROM gate_registrations WHERE id = ? AND version = ?`, registration.ID, registration.Version).Scan(&eventID); err != nil {
		return ""
	}
	return eventID
}

func hasCapability(capabilities []string, required string) bool {
	for _, capability := range capabilities {
		if capability == required {
			return true
		}
	}
	return false
}

func validateRegisterRequest(request RegisterRequest) error {
	if strings.TrimSpace(request.DefinitionID) == "" || request.DefinitionVersion < 1 || strings.TrimSpace(request.WorktreeID) == "" || strings.TrimSpace(request.Scope.IssueID) == "" {
		return &ExecutorPolicyUnsupportedError{Fields: []string{"binding"}}
	}
	return nil
}

func validateRepositoryBinding(request RegisterRequest, projectID string) error {
	repository := request.Repository
	if repository.CaptureSource != "git" || repository.ProjectID != projectID || repository.WorktreeID != request.WorktreeID || !validOID(repository.HeadOID) || !validOID(repository.TreeOID) || (repository.BaseOID != nil && !validOID(*repository.BaseOID)) || (repository.Dirty && (repository.DirtyDigest == nil || !validDigest(*repository.DirtyDigest))) || (!repository.Dirty && repository.DirtyDigest != nil) {
		return &ExecutorPolicyUnsupportedError{Fields: []string{"repository"}}
	}
	return nil
}

func validOID(value string) bool {
	return len(value) == 40 && strings.Trim(value, "0123456789abcdef") == ""
}

func validDigest(value string) bool {
	return strings.HasPrefix(value, "sha256:") && len(value) == 71 && strings.TrimPrefix(value, "sha256:") != ""
}

func gateBindingDigest(request RegisterRequest) (string, error) {
	payload, err := journal.MarshalCanonicalJSON(map[string]any{"repository": request.Repository, "scope": request.Scope, "worktree_id": request.WorktreeID})
	if err != nil {
		return "", fmt.Errorf("canonicalize gate binding: %w", err)
	}
	sum := sha256.Sum256(payload)
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}

func resultJSON(registration Registration, key string) (string, error) {
	payload, err := journal.MarshalCanonicalJSON(map[string]any{"registration": registration, "idempotency_key": key, "replayed": false})
	if err != nil {
		return "", err
	}
	return string(payload), nil
}

func repositoryEventValue(repository Repository) *journal.Repository {
	base := ""
	if repository.BaseOID != nil {
		base = *repository.BaseOID
	}
	return &journal.Repository{BaseOID: base, CaptureSource: repository.CaptureSource, Dirty: repository.Dirty, DirtyDigest: repository.DirtyDigest, HeadOID: repository.HeadOID, ProjectID: repository.ProjectID, TreeOID: repository.TreeOID, WorktreeID: repository.WorktreeID}
}
