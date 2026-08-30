package evidence

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"syscall"
	"unicode/utf8"
)

const (
	liveFrameSchema = "fid.live.v1"
	redactedText    = "[REDACTED]"
	maxLiveChunk    = 4096
)

// LiveFrame is the closed fid.live.v1 transport object.
type LiveFrame struct {
	FrameSchema string         `json:"frame_schema"`
	ExecutionID string         `json:"execution_id"`
	Sequence    int            `json:"sequence"`
	Stream      string         `json:"stream"`
	Data        map[string]any `json:"data"`
}

// ArtifactRejection is the public metadata retained for a rejected candidate.
type ArtifactRejection struct {
	Path   string `json:"path"`
	Reason string `json:"reason"`
}

// LiveResult contains bounded, redacted process observations. Durable Attempt
// and Receipt binding is deliberately supplied by the task-15 caller.
type LiveResult struct {
	Frames               []LiveFrame
	Outcome              string
	ExitCode             *int
	Signal               *string
	TimedOut             bool
	OutputBytes          int64
	OutputTruncated      bool
	RedactionApplied     bool
	RedactionReplaced    int
	RejectedArtifacts    []ArtifactRejection
	executionID          string
	terminalFrameEmitted bool
}

// LiveTerminal supplies the durable caller-owned terminal binding.
type LiveTerminal struct {
	AttemptID   string
	ReceiptID   *string
	Outcome     string
	ExitCode    *int
	Signal      *string
	TimedOut    bool
	ArtifactIDs []string
	Response    map[string]any
	Replayed    bool
}

type liveRead struct {
	stream string
	data   []byte
	err    error
}

// RunLive executes one validated definition, drains both process pipes, and
// emits only bounded redacted stdout/stderr/control frames.
func RunLive(ctx context.Context, workspaceRoot string, definition GateDefinition, executionID string) (LiveResult, error) {
	if strings.TrimSpace(executionID) == "" {
		return LiveResult{}, fmt.Errorf("run live: execution ID is required")
	}
	command, processContext, cancel, err := executorCommand(ctx, workspaceRoot, definition)
	if err != nil {
		return LiveResult{}, err
	}
	defer cancel()
	stdout, err := command.StdoutPipe()
	if err != nil {
		return LiveResult{}, fmt.Errorf("run live stdout pipe: %w", err)
	}
	stderr, err := command.StderrPipe()
	if err != nil {
		return LiveResult{}, fmt.Errorf("run live stderr pipe: %w", err)
	}
	if err := command.Start(); err != nil {
		return spawnFailedLiveResult(executionID), nil
	}

	reads := make(chan liveRead, 2)
	go readLivePipe(stdout, "stdout", reads)
	go readLivePipe(stderr, "stderr", reads)
	result := LiveResult{executionID: executionID, RejectedArtifacts: []ArtifactRejection{}}
	sources := redactionSources(definition)
	redactors := map[string]*liveRedactor{"stdout": newLiveRedactor(sources), "stderr": newLiveRedactor(sources)}
	completedReads := 0
	for completedReads < 2 {
		read := <-reads
		if len(read.data) > 0 {
			result.consumeLiveData(read.stream, redactors[read.stream].Transform(read.data, false), definition.OutputBytes)
		}
		if read.err != nil {
			return result, fmt.Errorf("run live %s: %w", read.stream, read.err)
		}
		if read.data == nil {
			result.consumeLiveData(read.stream, redactors[read.stream].Transform(nil, true), definition.OutputBytes)
			completedReads++
		}
	}
	if err := command.Wait(); err != nil {
		if processContext.Err() != nil {
			result.Outcome = "timed-out"
			result.TimedOut = true
		} else {
			result.Outcome = "failed"
		}
	} else {
		result.Outcome = "succeeded"
	}
	if command.ProcessState != nil {
		exitCode := command.ProcessState.ExitCode()
		result.ExitCode = &exitCode
		if status, ok := command.ProcessState.Sys().(syscall.WaitStatus); ok && status.Signaled() {
			signal := status.Signal().String()
			result.Signal = &signal
			result.ExitCode = nil
			if !result.TimedOut {
				result.Outcome = "signaled"
			}
		}
		if result.Outcome == "succeeded" && !containsExitCode(definition.SuccessExitCodes, exitCode) {
			result.Outcome = "failed"
		}
	}
	result.RedactionReplaced = redactors["stdout"].replacements + redactors["stderr"].replacements
	result.RedactionApplied = result.RedactionReplaced > 0
	if result.RedactionApplied {
		result.appendControl(map[string]any{"kind": "redaction-applied", "replacements": result.RedactionReplaced})
	}
	result.RejectedArtifacts = rejectedArtifacts(workspaceRoot, definition.CWD, definition.ArtifactPaths, sources)
	for _, rejected := range result.RejectedArtifacts {
		result.appendControl(map[string]any{"kind": "artifact-rejected", "path": rejected.Path, "reason": rejected.Reason})
	}
	return result, nil
}

func readLivePipe(reader io.Reader, stream string, reads chan<- liveRead) {
	buffer := make([]byte, 32*1024)
	for {
		count, err := reader.Read(buffer)
		if count > 0 {
			reads <- liveRead{stream: stream, data: append([]byte(nil), buffer[:count]...)}
		}
		if err == io.EOF {
			reads <- liveRead{stream: stream}
			return
		}
		if err != nil {
			reads <- liveRead{stream: stream, err: err}
			return
		}
	}
}

func spawnFailedLiveResult(executionID string) LiveResult {
	return LiveResult{executionID: executionID, Outcome: "spawn-failed", RejectedArtifacts: []ArtifactRejection{}}
}

func (result *LiveResult) consumeLiveData(stream string, text string, limit int64) {
	if text == "" {
		return
	}
	if result.OutputTruncated {
		return
	}
	result.OutputBytes += int64(len(text))
	emitted := result.emittedDataBytes()
	remaining := limit - emitted
	if remaining <= 0 {
		result.emitTruncation(limit)
		return
	}
	if int64(len(text)) > remaining {
		result.appendData(stream, truncateUTF8(text, int(remaining)))
		result.emitTruncation(limit)
		return
	}
	result.appendData(stream, text)
}

func (result *LiveResult) emittedDataBytes() int64 {
	var total int64
	for _, frame := range result.Frames {
		if frame.Stream == "stdout" || frame.Stream == "stderr" {
			total += int64(len(frame.Data["text"].(string)))
		}
	}
	return total
}

func (result *LiveResult) emitTruncation(limit int64) {
	if result.OutputTruncated {
		return
	}
	result.OutputTruncated = true
	result.appendControl(map[string]any{"kind": "output-truncated", "limit_bytes": limit, "observed_bytes_lower_bound": result.OutputBytes})
}

func (result *LiveResult) appendData(stream, text string) {
	for len(text) > 0 {
		chunk := truncateUTF8(text, maxLiveChunk)
		result.appendFrame(stream, map[string]any{"encoding": "utf-8", "text": chunk})
		text = text[len(chunk):]
	}
}

func (result *LiveResult) appendControl(data map[string]any) { result.appendFrame("control", data) }

func (result *LiveResult) appendFrame(stream string, data map[string]any) {
	result.Frames = append(result.Frames, LiveFrame{FrameSchema: liveFrameSchema, ExecutionID: result.executionID, Sequence: len(result.Frames) + 1, Stream: stream, Data: data})
}

// AppendTerminal adds the one required final frame after the caller has bound
// the observed execution to its durable Attempt and optional Receipt.
func (result *LiveResult) AppendTerminal(terminal LiveTerminal) error {
	if result.terminalFrameEmitted || strings.TrimSpace(terminal.AttemptID) == "" || terminal.Response == nil {
		return fmt.Errorf("append live terminal: one attempt ID and response are required")
	}
	if terminal.Outcome != result.Outcome || terminal.TimedOut != result.TimedOut || !sameIntPointer(terminal.ExitCode, result.ExitCode) || !sameStringPointer(terminal.Signal, result.Signal) {
		return fmt.Errorf("append live terminal: terminal observation does not match process result")
	}
	artifacts := append([]string{}, terminal.ArtifactIDs...)
	sort.Strings(artifacts)
	rejected := append([]ArtifactRejection(nil), result.RejectedArtifacts...)
	result.appendFrame("terminal", map[string]any{
		"attempt_id": terminal.AttemptID, "receipt_id": terminal.ReceiptID, "outcome": terminal.Outcome,
		"exit_code": terminal.ExitCode, "signal": terminal.Signal, "timed_out": terminal.TimedOut,
		"output_bytes": result.OutputBytes, "output_truncated": result.OutputTruncated,
		"redaction_applied": result.RedactionApplied, "artifact_ids": artifacts,
		"rejected_artifacts": rejected, "response": terminal.Response, "replayed": terminal.Replayed,
	})
	result.terminalFrameEmitted = true
	return nil
}

// AppendFidError closes a started stream when fid itself cannot complete the
// persistence path. It intentionally creates neither a Receipt nor an
// idempotency result.
func (result *LiveResult) AppendFidError(attemptID string, response map[string]any) error {
	if result.terminalFrameEmitted || len(result.Frames) == 0 || strings.TrimSpace(attemptID) == "" || response == nil {
		return fmt.Errorf("append live fid error: started frames, attempt ID, and response are required")
	}
	result.appendFrame("terminal", map[string]any{
		"attempt_id": attemptID, "receipt_id": nil, "outcome": "fid-error", "exit_code": nil, "signal": nil,
		"timed_out": false, "output_bytes": result.OutputBytes, "output_truncated": result.OutputTruncated,
		"redaction_applied": result.RedactionApplied, "artifact_ids": []string{}, "rejected_artifacts": result.RejectedArtifacts,
		"response": response, "replayed": false,
	})
	result.terminalFrameEmitted = true
	return nil
}

// JSONLines returns canonical newline-delimited frame objects.
func (result LiveResult) JSONLines() []string {
	lines := make([]string, 0, len(result.Frames))
	for _, frame := range result.Frames {
		payload, err := json.Marshal(map[string]any{"data": frame.Data, "execution_id": frame.ExecutionID, "frame_schema": frame.FrameSchema, "sequence": frame.Sequence, "stream": frame.Stream})
		if err != nil {
			panic(fmt.Sprintf("marshal live frame: %v", err))
		}
		lines = append(lines, string(payload))
	}
	return lines
}

// ControlKinds returns only the closed control variants in observed order.
func (result LiveResult) ControlKinds() []string {
	var kinds []string
	for _, frame := range result.Frames {
		if frame.Stream == "control" {
			kinds = append(kinds, frame.Data["kind"].(string))
		}
	}
	return kinds
}

func containsExitCode(values []int, value int) bool {
	for _, candidate := range values {
		if candidate == value {
			return true
		}
	}
	return false
}

func sameIntPointer(left, right *int) bool {
	return left == nil && right == nil || left != nil && right != nil && *left == *right
}

func sameStringPointer(left, right *string) bool {
	return left == nil && right == nil || left != nil && right != nil && *left == *right
}

func truncateUTF8(text string, maximum int) string {
	if len(text) <= maximum {
		return text
	}
	end := maximum
	for end > 0 && text[end]&0xc0 == 0x80 {
		end--
	}
	if end == 0 {
		return ""
	}
	return text[:end]
}

type liveRedactor struct {
	sources      [][]byte
	pending      []byte
	replacements int
}

func newLiveRedactor(sources []string) *liveRedactor {
	values := make([][]byte, 0, len(sources))
	for _, source := range sources {
		if source != "" {
			values = append(values, []byte(source))
		}
	}
	sort.Slice(values, func(i, j int) bool { return len(values[i]) > len(values[j]) })
	return &liveRedactor{sources: values}
}

func redactionSources(definition GateDefinition) []string {
	sources := append([]string(nil), definition.RedactLiterals...)
	for _, entry := range definition.Environment {
		if entry.Source == "credential" {
			sources = append(sources, os.Getenv("FID_CREDENTIAL_"+entry.Name))
		}
	}
	return sources
}

func (redactor *liveRedactor) Transform(data []byte, final bool) string {
	redactor.pending = append(redactor.pending, data...)
	var output []byte
	for len(redactor.pending) > 0 {
		if source := redactor.matchingSource(); source != nil {
			output = append(output, redactedText...)
			redactor.pending = redactor.pending[len(source):]
			redactor.replacements++
			continue
		}
		if !final && redactor.pendingCouldBeSource() {
			break
		}
		output = append(output, redactor.pending[0])
		redactor.pending = redactor.pending[1:]
	}
	return validUTF8(output)
}

func (redactor *liveRedactor) matchingSource() []byte {
	for _, source := range redactor.sources {
		if bytes.HasPrefix(redactor.pending, source) {
			return source
		}
	}
	return nil
}

func (redactor *liveRedactor) pendingCouldBeSource() bool {
	for _, source := range redactor.sources {
		if bytes.HasPrefix(source, redactor.pending) {
			return true
		}
	}
	return false
}

func validUTF8(value []byte) string {
	var builder strings.Builder
	for len(value) > 0 {
		runeValue, size := utf8.DecodeRune(value)
		if runeValue == utf8.RuneError && size == 1 {
			runeValue = '\uFFFD'
		}
		builder.WriteRune(runeValue)
		value = value[size:]
	}
	return builder.String()
}

func rejectedArtifacts(workspaceRoot, cwd string, paths, sources []string) []ArtifactRejection {
	rejected := make([]ArtifactRejection, 0)
	for _, path := range paths {
		reason := artifactRejectionReason(workspaceRoot, cwd, path)
		if reason != "" {
			rejected = append(rejected, ArtifactRejection{Path: redactText(path, sources), Reason: reason})
		}
	}
	sort.Slice(rejected, func(i, j int) bool {
		if rejected[i].Path == rejected[j].Path {
			return rejected[i].Reason < rejected[j].Reason
		}
		return rejected[i].Path < rejected[j].Path
	})
	return rejected
}

func redactText(text string, sources []string) string {
	redactor := newLiveRedactor(sources)
	return redactor.Transform([]byte(text), true)
}

func artifactRejectionReason(workspaceRoot, cwd, path string) string {
	_, reason := artifactCandidatePath(workspaceRoot, cwd, path)
	return reason
}
