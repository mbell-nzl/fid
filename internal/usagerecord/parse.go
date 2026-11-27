// Package usagerecord defines session usage and reads agent output without
// consulting a ledger. Parser revisions name the rules used for old receipts.
package usagerecord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
)

type Counts struct {
	Input      int64 `json:"uncached_input"`
	CacheRead  int64 `json:"cache_read"`
	CacheWrite int64 `json:"cache_write"`
	Output     int64 `json:"output"`
	Reasoning  int64 `json:"reasoning"`
	Total      int64 `json:"total"`
}

type Parsed struct {
	Format         string            `json:"format"`
	ParserRevision string            `json:"parser_revision"`
	ToolVersion    string            `json:"tool_version"`
	VersionSource  string            `json:"version_source"`
	Counts         Counts            `json:"counts"`
	Models         map[string]Counts `json:"models"`
	Reported       []json.RawMessage `json:"reported"`
}

func add(a, b Counts) (Counts, error) {
	pairs := [][2]int64{{a.Input, b.Input}, {a.CacheRead, b.CacheRead}, {a.CacheWrite, b.CacheWrite}, {a.Output, b.Output}, {a.Reasoning, b.Reasoning}, {a.Total, b.Total}}
	values := make([]int64, len(pairs))
	for i, p := range pairs {
		if p[0] < 0 || p[1] < 0 || p[0] > math.MaxInt64-p[1] {
			return Counts{}, fmt.Errorf("token count overflow or negative count")
		}
		values[i] = p[0] + p[1]
	}
	return Counts{values[0], values[1], values[2], values[3], values[4], values[5]}, nil
}

// Parse accepts JSON values separated by whitespace (including JSONL). Every
// value must be an object, even when its event type is not used for usage.
func Parse(format string, raw []byte, model string) (Parsed, error) {
	p := Parsed{Format: format, ParserRevision: format + "/1", Models: map[string]Counts{}}
	if format != "codex" && format != "claude" {
		return Parsed{}, fmt.Errorf("unknown usage format %q (want codex or claude)", format)
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	var result map[string]json.RawMessage
	for {
		var event map[string]json.RawMessage
		if err := dec.Decode(&event); err == io.EOF {
			break
		} else if err != nil {
			return Parsed{}, fmt.Errorf("malformed %s output: %w", format, err)
		}
		if event == nil {
			return Parsed{}, fmt.Errorf("%s output requires JSON objects", format)
		}
		var kind string
		if err := json.Unmarshal(event["type"], &kind); err != nil {
			return Parsed{}, fmt.Errorf("output event needs a type: %w", err)
		}
		if format == "codex" && kind == "turn.completed" {
			counts, err := codexCounts(event["usage"])
			if err != nil {
				return Parsed{}, err
			}
			if err = p.accumulate(model, counts); err != nil {
				return Parsed{}, err
			}
			p.Reported = append(p.Reported, event["usage"])
		}
		if format == "claude" {
			if kind == "system" && event["claude_code_version"] != nil {
				if err := json.Unmarshal(event["claude_code_version"], &p.ToolVersion); err != nil || p.ToolVersion == "" {
					return Parsed{}, fmt.Errorf("invalid claude_code_version")
				}
				p.VersionSource = "output"
			}
			if kind == "result" {
				if result != nil {
					return Parsed{}, fmt.Errorf("multiple Claude results")
				}
				result = event
			}
		}
	}
	if format == "claude" && result != nil {
		if err := p.claudeResult(result, model); err != nil {
			return Parsed{}, err
		}
	}
	if len(p.Reported) == 0 {
		return Parsed{}, fmt.Errorf("%s output has no usage", format)
	}
	return p, nil
}

func (p *Parsed) accumulate(model string, c Counts) error {
	var err error
	p.Counts, err = add(p.Counts, c)
	if err != nil {
		return err
	}
	p.Models[model], err = add(p.Models[model], c)
	return err
}

func numbers(raw json.RawMessage, names ...string) ([]int64, error) {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(raw, &fields); err != nil || fields == nil {
		return nil, fmt.Errorf("usage requires an object")
	}
	counts := make([]int64, len(names))
	for i, name := range names {
		value, ok := fields[name]
		if !ok || bytes.Equal(value, []byte("null")) {
			return nil, fmt.Errorf("usage missing required field %s", name)
		}
		if err := json.Unmarshal(value, &counts[i]); err != nil || counts[i] < 0 {
			return nil, fmt.Errorf("usage field %s requires a nonnegative integer", name)
		}
	}
	return counts, nil
}

func counted(input, read, write, output, reasoning int64) (Counts, error) {
	c := Counts{Input: input, CacheRead: read, CacheWrite: write, Output: output, Reasoning: reasoning}
	if reasoning > output {
		return Counts{}, fmt.Errorf("reasoning count exceeds output")
	}
	for _, n := range []int64{input, read, write, output} {
		if n < 0 || c.Total > math.MaxInt64-n {
			return Counts{}, fmt.Errorf("token count overflow or negative count")
		}
		c.Total += n
	}
	return c, nil
}
func codexCounts(raw json.RawMessage) (Counts, error) {
	n, err := numbers(raw, "input_tokens", "cached_input_tokens", "cache_write_input_tokens", "output_tokens", "reasoning_output_tokens")
	if err != nil {
		return Counts{}, err
	}
	// Codex input includes cache reads; cache writes are reported separately.
	return counted(n[0]-n[1], n[1], n[2], n[3], n[4])
}
func claudeCounts(raw json.RawMessage, perModel bool) (Counts, error) {
	names := []string{"input_tokens", "cache_read_input_tokens", "cache_creation_input_tokens", "output_tokens"}
	if perModel {
		names = []string{"inputTokens", "cacheReadInputTokens", "cacheCreationInputTokens", "outputTokens"}
	}
	n, err := numbers(raw, names...)
	if err != nil {
		return Counts{}, err
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(raw, &fields); err != nil {
		return Counts{}, err
	}
	reasoning := int64(0)
	if perModel {
		if r, ok := fields["thinkingTokens"]; ok {
			values, err := numbers(json.RawMessage(`{"r":`+string(r)+`}`), "r")
			if err != nil {
				return Counts{}, err
			}
			reasoning = values[0]
		}
	} else if r, ok := fields["output_tokens_details"]; ok {
		values, err := numbers(r, "thinking_tokens")
		if err != nil {
			return Counts{}, err
		}
		reasoning = values[0]
	}
	return counted(n[0], n[1], n[2], n[3], reasoning)
}
func (p *Parsed) claudeResult(result map[string]json.RawMessage, model string) error {
	if raw, ok := result["modelUsage"]; ok {
		var models map[string]json.RawMessage
		if err := json.Unmarshal(raw, &models); err != nil || len(models) == 0 {
			return fmt.Errorf("modelUsage requires a nonempty object")
		}
		for name, raw := range models {
			if name == "" {
				return fmt.Errorf("modelUsage needs model names")
			}
			counts, err := claudeCounts(raw, true)
			if err != nil {
				return err
			}
			if err = p.accumulate(name, counts); err != nil {
				return err
			}
		}
		p.Reported = append(p.Reported, raw)
	} else {
		raw := result["usage"]
		counts, err := claudeCounts(raw, false)
		if err != nil {
			return err
		}
		if err = p.accumulate(model, counts); err != nil {
			return err
		}
		p.Reported = append(p.Reported, raw)
	}
	// Keep both the aggregate usage and per-model report when supplied.
	if result["modelUsage"] != nil && result["usage"] != nil {
		p.Reported = append(p.Reported, result["usage"])
	}
	return nil
}
