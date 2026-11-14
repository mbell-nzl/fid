package main

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestP24KeyedLeafCorpusMatchesLiveTree(t *testing.T) {
	t.Parallel()
	data, err := os.ReadFile(filepath.Join("integration", "testdata", "p24_keyed_leaves.tsv"))
	if err != nil {
		t.Fatal(err)
	}
	rows := map[string]string{}
	families := map[string]bool{}
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		parts := strings.Split(line, "\t")
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			t.Fatalf("invalid keyed-leaf row %q", line)
		}
		if _, exists := rows[parts[0]]; exists {
			t.Fatalf("duplicate keyed-leaf row %q", parts[0])
		}
		rows[parts[0]] = parts[1]
		families[strings.Fields(parts[0])[0]] = true
	}
	if len(rows) != 16 || len(families) != 10 {
		t.Fatalf("keyed-leaf corpus has %d rows in %d families; want 16/10", len(rows), len(families))
	}
	root := newRootCommandAt(t.TempDir())
	got := make([]string, 0, len(rows))
	var walk func(*cobra.Command)
	walk = func(command *cobra.Command) {
		if command != root && (command.Run != nil || command.RunE != nil) && command.Flags().Lookup("idempotency-key") != nil {
			name := rootCommandKey(command)
			if name != "export import" {
				got = append(got, name)
				field, exists := rows[name]
				if !exists {
					t.Errorf("live keyed leaf %s has no one-field conflict row", name)
				} else if strings.HasPrefix(field, "--") && command.Flags().Lookup(strings.TrimPrefix(field, "--")) == nil {
					t.Errorf("%s conflict field %s is not a leaf flag", name, field)
				}
			}
		}
		for _, child := range command.Commands() {
			walk(child)
		}
	}
	walk(root)
	want := make([]string, 0, len(rows))
	for name := range rows {
		want = append(want, name)
	}
	sort.Strings(got)
	sort.Strings(want)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("live keyed leaves = %v; corpus rows = %v", got, want)
	}
}
