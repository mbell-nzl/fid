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
	return UserConfigPathFrom(os.Getenv)
}

// UserConfigPathFrom is UserConfigPath with getenv in place of the process
// environment.
func UserConfigPathFrom(getenv func(string) string) (string, error) {
	directory, err := userDirectory(getenv, "XDG_CONFIG_HOME", ".config")
	if err != nil {
		return "", fmt.Errorf("locate user config: %w", err)
	}
	return filepath.Join(directory, "fid", ConfigFileName), nil
}

// UserCacheDirFrom returns fid's user cache directory, $XDG_CACHE_HOME/fid
// (~/.cache/fid by default), looking variables up with getenv.
func UserCacheDirFrom(getenv func(string) string) (string, error) {
	directory, err := userDirectory(getenv, "XDG_CACHE_HOME", ".cache")
	if err != nil {
		return "", fmt.Errorf("locate user cache: %w", err)
	}
	return filepath.Join(directory, "fid"), nil
}

// userDirectory returns the XDG base directory named by variable, or
// $HOME/fallback when it is unset or empty. It must be absolute.
func userDirectory(getenv func(string) string, variable, fallback string) (string, error) {
	directory := getenv(variable)
	if directory == "" {
		variable = "HOME"
		home := getenv(variable)
		if home == "" {
			return "", fmt.Errorf("$HOME is not defined")
		}
		directory = filepath.Join(home, fallback)
	}
	if !filepath.IsAbs(directory) {
		return "", fmt.Errorf("$%s must be an absolute path", variable)
	}
	return directory, nil
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
