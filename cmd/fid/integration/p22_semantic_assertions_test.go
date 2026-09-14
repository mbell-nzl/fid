//go:build integration

package integration

import (
	"encoding/json"
	"strconv"
	"strings"
	"testing"
)

func p22ResultPath(result map[string]any, path string) any {
	var value any = result
	for _, field := range strings.Split(path, ".") {
		if list, ok := value.([]any); ok {
			if field == "#" {
				value = float64(len(list))
				continue
			}
			index, err := strconv.Atoi(field)
			if err != nil || index < 0 || index >= len(list) {
				return nil
			}
			value = list[index]
			continue
		}
		object, ok := value.(map[string]any)
		if !ok {
			return nil
		}
		value = object[field]
	}
	return value
}

func assertP22Replay(t *testing.T, fixture p22ConformanceCase, original map[string]any, response p22ConformanceResponse) {
	t.Helper()
	if response.ExitCode != 0 || response.Stderr != "" {
		t.Fatalf("%s same-key replay exit=%d stderr=%q", fixture.Operation, response.ExitCode, response.Stderr)
	}
	var replay map[string]any
	if err := json.Unmarshal([]byte(response.Stdout), &replay); err != nil {
		t.Fatal(err)
	}
	result, ok := replay["result"].(map[string]any)
	if !ok || result["replayed"] != true {
		t.Fatalf("%s replay = %#v, want boolean true", fixture.Operation, replay)
	}
	result["replayed"] = false
	if !p22SameResult(original, replay) {
		t.Fatalf("%s replay changed durable result/event: first=%#v replay=%#v", fixture.Operation, original, replay)
	}
}
