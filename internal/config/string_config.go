package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode/utf8"
)

// readStringConfig reads the single string setting in an identity or user
// config. Like the ledger config parser, it refuses unknown and duplicate
// keys. These files support single-line TOML basic and literal strings.
func readStringConfig(path, section, key string) (string, bool, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return "", false, fmt.Errorf("inspect config %q: %w", path, err)
	}
	if !info.Mode().IsRegular() {
		return "", false, fmt.Errorf("read config %q: expected a regular file", path)
	}
	file, err := os.Open(path)
	if err != nil {
		return "", false, fmt.Errorf("open config %q: %w", path, err)
	}
	defer closeConfigFile(file)
	scanner := bufio.NewScanner(file)
	var value string
	var seenSection, seenKey bool
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "[") {
			header, _, _ := strings.Cut(line, "#")
			if strings.TrimSpace(header) != "["+section+"]" || seenSection {
				return "", false, fmt.Errorf("parse config %q: unsupported or duplicate section", path)
			}
			seenSection = true
			continue
		}
		name, raw, ok := splitAssignment(line)
		if !ok || !seenSection || name != key || seenKey {
			return "", false, fmt.Errorf("parse config %q: unsupported, duplicate or invalid key %q", path, name)
		}
		value, err = parseTOMLString(raw)
		if err != nil {
			return "", false, fmt.Errorf("parse config %q %s.%s: %w", path, section, key, err)
		}
		seenKey = true
	}
	if err := scanner.Err(); err != nil {
		return "", false, fmt.Errorf("read config %q: %w", path, err)
	}
	return value, seenKey, nil
}

func parseTOMLString(raw string) (string, error) {
	if len(raw) < 2 || (raw[0] != '"' && raw[0] != '\'') || !utf8.ValidString(raw) {
		return "", fmt.Errorf("expected a single-line quoted string")
	}
	for index := 1; index < len(raw); index++ {
		character := raw[index]
		if character == raw[0] {
			tail := strings.TrimSpace(raw[index+1:])
			if tail != "" && !strings.HasPrefix(tail, "#") {
				break
			}
			if raw[0] == '\'' {
				return raw[1:index], nil
			}
			// Go and TOML share these escapes; a raw tab is legal TOML.
			value, err := strconv.Unquote(strings.ReplaceAll(raw[:index+1], "\t", `\t`))
			if err != nil {
				return "", fmt.Errorf("invalid string: %w", err)
			}
			return value, nil
		}
		if character < 0x20 && character != '\t' || character == 0x7f {
			break
		}
		if raw[0] == '"' && character == '\\' {
			index++
			if index >= len(raw) || !strings.ContainsRune(`btnfr"\uU`, rune(raw[index])) {
				break
			}
		}
	}
	return "", fmt.Errorf("invalid single-line TOML string")
}

func stringConfigContents(section, key, value string) string {
	// JSON's string escapes are valid in TOML basic strings.
	quoted, _ := json.Marshal(value)
	return fmt.Sprintf("[%s]\n%s = %s\n", section, key, quoted)
}

// writeStringConfig publishes a complete file. A new project identity uses a
// hard link so a concurrent init cannot overwrite an identity already there.
func writeStringConfig(path, contents string, mode fs.FileMode, replace bool) error {
	file, err := os.CreateTemp(filepath.Dir(path), ".fid-config-*")
	if err != nil {
		return fmt.Errorf("stage config %q: %w", path, err)
	}
	defer func() { _ = os.Remove(file.Name()) }()
	defer closeConfigFile(file)
	if err := file.Chmod(mode); err != nil {
		return fmt.Errorf("set config permissions %q: %w", path, err)
	}
	if _, err := file.WriteString(contents); err != nil {
		return fmt.Errorf("write config %q: %w", path, err)
	}
	if err := file.Sync(); err != nil {
		return fmt.Errorf("sync config %q: %w", path, err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close config %q: %w", path, err)
	}
	if replace {
		err = os.Rename(file.Name(), path)
	} else {
		err = os.Link(file.Name(), path)
	}
	if err != nil {
		return fmt.Errorf("publish config %q: %w", path, err)
	}
	return nil
}
