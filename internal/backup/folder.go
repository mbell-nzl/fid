package backup

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/journal"
)

// genesisHash returns the hash of the ledger's genesis Event, which names
// the ledger in its folder's manifest. It is the journal's first Event,
// the project.initialized Event fid init appends, and stays first because
// the journal never appends before its latest date. Every fid init appends
// its own, and two differ in their time, to the nanosecond, even for one
// project. An import keeps the journal byte for byte, so a restored ledger
// keeps its genesis.
func genesisHash(workspaceRoot string) (string, error) {
	directory := filepath.Join(config.LedgerDir(workspaceRoot), config.JournalDirName)
	files, err := os.ReadDir(directory)
	if err != nil {
		return "", fmt.Errorf("read the journal: %w", err)
	}
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".jsonl") {
			continue
		}
		line, err := firstLine(filepath.Join(directory, file.Name()))
		if err != nil {
			return "", err
		}
		if len(line) == 0 {
			continue
		}
		var event journal.Event
		if err := event.UnmarshalJSON(line); err != nil {
			return "", fmt.Errorf("read the journal's first Event: %w", err)
		}
		return journal.HashEvent(event)
	}
	return "", errors.New("the journal has no Events")
}

func firstLine(path string) ([]byte, error) {
	input, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open the journal: %w", err)
	}
	defer func() { _ = input.Close() }()
	line, err := bufio.NewReader(input).ReadBytes('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("read the journal: %w", err)
	}
	return bytes.TrimSpace(line), nil
}

// checkFolder refuses to replace the project's folder at parent unless this
// ledger wrote it: its manifest must name this ledger's genesis Event.
// Clones of one project each make a ledger with the project's ID, as does
// fid init on a new machine, and a backup of one would otherwise replace
// the other's, which may be its only copy. A folder that names no ledger,
// or can't be read as one, is refused too.
func (run *attempt) checkFolder(repository clone, parent string) error {
	if parent == "" {
		return nil
	}
	listed, err := repository.text("ls-tree", "-z", parent, "--", run.projectID)
	if err != nil {
		return failure(KindInternal, err, "read the backup repository's tree")
	}
	if listed == "" {
		return nil
	}
	var written struct {
		GenesisEventHash string `json:"genesis_event_hash"`
	}
	detail := fmt.Sprintf("%s/%s does not name the ledger that wrote it; this ledger's genesis Event is %s", run.projectID, ManifestName, run.genesis)
	if payload, err := repository.run(nil, "cat-file", "blob", parent+":"+run.projectID+"/"+ManifestName); err == nil &&
		json.Unmarshal([]byte(payload), &written) == nil && written.GenesisEventHash != "" {
		if written.GenesisEventHash == run.genesis {
			return nil
		}
		detail = fmt.Sprintf("%s/%s names the ledger whose genesis Event is %s; this ledger's is %s", run.projectID, ManifestName, written.GenesisEventHash, run.genesis)
	}
	return &Error{Kind: KindDifferentLedger, Detail: detail, Message: fmt.Sprintf(
		"the backup repository's %s/ holds another ledger's backup, so fid won't replace it; if that backup is no longer needed, remove %s/ from the backup repository's %s",
		run.projectID, run.projectID, Branch)}
}
