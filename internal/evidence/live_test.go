package evidence

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"
)

func TestRunLiveFramesInterleavedOutputControlsAndFailure(t *testing.T) {
	workspace := t.TempDir()
	t.Setenv("FID_CREDENTIAL_GATE_TOKEN", "CANARY-live-secret")
	literal := "literal-live-secret"
	definition := GateDefinition{
		ArtifactBytes:     1024,
		ArtifactPaths:     []string{"missing-artifact.txt"},
		Argv:              []string{os.Args[0], "-test.run=^TestLiveProcessChild$", "--", "--live-process-child"},
		CredentialNames:   []string{"GATE_TOKEN"},
		CWD:               ".",
		DefinitionID:      "live-frame-probe",
		DefinitionVersion: 1,
		Environment: []EnvironmentEntry{
			{Name: "GATE_TOKEN", Source: "credential"},
			{Name: "LITERAL", Source: "literal", Value: &literal},
		},
		ExcerptBytes:     64,
		Network:          "allowed",
		OutputBytes:      64,
		RedactLiterals:   []string{"literal-live-secret"},
		Sandbox:          "none",
		Shell:            false,
		SuccessExitCodes: []int{0},
		TimeoutMS:        5000,
		Write:            "allowed",
	}

	result, err := RunLive(context.Background(), workspace, definition, "EXE-live-probe")
	if err != nil {
		t.Fatalf("RunLive() error = %v", err)
	}
	if result.Outcome != "failed" || result.ExitCode == nil || *result.ExitCode != 23 {
		t.Fatalf("RunLive() terminal outcome = (%q, %v), want failed exit 23", result.Outcome, result.ExitCode)
	}
	if !result.OutputTruncated || !result.RedactionApplied {
		t.Fatalf("RunLive() controls = truncated:%t redacted:%t, want both true", result.OutputTruncated, result.RedactionApplied)
	}
	if err := result.AppendTerminal(LiveTerminal{
		AttemptID: "ATT-live-probe",
		Outcome:   result.Outcome,
		ExitCode:  result.ExitCode,
		TimedOut:  result.TimedOut,
		Response: map[string]any{
			"contract_version": "fid-cli.v1",
			"error":            map[string]any{"category": "gate-unsuccessful", "code": "gate-unsuccessful"},
			"ok":               false,
			"operation":        "gate.run",
		},
	}); err != nil {
		t.Fatalf("AppendTerminal() error = %v", err)
	}
	if got := strings.Join(result.JSONLines(), "\n"); strings.Contains(got, "CANARY-live-secret") || strings.Contains(got, "literal-live-secret") {
		t.Fatalf("live frames leak redaction source: %q", got)
	}
	if len(result.Frames) < 6 || result.Frames[len(result.Frames)-1].Stream != "terminal" {
		t.Fatalf("live frames = %#v, want data, three controls, and terminal", result.Frames)
	}
	for index, frame := range result.Frames {
		if frame.Sequence != index+1 {
			t.Fatalf("frame[%d] sequence = %d, want %d", index, frame.Sequence, index+1)
		}
		if frame.FrameSchema != "fid.live.v1" || frame.ExecutionID != "EXE-live-probe" {
			t.Fatalf("frame[%d] header = %#v, want fid.live.v1 / EXE-live-probe", index, frame)
		}
	}
	if got := result.ControlKinds(); strings.Join(got, ",") != "output-truncated,redaction-applied,artifact-rejected" {
		t.Fatalf("control kinds = %v, want exact closed sequence", got)
	}
}

func TestLiveProcessChild(t *testing.T) {
	if !containsArgument(os.Args, "--live-process-child") {
		return
	}
	_, _ = os.Stdout.WriteString("stdout-one CANARY-live-secret\x1b[31m\r")
	time.Sleep(20 * time.Millisecond)
	_, _ = os.Stderr.WriteString("stderr-one literal-live-secret\a\n")
	time.Sleep(20 * time.Millisecond)
	_, _ = os.Stdout.WriteString("stdout-two ")
	time.Sleep(20 * time.Millisecond)
	_, _ = os.Stderr.WriteString("stderr-two ")
	time.Sleep(100 * time.Millisecond)
	_, _ = os.Stdout.WriteString(strings.Repeat("x", 80))
	os.Exit(23)
}

func TestRunLiveDrainsBothPipesAfterTruncationAndPreservesOutcome(t *testing.T) {
	workspace := t.TempDir()
	definition := GateDefinition{
		ArtifactBytes:     1,
		Argv:              []string{os.Args[0], "-test.run=^TestLiveFloodChild$", "--", "--live-flood-child"},
		CWD:               ".",
		DefinitionID:      "live-flood-probe",
		DefinitionVersion: 1,
		Network:           "allowed",
		OutputBytes:       64,
		Sandbox:           "none",
		Shell:             false,
		SuccessExitCodes:  []int{0},
		TimeoutMS:         5000,
		Write:             "allowed",
	}

	result, err := RunLive(context.Background(), workspace, definition, "EXE-live-flood")
	if err != nil {
		t.Fatalf("RunLive() error = %v", err)
	}
	if result.Outcome != "failed" || result.ExitCode == nil || *result.ExitCode != 24 {
		t.Fatalf("RunLive() outcome = (%q, %v), want failed exit 24", result.Outcome, result.ExitCode)
	}
	if !result.OutputTruncated || result.OutputBytes <= definition.OutputBytes {
		t.Fatalf("RunLive() truncation = %t / %d bytes, want true / observed bytes above %d", result.OutputTruncated, result.OutputBytes, definition.OutputBytes)
	}
	if got := strings.Join(result.ControlKinds(), ","); got != "output-truncated" {
		t.Fatalf("control kinds = %q, want one output-truncated control", got)
	}
	var emitted int
	for _, frame := range result.Frames {
		if frame.Stream == "stdout" || frame.Stream == "stderr" {
			emitted += len(frame.Data["text"].(string))
		}
	}
	if emitted > int(definition.OutputBytes) {
		t.Fatalf("emitted process data = %d, want at most %d", emitted, definition.OutputBytes)
	}
}

func TestLiveFloodChild(t *testing.T) {
	if !containsArgument(os.Args, "--live-flood-child") {
		return
	}
	payload := []byte(strings.Repeat("f", 131072))
	_, _ = os.Stdout.Write(payload)
	_, _ = os.Stderr.Write(payload)
	os.Exit(24)
}
