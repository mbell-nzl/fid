package journal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"unicode/utf8"
)

func validateCurrentJSONText(raw string) error {
	if !utf8.ValidString(raw) {
		return fmt.Errorf("invalid UTF-8 JSON")
	}
	decoder := json.NewDecoder(bytes.NewBufferString(raw))
	decoder.UseNumber()
	if err := scanCurrentJSONValue(decoder); err != nil {
		return err
	}
	if _, err := decoder.Token(); !errors.Is(err, io.EOF) {
		if err == nil {
			return fmt.Errorf("multiple JSON values")
		}
		return err
	}
	return nil
}

func decodeCurrentJSON(raw string) (any, error) {
	if err := validateCurrentJSONText(raw); err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(bytes.NewBufferString(raw))
	decoder.UseNumber()
	var value any
	if err := decoder.Decode(&value); err != nil {
		return nil, err
	}
	return value, nil
}

func scanCurrentJSONValue(decoder *json.Decoder) error {
	token, err := decoder.Token()
	if err != nil {
		return err
	}
	delimiter, ok := token.(json.Delim)
	if !ok {
		return nil
	}
	switch delimiter {
	case '{':
		keys := map[string]bool{}
		for decoder.More() {
			keyValue, err := decoder.Token()
			if err != nil {
				return err
			}
			key, ok := keyValue.(string)
			if !ok {
				return fmt.Errorf("object key is not text")
			}
			if keys[key] {
				return fmt.Errorf("duplicate key %q", key)
			}
			keys[key] = true
			if err := scanCurrentJSONValue(decoder); err != nil {
				return err
			}
		}
		_, err := decoder.Token()
		return err
	case '[':
		for decoder.More() {
			if err := scanCurrentJSONValue(decoder); err != nil {
				return err
			}
		}
		_, err := decoder.Token()
		return err
	default:
		return fmt.Errorf("unexpected JSON delimiter %q", delimiter)
	}
}
