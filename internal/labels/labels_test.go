package labels

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/bootstrap"
	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestAddStoresOptionalFieldsAndRejectsDuplicateBeforeAppend(t *testing.T) {
	t.Parallel()

	workspace := newLabelsWorkspace(t)
	colour := "#4A90D9"
	description := "Core infrastructure"

	beforeEvents := countLabelJournalEvents(t, workspace)

	result, err := Add(AddOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testLabelsMigrationsDir(t),
		Name:          "core",
		Colour:        &colour,
		Description:   &description,
		Actor:         testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
	})
	if err != nil {
		t.Fatalf("Add() error = %v", err)
	}
	if got := countLabelJournalEvents(t, workspace); got != beforeEvents+1 {
		t.Fatalf("journal events after Add() = %d, want %d", got, beforeEvents+1)
	}
	if result.Label.Name != "core" {
		t.Fatalf("result label name = %q, want %q", result.Label.Name, "core")
	}
	if dereferenceTestString(result.Label.Colour) != colour {
		t.Fatalf("result label colour = %q, want %q", dereferenceTestString(result.Label.Colour), colour)
	}
	if dereferenceTestString(result.Label.Description) != description {
		t.Fatalf("result label description = %q, want %q", dereferenceTestString(result.Label.Description), description)
	}
	if result.EventID == "" || result.OccurredAt == "" {
		t.Fatalf("result Event identity = (%q, %q), want populated", result.EventID, result.OccurredAt)
	}

	conn := openLabelsDB(t, workspace)
	defer closeLabelsDB(t, conn)

	var (
		name              string
		storedColour      sql.NullString
		storedDescription sql.NullString
	)
	if err := conn.QueryRow(`SELECT name, colour, description FROM labels WHERE name = ?`, "core").Scan(&name, &storedColour, &storedDescription); err != nil {
		t.Fatalf("query label row: %v", err)
	}
	if name != "core" {
		t.Fatalf("stored label name = %q, want %q", name, "core")
	}
	if storedColour.String != colour {
		t.Fatalf("stored colour = %q, want %q", storedColour.String, colour)
	}
	if storedDescription.String != description {
		t.Fatalf("stored description = %q, want %q", storedDescription.String, description)
	}

	duplicateBefore := countLabelJournalEvents(t, workspace)
	_, err = Add(AddOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testLabelsMigrationsDir(t),
		Name:          "core",
		Actor:         testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
	})
	if err == nil {
		t.Fatal("duplicate Add() error = nil, want rejection")
	}
	var duplicateErr *DuplicateLabelError
	if !errors.As(err, &duplicateErr) {
		t.Fatalf("duplicate Add() error = %T, want *DuplicateLabelError", err)
	}
	if duplicateErr.Name != "core" {
		t.Fatalf("duplicate label name = %q, want %q", duplicateErr.Name, "core")
	}
	if got := countLabelJournalEvents(t, workspace); got != duplicateBefore {
		t.Fatalf("journal events after duplicate Add() = %d, want %d", got, duplicateBefore)
	}
}

func TestAddAppliesClosedLabelValueBoundariesBeforeAppend(t *testing.T) {
	t.Parallel()

	t.Run("accepts exact maxima", func(t *testing.T) {
		workspace := newLabelsWorkspace(t)
		colour := strings.Repeat("c", 64)
		description := strings.Repeat("d", 4096)
		result, err := Add(AddOptions{WorkspaceRoot: workspace, MigrationsDir: testLabelsMigrationsDir(t), Name: strings.Repeat("n", 255), Colour: &colour, Description: &description, Actor: testActorForRole("developer"), IdempotencyKey: testIdempotencyKey()})
		if err != nil || len(result.Label.Name) != 255 || len(*result.Label.Colour) != 64 || len(*result.Label.Description) != 4096 {
			t.Fatalf("Add() = (%#v, %v), want exact maxima", result, err)
		}
	})

	validName := "core"
	tests := []struct {
		name        string
		label       string
		colour      *string
		description *string
		wantErr     string
	}{
		{name: "name over maximum", label: strings.Repeat("n", 256), wantErr: "label name exceeds 255 bytes"},
		{name: "name control character", label: "bad\x00name", wantErr: "label name contains ASCII control character"},
		{name: "name invalid UTF-8", label: string([]byte{0xff}), wantErr: "label name must be valid UTF-8"},
		{name: "colour explicitly empty", label: validName, colour: testStringPointer(" "), wantErr: "label colour is required"},
		{name: "colour over maximum", label: validName, colour: testStringPointer(strings.Repeat("c", 65)), wantErr: "label colour exceeds 64 bytes"},
		{name: "colour control character", label: validName, colour: testStringPointer("red\x1f"), wantErr: "label colour contains ASCII control character"},
		{name: "description explicitly empty", label: validName, description: testStringPointer(" "), wantErr: "label description is required"},
		{name: "description over maximum", label: validName, description: testStringPointer(strings.Repeat("d", 4097)), wantErr: "label description exceeds 4096 bytes"},
		{name: "description invalid UTF-8", label: validName, description: testStringPointer(string([]byte{0xff})), wantErr: "label description must be valid UTF-8"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			workspace := newLabelsWorkspace(t)
			before := countLabelJournalEvents(t, workspace)
			_, err := Add(AddOptions{WorkspaceRoot: workspace, MigrationsDir: testLabelsMigrationsDir(t), Name: test.label, Colour: test.colour, Description: test.description, Actor: testActorForRole("developer"), IdempotencyKey: testIdempotencyKey()})
			if err == nil || err.Error() != test.wantErr {
				t.Fatalf("Add() error = %v, want %q", err, test.wantErr)
			}
			if got := countLabelJournalEvents(t, workspace); got != before {
				t.Fatalf("journal events after rejection = %d, want %d", got, before)
			}
		})
	}
}

func TestListRejectsOverTopLevelCollectionLimit(t *testing.T) {
	t.Parallel()
	workspace := newLabelsWorkspace(t)
	conn := openLabelsDB(t, workspace)
	tx, err := conn.Begin()
	if err != nil {
		t.Fatalf("Begin() error = %v", err)
	}
	for index := 0; index < 10000; index++ {
		if _, err := tx.Exec(`INSERT INTO labels (name) VALUES (?)`, fmt.Sprintf("label-%05d", index)); err != nil {
			t.Fatalf("insert label %d: %v", index, err)
		}
	}
	if err := tx.Commit(); err != nil {
		t.Fatalf("Commit() error = %v", err)
	}
	closeLabelsDB(t, conn)
	atLimit, err := List(ListOptions{WorkspaceRoot: workspace, MigrationsDir: testLabelsMigrationsDir(t)})
	if err != nil || len(atLimit) != 10000 {
		t.Fatalf("List(at limit) = (%d, %v), want (10000, nil)", len(atLimit), err)
	}

	conn = openLabelsDB(t, workspace)
	if _, err := conn.Exec(`INSERT INTO labels (name) VALUES (?)`, "label-10000"); err != nil {
		t.Fatalf("insert over-limit label: %v", err)
	}
	closeLabelsDB(t, conn)

	_, err = List(ListOptions{WorkspaceRoot: workspace, MigrationsDir: testLabelsMigrationsDir(t)})
	var limit *ResultLimitExceededError
	if !errors.As(err, &limit) || limit.Collection != "labels" || limit.Limit != 10000 || limit.Observed != 10001 {
		t.Fatalf("List() error = %#v, want labels 10001/10000 limit", err)
	}
}

func TestListReturnsLabelsInNameOrderWithNilOptionalFields(t *testing.T) {
	t.Parallel()

	workspace := newLabelsWorkspace(t)
	description := "Journal subsystem"
	if _, err := Add(AddOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testLabelsMigrationsDir(t),
		Name:          "journal",
		Description:   &description,
		Actor:         testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
	}); err != nil {
		t.Fatalf("Add(journal) error = %v", err)
	}
	if _, err := Add(AddOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testLabelsMigrationsDir(t),
		Name:          "core",
		Actor:         testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
	}); err != nil {
		t.Fatalf("Add(core) error = %v", err)
	}

	got, err := List(ListOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testLabelsMigrationsDir(t),
	})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	want := []Label{
		{Name: "core"},
		{Name: "journal", Description: &description},
	}
	if len(got) != len(want) {
		t.Fatalf("len(List()) = %d, want %d (%#v)", len(got), len(want), got)
	}
	for index := range want {
		if !equalLabels(got[index], want[index]) {
			t.Fatalf("List()[%d] = %#v, want %#v", index, got[index], want[index])
		}
	}
}

func TestListReturnsEmptySliceForWorkspaceWithNoLabels(t *testing.T) {
	t.Parallel()

	workspace := newLabelsWorkspace(t)

	got, err := List(ListOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testLabelsMigrationsDir(t),
	})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("List() = %#v, want empty slice", got)
	}
}

func TestAddAndListValidateRequiredInputs(t *testing.T) {
	t.Parallel()

	workspace := newLabelsWorkspace(t)
	migrationsDir := testLabelsMigrationsDir(t)

	tests := []struct {
		name string
		run  func() error
		want string
	}{
		{
			name: "add missing workspace",
			run: func() error {
				_, err := Add(AddOptions{
					MigrationsDir: migrationsDir,
					Name:          "core",
					Actor:         testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
				})
				return err
			},
			want: "label add: workspace root is required",
		},
		{
			name: "add missing migrations dir",
			run: func() error {
				_, err := Add(AddOptions{
					WorkspaceRoot: workspace,
					Name:          "core",
					Actor:         testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
				})
				return err
			},
			want: "label add: migrations dir is required",
		},
		{
			name: "add missing label name",
			run: func() error {
				_, err := Add(AddOptions{
					WorkspaceRoot: workspace,
					MigrationsDir: migrationsDir,
					Name:          "   ",
					Actor:         testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
				})
				return err
			},
			want: "label add: label name is required",
		},
		{
			name: "list missing workspace",
			run: func() error {
				_, err := List(ListOptions{
					MigrationsDir: migrationsDir,
				})
				return err
			},
			want: "label list: workspace root is required",
		},
		{
			name: "list missing migrations dir",
			run: func() error {
				_, err := List(ListOptions{
					WorkspaceRoot: workspace,
				})
				return err
			},
			want: "label list: migrations dir is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.run()
			if err == nil {
				t.Fatal("error = nil, want validation error")
			}
			if err.Error() != tt.want {
				t.Fatalf("error = %q, want %q", err.Error(), tt.want)
			}
		})
	}
}

func newLabelsWorkspace(t *testing.T) string {
	t.Helper()

	workspace := t.TempDir()
	if _, err := bootstrap.InitializeWorkspace(bootstrap.InitOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testLabelsMigrationsDir(t),
	}); err != nil {
		t.Fatalf("InitializeWorkspace() error = %v", err)
	}
	return workspace
}

func testLabelsMigrationsDir(t *testing.T) string {
	t.Helper()

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller(0) failed")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", "migrations"))
}

func countLabelJournalEvents(t *testing.T, workspace string) int {
	t.Helper()

	entries, err := os.ReadDir(filepath.Join(workspace, ".fid", "journal"))
	if err != nil {
		t.Fatalf("ReadDir(journal) error = %v", err)
	}
	total := 0
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".jsonl") {
			continue
		}
		payload, err := os.ReadFile(filepath.Join(workspace, ".fid", "journal", entry.Name()))
		if err != nil {
			t.Fatalf("ReadFile(%q) error = %v", entry.Name(), err)
		}
		trimmed := strings.TrimSpace(string(payload))
		if trimmed == "" {
			continue
		}
		total += len(strings.Split(trimmed, "\n"))
	}
	return total
}

func openLabelsDB(t *testing.T, workspace string) *sql.DB {
	t.Helper()

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	return conn
}

func closeLabelsDB(t *testing.T, conn *sql.DB) {
	t.Helper()
	if err := conn.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
}

func dereferenceTestString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func testStringPointer(value string) *string { return &value }

func equalLabels(left, right Label) bool {
	return left.Name == right.Name &&
		dereferenceTestString(left.Colour) == dereferenceTestString(right.Colour) &&
		dereferenceTestString(left.Description) == dereferenceTestString(right.Description)
}
