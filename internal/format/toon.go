package format

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

func writeTOON(writer io.Writer, value any) error {
	structure, err := genericStructure(value)
	if err != nil {
		return err
	}
	if err := writeTOONValue(writer, structure); err != nil {
		return err
	}
	if _, err := io.WriteString(writer, "\n"); err != nil {
		return fmt.Errorf("write toon newline: %w", err)
	}
	return nil
}

func writeTOONValue(writer io.Writer, value any) error {
	switch typed := value.(type) {
	case nil:
		return writeTOONToken(writer, "null")
	case bool:
		if typed {
			return writeTOONToken(writer, "true")
		}
		return writeTOONToken(writer, "false")
	case string:
		return writeTOONString(writer, typed)
	case json.Number:
		return writeTOONToken(writer, typed.String())
	case []any:
		return writeTOONArray(writer, typed)
	case map[string]any:
		return writeTOONObject(writer, typed)
	default:
		return fmt.Errorf("unsupported toon value type %T", value)
	}
}

func writeTOONObject(writer io.Writer, value map[string]any) error {
	if err := writeTOONToken(writer, "{"); err != nil {
		return err
	}

	keys := make([]string, 0, len(value))
	for key := range value {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for index, key := range keys {
		if index > 0 {
			if err := writeTOONToken(writer, ","); err != nil {
				return err
			}
		}
		if err := writeTOONKey(writer, key); err != nil {
			return err
		}
		if err := writeTOONToken(writer, ":"); err != nil {
			return err
		}
		if err := writeTOONValue(writer, value[key]); err != nil {
			return err
		}
	}

	return writeTOONToken(writer, "}")
}

func writeTOONArray(writer io.Writer, value []any) error {
	if err := writeTOONToken(writer, "["); err != nil {
		return err
	}
	for index, item := range value {
		if index > 0 {
			if err := writeTOONToken(writer, ","); err != nil {
				return err
			}
		}
		if err := writeTOONValue(writer, item); err != nil {
			return err
		}
	}
	return writeTOONToken(writer, "]")
}

func writeTOONKey(writer io.Writer, value string) error {
	if isBareTOONToken(value) {
		return writeTOONToken(writer, value)
	}
	return writeQuotedTOONString(writer, value)
}

func writeTOONString(writer io.Writer, value string) error {
	if isBareTOONToken(value) {
		return writeTOONToken(writer, value)
	}
	return writeQuotedTOONString(writer, value)
}

func writeQuotedTOONString(writer io.Writer, value string) error {
	encoded, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshal toon string: %w", err)
	}
	return writeTOONBytes(writer, encoded)
}

func isBareTOONToken(value string) bool {
	if value == "" {
		return false
	}
	for _, r := range value {
		switch {
		case r >= 'a' && r <= 'z':
		case r >= 'A' && r <= 'Z':
		case r >= '0' && r <= '9':
		case r == '_' || r == '-' || r == '.':
		default:
			return false
		}
	}
	return true
}

func writeTOONToken(writer io.Writer, value string) error {
	return writeTOONBytes(writer, []byte(value))
}

func writeTOONBytes(writer io.Writer, value []byte) error {
	if _, err := writer.Write(value); err != nil {
		return fmt.Errorf("write toon output: %w", err)
	}
	return nil
}
