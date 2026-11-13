package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const ProjectIdentityFileName = ".fid.toml"

// ReadProjectIdentity reads the committed identity used to key backups. It
// deliberately does not change the older ledger config's opaque-ID contract.
func ReadProjectIdentity(root string) (string, error) {
	path := filepath.Join(root, ProjectIdentityFileName)
	info, err := os.Lstat(path)
	if err == nil && !info.Mode().IsRegular() {
		err = fmt.Errorf("expected a regular file")
	}
	var data []byte
	if err == nil {
		data, err = os.ReadFile(path)
	}
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", err
		}
		return "", &InvalidConfigError{Path: path, Err: err}
	}
	document, err := DecodeProject(data)
	if err != nil {
		return "", &InvalidConfigError{Path: path, Err: err}
	}
	if !validBackupProjectID(document.ProjectID) {
		return "", &InvalidConfigError{Path: path, Err: fmt.Errorf("read config %q: project.id must be project- followed by 32 lowercase hexadecimal digits", path)}
	}
	return document.ProjectID, nil
}

func validBackupProjectID(id string) bool {
	hex, ok := strings.CutPrefix(id, "project-")
	if !ok || len(hex) != 32 {
		return false
	}
	for _, character := range hex {
		if (character < '0' || character > '9') && (character < 'a' || character > 'f') {
			return false
		}
	}
	return true
}

// EnsureProjectIdentity keeps an existing valid file byte for byte, and
// refuses an identity that changed while init was staging its ledger.
func EnsureProjectIdentity(root, id string) error {
	existing, err := ReadProjectIdentity(root)
	if errors.Is(err, os.ErrNotExist) {
		path := filepath.Join(root, ProjectIdentityFileName)
		if err := writeStringConfig(path, stringConfigContents("project", "id", id), 0o644, false); err != nil && !errors.Is(err, os.ErrExist) {
			return err
		}
		existing, err = ReadProjectIdentity(root)
	}
	if err != nil {
		return err
	}
	if existing != id {
		return fmt.Errorf("project identity in %q changed during initialization; retry init", filepath.Join(root, ProjectIdentityFileName))
	}
	return nil
}
