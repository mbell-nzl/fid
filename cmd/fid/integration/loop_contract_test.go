//go:build integration

package integration

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"
)

type loopDeclaration struct {
	Version string              `json:"contract_version"`
	Kinds   []string            `json:"next_kinds"`
	Waiting []string            `json:"waiting_kinds"`
	Fields  map[string][]string `json:"fields"`
}

func readLoopDeclaration(t *testing.T) loopDeclaration {
	t.Helper()
	raw := readP21File(t, "../../../docs/specs/LOOP_CONTRACT.md")
	_, block, found := strings.Cut(string(raw), "```json loop-contract\n")
	if !found {
		t.Fatal("missing machine-readable contract declaration")
	}
	block, _, found = strings.Cut(block, "\n```")
	if !found {
		t.Fatal("unclosed contract declaration")
	}
	var d loopDeclaration
	if err := json.Unmarshal([]byte(block), &d); err != nil {
		t.Fatal(err)
	}
	return d
}

// Select declared paths strictly: a missing required key fails even if another
// case has a null parent. ? marks an optional field; [] maps an array.
func loopField(t *testing.T, value any, path []string) any {
	t.Helper()
	if len(path) == 0 || value == nil {
		return value
	}
	key := path[0]
	if strings.HasSuffix(key, "[]") {
		items := loopField(t, value, []string{strings.TrimSuffix(key, "[]")})
		if items == nil && strings.HasSuffix(strings.TrimSuffix(key, "[]"), "?") {
			return nil
		}
		array, ok := items.([]any)
		if !ok {
			t.Fatalf("%s = %v, want array", key, items)
		}
		result := make([]any, 0, len(array))
		for _, item := range array {
			result = append(result, loopField(t, item, path[1:]))
		}
		return result
	}
	object, ok := value.(map[string]any)
	if !ok {
		t.Fatalf("%s traverses %v, want object", key, value)
	}
	item, found := object[strings.TrimSuffix(key, "?")]
	if !found && !strings.HasSuffix(key, "?") {
		t.Fatalf("missing declared field %s in %v", key, object)
	}
	return loopField(t, item, path[1:])
}

var loopVolatile = regexp.MustCompile(`^(?:(?:FID|HLD|RCL|USG|HOF|CHK|RND|FND|FRS|QST|PDN|CMP|AC)-[0-9a-f]{6}(?:-[0-9]+)?|[0-9a-f]{40}|\d{4}-\d\d-\d\dT\d\d:\d\d:\d\d(?:\.\d+)?Z)$`)
var loopEmbeddedFinding = regexp.MustCompile(`RND-[0-9a-f]{6}(?:-[0-9]+)?`)

func loopStable(value any) any {
	switch v := value.(type) {
	case string:
		if loopVolatile.MatchString(v) {
			if strings.Contains(v, "T") && strings.HasSuffix(v, "Z") {
				return "<time>"
			}
			if prefix, _, ok := strings.Cut(v, "-"); ok {
				return "<" + prefix + ">"
			}
			return "<commit>"
		}
		return loopEmbeddedFinding.ReplaceAllString(v, "<RND>")
	case []any:
		stable := make([]any, len(v))
		for i := range v {
			stable[i] = loopStable(v[i])
		}
		return stable
	}
	return value
}

type loopContract struct {
	t              *testing.T
	decl           loopDeclaration
	kinds, waiting []string
}

func (c *loopContract) observe(name string, envelope map[string]any, exit int) {
	c.t.Helper()
	if envelope["contract_version"] != c.decl.Version {
		c.t.Fatalf("contract version = %v, want %s", envelope["contract_version"], c.decl.Version)
	}
	if envelope["ok"] == true && (exit != 0 || envelope["error"] != nil) || envelope["ok"] == false && (exit == 0 || envelope["result"] != nil) {
		c.t.Fatalf("exit/envelope disagree: exit=%d, envelope=%v", exit, envelope)
	}
	got := map[string]any{"exit": exit, "contract_version": envelope["contract_version"], "ok": envelope["ok"], "operation": envelope["operation"]}
	if envelope["ok"] == true {
		operation := envelope["operation"].(string)
		fields, found := c.decl.Fields[operation]
		if !found {
			c.t.Fatalf("undeclared operation %s", operation)
		}
		projected := map[string]any{}
		for _, path := range fields {
			projected[path] = loopStable(loopField(c.t, envelope["result"], strings.Split(path, ".")))
		}
		got["result"] = projected
		if strings.HasPrefix(operation, "next.") {
			c.kinds = append(c.kinds, p21PathString(c.t, envelope, "result", "kind"))
			if w := p21Path(c.t, envelope, "result", "waiting"); w != nil {
				c.waiting = append(c.waiting, w.(map[string]any)["kind"].(string))
			}
		}
	} else {
		got["error.code"], got["error.retryable"] = loopField(c.t, envelope, []string{"error", "code"}), loopField(c.t, envelope, []string{"error", "retryable"})
	}
	path := filepath.Join("testdata", "loop-contract", name+".json")
	raw := p21JSON(c.t, got)
	if os.Getenv("UPDATE_LOOP_GOLDEN") == "1" {
		if err := os.WriteFile(path, raw, 0o600); err != nil {
			c.t.Fatal(err)
		}
	}
	if want := readP21File(c.t, path); string(want) != string(raw) {
		c.t.Fatalf("%s contract differs:\ngot %s\nwant %s", name, raw, want)
	}
}

func (c *loopContract) sets() {
	c.t.Helper()
	for _, test := range []struct {
		name      string
		got, want []string
	}{{"next.kind", c.kinds, c.decl.Kinds}, {"waiting.kind", c.waiting, c.decl.Waiting}} {
		slices.Sort(test.got)
		slices.Sort(test.want)
		if got := slices.Compact(test.got); !slices.Equal(got, test.want) {
			c.t.Errorf("produced %s = %v, declaration = %v", test.name, got, test.want)
		}
	}
}

func TestLoopContractDeclaration(t *testing.T) {
	t.Parallel()
	d := readLoopDeclaration(t)
	if d.Version != "fid-cli.v2" || len(d.Fields) != 18 {
		t.Fatalf("incomplete declaration: %+v", d)
	}
	for _, operation := range []string{"next.get", "next.claim"} {
		var allowance []string
		for _, path := range d.Fields[operation] {
			if strings.HasPrefix(path, "rework_allowance") {
				allowance = append(allowance, path)
			}
		}
		want := []string{"rework_allowance?.decision_id", "rework_allowance?.rounds", "rework_allowance?.reason", "rework_allowance?.actor", "rework_allowance?.at"}
		slices.Sort(allowance)
		slices.Sort(want)
		if !slices.Equal(allowance, want) {
			t.Errorf("%s allowance fields = %v, want %v", operation, allowance, want)
		}
	}
}
