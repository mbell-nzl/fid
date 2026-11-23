package gittest

import (
	"os"
	"path/filepath"
)

// BlockGH prevents tests and their fid children from reaching a real gh.
// Individual transport tests put their own fake earlier on PATH.
func BlockGH() (func(), error) {
	dir, err := os.MkdirTemp("", "fid-test-gh-")
	if err != nil {
		return nil, err
	}
	cleanup := func() { _ = os.RemoveAll(dir) }
	if err := os.WriteFile(filepath.Join(dir, "gh"), []byte("#!/bin/sh\nexit 1\n"), 0o700); err != nil {
		cleanup()
		return nil, err
	}
	if err := os.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH")); err != nil {
		cleanup()
		return nil, err
	}
	return cleanup, nil
}
