package backup

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/export"
	"github.com/mbell-nzl/fid/internal/gittest"
)

func TestManifestWriterStopsAtTheLimit(t *testing.T) {
	t.Parallel()
	for _, delta := range []int{-1, 0, 1} {
		t.Run(fmt.Sprint(delta), func(t *testing.T) {
			var output bytes.Buffer
			writer := &folderBlobWriter{output: &output, remaining: 16}
			if _, err := writer.Write([]byte(strings.Repeat("x", 15))); err != nil {
				t.Fatal(err)
			}
			n, err := writer.Write([]byte(strings.Repeat("x", 1+delta)))
			if delta <= 0 && (err != nil || output.Len() != 16+delta) {
				t.Fatalf("write = %d, %v, output %d bytes", n, err, output.Len())
			}
			if delta > 0 && (!errors.Is(err, errFolderResourceLimit) || n != 0 || output.Len() != 15 || !writer.exceeded) {
				t.Fatalf("write = %d, %v, output %d bytes, want bounded refusal", n, err, output.Len())
			}
		})
	}
}

func TestFolderInventoryStreamsMetadataAndUses64BitAggregate(t *testing.T) {
	t.Parallel()
	limits := export.ArchiveResourceLimits()
	inventory := &folderInventory{id: "project", limits: limits}
	for i := range 8 {
		// Each 1 GiB Artifact fits exactly. Split records across writes,
		// like the process pipe can, without constructing any large blob.
		record := fmt.Sprintf("100644 blob %s %d\tproject/artifacts/%064d\x00", strings.Repeat("a", 40), limits.ArtifactBytes, i)
		if _, err := inventory.Write([]byte(record[:13])); err != nil {
			t.Fatal(err)
		}
		if _, err := inventory.Write([]byte(record[13:])); err != nil {
			t.Fatal(err)
		}
	}
	if inventory.total != 8<<30 || len(inventory.members) != 8 {
		t.Fatalf("inventory = %d bytes, %d members, want exactly 8 GiB", inventory.total, len(inventory.members))
	}
	// The first byte over the aggregate ceiling stops the stream before
	// the additional tail can be retained as another pending record.
	record := "100644 blob " + strings.Repeat("b", 40) + " 1\tproject/manifest.json\x00" + strings.Repeat("x", 4096)
	_, err := inventory.Write([]byte(record))
	var invalid *export.InvalidArchiveError
	if !errors.As(err, &invalid) || invalid.Checks[0].Code != "resource-limit" || len(inventory.members) != 8 || inventory.pending.Len() >= 4096 {
		t.Fatalf("overflow = %v, members %d, pending %d bytes", err, len(inventory.members), inventory.pending.Len())
	}
}

func TestReconstructionChecksResourceLimitsBeforeCreatingArchive(t *testing.T) {
	t.Parallel()
	limits := export.ResourceLimits{NonArtifactBytes: 16, ArtifactBytes: 32, ExpandedBytes: 48, Entries: 3}
	for _, kind := range []string{"non-Artifact", "Artifact", "aggregate", "entries"} {
		for _, delta := range []int{-1, 0, 1} {
			t.Run(fmt.Sprintf("%s/%+d", kind, delta), func(t *testing.T) {
				t.Parallel()
				files := map[string]int{ManifestName: 10}
				switch kind {
				case "non-Artifact":
					files["journal/00000001.jsonl"] = 16 + delta
				case "Artifact":
					files["artifacts/"+strings.Repeat("a", 64)] = 32 + delta
				case "aggregate":
					files["journal/00000001.jsonl"] = 16
					files["artifacts/"+strings.Repeat("a", 64)] = 31 + delta
					files["manifest.json"] = 1
				case "entries":
					for i := range 3 + delta {
						files[fmt.Sprintf("journal/%08d.jsonl", i+1)] = 0
					}
				}
				repository, tree := restoreLimitFolder(t, files)
				scratch := t.TempDir()
				_, err := zipFolderWithLimits(repository, tree, "project", "mismatched", scratch, limits)
				var invalid *export.InvalidArchiveError
				if !errors.As(err, &invalid) || len(invalid.Checks) != 1 {
					t.Fatalf("reconstruction = %v, want one archive check", err)
				}
				want := CheckArchiveDigest
				if delta > 0 {
					want = "resource-limit"
				}
				if invalid.Checks[0].Code != want {
					t.Fatalf("check = %+v, want %s", invalid.Checks[0], want)
				}
				_, statErr := os.Stat(filepath.Join(scratch, ArchiveName))
				if delta > 0 && !os.IsNotExist(statErr) || delta <= 0 && statErr != nil {
					t.Fatalf("archive stat = %v, want creation only within the limits", statErr)
				}
			})
		}
	}
}

func restoreLimitFolder(t *testing.T, files map[string]int) (clone, string) {
	t.Helper()
	root := gittest.NewWorkspace(t)
	repository := clone{gitDir: filepath.Join(root, ".git"), workingDir: root,
		environment: gitEnvironment(append(os.Environ(), "HOME="+t.TempDir(), "XDG_CONFIG_HOME="+t.TempDir()))}
	tree := &memberTree{}
	for name, size := range files {
		blob, err := repository.output(strings.NewReader(strings.Repeat("x", size)), "hash-object", "-w", "--stdin")
		if err != nil {
			t.Fatal(err)
		}
		tree.add("project/"+name, blob)
	}
	oid, err := tree.write(repository)
	if err != nil {
		t.Fatal(err)
	}
	return repository, oid
}
