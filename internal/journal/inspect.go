package journal

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type FoundEvent struct {
	Event       Event
	JournalFile string
	Line        int
	Hash        string
}

type EventNotFoundError struct {
	EventID string
}

func (err *EventNotFoundError) Error() string {
	if err == nil {
		return "<nil>"
	}
	return fmt.Sprintf("journal event %q not found", err.EventID)
}

type Verdict struct {
	Eligible     bool           `json:"eligible"`
	Mode         string         `json:"mode,omitempty"`
	Restores     map[string]any `json:"restores,omitempty"`
	Precondition map[string]any `json:"precondition,omitempty"`
	Correction   string         `json:"correction,omitempty"`
	Note         string         `json:"note,omitempty"`
}

type journalInspection struct {
	checkpointSeen bool
	checkpoint     journalCheckpoint
	tail           journalTail
}

type journalCheckpoint struct {
	EventID   string
	EventHash string
}

type journalTail struct {
	HasEvents bool
	File      string
	Offset    int64
}

func FindEvent(journalDir, eventID string) (FoundEvent, error) {
	if strings.TrimSpace(journalDir) == "" {
		return FoundEvent{}, fmt.Errorf("find event: journal dir is required")
	}
	if strings.TrimSpace(eventID) == "" {
		return FoundEvent{}, fmt.Errorf("find event: event id is required")
	}

	files, err := orderedJournalFiles(journalDir)
	if err != nil {
		return FoundEvent{}, err
	}

	for _, fileName := range files {
		found, ok, err := findEventInFile(journalDir, fileName, eventID)
		if err != nil {
			return FoundEvent{}, err
		}
		if ok {
			return found, nil
		}
	}

	return FoundEvent{}, &EventNotFoundError{EventID: eventID}
}

func inspectJournal(journalDir, checkpointFile string, checkpointOffset int64) (journalInspection, error) {
	files, err := orderedJournalFiles(journalDir)
	if err != nil {
		return journalInspection{}, err
	}

	var inspection journalInspection
	for _, fileName := range files {
		if err := inspectJournalFile(journalDir, fileName, checkpointFile, checkpointOffset, &inspection); err != nil {
			return journalInspection{}, err
		}
	}
	return inspection, nil
}

func findEventInFile(journalDir, fileName, eventID string) (FoundEvent, bool, error) {
	path := filepath.Join(journalDir, fileName)
	file, err := os.Open(path)
	if err != nil {
		return FoundEvent{}, false, fmt.Errorf("open journal file %q: %w", path, err)
	}
	defer closeReplayFile(file)

	reader := bufio.NewReader(file)
	for lineNumber := 1; ; lineNumber++ {
		lineBytes, readErr := reader.ReadBytes('\n')
		if readErr != nil && !errors.Is(readErr, io.EOF) {
			return FoundEvent{}, false, fmt.Errorf("read journal file %q: %w", path, readErr)
		}
		if len(lineBytes) == 0 && errors.Is(readErr, io.EOF) {
			return FoundEvent{}, false, nil
		}

		found, ok, err := inspectJournalLine(fileName, lineNumber, lineBytes, eventID)
		if err != nil {
			return FoundEvent{}, false, err
		}
		if ok {
			return found, true, nil
		}

		if errors.Is(readErr, io.EOF) {
			return FoundEvent{}, false, nil
		}
	}
}

func inspectJournalFile(journalDir, fileName, checkpointFile string, checkpointOffset int64, inspection *journalInspection) error {
	path := filepath.Join(journalDir, fileName)
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open journal file %q: %w", path, err)
	}
	defer closeReplayFile(file)

	reader := bufio.NewReader(file)
	var offset int64
	for lineNumber := 1; ; lineNumber++ {
		lineBytes, readErr := reader.ReadBytes('\n')
		if readErr != nil && !errors.Is(readErr, io.EOF) {
			return fmt.Errorf("read journal file %q: %w", path, readErr)
		}
		if len(lineBytes) == 0 && errors.Is(readErr, io.EOF) {
			return nil
		}

		offset += int64(len(lineBytes))
		checkpoint, hasEvent, err := inspectCheckpointLine(fileName, lineNumber, offset, lineBytes)
		if err != nil {
			return err
		}
		if hasEvent {
			inspection.tail = journalTail{HasEvents: true, File: fileName, Offset: offset}
		}
		if checkpoint.FileMatches(fileName, checkpointFile, offset, checkpointOffset) {
			inspection.checkpointSeen = true
			inspection.checkpoint = journalCheckpoint(checkpoint)
		}

		if errors.Is(readErr, io.EOF) {
			return nil
		}
	}
}

type checkpointLine struct {
	EventID   string
	EventHash string
}

func (line checkpointLine) FileMatches(fileName, checkpointFile string, offset, checkpointOffset int64) bool {
	return fileName == checkpointFile && offset == checkpointOffset
}

func inspectCheckpointLine(fileName string, lineNumber int, offset int64, lineBytes []byte) (checkpointLine, bool, error) {
	trimmed := bytesTrimSpace(lineBytes)
	if len(trimmed) == 0 {
		return checkpointLine{}, false, nil
	}

	var stored Event
	if err := json.Unmarshal(trimmed, &stored); err != nil {
		return checkpointLine{}, false, &CorruptJournalError{
			File: filepath.ToSlash(filepath.Join(".fid", "journal", fileName)),
			Line: lineNumber,
			Err:  err,
		}
	}

	eventHash, err := HashEvent(stored)
	if err != nil {
		return checkpointLine{}, false, fmt.Errorf("hash event %q from %s:%d: %w", stored.EventID, fileName, lineNumber, err)
	}
	_ = offset
	return checkpointLine{EventID: stored.EventID, EventHash: eventHash}, true, nil
}

func inspectJournalLine(fileName string, lineNumber int, lineBytes []byte, eventID string) (FoundEvent, bool, error) {
	trimmed := bytesTrimSpace(lineBytes)
	if len(trimmed) == 0 {
		return FoundEvent{}, false, nil
	}

	var stored Event
	if err := json.Unmarshal(trimmed, &stored); err != nil {
		return FoundEvent{}, false, &CorruptJournalError{
			File: filepath.ToSlash(filepath.Join(".fid", "journal", fileName)),
			Line: lineNumber,
			Err:  err,
		}
	}
	if stored.EventID != eventID {
		return FoundEvent{}, false, nil
	}

	eventHash, err := HashEvent(stored)
	if err != nil {
		return FoundEvent{}, false, fmt.Errorf("hash event %q from %s:%d: %w", stored.EventID, fileName, lineNumber, err)
	}

	return FoundEvent{
		Event:       stored,
		JournalFile: filepath.ToSlash(filepath.Join(".fid", "journal", fileName)),
		Line:        lineNumber,
		Hash:        eventHash,
	}, true, nil
}

func RevertVerdict(event Event) Verdict {
	switch event.EventType {
	case "issue.transitioned":
		return inverseUpdateVerdict(event, "reverts this transition to the prior status; blocked if the issue has since moved on")
	case "issue.updated":
		return inverseUpdateVerdict(event, "restores the prior issue field values if the current state still matches this event")
	case "issue.concern_addressed":
		return inverseUpdateVerdict(event, "reopens the concern if the current concern state still matches this event")
	case "issue.label.added":
		return inverseRowPresenceVerdict(event, "issue_labels", "labels_absent", "labels_present", "removes the added labels if they are still present on the issue")
	case "issue.label.removed":
		return inverseRowPresenceVerdict(event, "issue_labels", "labels_present", "labels_absent", "restores the removed labels if they are still absent from the issue")
	case "link.added":
		return inverseRowPresenceVerdict(event, "links", "links_absent", "links_present", "removes the declared and inverse links if both are still present")
	case "link.removed":
		return inverseRowPresenceVerdict(event, "links", "links_present", "links_absent", "restores the declared and inverse links if they are still absent")
	case "issue.ac_added":
		return retractionVerdict("retracts this acceptance criterion fact while preserving the journal row")
	case "issue.ac_verified":
		return retractionVerdict("retracts this verification fact while preserving the journal row")
	case "issue.concern_raised":
		return retractionVerdict("retracts this concern fact while preserving the journal row")
	case "dev.gate.recorded":
		return retractionVerdict("retracts this recorded gate result while preserving the journal row")
	case "issue.created":
		return notEligibleVerdict("cancel the issue with fid issue transition <id> --to cancelled --by <role>", "creating the issue was a genesis action and cannot be removed from append-only history")
	case "project.initialized":
		return notEligibleVerdict("create a fresh workspace if the project bootstrap was wrong", "project initialization is a workspace genesis action and is not revertible")
	case "label.added":
		return notEligibleVerdict("stop applying the label, or create the corrected label and move usage forward", "label creation is treated as a catalog creation, not a reversible state flip")
	case "issue.completed":
		return notEligibleVerdict("move the issue forward with the appropriate workflow command instead of removing the completion history", "completion events combine durable delivery facts that are not reverted in v0.1")
	case "issue.claim.acquired", "issue.claim.renewed", "issue.claim.released":
		return notEligibleVerdict("use fid issue claim or fid issue unclaim to move the claim state forward", "claim events are coordination history and are not reverted generically")
	case "run.recorded":
		return notEligibleVerdict("record a correcting run entry or leave the telemetry history intact", "telemetry entries are append-only audit facts and are not reverted in v0.1")
	case "requirement.synced":
		return notEligibleVerdict("rerun fid requirement sync from the corrected source document", "requirement sync is a bulk import surface and is not reverted by event inspection")
	case "dev.started", "dev.finished", "dev.blocked":
		return notEligibleVerdict("use the issue and dev commands to move workflow state forward instead of reverting the composite event", "dev composite events combine multiple state changes and audit facts in one journal row")
	default:
		return notEligibleVerdict("use the domain command that moves the affected state forward", "unknown event types default to a safe non-eligible verdict until a specific revert strategy is defined")
	}
}

func inverseUpdateVerdict(event Event, note string) Verdict {
	restores, precondition := updateRestoreAndPrecondition(event.Effects)
	if len(restores) == 0 {
		restores = map[string]any{"prior_state": "captured in the event"}
		precondition = map[string]any{"current_state": "must still match the values written by this event"}
	}
	return Verdict{
		Eligible:     true,
		Mode:         "inverse",
		Restores:     restores,
		Precondition: precondition,
		Note:         note,
	}
}

func inverseRowPresenceVerdict(event Event, tableName, restoreKey, preconditionKey, note string) Verdict {
	rows := effectPKsForTable(event.Effects, tableName)
	if len(rows) == 0 {
		rows = []map[string]any{{"table": tableName}}
	}
	return Verdict{
		Eligible: true,
		Mode:     "inverse",
		Restores: map[string]any{
			restoreKey: rows,
		},
		Precondition: map[string]any{
			preconditionKey: rows,
		},
		Note: note,
	}
}

func updateRestoreAndPrecondition(effects []Effect) (map[string]any, map[string]any) {
	restores := map[string]any{}
	precondition := map[string]any{}
	for _, effect := range effects {
		if effect.Op != "update" {
			continue
		}
		for key, prior := range effect.Prior {
			restores[key] = prior
			if current, ok := effect.Values[key]; ok {
				precondition[key] = current
				continue
			}
			if current, ok := effect.Preconditions[key]; ok {
				precondition[key] = current
			}
		}
	}
	return restores, precondition
}

func effectPKsForTable(effects []Effect, tableName string) []map[string]any {
	rows := make([]map[string]any, 0, len(effects))
	for _, effect := range effects {
		if effect.Table != tableName {
			continue
		}
		rows = append(rows, cloneEffectMap(effect.PK))
	}
	return rows
}

func cloneEffectMap(input map[string]any) map[string]any {
	if len(input) == 0 {
		return map[string]any{}
	}
	output := make(map[string]any, len(input))
	for key, value := range input {
		output[key] = value
	}
	return output
}

func retractionVerdict(note string) Verdict {
	return Verdict{
		Eligible: true,
		Mode:     "retraction",
		Note:     note,
	}
}

func notEligibleVerdict(correction, note string) Verdict {
	return Verdict{
		Eligible:   false,
		Correction: correction,
		Note:       note,
	}
}
