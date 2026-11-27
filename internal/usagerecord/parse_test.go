package usagerecord

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	t.Parallel()
	codex, err := os.ReadFile("testdata/codex.jsonl")
	if err != nil {
		t.Fatal(err)
	}
	claude, err := os.ReadFile("testdata/claude.json")
	if err != nil {
		t.Fatal(err)
	}
	turn := `{"type":"turn.completed","usage":{"input_tokens":100,"cached_input_tokens":40,"cache_write_input_tokens":10,"output_tokens":20,"reasoning_output_tokens":5}}`
	multi := strings.Replace(string(claude), `"modelUsage": {`, `"modelUsage": {"second":{"inputTokens":10,"outputTokens":20,"cacheReadInputTokens":30,"cacheCreationInputTokens":40,"thinkingTokens":5},`, 1)
	for _, tc := range []struct {
		name, format, raw, model, version string
		want                              Counts
		models                            int
		bad                               bool
	}{
		{"codex sample", "codex", string(codex), "openai/sol", "", Counts{Input: 15747, Output: 5, Total: 15752}, 1, false},
		{"multi turn", "codex", string(codex) + turn, "openai/sol", "", Counts{Input: 15807, CacheRead: 40, CacheWrite: 10, Output: 25, Reasoning: 5, Total: 15882}, 1, false},
		{"claude sample", "claude", string(claude), "hold/model", "", Counts{Input: 2, CacheRead: 9444, CacheWrite: 10226, Output: 4, Total: 19676}, 1, false},
		{"multi model", "claude", multi, "hold/model", "", Counts{Input: 12, CacheRead: 9474, CacheWrite: 10266, Output: 24, Reasoning: 5, Total: 19776}, 2, false},
		{"stream", "claude", `{"type":"system","claude_code_version":"2.1.332"}` + "\n" + `{"type":"assistant"}` + "\n" + strings.Join(strings.Fields(string(claude)), ""), "hold/model", "2.1.332", Counts{Input: 2, CacheRead: 9444, CacheWrite: 10226, Output: 4, Total: 19676}, 1, false},
		{"fallback", "claude", `{"type":"result","usage":{"input_tokens":2,"cache_read_input_tokens":3,"cache_creation_input_tokens":4,"output_tokens":5}}`, "hold/model", "", Counts{Input: 2, CacheRead: 3, CacheWrite: 4, Output: 5, Total: 14}, 1, false},
		{"unknown", "other", string(codex), "m", "", Counts{}, 0, true},
		{"malformed", "codex", string(codex) + "{", "m", "", Counts{}, 0, true},
		{"malformed claude", "claude", `{"type":"result",`, "m", "", Counts{}, 0, true},
		{"no usage", "codex", `{"type":"turn.started"}`, "m", "", Counts{}, 0, true},
		{"no claude usage", "claude", `{"type":"result"}`, "m", "", Counts{}, 0, true},
		{"missing", "codex", strings.Replace(turn, `"output_tokens":20,`, "", 1), "m", "", Counts{}, 0, true},
		{"missing claude", "claude", strings.Replace(string(claude), `"inputTokens": 2,`, "", 1), "m", "", Counts{}, 0, true},
		{"reasoning exceeds output", "codex", strings.Replace(turn, `"reasoning_output_tokens":5`, `"reasoning_output_tokens":21`, 1), "m", "", Counts{}, 0, true},
		{"duplicate result", "claude", string(claude) + string(claude), "m", "", Counts{}, 0, true},
		{"negative", "codex", strings.Replace(turn, `:100`, `:-1`, 1), "m", "", Counts{}, 0, true},
		{"cache exceeds input", "codex", strings.Replace(turn, `:100`, `:1`, 1), "m", "", Counts{}, 0, true},
		{"fraction", "codex", strings.Replace(turn, `:100`, `:1.5`, 1), "m", "", Counts{}, 0, true},
		{"overflow wrapping positive", "codex", `{"type":"turn.completed","usage":{"input_tokens":9223372036854775807,"cached_input_tokens":0,"cache_write_input_tokens":9223372036854775807,"output_tokens":3,"reasoning_output_tokens":0}}`, "m", "", Counts{}, 0, true},
		{"overflow", "codex", strings.Replace(turn, `:100`, `:9223372036854775807`, 1), "m", "", Counts{}, 0, true},
		{"null", "codex", strings.Replace(turn, `:100`, `:null`, 1), "m", "", Counts{}, 0, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := Parse(tc.format, []byte(tc.raw), tc.model)
			if tc.bad {
				if tc.name == "unknown" && (err == nil || !strings.Contains(err.Error(), "unknown usage format")) {
					t.Fatalf("unclear unknown-format error: %v", err)
				}
				if err == nil {
					t.Fatalf("accepted %s", tc.raw)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if tc.format == "codex" || tc.name == "fallback" {
				if got.Models[tc.model] != tc.want {
					t.Fatalf("hold model attribution: %v", got.Models)
				}
			} else if _, ok := got.Models["claude-opus-5-5"]; !ok {
				t.Fatalf("reported model attribution: %v", got.Models)
			}
			if got.Counts != tc.want || len(got.Models) != tc.models || got.ToolVersion != tc.version || got.ParserRevision != tc.format+"/1" || len(got.Reported) == 0 {
				t.Fatalf("parsed %#v, want %#v", got, tc.want)
			}
		})
	}
}

func TestParseRequiresEachCount(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		format, raw string
		fields      []string
	}{
		{"codex", `{"input_tokens":1,"cached_input_tokens":0,"cache_write_input_tokens":0,"output_tokens":1,"reasoning_output_tokens":0}`, []string{"input_tokens", "cached_input_tokens", "cache_write_input_tokens", "output_tokens", "reasoning_output_tokens"}},
		{"claude", `{"input_tokens":1,"cache_read_input_tokens":0,"cache_creation_input_tokens":0,"output_tokens":1}`, []string{"input_tokens", "cache_read_input_tokens", "cache_creation_input_tokens", "output_tokens"}},
		{"claude-model", `{"inputTokens":1,"cacheReadInputTokens":0,"cacheCreationInputTokens":0,"outputTokens":1}`, []string{"inputTokens", "cacheReadInputTokens", "cacheCreationInputTokens", "outputTokens"}},
	} {
		for _, field := range tc.fields {
			for _, mode := range []string{"missing", "null"} {
				t.Run(tc.format+"/"+field+"/"+mode, func(t *testing.T) {
					t.Parallel()
					var fields map[string]json.RawMessage
					if err := json.Unmarshal([]byte(tc.raw), &fields); err != nil {
						t.Fatal(err)
					}
					if mode == "missing" {
						delete(fields, field)
					} else {
						fields[field] = json.RawMessage("null")
					}
					usage, err := json.Marshal(fields)
					if err != nil {
						t.Fatal(err)
					}
					format := tc.format
					raw := `{"type":"turn.completed","usage":` + string(usage) + `}`
					switch format {
					case "claude":
						raw = `{"type":"result","usage":` + string(usage) + `}`
					case "claude-model":
						format = "claude"
						raw = `{"type":"result","modelUsage":{"model":` + string(usage) + `}}`
					}
					if _, err := Parse(format, []byte(raw), "hold/model"); err == nil {
						t.Fatalf("missing %s was zero-filled", field)
					}
				})
			}
		}
	}
}

func TestReportedClaudeModelsRoundTrip(t *testing.T) {
	t.Parallel()
	raw := `{"type":"result","modelUsage":{"input_tokens":{"inputTokens":1,"outputTokens":2,"cacheReadInputTokens":3,"cacheCreationInputTokens":4}}}`
	parsed, err := Parse("claude", []byte(raw), "hold/model")
	if err != nil {
		t.Fatal(err)
	}
	r := Record{ID: "USG-test", TaskID: "FID-test", HoldID: "HLD-test", Actor: "agent:test", Role: "build", RecordedAt: "2026-11-27T00:00:00Z", Parsed: parsed}
	if err := r.Validate(); err != nil {
		t.Fatalf("reported model name must not select aggregate usage: %v", err)
	}
}
