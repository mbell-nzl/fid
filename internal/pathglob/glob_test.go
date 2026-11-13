package pathglob

import "testing"

func TestGlob(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		pattern, path string
		match         bool
	}{
		{"Makefile", "Makefile", true}, {"Makefile", "src/Makefile", false},
		{"*.go", "main.go", true}, {"*.go", "src/main.go", false},
		{"docs/*", "docs/a.md", true}, {"docs/*", "docs/sub/a.md", false},
		{"docs/**", "docs/sub/a.md", true}, {"docs/**", "other/docs/a.md", false},
		{"**/*.md", "README.md", true}, {"**/*.md", "docs/sub/readme.md", true},
		{"**/testdata/**", "testdata/a", true}, {"**/testdata/**", "internal/a/testdata/b/c", true},
		{"a/**/b", "a/b", true}, {"a/**/b", "a/x/y/b", true}, {"a/**/b", "a/x/b/c", false},
		{".github/workflows/**", ".github/workflows/ci.yml", true},
		{"file.+", "file.+", true}, {"file.+", "filexxx", false},
		{"folder space/**", "folder space/日本語.md", true},
		{"**", "../outside", false}, {"**", "/outside", false}, {"**", "", false},
	} {
		t.Run(test.pattern+"_"+test.path, func(t *testing.T) {
			glob, err := Compile(test.pattern)
			if err != nil {
				t.Fatal(err)
			}
			if got := glob.Match(test.path); got != test.match {
				t.Errorf("%q matches %q = %v, want %v", test.pattern, test.path, got, test.match)
			}
		})
	}
}

func TestInvalidGlob(t *testing.T) {
	t.Parallel()
	for _, pattern := range []string{"", "/a", "../a", "a/../b", "./a", "a//b", "a/", `a\b`, "a[bc]", "a?", "a{b,c}", "***", "a\x00b"} {
		t.Run(pattern, func(t *testing.T) {
			if _, err := Compile(pattern); err == nil {
				t.Errorf("accepted %q", pattern)
			}
		})
	}
}
