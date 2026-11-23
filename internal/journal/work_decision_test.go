package journal

import (
	"bytes"
	"errors"
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

// RESET 1.9c refuses legacy identity edits. Refusal preserves the creator
// index along with the Issue, every other projected row and the journal.
func TestCreatorIndexPreservesRefusedLegacyIdentity(t *testing.T) {
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
			reason := "an Issue cannot be deleted or recreated"
			if operation == "rename" {
				reason = "id is immutable"
			}
			for _, mode := range admissionModes(e) {
				t.Run(mode, func(t *testing.T) {
					w := newAdmissionWorkspace(t, admissionPrefix(), mode)
					reference := w
					if mode == "replay" {
						reference = newAdmissionWorkspace(t, admissionPrefix(), "commit")
					}
					before := stateProjection(t, reference.root)
					err := w.admit(t, mode, e)
					var refusal *RowIdentityError
					if !errors.As(err, &refusal) || refusal.Table != "issues" || refusal.Reason != reason {
						t.Fatalf("legacy identity edit = %v, want %s refusal", err, reason)
					}
					if !reflect.DeepEqual(before, stateProjection(t, w.root)) {
						t.Fatal("refused identity edit changed projection")
					}
					conn := openReplayDB(t, filepath.Join(config.LedgerDir(w.root), "fid.sqlite"))
					defer closeProjectDB(t, conn)
					var count int
					if err := conn.QueryRow(`SELECT count(*) FROM issue_creators WHERE issue_id='FID-b' AND actor='human:mark'`).Scan(&count); err != nil || count != 1 {
						t.Fatalf("creator after %s = %d/%v, want 1", mode, count, err)
					}
					live := mode == "dry-run" || mode == "project"
					want := w.journalWith(t, e)
					if live {
						want = w.prefix
					}
					stored, err := os.ReadFile(w.file)
					if err != nil || !bytes.Equal(stored, want) {
						t.Fatalf("refusal changed journal: %v", err)
					}
					if err := os.WriteFile(w.file, w.prefix, 0o600); err != nil {
						t.Fatal(err)
					}
					if err := w.admit(t, "commit", admissionLaneRaise()); err != nil {
						t.Fatalf("valid write after refusal: %v", err)
					}
				})
			}
		})
	}
}
