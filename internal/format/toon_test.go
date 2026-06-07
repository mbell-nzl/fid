package format

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/issues"
)

func TestWriteTOONMatchesJSONStructureForEmptyIssue(t *testing.T) {
	issue := emptyIssueFixture()

	var jsonOutput bytes.Buffer
	if err := Write(&jsonOutput, JSON, issue); err != nil {
		t.Fatalf("Write(JSON) error = %v", err)
	}

	var toonOutput bytes.Buffer
	if err := Write(&toonOutput, TOON, issue); err != nil {
		t.Fatalf("Write(TOON) error = %v", err)
	}

	want := decodeJSONStructure(t, jsonOutput.String())
	got := decodeTOONStructure(t, toonOutput.String())

	if !structuresEqual(want, got) {
		t.Fatalf("TOON structure = %#v, want %#v", got, want)
	}
}

func TestWriteTOONMatchesJSONStructureForPopulatedIssue(t *testing.T) {
	issue := populatedIssueFixture()

	var jsonOutput bytes.Buffer
	if err := Write(&jsonOutput, JSON, issue); err != nil {
		t.Fatalf("Write(JSON) error = %v", err)
	}

	var toonOutput bytes.Buffer
	if err := Write(&toonOutput, TOON, issue); err != nil {
		t.Fatalf("Write(TOON) error = %v", err)
	}

	want := decodeJSONStructure(t, jsonOutput.String())
	got := decodeTOONStructure(t, toonOutput.String())

	if !structuresEqual(want, got) {
		t.Fatalf("TOON structure = %#v, want %#v", got, want)
	}
}

func TestWriteTOONMatchesGoldenFixtureForEmptyIssue(t *testing.T) {
	issue := emptyIssueFixture()

	var toonOutput bytes.Buffer
	if err := Write(&toonOutput, TOON, issue); err != nil {
		t.Fatalf("Write(TOON) error = %v", err)
	}

	golden, err := os.ReadFile(toonGoldenPath(t, "issue-empty.golden.toon"))
	if err != nil {
		t.Fatalf("ReadFile() golden error = %v", err)
	}

	if toonOutput.String() != string(golden) {
		t.Fatalf("TOON output = %q, want golden %q", toonOutput.String(), string(golden))
	}
}

func decodeJSONStructure(t *testing.T, value string) any {
	t.Helper()

	decoder := json.NewDecoder(strings.NewReader(value))
	decoder.UseNumber()

	var result any
	if err := decoder.Decode(&result); err != nil {
		t.Fatalf("Decode JSON structure error = %v", err)
	}
	return result
}

func decodeTOONStructure(t *testing.T, value string) any {
	t.Helper()

	parser := toonParser{input: strings.TrimSpace(value)}
	result, err := parser.parseValue()
	if err != nil {
		t.Fatalf("parse TOON structure error = %v", err)
	}
	parser.skipWhitespace()
	if parser.position != len(parser.input) {
		t.Fatalf("parse TOON structure stopped at %d of %d", parser.position, len(parser.input))
	}
	return result
}

func structuresEqual(left any, right any) bool {
	switch leftValue := left.(type) {
	case map[string]any:
		rightValue, ok := right.(map[string]any)
		if !ok || len(leftValue) != len(rightValue) {
			return false
		}
		for key, nestedLeft := range leftValue {
			nestedRight, ok := rightValue[key]
			if !ok || !structuresEqual(nestedLeft, nestedRight) {
				return false
			}
		}
		return true
	case []any:
		rightValue, ok := right.([]any)
		if !ok || len(leftValue) != len(rightValue) {
			return false
		}
		for index := range leftValue {
			if !structuresEqual(leftValue[index], rightValue[index]) {
				return false
			}
		}
		return true
	case json.Number:
		return numericStructuresEqual(leftValue, right)
	default:
		return fmt.Sprintf("%v", left) == fmt.Sprintf("%v", right)
	}
}

func numericStructuresEqual(left json.Number, right any) bool {
	switch rightValue := right.(type) {
	case json.Number:
		return left.String() == rightValue.String()
	case int:
		return left.String() == strconv.Itoa(rightValue)
	case int8:
		return left.String() == strconv.FormatInt(int64(rightValue), 10)
	case int16:
		return left.String() == strconv.FormatInt(int64(rightValue), 10)
	case int32:
		return left.String() == strconv.FormatInt(int64(rightValue), 10)
	case int64:
		return left.String() == strconv.FormatInt(rightValue, 10)
	case uint:
		return left.String() == strconv.FormatUint(uint64(rightValue), 10)
	case uint8:
		return left.String() == strconv.FormatUint(uint64(rightValue), 10)
	case uint16:
		return left.String() == strconv.FormatUint(uint64(rightValue), 10)
	case uint32:
		return left.String() == strconv.FormatUint(uint64(rightValue), 10)
	case uint64:
		return left.String() == strconv.FormatUint(rightValue, 10)
	case float32:
		return left.String() == strconv.FormatFloat(float64(rightValue), 'f', -1, 32)
	case float64:
		return left.String() == strconv.FormatFloat(rightValue, 'f', -1, 64)
	default:
		return false
	}
}

type toonParser struct {
	input    string
	position int
}

func (parser *toonParser) parseValue() (any, error) {
	parser.skipWhitespace()
	if parser.position >= len(parser.input) {
		return nil, fmt.Errorf("unexpected end of input")
	}

	switch parser.input[parser.position] {
	case '{':
		return parser.parseObject()
	case '[':
		return parser.parseArray()
	case '"':
		return parser.parseString()
	default:
		return parser.parseScalar()
	}
}

func (parser *toonParser) parseObject() (map[string]any, error) {
	parser.position++
	result := map[string]any{}

	parser.skipWhitespace()
	if parser.consumeIf('}') {
		return result, nil
	}

	for {
		key, err := parser.parseKey()
		if err != nil {
			return nil, err
		}
		if !parser.consumeIf(':') {
			return nil, fmt.Errorf("expected ':' after key %q", key)
		}
		value, err := parser.parseValue()
		if err != nil {
			return nil, err
		}
		result[key] = value

		parser.skipWhitespace()
		if parser.consumeIf('}') {
			return result, nil
		}
		if !parser.consumeIf(',') {
			return nil, fmt.Errorf("expected ',' or '}' after key %q", key)
		}
	}
}

func (parser *toonParser) parseArray() ([]any, error) {
	parser.position++
	var result []any

	parser.skipWhitespace()
	if parser.consumeIf(']') {
		return result, nil
	}

	for {
		value, err := parser.parseValue()
		if err != nil {
			return nil, err
		}
		result = append(result, value)

		parser.skipWhitespace()
		if parser.consumeIf(']') {
			return result, nil
		}
		if !parser.consumeIf(',') {
			return nil, fmt.Errorf("expected ',' or ']' in array")
		}
	}
}

func (parser *toonParser) parseKey() (string, error) {
	parser.skipWhitespace()
	if parser.position >= len(parser.input) {
		return "", fmt.Errorf("unexpected end of input")
	}
	if parser.input[parser.position] == '"' {
		return parser.parseString()
	}

	start := parser.position
	for parser.position < len(parser.input) {
		switch parser.input[parser.position] {
		case ':', ' ', '\n', '\t', '\r':
			return parser.input[start:parser.position], nil
		default:
			parser.position++
		}
	}
	return parser.input[start:parser.position], nil
}

func (parser *toonParser) parseString() (string, error) {
	start := parser.position
	parser.position++
	escaped := false
	for parser.position < len(parser.input) {
		current := parser.input[parser.position]
		if escaped {
			escaped = false
			parser.position++
			continue
		}
		if current == '\\' {
			escaped = true
			parser.position++
			continue
		}
		if current == '"' {
			parser.position++
			var result string
			if err := json.Unmarshal([]byte(parser.input[start:parser.position]), &result); err != nil {
				return "", err
			}
			return result, nil
		}
		parser.position++
	}
	return "", fmt.Errorf("unterminated string")
}

func (parser *toonParser) parseScalar() (any, error) {
	start := parser.position
	for parser.position < len(parser.input) {
		switch parser.input[parser.position] {
		case ',', '}', ']', ' ', '\n', '\t', '\r':
			token := parser.input[start:parser.position]
			return parseTOONScalar(token)
		default:
			parser.position++
		}
	}
	return parseTOONScalar(parser.input[start:parser.position])
}

func (parser *toonParser) skipWhitespace() {
	for parser.position < len(parser.input) {
		switch parser.input[parser.position] {
		case ' ', '\n', '\t', '\r':
			parser.position++
		default:
			return
		}
	}
}

func (parser *toonParser) consumeIf(expected byte) bool {
	parser.skipWhitespace()
	if parser.position < len(parser.input) && parser.input[parser.position] == expected {
		parser.position++
		return true
	}
	return false
}

func parseTOONScalar(token string) (any, error) {
	switch token {
	case "null":
		return nil, nil
	case "true":
		return true, nil
	case "false":
		return false, nil
	}

	if _, err := strconv.ParseInt(token, 10, 64); err == nil {
		return json.Number(token), nil
	}
	if _, err := strconv.ParseFloat(token, 64); err == nil {
		return json.Number(token), nil
	}
	return token, nil
}

func emptyIssueFixture() issues.Issue {
	return issues.Issue{
		ID:                 "FID-empty01",
		Kind:               "task",
		Title:              "Render empty issue",
		Status:             "todo",
		StatusCategory:     "unstarted",
		WorkflowID:         "task-default",
		Labels:             []string{},
		InInbox:            false,
		ConcernsRaised:     false,
		OpenConcernsCount:  0,
		AcceptanceCriteria: []issues.AcceptanceCriterion{},
	}
}

func populatedIssueFixture() issues.Issue {
	parentID := "FID-epic001"
	priority := 1
	sequence := 2
	specMDPath := "project-docs/specs/task-1.md"
	estTokens := 987
	ownerRole := "developer"
	ownerActor := "mark"
	triageStatus := "blocked"
	createdAt := "2026-06-07T14:00:00Z"
	updatedAt := "2026-06-07T14:05:00Z"

	return issues.Issue{
		ID:                "FID-pop001",
		Kind:              "task",
		Title:             "Render populated issue",
		Status:            "qa_block",
		StatusCategory:    "started",
		WorkflowID:        "task-default",
		ParentID:          &parentID,
		Priority:          &priority,
		Sequence:          &sequence,
		SpecMDPath:        &specMDPath,
		EstTokens:         &estTokens,
		OwnerRole:         &ownerRole,
		OwnerActor:        &ownerActor,
		Labels:            []string{"backend", "phase-3"},
		InInbox:           true,
		ConcernsRaised:    true,
		OpenConcernsCount: 2,
		TriageStatus:      &triageStatus,
		AcceptanceCriteria: []issues.AcceptanceCriterion{
			{
				ID:      "AC-001",
				Ordinal: 1,
				Text:    "First nested check",
				LatestVerifications: map[string]any{
					"checks": []any{"unit", "smoke"},
					"status": "pass",
				},
			},
			{
				ID:      "AC-002",
				Ordinal: 2,
				Text:    "Second nested check",
				LatestVerifications: map[string]any{
					"checks": []any{"docs"},
					"status": "todo",
				},
			},
		},
		CreatedAt: &createdAt,
		UpdatedAt: &updatedAt,
	}
}

func toonGoldenPath(t *testing.T, name string) string {
	t.Helper()

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller(0) failed")
	}

	return filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", "testdata", "format", name))
}
