package journal

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"path"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

func validateCurrentValue(kind string, value any) error {
	if strings.HasPrefix(kind, "?") {
		if value == nil {
			return nil
		}
		return validateCurrentValue(kind[1:], value)
	}
	if strings.HasPrefix(kind, "[]") {
		values, ok := value.([]any)
		if !ok || len(values) > 10000 {
			return fmt.Errorf("%s requires a bounded array", kind)
		}
		for i, item := range values {
			if err := validateCurrentValue(kind[2:], item); err != nil {
				return fmt.Errorf("%s[%d]: %w", kind, i, err)
			}
		}
		return nil
	}
	if strings.HasPrefix(kind, "=") {
		if !currentEnum(value, strings.Split(kind[1:], ",")...) {
			return fmt.Errorf("invalid closed enum %s", kind)
		}
		return nil
	}
	if shape, ok := currentValueShapes[kind]; ok {
		_, err := validateCurrentShape(value, shape)
		if err != nil {
			return fmt.Errorf("%s: %w", kind, err)
		}
		return nil
	}
	valid := false
	text, stringValue := value.(string)
	switch kind {
	case "Text":
		valid = stringValue && utf8.ValidString(text) && len(text) <= 4096
	case "Nonempty":
		valid = stringValue && strings.TrimSpace(text) != "" && utf8.ValidString(text) && len(text) <= 4096
	case "ID":
		valid = currentID(value)
	case "Bool":
		_, valid = value.(bool)
	case "False":
		valid = value == false
	case "Null":
		valid = value == nil
	case "Integer", "Positive", "Nonnegative", "One":
		number, ok := currentIntegerValue(value)
		valid = ok
		if kind == "Positive" {
			valid = ok && number >= 1
		}
		if kind == "Nonnegative" {
			valid = ok && number >= 0
		}
		if kind == "One" {
			valid = ok && number == 1
		}
	case "Timestamp":
		valid = currentTimestamp(value)
	case "Digest":
		valid = stringValue && strings.HasPrefix(text, "sha256:") && currentLowerHex(strings.TrimPrefix(text, "sha256:"), 64)
	case "Hash":
		valid = stringValue && (currentLowerHex(text, 64) || (strings.HasPrefix(text, "sha256:") && currentLowerHex(text[7:], 64)))
	case "OID":
		valid = stringValue && currentLowerHex(text, 40)
	case "Path":
		valid = currentPath(value)
	case "Cwd":
		valid = value == "." || currentPath(value)
	case "EnvironmentName":
		valid = stringValue && currentEnvironmentName.MatchString(text)
	default:
		return fmt.Errorf("unregistered current value type %q", kind)
	}
	if !valid {
		return fmt.Errorf("invalid %s", kind)
	}
	return nil
}

var currentEnvironmentName = regexp.MustCompile(`^[A-Z_][A-Z0-9_]{0,127}$`)

func validateCurrentShape(value any, shape string) (map[string]any, error) {
	object, ok := value.(map[string]any)
	fields := strings.Fields(shape)
	if !ok || len(object) != len(fields) {
		return nil, fmt.Errorf("closed object has wrong shape")
	}
	for _, field := range fields {
		name, kind, _ := strings.Cut(field, ":")
		v, found := object[name]
		if !found {
			return nil, fmt.Errorf("missing required field %q", name)
		}
		if err := validateCurrentValue(kind, v); err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
	}
	return object, nil
}
func currentIntegerValue(value any) (int64, bool) {
	n, ok := value.(json.Number)
	if !ok {
		return 0, false
	}
	v, err := n.Int64()
	return v, err == nil
}
func currentID(value any) bool {
	v, ok := value.(string)
	if !ok || len(v) == 0 || len(v) > 128 {
		return false
	}
	for i := 0; i < len(v); i++ {
		if v[i] < 32 || v[i] > 126 {
			return false
		}
	}
	return true
}
func currentEnum(value any, allowed ...string) bool {
	v, ok := value.(string)
	if !ok {
		return false
	}
	for _, a := range allowed {
		if v == a {
			return true
		}
	}
	return false
}
func currentTimestamp(value any) bool {
	v, ok := value.(string)
	if !ok {
		return false
	}
	_, err := time.Parse(time.RFC3339Nano, v)
	return err == nil && strings.HasSuffix(v, "Z")
}
func currentLowerHex(value string, size int) bool {
	if len(value) != size || strings.ToLower(value) != value {
		return false
	}
	_, err := hex.DecodeString(value)
	return err == nil
}
func currentPath(value any) bool {
	v, ok := value.(string)
	if !ok || v == "" || len(v) > 512 || !utf8.ValidString(v) || strings.ContainsAny(v, "\\\x00") || strings.HasPrefix(v, "/") || path.Clean(v) != v {
		return false
	}
	for _, part := range strings.Split(v, "/") {
		if part == "." || part == ".." || part == "" {
			return false
		}
	}
	return true
}
