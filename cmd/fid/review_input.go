package main

import (
	"reflect"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/reviewrecord"
)

func reviewText(name, requirement, meaning, placeholder string) cli.InputField {
	return cli.InputField{Name: name, Type: "string", Requirement: requirement, Meaning: meaning, Placeholder: placeholder}
}

func reviewDescriptor() cli.InputDescriptor {
	profile := reviewText("profile", "required", "Profile that raised the finding.", "code")
	severity := reviewText("severity", "required", "major blocks; minor requires an answer; nit is optional.", "major")
	severity.Members = []string{"major", "minor", "nit"}
	stream := reviewText("stream", "required", "Assured work may also record a separate session.", "primary")
	stream.Members = []string{"primary", "separate"}
	verdict := reviewText("verdict", "required", "CHANGES exactly when a major remains open.", "APPROVE")
	verdict.Members = []string{"APPROVE", "CHANGES"}
	line := int64(1)
	return cli.InputDescriptor{Operation: "review.record", Purpose: "Record one attributed review session with every required profile.", Type: reflect.TypeFor[reviewrecord.Input](), Fields: []cli.InputField{
		{Name: "reviewer", Type: "object", Requirement: "required", Meaning: "Reviewer attribution, independent of the recorder.", Children: []cli.InputField{
			reviewText("actor", "required", "Typed reviewer actor.", "agent:reviewer"), reviewText("model", "optional", "Model as <family>/<model>; omitted means unknown.", "claude/opus-5.5"), reviewText("effort", "optional", "Declared reviewer effort.", "high"), reviewText("session_id", "required", "Attributed fresh session ID.", "review-session-1"),
		}}, stream,
		{Name: "profiles", Type: "array", Requirement: "required", Meaning: "Every profile this session reviewed.", Element: &cli.InputField{Type: "string", Placeholder: "code"}}, verdict,
		reviewText("claim", "required", "What the review checked.", "Persist confirmed merges and completed CI attempts."),
		{Name: "findings", Type: "array", Requirement: "required", Meaning: "Findings; use [] for none. fid assigns IDs.", Element: &cli.InputField{Type: "object", Children: []cli.InputField{profile, severity, reviewText("summary", "required", "Problem found.", "Conflicting merge retries have no regression coverage."), {Name: "location", Type: "object", Requirement: "nullable", Meaning: "File and positive line, or null.", Children: []cli.InputField{reviewText("file", "required", "Repository-relative file.", "cmd/fid/merge_records_test.go"), {Name: "line", Type: "integer", Requirement: "required", Meaning: "Line number.", MinInt: &line}}}, reviewText("evidence", "required", "Evidence supporting the finding.", "Disabling the conflict guard left the focused tests green."), reviewText("fix", "required", "Requested action.", "Retry the same key with changed task and PR; assert refusal and unchanged state.")}}},
		reviewText("verification", "required", "Commands, results, probes and limitations.", "Focused tests passed; conflict-guard probe was not caught."),
		{Name: "provenance", Type: "object", Requirement: "optional", Meaning: "Original review prompt, command and raw output.", Children: []cli.InputField{reviewText("prompt", "optional", "Launch prompt.", "Review PR 79 with code and test-adequacy."), reviewText("command", "optional", "Launch command.", "codex review"), reviewText("output", "optional", "Unchanged reviewer output.", "## Review: CHANGES")}},
	}}
}
