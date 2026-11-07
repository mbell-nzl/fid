package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

type UserConfig struct {
	BackupTarget string
}

func UserConfigPath() (string, error) {
	directory := os.Getenv("XDG_CONFIG_HOME")
	if directory == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("locate user config: %w", err)
		}
		directory = filepath.Join(home, ".config")
	}
	if !filepath.IsAbs(directory) {
		return "", fmt.Errorf("locate user config: config directory must be absolute")
	}
	return filepath.Join(directory, "fid", ConfigFileName), nil
}

func ReadUserConfig(path string) (UserConfig, error) {
	target, _, err := readStringConfig(path, "backup", "target")
	if errors.Is(err, os.ErrNotExist) {
		return UserConfig{}, nil
	}
	if err != nil {
		return UserConfig{}, &InvalidConfigError{Path: path, Err: err}
	}
	return UserConfig{BackupTarget: target}, nil
}

// SetBackupTarget stores the remote verbatim: it can be any git URL or local
// path. No git operation, expansion or existence check is performed here.
func SetBackupTarget(path, target string) error {
	if !utf8.ValidString(target) {
		return fmt.Errorf("backup target must be valid UTF-8")
	}
	if strings.TrimSpace(target) == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("create user config directory: %w", err)
	}
	return writeStringConfig(path, stringConfigContents("backup", "target", target), 0o600, true)
}
