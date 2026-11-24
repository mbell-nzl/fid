package journal

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestVisitCanonicalOrderAndStop(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"all", "stop", "callback-error", "corrupt"} {
		t.Run(mode, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			ids := []string{"EVT-first", "EVT-second", "EVT-third"}
			// Names determine replay order; the last file intentionally has no
			// final newline. Ignore files the replayer would ignore, too.
			writeWriterEvents(t, filepath.Join(dir, "2026-11-24.jsonl"), writerRoleEvent(ids[0]), writerRoleEvent(ids[1]))
			raw, err := json.Marshal(writerRoleEvent(ids[2]))
			if err != nil {
				t.Fatal(err)
			}
			if mode == "corrupt" {
				raw = []byte("{broken")
			}
			for name, data := range map[string][]byte{"2026-11-25.jsonl": append([]byte("\n"), raw...), "notes.jsonl": []byte("not journal history")} {
				if err := os.WriteFile(filepath.Join(dir, name), data, 0o600); err != nil {
					t.Fatal(err)
				}
			}
			callbackErr := errors.New("visitor failed")
			seen := []string{}
			err = Visit(dir, func(e Event) (bool, error) {
				seen = append(seen, e.EventID)
				if e.EventID == ids[1] {
					if mode == "callback-error" {
						return false, callbackErr
					}
					return mode == "stop", nil
				}
				return false, nil
			})
			want := ids[:2]
			if mode == "all" {
				want = ids
			}
			if !reflect.DeepEqual(seen, want) {
				t.Errorf("visited %v, want %v", seen, want)
			}
			var corrupt *CorruptJournalError
			switch mode {
			case "corrupt":
				if !errors.As(err, &corrupt) || corrupt.File != filepath.ToSlash(filepath.Join(dir, "2026-11-25.jsonl")) || corrupt.Line != 2 {
					t.Fatalf("corrupt location lost: %v", err)
				}
			case "callback-error":
				if !errors.Is(err, callbackErr) {
					t.Fatalf("visitor error lost: %v", err)
				}
			default:
				if err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}
