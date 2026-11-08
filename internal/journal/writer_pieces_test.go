package journal

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/gittest"
)

// An import installs the journal as numbered pieces, whose names carry no
// date. The journal's latest date must still cover the Events in them, so an
// append never lands before the last Event, wherever that Event is.
func TestAppendEventRefusesADateBeforeTheLastEventInNumberedPieces(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name   string
		files  func(t *testing.T) map[string]string
		latest string
	}{
		{"one piece", func(t *testing.T) map[string]string {
			return map[string]string{"00000001.jsonl": pieceLine(t, "2026-06-05", 0) + pieceLine(t, "2026-06-07", 0)}
		}, "2026-06-07"},
		{"the last piece holds the last Event", func(t *testing.T) map[string]string {
			return map[string]string{"00000001.jsonl": pieceLine(t, "2026-06-01", 0), "00000002.jsonl": pieceLine(t, "2026-06-07", 0)}
		}, "2026-06-07"},
		{"blank lines after the last Event", func(t *testing.T) map[string]string {
			return map[string]string{"00000001.jsonl": pieceLine(t, "2026-06-07", 0) + "\n \n"}
		}, "2026-06-07"},
		{"an empty last piece", func(t *testing.T) map[string]string {
			return map[string]string{"00000001.jsonl": pieceLine(t, "2026-06-07", 0), "00000002.jsonl": ""}
		}, "2026-06-07"},
		{"a last Event longer than one read", func(t *testing.T) map[string]string {
			return map[string]string{"00000001.jsonl": pieceLine(t, "2026-06-05", 0) + pieceLine(t, "2026-06-07", 200<<10)}
		}, "2026-06-07"},
		{"a dated file after the pieces", func(t *testing.T) map[string]string {
			return map[string]string{"00000001.jsonl": pieceLine(t, "2026-06-05", 0), "2026-06-07.jsonl": pieceLine(t, "2026-06-07", 0)}
		}, "2026-06-07"},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			workspace := gittest.NewWorkspace(t)
			journalDir := filepath.Join(config.LedgerDir(workspace), "journal")
			mustMkdirAll(t, journalDir)
			for name, contents := range test.files(t) {
				if err := os.WriteFile(filepath.Join(journalDir, name), []byte(contents), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			latest, err := time.Parse("2006-01-02", test.latest)
			if err != nil {
				t.Fatal(err)
			}

			_, err = appendTestEvent(workspace, latest.Add(-time.Second), "Before the last Event")
			if want := "system clock is behind latest journal date " + test.latest + "; refusing to append out of order"; err == nil || err.Error() != want {
				t.Fatalf("append before %s = %v, want %q", test.latest, err, want)
			}
			if _, err := os.Stat(filepath.Join(journalDir, latest.Add(-time.Second).Format("2006-01-02")+".jsonl")); !os.IsNotExist(err) {
				t.Fatalf("refused append made a journal file: %v", err)
			}
			if _, err := appendTestEvent(workspace, latest, "On the last Event's date"); err != nil {
				t.Fatalf("append on %s = %v, want it appended", test.latest, err)
			}
		})
	}
}

// pieceLine is one journal line: an Event on date, its title padded by
// padding bytes.
func pieceLine(t *testing.T, date string, padding int) string {
	t.Helper()
	at, err := time.Parse("2006-01-02", date)
	if err != nil {
		t.Fatal(err)
	}
	event := testAppendEvent("Imported" + strings.Repeat("x", padding))
	event.OccurredAt = at.Add(12 * time.Hour)
	if event.EventID, err = FormatEventID(event.OccurredAt, 1); err != nil {
		t.Fatal(err)
	}
	payload, err := MarshalCanonicalJSON(event)
	if err != nil {
		t.Fatal(err)
	}
	return string(payload) + "\n"
}
