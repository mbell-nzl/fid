package format

import (
	"bytes"
	"strings"
	"testing"
)

func TestWriteDispatchesStructuredFormats(t *testing.T) {
	value := map[string]any{
		"id":     "FID-a1b2c3",
		"labels": []any{},
	}

	for _, formatName := range []string{JSON, TOON, YAML} {
		t.Run(formatName, func(t *testing.T) {
			var output bytes.Buffer
			if err := Write(&output, formatName, value); err != nil {
				t.Fatalf("Write(%s) error = %v", formatName, err)
			}
			if output.Len() == 0 {
				t.Fatalf("Write(%s) produced empty output", formatName)
			}
		})
	}
}

func TestWriteDispatchesMarkdown(t *testing.T) {
	var output bytes.Buffer
	if err := Write(&output, Markdown, populatedIssueMarkdownFixture()); err != nil {
		t.Fatalf("Write(%s) error = %v", Markdown, err)
	}
	if output.Len() == 0 {
		t.Fatalf("Write(%s) produced empty output", Markdown)
	}
}

func TestWriteRejectsUnsupportedFormats(t *testing.T) {
	testCases := []struct {
		name       string
		formatName string
		want       string
	}{
		{
			name:       "unknown format names known formats",
			formatName: "xml",
			want:       `unsupported format "xml"; known formats: json, toon, markdown, yaml`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := Write(&bytes.Buffer{}, testCase.formatName, map[string]any{"id": "FID-a1b2c3"})
			if err == nil {
				t.Fatalf("Write(%q) error = nil, want rejection", testCase.formatName)
			}
			if !strings.Contains(err.Error(), testCase.want) {
				t.Fatalf("Write(%q) error = %q, want substring %q", testCase.formatName, err.Error(), testCase.want)
			}
		})
	}
}
