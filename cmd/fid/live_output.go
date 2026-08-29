package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/mbell-nzl/fid/internal/evidence"
)

func renderLiveHuman(result evidence.LiveResult, obligationStatus string) string {
	var lines []string
	for _, frame := range result.Frames {
		switch frame.Stream {
		case "stdout", "stderr":
			lines = append(lines, "["+frame.Stream+"] "+visibleLiveText(frame.Data["text"].(string)))
		case "control":
			lines = append(lines, renderLiveControl(frame.Data))
		case "terminal":
			lines = append(lines, renderLiveTerminal(frame.Data, obligationStatus))
		}
	}
	return strings.Join(lines, "\n") + "\n"
}

func renderLiveControl(data map[string]any) string {
	switch data["kind"] {
	case "redaction-applied":
		return fmt.Sprintf("[control] redaction-applied replacements=%d", data["replacements"])
	case "artifact-rejected":
		return fmt.Sprintf("[control] artifact-rejected path=%s reason=%s", data["path"], data["reason"])
	case "output-truncated":
		return fmt.Sprintf("[control] output-truncated limit_bytes=%d observed_bytes_lower_bound=%d", data["limit_bytes"], data["observed_bytes_lower_bound"])
	default:
		return "[control] unknown"
	}
}

func renderLiveTerminal(data map[string]any, obligationStatus string) string {
	receipt := "none"
	if value, ok := data["receipt_id"].(*string); ok && value != nil {
		receipt = *value
	}
	exit := "none"
	if value, ok := data["exit_code"].(*int); ok && value != nil {
		exit = strconv.Itoa(*value)
	}
	signal := "none"
	if value, ok := data["signal"].(*string); ok && value != nil {
		signal = *value
	}
	return fmt.Sprintf("[terminal] receipt=%s outcome=%s exit=%s signal=%s timed_out=%t artifacts=%s rejected_artifacts=%d obligation=%s",
		receipt, data["outcome"], exit, signal, data["timed_out"], liveArtifactIDs(data["artifact_ids"]), len(data["rejected_artifacts"].([]evidence.ArtifactRejection)), obligationStatus)
}

func liveArtifactIDs(value any) string {
	ids, ok := value.([]string)
	if !ok || len(ids) == 0 {
		return "none"
	}
	return strings.Join(ids, ",")
}

func visibleLiveText(text string) string {
	var builder strings.Builder
	for _, runeValue := range text {
		switch runeValue {
		case '\n':
			builder.WriteString("\\n")
		case '\r':
			builder.WriteString("\\r")
		case '\t':
			builder.WriteString("\\t")
		default:
			if unicode.IsControl(runeValue) {
				_, _ = fmt.Fprintf(&builder, "\\u%04x", runeValue)
			} else {
				builder.WriteRune(runeValue)
			}
		}
	}
	return builder.String()
}
