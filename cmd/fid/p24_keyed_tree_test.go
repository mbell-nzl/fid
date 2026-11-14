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

func TestP24EveryLiveKeyedLeafHasExecutableCorpusRow(t *testing.T) {
	t.Parallel()
	root := newRootCommandAt(t.TempDir())
	var live []string
	var walk func(*cobra.Command)
	walk = func(command *cobra.Command) {
		if len(command.Commands()) == 0 && command.Flags().Lookup("idempotency-key") != nil {
			name := strings.TrimPrefix(command.CommandPath(), root.Name()+" ")
			if name != "export import" {
				live = append(live, name)
			}
		}
		for _, child := range command.Commands() {
			walk(child)
		}
	}
	walk(root)
	sort.Strings(live)
	data, err := os.ReadFile(filepath.Join("integration", "testdata", "p24_keyed_leaves.tsv"))
	if err != nil {
		t.Fatal(err)
	}
	var planned []string
	for _, row := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		planned = append(planned, strings.SplitN(row, "\t", 2)[0])
	}
	sort.Strings(planned)
	if len(live) != 51 || !reflect.DeepEqual(live, planned) {
		t.Fatalf("live keyed leaves = %q, executable corpus = %q", live, planned)
	}
}
