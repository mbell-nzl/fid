package evidence

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/mbell-nzl/fid/internal/artifactstore"
	"github.com/mbell-nzl/fid/internal/artifactverify"
	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/journal"
	fidproject "github.com/mbell-nzl/fid/internal/project"
	"github.com/mbell-nzl/fid/internal/workflow"
)

type ArtifactVerifyRequest = artifactverify.Request
type ArtifactVerifyOptions = artifactverify.Options
type ArtifactVerifyResult = artifactverify.Result
type Artifact = artifactverify.Artifact

type retainedArtifact struct {
	Artifact
	DeclaredPath string
}

func retainArtifacts(workspaceRoot string, definition GateDefinition) ([]retainedArtifact, []ArtifactRejection, error) {
	project, err := config.ReadProjectConfig(filepath.Join(config.LedgerDir(workspaceRoot), config.ConfigFileName))
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
	fidRoot, err := filepath.EvalSymlinks(config.LedgerDir(workspaceRoot))
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
	return artifactstore.Publish(stagingPath, destination, digest)
}
func verifyArtifactDigest(path, digest string) error { return artifactstore.Verify(path, digest) }
func containedResolvedPath(root, candidate string) bool {
	return artifactstore.Contained(root, candidate)
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

// VerifyArtifact preserves kernel authority and readiness until RESET 1.5c2
// switches the command adapter to the lean verifier.
func VerifyArtifact(options ArtifactVerifyOptions) (ArtifactVerifyResult, error) {
	return artifactverify.VerifyWithHooks(options, artifactverify.Hooks{
		Authorize: func(conn *sql.DB) (journal.EventActor, error) {
			_, actor, err := loadActiveActor(conn, options.ActorID)
			if err != nil {
				return journal.EventActor{}, err
			}
			if !hasCapability(actor.Capabilities, "artifact.verify") {
				return journal.EventActor{}, &AuthorityDeniedError{ActorID: options.ActorID}
			}
			return journal.EventActor{ActorID: actor.ActorID, ActorType: actor.ActorType, Producer: actor.Producer, ProvenanceRef: actor.ProvenanceRef, SessionRef: actor.SessionRef}, nil
		},
		Restored: func(conn *sql.DB, artifactID string, actor journal.EventActor) ([]journal.Effect, error) {
			targets, err := artifactEnvelopeTargets(conn, artifactID)
			if err != nil || len(targets) == 0 {
				return nil, err
			}
			project, err := config.ReadProjectConfig(filepath.Join(config.LedgerDir(options.WorkspaceRoot), config.ConfigFileName))
			if err != nil {
				return nil, err
			}
			observation, err := fidproject.ObserveRepository(options.WorkspaceRoot)
			if err != nil {
				return nil, fmt.Errorf("verify Artifact: observe repository: %w", err)
			}
			head := observation.HeadOID
			repository := workflow.Repository{BaseOID: &head, CaptureSource: "git", Dirty: observation.Dirty, DirtyDigest: observation.DirtyDigest, HeadOID: observation.HeadOID, ProjectID: project.ProjectID, TreeOID: observation.TreeOID, WorktreeID: observation.WorktreeID}
			effects := []journal.Effect{}
			for _, target := range targets {
				overrides, err := workflow.ArtifactAvailabilityOverrides(conn, target.EnvelopeID, target.Revision, artifactID, repository)
				if err != nil {
					return nil, err
				}
				readiness, _, _, err := workflow.HumanOutcomeReadinessStopEffects(conn, target.EnvelopeID, target.Revision, actor.ActorID, repository, overrides)
				if err != nil {
					return nil, err
				}
				effects = append(effects, readiness...)
			}
			return effects, nil
		},
		Prepare: func(event journal.Event) { substituteEvidenceEventPlaceholders(event.Effects, event) },
	})
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
