package epicwork

import (
	"database/sql"
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"time"

	"github.com/mbell-nzl/fid/internal/artifactstore"
	"github.com/mbell-nzl/fid/internal/checkrecord"
	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/journal"
)

func keepArtifacts(o Options, q *sql.DB, root string, result *checkrecord.Result, definition *checkrecord.Definition, effects *[]journal.Effect) error {
	cfg, err := config.ReadProjectConfig(filepath.Join(config.LedgerDir(o.WorkspaceRoot), config.ConfigFileName))
	if err != nil {
		return err
	}
	ledger, err := filepath.EvalSymlinks(config.LedgerDir(o.WorkspaceRoot))
	if err != nil {
		return err
	}
	store := filepath.Join(ledger, cfg.ArtifactRoot)
	if err := os.MkdirAll(store, 0o700); err != nil {
		return err
	}
	store, err = filepath.EvalSymlinks(store)
	if err != nil {
		return err
	}
	if !artifactstore.Contained(ledger, store) {
		return fmt.Errorf("artifact store is outside the ledger")
	}
	for _, path := range definition.Artifacts {
		file, reason, err := artifactstore.Keep(root, store, path, checkrecord.ArtifactLimit)
		if err != nil {
			return err
		}
		if reason != "" {
			result.RefusedArtifacts = append(result.RefusedArtifacts, checkrecord.RefusedArtifact{Path: path, Reason: reason})
			continue
		}
		var id, availability string
		err = q.QueryRow(`SELECT id,availability FROM referenced_artifacts WHERE digest=?`, file.Digest).Scan(&id, &availability)
		if err != nil && err != sql.ErrNoRows {
			return err
		}
		for _, e := range *effects {
			if e.Table == "referenced_artifacts" && e.Op == "insert" && e.Values["digest"] == file.Digest {
				id = e.PK["id"].(string)
			}
		}
		if id == "" {
			id, err = journal.GenerateID("ART", func(candidate string) (bool, error) {
				var n int
				err := q.QueryRow(`SELECT count(*) FROM referenced_artifacts WHERE id=?`, candidate).Scan(&n)
				return n > 0, err
			})
			if err != nil {
				return err
			}
			media := mime.TypeByExtension(filepath.Ext(path))
			if media == "" {
				media = "application/octet-stream"
			}
			*effects = append(*effects, journal.Effect{Op: "insert", Table: "referenced_artifacts", PK: map[string]any{"id": id}, Values: map[string]any{"digest": file.Digest, "media_type": media, "size_bytes": file.SizeBytes, "retention_class": "evidence", "availability": "available", "location_policy": "content-addressed", "created_at": o.now().Format(time.RFC3339), "expires_at": nil}})
		} else if availability != "" && availability != "available" {
			*effects = append(*effects, journal.Effect{Op: "update", Table: "referenced_artifacts", PK: map[string]any{"id": id}, Values: map[string]any{"availability": "available"}, Preconditions: map[string]any{"availability": availability}})
		}
		result.Artifacts = append(result.Artifacts, checkrecord.Artifact{ID: id, Path: path, Digest: file.Digest, SizeBytes: file.SizeBytes})
	}
	return nil
}
