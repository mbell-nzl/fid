package reviewcomment

import "testing"

func TestCodeSpan(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ value, want string }{
		{"FID-123456", "`FID-123456`"},
		{"a`b", "``a`b``"},
		{"a``b`c", "```a``b`c```"},
		{"`start", "`` `start ``"},
		{"end`", "`` end` ``"},
		{"``", "``` `` ```"},
		{" a ", "`  a  `"},
		{"  ", "`  `"},
	} {
		t.Run(tc.value, func(t *testing.T) {
			t.Parallel()
			if got := codeSpan(tc.value); got != tc.want {
				t.Fatalf("codeSpan(%q)=%q, want %q", tc.value, got, tc.want)
			}
		})
	}
}
