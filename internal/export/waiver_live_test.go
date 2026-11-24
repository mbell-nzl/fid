package export

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/reviewwork"
)

func TestWaiverMalformedAppendAndCatchupPreserveState(t *testing.T) {
	t.Parallel()
	source := newExportSourceFixture(t, addWaiverExportFacts)
	for _, mode := range []string{"append", "catchup"} {
		for _, field := range []string{"actor", "reason", "binding", "latest", "lane"} {
			t.Run(mode+"/"+field, func(t *testing.T) {
				t.Parallel()
				root, _ := source.clone(t)
				files, err := filepath.Glob(filepath.Join(config.LedgerDir(root), "journal", "*.jsonl"))
				if err != nil || len(files) == 0 {
					t.Fatal(files, err)
				}
				path := files[len(files)-1]
				prefix, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				lines := bytes.Split(bytes.TrimSpace(prefix), []byte("\n"))
				var valid journal.Event
				if err := json.Unmarshal(lines[len(lines)-1], &valid); err != nil {
					t.Fatal(err)
				}
				valid.Idempotency = nil
				valid.Effects = valid.Effects[:1]
				valid.Aggregate.ID = "PDN-next"
				valid.Effects[0].PK = map[string]any{"id": "PDN-next"}
				seq, err := strconv.Atoi(valid.EventID[strings.LastIndex(valid.EventID, "-")+1:])
				if err != nil {
					t.Fatal(err)
				}
				valid.EventID, err = journal.FormatEventID(valid.OccurredAt, seq+1)
				if err != nil {
					t.Fatal(err)
				}
				raw, err := json.Marshal(valid)
				if err != nil {
					t.Fatal(err)
				}
				var forged journal.Event
				if err := json.Unmarshal(raw, &forged); err != nil {
					t.Fatal(err)
				}
				switch field {
				case "actor":
					forged.ActorV2.ActorType = "agent"
					forged.Effects[0].Values["actor"] = "agent:mark"
				case "reason":
					forged.Effects[0].Values["reason"] = " "
				case "binding":
					delete(forged.Effects[0].Values, "binding_json")
				default:
					var binding map[string]any
					if err := json.Unmarshal([]byte(forged.Effects[0].Values["binding_json"].(string)), &binding); err != nil {
						t.Fatal(err)
					}
					if field == "latest" {
						binding["handoff_id"] = "HOF-other"
					} else {
						binding["lane"] = "assured"
					}
					forged.Effects[0].Values["binding_json"] = encodedLean(binding)
				}
				conn, err := reviewwork.OpenProjection(root)
				if err != nil {
					t.Fatal(err)
				}
				defer func() { _ = conn.Close() }()
				tables, err := loadSemanticTables(conn)
				if err != nil {
					t.Fatal(err)
				}
				v := &semanticValidator{conn: conn, tables: tables}
				before := handoffProjectionSnapshot(t, v)
				if mode == "append" {
					err = journal.WithLock(root, journal.LockOwner{Command: "test waive"}, func(lock *journal.WriteLock) error { _, err := journal.CommitEvent(lock, conn, forged); return err })
				} else {
					tail, marshalErr := journal.MarshalCanonicalJSON(forged)
					if marshalErr != nil {
						t.Fatal(marshalErr)
					}
					if err := os.WriteFile(path, append(append([]byte{}, prefix...), append(tail, '\n')...), 0o600); err != nil {
						t.Fatal(err)
					}
					_, err = journal.EnsureProjectionFresh(root, testMigrationsFS(t))
				}
				if err == nil || handoffProjectionSnapshot(t, v) != before {
					t.Fatalf("refusal changed state: %v", err)
				}
				stored, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				if mode == "append" && !bytes.Equal(stored, prefix) {
					t.Fatal("refusal appended journal bytes")
				}
				if mode == "catchup" {
					if err := os.WriteFile(path, prefix, 0o600); err != nil {
						t.Fatal(err)
					}
				}
				if err := journal.WithLock(root, journal.LockOwner{Command: "test valid waive"}, func(lock *journal.WriteLock) error { _, err := journal.CommitEvent(lock, conn, valid); return err }); err != nil {
					t.Fatal("valid write after refusal", err)
				}
			})
		}
	}
}
