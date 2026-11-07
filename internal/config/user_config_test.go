package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestUserConfigPath(t *testing.T) { //nolint:paralleltest // t.Setenv isolates HOME and XDG_CONFIG_HOME.
	for _, useXDG := range []bool{true, false} {
		t.Run(map[bool]string{true: "XDG", false: "HOME fallback"}[useXDG], func(t *testing.T) {
			home, xdg := t.TempDir(), t.TempDir()
			t.Setenv("HOME", home)
			want := filepath.Join(xdg, "fid", "config.toml")
			if !useXDG {
				xdg = ""
				want = filepath.Join(home, ".config", "fid", "config.toml")
			}
			t.Setenv("XDG_CONFIG_HOME", xdg)
			got, err := UserConfigPath()
			if err != nil || got != want {
				t.Fatalf("UserConfigPath = %q, %v, want %q", got, err, want)
			}
		})
	}
}

func TestUserConfigTargetRoundTrip(t *testing.T) {
	t.Parallel()
	for _, target := range []string{
		"git@example.test:private/ledgers.git",
		"https://example.test/ledgers.git",
		"../backups with spaces/#ledger",
		`C:\Users\mark\backups`,
		"/tmp/a\"quote\"/日本語",
		"/tmp/a\tand\nnewline",
	} {
		t.Run(target, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "fid", "config.toml")
			if got, err := ReadUserConfig(path); err != nil || got.BackupTarget != "" {
				t.Fatalf("missing config = %+v, %v", got, err)
			}
			if err := SetBackupTarget(path, target); err != nil {
				t.Fatal(err)
			}
			got, err := ReadUserConfig(path)
			if err != nil || got.BackupTarget != target {
				t.Fatalf("round trip = %+v, %v, want %q", got, err, target)
			}
			info, err := os.Stat(path)
			if err != nil || info.Mode().Perm() != 0o600 {
				t.Fatalf("config mode = %v, %v, want 0600", info, err)
			}
			entries, err := os.ReadDir(filepath.Dir(path))
			if err != nil || len(entries) != 1 || entries[0].Name() != "config.toml" {
				t.Fatalf("config directory = %v, %v, want only config.toml", entries, err)
			}
		})
	}
}

func TestUserConfigParsing(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name, contents, want string
		valid                bool
	}{
		{name: "empty file", valid: true},
		{name: "empty table", contents: "[backup]\n", valid: true},
		{name: "comments and literal", contents: "# backup\n[backup] # target\ntarget = 'relative/#path' # comment\n", want: "relative/#path", valid: true},
		{name: "escaped unicode", contents: "[backup]\ntarget = \"\\u65e5\\U0000672c\"\n", want: "日本", valid: true},
		{name: "raw tab", contents: "[backup]\ntarget = \"a\tb\"\n", want: "a\tb", valid: true},
		{name: "unknown key", contents: "[backup]\nother = 'x'\n"},
		{name: "unknown table", contents: "[other]\ntarget = 'x'\n"},
		{name: "duplicate table", contents: "[backup]\n[backup]\n"},
		{name: "duplicate key", contents: "[backup]\ntarget = 'a'\ntarget = 'b'\n"},
		{name: "no table", contents: "target = 'x'\n"},
		{name: "unquoted", contents: "[backup]\ntarget = 5\n"},
		{name: "trailing value", contents: "[backup]\ntarget = 'x' junk\n"},
		{name: "unterminated string", contents: "[backup]\ntarget = \"x\n"},
		{name: "Go escape is not TOML", contents: "[backup]\ntarget = \"\\x41\"\n"},
	} {
		t.Run(test.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "config.toml")
			if err := os.WriteFile(path, []byte(test.contents), 0o600); err != nil {
				t.Fatal(err)
			}
			got, err := ReadUserConfig(path)
			if test.valid {
				if err != nil || got.BackupTarget != test.want {
					t.Fatalf("ReadUserConfig = %+v, %v, want %q", got, err, test.want)
				}
			} else if err == nil || !strings.Contains(err.Error(), path) {
				t.Fatalf("ReadUserConfig = %+v, %v, want a refusal naming %q", got, err, path)
			}
		})
	}
}
