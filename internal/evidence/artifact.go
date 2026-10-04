package evidence

import (
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/journal"
	fidproject "github.com/mbell-nzl/fid/internal/project"
	"github.com/mbell-nzl/fid/internal/workflow"
)

type ArtifactVerifyRequest struct {
	ArtifactID string `json:"artifact_id"`
}

type ArtifactVerifyOptions struct {
	WorkspaceRoot  string
	MigrationsDir  string
	ActorID        string
	IdempotencyKey string
	Request        ArtifactVerifyRequest
}

type ArtifactVerifyResult struct {
	Artifact        Artifact `json:"artifact"`
	Changed         bool     `json:"changed"`
	IdempotencyKey  string   `json:"idempotency_key"`
	Replayed        bool     `json:"replayed"`
	EventID         string
	EventOccurredAt string
}

// Artifact is immutable content-addressed material retained from a Gate run.
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

type retainedArtifact struct {
	Artifact
	DeclaredPath string
}

func retainArtifacts(workspaceRoot string, definition GateDefinition) ([]retainedArtifact, []ArtifactRejection, error) {
	project, err := config.ReadProjectConfig(filepath.Join(workspaceRoot, config.WorkspaceDirName, config.ConfigFileName))
	if err != nil {
		return nil, nil, fmt.Errorf("retain Artifacts: read project config: %w", err)
	}
	storeRoot, err := artifactStoreRoot(workspaceRoot, project.ArtifactRoot)
	if err != nil {
		return nil, nil, err
	}

	sources := redactionSources(definition)
	retained := make([]retainedArtifact, 0, len(definition.ArtifactPaths))
	rejected := make([]ArtifactRejection, 0)
	for _, declaredPath := range definition.ArtifactPaths {
		candidate, reason := artifactCandidatePath(workspaceRoot, definition.CWD, declaredPath)
		if reason != "" {
			rejected = append(rejected, ArtifactRejection{Path: redactText(declaredPath, sources), Reason: reason})
			continue
		}
		item, reason, err := publishArtifactCandidate(candidate, declaredPath, storeRoot, definition.ArtifactBytes, sources, nil)
		if err != nil {
			return nil, nil, err
		}
		if reason != "" {
			rejected = append(rejected, ArtifactRejection{Path: redactText(declaredPath, sources), Reason: reason})
			continue
		}
		retained = append(retained, item)
	}
	sort.Slice(rejected, func(i, j int) bool {
		if rejected[i].Path == rejected[j].Path {
			return rejected[i].Reason < rejected[j].Reason
		}
		return rejected[i].Path < rejected[j].Path
	})
	return retained, rejected, nil
}

func artifactStoreRoot(workspaceRoot, root string) (string, error) {
	fidRoot, err := filepath.EvalSymlinks(filepath.Join(workspaceRoot, config.WorkspaceDirName))
	if err != nil {
		return "", fmt.Errorf("retain Artifacts: resolve fid root: %w", err)
	}
	store := filepath.Join(fidRoot, root)
	if err := os.MkdirAll(store, 0o700); err != nil {
		return "", fmt.Errorf("retain Artifacts: create store: %w", err)
	}
	resolved, err := filepath.EvalSymlinks(store)
	if err != nil || !containedResolvedPath(fidRoot, resolved) {
		return "", fmt.Errorf("retain Artifacts: configured store escapes fid root")
	}
	return resolved, nil
}

func artifactCandidatePath(workspaceRoot, cwd, declaredPath string) (string, string) {
	workspace, err := filepath.EvalSymlinks(workspaceRoot)
	if err != nil {
		return "", "outside-cwd"
	}
	resolvedCWD, err := filepath.EvalSymlinks(filepath.Join(workspace, cwd))
	if err != nil || !containedResolvedPath(workspace, resolvedCWD) {
		return "", "outside-cwd"
	}
	candidate := filepath.Join(resolvedCWD, declaredPath)
	parent, err := filepath.EvalSymlinks(filepath.Dir(candidate))
	if err != nil || !containedResolvedPath(resolvedCWD, parent) {
		return "", "outside-cwd"
	}
	info, err := os.Lstat(candidate)
	if os.IsNotExist(err) {
		return "", "missing"
	}
	if err != nil {
		return "", "outside-cwd"
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return "", "symlink"
	}
	if info.IsDir() {
		return "", "directory"
	}
	if !info.Mode().IsRegular() {
		return "", "non-regular"
	}
	return candidate, ""
}

func publishArtifactCandidate(candidate, declaredPath, storeRoot string, limit int64, sources []string, beforeRevalidate func()) (retainedArtifact, string, error) {
	before, err := os.Lstat(candidate)
	if err != nil {
		return retainedArtifact{}, "changed-during-copy", nil
	}
	input, err := os.Open(candidate)
	if err != nil {
		return retainedArtifact{}, "changed-during-copy", nil
	}
	defer func() { _ = input.Close() }()

	staging, err := os.CreateTemp(storeRoot, ".artifact-")
	if err != nil {
		return retainedArtifact{}, "", fmt.Errorf("retain Artifact: create staging file: %w", err)
	}
	stagingPath := staging.Name()
	defer func() { _ = os.Remove(stagingPath) }()

	hash := sha256.New()
	writer := io.MultiWriter(staging, hash)
	redactor := newLiveRedactor(sources)
	buffer := make([]byte, 32*1024)
	var observed int64
	for {
		count, readErr := input.Read(buffer)
		if count > 0 {
			observed += int64(count)
			if observed > limit {
				_ = staging.Close()
				return retainedArtifact{}, "oversized", nil
			}
			if err := writeSanitizedArtifactChunk(writer, redactor, buffer[:count], false); err != nil {
				_ = staging.Close()
				return retainedArtifact{}, "", fmt.Errorf("retain Artifact: write staging file: %w", err)
			}
		}
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			_ = staging.Close()
			return retainedArtifact{}, "", fmt.Errorf("retain Artifact: read candidate: %w", readErr)
		}
	}
	if err := writeSanitizedArtifactChunk(writer, redactor, nil, true); err != nil {
		_ = staging.Close()
		return retainedArtifact{}, "", fmt.Errorf("retain Artifact: finalize sanitized content: %w", err)
	}
	if err := staging.Sync(); err != nil {
		_ = staging.Close()
		return retainedArtifact{}, "", fmt.Errorf("retain Artifact: sync staging file: %w", err)
	}
	if err := staging.Close(); err != nil {
		return retainedArtifact{}, "", fmt.Errorf("retain Artifact: close staging file: %w", err)
	}

	if beforeRevalidate != nil {
		beforeRevalidate()
	}
	after, err := os.Lstat(candidate)
	if err != nil || !os.SameFile(before, after) || before.Size() != after.Size() || !before.ModTime().Equal(after.ModTime()) {
		return retainedArtifact{}, "changed-during-copy", nil
	}
	if parent, err := filepath.EvalSymlinks(filepath.Dir(candidate)); err != nil || !containedResolvedPath(filepath.Dir(candidate), parent) {
		return retainedArtifact{}, "changed-during-copy", nil
	}

	digest := "sha256:" + fmt.Sprintf("%x", hash.Sum(nil))
	destination := filepath.Join(storeRoot, strings.TrimPrefix(digest, "sha256:"))
	if err := publishArtifactFile(stagingPath, destination, digest); err != nil {
		return retainedArtifact{}, "", err
	}
	info, err := os.Stat(destination)
	if err != nil {
		return retainedArtifact{}, "", fmt.Errorf("retain Artifact: stat published file: %w", err)
	}
	mediaType := mime.TypeByExtension(filepath.Ext(declaredPath))
	if mediaType == "" {
		mediaType = "application/octet-stream"
	}
	artifact := Artifact{Digest: digest, MediaType: mediaType, SizeBytes: info.Size(), RetentionClass: "evidence", Availability: "available", LocationPolicy: "content-addressed", DeclaredPath: declaredPath}
	return retainedArtifact{Artifact: artifact, DeclaredPath: declaredPath}, "", nil
}

func writeSanitizedArtifactChunk(writer io.Writer, redactor *liveRedactor, data []byte, final bool) error {
	if len(redactor.sources) == 0 && utf8.Valid(data) {
		_, err := writer.Write(data)
		return err
	}
	_, err := io.WriteString(writer, redactor.Transform(data, final))
	return err
}

func publishArtifactFile(stagingPath, destination, digest string) error {
	info, err := os.Lstat(destination)
	if err == nil {
		if !info.Mode().IsRegular() {
			return fmt.Errorf("retain Artifact: existing digest path is not regular")
		}
		return verifyArtifactDigest(destination, digest)
	}
	if !os.IsNotExist(err) {
		return fmt.Errorf("retain Artifact: inspect digest destination: %w", err)
	}
	if err := os.Rename(stagingPath, destination); err != nil {
		if !os.IsExist(err) {
			return fmt.Errorf("retain Artifact: publish by digest: %w", err)
		}
		return verifyArtifactDigest(destination, digest)
	}
	return nil
}

func verifyArtifactDigest(path, digest string) error {
	info, err := os.Lstat(path)
	if err != nil {
		return fmt.Errorf("verify Artifact digest: inspect: %w", err)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("verify Artifact digest: stored content is not regular")
	}
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("verify Artifact digest: open: %w", err)
	}
	defer func() { _ = file.Close() }()
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return fmt.Errorf("verify Artifact digest: read: %w", err)
	}
	if "sha256:"+fmt.Sprintf("%x", hash.Sum(nil)) != digest {
		return fmt.Errorf("verify Artifact digest: stored content does not match digest")
	}
	return nil
}

func containedResolvedPath(root, candidate string) bool {
	relative, err := filepath.Rel(root, candidate)
	return err == nil && relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator))
}

func receiptArtifactEffects(conn *sql.DB, retained []retainedArtifact, createdAt time.Time) ([]Artifact, []journal.Effect, error) {
	byDigest := make(map[string]Artifact)
	effects := make([]journal.Effect, 0, len(retained))
	for _, item := range retained {
		if _, exists := byDigest[item.Digest]; exists {
			continue
		}
		artifact, found, err := loadArtifactByDigest(conn, item.Digest)
		if err != nil {
			return nil, nil, err
		}
		if !found {
			id, err := journal.GenerateID("ART", func(candidate string) (bool, error) {
				var count int
				err := conn.QueryRow(`SELECT COUNT(*) FROM referenced_artifacts WHERE id = ?`, candidate).Scan(&count)
				return count > 0, err
			})
			if err != nil {
				return nil, nil, err
			}
			artifact = item.Artifact
			artifact.ID = id
			artifact.CreatedAt = createdAt.UTC().Format(time.RFC3339Nano)
			effects = append(effects, journal.Effect{Op: "insert", Table: "referenced_artifacts", PK: map[string]any{"id": artifact.ID}, Values: map[string]any{
				"digest": artifact.Digest, "media_type": artifact.MediaType, "size_bytes": artifact.SizeBytes,
				"retention_class": artifact.RetentionClass, "availability": artifact.Availability,
				"location_policy": artifact.LocationPolicy, "created_at": artifact.CreatedAt, "expires_at": nil,
			}})
		}
		artifact.DeclaredPath = item.DeclaredPath
		byDigest[item.Digest] = artifact
	}

	artifacts := make([]Artifact, 0, len(byDigest))
	for _, artifact := range byDigest {
		artifacts = append(artifacts, artifact)
	}
	sort.Slice(artifacts, func(i, j int) bool { return artifacts[i].ID < artifacts[j].ID })
	return artifacts, effects, nil
}

func loadArtifactByDigest(conn *sql.DB, digest string) (Artifact, bool, error) {
	var artifact Artifact
	err := conn.QueryRow(`SELECT id, digest, media_type, size_bytes, retention_class, availability, location_policy, created_at, expires_at FROM referenced_artifacts WHERE digest = ?`, digest).Scan(
		&artifact.ID, &artifact.Digest, &artifact.MediaType, &artifact.SizeBytes, &artifact.RetentionClass,
		&artifact.Availability, &artifact.LocationPolicy, &artifact.CreatedAt, &artifact.ExpiresAt,
	)
	if err == sql.ErrNoRows {
		return Artifact{}, false, nil
	}
	if err != nil {
		return Artifact{}, false, fmt.Errorf("load Artifact by digest: %w", err)
	}
	return artifact, true, nil
}

func loadReceiptArtifacts(conn *sql.DB, receiptID string) ([]Artifact, error) {
	rows, err := conn.Query(`SELECT a.id, a.digest, a.media_type, a.size_bytes, a.retention_class, a.availability, a.location_policy, a.created_at, a.expires_at, ra.declared_path FROM receipt_artifacts ra JOIN referenced_artifacts a ON a.id = ra.artifact_id WHERE ra.receipt_id = ? ORDER BY a.id`, receiptID)
	if err != nil {
		return nil, fmt.Errorf("load Receipt Artifacts: %w", err)
	}
	defer func() { _ = rows.Close() }()
	artifacts := make([]Artifact, 0)
	for rows.Next() {
		var artifact Artifact
		if err := rows.Scan(&artifact.ID, &artifact.Digest, &artifact.MediaType, &artifact.SizeBytes, &artifact.RetentionClass, &artifact.Availability, &artifact.LocationPolicy, &artifact.CreatedAt, &artifact.ExpiresAt, &artifact.DeclaredPath); err != nil {
			return nil, fmt.Errorf("load Receipt Artifact row: %w", err)
		}
		artifacts = append(artifacts, artifact)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("load Receipt Artifacts: %w", err)
	}
	return artifacts, nil
}

func containsArtifactRejection(values []ArtifactRejection, target ArtifactRejection) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func VerifyArtifact(options ArtifactVerifyOptions) (ArtifactVerifyResult, error) {
	if strings.TrimSpace(options.WorkspaceRoot) == "" || strings.TrimSpace(options.MigrationsDir) == "" || strings.TrimSpace(options.ActorID) == "" || strings.TrimSpace(options.Request.ArtifactID) == "" {
		return ArtifactVerifyResult{}, fmt.Errorf("verify Artifact: workspace root, migrations directory, actor, and Artifact are required")
	}
	var result ArtifactVerifyResult
	err := journal.WithLock(options.WorkspaceRoot, journal.LockOwner{Command: "fid artifact verify"}, func(lock *journal.WriteLock) error {
		if _, err := journal.EnsureProjectionFreshWithLock(lock, options.MigrationsDir); err != nil {
			return err
		}
		conn, err := openEvidenceDB(options.WorkspaceRoot)
		if err != nil {
			return err
		}
		defer func() { _ = conn.Close() }()
		project, err := config.ReadProjectConfig(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.ConfigFileName))
		if err != nil {
			return fmt.Errorf("verify Artifact: read project config: %w", err)
		}
		_, actor, err := loadActiveActor(conn, options.ActorID)
		if err != nil {
			return err
		}
		if !hasCapability(actor.Capabilities, "artifact.verify") {
			return &AuthorityDeniedError{ActorID: options.ActorID}
		}
		artifact, found, err := loadArtifactByID(conn, options.Request.ArtifactID)
		if err != nil {
			return err
		}
		if !found {
			return fmt.Errorf("verify Artifact: Artifact %q was not found", options.Request.ArtifactID)
		}
		digest, fields, err := journal.ComputeRequestDigest(map[string]any{"actor_id": options.ActorID, "artifact_id": artifact.ID, "project_id": project.ProjectID, "request": options.Request}, options.ActorID, nil)
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

		storeRoot, err := artifactStoreRoot(options.WorkspaceRoot, project.ArtifactRoot)
		if err != nil {
			return err
		}
		available := verifyArtifactDigest(filepath.Join(storeRoot, strings.TrimPrefix(artifact.Digest, "sha256:")), artifact.Digest) == nil
		observed := "missing"
		if available {
			observed = "available"
		}
		priorAvailability := artifact.Availability
		changed := priorAvailability != observed
		artifact.Availability = observed
		sequence, err := journal.NextAggregateSequence(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.JournalDirName), journal.Aggregate{ID: artifact.ID, Type: "artifact"})
		if err != nil {
			return err
		}
		effects := []journal.Effect{}
		if changed {
			effects = append(effects, journal.Effect{Op: "update", Table: "referenced_artifacts", PK: map[string]any{"id": artifact.ID}, Values: map[string]any{"availability": observed}, Preconditions: map[string]any{"availability": priorAvailability}})
		}
		if changed && observed == "available" {
			targets, err := artifactEnvelopeTargets(conn, artifact.ID)
			if err != nil {
				return err
			}
			if len(targets) > 0 {
				observation, err := fidproject.ObserveRepository(options.WorkspaceRoot)
				if err != nil {
					return fmt.Errorf("verify Artifact: observe repository: %w", err)
				}
				head := observation.HeadOID
				repository := workflow.Repository{BaseOID: &head, CaptureSource: "git", Dirty: observation.Dirty, DirtyDigest: observation.DirtyDigest, HeadOID: observation.HeadOID, ProjectID: project.ProjectID, TreeOID: observation.TreeOID, WorktreeID: observation.WorktreeID}
				for _, target := range targets {
					overrides, err := workflow.ArtifactAvailabilityOverrides(conn, target.EnvelopeID, target.Revision, artifact.ID, repository)
					if err != nil {
						return err
					}
					readinessEffects, _, _, err := workflow.HumanOutcomeReadinessStopEffects(conn, target.EnvelopeID, target.Revision, actor.ActorID, repository, overrides)
					if err != nil {
						return err
					}
					effects = append(effects, readinessEffects...)
				}
			}
		}
		event := journal.Event{EventSchema: journal.EventSchemaV2, EventType: "artifact.availability-observed", ActorV2: journal.EventActor{ActorID: actor.ActorID, ActorType: actor.ActorType, Producer: actor.Producer, ProvenanceRef: actor.ProvenanceRef, SessionRef: actor.SessionRef}, Aggregate: journal.Aggregate{ID: artifact.ID, Sequence: sequence, Type: "artifact"}, CommandName: "artifact.verify", Effects: effects}
		result = ArtifactVerifyResult{Artifact: artifact, Changed: changed, IdempotencyKey: options.IdempotencyKey, Replayed: false}
		resultJSON, err := journal.MarshalCanonicalJSON(result)
		if err != nil {
			return err
		}
		outcome, err := journal.CommitIdempotent(lock, conn, idempotency, event, func(prepared journal.Event) (journal.IdempotencyMutation, error) {
			substituteEvidenceEventPlaceholders(prepared.Effects, prepared)
			return journal.IdempotencyMutation{EntityID: artifact.ID, ResultJSON: string(resultJSON)}, nil
		})
		if err != nil {
			return err
		}
		result.EventID = outcome.EventID
		return nil
	})
	if err != nil {
		return ArtifactVerifyResult{}, err
	}
	if result.EventID != "" {
		occurredAt, err := gateEventOccurredAt(options.WorkspaceRoot, result.EventID)
		if err != nil {
			return ArtifactVerifyResult{}, err
		}
		result.EventOccurredAt = occurredAt
	}
	return result, nil
}

type artifactEnvelopeTarget struct {
	EnvelopeID string
	Revision   int
}

// artifactEnvelopeTargets finds every Envelope revision whose linked Evidence
// requires this Artifact, so a missing-to-available flip can re-evaluate each
// one for the human-outcome-required readiness Stop (DEC-224).
func artifactEnvelopeTargets(conn *sql.DB, artifactID string) ([]artifactEnvelopeTarget, error) {
	rows, err := conn.Query(`SELECT DISTINCT o.envelope_id,o.revision FROM evidence_artifact_links l JOIN evidence_records e ON e.id=l.evidence_id JOIN envelope_obligations o ON o.id=e.obligation_id WHERE l.artifact_id=? AND l.required_material=1 ORDER BY o.envelope_id,o.revision`, artifactID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	targets := []artifactEnvelopeTarget{}
	for rows.Next() {
		var target artifactEnvelopeTarget
		if err := rows.Scan(&target.EnvelopeID, &target.Revision); err != nil {
			return nil, err
		}
		targets = append(targets, target)
	}
	return targets, rows.Err()
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
