package policy

import (
	"reflect"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
)

func TestFloorUsesHighestChangedPathAndOnlyBuiltInHardFloors(t *testing.T) {
	t.Parallel()
	file, err := config.DecodeProject([]byte(`[lanes.routine]
paths=["docs/**", "**/*.md"]
[lanes.assured]
paths=["internal/trust/**", "docs/security/**"]
[checks.test]
run=["true"]
timeout="1m"
defined_by=["Makefile", ".github/scripts/**"]
`))
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name    string
		paths   []string
		lane    string
		setters []string
		hard    []string
	}{
		{"no changes", nil, "", []string{}, []string{}},
		{"routine", []string{"docs/a.md", "README.md"}, "routine", []string{"README.md", "docs/a.md"}, []string{}},
		{"standard fallback", []string{"docs/a.md", "main.go"}, "standard", []string{"main.go"}, []string{}},
		{"assured wins overlap", []string{"docs/security/a.md", "internal/trust/a.go"}, "assured", []string{"docs/security/a.md", "internal/trust/a.go"}, []string{}},
		{"check definition", []string{"Makefile", ".github/scripts/check.py"}, "assured", []string{".github/scripts/check.py", "Makefile"}, []string{}},
		{"policy hard floor", []string{".fid.toml", "internal/trust/a.go"}, "assured", []string{".fid.toml", "internal/trust/a.go"}, []string{".fid.toml"}},
		{"workflow hard floor", []string{".github/workflows/test.yml"}, "assured", []string{".github/workflows/test.yml"}, []string{".github/workflows/test.yml"}},
		{"workflow directory is ordinary", []string{".github/workflows-other/a.yml"}, "standard", []string{".github/workflows-other/a.yml"}, []string{}},
		{"unsupported paths", []string{"a\\b.md", "docs/a\nb.md", "docs/\xff.md"}, "assured", []string{"a\\b.md", "docs/a\nb.md", "docs/\ufffd.md"}, []string{}},
		{"stable unique setters", []string{"z.go", "a.go", "z.go"}, "standard", []string{"a.go", "z.go"}, []string{}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := ComputeFloor(Resolve(file), tc.paths)
			if err != nil {
				t.Fatal(err)
			}
			lane, hard := "", ""
			if got.Lane != nil {
				lane = *got.Lane
			}
			if got.HardFloor != nil {
				hard = *got.HardFloor
			}
			wantHard := ""
			if len(tc.hard) > 0 {
				wantHard = "assured"
			}
			if lane != tc.lane || hard != wantHard || !reflect.DeepEqual(got.Paths, tc.setters) || !reflect.DeepEqual(got.HardPaths, tc.hard) {
				t.Fatalf("floor=%+v lane=%s hard=%s, want %s/%s %q/%q", got, lane, hard, tc.lane, wantHard, tc.setters, tc.hard)
			}
			if got.Raise("assured") != "assured" {
				t.Fatal("floor lowered the recorded lane")
			}
		})
	}
}
