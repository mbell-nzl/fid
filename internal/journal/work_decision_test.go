package journal

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
)

func workDecisionEvent(allowance bool) Event {
	e := admissionEvent("EVT-20261110-000004", "work.accepted", "issue", "FID-a", 2, Effect{Op: "insert", Table: "person_decisions", PK: map[string]any{"id": "PDN-work"}, Values: map[string]any{"kind": "accept", "actor": "human:mark", "subject_type": "task", "subject_id": "FID-a", "binding_id": "FID-a", "reason": "Approved work", "decided_at": admissionStamp(60)}})
	e.CommandName = "task.accept"
	if allowance {
		e.EventType, e.CommandName, e.Aggregate.Type, e.Aggregate.Sequence = "rework.allowed", "task.rework-allow", "task", 1
		e.Effects[0].Values["kind"], e.Effects[0].Values["binding_json"] = "rework-allowance", `{"rounds":1}`
	}
	return e
}

func TestWorkDecisionsAcrossJournalPaths(t *testing.T) {
	t.Parallel()
	for _, allowance := range []bool{false, true} {
		e := workDecisionEvent(allowance)
		t.Run(e.EventType, func(t *testing.T) {
			t.Parallel()
			for _, mode := range admissionModes(e) {
				t.Run("valid/"+mode, func(t *testing.T) {
					w := newAdmissionWorkspace(t, admissionPrefix(), mode)
					if err := w.admit(t, mode, e); err != nil {
						t.Fatal(err)
					}
				})
			}
			for _, name := range []string{"reason", "binding", "subject", "count", "actor"} {
				t.Run("forged/"+name, func(t *testing.T) {
					t.Parallel()
					forged := workDecisionEvent(allowance)
					want := "work decision"
					switch name {
					case "reason":
						forged.Effects[0].Values["reason"] = "\t"
					case "binding":
						forged.Effects[0].Values["binding_id"] = "FID-b"
					case "subject":
						forged.Effects[0].Values["subject_type"] = "epic"
					case "count":
						extra := forged.Effects[0]
						extra.PK = map[string]any{"id": "PDN-extra"}
						forged.Effects = append(forged.Effects, extra)
					case "actor":
						forged.ActorV2.ActorType = "agent"
						want = "actor"
					}
					assertWorkDecisionRefused(t, forged, want)
				})
			}
		})
	}
}

// A full replay may retain the valid prefix (the existing journal contract),
// but the refused Event itself must change no domain row or checkpoint.
func assertWorkDecisionRefused(t *testing.T, e Event, reason string) {
	t.Helper()
	reference := newAdmissionWorkspace(t, admissionPrefix(), "commit")
	conn := openReplayDB(t, filepath.Join(config.LedgerDir(reference.root), "fid.sqlite"))
	prefix := currentProjectionSnapshot(t, conn)
	closeProjectDB(t, conn)
	for _, mode := range admissionModes(e) {
		t.Run(mode, func(t *testing.T) {
			w := newAdmissionWorkspace(t, admissionPrefix(), mode)
			beforeConn := openReplayDB(t, filepath.Join(config.LedgerDir(w.root), "fid.sqlite"))
			before := currentProjectionSnapshot(t, beforeConn)
			closeProjectDB(t, beforeConn)
			err := w.admit(t, mode, e)
			if err == nil || !strings.Contains(err.Error(), reason) {
				t.Fatalf("refusal=%v, want %s", err, reason)
			}
			conn := openReplayDB(t, filepath.Join(config.LedgerDir(w.root), "fid.sqlite"))
			after := currentProjectionSnapshot(t, conn)
			closeProjectDB(t, conn)
			expected := prefix
			if mode == "replay" && reason == "actor" {
				expected = before // shape refusal happens in the pre-scan
			}
			if !reflect.DeepEqual(expected, after) {
				t.Fatal("rejected Event changed the valid prefix")
			}
			stored, err := os.ReadFile(w.file)
			if err != nil {
				t.Fatal(err)
			}
			live := mode == "commit" || mode == "dry-run" || mode == "project"
			want := w.journalWith(t, e)
			if live {
				want = w.prefix
			}
			if !bytes.Equal(stored, want) {
				t.Fatal("refusal changed journal bytes")
			}
			if live {
				if err := w.admit(t, "commit", workDecisionEvent(e.EventType == "rework.allowed")); err != nil {
					t.Fatal("valid write after refusal", err)
				}
			}
		})
	}
}

// Legacy identity edits remain admitted until 1.9c owns row identity. The
// derived creator follows that identity, without changing admission timing.
func TestCreatorIndexFollowsLegacyIssueIdentity(t *testing.T) {
	t.Parallel()
	for _, operation := range []string{"rename", "delete"} {
		t.Run(operation, func(t *testing.T) {
			t.Parallel()
			effect := Effect{Op: "delete", Table: "issues", PK: map[string]any{"id": "FID-b"}}
			if operation == "rename" {
				effect.Op = "update"
				effect.Values = map[string]any{"id": "FID-renamed"}
			}
			e := stateLegacy(admissionEvent("EVT-20261110-000004", "issue.updated", "issue", "FID-b", 2, effect))
			for _, mode := range admissionModes(e) {
				t.Run(mode, func(t *testing.T) {
					w := newAdmissionWorkspace(t, admissionPrefix(), mode)
					if err := w.admit(t, mode, e); err != nil {
						t.Fatal("creator index changed legacy admission", err)
					}
					conn := openReplayDB(t, filepath.Join(config.LedgerDir(w.root), "fid.sqlite"))
					defer closeProjectDB(t, conn)
					id := "FID-renamed"
					want := 1
					if operation == "delete" {
						id = "FID-b"
						want = 0
					}
					if mode == "dry-run" {
						id = "FID-b"
						want = 1
					}
					var count int
					if err := conn.QueryRow(`SELECT count(*) FROM issue_creators WHERE issue_id=? AND actor='human:mark'`, id).Scan(&count); err != nil || count != want {
						t.Fatalf("creator after %s = %d/%v, want %d", mode, count, err, want)
					}
				})
			}
		})
	}
}
