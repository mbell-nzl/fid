package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestV01ReadCommandsSupportAllFormatsWithStructuralParity(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	seed := seedV01FormatMatrixWorkspace(t, workspace)

	cases := []v01FormatMatrixCase{
		{name: "issue get", args: []string{"issue", "get", seed.issueID}, markdownIdentity: []string{seed.issueID, "Format matrix task"}},
		{name: "issue list", args: []string{"issue", "list", "--order", "created"}, markdownIdentity: []string{seed.issueID, seed.peerIssueID}},
		{name: "issue ac list", args: []string{"issue", "ac", "list", seed.issueID}, markdownIdentity: []string{seed.acceptanceCriterionID, "Format matrix acceptance criterion"}},
		{name: "ready", args: []string{"ready"}, markdownIdentity: []string{seed.issueID, "Format matrix task"}},
		{name: "concern list", args: []string{"concern", "list"}, markdownIdentity: []string{seed.concernID, seed.concernIssueID}},
		{name: "claim list", args: []string{"claim", "list"}, markdownIdentity: []string{seed.issueID, "developer"}},
		{name: "link list", args: []string{"link", "list", seed.issueID}, markdownIdentity: []string{seed.peerIssueID, "blocks"}},
		{name: "label list", args: []string{"label", "list"}, markdownIdentity: []string{"matrix-label"}},
		{name: "policy show", args: []string{"policy", "show"}, markdownIdentity: []string{"issue.create"}},
		{name: "stats tokens", args: []string{"stats", "tokens"}, markdownIdentity: []string{"developer", "1200"}},
		{name: "stats cost", args: []string{"stats", "cost"}, markdownIdentity: []string{"developer", "1200"}},
		{name: "activity list", args: []string{"activity", "list", "--on", "issue/" + seed.issueID}, markdownIdentity: []string{seed.issueID, "issue_created"}},
		{name: "sql", args: []string{"sql", fmt.Sprintf("SELECT id, title FROM issues WHERE id = '%s'", seed.issueID)}, markdownIdentity: []string{seed.issueID, "Format matrix task"}},
		{name: "journal inspect", args: []string{"journal", "inspect", seed.journalEventID}, markdownIdentity: []string{seed.journalEventID, seed.journalEventType}},
		{name: "maintenance check", args: []string{"maintenance", "check"}, markdownIdentity: []string{"Status:", "warn"}},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			assertV01ReadFormatParity(t, workspace, testCase)
		})
	}
}

type v01FormatMatrixSeed struct {
	issueID               string
	peerIssueID           string
	acceptanceCriterionID string
	concernIssueID        string
	concernID             string
	journalEventID        string
	journalEventType      string
}

func seedV01FormatMatrixWorkspace(t *testing.T, workspace string) v01FormatMatrixSeed {
	t.Helper()

	clonePristineWorkspaceInto(t, workspace)

	_, issueStderr := runFIDCommand(t, workspace,
		"issue", "create", "--kind", "task", "--title", "Format matrix task", "--priority", "0", "--by", "developer",
	)
	issueID := parseCreatedIssueID(t, issueStderr)

	_, peerStderr := runFIDCommand(t, workspace,
		"issue", "create", "--kind", "task", "--title", "Format matrix linked peer", "--priority", "1", "--by", "developer",
	)
	peerIssueID := parseCreatedIssueID(t, peerStderr)

	_, acStderr := runFIDCommand(t, workspace,
		"issue", "ac", "add", issueID, "--text", "Format matrix acceptance criterion", "--ordinal", "1", "--by", "developer",
	)
	acceptanceCriterionID := parseAddedAcceptanceCriterionID(t, acStderr)

	runFIDCommand(t, workspace, "label", "add", "matrix-label", "--description", "Format matrix label", "--by", "developer")
	runFIDCommand(t, workspace, "issue", "label", "add", issueID, "matrix-label", "--by", "developer")
	runFIDCommand(t, workspace, "link", "add", issueID, "blocks", peerIssueID, "--by", "developer")
	runFIDCommand(t, workspace, "issue", "claim", issueID, "--by", "developer", "--note", "format matrix claim")
	runFIDCommand(t, workspace,
		"run", "record", "--issue", issueID, "--role", "developer", "--model", "claude-sonnet-4-6", "--tokens-in", "1200", "--tokens-out", "3400", "--by", "developer",
	)

	_, concernIssueStderr := runFIDCommand(t, workspace,
		"issue", "create", "--kind", "task", "--title", "Format matrix concern task", "--by", "developer",
	)
	concernIssueID := parseCreatedIssueID(t, concernIssueStderr)
	_, completionStderr := runFIDCommand(t, workspace, "issue", "complete", concernIssueID, "--by", "developer")
	completionID := parseCompletionID(t, completionStderr)
	_, concernStderr := runFIDCommand(t, workspace,
		"concern", "raise", "--on", concernIssueID, "--by", "developer", "--in-completion", completionID, "--summary", "Format matrix concern",
	)
	concernID := parseConcernID(t, concernStderr)

	event := mustReadLastJournalEvent(t, mustJournalPath(t, workspace))
	return v01FormatMatrixSeed{
		issueID:               issueID,
		peerIssueID:           peerIssueID,
		acceptanceCriterionID: acceptanceCriterionID,
		concernIssueID:        concernIssueID,
		concernID:             concernID,
		journalEventID:        event.EventID,
		journalEventType:      event.EventType,
	}
}

type v01FormatMatrixCase struct {
	name             string
	args             []string
	markdownIdentity []string
}

func assertV01ReadFormatParity(t *testing.T, workspace string, testCase v01FormatMatrixCase) {
	t.Helper()

	jsonOutput := runV01ReadFormat(t, workspace, testCase.args, "json")
	jsonStructure := decodeCommandJSONStructure(t, jsonOutput)

	for _, formatName := range []string{"toon", "yaml"} {
		output := runV01ReadFormat(t, workspace, testCase.args, formatName)
		var structure any
		if formatName == "toon" {
			structure = decodeCommandTOONStructure(t, output)
		} else {
			structure = decodeCommandYAMLStructure(t, output)
		}
		if !commandStructuresEqual(jsonStructure, structure) {
			t.Fatalf("%s structure = %#v, want JSON structure %#v", formatName, structure, jsonStructure)
		}
	}

	markdownOutput := runV01ReadFormat(t, workspace, testCase.args, "markdown")
	for _, identity := range testCase.markdownIdentity {
		if !strings.Contains(markdownOutput, identity) {
			t.Fatalf("markdown output = %q, want identity %q", markdownOutput, identity)
		}
	}
}

func runV01ReadFormat(t *testing.T, workspace string, args []string, formatName string) string {
	t.Helper()

	commandArgs := append(append([]string{}, args...), "--format", formatName)
	stdout, stderr := runFIDCommand(t, workspace, commandArgs...)
	if stderr != "" {
		t.Fatalf("%s stderr = %q, want empty", formatName, stderr)
	}
	return stdout
}
