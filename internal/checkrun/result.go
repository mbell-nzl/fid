package checkrun

import (
	"database/sql"
	"errors"
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"time"

	"github.com/mbell-nzl/fid/internal/artifactstore"
	"github.com/mbell-nzl/fid/internal/checkrecord"
	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/holds"
	"github.com/mbell-nzl/fid/internal/holdwork"
	"github.com/mbell-nzl/fid/internal/journal"
)

func asHoldRefusal(err error, target **holdwork.Refusal) bool { return errors.As(err, target) }

func (ex *execution) renew(id string) error {
	return ex.options.write(func(db store) error {
		if err := db.current(); err != nil {
			return err
		}
		hold, err := holds.CurrentByID(db.conn, ex.options.HoldID)
		if err != nil {
			return err
		}
		event, err := db.event("check.hold-renewed", id, []journal.Effect{{Op: "update", Table: "holds", PK: map[string]any{"id": hold.ID}, Values: map[string]any{"last_heartbeat_at": holds.Timestamp(db.now), "expires_at": holds.Timestamp(db.now.Add(holds.Lifetime))}, Preconditions: map[string]any{"ended_at": nil, "expires_at": holds.Timestamp(hold.ExpiresAt)}}})
		if err != nil {
			return err
		}
		_, err = journal.CommitEvent(db.lock, db.conn, event)
		return err
	})
}

func (ex *execution) commitProgress(db store, event journal.Event) error {
	prepared, err := journal.PrepareEvent(db.lock, event)
	if err != nil {
		return err
	}
	if ex.response.Complete {
		ex.response.EventID = prepared.EventID
		ex.response.OccurredAt = prepared.OccurredAt.UTC().Format(time.RFC3339Nano)
	}
	fields := encoded(ex.request.RequestFieldDigests)
	prepared.Idempotency = &journal.EventIdempotency{Key: ex.request.IdempotencyKey, Scope: ex.request.OperationScope, RequestDigest: ex.request.RequestDigest}
	prepared.Effects = append(prepared.Effects, journal.Effect{Op: "update", Table: "idempotency_results", PK: map[string]any{"operation_scope": ex.request.OperationScope, "idempotency_key": ex.request.IdempotencyKey}, Values: map[string]any{"result_entity_id": event.Aggregate.ID, "result_event_id": prepared.EventID, "aggregate_sequence": event.Aggregate.Sequence, "result_json": encoded(ex.response)}, Preconditions: map[string]any{"request_digest": ex.request.RequestDigest, "request_field_digests_json": fields}})
	if ex.response.Complete {
		for i := range prepared.Effects {
			if prepared.Effects[i].Table == "check_requests" {
				prepared.Effects[i].Values["response_json"] = encoded(ex.response)
			}
		}
	}
	_, err = journal.CommitEvent(db.lock, db.conn, prepared)
	return err
}

func (ex *execution) finish(db store, record checkrecord.Record, complete bool) error {
	effects := []journal.Effect{{Op: "insert", Table: "check_results", PK: map[string]any{"check_id": record.Start.ID}, Values: map[string]any{"result_json": encoded(record.Result)}}}
	if complete {
		effects = append(effects, journal.Effect{Op: "update", Table: "check_requests", PK: map[string]any{"id": ex.requestID}, Values: map[string]any{"response_json": encoded(ex.response)}, Preconditions: map[string]any{"response_json": nil}})
	}
	event, err := db.event("check.result-recorded", record.Start.ID, effects)
	if err != nil {
		return err
	}
	return ex.commitProgress(db, event)
}

func (ex *execution) keepFiles(db store, start checkrecord.Start, result *checkrecord.Result) error {
	cfg, err := config.ReadProjectConfig(filepath.Join(config.LedgerDir(ex.options.WorkspaceRoot), config.ConfigFileName))
	if err != nil {
		return err
	}
	ledger, err := filepath.EvalSymlinks(config.LedgerDir(ex.options.WorkspaceRoot))
	if err != nil {
		return err
	}
	storeRoot := filepath.Join(ledger, cfg.ArtifactRoot)
	if err := os.MkdirAll(storeRoot, 0o700); err != nil {
		return err
	}
	storeRoot, err = filepath.EvalSymlinks(storeRoot)
	if err != nil || !artifactstore.Contained(ledger, storeRoot) {
		return fmt.Errorf("artifact store is outside the ledger")
	}
	for _, path := range start.Definition.Artifacts {
		file, reason, err := artifactstore.Keep(ex.options.WorkspaceRoot, storeRoot, path, checkrecord.ArtifactLimit)
		if err != nil {
			return err
		}
		family := "check.artifact-refused"
		effects := []journal.Effect{}
		values := map[string]any{"artifact_id": nil, "reason": reason}
		if reason == "" {
			family = "check.artifact-kept"
			var id, availability string
			err := db.conn.QueryRow(`SELECT id,availability FROM referenced_artifacts WHERE digest=?`, file.Digest).Scan(&id, &availability)
			if err != nil && err != sql.ErrNoRows {
				return err
			}
			if id == "" {
				id, err = journal.GenerateID("ART", func(candidate string) (bool, error) {
					var n int
					err := db.conn.QueryRow(`SELECT count(*) FROM referenced_artifacts WHERE id=?`, candidate).Scan(&n)
					return n > 0, err
				})
				if err != nil {
					return err
				}
				media := mime.TypeByExtension(filepath.Ext(path))
				if media == "" {
					media = "application/octet-stream"
				}
				effects = append(effects, journal.Effect{Op: "insert", Table: "referenced_artifacts", PK: map[string]any{"id": id}, Values: map[string]any{"digest": file.Digest, "media_type": media, "size_bytes": file.SizeBytes, "retention_class": "evidence", "availability": "available", "location_policy": "content-addressed", "created_at": db.now.Format(time.RFC3339Nano), "expires_at": nil}})
			}
			if availability != "" && availability != "available" {
				effects = append(effects, journal.Effect{Op: "update", Table: "referenced_artifacts", PK: map[string]any{"id": id}, Values: map[string]any{"availability": "available"}, Preconditions: map[string]any{"availability": availability}})
			}
			values["artifact_id"], values["reason"] = id, nil
			result.Artifacts = append(result.Artifacts, checkrecord.Artifact{ID: id, Path: path, Digest: file.Digest, SizeBytes: file.SizeBytes})
		} else {
			result.RefusedArtifacts = append(result.RefusedArtifacts, checkrecord.RefusedArtifact{Path: path, Reason: reason})
		}
		effects = append(effects, journal.Effect{Op: "insert", Table: "check_files", PK: map[string]any{"check_id": start.ID, "path": path}, Values: values})
		event, err := db.event(family, start.ID, effects)
		if err != nil {
			return err
		}
		if _, err := journal.CommitEvent(db.lock, db.conn, event); err != nil {
			return err
		}
	}
	return nil
}

func summary(record checkrecord.Record) checkrecord.Summary {
	result := record.Result
	ids := []string{}
	for _, file := range result.Artifacts {
		ids = append(ids, file.ID)
	}
	return checkrecord.Summary{ID: record.Start.ID, Name: record.Start.Name, Outcome: result.Outcome, ExitCode: result.ExitCode, Counts: result.Counts, Reasons: result.Reasons, ArtifactIDs: ids}
}

func aggregateFailure(response checkrecord.Response) *checkrecord.Refusal {
	return checkrecord.FailureFor(response)
}
