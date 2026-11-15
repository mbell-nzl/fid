package github

import (
	"testing"
)

func TestRepositoryIdentity(t *testing.T) {
	t.Parallel()
	for _, remote := range []string{
		"https://github.com/mbell-nzl/fid", "https://GITHUB.COM/MBELL-NZL/FID.git",
		"git@github.com:mbell-nzl/fid.git", "ssh://git@GitHub.com/MBELL-NZL/fid.git",
	} {
		t.Run(remote, func(t *testing.T) {
			t.Parallel()
			got, err := ParseOrigin(remote)
			if err != nil || got.String() != "mbell-nzl/fid" {
				t.Fatalf("ParseOrigin = %v, %v", got, err)
			}
			pr, err := ParsePRURL("https://GITHUB.COM/MBELL-NZL/FID/pull/74")
			if err != nil || pr.Number != 74 || !got.Equal(pr.Repository) {
				t.Fatalf("ParsePRURL = %#v, %v", pr, err)
			}
		})
	}
}

func TestInvalidRepositoryIdentity(t *testing.T) {
	t.Parallel()
	for _, value := range []string{"", "/tmp/fid", "https://example.com/owner/fid", "https://github.com/a/b/extra", "https://github.com/a/b?x=y", "git@github.com:a", "https://github.com/a/..", "https://github.com/a/b#part"} {
		t.Run(value, func(t *testing.T) {
			t.Parallel()
			if got, err := ParseOrigin(value); err == nil {
				t.Fatalf("accepted origin %q: %v", value, got)
			}
		})
	}
	for _, value := range []string{"https://github.com/a/b/pull/0", "https://example.com/a/b/pull/1", "https://github.com/a/b/pull/1?x=y", "https://github.com/a/b/pull/9999999999999999999999999"} {
		t.Run(value, func(t *testing.T) {
			t.Parallel()
			if got, err := ParsePRURL(value); err == nil {
				t.Fatalf("accepted PR %q: %#v", value, got)
			}
		})
	}
}
