package policy

import (
	"reflect"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
)

func TestFloorRaiseUsesHigherLaneIncludingUnplannedTasks(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		floor, planned, want string
	}{
		{"", "", ""},
		{"", "routine", "routine"},
		{"", "standard", "standard"},
		{"", "assured", "assured"},
		{"routine", "", "routine"},
		{"standard", "", "standard"},
		{"assured", "", "assured"},
		{"routine", "routine", "routine"},
		{"standard", "routine", "standard"},
		{"assured", "routine", "assured"},
		{"routine", "standard", "standard"},
		{"standard", "standard", "standard"},
		{"assured", "standard", "assured"},
		{"routine", "assured", "assured"},
		{"standard", "assured", "assured"},
		{"assured", "assured", "assured"},
	} {
		t.Run(tc.floor+"/"+tc.planned, func(t *testing.T) {
			t.Parallel()
			floor := Floor{}
			if tc.floor != "" {
				floor.Lane = &tc.floor
			}
			if got := floor.Raise(tc.planned); got != tc.want {
				t.Fatalf("raise %q with floor %q = %q, want %q", tc.planned, tc.floor, got, tc.want)
			}
		})
	}
}

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
		{"workflow newline hard floor", []string{".github/workflows/odd\n.yml"}, "assured", []string{".github/workflows/odd\n.yml"}, []string{".github/workflows/odd\n.yml"}},
		{"workflow backslash hard floor", []string{".github/workflows/odd\\.yml"}, "assured", []string{".github/workflows/odd\\.yml"}, []string{".github/workflows/odd\\.yml"}},
		{"workflow invalid UTF-8 hard floor", []string{".github/workflows/odd\xff.yml"}, "assured", []string{".github/workflows/odd\ufffd.yml"}, []string{".github/workflows/odd\ufffd.yml"}},
		{"nested workflow unsupported name hard floor", []string{".github/workflows/odd\n/nested.yml"}, "assured", []string{".github/workflows/odd\n/nested.yml"}, []string{".github/workflows/odd\n/nested.yml"}},
		{"workflow directory is ordinary", []string{".github/workflows-other/a.yml"}, "standard", []string{".github/workflows-other/a.yml"}, []string{}},
		{"unsupported workflow lookalikes", []string{".github/workflows-other/odd\n.yml", ".github/workflows\\odd.yml", ".github/workflows\xff/odd.yml"}, "assured", []string{".github/workflows-other/odd\n.yml", ".github/workflows\\odd.yml", ".github/workflows\ufffd/odd.yml"}, []string{}},
		{"unsupported policy lookalikes", []string{".fid.toml\n", ".fid.toml\xff"}, "assured", []string{".fid.toml\n", ".fid.toml\ufffd"}, []string{}},
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

func TestFloorRetainsRecordedHardPatterns(t *testing.T) {
	t.Parallel()
	file, err := config.DecodeProject(nil)
	if err != nil {
		t.Fatal(err)
	}
	effective := Resolve(file)
	effective.HardAssured = []string{"legacy-policy/**"}
	got, err := ComputeFloor(effective, []string{".github/workflows/odd\n.yml", "legacy-policy/check.yml"})
	if err != nil {
		t.Fatal(err)
	}
	if got.HardFloor == nil || *got.HardFloor != "assured" || !reflect.DeepEqual(got.HardPaths, []string{"legacy-policy/check.yml"}) {
		t.Fatalf("floor ignored recorded hard patterns: %+v", got)
	}
}
