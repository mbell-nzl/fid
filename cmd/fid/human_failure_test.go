package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/spf13/cobra"
)

func TestP24FamilyFailureRenderer(t *testing.T) {
	t.Parallel()
	rows := []struct {
		name   string
		render func(*cobra.Command, error) error
	}{
		{"issue", func(c *cobra.Command, e error) error { return renderIssueMutationError(c, "issue.update", "ISS-1", e) }},
		{"issue-read", func(c *cobra.Command, e error) error { return renderIssueReadError(c, false, "issue.get", e) }},
		{"activation", func(c *cobra.Command, e error) error { return renderActivationError(c, "activation.prepare", e) }},
		{"claim", func(c *cobra.Command, e error) error { return renderClaimError(c, "claim.acquire", e) }},
		{"run", func(c *cobra.Command, e error) error { return renderRunError(c, "run.start", e) }},
		{"envelope", func(c *cobra.Command, e error) error { return renderEnvelopeError(c, "envelope.create", e, "actor") }},
		{"gate", func(c *cobra.Command, e error) error { return renderGateError(c, "gate.register", e) }},
		{"evidence", func(c *cobra.Command, e error) error { return renderEvidenceError(c, "evidence.link", e) }},
		{"stop", func(c *cobra.Command, e error) error { return renderStopError(c, "stop.raise", e) }},
		{"decision", func(c *cobra.Command, e error) error { return renderDecisionError(c, "decision.record-reference", e) }},
		{"human-outcome", func(c *cobra.Command, e error) error { return renderHumanOutcomeError(c, "human-outcome.record", e) }},
		{"concern", func(c *cobra.Command, e error) error { return renderConcernError(c, "concern.record", e) }},
		{"deviation", func(c *cobra.Command, e error) error { return renderDeviationError(c, "deviation.record", e) }},
		{"finding", func(c *cobra.Command, e error) error { return renderFindingError(c, "finding.address", e) }},
		{"review", renderReviewError}, {"review-list", renderReviewListError},
		{"link", func(c *cobra.Command, e error) error { return renderLinkError(c, "link.add", e) }},
		{"label", func(c *cobra.Command, e error) error { return renderLabelError(c, "label.add", e) }},
		{"workflow", func(c *cobra.Command, e error) error { return renderWorkflowError(c, "workflow.define", e) }},
		{"receipt", func(c *cobra.Command, e error) error { return renderReceiptError(c, "receipt.get", e) }},
		{"init", func(c *cobra.Command, e error) error { return renderInitFailure(c, "/fixture", e) }},
		{"rebuild", func(c *cobra.Command, e error) error { return renderRebuildOperationError(c, false, e) }},
		{"maintenance", renderMaintenanceCheckError}, {"capabilities", renderCapabilitiesError},
		{"export", func(c *cobra.Command, e error) error {
			return renderExportError(c, "export.validate", "archive.zip", e)
		}},
	}
	for _, row := range rows {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			c := &cobra.Command{Use: "fid"}
			c.PersistentFlags().Bool("json", false, "")
			var stderr bytes.Buffer
			c.SetErr(&stderr)
			c.SetOut(&bytes.Buffer{})
			err := row.render(c, errors.New("domain cause"))
			var failure *cli.FailureError
			if !errors.As(err, &failure) {
				t.Fatalf("not a typed failure: %v", err)
			}
			if failure.Code == cli.CodeInvalidArgument && !strings.HasPrefix(stderr.String(), "error: domain cause\n") {
				t.Fatalf("generic message lost cause: %q", stderr.String())
			}
			if !strings.HasPrefix(stderr.String(), "error: ") || !strings.Contains(stderr.String(), "\nreason: "+string(failure.Code)) {
				t.Fatalf("Human failure bypass: %q for %#v", stderr.String(), failure)
			}
		})
	}
}

func TestP24GenericFallbackMessage(t *testing.T) {
	t.Parallel()
	rows := []struct{ name, cause, want string }{
		{"512", strings.Repeat("a", 512), strings.Repeat("a", 512)},
		{"513", strings.Repeat("a", 513), strings.Repeat("a", 511) + "…"},
		{"non-ASCII", strings.Repeat("é", 513), strings.Repeat("é", 511) + "…"},
		{"ESC", "domain\x1brefused", `domain\u001brefused`},
		{"escaped 512", strings.Repeat("a", 506) + "\x1b", strings.Repeat("a", 506) + `\u001b`},
		{"escaped 513", "\x1b" + strings.Repeat("a", 507), `\u001b` + strings.Repeat("a", 505) + "…"},
		{"escape crosses cut", strings.Repeat("a", 510) + "\x1bX", strings.Repeat("a", 510) + "…"},
		{"empty", "", "invalid issue read argument"},
		{"whitespace", " \t\n ", "invalid issue read argument"},
	}
	for _, row := range rows {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			failure, err := issueReadFailure("issue.get", errors.New(row.cause))
			if err != nil {
				t.Fatal(err)
			}
			want := row.want
			if failure.Message != want {
				t.Fatalf("message=%q want=%q", failure.Message, want)
			}
			var out bytes.Buffer
			if err := cli.WriteFailure(&out, "issue.get", failure); err != nil {
				t.Fatal(err)
			}
			var payload map[string]any
			if err := json.Unmarshal(out.Bytes(), &payload); err != nil {
				t.Fatal(err)
			}
		})
	}
}

// Expected grammar text comes from the accepted class row and the live flag tree,
// never from the production renderer or the observed Human output.
func p24GrammarHumanText(t *testing.T, row grammarRow) string {
	t.Helper()
	root := newRootCommandAt(t.TempDir())
	leaf := root
	if row.operation != "fid" {
		var err error
		leaf, _, err = root.Find(strings.Split(row.operation, "."))
		if err != nil {
			t.Fatal(err)
		}
	}
	detail := "argument=" + row.argument + ", allowed=" + strings.Join(row.allowed, "|")
	if row.code == "operation-unavailable" {
		detail = "allowed=" + strings.Join(row.allowed, ",") + ", argument=command, value=" + fmt.Sprint(row.value)
	}
	reason := p24ExpectedReason("reason: " + row.code + "; " + detail)
	next := ""
	if row.code == "invalid-argument" && leaf.Flags().Lookup("input") != nil {
		next = "next: see " + leaf.CommandPath() + " --help; template: " + leaf.CommandPath() + " --example\n"
	}
	return "error: " + row.message + "\n" + reason + "\n" + next
}

func p24ExpectedReason(text string) string {
	escaped := cli.EscapeControlText(text)
	if utf8.RuneCountInString(escaped) <= 200 {
		return escaped
	}
	tokens := regexp.MustCompile(`\\(?:[abfnrtv\\"']|x[0-9a-fA-F]{2}|u[0-9a-fA-F]{4}|U[0-9a-fA-F]{8})|.`).FindAllString(escaped, -1)
	var prefix strings.Builder
	count := 0
	for _, token := range tokens {
		size := utf8.RuneCountInString(token)
		if count+size > 200 {
			break
		}
		prefix.WriteString(token)
		count += size
	}
	return prefix.String() + "…"
}
